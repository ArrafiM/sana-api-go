package controllers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sana-api/db"
	"sana-api/models"
	"sana-api/utils/token"
	"time"
	"strings"
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
		RemoveFile(user.Picture)
	}

	file, _ := c.FormFile("picture")

	path := "pictureuser"

	url := fileUrl(file, path)
	c.SaveUploadedFile(file, "public/"+url)

	if db.CON.Model(&user).Where("id = ?", user_id).Update("picture", url).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user not updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": url})
}

func fileUrl(file *multipart.FileHeader, path string) string {
	file.Filename = fmt.Sprint(time.Now().UnixNano()) + "-" + file.Filename
	log.Println(file.Filename)
	url := strings.ReplaceAll(path + "/" + file.Filename, " ", "")

	return url
}

func RemoveFile(filename string) int {
	e := os.Remove("public/" + filename)
	if e != nil {
		log.Fatal(e)
		return 0
	}
	return 1
}
