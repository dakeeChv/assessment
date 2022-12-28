package expense_test

import (
	"context"
	"database/sql"
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
		in := expn.Expense{
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id, title, amount, note, tags`)).
			ExpectQuery().
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
					AddRow(1, "strawberry smoothie", 79.00, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})),
			).
			WithArgs(in.Title, in.Amount, in.Note, pq.Array(in.Tags))

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

	t.Run("Failed to db scan row", func(t *testing.T) {
		want := errors.New(`sql: Scan error on column index 4, name "tags": unsupported Scan, storing driver.Value type string into type *[]string`)
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id, title, amount, note, tags`)).
			ExpectQuery().
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
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id, title, amount, note, tags`)).
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
}

func TestGetExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		want := expn.Expense{
			ID:     1,
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, amount, note, tags from expenses where id=$1")).
			WithArgs(want.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
					AddRow(1, "strawberry smoothie", 79.00, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"})),
			)

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Get(ctx, want.ID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, want.ID, got.ID)
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Amount, got.Amount)
		assert.Equal(t, want.Note, got.Note)
		assert.Equal(t, want.Tags, got.Tags)
	})

	t.Run("Error no row", func(t *testing.T) {
		want := expn.Expense{
			ID: 1,
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, amount, note, tags from expenses where id=$1")).
			WithArgs(want.ID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Get(ctx, want.ID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.NotEmpty(t, err)
		assert.Equal(t, expn.Expense{}, got)
		assert.ErrorIs(t, err, expn.ErrNoExpense)
	})

	t.Run("Some error", func(t *testing.T) {
		var id int64 = 1
		want := errors.New("some error")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, amount, note, tags from expenses where id=$1")).
			WithArgs(id).
			WillReturnError(want)

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Get(ctx, id)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.NotEmpty(t, err)
		assert.Equal(t, expn.Expense{}, got)
		assert.ErrorIs(t, err, want)
	})
}

func TestUpdateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		want := expn.Expense{
			ID:     123,
			Title:  "apple smoothie",
			Amount: 89,
			Note:   "no discount",
			Tags:   []string{"beverage"},
		}

		mock.ExpectQuery(regexp.QuoteMeta("UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5 RETURNING id, title, amount, note, tags")).
			WithArgs(want.Title, want.Amount, want.Note, pq.Array(want.Tags), want.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
					AddRow(123, "apple smoothie", 89.00, "no discount", pq.Array([]string{"beverage"})),
			)

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Update(ctx, want)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, want.ID, got.ID)
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Amount, got.Amount)
		assert.Equal(t, want.Note, got.Note)
		assert.Equal(t, want.Tags, got.Tags)
	})

	t.Run("Error no row", func(t *testing.T) {
		want := expn.Expense{
			ID:     123,
			Title:  "apple smoothie",
			Amount: 89,
			Note:   "no discount",
			Tags:   []string{"beverage"},
		}

		mock.ExpectQuery(regexp.QuoteMeta("UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5 RETURNING id, title, amount, note, tags")).
			WithArgs(want.Title, want.Amount, want.Note, pq.Array(want.Tags), want.ID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Update(ctx, want)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.ErrorIs(t, err, expn.ErrNoExpense)
		assert.Empty(t, got.ID)
		assert.Empty(t, got.Title)
		assert.Empty(t, got.Amount)
		assert.Empty(t, got.Note)
		assert.Empty(t, got.Tags)
	})

	t.Run("Some error", func(t *testing.T) {
		want := expn.Expense{
			ID:     123,
			Title:  "apple smoothie",
			Amount: 89,
			Note:   "no discount",
			Tags:   []string{"beverage"},
		}

		errwant := errors.New("some error")

		mock.ExpectQuery(regexp.QuoteMeta("UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5 RETURNING id, title, amount, note, tags")).
			WithArgs(want.Title, want.Amount, want.Note, pq.Array(want.Tags), want.ID).
			WillReturnError(errwant)

		ctx := context.Background()
		expense, _ := expn.NewService(ctx, db)

		got, err := expense.Update(ctx, want)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.ErrorIs(t, err, errwant)
		assert.Empty(t, got.ID)
		assert.Empty(t, got.Title)
		assert.Empty(t, got.Amount)
		assert.Empty(t, got.Note)
		assert.Empty(t, got.Tags)
	})
}
