'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter, usePathname } from 'next/navigation'
import { Language, translations, languageNames, languageFlags } from './translations'

export function useLanguage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const pathname = usePathname()
  
  const [language, setLanguage] = useState<Language>('EN') // English by default
  const [isLoading, setIsLoading] = useState(true)

  // Load language from query params or localStorage
  useEffect(() => {
    const langParam = searchParams.get('lang') as Language
    const savedLang = localStorage.getItem('preferred-language') as Language
    
    // Detect browser language if no param or saved language
    const detectBrowserLanguage = (): Language => {
      if (typeof window === 'undefined') return 'EN'
      const browserLang = navigator.language.toLowerCase()
      if (browserLang.startsWith('es')) return 'ES'
      if (browserLang.startsWith('pt')) return 'PT'
      return 'EN' // Default to English
    }
    
    if (langParam && ['EN', 'ES', 'PT'].includes(langParam)) {
      setLanguage(langParam)
      localStorage.setItem('preferred-language', langParam)
    } else if (savedLang && ['EN', 'ES', 'PT'].includes(savedLang)) {
      setLanguage(savedLang)
    } else {
      // Default to English, but can detect from browser
      const detected = detectBrowserLanguage()
      setLanguage(detected)
      localStorage.setItem('preferred-language', detected)
    }
    
    setIsLoading(false)
  }, [searchParams])

  const changeLanguage = (newLanguage: Language) => {
    setLanguage(newLanguage)
    localStorage.setItem('preferred-language', newLanguage)
    
    // Update URL with the new language
    const params = new URLSearchParams(searchParams.toString())
    params.set('lang', newLanguage)
    router.push(`${pathname}?${params.toString()}`)
  }

  const t = translations[language]
  const availableLanguages = Object.keys(languageNames) as Language[]

  return {
    language,
    changeLanguage,
    t,
    isLoading,
    availableLanguages,
    languageNames,
    languageFlags
  }
}

