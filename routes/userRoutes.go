package routes

import (
	"sana-api/controllers"
	"sana-api/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	routeGroup := r.Group("/api/users")
	routeGroup.Use(middlewares.JwtAuthMiddleware())
	routeGroup.GET("/", controllers.Users)
	routeGroup.PUT("/", controllers.UploadFile)
	routeGroup.GET("/me", controllers.CurrentUser)
	routeGroup.DELETE("/:id", controllers.DeleteUser)
	changepass := r.Group("/api/users-changepass")
	changepass.Use(middlewares.JwtAuthMiddleware())
	changepass.PUT("/", controllers.ChangePass)
}
