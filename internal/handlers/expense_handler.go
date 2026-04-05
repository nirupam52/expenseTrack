package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nirupam52/expenseTrack/internal/repository"
	"github.com/nirupam52/expenseTrack/internal/response"
)

type ExpenseHandler struct {
	repo *repository.ExpenseRepository
}

func NewExpenseHandler(repo *repository.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{repo: repo}
}

type CreateExpenseInput struct {
	PaidBy      int64   `json:"paid_by"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

type UpdateExpenseInput struct {
	Amount      *float64 `json:"amount,omitempty"`
	Date        *string  `json:"date,omitempty"`
	Description *string  `json:"description,omitempty"`
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input CreateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.PaidBy <= 0 {
		response.WriteError(w, http.StatusBadRequest, "paid_by is required")
		return
	}
	if input.Amount <= 0 {
		response.WriteError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if input.Date == "" {
		response.WriteError(w, http.StatusBadRequest, "date is required")
		return
	}
	if input.Description == "" {
		response.WriteError(w, http.StatusBadRequest, "description is required")
		return
	}

	ctx := context.Background()
	err := h.repo.CreateExpense(ctx, input.PaidBy, input.Description, input.Amount, input.Date)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create expense")
		return
	}

	response.WriteSuccess(w, http.StatusCreated, map[string]string{"message": "expense created"})
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	ctx := context.Background()
	expense, err := h.repo.GetExpenseByID(ctx, id)
	if err != nil {
		if err.Error() == "expense not found" {
			response.WriteError(w, http.StatusNotFound, "expense not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to get expense")
		return
	}

	response.WriteSuccess(w, http.StatusOK, expense)
}

func (h *ExpenseHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		response.WriteError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	ctx := context.Background()
	expenses, err := h.repo.ListExpensesByUser(ctx, userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to list expenses")
		return
	}

	response.WriteList(w, expenses)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	var input UpdateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := context.Background()
	existing, err := h.repo.GetExpenseByID(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "expense not found")
		return
	}

	description := existing.Description
	if input.Description != nil {
		description = *input.Description
	}

	amount := existing.Amount
	if input.Amount != nil {
		amount = *input.Amount
	}

	date := existing.Date
	if input.Date != nil {
		date = *input.Date
	}

	err = h.repo.UpdateExpense(ctx, id, description, amount, date)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to update expense")
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "expense updated"})
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	ctx := context.Background()
	err = h.repo.DeleteExpense(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "expense not found")
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "expense deleted"})
}

// to clean up the main file
func (h *ExpenseHandler) RegisterRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("POST /expenses", h.Create)
	mux.HandleFunc("GET /expenses/{id}", h.GetByID)
	mux.HandleFunc("GET /expenses", h.ListByUser)
	mux.HandleFunc("PUT /expenses/{id}", h.Update)
	mux.HandleFunc("DELETE /expenses/{id}", h.Delete)
	return nil
}
