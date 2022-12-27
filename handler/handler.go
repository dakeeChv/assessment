package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	expn "github.com/dakeeChv/assessment/expense"
)

// Handler manages http transports.
type Handler struct {
	expense *expn.Service
}

// NewHandler returns handler instance.
func NewHandler(_ context.Context, expense *expn.Service) (*Handler, error) {
	return &Handler{
		expense: expense,
	}, nil
}

func (h *Handler) SetupRoute(e *echo.Echo) {
	// Sample authentication with pare data value.
	cmdw := func() echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				val := c.Request().Header.Get(echo.HeaderAuthorization)
				_, err = time.Parse("January 02, 2006", val)
				if err == nil {
					next(c)
				}
				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "Unauthorized",
					Internal: err,
				}
			}
		}
	}

	v1 := e.Group("")
	v1.Use(cmdw())
	v1.POST("/expenses", h.CreateExpense)
}
