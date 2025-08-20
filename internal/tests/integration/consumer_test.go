//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"order-service/internal/config"
	"order-service/internal/infra/broker"
	"order-service/internal/infra/broker/handler"
	"order-service/internal/infra/broker/kafka"
	"order-service/internal/infra/broker/retry"
	"order-service/internal/infra/cache"
	pg "order-service/internal/infra/repo/postgres"
	"order-service/internal/lib/logger"
	"order-service/internal/usecase"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	kafkago "github.com/segmentio/kafka-go"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestBrokerIntegration(t *testing.T) {
	ctx := context.Background()

	dbDSN := "postgres://test:test@localhost:5432/test?sslmode=disable"

	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	cwd, _ := os.Getwd()
	migrationsDir := filepath.Join(cwd, "migrations")
	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	addr := "localhost:9092"
	cfg := createKafkaCfg(addr)
	require.NoError(t, waitForKafka(cfg.Broker, 60*time.Second))
	createTopics(t, *cfg)

	logger, err := logger.InitLogger("test")
	if err != nil {
		t.Fatalf("error creating logger %v", err)
	}
	if logger == nil {
		t.Fatalf("logger is nil")
	}
	if err != nil {
		t.Fatalf("error creating usecase: %v", err)
	}

	msgs := []kafkago.Message{
		{Key: []byte("1"), Value: []byte(`{"order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`)}, {Key: []byte("2"), Value: []byte(`{
   "order_uid": "b563feb7b2b84bas6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`)}, {Key: []byte("3"), Value: []byte(`{
   "order_uid": "b563feb7b2b84bas6testt",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`)},
		{Key: []byte("invalid1"), Value: []byte(`{"order_uid:"","amount":-100}`)},
		{Key: []byte("invalid1"), Value: []byte(`{"order_uid:"wwe","amount":-100}`)},
	}
	repo := pg.NewPostgresDB(db)
	if repo == nil {
		t.Fatalf("repo is nil")
	}
	cache, err := cache.NewLRUCache(1000)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	uc := usecase.NewOrderUseCase(repo, cache)
	if uc == nil {
		t.Fatalf("uc is nil")
	}
	processor := handler.NewMessageProcessor(uc, logger)
	if processor == nil {
		t.Fatalf("processor id nil")
	}
	retry := retry.NewRetry(*cfg)
	if retry == nil {
		t.Fatalf("retry is nil")
	}

	if logger == nil {
		t.Fatalf("logger is nil")
	}

	consumer := kafka.NewKafkaConsumer(cfg, processor, retry, logger)
	if consumer == nil {
		t.Fatalf("conumer id nil")
	}

	broker := broker.NewBroker(consumer)
	go broker.Run(ctx)

	if err := produceTestMessages(cfg, msgs); err != nil {
		t.Fatalf("error producing messages: %v", err)
	}

	waitForProcessing(t, db, 3, 60*time.Second)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	require.NoError(t, err)
	rows, err := db.Query("SELECT order_uid FROM orders")
	if err != nil {
		require.NoError(t, err)
	}
	defer rows.Close()

	var uids []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			require.NoError(t, err)
		}
		uids = append(uids, uid)
	}

	if err := rows.Err(); err != nil {
		require.NoError(t, err)
	}
	t.Logf("orders loaded: %v", strings.Join(uids, ", "))
	require.NoError(t, err)
	require.Equal(t, 3, count)
}

func buildUseCase(db *sql.DB) (*usecase.OrderUseCase, error) {
	repo := pg.NewPostgresDB(db)
	cache, err := cache.NewLRUCache(1000)
	if err != nil {
		return nil, err
	}
	uc := usecase.NewOrderUseCase(repo, cache)
	return uc, nil
}

func createKafkaCfg(addr string) *config.KafkaConfig {
	return &config.KafkaConfig{
		Broker:             addr,
		RetryMaxAttempts:   10,
		BackoffDurationMin: 10,
		BackoffDurationMax: 10,
		OrderTopicCfg: config.OrderTopicConfig{
			KafkaTopic: "orders",
			GroupID:    "test-orders",
		},
		RetryTopicCfg: config.RetryTopicConfig{
			KafkaTopic: "retry-order",
			GroupID:    "test-retry",
		},
		DLQTopicCfg: config.DLQTopicConfig{
			KafkaTopic: "dlq",
			GroupID:    "test-dlq",
		},
	}
}

func produceTestMessages(cfg *config.KafkaConfig, messages []kafkago.Message) error {
	w := &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Broker),
		Topic:        cfg.OrderTopicCfg.KafkaTopic,
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
		WriteTimeout: 10 * time.Second,
	}

	ctx := context.Background()
	return w.WriteMessages(ctx, messages...)
}

func createTopics(t *testing.T, cfg config.KafkaConfig) {
	conn, err := kafkago.Dial("tcp", cfg.Broker)
	require.NoError(t, err)
	defer conn.Close()

	topics := []kafkago.TopicConfig{
		{Topic: cfg.OrderTopicCfg.KafkaTopic, NumPartitions: 1, ReplicationFactor: 1},
		{Topic: cfg.RetryTopicCfg.KafkaTopic, NumPartitions: 1, ReplicationFactor: 1},
		{Topic: cfg.DLQTopicCfg.KafkaTopic, NumPartitions: 1, ReplicationFactor: 1},
	}

	err = conn.CreateTopics(topics...)
	require.NoError(t, err)
}

func readDLQMessages(t *testing.T, cfg *config.KafkaConfig, expectedCount int, timeoutSec int) []kafkago.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	conn, err := kafkago.DialLeader(ctx, "tcp", cfg.Broker, cfg.DLQTopicCfg.KafkaTopic, 0)
	require.NoError(t, err)
	defer conn.Close()

	var messages []kafkago.Message
	for len(messages) < expectedCount {
		msg, err := conn.ReadMessage(1e6)
		require.NoError(t, err)
		messages = append(messages, msg)
	}

	return messages
}

func waitForProcessing(t *testing.T, db *sql.DB, expectedCount int, timeout time.Duration) {
	ctx := context.Background()
	timer := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timer:
			t.Fatalf("timeout waiting for processing")
		case <-ticker.C:
			var count int
			err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
			require.NoError(t, err)
			if count >= expectedCount {
				return
			}
		}
	}
}

func waitForKafka(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		conn, err := kafkago.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for Kafka at %s: %w", addr, err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
