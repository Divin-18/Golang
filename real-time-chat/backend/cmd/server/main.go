package main

import (
	"log"
	"real-time-chat/internal/config"
	"real-time-chat/internal/database"
	"real-time-chat/internal/handlers"
	"real-time-chat/internal/middleware"
	"real-time-chat/internal/repository"
	"real-time-chat/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL database")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Initialize WebSocket hub
	hub := websocket.NewHub(messageRepo, userRepo)
	go hub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	roomHandler := handlers.NewRoomHandler(roomRepo, messageRepo)
	wsHandler := handlers.NewWebSocketHandler(hub, messageRepo, roomRepo)

	// Setup Gin router
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User routes
		protected.GET("/me", authHandler.GetCurrentUser)

		// Room routes
		protected.GET("/rooms", roomHandler.GetRooms)
		protected.POST("/rooms", roomHandler.CreateRoom)
		protected.GET("/rooms/my", roomHandler.GetUserRooms)
		protected.GET("/rooms/:id", roomHandler.GetRoom)
		protected.POST("/rooms/:id/join", roomHandler.JoinRoom)
		protected.POST("/rooms/:id/leave", roomHandler.LeaveRoom)
		protected.GET("/rooms/:id/members", roomHandler.GetRoomMembers)
		protected.GET("/rooms/:id/messages", roomHandler.GetRoomMessages)
	}

	// WebSocket route (with auth)
	router.GET("/ws", middleware.AuthMiddleware(cfg.JWTSecret), wsHandler.HandleWebSocket)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
