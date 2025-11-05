package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	userSrv "github.com/Xiancel/ecommerce/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserSrv userSrv.UserService
}

func NewUserHandler(srv userSrv.UserService) *UserHandler {
	return &UserHandler{UserSrv: srv}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/user", h.ListUser)
	})
	r.Group(func(r chi.Router) {
		r.Get("/admin/user/{id}", h.GetUser)
		r.Delete("/admin/user{id}", h.DeleteUser)
		r.Put("/admin/user{id}", h.UpdateUser)
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "InvalidID")
		return
	}

	user, err := h.UserSrv.GetUser(r.Context(), id)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.UserSrv.ListUser(r.Context(), filter)
	if err != nil {
		handlerServiceUserError(w, err)
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	updUser, err := h.UserSrv.UpdateUser(r.Context(), id, req)
	if err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := h.UserSrv.DeleteUser(r.Context(), id); err != nil {
		handlerServiceUserError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted succesfully",
	})
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
