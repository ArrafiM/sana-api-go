package models

import (
	"gorm.io/gorm"
)

type UserLocation struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Location string    `gorm:"type:geometry(POINT,4326)" json:"locations"`
	User     *User     `json:"user"`
	Merchant *Merchant `json:"merchant" gorm:"foreignKey:UserID;references:UserID"`
}

type CustomLocation struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
	Distance  *string `json:"distance"`
	UserMerchant
}

func (CustomLocation) TableName() string {
	return "user_locations"
}

type UserMerchant struct {
	User     *User     `json:"user"`
	Merchant *Merchant `json:"merchant" gorm:"foreignKey:UserID;references:UserID"`
}
