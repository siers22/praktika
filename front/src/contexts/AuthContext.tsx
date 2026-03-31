'use client'

import { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { AuthUser, Role } from '@/types'

interface AuthContextType {
  user: AuthUser | null
  loading: boolean
  login: (accessToken: string, refreshToken: string) => void
  logout: () => void
  hasRole: (...roles: Role[]) => boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

function decodeToken(token: string): AuthUser | null {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return { id: payload.user_id, role: payload.role }
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('access_token')
    if (token) {
      setUser(decodeToken(token))
    }
    setLoading(false)
  }, [])

  const login = (accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)
    setUser(decodeToken(accessToken))
  }

  const logout = () => {
    const rt = localStorage.getItem('refresh_token')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    setUser(null)
    if (rt) {
      fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/auth/logout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: rt }),
      }).catch(() => {})
    }
  }

  const hasRole = (...roles: Role[]) => {
    if (!user) return false
    return roles.includes(user.role)
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, hasRole }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
