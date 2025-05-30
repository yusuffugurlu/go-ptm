package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id           uint   `gorm:"primaryKey"`
	Username     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex; not null"`
	PasswordHash string `gorm:"not null"`
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	balance := Balance{
		UserId:        u.Id,
		Amount:        0.0,
		LastUpdatedAt: time.Now(),
	}

	if err = tx.Create(&balance).Error; err != nil {
		return err
	}

	return nil
}


func NewUser(username, email, password, role string) (*User, error) {
	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
		Role:         role,
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	return user, nil
}