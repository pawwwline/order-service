package dto

import "time"

type OrderResponse struct {
	OrderUID        string           `json:"order_uid"`
	TrackNumber     string           `json:"track_number"`
	Delivery        DeliveryResponse `json:"delivery"`
	Payment         PaymentResponse  `json:"payment"`
	Items           []ItemResponse   `json:"items"`
	CustomerID      string           `json:"customer_id"`
	DeliveryService string           `json:"delivery_service"`
	DateCreated     time.Time        `json:"date_created"`
}

type DeliveryResponse struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentResponse struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type ItemResponse struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
