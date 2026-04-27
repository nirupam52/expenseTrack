package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// GetCredentialsByEmail returns the user ID and bcrypt hash for login verification.
// Kept separate from GetUserByEmail so the hash never leaks into general user queries.
func (r *AuthRepository) GetCredentialsByEmail(ctx context.Context, email string) (userID int64, passwordHash string, err error) {
	query := `SELECT id, password_hash FROM users WHERE email = ?`
	err = r.db.QueryRowContext(ctx, query, email).Scan(&userID, &passwordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", ErrNotFound
	}
	if err != nil {
		return 0, "", fmt.Errorf("get credentials by email: %w", err)
	}
	return userID, passwordHash, nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetSessionByToken(ctx context.Context, token string) (int64, error) {
	query := `SELECT user_id, expires_at FROM sessions WHERE token = ?`
	var userID int64
	var expiresAt string
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("get session by token: %w", err)
	}

	exp, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return 0, fmt.Errorf("parse session expiry: %w", err)
	}
	if time.Now().UTC().After(exp) {
		return 0, ErrNotFound
	}

	return userID, nil
}

func (r *AuthRepository) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete session rows affected: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
