package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sana-api/models"
	"time"

	"strconv"

	"sana-api/db"

	"github.com/gin-gonic/gin"
)

func GetMerchandises(c *gin.Context) {
	var merchant []models.GetMerchandiseImage
	db.CON.
		Preload("Images").
		Find(&merchant)
	c.JSON(http.StatusOK, gin.H{"message": "All merchandise data", "data": merchant})
}

func GetMerchandiseId(c *gin.Context) {
	id := c.Param("id")
	// Fetch existing merchandise
	var existingMerchandise models.GetMerchandiseImage
	if err := db.CON.
		Preload("Images").
		First(&existingMerchandise, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchandise not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All merchandise data", "data": existingMerchandise})
}

func CreateMerchandise(c *gin.Context) {
	var payload models.MerchandiseCreate
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, _ := c.FormFile("picture")
	path := "merchandisepic"
	url := fileUrl(file, path)

	c.SaveUploadedFile(file, "public/"+url)

	price, _ := strconv.Atoi(payload.Price)
	merchantID, _ := strconv.Atoi(payload.MerchantID)

	merchandise := models.Merchandise{
		MerchantID:  merchantID,
		Price:       price,
		Name:        payload.Name,
		Description: payload.Description,
		Picture:     url,
	}

	db.CON.Create(&merchandise)
	c.JSON(http.StatusCreated, gin.H{"message": "Merchandise created", "data": merchandise})
}

func MerchandiseUploadImages(c *gin.Context) {
	var payload models.MerchandiseUploadImages
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	var images []models.MerchandiseImages
	merchandiseId, err := strconv.Atoi(payload.MerchandiseID)

	// Fetch existing merchandise
	var existingMerchandise models.Merchandise
	if err := db.CON.First(&existingMerchandise, merchandiseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchandise not found"})
		return
	}

	for _, file := range files {
		file.Filename = fmt.Sprint(time.Now().UnixNano()) + "-" + file.Filename
		log.Println(file.Filename)
		path := "merchandiseimages"
		url := path + "/" + file.Filename
		c.SaveUploadedFile(file, "public/"+url)
		if err != nil {
			continue
		}
		images = append(images, models.MerchandiseImages{
			MerchandiseID: uint(merchandiseId),
			Url:           url,
		})
	}
	db.CON.CreateInBatches(&images, len(images))

	c.JSON(http.StatusCreated, gin.H{"message": "merchindise multiple image uploaded", "data": images})
}
