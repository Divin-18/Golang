import { useState, useEffect } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { useWebSocket } from '../contexts/WebSocketContext'
import api from '../services/api'
import Sidebar from '../components/Sidebar'
import ChatRoom from '../components/ChatRoom'
import CreateRoomModal from '../components/CreateRoomModal'
import './Chat.css'

export default function Chat() {
  const { user, logout } = useAuth()
  const { isConnected, onlineUsers } = useWebSocket()
  const [rooms, setRooms] = useState([])
  const [selectedRoom, setSelectedRoom] = useState(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [sidebarOpen, setSidebarOpen] = useState(true)

  useEffect(() => {
    loadRooms()
  }, [])

  const loadRooms = async () => {
    try {
      const data = await api.getRooms()
      setRooms(data)
    } catch (error) {
      console.error('Failed to load rooms:', error)
    }
  }

  const handleRoomCreated = (room) => {
    setRooms(prev => [room, ...prev])
    setSelectedRoom(room)
    setShowCreateModal(false)
  }

  const handleSelectRoom = (room) => {
    setSelectedRoom(room)
    if (window.innerWidth < 768) {
      setSidebarOpen(false)
    }
  }

  return (
    <div className="chat-layout">
      {/* Mobile Header */}
      <header className="mobile-header">
        <button 
          className="menu-btn"
          onClick={() => setSidebarOpen(!sidebarOpen)}
        >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M3 12h18M3 6h18M3 18h18" />
          </svg>
        </button>
        <div className="header-title">
          {selectedRoom ? selectedRoom.name : 'Select a room'}
        </div>
        <div className={`connection-status ${isConnected ? 'connected' : ''}`}>
          <span className="status-dot"></span>
        </div>
      </header>

      {/* Sidebar */}
      <Sidebar
        rooms={rooms}
        selectedRoom={selectedRoom}
        onSelectRoom={handleSelectRoom}
        onCreateRoom={() => setShowCreateModal(true)}
        onlineUsers={onlineUsers}
        user={user}
        onLogout={logout}
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
      />

      {/* Main Chat Area */}
      <main className="chat-main">
        {selectedRoom ? (
          <ChatRoom room={selectedRoom} />
        ) : (
          <div className="no-room-selected">
            <div className="no-room-content">
              <div className="no-room-icon">
                <svg width="80" height="80" viewBox="0 0 80 80" fill="none">
                  <rect width="80" height="80" rx="20" fill="url(#noRoomGradient)" fillOpacity="0.1" />
                  <path d="M25 30C25 27.2386 27.2386 25 30 25H50C52.7614 25 55 27.2386 55 30V43.3333C55 46.0948 52.7614 48.3333 50 48.3333H36.6667L28.3333 55V48.3333H30C27.2386 48.3333 25 46.0948 25 43.3333V30Z" stroke="url(#noRoomGradient)" strokeWidth="2.5" />
                  <circle cx="35" cy="36.6667" r="2.5" fill="url(#noRoomGradient)" />
                  <circle cx="45" cy="36.6667" r="2.5" fill="url(#noRoomGradient)" />
                  <defs>
                    <linearGradient id="noRoomGradient" x1="0" y1="0" x2="80" y2="80">
                      <stop stopColor="#6366f1" />
                      <stop offset="1" stopColor="#8b5cf6" />
                    </linearGradient>
                  </defs>
                </svg>
              </div>
              <h2>Welcome to Chat</h2>
              <p>Select a room from the sidebar or create a new one to start chatting</p>
              <button 
                className="btn btn-primary"
                onClick={() => setShowCreateModal(true)}
              >
                <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
                </svg>
                Create Room
              </button>
            </div>
          </div>
        )}
      </main>

      {/* Create Room Modal */}
      {showCreateModal && (
        <CreateRoomModal
          onClose={() => setShowCreateModal(false)}
          onCreated={handleRoomCreated}
        />
      )}
    </div>
  )
}
