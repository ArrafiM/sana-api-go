package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {
	routeGroup := r.Group("/api/chats")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("", controllers.AllChat)
	routeGroup.POST("", controllers.StoreChat)
	r.GET("/api/chatrooms", controllers.ChatRoom)
}
