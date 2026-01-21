package models

import (
	"time"

	"gorm.io/datatypes"
)

//ID
//CreatedAt
//UpdatedAt
//Name
//Description
//Image
//Skills

type Project struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Image       string
	Skills      datatypes.JSONSlice[string] `gorm:"type:json"`
}
