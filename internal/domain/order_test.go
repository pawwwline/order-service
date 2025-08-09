package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewOrderInvalid(t *testing.T) {
	params := OrderParams{
		OrderUID:   "",
		Items:      nil,
		Locale:     "en",
		Delivery:   DeliveryParams{},
		Payment:    PaymentParams{},
		CustomerID: "",
		Shardkey:   "",
		SmID:       0,
	}
	order, err := NewOrder(params)
	if order != nil {
		t.Fatalf("expected nil, got %v", order)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrMissingRequiredField) && !errors.Is(err, ErrValueBelowZero) && !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected error to be ErrMissingRequiredField or be ErrValueBelowZero or be ErrInvalidState, got %v", err)
	}
}

func TestNewOrderValid(t *testing.T) {
	params := OrderParams{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: DeliveryParams{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: PaymentParams{
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
		Items: []ItemParams{
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

	order, err := NewOrder(params)
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}
	if order == nil {
		t.Fatalf("expected not nil, got %v", order)
	}

}

func TestNewOrderEmptyStructs(t *testing.T) {
	params := OrderParams{
		OrderUID:          "b563feb7b2b84b6test",
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Delivery:          DeliveryParams{},
		Payment:           PaymentParams{},
		Items:             []ItemParams{},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 0, 0, time.UTC),
		OofShard:          "9",
	}

	order, err := NewOrder(params)
	if order != nil {
		t.Fatalf("expected nil, got %v", order)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrMissingRequiredField) {
		t.Fatalf("expected error to be ErrMissingRequiredField, got %v", err)
	}
}
