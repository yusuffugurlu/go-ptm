package repositories

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"gorm.io/gorm"
)

type BalancesRepository interface {
	Create(balance *models.Balance) error
	GetByUserId(userId uint) (*models.Balance, error)
	Deposit(userId uint, amount float64) error
	Withdraw(userId uint, amount float64) error
}

type balancesRepository struct {
	db *gorm.DB
	mu sync.RWMutex
}

func NewBalancesRepository(db *gorm.DB) BalancesRepository {
	return &balancesRepository{
		db: db,
		mu: sync.RWMutex{},
	}
}

func (b *balancesRepository) Create(balance *models.Balance) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.db.Create(balance).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to create balance")
	}
	return nil
}

func (b *balancesRepository) GetByUserId(userId uint) (*models.Balance, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var balance models.Balance
	if err := b.db.Where("user_id = ?", userId).First(&balance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound(err, fmt.Sprintf("balance not found for user %d", userId))
		}
		return nil, appErrors.NewDatabaseError(err, fmt.Sprintf("failed to get balance for user %d", userId))
	}
	return &balance, nil
}

func (b *balancesRepository) Deposit(userId uint, amount float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if amount <= 0 {
		return appErrors.NewBadRequest(nil, "amount must be positive")
	}

	return b.db.Transaction(func(tx *gorm.DB) error {
		var balance models.Balance

		if err := tx.Where("user_id = ?", userId).First(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.NewNotFound(err, fmt.Sprintf("balance not found for user %d, cannot deposit", userId))
			}
			return appErrors.NewDatabaseError(err, "failed to get balance for deposit")
		}

		balance.Amount += amount
		if err := tx.Save(&balance).Error; err != nil {
			return appErrors.NewDatabaseError(err, "failed to update balance after deposit")
		}
		return nil
	})
}

func (b *balancesRepository) Withdraw(userId uint, amount float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	return b.db.Transaction(func(tx *gorm.DB) error {
		var balance models.Balance

		if err := tx.Where("user_id = ?", userId).First(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.NewNotFound(err, fmt.Sprintf("balance not found for user %d, cannot withdraw", userId))
			}
			return appErrors.NewDatabaseError(err, "failed to get balance for withdrawal")
		}

		if balance.Amount < amount {
			return appErrors.NewConflict(nil, fmt.Sprintf("insufficient funds for user %d: requested %.2f, available %.2f", userId, amount, balance.Amount))
		}

		balance.Amount -= amount
		if err := tx.Save(&balance).Error; err != nil {
			return appErrors.NewDatabaseError(err, "failed to update balance after withdrawal")
		}
		return nil
	})
}