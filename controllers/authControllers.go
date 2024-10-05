package controllers

import (
	"net/http"
	"sana-api/db"
	"sana-api/models"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	ConfimPassword string `json:"confirm_password" binding:"required"`
	// Phone          string `json:"phone" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Password != input.ConfimPassword {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password not match"})
		return
	}

	u := models.User{}

	if err := db.CON.Where("email ILIKE ?", "%"+input.Email+"%").First(&u).Error; err == nil {
		println(u.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email telah terdaftar"})
		return
	}

	u.Email = input.Email
	u.Password = input.Password
	// u.Phone = input.Phone
	u.Name = input.Name
	u.Role_id = 2
	u.BeforeSave()
	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success", "data": u})
}

type LoginInput struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Password = input.Password

	loginRes, err := models.LoginCheck(u.Email, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": loginRes.Token, "user_id": loginRes.UserId})

}
