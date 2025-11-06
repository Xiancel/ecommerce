package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	orderSrv "github.com/Xiancel/ecommerce/internal/service/order"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	OrderSrv orderSrv.OrderService
}

func NewOrderhandler(srv orderSrv.OrderService) *OrderHandler {
	return &OrderHandler{OrderSrv: srv}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/orders/{id}", h.GetOrder)
		r.Get("/orders", h.ListOrder)
		r.Post("/orders", h.CreateOrder)
		r.Put("/orders/{id}/cancel", h.CancelOrder)
	})
	r.Group(func(r chi.Router) {
		r.Get("/admin/orders", h.ListAllOrder)
		r.Put("/admin/orders/{id}/status", h.UpdateOrderStatus)
	})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	user, err := h.OrderSrv.GetOrder(r.Context(), id)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	var req orderSrv.CreateOrderRequset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.UserID = userID

	order, err := h.OrderSrv.CreateOrder(r.Context(), req)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}
	if err := h.OrderSrv.CancelOrder(r.Context(), id); err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Order canceled!",
	})
}

func (h *OrderHandler) ListOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	filter := orderSrv.OrderFilter{
		UserID: &userID,
		Limit:  20,
		Offset: 0,
	}

	filter.Status = r.URL.Query().Get("status")

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid Offset")
		}
		filter.Offset = offset
	}
	orders, err := h.OrderSrv.ListOrder(r.Context(), filter)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

//ADMIN

func (h *OrderHandler) ListAllOrder(w http.ResponseWriter, r *http.Request) {
	filter := orderSrv.OrderFilter{
		Limit:  20,
		Offset: 0,
	}
	filter.Status = r.URL.Query().Get("status")

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "Invalid limit")
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "Invalid Offset")
		}
		filter.Offset = offset
	}

	orders, err := h.OrderSrv.ListOrder(r.Context(), filter)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}
	var req orderSrv.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updOrders, err := h.OrderSrv.UpdateOrderStatus(r.Context(), id, req)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updOrders)
}

func handlerOrderError(w http.ResponseWriter, err error) {
	switch err {
	case orderSrv.ErrOrderNotFound:
		respondError(w, http.StatusNotFound, err.Error())

	case orderSrv.ErrOrderIDRequired,
		orderSrv.ErrUserIDRequired,
		orderSrv.ErrShippingAddrReq,
		orderSrv.ErrInvalidPayment,
		orderSrv.ErrProductIDRequired,
		orderSrv.ErrStatusRequired,
		orderSrv.ErrInvalidStatus,
		orderSrv.ErrOrderEmpty:
		respondError(w, http.StatusBadRequest, err.Error())

	case orderSrv.ErrOrderAlreadyCanceled,
		orderSrv.ErrCannotCancelDelivered:
		respondError(w, http.StatusConflict, err.Error())

	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
