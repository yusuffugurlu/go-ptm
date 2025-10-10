package repositories

import (
	"time"

	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"gorm.io/gorm"
)

type ScheduledTransactionRepository interface {
	Create(tx *models.ScheduledTransaction) error
	GetDue(now time.Time) ([]*models.ScheduledTransaction, error)
	MarkProcessed(id uint) error
}

type scheduledTransactionRepository struct {
	db *gorm.DB
}

func NewScheduledTransactionRepository(db *gorm.DB) ScheduledTransactionRepository {
	return &scheduledTransactionRepository{db: db}
}

func (r *scheduledTransactionRepository) Create(tx *models.ScheduledTransaction) error {
	if err := r.db.Create(tx).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to create scheduled transaction")
	}
	return nil
}

func (r *scheduledTransactionRepository) GetDue(now time.Time) ([]*models.ScheduledTransaction, error) {
	var list []*models.ScheduledTransaction
	if err := r.db.Where("scheduled_at <= ? AND status = ?", now, "pending").Find(&list).Error; err != nil {
		return nil, appErrors.NewDatabaseError(err, "failed to query due scheduled transactions")
	}
	return list, nil
}

func (r *scheduledTransactionRepository) MarkProcessed(id uint) error {
	if err := r.db.Model(&models.ScheduledTransaction{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "processed"}).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to mark scheduled transaction processed")
	}
	return nil
}
