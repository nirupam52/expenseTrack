package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nirupam52/expenseTrack/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, name string, email string, passwordHash string) error {
	query := "INSERT INTO users (name, email, password_hash) VALUES (?,?,?)"
	_, err := r.db.ExecContext(ctx, query, name, email, passwordHash)
	if err != nil {
		return fmt.Errorf("Error when creating user")
	}
	return nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
	var user models.User
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User not found")
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching by user ID : %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, name, email, created_at FROM users WHERE email = ?"
	var user models.User
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User not found")
	}
	if err != nil {
		return nil, fmt.Errorf("Error fetching by user email : %w", err)
	}
	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := "SELECT id, name, email, created_at FROM users"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(" list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user : %w", err)
		}
		users = append(users, &u)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list users rows: %w", err)
	}
	return users, nil

}
