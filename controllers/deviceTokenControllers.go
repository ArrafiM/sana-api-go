package controllers

import (
	"net/http"
	"sana-api/db"
	"sana-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDeviceTokens(c *gin.Context) {
	userId := c.Query("userid")
	var deviceTokens []models.DeviceToken
	if userId != "" {
		// Find the first device token by user_id
		var deviceToken models.DeviceToken
		if err := db.CON.Where("user_id = ?", userId).First(&deviceToken).Error; err != nil {
			// Handle error if no record found or other issues
			c.JSON(http.StatusNotFound, gin.H{"message": "device token not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "device token found", "data": deviceToken})
		return
	}
	db.CON.Find(&deviceTokens)
	c.JSON(http.StatusOK, gin.H{"message": "all device token", "data": deviceTokens})
}

func StoreDeviceToken(c *gin.Context) {
	var tokenCreate models.DeviceTokenCreate

	if err := c.ShouldBindJSON(&tokenCreate); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if a user with the same email exists
	var token models.DeviceToken
	errCek := db.CON.Where("user_id = ? and token = ?", tokenCreate.UserID, tokenCreate.Token).First(&token).Error
	if errCek == nil {
		// If a token is found, return "already exist" message
		c.JSON(http.StatusConflict, gin.H{"message": "Token already exists"})
		return
	}
	err := db.CON.Where("user_id = ?", tokenCreate.UserID).First(&token).Error
	if err != nil {
		// If no record found, create a new one
		if err == gorm.ErrRecordNotFound {
			if err := db.CON.Create(&tokenCreate).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device token"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"message": "device token created", "data": tokenCreate})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// If record exists, update it with new data
	token.Token = tokenCreate.Token
	if err := db.CON.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device token updated", "data": token})
}
