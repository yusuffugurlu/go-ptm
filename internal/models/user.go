package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex; not null"`
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