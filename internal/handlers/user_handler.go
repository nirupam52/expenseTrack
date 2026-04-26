package handlers

import (
	"encoding/json"
	"errors"
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

	_, err := h.repo.GetUserByEmail(r.Context(), input.Email)
	if err == nil {
		response.WriteError(w, http.StatusConflict, "email already registered")
		return
	}
	if !errors.Is(err, repository.ErrNotFound) {
		response.WriteError(w, http.StatusInternalServerError, "failed to check email")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to process password")
		return
	}

	id, err := h.repo.CreateUser(r.Context(), input.Name, input.Email, string(hashedPassword))
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	response.WriteSuccess(w, http.StatusCreated, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	response.WriteSuccess(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers(r.Context())
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
