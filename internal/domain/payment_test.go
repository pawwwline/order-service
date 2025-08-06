package domain

import (
	"errors"
	"testing"
)

func TestNewPaymentInvalidParams(t *testing.T) {
	params := PaymentParams{
		Transaction:  "",
		RequestID:    "",
		Currency:     "",
		Provider:     "",
		Amount:       -200,
		PaymentDt:    0,
		Bank:         "",
		DeliveryCost: 0,
		GoodsTotal:   -100,
		CustomFee:    0,
	}

	payment, err := NewPayment(params)
	if payment != nil {
		t.Fatalf("expected nil, got %v", payment)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrValueBelowZero) && !errors.Is(err, ErrMissingRequiredField) {
		t.Fatalf("expected error to be ErrValueBelowZero or be ErrMissingRequiredField, got %v", err)
	}

}

func TestNewPaymentValidParams(t *testing.T) {
	params := PaymentParams{
		Transaction:  "12345",
		RequestID:    "12345",
		Currency:     "USD",
		Provider:     "PayPal",
		Amount:       100,
		PaymentDt:    1636739200,
		Bank:         "PayPal",
		DeliveryCost: 10,
	}
	payment, err := NewPayment(params)
	if payment == nil {
		t.Fatalf("expected not nil, got %v", payment)
	}
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}
}
