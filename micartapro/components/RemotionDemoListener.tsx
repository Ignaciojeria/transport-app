'use client'

import { useEffect } from 'react'

/**
 * En modo demo (?demo=1), escucha postMessage de Remotion para simular clics.
 * Usado cuando la landing se muestra en un iframe dentro del video de marketing.
 */
export function RemotionDemoListener() {
  useEffect(() => {
    if (typeof window === 'undefined') return
    const params = new URLSearchParams(window.location.search)
    if (params.get('demo') !== '1') return

    const handler = (e: MessageEvent) => {
      const allowedOrigins = ['http://localhost:', 'http://127.0.0.1:', 'https://localhost:', 'https://127.0.0.1:']
      if (!allowedOrigins.some(o => e.origin?.startsWith(o))) return

      if (e.data?.clickIniciaGratis) {
        // Buscar el primer botón que contenga "Inicia Gratis", "Start Free" o "Comece Grátis"
        const btns = Array.from(document.querySelectorAll('button, a'))
        const btn = btns.find(b => {
          const text = (b.textContent || '').toLowerCase()
          return text.includes('inicia gratis') || text.includes('start free') || text.includes('comece grátis')
        })
        if (btn) {
          btn.classList.add('ring-4', 'ring-orange-400', 'ring-opacity-75')
          setTimeout(() => btn.classList.remove('ring-4', 'ring-orange-400', 'ring-opacity-75'), 200)
          btn.click() // Redirige a auth-ui/consola
        }
      }
    }

    window.addEventListener('message', handler)
    return () => window.removeEventListener('message', handler)
  }, [])

  return null
}
