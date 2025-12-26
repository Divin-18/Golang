import { useState, useEffect, useRef, useCallback } from 'react'
import { useWebSocket } from '../contexts/WebSocketContext'
import { useAuth } from '../contexts/AuthContext'
import api from '../services/api'
import './ChatRoom.css'

export default function ChatRoom({ room }) {
  const { user } = useAuth()
  const { messages, joinRoom, sendChatMessage, sendTyping, typingUsers, addMessageToRoom } = useWebSocket()
  const [newMessage, setNewMessage] = useState('')
  const [loading, setLoading] = useState(true)
  const messagesEndRef = useRef(null)
  const typingTimeoutRef = useRef(null)
  const roomMessages = messages[room.id] || []
  const roomTyping = typingUsers[room.id] || []

  useEffect(() => {
    loadMessages()
    joinRoom(room.id)
  }, [room.id])

  useEffect(() => {
    scrollToBottom()
  }, [roomMessages])

  const loadMessages = async () => {
    setLoading(true)
    try {
      const data = await api.getRoomMessages(room.id)
      addMessageToRoom(room.id, data)
    } catch (error) {
      console.error('Failed to load messages:', error)
    } finally {
      setLoading(false)
    }
  }

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const handleTyping = useCallback((e) => {
    setNewMessage(e.target.value)
    sendTyping(room.id, true)
    if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current)
    typingTimeoutRef.current = setTimeout(() => sendTyping(room.id, false), 2000)
  }, [room.id, sendTyping])

  const handleSend = (e) => {
    e.preventDefault()
    if (!newMessage.trim()) return
    sendChatMessage(room.id, newMessage.trim())
    setNewMessage('')
    sendTyping(room.id, false)
  }

  const formatTime = (dateString) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  const formatDate = (dateString) => {
    const date = new Date(dateString)
    const today = new Date()
    if (date.toDateString() === today.toDateString()) return 'Today'
    const yesterday = new Date(today)
    yesterday.setDate(yesterday.getDate() - 1)
    if (date.toDateString() === yesterday.toDateString()) return 'Yesterday'
    return date.toLocaleDateString()
  }

  let lastDate = ''

  return (
    <div className="chat-room">
      <div className="chat-header">
        <div className="room-title">
          <span className="room-hash">#</span>
          <h2>{room.name}</h2>
        </div>
        {room.description && <p className="room-desc">{room.description}</p>}
      </div>

      <div className="messages-container">
        {loading ? (
          <div className="messages-loading"><div className="spinner"></div><p>Loading messages...</p></div>
        ) : roomMessages.length === 0 ? (
          <div className="no-messages"><p>No messages yet. Start the conversation!</p></div>
        ) : (
          roomMessages.map((msg, index) => {
            const msgDate = formatDate(msg.created_at)
            const showDate = msgDate !== lastDate
            lastDate = msgDate
            const isOwn = msg.user_id === user?.id

            return (
              <div key={msg.id || index}>
                {showDate && <div className="date-divider"><span>{msgDate}</span></div>}
                <div className={`message ${isOwn ? 'own' : ''}`}>
                  {!isOwn && <div className="msg-avatar">{msg.username?.charAt(0).toUpperCase()}</div>}
                  <div className="msg-content">
                    {!isOwn && <span className="msg-username">{msg.username}</span>}
                    <div className="msg-bubble">
                      <p>{msg.content}</p>
                      <span className="msg-time">{formatTime(msg.created_at)}</span>
                    </div>
                  </div>
                </div>
              </div>
            )
          })
        )}
        {roomTyping.length > 0 && (
          <div className="typing-indicator">
            <span>{roomTyping.join(', ')} {roomTyping.length === 1 ? 'is' : 'are'} typing...</span>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      <form className="message-form" onSubmit={handleSend}>
        <input
          type="text"
          className="input message-input"
          placeholder={`Message #${room.name}`}
          value={newMessage}
          onChange={handleTyping}
        />
        <button type="submit" className="btn btn-primary send-btn" disabled={!newMessage.trim()}>
          Send
        </button>
      </form>
    </div>
  )
}
