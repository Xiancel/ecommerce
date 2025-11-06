package http

import (
	"encoding/json"
	"net/http"

	cartSrv "github.com/Xiancel/ecommerce/internal/service/cart"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CartHandler struct {
	CartSrv cartSrv.CartService
}

func NewCartHandler(srv cartSrv.CartService) *CartHandler {
	return &CartHandler{CartSrv: srv}
}

func (h *CartHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/cart", h.ListItems)
		r.Post("/cart/items", h.AddItem)
		r.Put("/cart/items/{id}", h.UpdateItem)
		r.Delete("/cart/items/{id}", h.DeleteItem)
		r.Delete("/cart", h.ClearCart)
	})
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}
	var req cartSrv.AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.CartSrv.AddItem(r.Context(), userID, req)
	if err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req cartSrv.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.CartSrv.UpdateItem(r.Context(), userID, id, req)
	if err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.CartSrv.DeleteItem(r.Context(), userID, id); err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "item deleted",
	})
}
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}
	if err := h.CartSrv.ClearItem(r.Context(), userID); err != nil {
		handlerCartError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "cart clear succesfully",
	})
}
func (h *CartHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}
	items, err := h.CartSrv.ListItem(r.Context(), userID)
	if err != nil {
		handlerCartError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func handlerCartError(w http.ResponseWriter, err error) {
	switch err {
	case cartSrv.ErrItemNotFound:
		respondError(w, http.StatusNotFound, err.Error())
	case cartSrv.ErrInvalidQuantity,
		cartSrv.ErrProductNotAvailable,
		cartSrv.ErrInvalidProductID:
		respondError(w, http.StatusBadRequest, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
