package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
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
	var req dtos.TransactionRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	
}

func (t *transactionController) Transfer(e echo.Context) error {
	panic("unimplemented")
}

func (t *transactionController) Withdraw(e echo.Context) error {
	panic("unimplemented")
}
