package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/nirupam52/expenseTrack/internal/repository"
	"github.com/nirupam52/expenseTrack/internal/response"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Name == "" {
		response.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if input.Email == "" {
		response.WriteError(w, http.StatusBadRequest, "email is required")
		return
	}
	if input.Password == "" {
		response.WriteError(w, http.StatusBadRequest, "password is required")
		return
	}

	ctx := context.Background()

	_, err := h.repo.GetUserByEmail(ctx, input.Email)
	if err == nil {
		response.WriteError(w, http.StatusConflict, "email already registered")
		return
	}
	if err != nil && (err != sql.ErrNoRows && err.Error() != "User not found") {
		response.WriteError(w, http.StatusInternalServerError, "failed to check email")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to process password")
		return
	}

	err = h.repo.CreateUser(ctx, input.Name, input.Email, string(hashedPassword))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	user, err := h.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	userResp := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	response.WriteSuccess(w, http.StatusCreated, userResp)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	ctx := context.Background()
	user, err := h.repo.GetUserById(ctx, id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	userResp := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	response.WriteSuccess(w, http.StatusOK, userResp)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := context.Background()
	users, err := h.repo.ListUsers(ctx)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	var userResponses []UserResponse
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}

	response.WriteList(w, userResponses)
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("POST /users/register", h.Register)
	mux.HandleFunc("GET /users/{id}", h.GetByID)
	mux.HandleFunc("GET /users", h.List)
	return nil
}
