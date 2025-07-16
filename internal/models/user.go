package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint           `gorm:"primaryKey"                     json:"id"`
	Username     string         `gorm:"size:32;uniqueIndex;not null"   json:"username"`
	Email        string         `gorm:"size:254;uniqueIndex;not null"  json:"email"`
	PasswordHash string         `gorm:"size:60;not null"               json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                          json:"-"`
	Todos        []Todo         `gorm:"constraint:OnDelete:CASCADE;"   json:"-"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email"    binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
