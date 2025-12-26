import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import './Auth.css'

export default function Login({ onSwitch }) {
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await login(email, password)
    } catch (err) {
      setError(err.message || 'Failed to login')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-card animate-fadeIn">
      <div className="auth-header">
        <div className="auth-logo">
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <rect width="48" height="48" rx="12" fill="url(#gradient)" />
            <path d="M15 18C15 16.3431 16.3431 15 18 15H30C31.6569 15 33 16.3431 33 18V26C33 27.6569 31.6569 29 30 29H22L17 33V29H18C16.3431 29 15 27.6569 15 26V18Z" fill="white" fillOpacity="0.9" />
            <circle cx="21" cy="22" r="2" fill="url(#gradient)" />
            <circle cx="27" cy="22" r="2" fill="url(#gradient)" />
            <defs>
              <linearGradient id="gradient" x1="0" y1="0" x2="48" y2="48">
                <stop stopColor="#6366f1" />
                <stop offset="1" stopColor="#8b5cf6" />
              </linearGradient>
            </defs>
          </svg>
        </div>
        <h1>Welcome Back</h1>
        <p>Sign in to continue chatting</p>
      </div>

      <form onSubmit={handleSubmit} className="auth-form">
        {error && (
          <div className="error-message animate-fadeIn">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            {error}
          </div>
        )}

        <div className="input-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            className="input"
            placeholder="Enter your email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <div className="input-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            className="input"
            placeholder="Enter your password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <button 
          type="submit" 
          className="btn btn-primary auth-btn"
          disabled={loading}
        >
          {loading ? (
            <>
              <div className="spinner"></div>
              Signing in...
            </>
          ) : (
            'Sign In'
          )}
        </button>
      </form>

      <div className="auth-footer">
        <p>
          Don't have an account?{' '}
          <button onClick={onSwitch} className="link-btn">
            Create one
          </button>
        </p>
      </div>
    </div>
  )
}
