package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"order-service/internal/controller/http/dto"
	"order-service/internal/domain"
	"order-service/internal/infra/repo"
	"order-service/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type HTTPHandler struct {
	service *usecase.OrderUseCase
}

func NewHTTPHandler(service *usecase.OrderUseCase) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/order/{uid}", h.GetOrderHandler)
}

func (h *HTTPHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uuid")
	if uid == "" {
		http.Error(w, "uuid required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), uid)
	orderDTO := orderToResponse(order)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidState):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, repo.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderDTO)
	}
}

func orderToResponse(order *domain.Order) dto.OrderResponse {
	items := make([]dto.ItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.ItemResponse{
			Name:       item.Name,
			Price:      item.Price,
			Sale:       item.Sale,
			Size:       item.Size,
			TotalPrice: item.TotalPrice,
			Brand:      item.Brand,
			Status:     item.Status,
		}
	}

	var delivery dto.DeliveryResponse
	if order.Delivery != nil {
		delivery = dto.DeliveryResponse{
			Name:    order.Delivery.Name,
			Phone:   order.Delivery.Phone,
			Zip:     order.Delivery.Zip,
			City:    order.Delivery.City,
			Address: order.Delivery.Address,
			Region:  order.Delivery.Region,
			Email:   order.Delivery.Email,
		}
	}

	var payment dto.PaymentResponse
	if order.Payment != nil {
		payment = dto.PaymentResponse{
			Transaction:  order.Payment.Transaction,
			RequestID:    order.Payment.RequestID,
			Currency:     order.Payment.Currency,
			Provider:     order.Payment.Provider,
			Amount:       order.Payment.Amount,
			Bank:         order.Payment.Bank,
			DeliveryCost: order.Payment.DeliveryCost,
			GoodsTotal:   order.Payment.GoodsTotal,
			CustomFee:    order.Payment.CustomFee,
		}
	}

	return dto.OrderResponse{
		OrderUID:        order.OrderUID,
		TrackNumber:     order.TrackNumber,
		DeliveryService: order.DeliveryService,
		CustomerID:      order.CustomerID,
		DateCreated:     order.DateCreated,
		Delivery:        delivery,
		Payment:         payment,
		Items:           items,
	}
}
