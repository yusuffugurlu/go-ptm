package services

import (
	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/jwt"
	"github.com/yusuffugurlu/go-project/pkg/metrics"
)


type AuthService interface {
	Login(email string, password string) (string, error)
	Register(username string, email string, password string) (*models.User, error)
}

type authService struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &authService{userService: userService}
}

func (a *authService) Login(email string, password string) (string, error) {
	user, err := a.userService.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user.CheckPassword(password) {
		return "", appErrors.NewUnauthorized(nil, "invalid email or password")
	}

	token, err := jwt.GenerateJWT(int(user.Id), user.Email, user.Role)
	if err != nil {
		return "", err
	}

	metrics.IncrementUserLogin();

	return token, nil
}

func (a *authService) Register(username string, email string, password string) (*models.User, error) {
	user, err := models.NewUser(username, email, password, "User")
	if err != nil {
		return nil, appErrors.NewInternalServerError(err)
	}

	if _, err := a.userService.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}