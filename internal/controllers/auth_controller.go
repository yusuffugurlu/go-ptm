package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/services"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/response"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

type AuthController interface {
	Login(e echo.Context) error
	Register(e echo.Context) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

func (a *authController) Login(e echo.Context) error {
	var req dtos.LoginRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	token, err := a.authService.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, map[string]string{"token": token})
}

func (a *authController) Register(e echo.Context) error {
	var req dtos.RegisterRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format");
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	user, err := a.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, user)
}
