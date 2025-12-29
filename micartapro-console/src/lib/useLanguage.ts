import { writable, derived, get } from 'svelte/store'
import type { Language } from './translations'
import { translations, languageNames, languageFlags } from './translations'

// Función para detectar el idioma del navegador
function detectBrowserLanguage(): Language {
  if (typeof window === 'undefined') return 'ES'
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('es')) return 'ES'
  if (browserLang.startsWith('pt')) return 'PT'
  if (browserLang.startsWith('en')) return 'EN'
  return 'ES' // Default to Spanish
}

// Store reactivo para el idioma
function createLanguageStore() {
  let initialLanguage: Language = 'ES'
  let initialLoading = true

  if (typeof window !== 'undefined') {
    // Leer del query param
    const urlParams = new URLSearchParams(window.location.search)
    const langParam = urlParams.get('lang') as Language
    
    // Leer de localStorage
    const savedLang = localStorage.getItem('preferred-language') as Language

    // Prioridad: query param > localStorage > browser > default
    if (langParam && ['EN', 'ES', 'PT'].includes(langParam)) {
      initialLanguage = langParam
      localStorage.setItem('preferred-language', langParam)
    } else if (savedLang && ['EN', 'ES', 'PT'].includes(savedLang)) {
      initialLanguage = savedLang
    } else {
      initialLanguage = detectBrowserLanguage()
      localStorage.setItem('preferred-language', initialLanguage)
    }
    
    initialLoading = false
  }

  const language = writable<Language>(initialLanguage)
  const loading = writable<boolean>(initialLoading)

  return {
    language,
    loading
  }
}

const languageStore = createLanguageStore()

// Función para inicializar el idioma (para compatibilidad)
export function initLanguage() {
  if (typeof window === 'undefined') {
    languageStore.language.set('ES')
    languageStore.loading.set(false)
    return
  }

  // Leer del query param
  const urlParams = new URLSearchParams(window.location.search)
  const langParam = urlParams.get('lang') as Language
  
  // Leer de localStorage
  const savedLang = localStorage.getItem('preferred-language') as Language

  // Prioridad: query param > localStorage > browser > default
  let newLanguage: Language = 'ES'
  if (langParam && ['EN', 'ES', 'PT'].includes(langParam)) {
    newLanguage = langParam
    localStorage.setItem('preferred-language', langParam)
  } else if (savedLang && ['EN', 'ES', 'PT'].includes(savedLang)) {
    newLanguage = savedLang
  } else {
    newLanguage = detectBrowserLanguage()
    localStorage.setItem('preferred-language', newLanguage)
  }
  
  languageStore.language.set(newLanguage)
  languageStore.loading.set(false)
}

// Función para cambiar el idioma
export function changeLanguage(newLanguage: Language) {
  if (!['EN', 'ES', 'PT'].includes(newLanguage)) return
  
  languageStore.language.set(newLanguage)
  if (typeof window !== 'undefined') {
    localStorage.setItem('preferred-language', newLanguage)
    
    // Actualizar URL sin recargar la página
    const url = new URL(window.location.href)
    url.searchParams.set('lang', newLanguage)
    window.history.replaceState({}, '', url.toString())
  }
}

// Exportar stores
export const language = languageStore.language
export const t = derived(languageStore.language, ($language) => translations[$language])
export const loading = languageStore.loading
export const availableLanguages = Object.keys(languageNames) as Language[]
export { languageNames, languageFlags }

// Inicializar al cargar el módulo
if (typeof window !== 'undefined') {
  initLanguage()
}

