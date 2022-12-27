package handler

import (
	"context"

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
