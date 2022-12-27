package handler

import (
	"context"

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
	v1 := e.Group("/expenses")
	v1.POST("/", h.CreateExpense)
}
