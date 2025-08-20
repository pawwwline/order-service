package handler

import (
	"context"
	"order-service/internal/domain"
)

type OrderCreatorUseCase interface {
	CreateOrder(ctx context.Context, params domain.OrderParams) error
}
