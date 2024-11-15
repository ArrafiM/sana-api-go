package models

import (
	"gorm.io/gorm"
)

type DeviceToken struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

type DeviceTokenCreate struct {
	UserID int `json:"user_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func (DeviceTokenCreate) TableName() string {
	return "device_tokens"
}
