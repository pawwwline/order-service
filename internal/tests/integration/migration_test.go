package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMigrationWithDB(t *testing.T) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}()

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %s", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, host, port.Port(), dbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to connect to db: %s", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close db: %v", err)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %s", err)
	}

	cwd, _ := os.Getwd()

	migrationsDir := filepath.Join(cwd, "migrations")

	t.Run("apply migrations", func(t *testing.T) {
		if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
			t.Fatalf("failed to run migrations: %s", err)
		}
	})

	t.Run("idempotent", func(t *testing.T) {
		if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
			t.Fatalf("migration failed: %s", err)
		}

	})

	t.Run("check schema", func(t *testing.T) {
		tables := []string{"orders", "order_items", "payments", "deliveries"}
		for _, table := range tables {
			var exists bool
			err := db.QueryRow(
				"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name=$1)",
				table,
			).Scan(&exists)
			if err != nil || !exists {
				t.Fatalf("table '%s' not created", table)
			}
		}
	})

	t.Run("rollback", func(t *testing.T) {
		if err := goose.DownContext(ctx, db, migrationsDir); err != nil {
			t.Fatalf("failed to rollback migrations: %s", err)
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='users')").Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check table after rollback: %s", err)
		}
		if exists {
			t.Fatalf("table 'users' still exists after rollback")
		}
	})
}
