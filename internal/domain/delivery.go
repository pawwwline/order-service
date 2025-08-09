package domain

import "fmt"

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type DeliveryParams struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

func NewDelivery(params DeliveryParams) (*Delivery, error) {
	if err := validateDeliveryNotEmpty(params); err != nil {
		return nil, err
	}
	return &Delivery{
		Name:    params.Name,
		Phone:   params.Phone,
		Zip:     params.Zip,
		City:    params.City,
		Address: params.Address,
		Region:  params.Region,
		Email:   params.Email,
	}, nil
}

func validateDeliveryNotEmpty(params DeliveryParams) error {
	if params.Phone == "" {
		return fmt.Errorf("phone is empty: %w", ErrInvalidState)
	}
	if params.Zip == "" {
		return fmt.Errorf("zip is empty: %w", ErrInvalidState)
	}
	if params.City == "" {
		return fmt.Errorf("city is empty: %w", ErrInvalidState)
	}
	if params.Address == "" {
		return fmt.Errorf("address is empty: %w", ErrInvalidState)
	}
	if params.Region == "" {
		return fmt.Errorf("region is empty: %w", ErrInvalidState)
	}

	return nil
}
