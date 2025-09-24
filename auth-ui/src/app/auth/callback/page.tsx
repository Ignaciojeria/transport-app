'use client'

import { useSearchParams, useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'

interface GoogleUserInfo {
  id: string
  email: string
  verified_email: boolean
  name: string
  given_name: string
  family_name: string
  picture: string
  locale: string
}

interface GoogleExchangeResponse {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token: string
  user: GoogleUserInfo
  error?: string
}

export default function AuthCallback() {
  const [status, setStatus] = useState('Procesando autenticación...')
  const [isLoading, setIsLoading] = useState(true)
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

        // Validar state CSRF
        const savedState = localStorage.getItem('oauth_state')
        localStorage.removeItem('oauth_state')
        console.log('🔍 Validando state CSRF:', { received: state, saved: savedState })
        
        if (state !== savedState) {
          console.error('❌ Estado de seguridad inválido')
          setStatus('Error: Estado de seguridad inválido')
          setIsLoading(false)
          return
        }

        setStatus('Intercambiando código con el backend...')

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

        const backendResponse = await fetch('https://einar-main-f0820bc.d2.zuplo.dev/auth/google/exchange', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(requestBody),
        })

        console.log('📡 Respuesta del backend:', {
          status: backendResponse.status,
          statusText: backendResponse.statusText
        })

        if (!backendResponse.ok) {
          const errorText = await backendResponse.text()
          console.error('❌ Error HTTP del backend:', {
            status: backendResponse.status,
            statusText: backendResponse.statusText,
            body: errorText
          })
          throw new Error(`Error del backend: ${backendResponse.status} ${backendResponse.statusText}`)
        }

        const result: GoogleExchangeResponse = await backendResponse.json()
        console.log('📦 Respuesta JSON del backend:', result)

        if (result.error) {
          console.error('❌ Error en respuesta del backend:', result.error)
          throw new Error(result.error)
        }

        if (result.access_token && result.user) {
          setStatus(`¡Bienvenido ${result.user.name}! Redirigiendo a la aplicación...`)
          
          console.log('✅ Autenticación exitosa, preparando redirect:', {
            access_token: `${result.access_token.substring(0, 20)}...`,
            refresh_token: `${result.refresh_token.substring(0, 20)}...`,
            user: result.user
          })
          
          // Preparar datos para el fragment
          const authData = {
            access_token: result.access_token,
            refresh_token: result.refresh_token,
            token_type: result.token_type,
            expires_in: result.expires_in,
            user: result.user,
            timestamp: Date.now()
          }
          
          // Encodear en base64 para el fragment
          const encodedAuth = btoa(JSON.stringify(authData))
          
          setIsLoading(false)
          
          // Redirect a la app con los datos en el fragment
          setTimeout(() => {
            const redirectUrl = `https://app.transport-app.com#auth=${encodedAuth}`
            console.log('🚀 Redirigiendo a:', redirectUrl)
            window.location.href = redirectUrl
          }, 2000) // 2 segundos para mostrar el mensaje
          
        } else {
          throw new Error('Respuesta incompleta del servidor')
        }

      } catch (err: any) {
        console.error('❌ Error en autenticación:', {
          error: err,
          message: err.message,
          stack: err.stack
        })
        setStatus(`Error: ${err.message || 'Error desconocido durante la autenticación.'}`)
        setIsLoading(false)
      }
    }

    processAuth()
  }, [searchParams, router])

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex flex-col items-center justify-center p-6">
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
              {isLoading ? 'Iniciando Sesión' : '¡Autenticación Exitosa!'}
            </h2>

            {/* Estado */}
            <p className="text-gray-600 mb-6">
              {status}
            </p>

            {/* Indicador de progreso */}
            {isLoading && (
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-blue-600 h-2 rounded-full animate-pulse" style={{width: '75%'}}></div>
              </div>
            )}

            {/* Mensaje adicional */}
            {!isLoading && (
              <div className="text-sm text-gray-500">
                Serás redirigido automáticamente...
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
  )
}