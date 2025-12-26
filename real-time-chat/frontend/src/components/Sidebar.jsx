import { useWebSocket } from '../contexts/WebSocketContext'
import './Sidebar.css'

export default function Sidebar({ rooms, selectedRoom, onSelectRoom, onCreateRoom, onlineUsers, user, onLogout, isOpen, onClose }) {
  const { isConnected } = useWebSocket()

  return (
    <>
      <div className={`sidebar-overlay ${isOpen ? 'active' : ''}`} onClick={onClose} />
      <aside className={`sidebar ${isOpen ? 'open' : ''}`}>
        <div className="sidebar-header">
          <div className="brand">
            <div className="brand-logo">💬</div>
            <div className="brand-info">
              <h1>ChatApp</h1>
              <div className={`connection-badge ${isConnected ? 'online' : ''}`}>
                <span className="badge-dot"></span>
                {isConnected ? 'Connected' : 'Connecting...'}
              </div>
            </div>
          </div>
        </div>

        <div className="sidebar-section">
          <div className="section-header">
            <h3>Rooms</h3>
            <button className="icon-btn" onClick={onCreateRoom}>+</button>
          </div>
          <div className="rooms-list">
            {rooms.length === 0 ? (
              <div className="empty-rooms"><p>No rooms yet</p></div>
            ) : (
              rooms.map(room => (
                <button key={room.id} className={`room-item ${selectedRoom?.id === room.id ? 'active' : ''}`} onClick={() => onSelectRoom(room)}>
                  <div className="room-icon">#</div>
                  <div className="room-info">
                    <span className="room-name">{room.name}</span>
                  </div>
                </button>
              ))
            )}
          </div>
        </div>

        <div className="sidebar-section">
          <div className="section-header">
            <h3>Online</h3>
            <span className="online-count">{onlineUsers?.length || 0}</span>
          </div>
          <div className="users-list">
            {onlineUsers?.map(u => (
              <div key={u.id} className="user-item">
                <div className="user-avatar">{u.username.charAt(0).toUpperCase()}<span className="status-indicator"></span></div>
                <span className="user-name">{u.username}{u.id === user?.id && ' (You)'}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="sidebar-footer">
          <div className="user-profile">
            <div className="profile-avatar">{user?.username?.charAt(0).toUpperCase()}</div>
            <div className="profile-info">
              <span className="profile-name">{user?.username}</span>
              <span className="profile-email">{user?.email}</span>
            </div>
            <button className="icon-btn logout-btn" onClick={onLogout}>⏻</button>
          </div>
        </div>
      </aside>
    </>
  )
}
