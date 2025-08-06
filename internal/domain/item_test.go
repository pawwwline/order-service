package domain

import (
	"errors"
	"testing"
)

func TestNewItemListInvalid(t *testing.T) {
	itemParams := []ItemParams{
		{
			ChrtID:      1,
			TrackNumber: "",
			Price:       -1000,
			Rid:         "12345",
			Name:        "test",
			Sale:        10,
			Size:        "",
			TotalPrice:  -100,
			NmID:        1,
		},
		{
			ChrtID:      1,
			TrackNumber: "12345",
			Price:       100,
			Rid:         "12345",
			Name:        "test",
			Sale:        10,
		},
	}
	items, err := NewItemList(itemParams)
	if items != nil {
		t.Fatalf("expected nil, got %v", items)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrValueBelowZero) && !errors.Is(err, ErrMissingRequiredField) {
		t.Fatalf("expected error to be ErrValueBelowZero or be ErrMissingRequiredField, got %v", err)
	}

}

func TestNewItemListValid(t *testing.T) {
	itemParams := []ItemParams{
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
	}
	items, err := NewItemList(itemParams)
	if items == nil {
		t.Fatalf("expected not nil, got %v", items)
	}
	if err != nil {
		t.Fatalf("expected err nil, got %v", err)
	}

}

func TestNewItemListEmpty(t *testing.T) {
	var itemParams []ItemParams
	items, err := NewItemList(itemParams)
	if items != nil {
		t.Fatalf("expected nil, got %v", items)
	}
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected error to be ErrInvalidState, got %v", err)
	}
}
