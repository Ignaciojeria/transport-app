/**
 * Configuración de la aplicación
 */

// Detectar si estamos en desarrollo local
const isLocalDev = typeof window !== 'undefined' && 
  (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')

// URL del backend - se detecta automáticamente según el entorno
// En local: http://localhost:8082
// En producción: https://micartapro-backend-27303662337.us-central1.run.app
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 
  (isLocalDev ? 'http://localhost:8082' : 'https://micartapro-backend-27303662337.us-central1.run.app')

