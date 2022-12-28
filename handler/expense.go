package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	expn "github.com/dakeeChv/assessment/expense"
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
	resp, err := h.expense.Create(ctx, req)
	if err != nil {
		ref := uuid.New()
		log.Printf("\nlogId: %s, %v\n", ref, err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    500,
			"status":  "Internal Server Error",
			"Message": fmt.Sprintf("failed to processing request, refer: %s", ref),
		})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetExpense(c echo.Context) error {
	rid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"status":  "Bad Request",
			"Message": "failed to binding param, Please pass a valid param",
		})
	}

	var id int64 = int64(rid)
	ctx := c.Request().Context()
	resp, err := h.expense.Get(ctx, id)
	if errors.Is(err, expn.ErrNoExpense) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"code":    404,
			"status":  "Not Found",
			"Message": fmt.Sprintf("Not Found, a expense with ID: %d", id),
		})
	}

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

func (h *Handler) UpdateExpense(c echo.Context) error {
	var req expn.Expense
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"status":  "Bad Request",
			"Message": "failed to binding json body, Please pass a valid json body",
		})
	}

	rid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"status":  "Bad Request",
			"Message": "failed to binding param, Please pass a valid param",
		})
	}
	req.ID = int64(rid)

	ctx := c.Request().Context()
	resp, err := h.expense.Update(ctx, req)
	if errors.Is(err, expn.ErrNoExpense) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"code":    404,
			"status":  "Not Found",
			"Message": fmt.Sprintf("Not Found, a expense with ID: %d", req.ID),
		})
	}

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
