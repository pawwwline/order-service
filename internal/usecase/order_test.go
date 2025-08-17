package usecase

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedOrder = &domain.Order{
	OrderUID:    "b563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: &domain.Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: &domain.Payment{
		Transaction:  "b563feb7b2b84b6test",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
	},
	Items: []*domain.Item{
		{
			ChrtID:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmID:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
	},
	Locale:          "en",
	CustomerID:      "test",
	DeliveryService: "meest",
	Shardkey:        "9",
	SmID:            99,
	DateCreated:     time.Date(2021, 11, 26, 6, 22, 0, 0, time.UTC),
	OofShard:        "9",
}

type MockOrderRepo struct {
	validOrder domain.Order
	saveErr    error
	getErr     error
	idempErr   error
	called     bool
}

func (m *MockOrderRepo) SaveOrder(ctx context.Context, order *domain.Order) error {
	return m.saveErr
}

func (m *MockOrderRepo) GetOrderByUid(ctx context.Context, uid string) (*domain.Order, error) {
	m.called = true
	if m.getErr != nil {
		return nil, m.getErr
	}
	return &m.validOrder, nil
}

func (m *MockOrderRepo) CheckIdempotencyKey(ctx context.Context, key string) (bool, error) {
	m.called = true
	if m.idempErr != nil {
		return false, m.idempErr
	}
	if key == m.validOrder.OrderUID {
		return false, nil
	}
	return true, nil
}

func (m *MockOrderRepo) GetLastOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0, limit)
	for i := 0; i < limit; i++ {
		orders = append(orders, &domain.Order{
			OrderUID:    fmt.Sprintf("order-%d", i),
			DateCreated: time.Now().Add(-time.Duration(i) * time.Minute),
		})
	}
	return orders, nil
}

type MockCache struct {
	cache  map[string]*domain.Order
	called bool
}

func NewMockCache() *MockCache { return &MockCache{cache: make(map[string]*domain.Order)} }

func (mc *MockCache) Get(uid string) (*domain.Order, bool) {
	order, ok := mc.cache[uid]
	mc.called = true
	return order, ok
}

func (mc *MockCache) Set(order *domain.Order) {
	mc.called = true
	mc.cache[order.OrderUID] = order
}

func setupUseCase(repo *MockOrderRepo) (*OrderUseCase, *MockCache) {
	cache := NewMockCache()
	return NewOrderUseCase(repo, cache), cache
}

func TestOrderUseCase_CreateOrder(t *testing.T) {
	repo := &MockOrderRepo{validOrder: *expectedOrder}
	uc, cache := setupUseCase(repo)

	validParams := domain.OrderParams{
		OrderUID:    expectedOrder.OrderUID,
		TrackNumber: expectedOrder.TrackNumber,
		Entry:       expectedOrder.Entry,
		Delivery: domain.DeliveryParams{
			Name:    expectedOrder.Delivery.Name,
			Phone:   expectedOrder.Delivery.Phone,
			Zip:     expectedOrder.Delivery.Zip,
			City:    expectedOrder.Delivery.City,
			Address: expectedOrder.Delivery.Address,
			Region:  expectedOrder.Delivery.Region,
			Email:   expectedOrder.Delivery.Email,
		},
		Payment: domain.PaymentParams{
			Transaction:  expectedOrder.Payment.Transaction,
			Currency:     expectedOrder.Payment.Currency,
			Provider:     expectedOrder.Payment.Provider,
			Amount:       expectedOrder.Payment.Amount,
			PaymentDt:    expectedOrder.Payment.PaymentDt,
			Bank:         expectedOrder.Payment.Bank,
			DeliveryCost: expectedOrder.Payment.DeliveryCost,
			GoodsTotal:   expectedOrder.Payment.GoodsTotal,
		},
		Items: []domain.ItemParams{
			{
				ChrtID:      expectedOrder.Items[0].ChrtID,
				TrackNumber: expectedOrder.Items[0].TrackNumber,
				Price:       expectedOrder.Items[0].Price,
				Rid:         expectedOrder.Items[0].Rid,
				Name:        expectedOrder.Items[0].Name,
				Sale:        expectedOrder.Items[0].Sale,
				Size:        expectedOrder.Items[0].Size,
				TotalPrice:  expectedOrder.Items[0].TotalPrice,
				NmID:        expectedOrder.Items[0].NmID,
				Brand:       expectedOrder.Items[0].Brand,
				Status:      expectedOrder.Items[0].Status,
			},
		},
		Locale:          expectedOrder.Locale,
		CustomerID:      expectedOrder.CustomerID,
		DeliveryService: expectedOrder.DeliveryService,
		Shardkey:        expectedOrder.Shardkey,
		SmID:            expectedOrder.SmID,
		DateCreated:     expectedOrder.DateCreated,
		OofShard:        expectedOrder.OofShard,
	}

	t.Run("valid params success", func(t *testing.T) {
		err := uc.CreateOrder(context.Background(), validParams)
		assert.NoError(t, err)
		cached, ok := cache.Get(expectedOrder.OrderUID)
		assert.True(t, ok)
		assert.Equal(t, expectedOrder.OrderUID, cached.OrderUID)
	})
}

func TestOrderUseCase_GetOrder(t *testing.T) {
	repo := &MockOrderRepo{validOrder: *expectedOrder}
	uc, cache := setupUseCase(repo)
	ctx := context.Background()

	t.Run("invalid uid", func(t *testing.T) {
		order, err := uc.GetOrder(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("cache hit", func(t *testing.T) {
		cache.Set(expectedOrder)
		repo.called = false
		order, err := uc.GetOrder(ctx, expectedOrder.OrderUID)
		assert.NoError(t, err)
		assert.False(t, repo.called)
		assert.Equal(t, expectedOrder.OrderUID, order.OrderUID)
	})

	t.Run("cache miss -> repo success", func(t *testing.T) {
		cache = NewMockCache()
		uc.cache = cache
		repo.called = false

		order, err := uc.GetOrder(ctx, expectedOrder.OrderUID)
		assert.NoError(t, err)
		assert.True(t, repo.called)

		cached, ok := cache.Get(expectedOrder.OrderUID)
		assert.True(t, ok)
		assert.Equal(t, expectedOrder.OrderUID, cached.OrderUID)
		assert.Equal(t, expectedOrder.OrderUID, order.OrderUID)
	})

	t.Run("cache miss -> repo error", func(t *testing.T) {
		cache = NewMockCache()
		repo = &MockOrderRepo{getErr: errors.New("repo error")}
		uc = NewOrderUseCase(repo, cache)

		order, err := uc.GetOrder(ctx, expectedOrder.OrderUID)
		assert.Error(t, err)
		assert.Nil(t, order)
	})
}

func TestOrderUseCase_LoadOrdersCache(t *testing.T) {
	repo := &MockOrderRepo{}
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	limit := 5
	err := uc.LoadOrdersCache(context.Background(), limit)
	assert.NoError(t, err)

	orders, _ := repo.GetLastOrders(context.Background(), limit)
	for _, order := range orders {
		cached, ok := cache.Get(order.OrderUID)
		assert.True(t, ok)
		assert.Equal(t, order.OrderUID, cached.OrderUID)
	}
}
