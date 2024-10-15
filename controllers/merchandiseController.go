package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sana-api/models"
	"time"

	"reflect"
	"strconv"

	"sana-api/db"
	"sana-api/utils/token"

	"github.com/gin-gonic/gin"
)

func GetMerchandises(c *gin.Context) {
	var merchandise []models.GetMerchandiseImage
	db.CON.
		Preload("Images").
		Find(&merchandise)
	c.JSON(http.StatusOK, gin.H{"message": "All merchandise data", "data": merchandise})
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
	userId, _ := token.ExtractTokenID(c)
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
	broadCastMerchant(userId)
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

func MerchandiseUpdate(c *gin.Context) {
	userId, _ := token.ExtractTokenID(c)
	var merchandise models.MerchandiseUpdate

	id := c.Param("id")

	// Bind multipart/form-data to struct
	if err := c.ShouldBind(&merchandise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch existing merchant
	var existingMerchandise models.Merchandise
	if err := db.CON.First(&existingMerchandise, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchandise not found"})
		return
	}

	// Handle picture upload if it exists
	file, err := c.FormFile("picture")
	if err == nil {
		// File "picture" exists in the request
		if existingMerchandise.Picture != "" {
			RemoveFile(existingMerchandise.Picture)
		}
		path := "merchandiseimages"
		url := fileUrl(file, path)
		c.SaveUploadedFile(file, "public/"+url)
		existingMerchandise.Picture = url
		merchandise.Picture = nil
	}
	// Use reflection to iterate over fields and update only non-nil fields
	merchandiseValue := reflect.ValueOf(merchandise)
	merchandiseType := merchandiseValue.Type()
	existingMerchandiseValue := reflect.ValueOf(&existingMerchandise).Elem()

	for i := 0; i < merchandiseValue.NumField(); i++ {
		field := merchandiseValue.Field(i)
		if !field.IsNil() {
			existingField := existingMerchandiseValue.FieldByName(merchandiseType.Field(i).Name)
			if existingField.Kind() == reflect.Ptr {
				existingField.Set(field)
			} else {
				existingField.Set(reflect.Indirect(field))
			}
		}
	}

	// Save updated merchandise
	if err := db.CON.Save(&existingMerchandise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update merchandise"})
		return
	}

	broadCastMerchant(userId)

	c.JSON(http.StatusOK, gin.H{"message": "merchandise updated", "data": existingMerchandise})
}

func MerchandiseDelete(c *gin.Context) {
	id := c.Param("id")
	var merchandise models.Merchandise

	if err := db.CON.Where("id = ?", id).First(&merchandise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchandise not found", "data": false})
		return
	}
	//remove picture
	if merchandise.Picture != "" {
		RemoveFile(merchandise.Picture)
	}
	//delete permanently
	db.CON.Unscoped().Delete(&merchandise)
	c.JSON(http.StatusNotFound, gin.H{"error": "Merchandise deleted", "data": true})
}
