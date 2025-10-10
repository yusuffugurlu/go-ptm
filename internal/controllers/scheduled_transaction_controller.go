package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/services"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

type ScheduledTransactionController struct {
	svc services.ScheduledTransactionService
}

func NewScheduledTransactionController() *ScheduledTransactionController {
	return &ScheduledTransactionController{svc: services.NewScheduledTransactionService()}
}

type scheduleReq struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Date   string  `json:"date" validate:"required"`
}

func (c *ScheduledTransactionController) Schedule(ctx echo.Context) error {
	var req scheduleReq
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := validator.New().Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	at, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format, use RFC3339"})
	}

	st, err := c.svc.Schedule(1, req.Amount, at) // todo: replace 1 with authenticated user id
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to schedule transaction"})
	}
	return ctx.JSON(http.StatusCreated, st)
}
