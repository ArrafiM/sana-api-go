package models

import (
	"gorm.io/gorm"
)

type UserLocation struct {
	gorm.Model
	User_id  uint   `json:"user_id"`
	Location string `gorm:"type:geometry(POINT,4326)" json:"locations"`
	User     *User  `json:"user"`
}

type CustomLocation struct {
	gorm.Model
	User_id   uint `json:"user_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
