package services

import (
	"errors"
	"fmt"

	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	"gorm.io/gorm"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(id int, user *dtos.UpdateUserRequest) (*models.User, error)
	DeleteUser(userId int) error
}

type userService struct {
	userRepo repositories.UserRepository
	logService AuditLogService
}

func NewUserService(userRepository repositories.UserRepository, logService AuditLogService) UserService {
	return &userService{
		userRepo: userRepository,
		logService: logService,
	}
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

	if err := u.logService.CreateAuditLog(int(user.ID), "user", "create", fmt.Sprintf("user %d created", user.ID)); err != nil {
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
		return err
	}

	if err := u.logService.CreateAuditLog(userId, "user", "delete", fmt.Sprintf("user %d deleted", userId)); err != nil {
		return err
	}

	return nil
}

func (u *userService) GetAllUsers() ([]models.User, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
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
		return nil, err
	}

	if err := u.logService.CreateAuditLog(id, "user", "update", fmt.Sprintf("user %d updated", id)); err != nil {
		return nil, err
	}

	return existingUser, nil
}
