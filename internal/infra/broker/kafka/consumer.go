package kafka

import (
	"context"
	"log/slog"
	"order-service/internal/config"
	"order-service/internal/infra/broker/handler"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	broker       string
	handler      Handler
	retryHandler RetryHandler
	logger       *slog.Logger
	orderReader  *kafka.Reader
	retryReader  *kafka.Reader
	retryWriter  *kafka.Writer
	DLQWriter    *kafka.Writer
}

func NewKafkaConsumer(cfg *config.KafkaConfig, handler Handler, retry RetryHandler, logger *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		broker: cfg.Broker,
		orderReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Broker},
			GroupID: cfg.OrderTopicCfg.GroupID,
			Topic:   cfg.OrderTopicCfg.KafkaTopic,
		}),
		retryReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Broker},
			GroupID: cfg.RetryTopicCfg.GroupID,
			Topic:   cfg.RetryTopicCfg.KafkaTopic,
		}),
		retryWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{cfg.Broker},
			Topic:   cfg.RetryTopicCfg.KafkaTopic,
		}),
		DLQWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{cfg.Broker},
			Topic:   cfg.DLQTopicCfg.KafkaTopic,
		}),
	}
}

func (kc *KafkaConsumer) ReadOrderMsg(ctx context.Context) {
	msg, err := kc.orderReader.ReadMessage(ctx)
	if err != nil {
		kc.logger.Error("failed to read message from Kafka", "error", err)
		return
	}
	if len(msg.Value) == 0 {
		kc.logger.Error("no data to process")
		return
	}

	res := kc.handler.ProcessOrderMessage(ctx, msg.Value)
	switch res {
	case handler.Success:
		if err := kc.orderReader.CommitMessages(ctx, msg); err != nil {
			kc.logger.Error("failed to commit message", "error", err)
		}
		kc.logger.Info("message processed", "order_uid", msg.Value)
		return
	case handler.Retry:
		if err := kc.WriteRetryTopic(ctx, msg); err != nil {
			kc.logger.Error("failed to write to retry topic", "error", err)
		}
		return
	case handler.DLQ:
		if err := kc.WriteDLQTopic(ctx, msg); err != nil {
			kc.logger.Error("failed to write to dead letter queue", "error", err)
		}
		return
	}

}

func (kc *KafkaConsumer) WriteRetryTopic(ctx context.Context, msg kafka.Message) error {
	if err := kc.retryWriter.WriteMessages(ctx, msg); err != nil {
		kc.logger.Error("error writing to retry topic", "err", err)
	}
	return nil
}

func (kc *KafkaConsumer) ReadRetryMsg(ctx context.Context) {
	msg, err := kc.retryReader.ReadMessage(ctx)
	if err != nil {
		kc.logger.Error("failed to read message from Kafka", "error", err)
		return
	}
	if len(msg.Value) == 0 {
		kc.logger.Error("no data to process")
	}

	res := kc.retryHandler.RetryWrapper(ctx, func() handler.Result {
		return kc.handler.ProcessOrderMessage(ctx, msg.Value)
	})

	switch res {
	case handler.Success:
		if err := kc.retryReader.CommitMessages(ctx, msg); err != nil {
			kc.logger.Error("failed to commit messages", "err", err)
		}
		return
	case handler.DLQ:
		if err := kc.WriteDLQTopic(ctx, msg); err != nil {
			kc.logger.Error("failed to write to dlq topuc", "err", err)
		}
		return
	}
	kc.logger.Info("message processed", "order_uid", msg.Value)

}

func (kc *KafkaConsumer) WriteDLQTopic(ctx context.Context, msg kafka.Message) error {
	return kc.DLQWriter.WriteMessages(ctx, msg)
}

func (kc *KafkaConsumer) ShutDown() error {
	if err := kc.orderReader.Close(); err != nil {
		return err
	}
	if err := kc.retryReader.Close(); err != nil {
		return err
	}
	if err := kc.retryWriter.Close(); err != nil {
		return err
	}
	if err := kc.DLQWriter.Close(); err != nil {
		return err
	}
	return nil
}
