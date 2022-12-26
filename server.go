package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	emdw "github.com/labstack/echo/v4/middleware"
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
