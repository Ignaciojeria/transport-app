import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Obtiene la URL del auth-ui según el entorno
 * - Desarrollo local: http://localhost:3003
 * - Producción: https://auth.micartapro.com
 */
export function getAuthUiUrl(): string {
  if (typeof window === 'undefined') {
    return 'https://auth.micartapro.com'
  }
  
  const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
  return isLocalDev ? 'http://localhost:3003' : 'https://auth.micartapro.com'
}

