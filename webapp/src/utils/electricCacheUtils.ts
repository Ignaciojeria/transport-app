/**
 * Utilidades para manejar el caché de Electric SQL
 * Con LiveQuery, la sincronización es automática, pero estas utilidades
 * siguen siendo útiles para debug y casos edge
 */

/**
 * Limpia el caché local de Electric SQL
 * Útil para debug o cuando LiveQuery no funciona correctamente
 */
export const clearElectricCache = (): void => {
  try {
    // Limpiar localStorage relacionado con Electric SQL
    const keysToRemove: string[] = []
    
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && (
        key.includes('electric') || 
        key.includes('electric-sql') ||
        key.includes('transport_auth')
      )) {
        keysToRemove.push(key)
      }
    }
    
    keysToRemove.forEach(key => {
      console.log('🧹 Eliminando del caché:', key)
      localStorage.removeItem(key)
    })
    
    console.log('✅ Caché de Electric SQL limpiado')
  } catch (error) {
    console.error('❌ Error al limpiar caché de Electric SQL:', error)
  }
}

/**
 * Verifica si hay datos en caché de Electric SQL
 * Útil para debug
 */
export const hasElectricCache = (): boolean => {
  try {
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.includes('electric')) {
        return true
      }
    }
    return false
  } catch (error) {
    console.error('❌ Error al verificar caché:', error)
    return false
  }
}

/**
 * Obtiene información del caché de Electric SQL para debugging
 * Útil para diagnosticar problemas de sincronización
 */
export const getElectricCacheInfo = (): { keys: string[], size: number } => {
  try {
    const electricKeys: string[] = []
    let totalSize = 0
    
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.includes('electric')) {
        electricKeys.push(key)
        const value = localStorage.getItem(key)
        if (value) {
          totalSize += value.length
        }
      }
    }
    
    return {
      keys: electricKeys,
      size: totalSize
    }
  } catch (error) {
    console.error('❌ Error al obtener info del caché:', error)
    return { keys: [], size: 0 }
  }
}

/**
 * Fuerza una recarga completa de la aplicación
 * Útil como último recurso cuando LiveQuery no funciona
 */
export const forceAppReload = (): void => {
  console.log('🔄 Forzando recarga completa de la aplicación...')
  clearElectricCache()
  window.location.reload()
}