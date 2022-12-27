package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	emdw "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	expn "github.com/dakeeChv/assessment/expense"
	handler "github.com/dakeeChv/assessment/handler"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	PORT   = GetEnv("PORT", "2565")
	PG_URL = os.Getenv("DATABASE_URL")
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("execute(): %v", err)
	}
}

func execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.Open("postgres", PG_URL)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	defer db.Close()

	//Auto initial migration.
	if err := migrateDB(ctx, db); err != nil {
		return fmt.Errorf("failed to initialize db schema: %v", err)
	}

	expense, _ := expn.NewService(ctx, db)
	h, _ := handler.NewHandler(ctx, expense)

	e := newEchoServer()
	h.SetupRoute(e)

	cerr := make(chan error, 1)
	cerr <- e.Start(fmt.Sprintf(":%s", PORT))
	<-cerr

	return nil
}

func newEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(
		emdw.Logger(),
		emdw.Recover(),
		emdw.CORS(),
		emdw.Secure(),
	)
	e.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	return e
}

func migrateDB(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	)`

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
