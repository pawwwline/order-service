package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var ErrCfgInvalid = errors.New("invalid configuration")

type Config struct {
	Env   string `env:"APP_ENV"`
	Kafka KafkaConfig
	DB    DBConfig
}

type DBConfig struct {
	Host         string `env:"DB_PORT"`
	Port         int    `env:"DB_PORT"`
	User         string `env:"DB_USER"`
	Password     string `env:"DB_PASSWORD"`
	Name         string `env:"DB_NAME"`
	SSLMode      string `env:"DB_SSL"`
	MaxOpenConns int    `env:"DB_MAX_CONNS"`
	MaxIdleConns int    `env:"DB_IDLE_CONNS"`
}

type OrderTopicConfig struct {
	KafkaTopic string `env:"KAFKA_ORDER_TOPIC"`
	GroupID    string `env:"KAFKA_ORDER_GROUP_ID"`
}

type DLQTopicConfig struct {
	KafkaTopic string `env:"KAFKA_DLQ_TOPIC"`
	GroupID    string `env:"KAFKA_DLQ_GROUP_ID"`
}

type RetryTopicConfig struct {
	KafkaTopic string `env:"KAFKA_RETRY_TOPIC"`
	GroupID    string `env:"KAFKA_RETRY_GROUP_ID"`
}

type KafkaConfig struct {
	OrderTopicCfg      OrderTopicConfig
	DLQTopicCfg        DLQTopicConfig
	RetryTopicCfg      RetryTopicConfig
	Broker             string `env:"KAFKA_BROKER"`
	RetryMaxAttempts   int    `env:"KAFKA_RETRY_MAX"`
	BackoffDurationMin int    `env:"KAFKA_BACKOFF_MAX"` // in milliseconds
	BackoffDurationMax int    `env:"KAFKA_BACKOFF_MIN"` // in milliseconds
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg, err := mapStructs()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func mapStructs() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
