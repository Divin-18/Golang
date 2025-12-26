# Real-Time Chat Application

A modern real-time chat application built with React (frontend), Go (backend), and PostgreSQL (database).

## Features

- ✨ Real-time messaging with WebSocket
- 🔐 JWT-based authentication
- 🏠 Multiple chat rooms
- 👥 Online user presence
- ⌨️ Typing indicators
- 📱 Responsive design
- 🎨 Modern glassmorphism UI

## Project Structure

```
real-time-chat/
├── backend/                 # Go backend
│   ├── cmd/
│   │   └── server/         # Main entry point
│   └── internal/
│       ├── auth/           # JWT authentication
│       ├── config/         # Configuration
│       ├── database/       # PostgreSQL connection
│       ├── handlers/       # HTTP handlers
│       ├── middleware/     # Auth middleware
│       ├── models/         # Data models
│       ├── repository/     # Database operations
│       └── websocket/      # WebSocket hub & client
└── frontend/               # React frontend
    └── src/
        ├── components/     # React components
        ├── contexts/       # Auth & WebSocket contexts
        ├── pages/          # Page components
        └── services/       # API service
```

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

## Setup

### Database

1. Create a PostgreSQL database:

```sql
CREATE DATABASE chat_db;
```

### Backend

1. Navigate to backend:

```bash
cd backend
```

2. Copy environment file and update values:

```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Install dependencies:

```bash
go mod tidy
```

4. Run the server:

```bash
go run cmd/server/main.go
```

### Frontend

1. Navigate to frontend:

```bash
cd frontend
```

2. Install dependencies:

```bash
npm install
```

3. Run development server:

```bash
npm run dev
```

## API Endpoints

### Authentication

- `POST /api/register` - Register new user
- `POST /api/login` - Login user
- `GET /api/me` - Get current user

### Rooms

- `GET /api/rooms` - Get all public rooms
- `POST /api/rooms` - Create new room
- `GET /api/rooms/:id` - Get room by ID
- `POST /api/rooms/:id/join` - Join room
- `POST /api/rooms/:id/leave` - Leave room
- `GET /api/rooms/:id/messages` - Get room messages
- `GET /api/rooms/:id/members` - Get room members

### WebSocket

- `GET /ws?token=<JWT>` - WebSocket connection

## WebSocket Events

### Client to Server

- `join_room` - Join a chat room
- `leave_room` - Leave a chat room
- `send_message` - Send a message
- `typing` - Typing indicator

### Server to Client

- `new_message` - New message received
- `user_joined` - User joined room
- `user_left` - User left room
- `online_users` - Online users list
- `typing` - User typing status

## Tech Stack

### Backend

- Go with Gin framework
- Gorilla WebSocket
- PostgreSQL with pq driver
- JWT for authentication
- bcrypt for password hashing

### Frontend

- React 18 with Vite
- Context API for state management
- WebSocket for real-time communication
- CSS with modern design patterns
