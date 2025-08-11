package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"order-service/internal/domain"
)

var ErrIdempotencyKeyExists = errors.New("idempotency key already exists")

type OrderUseCase struct {
	logger     *slog.Logger
	repository OrderRepository
}

func NewOrderUseCase(logger *slog.Logger, repository OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		logger:     logger,
		repository: repository,
	}
}

func (c *OrderUseCase) CreateOrder(ctx context.Context, params domain.OrderParams) error {
	if err := c.checkIdempotency(params.OrderUID); err != nil {
		if errors.Is(err, ErrIdempotencyKeyExists) {
			c.logger.Error("idempotency key already exists", "key", params.OrderUID)
			return err
		}
		return err
	}
	order, err := domain.NewOrder(params)
	if err != nil {
		c.logger.Error("failed to create order", "error:", err, "params:", params)
		return err
	}
	if err = c.repository.SaveOrder(ctx, order); err != nil {
		return err
	}
	return nil

}

func (c *OrderUseCase) GetOrder(ctx context.Context, uuid string) (*domain.Order, error) {
	if uuid == "" {
		c.logger.Error("missing required field", "uuid", uuid)
		return nil, domain.ErrInvalidState
	}

	order, err := c.repository.GetOrderByUid(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (c *OrderUseCase) checkIdempotency(uuid string) error {
	if uuid == "" {
		c.logger.Error("missing required field", "uuid", uuid)
		return domain.ErrInvalidState
	}

	_, err := c.repository.GetIdempotencyKey(context.Background(), uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		c.logger.Error("failed to check idempotency key", "error", err)
		return err
	}
	return ErrIdempotencyKeyExists

}
