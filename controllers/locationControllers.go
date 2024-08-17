package controllers

import (
	"fmt"
	"net/http"
	"sana-api/models"

	"sana-api/db"
	"sana-api/utils/token"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm/clause"
	// "gorm.io/gorm"
	// "log"
)

type CreateLocationInput struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"Long" binding:"required"`
}

func StoreLocation(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input CreateLocationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	location := fmt.Sprintf("POINT(%f %f)", input.Long, input.Lat)
	var UserLocation models.UserLocation
	//upsert
	// db.CON.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "user_id"}},
	// 	DoUpdates: clause.Assignments(map[string]interface{}{"location": location}),
	// }).Create(&UserLocation)

	db.CON.Where("user_id = ?", user_id).Delete(&UserLocation)
	createLocation := models.UserLocation{
		User_id:  user_id,
		Location: location,
	}
	db.CON.Create(&createLocation)

	c.JSON(http.StatusOK, gin.H{"message": "user location added", "data": createLocation})
}

func GetUserLocations(c *gin.Context) {
	my := c.Query("my")
	if my == "1" {
		GetMyLocation(c)
		return
	}
	var locations []models.UserLocation

	db.CON.Model(&locations).Where("tracked", true).Preload("User").Find(&locations)
	c.JSON(http.StatusOK, gin.H{"message": "all user location", "data": locations})
}

func GetMyLocation(c *gin.Context) {
	// userId, err := token.ExtractTokenID(c)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	var location models.CustomLocation
	db.CON.Raw("SELECT id, user_id, ST_X(location::geometry) AS longitude, ST_Y(location::geometry) AS latitude FROM user_locations where id = 5").Scan(&location)
	// db.CON.Where("user_id = ?", userId).First(&location)

	c.JSON(http.StatusOK, gin.H{"message": "your location", "data": location})
}

type NearestModel struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Color     string  `json:"color"`
	Title     string  `json:"title"`
}

var nearestData = [4]NearestModel{
	{
		ID: 1, Latitude: -6.9092216, Longitude: 107.593377,
		Color: "#FF5733", Title: "Strowberry",
	},
	{
		ID: 2, Latitude: -6.9100445, Longitude: 107.5939065,
		Color: "#7e4ac2", Title: "Mie Pedas",
	},
	{
		ID: 3, Latitude: -6.9090957, Longitude: 107.5937257,
		Color: "#e359cc", Title: "Soto",
	},
	{
		ID: 4, Latitude: -6.9097981, Longitude: 107.5936367,
		Color: "#e39c59", Title: "Nasi Lemak",
	},
}

func NearestLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your nearest location", "data": nearestData})
}
