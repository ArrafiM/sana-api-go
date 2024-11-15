package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func DeviceTokenRoutes(r *gin.Engine) {
	routeGroup := r.Group("/api/devicetokens")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("/", controllers.GetDeviceTokens)
	routeGroup.POST("/", controllers.StoreDeviceToken)
}
