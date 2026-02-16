'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { supabase } from '@/lib/supabase'
import { useLanguage } from '@/lib/useLanguage'
import { Language } from '@/lib/translations'
import { getOrCreateMenuId, createMenuWithSlug } from '@/lib/menuId'
import '@/lib/diagnoseMenuId' // Cargar función de diagnóstico en window

export default function AuthCallback() {
  const [status, setStatus] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const [showSlugForm, setShowSlugForm] = useState(false)
  const [slug, setSlug] = useState('')
  const [slugError, setSlugError] = useState('')
  const router = useRouter()
  const { t, isLoading: langLoading, language } = useLanguage()
  
  // Initialize status with translated text
  useEffect(() => {
    if (!langLoading) {
      setStatus(t.callback.processing)
    }
  }, [langLoading, t])

  useEffect(() => {
    // No ejecutar hasta que el idioma esté cargado
    if (langLoading) return

    const handleAuthCallback = async () => {
      try {
        // Leer el idioma directamente del query param o localStorage para asegurar que esté actualizado
        const urlParams = new URLSearchParams(window.location.search)
        const langParam = urlParams.get('lang') as Language
        const savedLang = typeof window !== 'undefined' 
          ? localStorage.getItem('preferred-language') as Language
          : null
        
        // Determinar el idioma a usar (prioridad: query param > localStorage > language del hook)
        const currentLanguage = (langParam && ['EN', 'ES', 'PT'].includes(langParam))
          ? langParam
          : (savedLang && ['EN', 'ES', 'PT'].includes(savedLang))
          ? savedLang
          : language

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
          
          let menuId: string | null
          try {
            await new Promise(resolve => setTimeout(resolve, 100))
            console.log('🔄 Obteniendo menuID para usuario:', session.user.id)
            menuId = await getOrCreateMenuId(session.user.id)
            if (menuId) {
              console.log('✅ MenuID obtenido:', menuId)
            } else {
              console.log('📝 Usuario nuevo, mostrar formulario de slug')
              setShowSlugForm(true)
              setIsLoading(false)
              return
            }
          } catch (menuError: any) {
            console.error('❌ Error al obtener menuID:', menuError)
            setStatus(`Error: ${menuError?.message || 'Error desconocido'}`)
            setIsLoading(false)
            setTimeout(() => router.push('/'), 5000)
            return
          }
          
          setStatus('Redirigiendo...')
          
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
          // Preservar el parámetro de idioma en la redirección (usar currentLanguage que leímos directamente)
          const redirectUrl = `${baseUrl}?lang=${currentLanguage}#auth=${encodedAuth}`
          
          console.log('🌐 Idioma detectado:', currentLanguage)
          console.log('🚀 Redirigiendo a:', redirectUrl)
          
          // Redirigir directamente a console después de un breve delay
          setTimeout(() => {
            window.location.href = redirectUrl
          }, 1000)
        } else {
          setStatus(t.callback.noSession)
          setTimeout(() => router.push('/'), 2000)
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
  }, [router, t, language, langLoading])

  const handleCreateWithSlug = async () => {
    const { data: { session } } = await supabase.auth.getSession()
    if (!session?.user || !session.access_token) return
    const s = slug.trim()
    if (!s) {
      setSlugError(t.callback.slugLabel || 'El slug es requerido')
      return
    }
    const normalized = s.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '')
    if (!normalized) {
      setSlugError(t.callback.slugHint || 'Solo letras, números y guiones')
      return
    }
    setSlugError('')
    setIsLoading(true)
    try {
      await createMenuWithSlug(session.user.id, normalized, session.access_token)
      const urlParams = new URLSearchParams(window.location.search)
      const langParam = urlParams.get('lang') as Language
      const savedLang = localStorage.getItem('preferred-language') as Language
      const currentLanguage = (langParam && ['EN', 'ES', 'PT'].includes(langParam))
        ? langParam
        : (savedLang && ['EN', 'ES', 'PT'].includes(savedLang))
        ? savedLang
        : language
      const userMetadata = session.user.user_metadata || {}
      const authData = {
        access_token: session.access_token,
        token_type: 'Bearer',
        expires_in: session.expires_in || 3600,
        refresh_token: session.refresh_token || '',
        user: {
          id: session.user.id,
          email: session.user.email || '',
          verified_email: !!session.user.email_confirmed_at,
          name: userMetadata.full_name || userMetadata.name || session.user.email?.split('@')[0] || '',
          given_name: userMetadata.given_name || userMetadata.name?.split(' ')[0] || '',
          family_name: userMetadata.family_name || userMetadata.name?.split(' ').slice(1).join(' ') || '',
          picture: userMetadata.avatar_url || userMetadata.picture || '',
          locale: userMetadata.locale || 'es',
        },
        timestamp: Date.now(),
        provider: 'supabase',
      }
      const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
      const baseUrl = isLocalDev ? 'http://localhost:5174' : 'https://console.micartapro.com'
      window.location.href = `${baseUrl}?lang=${currentLanguage}#auth=${btoa(JSON.stringify(authData))}`
    } catch (e: any) {
      setSlugError(e?.message === 'SLUG_EXISTS' ? 'Ese identificador ya está en uso' : (e?.message || 'Error al crear'))
      setIsLoading(false)
    }
  }

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

            {showSlugForm && (
              <div className="mb-6 text-left">
                <h3 className="text-sm font-semibold text-gray-800 mb-3">
                  {t.callback.chooseSlug ?? 'Elige el identificador de tu menú'}
                </h3>
                <label htmlFor="slug" className="block text-sm font-medium text-gray-700 mb-2">
                  {t.callback.slugLabel ?? 'Identificador (slug)'}
                </label>
                <input
                  id="slug"
                  type="text"
                  value={slug}
                  onChange={(e) => { setSlug(e.target.value); setSlugError('') }}
                  placeholder={t.callback.slugPlaceholder ?? 'ej: mi-restaurante'}
                  className="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none mb-2"
                  onKeyDown={(e) => e.key === 'Enter' && handleCreateWithSlug()}
                />
                <p className="text-xs text-gray-500 mb-3">
                  {t.callback.slugHint ?? 'Solo letras minúsculas, números y guiones.'}
                </p>
                {slugError && <p className="text-sm text-red-600 mb-2">{slugError}</p>}
                <button
                  type="button"
                  onClick={handleCreateWithSlug}
                  disabled={isLoading}
                  className="w-full py-3 rounded-xl font-semibold text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-60"
                >
                  {isLoading ? (t.callback.creating ?? 'Creando...') : (t.callback.create ?? 'Crear')}
                </button>
              </div>
            )}

            {isLoading && !showSlugForm && (
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-blue-600 h-2 rounded-full animate-pulse" style={{width: '75%'}}></div>
              </div>
            )}

            {!isLoading && !showSlugForm && (
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
