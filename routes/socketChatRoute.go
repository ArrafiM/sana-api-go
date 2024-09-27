package routes

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SocketChatRoute(r *gin.Engine) {
	//websoket
	r.GET("/ws", handleWebSocket)
	// r.GET()
}

// Menggunakan upgrader untuk meng-upgrade HTTP connection menjadi WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client struct untuk memegang informasi WebSocket connection dan ID pengguna
type Client struct {
	conn *websocket.Conn
	send chan []byte
	id   string
}

// Hub struct untuk mengelola semua client yang terhubung dan map ID pengguna ke client
type Hub struct {
	clients    map[string]*Client
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

// Message struct untuk mewakili pesan yang dikirim
type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

var hub = Hub{
	clients:    make(map[string]*Client),
	broadcast:  make(chan Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Handle WebSocket connection
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Ambil ID pengguna dari query parameter (atau bisa dari header/token)
	userID := c.Request.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("user_id is required")
		conn.Close()
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte),
		id:   userID,
	}
	hub.register <- client

	go client.readPump()
	go client.writePump()
}

// ReadPump untuk membaca pesan dari client
func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg.SenderID = c.id

		// Simpan pesan ke database
		// saveMessage(msg)

		// Kirim pesan ke client tujuan
		hub.broadcast <- msg
	}
}

// WritePump untuk menulis pesan kepada client
func (c *Client) writePump() {
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	c.conn.Close()
}

// Menjalankan hub untuk mengelola client dan pesan
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Kirim pesan ke penerima tertentu
			if receiver, ok := h.clients[message.ReceiverID]; ok {
				select {
				case receiver.send <- []byte(message.Content):
				default:
					close(receiver.send)
					delete(h.clients, message.ReceiverID)
				}
			}
		}
	}
}

// Simpan pesan ke PostgreSQL
// func saveMessage(msg Message) {
// 	_, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, content) VALUES ($1, $2, $3)", msg.SenderID, msg.ReceiverID, msg.Content)
// 	if err != nil {
// 		log.Println("Error saving message:", err)
// 	}
// }