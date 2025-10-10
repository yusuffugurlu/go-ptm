package models

import "time"

type ScheduledTransaction struct {
	Id          uint `gorm:"primaryKey"`
	UserId      uint
	Amount      float64
	ScheduledAt time.Time
	Status      string // pending | processed | cancelled
	CreatedAt   time.Time
}
