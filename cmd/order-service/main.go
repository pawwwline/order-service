package main

import (
	"context"
	"fmt"
	"log"
	"order-service/internal/app"
	"order-service/internal/config"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Printf("config err: %v", err)
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := app.BuildApp(cfg)
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Printf("app run error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()

	if err := app.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown errors: %v", err)
	}
}
