package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sana-api/models"

	"sana-api/db"
	"sana-api/utils/token"

	"time"

	// "gorm.io/gorm"

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
		UserID:   user_id,
		Location: location,
	}
	db.CON.Create(&createLocation)

	c.JSON(http.StatusOK, gin.H{"message": "user location added", "data": createLocation})
}

func GetUserLocations(c *gin.Context) {
	// my := c.Query("my")
	// if my == "1" {
	// 	GetMyLocation(c)
	// 	return
	// }
	var locations []models.UserLocation

	db.CON.Model(&locations).
		// Where("tracked", true).
		Preload("Merchant").
		Preload("User").
		Find(&locations)
	c.JSON(http.StatusOK, gin.H{"message": "all user location", "data": locations})
}

func GetNearestPoint(c *gin.Context) {
	userId, _ := token.ExtractTokenID(c)
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	radius := c.Query("radius")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	merchandise := c.DefaultQuery("merchandise", "false") == "true"
	itemName := c.Query("itemname")
	excludeMy := c.DefaultQuery("excludemy", "false") == "true"
	var location []models.CustomLocation
	distance := fmt.Sprintf("ST_Distance(location, ST_SetSRID(ST_MakePoint(%s, %s), 4326)::geography)/1000 AS distance",
		longitude, latitude)
	query := db.CON.
		Scopes(db.Paginate(page, pageSize)).
		Select("id", "user_id",
			"ST_X(location::geography::geometry) AS longitude",
			"ST_Y(location::geography::geometry) AS latitude", distance).
		Where(`ST_DWithin(
				location,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
				?
			)`, longitude, latitude, radius).
		Order("distance ASC").
		Preload("User").
		Preload("Merchant")
	if excludeMy {
		query.Where("user_id != ?", userId)
	}
	// query.Joins("merchants ON merchants.user_id = user_locations.user_id")
	// if itemName != ""{
	// }
	query.Find(&location)
	if merchandise {
		//merchandise data
		for i, val := range location {
			merchant := val.Merchant
			merchId := merchant.ID
			var itemData []models.Merchandise
			query := db.CON.Where("merchant_id = ?", merchId)
			if itemName != "" {
				query.Where("name ilike ?", "%"+itemName+"%")
			}
			query.Order("id DESC").
				Limit(3).Find(&itemData)
			location[i].Merchant.Merchandise = &itemData
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "your location", "data": location})
}

type NearestModel struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Color     string  `json:"color"`
	Title     string  `json:"title"`
}

func NewLocation(c *gin.Context) {
	lat := c.Query("lat")
	long := c.Query("long")
	user_id, _ := token.ExtractTokenID(c)
	location := fmt.Sprintf("POINT(%s %s)", long, lat)
	var UserLocation models.UserLocation

	db.CON.Where("user_id = ?", user_id).Delete(&UserLocation)
	createLocation := models.UserLocation{
		UserID:   user_id,
		Location: location,
	}
	db.CON.Create(&createLocation)

	c.JSON(http.StatusOK, gin.H{"message": "user location added", "data": createLocation})
}

func broadCastLocation(userId string) {
	msg := Message{
		SenderID:   userId,
		ReceiverID: userId,
		Content:    fmt.Sprintf("postMyLocationuser%s", userId),
	}
	BroadcastMessage(msg)
}

func postLocation(msg Message) {
	senderId := msg.SenderID
	if msg.Location != nil {
		log.Printf("location msg socket: %s", *msg.Location)
	}
	log.Printf("job started delay 1 menit userId: %s", senderId)
	time.Sleep(1 * time.Minute)
	log.Printf("bakcground job id: %s, finish", senderId)
	broadCastLocation(senderId)
}

func bakcgroundLocation(msg Message) {
	fmt.Println("Starting background job with delay...")
	go func() { postLocation(msg) }()
}
