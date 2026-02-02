import { get } from 'svelte/store';
import { language } from './useLanguage';
import type { Language } from './translations';

/**
 * Tipo para texto multiidioma
 */
export type MultilingualText = {
  base: string;
  languages?: {
    es?: string;
    en?: string;
    pt?: string;
  };
};

/**
 * Obtiene el texto en el idioma actual, o el texto base si no está disponible
 * @param text - Texto multiidioma o string simple (para retrocompatibilidad)
 * @param lang - Idioma deseado (opcional, usa el idioma actual si no se especifica)
 * @returns El texto en el idioma especificado o el texto base
 */
export function getMultilingualText(text: MultilingualText | string | undefined | null, lang?: Language): string {
  if (!text) return '';
  
  // Si es un string simple (retrocompatibilidad), devolverlo directamente
  if (typeof text === 'string') {
    return text;
  }
  
  // Si es un objeto MultilingualText
  const targetLang = lang || (typeof window !== 'undefined' ? get(language) : 'ES');
  const langKey = targetLang.toLowerCase() as 'es' | 'en' | 'pt';
  
  // Intentar obtener el texto en el idioma especificado
  if (text.languages && text.languages[langKey]) {
    return text.languages[langKey]!;
  }
  
  // Si no está disponible, usar el texto base
  return text.base || '';
}

/**
 * Obtiene el texto base de un texto multiidioma
 * @param text - Texto multiidioma o string simple
 * @returns El texto base o el string original
 */
export function getBaseText(text: MultilingualText | string | undefined | null): string {
  if (!text) return '';
  if (typeof text === 'string') return text;
  
  // Si es un objeto, verificar que tenga la propiedad base
  if (typeof text === 'object' && text !== null) {
    if ('base' in text && typeof text.base === 'string') {
      return text.base;
    }
  }
  
  return '';
}
