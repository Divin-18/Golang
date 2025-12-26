const API_BASE_URL = 'http://localhost:8080/api'

class ApiService {
  constructor() {
    this.baseUrl = API_BASE_URL
  }

  getToken() {
    return localStorage.getItem('token')
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`
    const token = this.getToken()

    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'An error occurred')
    }

    return data
  }

  // Auth endpoints
  async register(username, email, password) {
    return this.request('/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    })
  }

  async login(email, password) {
    return this.request('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  }

  async getCurrentUser() {
    return this.request('/me')
  }

  // Room endpoints
  async getRooms() {
    return this.request('/rooms')
  }

  async createRoom(name, description = '', isPrivate = false) {
    return this.request('/rooms', {
      method: 'POST',
      body: JSON.stringify({ name, description, is_private: isPrivate }),
    })
  }

  async getRoom(roomId) {
    return this.request(`/rooms/${roomId}`)
  }

  async getUserRooms() {
    return this.request('/rooms/my')
  }

  async joinRoom(roomId) {
    return this.request(`/rooms/${roomId}/join`, {
      method: 'POST',
    })
  }

  async leaveRoom(roomId) {
    return this.request(`/rooms/${roomId}/leave`, {
      method: 'POST',
    })
  }

  async getRoomMembers(roomId) {
    return this.request(`/rooms/${roomId}/members`)
  }

  async getRoomMessages(roomId, limit = 50, offset = 0) {
    return this.request(`/rooms/${roomId}/messages?limit=${limit}&offset=${offset}`)
  }
}

export const api = new ApiService()
export default api
