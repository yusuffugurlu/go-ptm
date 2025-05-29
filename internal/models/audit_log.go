package models

import "time"

type AuditLog struct {
	Id           uint   `gorm:"primaryKey"`
	EntityType  string `gorm:"not null"`
	EntityId   uint   `gorm:"not null"`
	Action      string `gorm:"not null"`
	Details      string `gorm:"not null"`
	CreatedAt   time.Time
}

func NewAuditLog(entityType string, entityId uint, action string, details string) *AuditLog {
	return &AuditLog{
		EntityType: entityType,
		EntityId:   entityId,
		Action:     action,
		Details:    details,
		CreatedAt:  time.Now(),
	}
}