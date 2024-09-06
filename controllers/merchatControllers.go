package controllers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sana-api/db"
	"sana-api/models"
	"sana-api/utils/token"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMerchants(c *gin.Context) {
	var merchant []models.MerchantUser
	db.CON.Preload("User").
		Preload("LandingImages").
		Find(&merchant)
	c.JSON(http.StatusOK, gin.H{"message": "All merchant data", "data": merchant})
}

func GetMerchantId(c *gin.Context) {
	id := c.Param("id")
	// Fetch existing merchant
	var existingMerchant models.MerchantDtl
	if err := db.CON.
		Preload("User").
		Preload("LandingImages").
		Preload("Merchandise").
		First(&existingMerchant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "detail merchant data", "data": existingMerchant})
}

func GetMyMerchant(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	println("user_id", user_id)
	var merchant models.MerchantDtl
	if err := db.CON.Where("user_id = ?", user_id).
		Preload("User").
		Preload("LandingImages").
		Preload("Merchandise").
		First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Your merchant data", "data": merchant})
}

func CreateMerchant(c *gin.Context) {
	var payload models.MerchantCreate
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf(payload.Name)
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file := payload.Picture
	path := "merchantpic"
	url := fileUrl(file, path)

	c.SaveUploadedFile(file, "public/"+url)

	merchant := models.Merchant{
		UserID:      userId,
		Name:        payload.Name,
		Description: payload.Description,
		Picture:     url,
	}

	db.CON.Create(&merchant)
	c.JSON(http.StatusCreated, gin.H{"message": "Merchant created", "data": merchant})
}

func MerchantUploadLandingImage(c *gin.Context) {
	// var payload models.MerchantLandingImageCreate
	// if err := c.ShouldBindJSON(&payload); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	var images []models.MerchantLandingImage
	for _, file := range files {
		file.Filename = fmt.Sprint(time.Now().UnixNano()) + "-" + file.Filename
		log.Println(file.Filename)
		path := "merchantlanding"
		url := path + "/" + file.Filename
		c.SaveUploadedFile(file, "public/"+url)
		merchantId, err := strconv.Atoi(form.Value["merchant_id"][0])
		if err != nil {
			continue
		}
		images = append(images, models.MerchantLandingImage{
			MerchantId: uint(merchantId),
			Url:        url,
		})
	}
	db.CON.CreateInBatches(&images, len(images))

	c.JSON(http.StatusCreated, gin.H{"message": "multiple file uploaded", "data": images})
}

func MerchantUpdate(c *gin.Context) {
	var merchant models.MerchantUpdate

	id := c.Param("id")

	// Bind multipart/form-data to struct
	if err := c.ShouldBind(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch existing merchant
	var existingMerchant models.Merchant
	if err := db.CON.First(&existingMerchant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	// Handle picture upload if it exists
	file, err := c.FormFile("picture")
	if err == nil {
		// File "picture" exists in the request
		if existingMerchant.Picture != "" {
			RemoveFile(existingMerchant.Picture)
		}
		path := "merchantpic"
		url := fileUrl(file, path)
		c.SaveUploadedFile(file, "public/"+url)
		existingMerchant.Picture = url
		merchant.Picture = nil
	}
	// Use reflection to iterate over fields and update only non-nil fields
	merchantValue := reflect.ValueOf(merchant)
	merchantType := merchantValue.Type()
	existingMerchantValue := reflect.ValueOf(&existingMerchant).Elem()

	for i := 0; i < merchantValue.NumField(); i++ {
		field := merchantValue.Field(i)
		if !field.IsNil() {
			existingField := existingMerchantValue.FieldByName(merchantType.Field(i).Name)
			if existingField.Kind() == reflect.Ptr {
				existingField.Set(field)
			} else {
				existingField.Set(reflect.Indirect(field))
			}
		}
	}

	// Save updated merchant
	if err := db.CON.Save(&existingMerchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update merchant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "merchant updated", "data": existingMerchant})
}
