package cache

import (
	"order-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLRUCache_SetGet(t *testing.T) {
	l, err := NewLRUCache(2)
	require.NoError(t, err)
	if l == nil {
		t.Fatalf("l is nil")
	}

	order1 := &domain.Order{OrderUID: "1"}
	order2 := &domain.Order{OrderUID: "2"}

	l.Set(order1)
	l.Set(order2)

	got, ok := l.Get("1")
	assert.True(t, ok)
	assert.Equal(t, order1.OrderUID, got.OrderUID)

	got, ok = l.Get("2")
	assert.True(t, ok)
	assert.Equal(t, order2.OrderUID, got.OrderUID)

	got, ok = l.Get("3")
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestLRUCache_Eviction(t *testing.T) {
	l, err := NewLRUCache(2)
	require.NoError(t, err)

	order1 := &domain.Order{OrderUID: "1"}
	order2 := &domain.Order{OrderUID: "2"}
	order3 := &domain.Order{OrderUID: "3"}

	l.Set(order1)
	l.Set(order2)

	_, _ = l.Get("1")

	l.Set(order3)

	_, ok := l.Get("1")
	assert.True(t, ok)

	_, ok = l.Get("2")
	assert.False(t, ok)

	_, ok = l.Get("3")
	assert.True(t, ok)
}
