//go:build integration

package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	expn "github.com/dakeeChv/assessment/expense"
	"github.com/dakeeChv/assessment/handler"
)

const (
	port  = 3333
	pgdns = "postgresql://root:root@db/assessment?sslmode=disable"
)

func TestCreateExpense(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := echo.New()
	db, err := sql.Open("postgres", pgdns)
	if err != nil {
		log.Printf("failed to db open: %v\n", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		log.Printf("failed to db connect: %v\n", err)
	}

	expense, _ := expn.NewService(ctx, db)
	h, _ := handler.NewHandler(ctx, expense)
	h.SetupRoute(e)

	go func() {
		e.Start(fmt.Sprintf(":%d", port))
	}()

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 30*time.Second)
		if err != nil {
			log.Printf("failed to dial timeout: %v", err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	rbody := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`

	want := expn.Expense{}
	err = json.Unmarshal([]byte(rbody), &want)
	if err != nil {
		log.Printf("failed to ummarshal want: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", port), strings.NewReader(rbody))
	assert.NoError(t, err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "November 10, 2009")

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NoError(t, err)

	got := expn.Expense{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)
	defer resp.Body.Close()

	if assert.NoError(t, err) {
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Amount, got.Amount)
		assert.Equal(t, want.Note, got.Note)
		assert.Equal(t, want.Tags, got.Tags)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = e.Shutdown(ctx)
	assert.NoError(t, err)
}
