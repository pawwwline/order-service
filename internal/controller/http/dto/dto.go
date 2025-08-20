package dto

import "time"

type OrderResponse struct {
	OrderUID        string           `json:"order_uid" example:"b563feb7b2b84b6test"`
	TrackNumber     string           `json:"track_number" example:"WBILMTESTTRACK"`
	Delivery        DeliveryResponse `json:"delivery"`
	Payment         PaymentResponse  `json:"payment"`
	Items           []ItemResponse   `json:"items"`
	CustomerID      string           `json:"customer_id" example:"test"`
	DeliveryService string           `json:"delivery_service" example:"meest"`
	DateCreated     time.Time        `json:"date_created" example:"2021-11-26T06:22:19Z"`
}

type DeliveryResponse struct {
	Name    string `json:"name" example:"Test Testov"`
	Phone   string `json:"phone" example:"+9720000000"`
	Zip     string `json:"zip" example:"2639809"`
	City    string `json:"city" example:"Kiryat Mozkin"`
	Address string `json:"address" example:"Ploshad Mira 15"`
	Region  string `json:"region" example:"Kraiot"`
	Email   string `json:"email" example:"test@gmail.com"`
}

type PaymentResponse struct {
	Transaction  string `json:"transaction" example:"b563feb7b2b84b6test"`
	RequestID    string `json:"request_id" example:""`
	Currency     string `json:"currency" example:"USD"`
	Provider     string `json:"provider" example:"wbpay"`
	Amount       int    `json:"amount" example:"1817"`
	Bank         string `json:"bank" example:"alpha"`
	DeliveryCost int    `json:"delivery_cost" example:"1500"`
	GoodsTotal   int    `json:"goods_total" example:"317"`
	CustomFee    int    `json:"custom_fee" example:"0"`
}

type ItemResponse struct {
	ChrtID      int    `json:"chrt_id" example:"9934930"`
	TrackNumber string `json:"track_number" example:"WBILMTESTTRACK"`
	Price       int    `json:"price" example:"453"`
	Name        string `json:"name" example:"Mascaras"`
	Sale        int    `json:"sale" example:"30"`
	Size        string `json:"size" example:"0"`
	TotalPrice  int    `json:"total_price" example:"317"`
	Brand       string `json:"brand" example:"Vivienne Sabo"`
	Status      int    `json:"status" example:"202"`
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"order not found"`
}
