package domain

import "fmt"

type Payment struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type PaymentParams struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

func NewPayment(params PaymentParams) (*Payment, error) {
	err := validatePaymentParams(params)
	if err != nil {
		return nil, err
	}
	return &Payment{
		Transaction:  params.Transaction,
		RequestID:    params.RequestID,
		Currency:     params.Currency,
		Provider:     params.Provider,
		Amount:       params.Amount,
		PaymentDt:    params.PaymentDt,
		Bank:         params.Bank,
		DeliveryCost: params.DeliveryCost,
		GoodsTotal:   params.GoodsTotal,
		CustomFee:    params.CustomFee,
	}, nil
}

func validatePaymentParams(p PaymentParams) error {
	if p.Transaction == "" {
		return fmt.Errorf("transaction %w", ErrInvalidState)
	}
	if p.Currency == "" {
		return fmt.Errorf("currency %w", ErrInvalidState)
	}

	if p.Amount < 0 {
		return fmt.Errorf("amount is below zero: %w", ErrInvalidState)
	}
	if p.DeliveryCost < 0 {
		return fmt.Errorf("deliveryCost is below zero: %w", ErrInvalidState)
	}
	if p.GoodsTotal < 0 {
		return fmt.Errorf("goodsTotal is below zero: %w", ErrInvalidState)
	}
	if p.CustomFee < 0 {
		return fmt.Errorf("customFee is below zero: %w", ErrInvalidState)
	}

	return nil
}
