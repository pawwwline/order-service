package usecase

import (
	"context"
	"order-service/internal/domain"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *domain.Order) error
	GetOrderByUid(ctx context.Context, orderUID string) (*domain.Order, error)
	CheckIdempotencyKey(ctx context.Context, key string) (bool, error)
	GetLastOrders(ctx context.Context, limit int) ([]*domain.Order, error)
}
