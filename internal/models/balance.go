package models

import "time"

type Balance struct {
    UserID        uint    `gorm:"primaryKey"`
    Amount        float64
    LastUpdatedAt time.Time

    User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}