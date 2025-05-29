package models

import "time"

type Transaction struct {
	Id         uint `gorm:"primaryKey"`
	FromUserId uint
	ToUserId   uint
	Amount     float64
	Type       string
	Status     string
	CreatedAt  time.Time

	FromUser *User `gorm:"foreignKey:FromUserId"`
	ToUser   *User `gorm:"foreignKey:ToUserId"`
}
