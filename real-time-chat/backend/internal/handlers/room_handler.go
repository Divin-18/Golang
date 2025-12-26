package handlers

import (
	"net/http"
	"real-time-chat/internal/models"
	"real-time-chat/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomRepo    *repository.RoomRepository
	messageRepo *repository.MessageRepository
}

func NewRoomHandler(roomRepo *repository.RoomRepository, messageRepo *repository.MessageRepository) *RoomHandler {
	return &RoomHandler{
		roomRepo:    roomRepo,
		messageRepo: messageRepo,
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req models.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	room := &models.Room{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID.(int),
		IsPrivate:   req.IsPrivate,
	}

	if err := h.roomRepo.Create(room); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create room"})
		return
	}

	// Add creator as member
	h.roomRepo.AddMember(room.ID, userID.(int))

	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.roomRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch rooms"})
		return
	}

	if rooms == nil {
		rooms = []*models.Room{}
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid room ID"})
		return
	}

	room, err := h.roomRepo.GetByID(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) GetUserRooms(c *gin.Context) {
	userID, _ := c.Get("userID")

	rooms, err := h.roomRepo.GetUserRooms(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch rooms"})
		return
	}

	if rooms == nil {
		rooms = []*models.Room{}
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) JoinRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid room ID"})
		return
	}

	userID, _ := c.Get("userID")

	if err := h.roomRepo.AddMember(roomID, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to join room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined room successfully"})
}

func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid room ID"})
		return
	}

	userID, _ := c.Get("userID")

	if err := h.roomRepo.RemoveMember(roomID, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to leave room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left room successfully"})
}

func (h *RoomHandler) GetRoomMembers(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid room ID"})
		return
	}

	members, err := h.roomRepo.GetMembers(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch members"})
		return
	}

	if members == nil {
		members = []*models.User{}
	}

	c.JSON(http.StatusOK, members)
}

func (h *RoomHandler) GetRoomMessages(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid room ID"})
		return
	}

	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	messages, err := h.messageRepo.GetByRoomID(roomID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch messages"})
		return
	}

	if messages == nil {
		messages = []*models.Message{}
	}

	c.JSON(http.StatusOK, messages)
}
