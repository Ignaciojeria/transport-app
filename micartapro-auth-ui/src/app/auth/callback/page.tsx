'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { supabase } from '@/lib/supabase'

export default function AuthCallback() {
  const [status, setStatus] = useState('Procesando autenticación...')
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  useEffect(() => {
    const handleAuthCallback = async () => {
      try {
        // Supabase maneja automáticamente el callback de OAuth
        // Solo necesitamos verificar si hay una sesión
        const { data: { session }, error } = await supabase.auth.getSession()

        if (error) {
          console.error('❌ Error obteniendo sesión:', error)
          setStatus(`Error: ${error.message}`)
          setIsLoading(false)
          setTimeout(() => {
            router.push('/')
          }, 3000)
          return
        }

        if (session && session.user) {
          setStatus(`¡Bienvenido ${session.user.email}! Redirigiendo...`)
          
          // Esperar un momento para que el AuthProvider actualice el estado
          setTimeout(() => {
            // Redirigir a la página principal, que manejará la redirección final
            router.push('/')
          }, 1500)
        } else {
          setStatus('No se encontró sesión. Redirigiendo...')
          setTimeout(() => {
            router.push('/')
          }, 2000)
        }
      } catch (err: any) {
        console.error('❌ Error en callback:', err)
        setStatus(`Error: ${err.message || 'Error desconocido'}`)
        setIsLoading(false)
        setTimeout(() => {
          router.push('/')
        }, 3000)
      }
    }

    handleAuthCallback()
  }, [router])

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex flex-col items-center justify-center p-6">
      <div className="w-full max-w-md">
        <div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-8">
          <div className="text-center">
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

            <h2 className="text-2xl font-bold text-gray-800 mb-2">
              {isLoading ? 'Iniciando Sesión' : '¡Autenticación Exitosa!'}
            </h2>

            <p className="text-gray-600 mb-6">
              {status}
            </p>

            {isLoading && (
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-blue-600 h-2 rounded-full animate-pulse" style={{width: '75%'}}></div>
              </div>
            )}

            {!isLoading && (
              <div className="text-sm text-gray-500">
                Serás redirigido automáticamente...
              </div>
            )}
          </div>
        </div>

        <div className="mt-6 text-center">
          <p className="text-xs text-gray-500">
            Conectando de forma segura • SSL/TLS Encriptado
          </p>
        </div>
      </div>
    </div>
  )
}
