package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/yusuffugurlu/go-project/internal/models"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
)

type AuditLogRepository interface {
	GetAll() ([]models.AuditLog, error)
	Create(log *models.AuditLog) error
	GetByEntityType(entityType string) ([]models.AuditLog, error)
	Delete(id int) error
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (a *auditLogRepository) GetAll() ([]models.AuditLog, error) {
	var logs []models.AuditLog

	if err := a.db.Find(&logs).Error; err != nil {
		return nil, appErrors.NewDatabaseError(err, "failed to fetch all logs")
	}

	return logs, nil
}

func (a *auditLogRepository) Create(log *models.AuditLog) error {
	if err := a.db.Create(log).Error; err != nil {
		return appErrors.NewDatabaseError(err, "failed to create audit log")
	}
	return nil
}

func (a *auditLogRepository) GetByEntityType(entityType string) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	if err := a.db.Where("entity_type = ?", entityType).Find(&logs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewDatabaseError(err, fmt.Sprintf("no audit logs found for entity type %s", entityType))
		}
		return nil, appErrors.NewDatabaseError(err, "failed to fetch audit logs by entity type")
	}
	return logs, nil
}

func (a *auditLogRepository) Delete(id int) error {
	result := a.db.Delete(&models.AuditLog{}, id)
	if result.Error != nil {
		return appErrors.NewDatabaseError(result.Error, fmt.Sprintf("failed to delete audit log with id %d", id))
	}

	if result.RowsAffected == 0 {
		return appErrors.NewNotFound(nil, fmt.Sprintf("audit log with id %d not found", id))
	}

	return nil
}
