package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time ``
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
