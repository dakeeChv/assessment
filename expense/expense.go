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

func (s *Service) Get(ctx context.Context, id int64) (Expense, error) {
	query := `SELECT id, title, amount, note, tags from expenses where id=$1`

	var out Expense
	err := s.db.QueryRowContext(ctx, query, id).Scan(&out.ID, &out.Title, &out.Amount, &out.Note, pq.Array(&out.Tags))
	if err == sql.ErrNoRows {
		return Expense{}, ErrNoExpense
	}
	if err != nil {
		return Expense{}, fmt.Errorf("Get(): db scan row: %w", err)
	}

	return out, nil
}

func (s *Service) Update(ctx context.Context, in Expense) (Expense, error) {
	query := `UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5 RETURNING id, title, amount, note, tags`

	var out Expense
	err := s.db.QueryRowContext(ctx, query, in.Title, in.Amount, in.Note, pq.Array(in.Tags), in.ID).Scan(&out.ID, &out.Title, &out.Amount, &out.Note, pq.Array(&out.Tags))
	if err == sql.ErrNoRows {
		return Expense{}, ErrNoExpense
	}
	if err != nil {
		return Expense{}, fmt.Errorf("Update(): db scan row: %w", err)
	}

	return out, nil
}

func (s *Service) List(ctx context.Context) ([]Expense, error) {
	query := `SELECT id, title, amount, note, tags from expenses`

	out := make([]Expense, 0)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []Expense{}, fmt.Errorf("List(): db query context: %w", err)
	}

	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			return []Expense{}, fmt.Errorf("List(): db scan row: %w", err)
		}
		out = append(out, expense)
	}

	return out, nil
}
