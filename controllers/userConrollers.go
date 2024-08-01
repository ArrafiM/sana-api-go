package controllers

import (
	"net/http"
	"sana-api/db"
	"sana-api/models"
	"sana-api/utils/token"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Users(c *gin.Context) {
	showAll := c.Query("all")
	var users []models.User
	if showAll == "1" {
		db.CON.Unscoped().Find(&users)
	} else {
		db.CON.Find(&users)
	}
	c.JSON(http.StatusOK, gin.H{"message": "all users", "data": users})
}

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var u models.User

	db.CON.First(&u, user_id)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	var u models.User

	if err := db.CON.First(&u, userId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user notfound"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	//soft delete user
	db.CON.Delete(&u)
	c.JSON(http.StatusOK, gin.H{"message": "user soft deleted", "id": userId})
}
