package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	emdw "github.com/labstack/echo/v4/middleware"
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
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
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
