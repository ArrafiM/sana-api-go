package controllers

import (
	"net/http"
	"sana-api/db"
	"sana-api/models"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var roles []models.RoleGet
	db.CON.Find(&roles)
	c.JSON(http.StatusOK, gin.H{"message": "all items", "data": roles})
}

func Store(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	db.CON.Create(&role)
	c.JSON(http.StatusOK, gin.H{"message": "item created", "data": role})
}
