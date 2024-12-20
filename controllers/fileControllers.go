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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadFile(c *gin.Context) {
	// single file
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.CON.First(&user, userId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user notfound"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	file, _ := c.FormFile("picture")

	if user.Picture != "" && file != nil {
		RemoveFile(user.Picture)
	}

	if file != nil {
		path := "pictureuser"

		url := fileUrl(file, path)
		c.SaveUploadedFile(file, "public/"+url)

		if db.CON.Model(&user).Where("id = ?", userId).Update("picture", url).RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "user not updated"})
			return
		}
	}

	name := c.PostForm("name")
	if name != "" {
		user.Name = name
		db.CON.Save(&user)
	}
	broadCastMerchant(userId)
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

func fileUrl(file *multipart.FileHeader, path string) string {
	file.Filename = fmt.Sprint(time.Now().UnixNano()) + "-" + file.Filename
	log.Println(file.Filename)
	url := strings.ReplaceAll(path+"/"+file.Filename, " ", "")

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
