package models

import "time"

type Balance struct {
	UserId        uint `gorm:"primaryKey"`
	Amount        float64
	LastUpdatedAt time.Time

	User          *User   `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
}
