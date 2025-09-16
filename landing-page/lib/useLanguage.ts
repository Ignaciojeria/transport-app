'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter, usePathname } from 'next/navigation'
import { Language, translations, languageNames, languageFlags } from './translations'

export function useLanguage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const pathname = usePathname()
  
  const [language, setLanguage] = useState<Language>('CL')
  const [isLoading, setIsLoading] = useState(true)

  // Cargar idioma desde query params o localStorage
  useEffect(() => {
    const langParam = searchParams.get('lang') as Language
    const savedLang = localStorage.getItem('preferred-language') as Language
    
    if (langParam && ['CL', 'BR', 'EU'].includes(langParam)) {
      setLanguage(langParam)
      localStorage.setItem('preferred-language', langParam)
    } else if (savedLang && ['CL', 'BR', 'EU'].includes(savedLang)) {
      setLanguage(savedLang)
    } else {
      // Detectar idioma del navegador
      const browserLang = navigator.language.toLowerCase()
      if (browserLang.startsWith('pt')) {
        setLanguage('BR')
      } else if (browserLang.startsWith('en')) {
        setLanguage('EU')
      } else {
        setLanguage('CL') // Default a español
      }
    }
    
    setIsLoading(false)
  }, [searchParams])

  const changeLanguage = (newLanguage: Language) => {
    setLanguage(newLanguage)
    localStorage.setItem('preferred-language', newLanguage)
    
    // Actualizar URL con el nuevo idioma
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
