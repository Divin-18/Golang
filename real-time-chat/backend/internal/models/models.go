package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    string    `json:"avatar_url"`
	IsOnline     bool      `json:"is_online"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Room struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int       `json:"created_by"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
}

type RoomMember struct {
	ID       int       `json:"id"`
	RoomID   int       `json:"room_id"`
	UserID   int       `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

type Message struct {
	ID          int       `json:"id"`
	RoomID      int       `json:"room_id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ChatMessage struct {
	RoomID  int    `json:"room_id"`
	Content string `json:"content"`
}

type JoinRoom struct {
	RoomID int `json:"room_id"`
}

type LeaveRoom struct {
	RoomID int `json:"room_id"`
}

type TypingIndicator struct {
	RoomID   int    `json:"room_id"`
	Username string `json:"username"`
	IsTyping bool   `json:"is_typing"`
}

// API Request/Response types
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
