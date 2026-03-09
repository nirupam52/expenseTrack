package models

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type Expense struct {
	ID          int64   `json:"id"`
	PaidBy      int64   `json:"paid_by"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}
