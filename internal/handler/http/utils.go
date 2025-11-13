package http

import (
	"encoding/json"
	"net/http"
)

// допоміжні функції

// повертає помилку
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResponse{
		Error:  msg,
		Status: status,
	})
}

// структура помилки
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// вивід інформації в JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
