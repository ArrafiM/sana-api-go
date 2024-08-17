package models

import "gorm.io/gorm"

type Merchant struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'active'"`
}

type MerchantCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type MerchantUser struct {
	Merchant
	User User `json:"user"`
}

func (MerchantUser) TableName() string {
	return "merchants"
}
