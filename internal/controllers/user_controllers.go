package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/services"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/response"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

type UserController interface {
	GetAllUsers(e echo.Context) error
	GetUserById(e echo.Context) error
	CreateUser(e echo.Context) error
	UpdateUser(e echo.Context) error
	DeleteUser(e echo.Context) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (u *userController) GetAllUsers(e echo.Context) error {
	users, err := u.userService.GetAllUsers()
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, users)
}

func (u *userController) GetUserById(e echo.Context) error {
	userId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return appErrors.NewBadRequest(err, "invalid user id format")
	}

	user, err := u.userService.GetUserById(userId)
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, user)
}

func (u *userController) UpdateUser(e echo.Context) error {
	var req dtos.UpdateUserRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	userId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return appErrors.NewBadRequest(err, "invalid user id format")
	}

	updatedUser, err := u.userService.UpdateUser(userId, &req)
	if err != nil {
		return err
	}

	return response.Success(e, http.StatusOK, updatedUser)
}

func (u *userController) DeleteUser(e echo.Context) error {
	userId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return appErrors.NewBadRequest(err, "invalid user id format")
	}

	if err := u.userService.DeleteUser(userId); err != nil {
		return err
	}

	return response.NoContent(e)
}

func (u *userController) CreateUser(e echo.Context) error {
	var req dtos.CreateUserRequest
	if err := e.Bind(&req); err != nil {
		return appErrors.NewBadRequest(err, "invalid request format")
	}

	if err := e.Validate(req); err != nil {
		return validator.ProcessValidationErrors(err)
	}

	userModel, err := models.NewUser(
		req.Username,
		req.Email,
		req.Password,
		req.Role,
	)
	if err != nil {
		return appErrors.NewInternalServerError(err)
	}

	createdUser, err := u.userService.CreateUser(userModel)
	if err != nil {
		return err
	}

	return response.Created(e, createdUser)
}
