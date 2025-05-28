package models

import "time"

type AuditLog struct {
	ID           uint   `gorm:"primaryKey"`
	EntityType  string `gorm:"not null"`
	EntityID    uint   `gorm:"not null"`
	Action      string `gorm:"not null"`
	Details      string `gorm:"not null"`
	CreatedAt   time.Time
}

func NewAuditLog(entityType string, entityID uint, action string, details string) *AuditLog {
	return &AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Details:    details,
		CreatedAt:  time.Now(),
	}
}