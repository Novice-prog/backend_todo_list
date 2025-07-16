package models

import (
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	ID          uint           `gorm:"primaryKey"           json:"id"`
	UserID      uint           `gorm:"not null;index"       json:"user_id"`
	Title       string         `gorm:"size:255;not null"    json:"title"`
	Description string         `gorm:"type:text"            json:"description,omitempty"`
	Completed   bool           `gorm:"not null;default:false" json:"completed"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"               json:"-"`
}

type CreateTodoInput struct {
	Title       string `json:"title"       binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=4096"`
}

type UpdateTodoInput struct {
	Title       *string `json:"title,omitempty"       binding:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=4096"`
	Completed   *bool   `json:"completed,omitempty"`
}
