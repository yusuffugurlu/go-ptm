package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/process"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/response"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

type TransactionController interface {
	Deposit(e echo.Context) error
	Withdraw(e echo.Context) error
	Transfer(e echo.Context) error
}

type transactionController struct{}

func NewTransactionController() TransactionController {
	return &transactionController{}
}

func (t *transactionController) Deposit(e echo.Context) error {
	var req process.Transaction
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	process.JobQueue <- process.Transaction{
		Amount: req.Amount,
		UserId: req.UserId,
		Type:   process.DepositTransaction,
	}

	return response.Success(e, http.StatusOK, req)
}

func (t *transactionController) Transfer(e echo.Context) error {
	panic("unimplemented")
}

func (t *transactionController) Withdraw(e echo.Context) error {
	panic("unimplemented")
}
