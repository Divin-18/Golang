package handlers

import (
	"log"
	"net/http"
	"real-time-chat/internal/repository"
	"real-time-chat/internal/websocket"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WebSocketHandler struct {
	hub         *websocket.Hub
	messageRepo *repository.MessageRepository
	roomRepo    *repository.RoomRepository
}

func NewWebSocketHandler(hub *websocket.Hub, messageRepo *repository.MessageRepository, roomRepo *repository.RoomRepository) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username, _ := c.Get("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := websocket.NewClient(
		h.hub,
		conn,
		userID.(int),
		username.(string),
		h.messageRepo,
		h.roomRepo,
	)

	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (h *WebSocketHandler) GetHub() *websocket.Hub {
	return h.hub
}
