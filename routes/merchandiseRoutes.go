package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func MerchandiseRoute(r *gin.Engine) {
	routeGroup := r.Group("/api/merchandise")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("", controllers.GetMerchandises)
	routeGroup.GET("/:id", controllers.GetMerchandiseId)
	routeGroup.POST("", controllers.CreateMerchandise)
	routeGroup.PUT("/:id", controllers.MerchandiseUpdate)
	routeGroup.DELETE("/:id", controllers.MerchandiseDelete)
	routeGroup.POST("/uploadimages", controllers.MerchandiseUploadImages)
}
