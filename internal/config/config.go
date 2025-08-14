package config

import (
	"errors"
)

var ErrCfgInvalid = errors.New("invalid configuration")

type Config struct {
	Kafka KafkaConfig
}

type OrderTopicConfig struct {
	KafkaTopic string
	GroupID    string
}

type DLQTopicConfig struct {
	KafkaTopic string
	GroupID    string
}

type RetryTopicConfig struct {
	KafkaTopic string
	GroupID    string
}

type KafkaConfig struct {
	OrderTopicCfg      OrderTopicConfig
	DLQTopicCfg        DLQTopicConfig
	RetryTopicCfg      RetryTopicConfig
	Broker             string
	RetryMaxAttempts   int
	BackoffDurationMin int // in milliseconds
	BackoffDurationMax int // in milliseconds
}
