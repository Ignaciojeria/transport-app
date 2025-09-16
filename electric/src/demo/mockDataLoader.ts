import type { Route } from '../domain/route'
import type { Language } from '../lib/translations'
import { mockRouteData as mockRouteChile } from './mockData'
import { mockRouteBrazil } from './mockDataBrazil'
import { mockRouteUSA } from './mockDataUSA'

// Función para obtener mock data según el idioma
export const getMockDataByLanguage = (language: Language): Route => {
  switch (language) {
    case 'CL':
      return mockRouteChile as Route
    case 'BR':
      return mockRouteBrazil
    case 'EU':
      return mockRouteUSA
    default:
      return mockRouteChile as Route
  }
}

// Función para obtener información del país según el idioma
export const getCountryInfo = (language: Language) => {
  switch (language) {
    case 'CL':
      return {
        country: 'Chile',
        city: 'La Florida, Santiago',
        flag: '🇨🇱',
        currency: 'CLP',
        timeZone: 'America/Santiago'
      }
    case 'BR':
      return {
        country: 'Brasil',
        city: 'Ubatuba, São Paulo',
        flag: '🇧🇷',
        currency: 'BRL',
        timeZone: 'America/Sao_Paulo'
      }
    case 'EU':
      return {
        country: 'United States',
        city: 'Miami, Florida',
        flag: '🇺🇸',
        currency: 'USD',
        timeZone: 'America/New_York'
      }
    default:
      return {
        country: 'Chile',
        city: 'La Florida, Santiago',
        flag: '🇨🇱',
        currency: 'CLP',
        timeZone: 'America/Santiago'
      }
  }
}
