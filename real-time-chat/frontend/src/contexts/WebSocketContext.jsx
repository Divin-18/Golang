import { createContext, useContext, useState, useEffect, useRef, useCallback } from 'react'

const WebSocketContext = createContext(null)

const WS_URL = 'ws://localhost:8080/ws'

export function WebSocketProvider({ children }) {
  const [isConnected, setIsConnected] = useState(false)
  const [messages, setMessages] = useState({})
  const [onlineUsers, setOnlineUsers] = useState([])
  const [typingUsers, setTypingUsers] = useState({})
  const wsRef = useRef(null)
  const reconnectTimeoutRef = useRef(null)
  const messageHandlersRef = useRef([])

  const connect = useCallback(() => {
    const token = localStorage.getItem('token')
    if (!token) return

    try {
      wsRef.current = new WebSocket(`${WS_URL}?token=${token}`)

      wsRef.current.onopen = () => {
        console.log('WebSocket connected')
        setIsConnected(true)
      }

      wsRef.current.onclose = () => {
        console.log('WebSocket disconnected')
        setIsConnected(false)
        
        // Reconnect after 3 seconds
        reconnectTimeoutRef.current = setTimeout(() => {
          connect()
        }, 3000)
      }

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      wsRef.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          handleMessage(data)
        } catch (error) {
          console.error('Failed to parse message:', error)
        }
      }
    } catch (error) {
      console.error('Failed to connect:', error)
    }
  }, [])

  const handleMessage = (data) => {
    switch (data.type) {
      case 'new_message':
        const message = data.payload
        setMessages(prev => ({
          ...prev,
          [message.room_id]: [...(prev[message.room_id] || []), message]
        }))
        // Call registered handlers
        messageHandlersRef.current.forEach(handler => handler(message))
        break

      case 'online_users':
        setOnlineUsers(data.payload || [])
        break

      case 'user_joined':
        console.log(`${data.payload.username} joined room ${data.payload.room_id}`)
        break

      case 'user_left':
        console.log(`${data.payload.username} left room ${data.payload.room_id}`)
        break

      case 'typing':
        const typing = data.payload
        setTypingUsers(prev => ({
          ...prev,
          [typing.room_id]: typing.is_typing 
            ? [...(prev[typing.room_id] || []).filter(u => u !== typing.username), typing.username]
            : (prev[typing.room_id] || []).filter(u => u !== typing.username)
        }))
        break

      default:
        console.log('Unknown message type:', data.type)
    }
  }

  useEffect(() => {
    connect()

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [connect])

  const sendMessage = useCallback((type, payload) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type, payload }))
    }
  }, [])

  const joinRoom = useCallback((roomId) => {
    sendMessage('join_room', { room_id: roomId })
  }, [sendMessage])

  const leaveRoom = useCallback((roomId) => {
    sendMessage('leave_room', { room_id: roomId })
  }, [sendMessage])

  const sendChatMessage = useCallback((roomId, content) => {
    sendMessage('send_message', { room_id: roomId, content })
  }, [sendMessage])

  const sendTyping = useCallback((roomId, isTyping) => {
    sendMessage('typing', { room_id: roomId, is_typing: isTyping })
  }, [sendMessage])

  const addMessageToRoom = useCallback((roomId, newMessages) => {
    setMessages(prev => ({
      ...prev,
      [roomId]: newMessages
    }))
  }, [])

  const onMessage = useCallback((handler) => {
    messageHandlersRef.current.push(handler)
    return () => {
      messageHandlersRef.current = messageHandlersRef.current.filter(h => h !== handler)
    }
  }, [])

  return (
    <WebSocketContext.Provider value={{
      isConnected,
      messages,
      onlineUsers,
      typingUsers,
      joinRoom,
      leaveRoom,
      sendChatMessage,
      sendTyping,
      addMessageToRoom,
      onMessage
    }}>
      {children}
    </WebSocketContext.Provider>
  )
}

export function useWebSocket() {
  const context = useContext(WebSocketContext)
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider')
  }
  return context
}
