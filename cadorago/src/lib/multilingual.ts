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
 * @param text - Texto multiidioma o array de textos (description)
 * @param lang - Idioma deseado (opcional, usa el idioma actual si no se especifica)
 * @returns El texto en el idioma especificado o el texto base
 */
export function getMultilingualText(
  text: MultilingualText | MultilingualText[] | undefined | null,
  lang?: Language
): string {
  if (!text) return '';

  if (Array.isArray(text)) {
    const targetLang = lang || (typeof window !== 'undefined' ? get(language) : 'ES');
    const parts = text
      .map((t) => getMultilingualText(t, targetLang))
      .filter((s) => s !== '');
    return parts.join(' ');
  }

  const targetLang = lang || (typeof window !== 'undefined' ? get(language) : 'ES');
  const langKey = targetLang.toLowerCase() as 'es' | 'en' | 'pt';
  if (text.languages?.[langKey]) return text.languages[langKey]!;
  return text.base || '';
}

/**
 * Obtiene el texto base de un texto multiidioma o string
 */
export function getBaseText(text: MultilingualText | string | undefined | null): string {
  if (!text) return '';
  if (typeof text === 'string') return text;
  return text.base || '';
}
