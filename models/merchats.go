package models

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	UserID      uint           `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status" gorm:"default:'active'"`
	Picture     string         `json:"picture"`
	Color       string         `json:"color"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-" gorm:"autoUpdateTime"`
	Merchandise *[]Merchandise `json:"merchandise" gorm:"foreignKey:MerchantID"`
}

type MerchantCreate struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Picture     *multipart.FileHeader `form:"picture" binding:"required"`
	Color       *string               `form:"color"`
}

type MerchantUser struct {
	Merchant
	User          *User                   `json:"user"`
	LandingImages *[]MerchantLandingImage `json:"landing_images" gorm:"foreignKey:merchant_id"`
}

func (MerchantUser) TableName() string {
	return "merchants"
}

type MerchantDtl struct {
	MerchantUser
	Merchandise *[]Merchandise `json:"merchandise" gorm:"foreignKey:merchant_id"`
}

func (MerchantDtl) TableName() string {
	return "merchants"
}

type MerchantUpdate struct {
	UserID      *uint                 `form:"user_id"`
	Name        *string               `form:"name"`
	Description *string               `form:"description"`
	Status      *string               `form:"status" gorm:"default:'active'"`
	Picture     *multipart.FileHeader `form:"picture"`
}
