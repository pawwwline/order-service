package usecase

import (
	"order-service/internal/domain"
)

type Cache interface {
	Set(order *domain.Order)
	Get(uuid string) (*domain.Order, bool)
}
