package repositories

import (
	"errors"
	"fmt"

	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id uint) (*models.Transaction, error)
	GetByUserID(userID uint) ([]*models.Transaction, error)
	GetHistoryByUserID(userID uint, limit, offset int) ([]*models.Transaction, error)
	GetAll(limit, offset int) ([]*models.Transaction, error)
	Transfer(fromUserID, toUserID uint, amount float64) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	if err := r.db.Create(transaction).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to create transaction")
	}
	return nil
}

func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Preload("FromUser.Balance").Preload("ToUser.Balance").First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound(err, fmt.Sprintf("transaction with id %d not found", id))
		}
		return nil, appErrors.NewDatabaseError(err, "failed to get transaction")
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByUserID(userID uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Preload("FromUser.Balance").Preload("ToUser.Balance").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, appErrors.NewDatabaseError(err, "failed to get user transactions")
	}
	return transactions, nil
}

func (r *transactionRepository) GetHistoryByUserID(userID uint, limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	if err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Preload("FromUser.Balance").Preload("ToUser.Balance").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, appErrors.NewDatabaseError(err, "failed to get user transaction history")
	}
	return transactions, nil
}

func (r *transactionRepository) GetAll(limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := r.db.Preload("FromUser.Balance").Preload("ToUser.Balance").Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, appErrors.NewDatabaseError(err, "failed to get all transactions")
	}
	return transactions, nil
}

func (r *transactionRepository) Transfer(fromUserID, toUserID uint, amount float64) error {
	if amount <= 0 {
		return appErrors.NewBadRequest(nil, "amount must be positive")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var fromBalance, toBalance models.Balance

		if err := tx.Where("user_id = ?", fromUserID).First(&fromBalance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.NewNotFound(err, fmt.Sprintf("balance not found for user %d", fromUserID))
			}
			return appErrors.NewDatabaseError(err, "failed to get sender balance")
		}

		if err := tx.Where("user_id = ?", toUserID).First(&toBalance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.NewNotFound(err, fmt.Sprintf("balance not found for user %d", toUserID))
			}
			return appErrors.NewDatabaseError(err, "failed to get receiver balance")
		}

		if fromBalance.Amount < amount {
			return appErrors.NewConflict(nil, fmt.Sprintf("insufficient funds: requested %.2f, available %.2f", amount, fromBalance.Amount))
		}

		fromBalance.Amount -= amount
		toBalance.Amount += amount

		if err := tx.Save(&fromBalance).Error; err != nil {
			return appErrors.NewDatabaseError(err, "failed to update sender balance")
		}

		if err := tx.Save(&toBalance).Error; err != nil {
			return appErrors.NewDatabaseError(err, "failed to update receiver balance")
		}

		transaction := &models.Transaction{
			FromUserId: &fromUserID,
			ToUserId:   &toUserID,
			Amount:     amount,
			Type:       "transfer",
			Status:     "completed",
		}

		if err := tx.Create(transaction).Error; err != nil {
			return appErrors.NewDatabaseError(err, "failed to create transaction record")
		}

		return nil
	})
}
