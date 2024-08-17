package controllers

import (
	"net/http"
	"sana-api/db"
	"sana-api/models"
	"sana-api/utils/token"

	"github.com/gin-gonic/gin"
)

func GetMerchants(c *gin.Context) {
	var merchant []models.MerchantUser
	db.CON.Preload("User").Find(&merchant)
	c.JSON(http.StatusOK, gin.H{"message": "All merchant data", "data": merchant})
}

func CreateMerchant(c *gin.Context) {
	var payload models.MerchantCreate
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchant := models.Merchant{
		UserID:      userId,
		Name:        payload.Name,
		Description: payload.Description,
	}

	db.CON.Create(&merchant)
	c.JSON(http.StatusCreated, gin.H{"message": "Merchant created", "data": merchant})
}
