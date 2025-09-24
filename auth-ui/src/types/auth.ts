export interface GoogleUserInfo {
  id: string
  email: string
  verified_email: boolean
  name: string
  given_name: string
  family_name: string
  picture: string
  locale: string
}

export interface GoogleExchangeResponse {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token: string
  user: GoogleUserInfo
  error?: string
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_at: number // timestamp
  token_type: string
}

export interface AuthUser {
  id: string
  email: string
  name: string
  picture?: string
  verified_email: boolean
}

export interface AuthState {
  isAuthenticated: boolean
  user: AuthUser | null
  tokens: AuthTokens | null
}
