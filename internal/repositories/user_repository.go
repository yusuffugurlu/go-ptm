package repositories

import (
	"gorm.io/gorm"
	"errors"
	"fmt"

	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
)

type UserRepository interface {
	Create(user *models.User) error
	GetById(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(user *models.User) error {
	if err := u.db.Create(user).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to create user")
	}
	return nil
}

func (u *userRepository) Delete(id int) error {
	result := u.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return appErrors.NewDatabaseError(result.Error, fmt.Sprintf("failed to delete user with id %d", id))
	}

	if result.RowsAffected == 0 {
		return appErrors.NewNotFound(nil, fmt.Sprintf("user with id %d not found", id))
	}

	return nil
}

func (u *userRepository) GetAll() ([]models.User, error) {
	var users []models.User

	if err := u.db.Find(&users).Error; err != nil {
		return nil, appErrors.NewDatabaseError(err, "failed to fetch all users")
	}

	return users, nil
}

func (u *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound(err, fmt.Sprintf("user with email %s not found", email))
		}
		return nil, appErrors.NewDatabaseError(err, "failed to fetch user by email")
	}
	return &user, nil
}

func (u *userRepository) GetById(id int) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound(err, fmt.Sprintf("user with id %d not found", id))
		}
		return nil, appErrors.NewDatabaseError(err, "failed to fetch user by id")
	}
	return &user, nil
}

func (u *userRepository) Update(user *models.User) error {
	result := u.db.Save(user)
	if result.Error != nil {
		return appErrors.NewDatabaseError(result.Error, fmt.Sprintf("failed to update user with id %d", user.Id))
	}

	if result.RowsAffected == 0 {
		return appErrors.NewNotFound(nil, fmt.Sprintf("user with id %d not found", user.Id))
	}

	return nil
}
