package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func FeedbackRoutes(r *gin.Engine) {
	routeGroup := r.Group("/api/feedback")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("", controllers.GetFeedback)
	routeGroup.POST("", controllers.StoreFeedback)
}
