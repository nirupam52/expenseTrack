package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nirupam52/expenseTrack/internal/repository"
	"github.com/nirupam52/expenseTrack/internal/response"
)

type AuthHandler struct {
	authRepo *repository.AuthRepository
}

func NewAuthHandler(authRepo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{authRepo: authRepo}
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
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

	userID, passwordHash, err := h.authRepo.GetCredentialsByEmail(r.Context(), input.Email)
	if err != nil {
		// Return 401 for both not-found and DB errors to prevent user enumeration
		response.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := generateToken()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	if err := h.authRepo.CreateSession(r.Context(), userID, token, expiresAt); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	response.WriteSuccess(w, http.StatusOK, LoginResponse{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	err := h.authRepo.DeleteSession(r.Context(), token)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		response.WriteError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	response.WriteSuccess(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux, protect func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/logout", protect(h.Logout))
}
