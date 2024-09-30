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

	router.GET("/", index)

	routes.IndexRoutes(router)

	//websocket
	// routes.SocketRoute(router)
	router.Run("localhost:8080")

}

func index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"Welcome": "to sana API!"})
}
