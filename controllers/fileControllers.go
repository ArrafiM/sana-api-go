package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sana-api/db"
	"sana-api/models"
	"sana-api/utils/token"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadFile(c *gin.Context) {
	// single file
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, _ := c.FormFile("picture")
	file.Filename = fmt.Sprint(time.Now().UnixNano()) + "-" + file.Filename
	log.Println(file.Filename)
	var user models.User
	if err := db.CON.First(&user, user_id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user notfound"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	if user.Picture != "" {
		e := os.Remove("public/" + user.Picture)
		if e != nil {
			log.Fatal(e)
		}
	}

	path := "pictureuser/" + file.Filename
	c.SaveUploadedFile(file, "public/"+path)

	if db.CON.Model(&user).Where("id = ?", user_id).Update("picture", path).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user not updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": file.Filename})
}
