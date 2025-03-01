package controllers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sana-api/db"
	"sana-api/helpers"
	"sana-api/models"
	"sana-api/utils/token"
	"strconv"
	"strings"
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
	user := c.Query("user")
	image := c.Query("image")
	item := c.Query("item")
	// Fetch existing merchant
	var existingMerchant models.MerchantDtl
	getMerchant := db.CON.Where("id = ?", id)
	if user == "true" {
		getMerchant.Preload("User")
	}
	if item == "true" {
		getMerchant.Preload("Merchandise")
	}
	if image == "true" {
		getMerchant.Preload("LandingImages")
	}
	if err := getMerchant.
		First(&existingMerchant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "detail merchant data", "data": existingMerchant})
}

func GetMyMerchant(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)
	cek := c.Query("cek")
	image := c.Query("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	println("user_id", user_id)
	var merchant models.MerchantDtl
	getMerchant := db.CON.Where("user_id = ?", user_id)
	if cek != "true" {
		getMerchant.
			Preload("User").
			Preload("Merchandise")
	}
	if image == "true" {
		getMerchant.Preload("LandingImages")
	}
	if err := getMerchant.
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

	Color := ""
	if payload.Color != nil {
		Color = *payload.Color
	} else {
		Color = helpers.GenerateHexColor()
	}

	merchant := models.Merchant{
		UserID:      userId,
		Name:        payload.Name,
		Description: payload.Description,
		Picture:     url,
		Color:       Color,
	}
	if payload.Status != nil {
		merchant.Status = *payload.Status
	}

	db.CON.Create(&merchant)
	// broadCastMerchant(userId)
	c.JSON(http.StatusCreated, gin.H{"message": "Merchant created", "data": merchant})
}

func strToArray(str string) []int {
	var intArr []int
	if str != "" {
		str = strings.Trim(str, "[]")
		strArr := strings.Split(str, ",")
		// Konversi setiap elemen string menjadi integer
		for _, s := range strArr {
			num, _ := strconv.Atoi(strings.TrimSpace(s)) // Hapus spasi dan konversi
			intArr = append(intArr, num)
		}
		return intArr
	}
	return intArr
}

func MerchantUploadLandingImage(c *gin.Context) {
	// userId, _ := token.ExtractTokenID(c)
	// Multipart form
	form, _ := c.MultipartForm()
	merchantId, _ := strconv.Atoi(form.Value["merchant_id"][0])
	removeId := form.Value["remove_id"]
	if removeId != nil {
		val := strToArray(removeId[0])
		if len(val) > 0 {
			var merchantImage []models.MerchantLandingImage
			db.CON.Where("merchant_id = ? and id IN ?", merchantId, val).Find(&merchantImage)
			//remove selected old image
			for _, image := range merchantImage {
				if image.Url != "" {
					RemoveFile(image.Url)
					db.CON.Delete(&image)
				}
			}
		}
	}
	//upload new image
	files := form.File["files[]"]
	var images []models.MerchantLandingImage
	for _, file := range files {
		file.Filename = fmt.Sprint(time.Now().UnixNano()) + "-" + file.Filename
		log.Println(file.Filename)
		path := "merchantlanding"
		url := fileUrl(file, path)
		c.SaveUploadedFile(file, "public/"+url)
		images = append(images, models.MerchantLandingImage{
			MerchantId: uint(merchantId),
			Url:        url,
		})
	}
	db.CON.CreateInBatches(&images, len(images))
	// broadCastMerchant(userId)
	c.JSON(http.StatusCreated, gin.H{"message": "multiple file uploaded", "data": images, "remove_id": removeId})
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

	// broadCastMerchant(existingMerchant.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "merchant updated", "data": existingMerchant})
}

// func broadCastMerchant(userId uint) {
// 	msg := Message{
// 		SenderID:   fmt.Sprintf("user%s", strconv.Itoa(int(userId))),
// 		ReceiverID: fmt.Sprintf("user%s", strconv.Itoa(int(userId))),
// 		Content:    fmt.Sprintf("myMerchant%s", strconv.Itoa(int(userId))),
// 	}
// 	BroadcastMessage(msg)
// }
