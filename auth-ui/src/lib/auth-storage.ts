import { AuthTokens, AuthUser, AuthState } from '@/types/auth'

const AUTH_TOKENS_KEY = 'auth_tokens'
const AUTH_USER_KEY = 'auth_user'
const OAUTH_STATE_KEY = 'oauth_state'

export class AuthStorage {
  /**
   * Guarda los tokens de autenticación
   */
  static saveTokens(accessToken: string, refreshToken: string, expiresIn: number): void {
    const tokens: AuthTokens = {
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_at: Date.now() + (expiresIn * 1000), // convertir segundos a timestamp
      token_type: 'Bearer'
    }
    
    if (typeof window !== 'undefined') {
      localStorage.setItem(AUTH_TOKENS_KEY, JSON.stringify(tokens))
    }
  }

  /**
   * Guarda la información del usuario
   */
  static saveUser(user: AuthUser): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user))
    }
  }

  /**
   * Obtiene los tokens almacenados
   */
  static getTokens(): AuthTokens | null {
    if (typeof window === 'undefined') return null
    
    const stored = localStorage.getItem(AUTH_TOKENS_KEY)
    if (!stored) return null
    
    try {
      const tokens: AuthTokens = JSON.parse(stored)
      
      // Verificar si el token ha expirado
      if (Date.now() >= tokens.expires_at) {
        this.clearTokens()
        return null
      }
      
      return tokens
    } catch {
      return null
    }
  }

  /**
   * Obtiene la información del usuario almacenada
   */
  static getUser(): AuthUser | null {
    if (typeof window === 'undefined') return null
    
    const stored = localStorage.getItem(AUTH_USER_KEY)
    if (!stored) return null
    
    try {
      return JSON.parse(stored)
    } catch {
      return null
    }
  }

  /**
   * Obtiene el estado completo de autenticación
   */
  static getAuthState(): AuthState {
    const tokens = this.getTokens()
    const user = this.getUser()
    
    return {
      isAuthenticated: !!(tokens && user),
      tokens,
      user
    }
  }

  /**
   * Verifica si el usuario está autenticado y el token es válido
   */
  static isAuthenticated(): boolean {
    const tokens = this.getTokens()
    const user = this.getUser()
    return !!(tokens && user)
  }

  /**
   * Obtiene el access token válido
   */
  static getAccessToken(): string | null {
    const tokens = this.getTokens()
    return tokens?.access_token || null
  }

  /**
   * Obtiene el refresh token
   */
  static getRefreshToken(): string | null {
    const tokens = this.getTokens()
    return tokens?.refresh_token || null
  }

  /**
   * Limpia todos los datos de autenticación
   */
  static clearAuth(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(AUTH_TOKENS_KEY)
      localStorage.removeItem(AUTH_USER_KEY)
    }
  }

  /**
   * Limpia solo los tokens
   */
  static clearTokens(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(AUTH_TOKENS_KEY)
    }
  }

  /**
   * Guarda el estado OAuth para validación CSRF
   */
  static saveOAuthState(state: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(OAUTH_STATE_KEY, state)
    }
  }

  /**
   * Obtiene y elimina el estado OAuth
   */
  static getAndClearOAuthState(): string | null {
    if (typeof window === 'undefined') return null
    
    const state = localStorage.getItem(OAUTH_STATE_KEY)
    if (state) {
      localStorage.removeItem(OAUTH_STATE_KEY)
    }
    return state
  }

  /**
   * Verifica si el token necesita renovación (expira en menos de 5 minutos)
   */
  static needsRefresh(): boolean {
    const tokens = this.getTokens()
    if (!tokens) return false
    
    const fiveMinutesFromNow = Date.now() + (5 * 60 * 1000)
    return tokens.expires_at <= fiveMinutesFromNow
  }
}
