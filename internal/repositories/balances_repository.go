package repositories

import (
	"fmt"

	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"gorm.io/gorm"
)

type BalancesRepository interface {
	Create(balance *models.Balance) error
	GetByUserId(userId uint) (*models.Balance, error)
}

type balancesRepository struct {
	db *gorm.DB
}

func NewBalancesRepository(db *gorm.DB) BalancesRepository {
	return &balancesRepository{}
}

func (b *balancesRepository) Create(balance *models.Balance) error {
	if err := b.db.Create(balance).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to create balance")
	}
	return nil
}

func (b *balancesRepository) GetByUserId(userId uint) (*models.Balance, error) {
	var balance models.Balance
	if err := b.db.Where("userId = ?", userId).First(&balance).Error; err != nil {
		return nil, appErrors.NewNotFound(err, fmt.Sprintf("balance for user with id %d not found", userId))
	}

	return &balance, nil
}
