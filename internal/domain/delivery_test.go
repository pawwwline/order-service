package domain

import (
	"errors"
	"testing"
)

func TestNewDeliveryInvalidParams(t *testing.T) {
	params := DeliveryParams{
		Name:    "",
		Phone:   "",
		Zip:     "",
		City:    "",
		Address: "",
		Region:  "",
		Email:   "",
	}

	delivery, err := NewDelivery(params)
	if delivery != nil {
		t.Fatalf("expected nil, got %v", delivery)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected error to be ErrInvalidState, got %v", err)
	}
}

func TestNewDeliveryValidParams(t *testing.T) {
	params := DeliveryParams{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "12345",
		City:    "test",
		Address: "test",
		Region:  "test",
		Email:   "test@gmail.com",
	}
	delivery, err := NewDelivery(params)
	if delivery == nil {
		t.Fatalf("expected not nil, got %v", delivery)
	}
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}

}
