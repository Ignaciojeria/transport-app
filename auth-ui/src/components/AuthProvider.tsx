'use client'

import { createContext, useContext, useEffect, useState } from 'react'
import { AuthState, AuthUser, AuthTokens } from '@/types/auth'
import { AuthStorage } from '@/lib/auth-storage'

interface AuthContextType extends AuthState {
  login: (accessToken: string, refreshToken: string, expiresIn: number, user: AuthUser) => void
  logout: () => void
  refreshTokens: () => Promise<boolean>
}

const AuthContext = createContext<AuthContextType | null>(null)

interface AuthProviderProps {
  children: React.ReactNode
}

export default function AuthProvider({ children }: AuthProviderProps) {
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    user: null,
    tokens: null
  })

  // Cargar estado inicial desde localStorage
  useEffect(() => {
    const initialState = AuthStorage.getAuthState()
    setAuthState(initialState)
  }, [])

  const login = (accessToken: string, refreshToken: string, expiresIn: number, user: AuthUser) => {
    AuthStorage.saveTokens(accessToken, refreshToken, expiresIn)
    AuthStorage.saveUser(user)
    
    const tokens: AuthTokens = {
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_at: Date.now() + (expiresIn * 1000),
      token_type: 'Bearer'
    }

    setAuthState({
      isAuthenticated: true,
      user,
      tokens
    })
  }

  const logout = () => {
    AuthStorage.clearAuth()
    setAuthState({
      isAuthenticated: false,
      user: null,
      tokens: null
    })
  }

  const refreshTokens = async (): Promise<boolean> => {
    const refreshToken = AuthStorage.getRefreshToken()
    if (!refreshToken) {
      logout()
      return false
    }

    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `grant_type=refresh_token&refresh_token=${refreshToken}`
      })

      if (!response.ok) {
        logout()
        return false
      }

      const result = await response.json()
      
      if (result.access_token && authState.user) {
        AuthStorage.saveTokens(result.access_token, result.refresh_token, result.expires_in)
        
        const newTokens: AuthTokens = {
          access_token: result.access_token,
          refresh_token: result.refresh_token,
          expires_at: Date.now() + (result.expires_in * 1000),
          token_type: 'Bearer'
        }

        setAuthState(prev => ({
          ...prev,
          tokens: newTokens
        }))

        return true
      }

      logout()
      return false
    } catch (error) {
      console.error('Error refreshing tokens:', error)
      logout()
      return false
    }
  }

  // Auto-refresh tokens si están cerca de expirar
  useEffect(() => {
    if (!authState.isAuthenticated) return

    const checkTokens = () => {
      if (AuthStorage.needsRefresh()) {
        refreshTokens()
      }
    }

    // Verificar cada 5 minutos
    const interval = setInterval(checkTokens, 5 * 60 * 1000)
    
    // Verificar inmediatamente si es necesario
    checkTokens()

    return () => clearInterval(interval)
  }, [authState.isAuthenticated])

  const contextValue: AuthContextType = {
    ...authState,
    login,
    logout,
    refreshTokens
  }

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
