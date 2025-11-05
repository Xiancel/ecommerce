package http

import (
	"encoding/json"
	"net/http"

	authSrv "github.com/Xiancel/ecommerce/internal/service/auth"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	AuthSrv authSrv.AuthService
}

func NewAuthHandler(srv authSrv.AuthService) *AuthHandler {
	return &AuthHandler{AuthSrv: srv}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/refresh", h.RefreshToken)
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authSrv.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	resp, err := h.AuthSrv.Register(r.Context(), req)
	if err != nil {
		handlerAuthError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authSrv.LoginRequset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	resp, err := h.AuthSrv.Login(r.Context(), req)
	if err != nil {
		handlerAuthError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req authSrv.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.AuthSrv.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		handlerAuthError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func handlerAuthError(w http.ResponseWriter, err error) {
	switch err {
	case authSrv.ErrInvalidCredentials:
		respondError(w, http.StatusUnauthorized, err.Error())
	case authSrv.ErrUserAlreadyExists:
		respondError(w, http.StatusConflict, err.Error())
	case authSrv.ErrInvalidToken, authSrv.ErrExpiredToken:
		respondError(w, http.StatusUnauthorized, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
