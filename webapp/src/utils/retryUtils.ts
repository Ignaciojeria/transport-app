/**
 * Utilidades para manejar reintentos con backoff exponencial
 */

export interface RetryOptions {
  maxAttempts?: number
  baseDelay?: number
  maxDelay?: number
  backoffMultiplier?: number
}

/**
 * Ejecuta una función con reintentos automáticos y backoff exponencial
 * @param fn - Función a ejecutar
 * @param options - Opciones de reintento
 * @returns Promise que resuelve con el resultado de la función
 */
export const retryWithBackoff = async <T>(
  fn: () => Promise<T>,
  options: RetryOptions = {}
): Promise<T> => {
  const {
    maxAttempts = 3,
    baseDelay = 1000,
    maxDelay = 10000,
    backoffMultiplier = 2
  } = options

  let lastError: Error | null = null

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      console.log(`🔄 Intento ${attempt}/${maxAttempts}`)
      const result = await fn()
      console.log(`✅ Éxito en intento ${attempt}`)
      return result
    } catch (error) {
      lastError = error instanceof Error ? error : new Error('Error desconocido')
      console.warn(`❌ Intento ${attempt} falló:`, lastError.message)

      if (attempt === maxAttempts) {
        console.error(`💥 Todos los intentos fallaron después de ${maxAttempts} intentos`)
        throw lastError
      }

      // Calcular delay con backoff exponencial
      const delay = Math.min(
        baseDelay * Math.pow(backoffMultiplier, attempt - 1),
        maxDelay
      )

      console.log(`⏳ Esperando ${delay}ms antes del siguiente intento...`)
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }

  throw lastError || new Error('Error en reintentos')
}

/**
 * Función específica para sincronizar con Electric SQL
 * @param syncFn - Función de sincronización
 * @returns Promise que resuelve cuando la sincronización es exitosa
 */
export const syncWithElectric = async (syncFn: () => Promise<void>): Promise<void> => {
  return retryWithBackoff(syncFn, {
    maxAttempts: 5,
    baseDelay: 2000,
    maxDelay: 15000,
    backoffMultiplier: 1.5
  })
}
