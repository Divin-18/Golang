package websocket

import (
	"encoding/json"
	"log"
	"real-time-chat/internal/models"
	"real-time-chat/internal/repository"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	rooms      map[int]map[*Client]bool
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
}

type BroadcastMessage struct {
	RoomID  int
	Message []byte
}

func NewHub(messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		rooms:       make(map[int]map[*Client]bool),
		broadcast:   make(chan *BroadcastMessage, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.userRepo.UpdateOnlineStatus(client.UserID, true)
			h.mutex.Unlock()
			log.Printf("Client connected: %s (ID: %d)", client.Username, client.UserID)
			h.broadcastOnlineUsers()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Remove from all rooms
				for roomID := range h.rooms {
					delete(h.rooms[roomID], client)
				}
				h.userRepo.UpdateOnlineStatus(client.UserID, false)
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected: %s (ID: %d)", client.Username, client.UserID)
			h.broadcastOnlineUsers()

		case message := <-h.broadcast:
			h.mutex.RLock()
			if clients, ok := h.rooms[message.RoomID]; ok {
				for client := range clients {
					select {
					case client.send <- message.Message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
	log.Printf("User %s joined room %d", client.Username, roomID)

	// Notify room members
	notification := models.WSMessage{
		Type: "user_joined",
		Payload: map[string]interface{}{
			"room_id":  roomID,
			"user_id":  client.UserID,
			"username": client.Username,
		},
	}
	h.notifyRoom(roomID, notification, client)
}

func (h *Hub) LeaveRoom(client *Client, roomID int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if clients, ok := h.rooms[roomID]; ok {
		delete(clients, client)
		log.Printf("User %s left room %d", client.Username, roomID)

		// Notify room members
		notification := models.WSMessage{
			Type: "user_left",
			Payload: map[string]interface{}{
				"room_id":  roomID,
				"user_id":  client.UserID,
				"username": client.Username,
			},
		}
		h.notifyRoomLocked(roomID, notification, client)
	}
}

func (h *Hub) BroadcastToRoom(roomID int, message *models.Message) {
	wsMessage := models.WSMessage{
		Type:    "new_message",
		Payload: message,
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: data,
	}
}

func (h *Hub) BroadcastTyping(roomID int, userID int, username string, isTyping bool) {
	wsMessage := models.WSMessage{
		Type: "typing",
		Payload: models.TypingIndicator{
			RoomID:   roomID,
			Username: username,
			IsTyping: isTyping,
		},
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		return
	}

	h.mutex.RLock()
	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			if client.UserID != userID {
				select {
				case client.send <- data:
				default:
				}
			}
		}
	}
	h.mutex.RUnlock()
}

func (h *Hub) notifyRoom(roomID int, message models.WSMessage, exclude *Client) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			if client != exclude {
				select {
				case client.send <- data:
				default:
				}
			}
		}
	}
}

func (h *Hub) notifyRoomLocked(roomID int, message models.WSMessage, exclude *Client) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			if client != exclude {
				select {
				case client.send <- data:
				default:
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) broadcastOnlineUsers() {
	users, err := h.userRepo.GetOnlineUsers()
	if err != nil {
		log.Printf("Error getting online users: %v", err)
		return
	}

	wsMessage := models.WSMessage{
		Type:    "online_users",
		Payload: users,
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		return
	}

	h.mutex.RLock()
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
		}
	}
	h.mutex.RUnlock()
}
