package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func MerchantRoute(r *gin.Engine) {
	routeGroup := r.Group("/api/merchants")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("/", controllers.GetMerchants)
	routeGroup.POST("/", controllers.CreateMerchant)
	routeGroup.PUT("/:id", controllers.MerchantUpdate)
	routeGroup.POST("/uploadlanding", controllers.MerchantUploadLandingImage)
}
