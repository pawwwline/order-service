//go:build integration

package integration

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"order-service/internal/config"
	srv "order-service/internal/controller/http"
	"order-service/internal/infra/cache"
	"order-service/internal/infra/repo/postgres"
	logger2 "order-service/internal/lib/logger"
	"order-service/internal/usecase"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedJSON = `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"name":"Mascaras","sale":30,"size":"0","total_price":317,"brand":"Vivienne Sabo","status":202}],"customer_id":"test","delivery_service":"meest","date_created":"2021-11-26T06:22:19Z"}`

func TestGetOrderHTTP(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t, ctx)
	defer teardownTestDB(t, db)
	uc := buildUCase(t, db)
	cfg := setupCfg()
	logger, err := logger2.InitLogger("test")
	require.NoError(t, err)
	server := srv.NewServer(cfg, uc, logger)
	order, err := uc.GetOrder(ctx, "b563feb7b2b84b6test")
	require.NoError(t, err)
	require.NotNil(t, order)
	ts := httptest.NewServer(server.Routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/order/b563feb7b2b84b6test")
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("body: %s", string(body))
	assert.JSONEq(t, expectedJSON, string(body))

}

func setupTestDB(t *testing.T, ctx context.Context) *sql.DB {
	dbDSN := "postgres://test:test@localhost:5432/test?sslmode=disable"
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	cwd, err := os.Getwd()
	require.NoError(t, err)
	migrationsDir := filepath.Join(cwd, "migrations")
	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func teardownTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		t.Logf("failed to clean db: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Logf("failed to close db: %v", err)
	}
}

func buildUCase(t *testing.T, db *sql.DB) *usecase.OrderUseCase {
	repo := postgres.NewPostgresDB(db)
	lruCache, err := cache.NewLRUCache(1000)
	require.NoError(t, err)

	return usecase.NewOrderUseCase(repo, lruCache)
}

func setupCfg() *config.HTTPConfig {
	return &config.HTTPConfig{
		Host:         "0.0.0.0",
		Port:         "8080",
		ReadTimeout:  5,
		WriteTimeout: 10,
		IdleTimeout:  60,
	}
}
