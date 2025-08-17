package usecase

import (
	"context"
	"database/sql"
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
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
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
	Locale:            "en",
	InternalSignature: "",
	CustomerID:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmID:              99,
	DateCreated:       time.Date(2021, 11, 26, 6, 22, 0, 0, time.UTC),
	OofShard:          "9",
}

var testParamsValid = domain.OrderParams{
	OrderUID:    "b563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: domain.DeliveryParams{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: domain.PaymentParams{
		Transaction:  "b563feb7b2b84b6test",
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Items: []domain.ItemParams{
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
	Locale:            "en",
	InternalSignature: "",
	CustomerID:        "test",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmID:              99,
	DateCreated:       time.Date(2021, 11, 26, 6, 22, 0, 0, time.UTC),
	OofShard:          "9",
}

var testParamsInvalid = domain.OrderParams{
	OrderUID:    "",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: domain.DeliveryParams{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: domain.PaymentParams{
		Transaction:  "b563feb7b2b84b6test",
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Items: []domain.ItemParams{
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
	Locale:            "en",
	InternalSignature: "",
	CustomerID:        "",
	DeliveryService:   "meest",
	Shardkey:          "9",
	SmID:              99,
	DateCreated:       time.Date(2021, 11, 26, 6, 22, 0, 0, time.UTC),
	OofShard:          "9",
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

func (m *MockOrderRepo) GetOrderByUid(ctx context.Context, orderUID string) (*domain.Order, error) {
	m.called = true
	if m.getErr != nil {
		return nil, m.getErr
	}
	return &m.validOrder, nil
}

func (m *MockOrderRepo) GetIdempotencyKey(ctx context.Context, key string) (string, error) {
	if m.idempErr != nil {
		return key, nil
	}
	return "", sql.ErrNoRows
}

func (m *MockOrderRepo) GetLastOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0, limit)
	for i := 0; i < limit; i++ {
		order := &domain.Order{OrderUID: fmt.Sprintf("order-%d", i), DateCreated: time.Now().Add(-time.Duration(i) * time.Minute)}
		orders = append(orders, order)
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

func setupUseCaseWithRepo(repo *MockOrderRepo) (*OrderUseCase, *MockCache) {
	cache := NewMockCache()
	return NewOrderUseCase(repo, cache), cache
}

func TestOrderUseCase_CreateOrder(t *testing.T) {
	validOrder := expectedOrder
	repo := &MockOrderRepo{validOrder: *validOrder}
	uc, cache := setupUseCaseWithRepo(repo)

	t.Run("invalid params", func(t *testing.T) {
		err := uc.CreateOrder(context.Background(), testParamsInvalid)
		assert.Error(t, err)
	})

	t.Run("valid params success", func(t *testing.T) {
		err := uc.CreateOrder(context.Background(), testParamsValid)
		assert.NoError(t, err)

		cached, ok := cache.Get(testParamsValid.OrderUID)
		assert.True(t, ok)
		assert.Equal(t, testParamsValid.OrderUID, cached.OrderUID)
	})

	t.Run("idempotency key exists", func(t *testing.T) {
		repo.idempErr = errors.New("key exists")
		err := uc.CreateOrder(context.Background(), testParamsValid)
		assert.ErrorIs(t, err, ErrIdempotencyKeyExists)
	})

	t.Run("repository save error", func(t *testing.T) {
		repo.saveErr = errors.New("repo error")
		err := uc.CreateOrder(context.Background(), testParamsValid)
		assert.Error(t, err)
	})
}

func TestOrderUseCase_GetOrder(t *testing.T) {
	validOrder := expectedOrder
	repo := &MockOrderRepo{validOrder: *validOrder}
	uc, cache := setupUseCaseWithRepo(repo)
	ctx := context.Background()
	if err := uc.CreateOrder(ctx, testParamsValid); err != nil {
		t.Fatalf("error creating order %v", err)
	}

	t.Run("invalid uid", func(t *testing.T) {
		order, err := uc.GetOrder(ctx, "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidState)
		assert.Nil(t, order)
	})

	t.Run("cache hit", func(t *testing.T) {
		cache.Set(validOrder)
		order, err := uc.GetOrder(ctx, validOrder.OrderUID)
		assert.NoError(t, err)
		assert.False(t, repo.called)
		assert.Equal(t, validOrder, order)
	})

	t.Run("cache miss -> repo success", func(t *testing.T) {
		cache = NewMockCache()
		uc.cache = cache
		order, err := uc.GetOrder(ctx, validOrder.OrderUID)
		assert.NoError(t, err)
		assert.Equal(t, validOrder.OrderUID, order.OrderUID)
		assert.True(t, repo.called)
		cached, ok := cache.Get(validOrder.OrderUID)
		assert.True(t, ok)
		assert.Equal(t, validOrder.OrderUID, cached.OrderUID)
	})

	t.Run("cache miss -> repo error", func(t *testing.T) {
		cache = NewMockCache()
		repo.getErr = errors.New("repo error")
		uc.cache = cache
		order, err := uc.GetOrder(ctx, validOrder.OrderUID)
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
