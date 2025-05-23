package services

import (
	"errors"
	"fmt"

	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserById(id int) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(id int, user *dtos.UpdateUserRequest) (*models.User, error)
	DeleteUser(userId int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepo: userRepository}
}

func (u *userService) CreateUser(user *models.User) (*models.User, error) {
	_, err := u.userRepo.GetByEmail(user.Email)
	if err == nil {
		return nil, appErrors.NewConflict(nil, fmt.Sprintf("Email %s is already in use", user.Email))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := user.HashPassword(); err != nil {
		return nil, appErrors.NewInternalServerError(err)
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) DeleteUser(userId int) error {
	_, err := u.userRepo.GetById(userId)
	if err != nil {
		return err
	}

	if err := u.userRepo.Delete(userId); err != nil {
		return appErrors.NewDatabaseError(err, fmt.Sprintf("Failed to delete user with ID %d", userId))
	}

	return nil
}

func (u *userService) GetAllUsers() ([]models.User, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, appErrors.NewDatabaseError(err, "Failed to fetch all users")
	}

	return users, nil
}

func (u *userService) GetUserById(id int) (*models.User, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) UpdateUser(id int, userData *dtos.UpdateUserRequest) (*models.User, error) {
	existingUser, err := u.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	if userData.Email != "" && userData.Email != existingUser.Email {
		user, err := u.userRepo.GetByEmail(userData.Email)
		if err == nil && user != nil {
			return nil, appErrors.NewConflict(nil, fmt.Sprintf("Email %s is already in use", userData.Email))
		} else if !errors.Is(err, gorm.ErrRecordNotFound) && !appErrors.IsAppError(err) {
			return nil, appErrors.NewDatabaseError(err, "Failed to check email uniqueness")
		}
	}

	if userData.Username != "" {
		existingUser.Username = userData.Username
	}

	if userData.Email != "" {
		existingUser.Email = userData.Email
	}

	if userData.Password != "" {
		existingUser.PasswordHash = userData.Password
		if err := existingUser.HashPassword(); err != nil {
			return nil, appErrors.NewInternalServerError(err)
		}
	}

	if err := u.userRepo.Update(existingUser); err != nil {
		return nil, appErrors.NewDatabaseError(err, fmt.Sprintf("Failed to update user with ID %d", id))
	}

	return existingUser, nil
}
