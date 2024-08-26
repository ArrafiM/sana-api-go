package models

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Merchandise struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Active      bool      `json:"active" gorm:"default:TRUE"`
	Picture     string    `json:"picture"`
	MerchantID  int       `json:"merchant_id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateTime"`
}

type GetMerchandiseImage struct {
	Merchandise
	Images *[]MerchandiseImages `json:"images" gorm:"foreignKey:MerchandiseID"`
}

func (GetMerchandiseImage) TableName() string {
	return "merchandises"
}

type MerchandiseCreate struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Price       string                `form:"price" binding:"required"`
	MerchantID  string                `form:"merchant_id" binding:"required"`
	Picture     *multipart.FileHeader `form:"picture" binding:"required"`
}

type MerchandiseImages struct {
	ID            uint   `json:"id" gorm:"primarykey"`
	Url           string `json:"url"`
	MerchandiseID uint   `json:"merchandise_id"`
	CreatedAt     time.Time
}

type MerchandiseUploadImages struct {
	MerchandiseID string                 `form:"merchandise_id" binding:"required"`
	Files         []multipart.FileHeader `form:"files[]" binding:"required"`
}
