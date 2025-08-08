package usecase

import (
	"context"
	"log/slog"
	"order-service/internal/domain"
)

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
		return nil, domain.ErrMissingRequiredField
	}

	order, err := c.repository.GetOrderByUid(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return order, nil
}
