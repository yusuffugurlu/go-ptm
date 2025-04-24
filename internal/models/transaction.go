package models

import "time"

type Transaction struct {
    ID         uint      `gorm:"primaryKey"`
    FromUserID uint
    ToUserID   uint
    Amount     float64
    Type       string
    Status     string
    CreatedAt  time.Time

    FromUser   *User `gorm:"foreignKey:FromUserID"`
    ToUser     *User `gorm:"foreignKey:ToUserID"`
}