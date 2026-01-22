package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAT time.Time
	Email     string `gorm:"unique" binding:"required,email"`
	Password  string `binding:"required,min=6"`
}
