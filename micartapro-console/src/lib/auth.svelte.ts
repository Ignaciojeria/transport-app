import { supabase } from './supabase'
import type { User, Session } from '@supabase/supabase-js'

// Estado reactivo usando runes de Svelte 5
class AuthState {
  user = $state<User | null>(null)
  session = $state<Session | null>(null)
  loading = $state<boolean>(true)
}

// Instancia única del estado
export const authState = new AuthState()

interface AuthData {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: {
    id: string
    email: string
    verified_email: boolean
    name: string
    given_name: string
    family_name: string
    picture: string
    locale: string
  }
  timestamp: number
  provider: string
}

// Procesar fragment de autenticación desde micartapro-auth-ui
export function processAuthFragment(): void {
  const fragment = window.location.hash
  
  if (fragment.startsWith('#auth=')) {
    try {
      console.log('🔍 Fragment de auth detectado:', fragment.substring(0, 50) + '...')
      
      // Decodificar el payload
      const encodedData = fragment.substring(6) // Remove '#auth='
      const authData: AuthData = JSON.parse(atob(encodedData))
      
      // Mostrar el contenido completo del fragment decodificado
      console.log('✅ AUTH PAYLOAD DECODIFICADO (completo):', JSON.stringify(authData, null, 2))
      console.log('📋 Estructura del fragment:')
      console.log('  - provider:', authData.provider)
      console.log('  - access_token:', authData.access_token ? `${authData.access_token.substring(0, 30)}...` : 'N/A')
      console.log('  - refresh_token:', authData.refresh_token ? `${authData.refresh_token.substring(0, 30)}...` : 'N/A')
      console.log('  - token_type:', authData.token_type)
      console.log('  - expires_in:', authData.expires_in)
      console.log('  - user:', {
        id: authData.user?.id,
        email: authData.user?.email,
        name: authData.user?.name,
        verified_email: authData.user?.verified_email,
        picture: authData.user?.picture,
      })
      console.log('  - timestamp:', new Date(authData.timestamp).toISOString())
      
      // Si viene de Supabase, establecer la sesión
      if (authData.provider === 'supabase' && authData.access_token) {
        supabase.auth.setSession({
          access_token: authData.access_token,
          refresh_token: authData.refresh_token,
        }).then(({ data: sessionData, error }) => {
          if (error) {
            console.error('❌ Error estableciendo sesión:', error)
          } else {
            console.log('✅ Sesión establecida:', sessionData.session?.user?.email)
            authState.session = sessionData.session
            authState.user = sessionData.session?.user ?? null
          }
        })
      }
      
      // Limpiar fragment de la URL
      window.history.replaceState({}, '', window.location.pathname)
      
    } catch (err) {
      console.error('❌ Error procesando fragment:', err)
    }
  }
}

// Inicializar autenticación
export function initAuth(): void {
  // Procesar fragment primero
  processAuthFragment()
  
  // Obtener sesión actual
  supabase.auth.getSession().then(({ data: { session: currentSession } }) => {
    authState.session = currentSession
    authState.user = currentSession?.user ?? null
    authState.loading = false
  })
  
  // Escuchar cambios en la autenticación
  supabase.auth.onAuthStateChange((_event, newSession) => {
    authState.session = newSession
    authState.user = newSession?.user ?? null
    authState.loading = false
  })
}

// Cerrar sesión
export async function signOut(): Promise<void> {
  try {
    // Cerrar sesión en Supabase
    const { error } = await supabase.auth.signOut()
    if (error) {
      console.error('Error al cerrar sesión:', error)
      // Continuar con la redirección incluso si hay error
    }
    
    // Limpiar estado local
    authState.user = null
    authState.session = null
    
    // Limpiar localStorage de Supabase
    const supabaseKeys: string[] = []
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && (
        key.startsWith('sb-') || 
        key.startsWith('supabase.') ||
        key.includes('supabase') ||
        key.startsWith('sb-auth-')
      )) {
        supabaseKeys.push(key)
      }
    }
    supabaseKeys.forEach(key => {
      try {
        localStorage.removeItem(key)
      } catch (e) {
        console.warn('No se pudo eliminar la clave:', key, e)
      }
    })
    
    // Limpiar caché de clientes autenticados
    if (typeof window !== 'undefined' && (window as any).clearAuthenticatedClientsCache) {
      (window as any).clearAuthenticatedClientsCache()
    }
    
    // Redirigir a auth-ui con parámetro de logout para que cierre la sesión allí también
    const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    const authUiUrl = isLocalDev ? 'http://localhost:3003' : 'https://auth.micartapro.com'
    window.location.replace(`${authUiUrl}?logout=true`)
  } catch (error) {
    console.error('Error al cerrar sesión:', error)
    // Asegurarse de limpiar todo incluso si hay un error
    authState.user = null
    authState.session = null
    throw error
  }
}

