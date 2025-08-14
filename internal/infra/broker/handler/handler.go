package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"order-service/internal/domain"
	"order-service/internal/usecase"
)

type MessageProcessor struct {
	useCase OrderCreatorUseCase
	logger  *slog.Logger
}

func NewMessageProcessor(useCase OrderCreatorUseCase, logger *slog.Logger) *MessageProcessor {
	return &MessageProcessor{
		useCase: useCase,
		logger:  logger,
	}
}

type Result int

const (
	Success Result = iota
	Retry
	DLQ
)

func (p *MessageProcessor) ProcessOrderMessage(ctx context.Context, data []byte) Result {
	var params domain.OrderParams
	if err := json.Unmarshal(data, &params); err != nil {
		p.logger.Error("failed to unmarshal order message", "error", err)
		return DLQ
	}
	if err := p.useCase.CreateOrder(ctx, params); err != nil {
		p.logger.Error("failed to create order", "error", err, "order_uid", params.OrderUID)
		if p.shouldRetryErr(err) {
			return Retry
		}
		return DLQ
	}
	p.logger.Info("order processed successfully", "order_uid", params.OrderUID)
	return Success
}

func (p *MessageProcessor) shouldRetryErr(err error) bool {
	return !errors.Is(err, domain.ErrInvalidState) && !errors.Is(err, usecase.ErrIdempotencyKeyExists)
}
