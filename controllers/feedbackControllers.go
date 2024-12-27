package controllers

import (
	"net/http"
	"sana-api/db"
	"sana-api/models"

	"sana-api/utils/token"

	"github.com/gin-gonic/gin"
)

func GetFeedback(c *gin.Context) {
	var feedbacks []models.Feedback
	page := c.Query("page")
	pageSize := c.Query("page_size")
	db.CON.Scopes(db.Paginate(page, pageSize)).
		Order("id desc").
		Find(&feedbacks)
	c.JSON(http.StatusOK, gin.H{"message": "all feedback", "data": feedbacks})
}

func StoreFeedback(c *gin.Context) {
	var payload models.FeedbackCreate
	userId, _ := token.ExtractTokenID(c)

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	feedback := models.Feedback{
		UserID:     userId,
		Email:      payload.Email,
		Properties: payload.Properties,
	}

	db.CON.Create(&feedback)
	c.JSON(http.StatusOK, gin.H{"message": "Thank you for your feedback!", "data": feedback})
}
