package usecase

import (
	"context"
	"order-service/internal/domain"
)

type Cache interface {
	Set(ctx context.Context, uuid string, order *domain.Order)
	Get(ctx context.Context, uuid string) (*domain.Order, bool)
}
