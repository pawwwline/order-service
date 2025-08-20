package domain

import "fmt"

type Item struct {
	ChrtID      int
	TrackNumber string
	Price       int
	Rid         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int
	Brand       string
	Status      int
}

type ItemParams struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func NewItem(p ItemParams) (*Item, error) {
	if err := validateItemParams(p); err != nil {
		return nil, err
	}
	return &Item{
		ChrtID:      p.ChrtID,
		TrackNumber: p.TrackNumber,
		Price:       p.Price,
		Brand:       p.Brand,
		Rid:         p.Rid,
		Name:        p.Name,
		Sale:        p.Sale,
		Size:        p.Size,
		TotalPrice:  p.TotalPrice,
		NmID:        p.NmID,
		Status:      p.Status,
	}, nil
}

func NewItemList(p []ItemParams) ([]*Item, error) {
	err := validateItemLength(p)
	if err != nil {
		return nil, err
	}
	var items []*Item
	for _, item := range p {
		i, err := NewItem(item)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func validateItemLength(p []ItemParams) error {
	if len(p) == 0 {
		return fmt.Errorf("no items : %w", ErrInvalidState)
	}
	return nil
}

func validateItemParams(p ItemParams) error {
	if p.Price < 0 {
		return fmt.Errorf("price is below zero: %w", ErrInvalidState)
	}
	if p.TotalPrice < 0 {
		return fmt.Errorf("price is below zero: %w", ErrInvalidState)
	}
	return nil
}
