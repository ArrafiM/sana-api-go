package models

import (
	"time"

	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	Email      string    `json:"email"`
	Properties string    `json:"properties"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-" gorm:"autoUpdateTime"`
}

type FeedbackCreate struct {
	Email      string `json:"email" binding:"required"`
	Properties string `json:"properties" binding:"required"`
}
