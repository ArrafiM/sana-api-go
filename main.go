package main

import (
	"sana-api/db"
	"sana-api/routes"

	// "time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	router := gin.Default()
	db.ConnecDatabase()

	router.Static("/public", "./public")

	router.GET("/api", index)

	routes.IndexRoutes(router)

	//websocket
	// routes.SocketRoute(router)
	// router.Run("192.168.1.2:8080")
	// router.Run("192.168.18.32:8080")
	// router.Run("192.168.19.108:8080")
	// router.Run("172.20.10.3:8080")
	// router.Run("192.168.19.14:8080")
	router.Run("localhost:8080")

}

func index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"Welcome": "to sana API!"})
}
