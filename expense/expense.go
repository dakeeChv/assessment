package expense

import (
	"context"
	"database/sql"
	"fmt"
)

// Expense is  Expense tracking model.
type Expense struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Service struct {
	db *sql.DB
}

// NewService returns expense service.
func NewService(_ context.Context, db *sql.DB) (*Service, error) {
	return &Service{db: db}, nil
}

func (s *Service) Create(ctx context.Context, in Expense) (Expense, error) {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4)`)
	if err != nil {
		return Expense{}, fmt.Errorf("Create(): db prepare context failure: %w", err)
	}
	if err != nil {
		return Expense{}, fmt.Errorf("Create(): db exec statement failure: %w", err)
	}

	out := in
	out.ID, _ = result.LastInsertId()

	return out, nil
}
