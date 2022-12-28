package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var ErrNoExpense = errors.New("no expense")

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
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id, title, amount, note, tags`)
	if err != nil {
		return Expense{}, fmt.Errorf("Create(): db prepare context failure: %w", err)
	}

	err = stmt.QueryRowContext(ctx, in.Title, in.Amount, in.Note, pq.Array(in.Tags)).Scan(&in.ID, &in.Title, &in.Amount, &in.Note, pq.Array(&in.Tags))
	if err != nil {
		return Expense{}, fmt.Errorf("Create(): db scan row: %w", err)
	}

	out := in

	return out, nil
}
