package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func LocationRoutes(r *gin.Engine) {
	routeGroup := r.Group("/api/locations")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.POST("/", controllers.StoreLocation)
	routeGroup.GET("/", controllers.GetUserLocations)
	routeGroup.GET("/nearest", controllers.GetNearestPoint)
	routeGroup.GET("/new", controllers.NewLocation)
}
