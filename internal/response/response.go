package response

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type ListResponse[T any] struct {
	Success bool `json:"success"`
	Data    []T  `json:"data"`
	Meta    Meta `json:"meta"`
}

type Meta struct {
	Count int `json:"count"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func WriteSuccess[T any](w http.ResponseWriter, status int, data T) error {
	return WriteJSON(w, status, Response[T]{
		Success: true,
		Data:    data,
	})
}

func WriteList[T any](w http.ResponseWriter, data []T) error {
	return WriteJSON(w, http.StatusOK, ListResponse[T]{
		Success: true,
		Data:    data,
		Meta: Meta{
			Count: len(data),
		},
	})
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, ErrorResponse{
		Success: false,
		Error: message,
	})
}
