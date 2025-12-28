'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { supabase } from '@/lib/supabase'
import { useLanguage } from '@/lib/useLanguage'
import { getOrCreateMenuId } from '@/lib/menuId'

export default function AuthCallback() {
  const [status, setStatus] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()
  const { t, isLoading: langLoading, language } = useLanguage()
  
  // Initialize status with translated text
  useEffect(() => {
    if (!langLoading) {
      setStatus(t.callback.processing)
    }
  }, [langLoading, t])

  useEffect(() => {
    const handleAuthCallback = async () => {
      try {
        // Supabase maneja automáticamente el callback de OAuth
        // Solo necesitamos verificar si hay una sesión
        const { data: { session }, error } = await supabase.auth.getSession()

        if (error) {
          console.error('❌ Error obteniendo sesión:', error)
          setStatus(t.callback.error.replace('{message}', error.message))
          setIsLoading(false)
          setTimeout(() => {
            router.push('/')
          }, 3000)
          return
        }

        if (session && session.user) {
          setStatus(t.callback.welcome.replace('{email}', session.user.email || ''))
          
          // Obtener o crear menuID para el usuario
          try {
            const menuId = await getOrCreateMenuId(session.user.id)
            console.log('✅ MenuID procesado:', menuId)
          } catch (menuError) {
            console.error('⚠️ Error al obtener/crear menuID (continuando de todas formas):', menuError)
            // No bloqueamos el flujo si falla la creación del menuID
          }
          
          // Preparar datos para el fragment
          const userMetadata = session.user.user_metadata || {}
          const userInfo = {
            id: session.user.id,
            email: session.user.email || '',
            verified_email: session.user.email_confirmed_at ? true : false,
            name: userMetadata.full_name || userMetadata.name || session.user.email?.split('@')[0] || '',
            given_name: userMetadata.given_name || userMetadata.name?.split(' ')[0] || '',
            family_name: userMetadata.family_name || userMetadata.name?.split(' ').slice(1).join(' ') || '',
            picture: userMetadata.avatar_url || userMetadata.picture || '',
            locale: userMetadata.locale || 'es',
          }

          const authData = {
            access_token: session.access_token,
            token_type: 'Bearer',
            expires_in: session.expires_in || 3600,
            refresh_token: session.refresh_token || '',
            user: userInfo,
            timestamp: Date.now(),
            provider: 'supabase',
          }
          
          const encodedAuth = btoa(JSON.stringify(authData))
          
          // Detectar si estamos en desarrollo local
          const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
          const baseUrl = isLocalDev ? 'http://localhost:5174' : 'https://console.micartapro.com'
          // Preservar el parámetro de idioma en la redirección
          const redirectUrl = `${baseUrl}?lang=${language}#auth=${encodedAuth}`
          
          // Redirigir directamente a console después de un breve delay
          setTimeout(() => {
            window.location.href = redirectUrl
          }, 1000)
        } else {
          setStatus(t.callback.noSession)
          setTimeout(() => {
            router.push('/')
          }, 2000)
        }
      } catch (err: any) {
        console.error('❌ Error en callback:', err)
        setStatus(t.callback.error.replace('{message}', err.message || 'Error desconocido'))
        setIsLoading(false)
        setTimeout(() => {
          router.push('/')
        }, 3000)
      }
    }

    handleAuthCallback()
  }, [router, t, language])

  // Show loading while language is loading
  if (langLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }

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
              {isLoading ? t.callback.signingIn : t.callback.success}
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
                {t.callback.redirecting}
              </div>
            )}
          </div>
        </div>

        <div className="mt-6 text-center">
          <p className="text-xs text-gray-500">
            {t.callback.secureConnection}
          </p>
        </div>
      </div>
    </div>
  )
}
