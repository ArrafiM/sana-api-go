package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func RoleRoute(r *gin.Engine) {
	routeGroup := r.Group("/api/roles")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("/", controllers.Index)
	routeGroup.POST("/", controllers.Store)
}
