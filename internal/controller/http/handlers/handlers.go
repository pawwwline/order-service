package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	_ "order-service/docs"
	"order-service/internal/controller/http/dto"
	"order-service/internal/domain"
	"order-service/internal/infra/repo"
	"order-service/internal/usecase"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type HTTPHandler struct {
	service *usecase.OrderUseCase
}

func NewHTTPHandler(service *usecase.OrderUseCase) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/api/v1/order/{uid}", h.GetOrderHandler)
}

func (h *HTTPHandler) RegisterStaticRoutes(r *chi.Mux) {
	fileServer := http.FileServer(http.Dir("public"))
	r.Handle("/api/v1/swagger/*", httpSwagger.WrapHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
}

// GetOrderHandler @Summary Get order
// @Description Get order by UID
// @Tags orders
// @Param uid path string true "Order UID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{uid} [get]
func (h *HTTPHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "uuid required", http.StatusBadRequest)

		return
	}

	order, err := h.service.GetOrder(r.Context(), uid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case errors.Is(err, domain.ErrInvalidState):
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(dto.ErrorResponse{Code: http.StatusBadRequest, Message: "invalid state"})
			if err != nil {
				return
			}
		case errors.Is(err, repo.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(dto.ErrorResponse{Code: http.StatusNotFound, Message: err.Error()})
			if err != nil {
				return
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(dto.ErrorResponse{Code: http.StatusInternalServerError, Message: "internal server error"})
			if err != nil {
				return
			}
		}

		return
	}

	orderDTO := orderToResponse(order)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orderDTO)
	if err != nil {

		return
	}
}

func orderToResponse(order *domain.Order) dto.OrderResponse {
	if order == nil {
		return dto.OrderResponse{} // или можно возвращать ошибку
	}
	items := make([]dto.ItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.ItemResponse{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Name:        item.Name,
			Price:       item.Price,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			Brand:       item.Brand,
			Status:      item.Status,
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
