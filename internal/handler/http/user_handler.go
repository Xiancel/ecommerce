package http

import (
	"encoding/json"
	"net/http"

	userSrv "github.com/Xiancel/ecommerce/internal/service/user"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserSrv userSrv.UserService
}

func NewUserHandler(srv userSrv.UserService) *UserHandler {
	return &UserHandler{UserSrv: srv}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/users", h.ListUser)
		r.Put("/users", h.UpdateUser)
	})

}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}
	user, err := h.UserSrv.GetUser(r.Context(), userID)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "User not authorized")
		return
	}
	var req userSrv.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updUser, err := h.UserSrv.UpdateUser(r.Context(), userID, req)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updUser)
}

func handlerServiceUserError(w http.ResponseWriter, err error) {
	switch err {
	case userSrv.ErrUserNotFound:
		respondError(w, http.StatusNotFound, err.Error())

	case userSrv.ErrUserIDRequired,
		userSrv.ErrInvalidEmail,
		userSrv.ErrInvalidRole,
		userSrv.ErrPasswordRequired,
		userSrv.ErrNoFields:
		respondError(w, http.StatusBadRequest, err.Error())

	case userSrv.ErrEmailAlreadyExists:
		respondError(w, http.StatusConflict, err.Error())

	default:
		respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
