package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/services"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/jwt"
	"github.com/yusuffugurlu/go-project/pkg/response"
)

type BalanceController interface {
	GetCurrentBalance(e echo.Context) error
	GetHistoricalBalances(e echo.Context) error
}

type balanceController struct {
	service services.BalanceService
}

func NewBalanceController() BalanceController {
	return &balanceController{
		service: services.NewBalanceService(),
	}
}

func (b *balanceController) GetCurrentBalance(e echo.Context) error {
	userClaims, ok := e.Get("user").(*jwt.UserClaims)
	if !ok {
		return appErrors.NewUnauthorized(nil, "invalid user context")
	}

	balance, err := b.service.GetUserBalance(uint(userClaims.Id))
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, balance)
}

func (b *balanceController) GetHistoricalBalances(e echo.Context) error {
	userClaims, ok := e.Get("user").(*jwt.UserClaims)
	if !ok {
		return appErrors.NewUnauthorized(nil, "invalid user context")
	}

	historicalBalances, err := b.service.GetHistoricalBalances(uint(userClaims.Id))
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, historicalBalances)
}
