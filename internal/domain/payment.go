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
		return fmt.Errorf("transaction %w", ErrMissingRequiredField)
	}
	if p.Currency == "" {
		return fmt.Errorf("currency %w", ErrMissingRequiredField)
	}

	if p.Amount < 0 {
		return fmt.Errorf("amount %w", ErrValueBelowZero)
	}
	if p.DeliveryCost < 0 {
		return fmt.Errorf("deliveryCost %w", ErrValueBelowZero)
	}
	if p.GoodsTotal < 0 {
		return fmt.Errorf("goodsTotal %w", ErrValueBelowZero)
	}
	if p.CustomFee < 0 {
		return fmt.Errorf("customFee %w", ErrValueBelowZero)
	}

	return nil
}
