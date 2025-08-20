package cache

import (
	"order-service/internal/domain"

	lru "github.com/hashicorp/golang-lru"
)

type LRUCache struct {
	cache *lru.Cache
}

func NewLRUCache(size int) (*LRUCache, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &LRUCache{
		cache: cache,
	}, nil
}

func (l *LRUCache) Get(key string) (*domain.Order, bool) {
	v, ok := l.cache.Get(key)
	if !ok {
		return nil, false
	}
	order, ok := v.(*domain.Order)
	return order, ok
}

func (l *LRUCache) Set(order *domain.Order) {
	l.cache.Add(order.OrderUID, order)
}
