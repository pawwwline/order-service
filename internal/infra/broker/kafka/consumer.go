package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"order-service/internal/config"
	"order-service/internal/domain"
	"order-service/internal/infra/broker/handler"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	broker       string
	handler      Handler
	retryHandler RetryHandler
	logger       *slog.Logger
	cfg          *config.KafkaConfig
	orderReader  *kafka.Reader
	retryReader  *kafka.Reader
	retryWriter  *kafka.Writer
	DLQWriter    *kafka.Writer
	ready        chan struct{}
}

var ErrNotInitialized = errors.New("not initialized")

func NewKafkaConsumer(cfg *config.KafkaConfig, handler Handler, retryHandler RetryHandler, logger *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		broker:       cfg.Broker,
		cfg:          cfg,
		handler:      handler,
		retryHandler: retryHandler,
		logger:       logger,
		ready:        make(chan struct{}),
	}
}

func (kc *KafkaConsumer) Init() error {
	kc.orderReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kc.broker},
		GroupID: kc.cfg.OrderTopicCfg.GroupID,
		Topic:   kc.cfg.OrderTopicCfg.KafkaTopic,
	})
	kc.retryReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kc.broker},
		GroupID: kc.cfg.RetryTopicCfg.GroupID,
		Topic:   kc.cfg.RetryTopicCfg.KafkaTopic,
	})
	kc.retryWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kc.broker},
		Topic:   kc.cfg.RetryTopicCfg.KafkaTopic,
	})
	kc.DLQWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kc.broker},
		Topic:   kc.cfg.DLQTopicCfg.KafkaTopic,
	})

	conn, err := kafka.Dial("tcp", kc.broker)
	if err != nil {
		return err
	}
	defer func(conn *kafka.Conn) {
		err := conn.Close()
		if err != nil {
			kc.logger.Error("failed to close connection", "error", err)
		}
	}(conn)

	topics := []struct {
		topic string
	}{
		{kc.cfg.OrderTopicCfg.KafkaTopic},
		{kc.cfg.RetryTopicCfg.KafkaTopic},
		{kc.cfg.DLQTopicCfg.KafkaTopic},
	}

	for _, t := range topics {
		_, err := kafka.DialLeader(context.Background(), "tcp", kc.broker, t.topic, 0)
		if err != nil {
			return err
		}
	}

	close(kc.ready)

	return nil
}

func (kc *KafkaConsumer) Ready() <-chan struct{} {
	return kc.ready
}
func (kc *KafkaConsumer) ReadOrderMsg(ctx context.Context) error {
	if kc.orderReader == nil {
		kc.logger.Error("orderReader is not initialized")
		return ErrNotInitialized
	}
	msg, err := kc.orderReader.ReadMessage(ctx)
	if err != nil {
		kc.logger.Error("failed to read message from Kafka", "error", err)
		return err
	}
	if len(msg.Value) == 0 {
		kc.logger.Error("no data to process")
		return nil
	}

	res := kc.handler.ProcessOrderMessage(ctx, msg.Value)
	var order domain.OrderParams
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		kc.logger.Error("failed to unmarshal order message", "error", err)
		return err
	}
	switch res {
	case handler.Success:
		if err := kc.orderReader.CommitMessages(ctx, msg); err != nil {
			kc.logger.Error("failed to commit message", "error", err, "uid", order.OrderUID)
			return err
		}
		kc.logger.Info("message processed", "uid", order.OrderUID)
		return nil
	case handler.Retry:
		if kc.retryWriter == nil {
			kc.logger.Error("retry writer is not initialized")
			return ErrNotInitialized
		}
		if err := kc.WriteRetryTopic(ctx, kafka.Message{
			Key:   msg.Key,
			Value: msg.Value,
		}); err != nil {
			kc.logger.Error("failed to write to retryHandler topic", "error", err, "uid", order.OrderUID)
			return err
		}
		kc.logger.Debug("written to retry")
		return nil
	case handler.DLQ:
		if kc.DLQWriter == nil {
			kc.logger.Error("orderReader is not initialized")
			return ErrNotInitialized
		}
		if err := kc.WriteDLQTopic(ctx, kafka.Message{
			Key:   msg.Key,
			Value: msg.Value,
		}); err != nil {

			kc.logger.Error("failed to write to dead letter queue", "error", err, "uid", order.OrderUID)
		}
		kc.logger.Debug("written to dlq", "uid", order.OrderUID)

		return nil
	}

	return nil

}

func (kc *KafkaConsumer) WriteRetryTopic(ctx context.Context, msg kafka.Message) error {
	if kc.retryWriter == nil {
		kc.logger.Error("retry writer is not initialized")
		return ErrNotInitialized
	}
	return kc.retryWriter.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
	})
}

func (kc *KafkaConsumer) ReadRetryMsg(ctx context.Context) error {
	msg, err := kc.retryReader.ReadMessage(ctx)
	if err != nil {
		kc.logger.Error("failed to read message from Kafka", "error", err)
		return err
	}
	if len(msg.Value) == 0 {
		kc.logger.Error("no data to process")
	}

	res := kc.retryHandler.RetryWrapper(ctx, func() handler.Result {
		return kc.handler.ProcessOrderMessage(ctx, msg.Value)
	})
	var order domain.OrderParams
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		kc.logger.Error("failed to unmarshal order message", "error", err)
		return err
	}

	switch res {
	case handler.Success:
		if err := kc.retryReader.CommitMessages(ctx, msg); err != nil {
			kc.logger.Error("failed to commit messages", "err", err, "uid", order.OrderUID)
		}
		kc.logger.Debug("succesfully created")
		return err
	case handler.DLQ:
		if err := kc.WriteDLQTopic(ctx, kafka.Message{
			Key:   msg.Key,
			Value: msg.Value,
		}); err != nil {
			kc.logger.Error("failed to write to dlq topuc", "err", err)
			return err
		}
		kc.logger.Debug("written to dlq", "uid", order.OrderUID)
		return nil

	}
	kc.logger.Info("message processed", "order_uid", order.OrderUID)
	return nil
}

func (kc *KafkaConsumer) WriteDLQTopic(ctx context.Context, msg kafka.Message) error {
	if kc.DLQWriter == nil {
		kc.logger.Error("dlq writer is not initialized")
		return ErrNotInitialized
	}
	return kc.DLQWriter.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
	})
}

func (kc *KafkaConsumer) ShutDown() error {
	var errs []error

	if kc.orderReader != nil {
		if err := kc.orderReader.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if kc.retryReader != nil {
		if err := kc.retryReader.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if kc.retryWriter != nil {
		if err := kc.retryWriter.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if kc.DLQWriter != nil {
		if err := kc.DLQWriter.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
