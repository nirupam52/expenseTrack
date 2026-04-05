package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nirupam52/expenseTrack/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) CreateExpense(ctx context.Context, groupID *int64, userID int64, description string, amount float64, date string) error {
	query := `INSERT INTO expenses (group_id, paid_by, description, amount, date) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, groupID, userID, description, amount, date)
	if err != nil {
		return fmt.Errorf("Error creating expense: %w", err)
	}

	return nil
}

func (r *ExpenseRepository) GetExpenseByID(ctx context.Context, id int64) (*models.Expense, error) {
	query := `SELECT id, group_id, paid_by, description, amount, date, created_at FROM expenses WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var e models.Expense
	err := row.Scan(&e.ID, &e.GroupID, &e.PaidBy, &e.Description, &e.Amount, &e.Date, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("expense not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get expense by id: %w", err)
	}

	return &e, nil
}

func (r *ExpenseRepository) ListExpensesByUser(ctx context.Context, userID int64) ([]*models.Expense, error) {
	query := `SELECT id, group_id, paid_by, description, amount, date, created_at FROM expenses WHERE paid_by = ? ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.GroupID, &e.PaidBy, &e.Description, &e.Amount, &e.Date, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}
		expenses = append(expenses, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list expenses rows: %w", err)
	}

	return expenses, nil
}

func (r *ExpenseRepository) UpdateExpense(ctx context.Context, id int64, description string, amount float64, date string) error {
	query := `UPDATE expenses SET description = ?, amount = ?, date = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, description, amount, date, id)
	if err != nil {
		return fmt.Errorf("update expense: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update expense, rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

func (r *ExpenseRepository) DeleteExpense(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete expense: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete expense, rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}
