package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/process"
	"github.com/yusuffugurlu/go-project/internal/services"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/jwt"
	"github.com/yusuffugurlu/go-project/pkg/response"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

type TransactionController interface {
	Deposit(e echo.Context) error
	Withdraw(e echo.Context) error
	Transfer(e echo.Context) error
	Debit(e echo.Context) error
	GetHistory(e echo.Context) error
	GetByID(e echo.Context) error
	GetAllTransactions(e echo.Context) error
}

type transactionController struct {
	service services.TransactionService
}

func NewTransactionController() TransactionController {
	return &transactionController{
		service: services.NewTransactionService(),
	}
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

	req.Date = time.Now()

	return response.Success(e, http.StatusOK, req)
}

func (t *transactionController) Withdraw(e echo.Context) error {
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
		Type:   process.WithdrawTransaction,
	}

	return response.Success(e, http.StatusOK, req)
}

func (t *transactionController) Transfer(e echo.Context) error {
	var req dtos.TransferRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	userClaims, ok := e.Get("user").(*jwt.UserClaims)
	if !ok {
		return appErrors.NewUnauthorized(nil, "invalid user context")
	}

	if err := t.service.TransferBetweenUsers(uint(userClaims.Id), req.ToUserID, req.Amount); err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, map[string]interface{}{
		"message":    "transfer completed successfully",
		"amount":     req.Amount,
		"to_user_id": req.ToUserID,
	})
}

func (t *transactionController) GetHistory(e echo.Context) error {
	userClaims, ok := e.Get("user").(*jwt.UserClaims)
	if !ok {
		return appErrors.NewUnauthorized(nil, "invalid user context")
	}

	limitStr := e.QueryParam("limit")
	offsetStr := e.QueryParam("offset")

	limit := 10 // default
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	transactions, err := t.service.GetTransactionHistory(uint(userClaims.Id), limit, offset)
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, transactions)
}

func (t *transactionController) GetByID(e echo.Context) error {
	idStr := e.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return appErrors.NewBadRequest(err, "invalid transaction id")
	}

	transaction, err := t.service.GetTransactionByID(uint(id))
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, transaction)
}

func (t *transactionController) GetAllTransactions(e echo.Context) error {
	limitStr := e.QueryParam("limit")
	offsetStr := e.QueryParam("offset")

	limit := 50 // default for admin
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	transactions, err := t.service.GetAllTransactions(limit, offset)
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, map[string]interface{}{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
		"count":        len(transactions),
	})
}

func (t *transactionController) Debit(e echo.Context) error {
	var req dtos.DebitRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	userClaims, ok := e.Get("user").(*jwt.UserClaims)
	if !ok {
		return appErrors.NewUnauthorized(nil, "invalid user context")
	}

	if err := t.service.DebitFromUser(uint(userClaims.Id), req.Amount); err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, map[string]interface{}{
		"message": "debit completed successfully",
		"amount":  req.Amount,
	})
}
