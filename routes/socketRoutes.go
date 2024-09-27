package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SocketRoute(r *gin.Engine) {
	//websoket
	r.GET("/ws1", socket)
	// r.GET()
}

var upgrader1 = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(c *http.Request) bool {
		return true
	},
}

func socket(c *gin.Context) {
	conn, err := upgrader1.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	log.Println("Client Connected")

	// for {
	// 	conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket Sana!"))
	// 	// reader(conn)
	// 	time.Sleep(5 * time.Second)
	// }
	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		fmt.Println("message from client: ", string(p))
		fmt.Println(messageType)

		message := []byte("your message recived!!")

		// for {
		// time.Sleep(5 * time.Second)
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			// break
		}
		// }

	}
}
