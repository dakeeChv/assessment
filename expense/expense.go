package expense

import (
	"context"
	"database/sql"
	"fmt"
)

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

func (s *Service) Create(ctx context.Context, in Expense) (Expense, error) {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4)`)
	if err != nil {
		return Expense{}, fmt.Errorf("Create(): db prepare context failure: %w", err)
	}
	result, err := stmt.ExecContext(ctx, in.Title, in.Amount, in.Note, in.Tags)
	if err != nil {
		return Expense{}, fmt.Errorf("Create(): db exec statement failure: %w", err)
	}

	out := in
	out.ID, _ = result.LastInsertId()

	return out, nil
}
