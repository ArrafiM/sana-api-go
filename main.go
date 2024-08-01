package main

import (
	"sana-api/db"
	"sana-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	router := gin.Default()
	db.ConnecDatabase()

	router.Static("/public", "./public")

	routes.AlbumRoutes(router)
	routes.RoleRoute(router)
	routes.UserRoute(router)
	routes.AuthRoute(router)

	router.Run("localhost:8080")
}
