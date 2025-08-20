package usecase

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/domain"
	"order-service/internal/infra/repo"
)

var ErrIdempotencyKeyExists = errors.New("idempotency key already exists")

type OrderUseCase struct {
	repository OrderRepository
	cache      Cache
}

func NewOrderUseCase(repository OrderRepository, cache Cache) *OrderUseCase {
	return &OrderUseCase{
		repository: repository,
		cache:      cache,
	}
}

func (c *OrderUseCase) CreateOrder(ctx context.Context, params domain.OrderParams) error {
	if err := c.checkIdempotency(params.OrderUID); err != nil {
		return err
	}

	order, err := domain.NewOrder(params)
	if err != nil {
		return err
	}
	if err = c.repository.SaveOrder(ctx, order); err != nil {
		return err
	}

	c.cache.Set(order)

	return nil

}

func (c *OrderUseCase) GetOrder(ctx context.Context, uid string) (*domain.Order, error) {
	if uid == "" {
		return nil, fmt.Errorf("uid is empty: %w", domain.ErrInvalidState)
	}

	order, ok := c.cache.Get(uid)
	if ok {
		return order, nil
	}

	order, err := c.repository.GetOrderByUid(ctx, uid)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, fmt.Errorf("order_uid %s: %w", uid, err)
		}
		return nil, err
	}
	c.cache.Set(order)

	return order, nil
}

func (c *OrderUseCase) LoadOrdersCache(ctx context.Context, limit int) error {
	orders, err := c.repository.GetLastOrders(ctx, limit)
	if err != nil {
		return err
	}
	for _, order := range orders {
		c.cache.Set(order)
	}

	return nil

}

func (c *OrderUseCase) checkIdempotency(uid string) error {
	if uid == "" {
		return fmt.Errorf("uid is empty %w", domain.ErrInvalidState)
	}

	exists, err := c.repository.CheckIdempotencyKey(context.Background(), uid)
	if err != nil {
		return fmt.Errorf("idempotnecy check failed %w", err)
	}
	if exists {
		return ErrIdempotencyKeyExists
	}

	return nil
}
