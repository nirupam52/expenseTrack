package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nirupam52/expenseTrack/internal/repository"
	"github.com/nirupam52/expenseTrack/internal/response"
)

type contextKey string

const contextKeyUserID contextKey = "user_id"

func NewAuthMiddleware(authRepo *repository.AuthRepository) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.WriteError(w, http.StatusUnauthorized, "authorization required")
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")

			userID, err := authRepo.GetSessionByToken(r.Context(), token)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					response.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
					return
				}
				response.WriteError(w, http.StatusInternalServerError, "failed to validate token")
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
			next(w, r.WithContext(ctx))
		}
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(contextKeyUserID).(int64)
	if !ok {
		return 0, fmt.Errorf("user id not in context")
	}
	return userID, nil
}
