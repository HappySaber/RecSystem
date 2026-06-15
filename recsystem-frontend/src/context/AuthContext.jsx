import { createContext, useContext, useMemo, useState } from 'react'
import { getUserFromToken } from '../utils/jwt'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('access_token'))
  const [profile, setProfile] = useState(() => {
    const stored = localStorage.getItem('access_token')
    return stored ? getUserFromToken(stored) : null
  })

  const login = (accessToken, extra = {}) => {
    localStorage.setItem('access_token', accessToken)
    setToken(accessToken)
    const fromToken = getUserFromToken(accessToken)
    setProfile({ ...fromToken, ...extra })
  }

  const logout = () => {
    localStorage.removeItem('access_token')
    setToken(null)
    setProfile(null)
  }

  const value = useMemo(
    () => ({
      token,
      user: profile,
      isAuthenticated: Boolean(token),
      login,
      logout,
    }),
    [token, profile]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
