package services

import (
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/repositories"
)

type AuditLogService interface {
	GetAllAuditLogs() ([]models.AuditLog, error)
	CreateAuditLog(entity_id int, entity_type string, action string, details string) error
	GetAuditLogsByEntityType(entityType string) ([]models.AuditLog, error)
	DeleteAuditLog(id int) error
}

type auditLogService struct {
	auditLogRepo repositories.AuditLogRepository
}

func NewAuditLogService(auditLogRepo repositories.AuditLogRepository) AuditLogService {
	return &auditLogService{auditLogRepo: auditLogRepo}
}

func (a *auditLogService) GetAllAuditLogs() ([]models.AuditLog, error) {
	logs, err := a.auditLogRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (a *auditLogService) CreateAuditLog(entity_id int, entity_type string, action string, details string) error {
	log := models.NewAuditLog(entity_type, uint(entity_id), action, details)
	if err := a.auditLogRepo.Create(log); err != nil {
		return err
	}

	return nil
}

func (a *auditLogService) GetAuditLogsByEntityType(entityType string) ([]models.AuditLog, error) {
	logs, err := a.auditLogRepo.GetByEntityType(entityType)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (a *auditLogService) DeleteAuditLog(id int) error {
	if err := a.auditLogRepo.Delete(id); err != nil {
		return err
	}

	return nil
}