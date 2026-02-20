'use client'

import { useEffect } from 'react'
import { getDemoAuthData } from '@/lib/demoAuth'

/**
 * En modo demo (?demo=1), escucha postMessage de Remotion para simular clics.
 * Acepta accessToken vía postMessage: { clickIniciarSesion: true, accessToken: "eyJ..." }
 * Si viene el token, redirige a la consola. Si no, solo efecto visual.
 */
export function RemotionDemoListener() {
  useEffect(() => {
    if (typeof window === 'undefined') return
    const params = new URLSearchParams(window.location.search)
    if (params.get('demo') !== '1') return

    const handler = (e: MessageEvent) => {
      const allowedOrigins = ['http://localhost:', 'http://127.0.0.1:', 'https://localhost:', 'https://127.0.0.1:']
      if (!allowedOrigins.some(o => e.origin?.startsWith(o))) return

      if (e.data?.clickIniciarSesion) {
        const accessToken = e.data.accessToken as string | undefined
        const refreshToken = e.data.refreshToken as string | undefined

        const btns = Array.from(document.querySelectorAll('button'))
        const btn = btns.find(b => {
          const text = (b.textContent || '').toLowerCase()
          return (
            text.includes('iniciar sesión') ||
            text.includes('sign in') ||
            text.includes('entrar com') ||
            text.includes('iniciar sesión con google') ||
            text.includes('sign in with google')
          )
        })
        if (btn) {
          btn.classList.add('ring-4', 'ring-orange-400', 'ring-opacity-75')
          setTimeout(() => btn.classList.remove('ring-4', 'ring-orange-400', 'ring-opacity-75'), 400)

          if (accessToken) {
            const authData = getDemoAuthData(accessToken, refreshToken)
            const encodedAuth = btoa(JSON.stringify(authData))
            const lang = params.get('lang') || 'ES'
            const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
            const baseUrl = isLocalDev ? 'http://localhost:5174' : 'https://console.micartapro.com'
            const redirectUrl = `${baseUrl}?lang=${lang}#auth=${encodedAuth}`
            window.location.href = redirectUrl
          }
        }
      }
    }

    window.addEventListener('message', handler)
    return () => window.removeEventListener('message', handler)
  }, [])

  return null
}
