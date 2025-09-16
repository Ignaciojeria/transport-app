import { useState, useEffect } from 'react'
import type { Language } from '../lib/translations'
import { translations } from '../lib/translations'

export function useLanguage() {
  // Leer idioma directamente de la URL cada vez que se llama
  const getLanguageFromURL = (): Language => {
    const urlParams = new URLSearchParams(window.location.search)
    const langParam = urlParams.get('lang')?.toUpperCase() as Language
    
    console.log('🌍 useLanguage: Reading from URL:', {
      search: window.location.search,
      langParam,
      href: window.location.href
    })
    
    if (langParam && ['CL', 'BR', 'EU'].includes(langParam)) {
      console.log('✅ useLanguage: Valid language found:', langParam)
      return langParam
    }
    
    console.log('🔄 useLanguage: No valid lang param, defaulting to CL')
    return 'CL'
  }

  const [language, setLanguage] = useState<Language>(getLanguageFromURL())

  useEffect(() => {
    const updateLanguage = () => {
      const newLang = getLanguageFromURL()
      if (newLang !== language) {
        console.log('🔄 useLanguage: Language changed:', language, '→', newLang)
        setLanguage(newLang)
      }
    }

    // Verificar cambios cada segundo (más agresivo)
    const interval = setInterval(updateLanguage, 1000)

    // Escuchar eventos de navegación
    const handlePopState = updateLanguage
    
    window.addEventListener('popstate', handlePopState)
    
    // Interceptar pushState y replaceState
    const originalPushState = history.pushState
    const originalReplaceState = history.replaceState
    
    history.pushState = function(...args) {
      originalPushState.apply(history, args)
      setTimeout(updateLanguage, 100)
    }
    
    history.replaceState = function(...args) {
      originalReplaceState.apply(history, args)
      setTimeout(updateLanguage, 100)
    }

    return () => {
      clearInterval(interval)
      window.removeEventListener('popstate', handlePopState)
      history.pushState = originalPushState
      history.replaceState = originalReplaceState
    }
  }, [language])

  const currentLanguage = getLanguageFromURL() // Siempre leer el más actual

  return {
    language: currentLanguage,
    t: translations[currentLanguage]
  }
}
