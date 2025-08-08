package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"order-service/internal/domain"
	logger2 "order-service/internal/lib/logger"
	"testing"
	"time"
)

type MockOrderRepository struct {
	validOrder domain.Order
}

func (m *MockOrderRepository) SaveOrder(ctx context.Context, order *domain.Order) error {
	return nil
}

func (m *MockOrderRepository) GetOrderByUid(ctx context.Context, orderUID string) (*domain.Order, error) {
	m.validOrder = domain.Order{
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
		Items: []domain.Item{
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
	return &m.validOrder, nil

}

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
	Items: []domain.Item{
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

func TestCreateOrderInvalidParams(t *testing.T) {
	repo := &MockOrderRepository{}
	logger, err := logger2.InitLogger("test")
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}
	orderUseCase := NewOrderUseCase(logger, repo)
	ctx := context.Background()
	err = orderUseCase.CreateOrder(ctx, testParamsInvalid)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCreateOrderValidParams(t *testing.T) {
	repo := &MockOrderRepository{}
	ctx := context.Background()
	logger, err := logger2.InitLogger("test")
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}
	orderUseCase := NewOrderUseCase(logger, repo)
	err = orderUseCase.CreateOrder(ctx, testParamsValid)
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}

}

func TestGetOrderInvalidParams(t *testing.T) {
	repo := &MockOrderRepository{}
	logger, err := logger2.InitLogger("test")
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}
	orderUseCase := NewOrderUseCase(logger, repo)
	ctx := context.Background()

	uuid := ""
	order, err := orderUseCase.GetOrder(ctx, uuid)
	if err == nil {
		t.Fatalf("expected err got nil")
	}
	if order != nil {
		t.Fatalf("expected order to be nil, got %v", order)
	}

}

func TestGetOrderValidParams(t *testing.T) {
	repo := &MockOrderRepository{}
	logger, err := logger2.InitLogger("test")
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}
	orderUseCase := NewOrderUseCase(logger, repo)
	uuid := testParamsValid.OrderUID
	t.Logf("uuid: %v", uuid)
	ctx := context.Background()
	t.Logf("expected: %v", expectedOrder)
	order, err := orderUseCase.GetOrder(ctx, uuid)
	if err != nil {
		t.Fatalf("expected err to be nil , got %v", err)
	}
	if !assert.Equal(t, expectedOrder, order) {
		t.Fatalf("expected order to be equal to %v, got %v", repo.validOrder, order)
	}
}
