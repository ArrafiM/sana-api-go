package routes

import (
	"sana-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {
	routeGroup := r.Group("/api/auth")
	routeGroup.POST("/register", controllers.Register)
	routeGroup.POST("/login", controllers.Login)
}
