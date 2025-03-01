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
	routeGroup.PUT("/update", controllers.UploadFile)
	routeGroup.GET("/me", controllers.CurrentUser)
	routeGroup.DELETE("/:id", controllers.DeleteUser)
	routeGroup.PUT("/changepass", controllers.ChangePass)
}
