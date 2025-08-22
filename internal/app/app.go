package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"order-service/internal/config"
	server "order-service/internal/controller/http"
	"order-service/internal/infra/broker"
	"order-service/internal/infra/broker/handler"
	"order-service/internal/infra/broker/kafka"
	"order-service/internal/infra/broker/retry"
	"order-service/internal/infra/cache"
	"order-service/internal/infra/repo/postgres"
	"order-service/internal/lib/logger"
	"order-service/internal/usecase"
)

type App struct {
	httpServer *server.Server
	broker     *broker.Broker
	db         *postgres.PostgresDB
	usecase    *usecase.OrderUseCase
	logger     *slog.Logger
}

func BuildApp(cfg *config.Config) (*App, error) {
	logger, err := buildLogger(cfg.Env)
	if err != nil {
		return nil, err
	}
	db, err := buildDB(&cfg.DB)
	if err != nil {
		return nil, err
	}

	cache, err := buildCache(&cfg.Cache)
	if err != nil {
		return nil, err
	}
	usecase := buildUseCase(db, cache)
	broker := buildBroker(&cfg.Kafka, usecase, logger)
	httpServer := buildHTTP(&cfg.HTTP, usecase, logger)

	return &App{
		httpServer: httpServer,
		broker:     broker,
		db:         db,
		usecase:    usecase,
		logger:     logger,
	}, nil
}

func buildLogger(env string) (*slog.Logger, error) {
	return logger.InitLogger(env)
}

func buildDB(cfg *config.DBConfig) (*postgres.PostgresDB, error) {
	dbConn, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}
	return postgres.NewPostgresDB(dbConn), nil
}

func buildCache(cfg *config.CacheConfig) (*cache.LRUCache, error) {
	return cache.NewLRUCache(cfg.Limit)
}

func buildUseCase(db *postgres.PostgresDB, cache *cache.LRUCache) *usecase.OrderUseCase {
	return usecase.NewOrderUseCase(db, cache)
}

func buildBroker(cfg *config.KafkaConfig, uc *usecase.OrderUseCase, logger *slog.Logger) *broker.Broker {
	processor := handler.NewMessageProcessor(uc, logger)
	retry := retry.NewRetry(*cfg)
	consumer := kafka.NewKafkaConsumer(cfg, processor, retry, logger)
	return broker.NewBroker(consumer, logger)
}

func buildHTTP(cfg *config.HTTPConfig, uc *usecase.OrderUseCase, logger *slog.Logger) *server.Server {
	return server.NewServer(cfg, uc, logger)
}

func (a *App) Run(ctx context.Context) error {
	if err := a.usecase.LoadOrdersCache(ctx, 1000); err != nil {
		return err
	}
	a.logger.Info("orders cache loaded")

	go func() {
		a.httpServer.Run()
	}()
	go func() {
		a.broker.Run(ctx)
	}()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var errList []error

	if err := a.httpServer.Shutdown(ctx); err != nil {
		errList = append(errList, err)
	}
	a.logger.Info("http server shutdown")

	if err := a.broker.Shutdown(); err != nil {
		errList = append(errList, err)
	}
	a.logger.Info("broker shutdown")

	if err := a.db.Close(); err != nil {
		errList = append(errList, err)
	}
	a.logger.Info("db shutdown")

	if len(errList) > 0 {
		return fmt.Errorf("shutdown errors: %v", errList)
	}

	return nil
}
