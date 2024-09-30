package routes

import (
	"sana-api/controllers"

	"github.com/gin-gonic/gin"
)

func WebscoketRoute(r *gin.Engine) {
	//websoket
	r.GET("/ws", controllers.HandleWebSocket)
	go controllers.RunHub()
	// r.GET()
}
