/**
 * Helpers para moneda según businessInfo.currency del menú.
 */

export type CurrencyCode = 'USD' | 'CLP' | 'BRL';

/** Locale por moneda para formateo de números */
const LOCALE_BY_CURRENCY: Record<CurrencyCode, string> = {
  CLP: 'es-CL',
  USD: 'en-US',
  BRL: 'pt-BR'
};

/** Decimales por moneda: CLP sin decimales, USD/BRL con 2 */
function getFractionDigits(currency: string): number {
  return currency === 'CLP' ? 0 : 2;
}

/**
 * Obtiene la moneda efectiva del negocio; si no existe retorna CLP por defecto.
 */
export function getEffectiveCurrency(restaurantData: { businessInfo?: { currency?: string } } | null | undefined): CurrencyCode {
  const code = restaurantData?.businessInfo?.currency?.toUpperCase();
  if (code === 'USD' || code === 'CLP' || code === 'BRL') return code;
  return 'CLP';
}

/**
 * Formatea un monto con la moneda del negocio (símbolo/código y separadores según locale).
 */
export function formatPrice(amount: number, currency: string = 'CLP'): string {
  const code = (currency?.toUpperCase() || 'CLP') as CurrencyCode;
  const locale = LOCALE_BY_CURRENCY[code] ?? 'es-CL';
  const fractionDigits = getFractionDigits(code);
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: code,
    minimumFractionDigits: fractionDigits,
    maximumFractionDigits: fractionDigits
  }).format(amount);
}
