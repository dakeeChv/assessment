package handler

import (
	"fmt"
	"log"
	"net/http"

	expn "github.com/dakeeChv/assessment/expense"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateExpense(c echo.Context) error {
	var req expn.Expense
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"status":  "Bad Request",
			"Message": "failed to binding json body, Please pass a valid json body",
		})
	}

	ctx := c.Request().Context()
	resp, err := h.svc.Create(ctx, req)
	if err != nil {
		ref := uuid.New()
		log.Printf("\nlogId: %s, %v\n", ref, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    500,
			"status":  "Internal Server Error",
			"Message": fmt.Sprintf("failed to processing request, refer: %s", ref),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
