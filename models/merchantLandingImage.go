package models

import (
	"time"
)

type MerchantLandingImage struct {
	ID         uint   `gorm:"primarykey"`
	MerchantId uint   `json:"merchant_id" binding:"required"`
	Url        string `json:"url" binding:"required"`
	CreatedAt  time.Time
}

func (MerchantLandingImage) TableName() string {
	return "merchant_landingimages"
}

type MerchantLandingImageCreate struct {
	MerchantId uint   `form:"merchant_id" json:"merchant_id" binding:"required"`
	Url        string `form:"files" json:"url" binding:"required"`
}
