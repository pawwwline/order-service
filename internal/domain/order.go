package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidState = errors.New("invalid domain state")
)

type Order struct {
	Id                int
	OrderUID          string
	TrackNumber       string
	Entry             string
	Delivery          *Delivery
	Payment           *Payment
	Items             []*Item
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SmID              int
	DateCreated       time.Time
	OofShard          string
}

type OrderParams struct {
	OrderUID          string         `json:"order_uid"`
	TrackNumber       string         `json:"track_number"`
	Entry             string         `json:"entry"`
	Delivery          DeliveryParams `json:"delivery"`
	Payment           PaymentParams  `json:"payment"`
	Items             []ItemParams   `json:"items"`
	Locale            string         `json:"locale"`
	InternalSignature string         `json:"internal_signature"`
	CustomerID        string         `json:"customer_id"`
	DeliveryService   string         `json:"delivery_service"`
	Shardkey          string         `json:"shardkey"`
	SmID              int            `json:"sm_id"`
	DateCreated       time.Time      `json:"date_created"`
	OofShard          string         `json:"oof_shard"`
}

func NewOrder(p OrderParams) (*Order, error) {
	err := validateOrder(p)
	if err != nil {
		return nil, err
	}
	delivery, err := NewDelivery(p.Delivery)
	if err != nil {
		return nil, err
	}
	payment, err := NewPayment(p.Payment)
	if err != nil {
		return nil, err
	}
	items, err := NewItemList(p.Items)
	if err != nil {
		return nil, err
	}
	return &Order{
		OrderUID:          p.OrderUID,
		TrackNumber:       p.TrackNumber,
		Entry:             p.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            p.Locale,
		InternalSignature: p.InternalSignature,
		CustomerID:        p.CustomerID,
		DeliveryService:   p.DeliveryService,
		Shardkey:          p.Shardkey,
		SmID:              p.SmID,
		DateCreated:       p.DateCreated,
		OofShard:          p.OofShard,
	}, nil
}

func validateOrder(p OrderParams) error {
	if p.OrderUID == "" {
		return fmt.Errorf("orderId is missing: %w", ErrInvalidState)
	}
	if p.TrackNumber == "" {
		return fmt.Errorf("tracknumber is missing: %w", ErrInvalidState)
	}
	if p.CustomerID == "" {
		return fmt.Errorf("customerID is missing: %w", ErrInvalidState)
	}
	return nil
}
