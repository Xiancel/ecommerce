package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	orderSrv "github.com/Xiancel/ecommerce/internal/service/order"
	productSrv "github.com/Xiancel/ecommerce/internal/service/product"
	userSrv "github.com/Xiancel/ecommerce/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AdminHandler struct {
	productSrv productSrv.ProductService
	orderSrv   orderSrv.OrderService
	userSrv    userSrv.UserService
}

func NewAdminHandler(productSrv productSrv.ProductService, orderSrv orderSrv.OrderService, userSrv userSrv.UserService) *AdminHandler {
	return &AdminHandler{productSrv: productSrv,
		orderSrv: orderSrv,
		userSrv:  userSrv}
}

func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		//Product
		r.Post("/products", h.CreateProduct)
		r.Put("/products/{id}", h.UpdateProduct)
		//r.Delete("products/{id}", h.DeleteProduct)

		//Order
		r.Get("/orders", h.ListAllOrder)
		r.Put("/orders/{id}/status", h.UpdateOrderStatus)

		//User
		r.Get("/users/{id}", h.GetUser)
		r.Get("/users", h.ListUsers)
		r.Delete("/users/{id}", h.DeleteUser)
		r.Put("/users/{id}", h.UpdateUser)
	})
}

// Products
func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productSrv.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createProd, err := h.productSrv.CreateProduct(r.Context(), req)
	if err != nil {
		handlerServiceProductError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, createProd)
}

func (h *AdminHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req productSrv.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updProd, err := h.productSrv.UpdateProduct(r.Context(), id, req)
	if err != nil {
		handlerServiceProductError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updProd)
}

// Orders
func (h *AdminHandler) ListAllOrder(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.orderSrv.ListOrder(r.Context(), filter)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, orders)
}

func (h *AdminHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
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
	updOrders, err := h.orderSrv.UpdateOrderStatus(r.Context(), id, req)
	if err != nil {
		handlerOrderError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updOrders)
}

// Users
func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	user, err := h.userSrv.GetUser(r.Context(), id)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	filter := userSrv.UserFilter{
		Limit:  20,
		Offset: 0,
	}
	filter.Search = r.URL.Query().Get("search")

	if role := r.URL.Query().Get("role"); role != "" {
		filter.Role = &role
	}

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

	response, err := h.userSrv.ListUser(r.Context(), filter)
	if err != nil {
		handlerServiceUserError(w, err)
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := h.userSrv.DeleteUser(r.Context(), id); err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted succesfully",
	})
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var req userSrv.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updUser, err := h.userSrv.UpdateUser(r.Context(), id, req)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updUser)
}
