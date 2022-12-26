package expense_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/assert"

	expn "github.com/dakeeChv/assessment/expense"
)

func TestCreateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4)`)).
			ExpectExec().
			WithArgs("strawberry smoothie", 79.00, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
			WillReturnResult(sqlmock.NewResult(1, 1))

		in := expn.Expense{
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}
		want := in

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Create(ctx, in)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Nil(t, err)
		assert.NotEmpty(t, got.ID)
		assert.NotEqual(t, want.ID, got.ID)
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Amount, got.Amount)
		assert.Equal(t, want.Note, got.Note)
		assert.Equal(t, want.Tags, got.Tags)
	})

	t.Run("Failed to db exec", func(t *testing.T) {
		want := errors.New("sql: converting argument $4 type: unsupported type []string, a slice of string")
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4)`)).
			ExpectExec().
			WithArgs("strawberry smoothie", 79.00, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})).
			WillReturnError(want)

		in := expn.Expense{
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Create(ctx, in)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t, expn.Expense{}, got)
		assert.ErrorIs(t, err, want)
	})

	t.Run("Failed to db prepare", func(t *testing.T) {
		want := errors.New("call to Prepare statement with query 'INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3)', was not expected")
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4)`)).WillReturnError(want)

		in := expn.Expense{
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Create(ctx, in)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t, expn.Expense{}, got)
		assert.ErrorIs(t, err, want)
	})
}
