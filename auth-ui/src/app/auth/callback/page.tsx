'use client'

import { useSearchParams, useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'

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

        if (error) {
          setStatus(`Error: ${error}`)
          setIsLoading(false)
          setTimeout(() => router.push('/'), 3000)
          return
        }

        if (!code) {
          setStatus('Error: Código de autorización faltante')
          setIsLoading(false)
          setTimeout(() => router.push('/'), 3000)
          return
        }

        // Validar state
        const savedState = localStorage.getItem('oauth_state')
        if (state !== savedState) {
          setStatus('Error: Estado de seguridad inválido')
          setIsLoading(false)
          setTimeout(() => router.push('/'), 3000)
          return
        }

        setStatus('Enviando código al backend...')

        // Enviar código al backend para que haga el intercambio seguro
        const backendResponse = await fetch('https://einar-main-f0820bc.d2.zuplo.dev/auth/google/exchange', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            code: code,
            state: state,
            redirect_uri: window.location.origin + '/auth/callback',
          }),
        })

        if (!backendResponse.ok) {
          throw new Error('Error validando con backend')
        }

        const result = await backendResponse.json()

        if (result.token) {
          localStorage.setItem('auth_token', result.token)
          localStorage.removeItem('oauth_state')
          
          setStatus('¡Autenticación exitosa! Redirigiendo...')
          setIsLoading(false)
          
          setTimeout(() => {
            router.push('/?success=true')
          }, 1500)
        } else {
          throw new Error('No se recibió token del backend')
        }

      } catch (err) {
        console.error('Error en autenticación:', err)
        setStatus(`Error: ${err.message}`)
        setIsLoading(false)
        setTimeout(() => router.push('/?error=auth_failed'), 3000)
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
