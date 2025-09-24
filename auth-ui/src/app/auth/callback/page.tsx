'use client'

import { useSearchParams, useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { GoogleExchangeResponse } from '@/types/auth'
import { AuthStorage } from '@/lib/auth-storage'

export default function AuthCallback() {
  const [status, setStatus] = useState('Procesando autenticación...')
  const [isLoading, setIsLoading] = useState(true)
  const [debugInfo, setDebugInfo] = useState<any>(null)
  const [authSuccess, setAuthSuccess] = useState(false)
  const searchParams = useSearchParams()
  const router = useRouter()

  useEffect(() => {
    const processAuth = async () => {
      try {
        const code = searchParams.get('code')
        const state = searchParams.get('state')
        const error = searchParams.get('error')

        console.log('🔍 Parámetros de callback recibidos:', {
          code: code ? `${code.substring(0, 20)}...` : null,
          state,
          error,
          full_url: window.location.href
        })

        if (error) {
          console.error('❌ Error en callback de Google:', error)
          setStatus(`Error: ${error}`)
          setIsLoading(false)
          return
        }

        if (!code) {
          console.error('❌ Código de autorización faltante')
          setStatus('Error: Código de autorización faltante')
          setIsLoading(false)
          return
        }

        // Validar state
        const savedState = AuthStorage.getAndClearOAuthState()
        console.log('🔍 Validando state CSRF:', { received: state, saved: savedState })
        
        if (state !== savedState) {
          console.error('❌ Estado de seguridad inválido')
          setStatus('Error: Estado de seguridad inválido')
          setIsLoading(false)
          return
        }

        setStatus('Enviando código al backend...')

        const requestBody = {
          code: code,
          state: state,
          redirect_uri: window.location.origin + '/auth/callback',
        }

        console.log('🚀 Enviando request al backend:', {
          url: 'https://einar-main-f0820bc.d2.zuplo.dev/auth/google/exchange',
          method: 'POST',
          body: requestBody
        })

        // Enviar código al backend para que haga el intercambio seguro
        const backendResponse = await fetch('https://einar-main-f0820bc.d2.zuplo.dev/auth/google/exchange', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(requestBody),
        })

        console.log('📡 Respuesta del backend:', {
          status: backendResponse.status,
          statusText: backendResponse.statusText,
          headers: Object.fromEntries(backendResponse.headers.entries())
        })

        if (!backendResponse.ok) {
          const errorText = await backendResponse.text()
          console.error('❌ Error HTTP del backend:', {
            status: backendResponse.status,
            statusText: backendResponse.statusText,
            body: errorText
          })
          throw new Error(`Error validando con backend: ${backendResponse.status} ${backendResponse.statusText}`)
        }

        const result: GoogleExchangeResponse = await backendResponse.json()
        console.log('📦 Respuesta JSON del backend:', result)

        if (result.error) {
          console.error('❌ Error en respuesta del backend:', result.error)
          throw new Error(result.error)
        }

        if (result.access_token && result.user) {
          // Guardar tokens usando el nuevo sistema
          AuthStorage.saveTokens(
            result.access_token,
            result.refresh_token,
            result.expires_in
          )

          // Guardar información del usuario
          AuthStorage.saveUser({
            id: result.user.id,
            email: result.user.email,
            name: result.user.name,
            picture: result.user.picture,
            verified_email: result.user.verified_email
          })
          
          setStatus(`¡Bienvenido ${result.user.name}! Autenticación completada.`)
          setIsLoading(false)
          setAuthSuccess(true)
          
          // Información de debug para la interfaz
          const debugData = {
            backend_response: result,
            tokens_stored: AuthStorage.getTokens(),
            user_stored: AuthStorage.getUser(),
            is_authenticated: AuthStorage.isAuthenticated(),
            token_expires_at: new Date(Date.now() + (result.expires_in * 1000)).toISOString(),
            needs_refresh: AuthStorage.needsRefresh()
          }
          
          setDebugInfo(debugData)
          
          // Log para inspección en consola
          console.log('✅ Autenticación exitosa:', debugData)
          
          // NO redirigir automáticamente - mantener en esta pantalla para inspección
        } else {
          throw new Error('Respuesta incompleta del servidor')
        }

      } catch (err: any) {
        console.error('❌ Error en autenticación:', {
          error: err,
          message: err.message,
          stack: err.stack,
          response_received: err.response,
          network_error: err.cause
        })
        setStatus(`Error: ${err.message || 'Error desconocido durante la autenticación.'}`)
        setIsLoading(false)
        // NO redirigir automáticamente para inspeccionar errores
      }
    }

    processAuth()
  }, [searchParams, router])

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 particles-container flex flex-col">
      {/* Partículas de fondo (mismo estilo que la página principal) */}
      <div className="auth-particles">
        {[...Array(20)].map((_, i) => (
          <div
            key={i}
            className="particle"
            style={{
              left: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 15}s`,
              animationDuration: `${15 + Math.random() * 10}s`
            }}
          />
        ))}
      </div>

      {/* Header con logo (igual que página principal) */}
      <header className="relative z-10 p-6">
        <div className="flex items-center space-x-3">
          <div className="bg-blue-600 p-2 rounded-xl">
            <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-800">TransportApp</h1>
            <p className="text-sm text-gray-600">Sistema de gestión logística</p>
          </div>
        </div>
      </header>

      {/* Contenido principal */}
      <div className="flex-1 flex items-center justify-center p-6">
        <div className="w-full max-w-md">
          <div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-8">
            <div className="text-center">
              {/* Icono de loading o éxito */}
              <div className="mb-6">
                {isLoading ? (
                  <div className="relative">
                    <div className="animate-spin rounded-full h-16 w-16 border-4 border-blue-600/20 border-t-blue-600 mx-auto"></div>
                    <div className="absolute inset-0 flex items-center justify-center">
                      <svg className="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                      </svg>
                    </div>
                  </div>
                ) : (
                  <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto">
                    <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                )}
              </div>

              {/* Título */}
              <h2 className="text-2xl font-bold text-gray-800 mb-2">
                {isLoading ? 'Iniciando Sesión' : '¡Bienvenido!'}
              </h2>

              {/* Estado */}
              <p className="text-gray-600 mb-6">
                {status}
              </p>

              {/* Indicador de progreso */}
              {isLoading && (
                <div className="w-full bg-gray-200 rounded-full h-2 mb-4">
                  <div className="bg-blue-600 h-2 rounded-full animate-pulse" style={{width: '60%'}}></div>
                </div>
              )}

              {/* Mensaje adicional */}
              {isLoading && (
                <div className="text-sm text-gray-500">
                  Validando credenciales con Google...
                </div>
              )}

              {/* Información de debug cuando la autenticación es exitosa */}
              {authSuccess && debugInfo && (
                <div className="mt-6 space-y-4">
                  <div className="bg-green-50 border border-green-200 rounded-xl p-4">
                    <h3 className="font-medium text-green-900 mb-2">✅ Debug Info - Revisa la consola para más detalles</h3>
                    <div className="text-xs text-green-800 space-y-1">
                      <div>🎯 Usuario: {debugInfo.user_stored?.name} ({debugInfo.user_stored?.email})</div>
                      <div>🔑 Access Token: {debugInfo.tokens_stored?.access_token ? '✅ Guardado' : '❌ No guardado'}</div>
                      <div>🔄 Refresh Token: {debugInfo.backend_response?.refresh_token ? '✅ Recibido' : '❌ No recibido'}</div>
                      <div>⏰ Expira: {debugInfo.token_expires_at}</div>
                      <div>🔐 Autenticado: {debugInfo.is_authenticated ? '✅ Sí' : '❌ No'}</div>
                    </div>
                  </div>
                  
                  <div className="space-y-3">
                    <button
                      onClick={() => router.push('/dashboard')}
                      className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 px-4 rounded-xl font-medium transition-colors"
                    >
                      Continuar al Dashboard
                    </button>
                    
                    <button
                      onClick={() => router.push('/')}
                      className="w-full bg-gray-100 hover:bg-gray-200 text-gray-700 py-3 px-4 rounded-xl font-medium transition-colors"
                    >
                      Volver al Inicio
                    </button>
                    
                    <button
                      onClick={() => {
                        console.log('🔍 Estado actual completo:', {
                          searchParams: Object.fromEntries(searchParams.entries()),
                          localStorage_keys: Object.keys(localStorage),
                          auth_tokens: AuthStorage.getTokens(),
                          auth_user: AuthStorage.getUser(),
                          auth_state: AuthStorage.getAuthState()
                        })
                      }}
                      className="w-full bg-purple-100 hover:bg-purple-200 text-purple-700 py-2 px-4 rounded-xl text-sm font-medium transition-colors"
                    >
                      🔍 Log Estado Completo en Consola
                    </button>
                  </div>
                </div>
              )}

              {/* Información de debug cuando hay error */}
              {!isLoading && !authSuccess && (
                <div className="mt-6 space-y-4">
                  <div className="bg-red-50 border border-red-200 rounded-xl p-4">
                    <h3 className="font-medium text-red-900 mb-2">❌ Error - Revisa la consola para más detalles</h3>
                    <div className="text-xs text-red-800">
                      Revisa la pestaña Network y Console en DevTools
                    </div>
                  </div>
                  
                  <button
                    onClick={() => router.push('/')}
                    className="w-full bg-gray-100 hover:bg-gray-200 text-gray-700 py-3 px-4 rounded-xl font-medium transition-colors"
                  >
                    Volver a Intentar
                  </button>
                </div>
              )}
            </div>
          </div>

          {/* Footer info */}
          <div className="mt-6 text-center">
            <p className="text-xs text-gray-500">
              Conectando de forma segura • SSL/TLS Encriptado
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
