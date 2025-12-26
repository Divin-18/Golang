package websocket

import (
	"encoding/json"
	"log"
	"real-time-chat/internal/models"
	"real-time-chat/internal/repository"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	UserID      int
	Username    string
	messageRepo *repository.MessageRepository
	roomRepo    *repository.RoomRepository
}

func NewClient(hub *Hub, conn *websocket.Conn, userID int, username string, 
	messageRepo *repository.MessageRepository, roomRepo *repository.RoomRepository) *Client {
	return &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		UserID:      userID,
		Username:    username,
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Batch pending messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(data []byte) {
	var wsMessage models.WSMessage
	if err := json.Unmarshal(data, &wsMessage); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch wsMessage.Type {
	case "join_room":
		c.handleJoinRoom(wsMessage.Payload)
	case "leave_room":
		c.handleLeaveRoom(wsMessage.Payload)
	case "send_message":
		c.handleSendMessage(wsMessage.Payload)
	case "typing":
		c.handleTyping(wsMessage.Payload)
	default:
		log.Printf("Unknown message type: %s", wsMessage.Type)
	}
}

func (c *Client) handleJoinRoom(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var joinRoom models.JoinRoom
	if err := json.Unmarshal(data, &joinRoom); err != nil {
		return
	}

	// Add user to room membership in database
	c.roomRepo.AddMember(joinRoom.RoomID, c.UserID)
	c.hub.JoinRoom(c, joinRoom.RoomID)
}

func (c *Client) handleLeaveRoom(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var leaveRoom models.LeaveRoom
	if err := json.Unmarshal(data, &leaveRoom); err != nil {
		return
	}

	c.hub.LeaveRoom(c, leaveRoom.RoomID)
}

func (c *Client) handleSendMessage(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var chatMessage models.ChatMessage
	if err := json.Unmarshal(data, &chatMessage); err != nil {
		return
	}

	// Save message to database
	message := &models.Message{
		RoomID:      chatMessage.RoomID,
		UserID:      c.UserID,
		Username:    c.Username,
		Content:     chatMessage.Content,
		MessageType: "text",
	}

	if err := c.messageRepo.Create(message); err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	// Broadcast to room
	c.hub.BroadcastToRoom(chatMessage.RoomID, message)
}

func (c *Client) handleTyping(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var typing models.TypingIndicator
	if err := json.Unmarshal(data, &typing); err != nil {
		return
	}

	c.hub.BroadcastTyping(typing.RoomID, c.UserID, c.Username, typing.IsTyping)
}
