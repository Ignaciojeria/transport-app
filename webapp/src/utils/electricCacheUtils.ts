/**
 * Utilidades para manejar el caché de Electric SQL
 */

/**
 * Limpia el caché local de Electric SQL
 * Esto puede ayudar cuando los datos no se sincronizan correctamente
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
 * Fuerza una sincronización completa con Electric SQL
 * Esto puede ser útil cuando los datos no están actualizados
 */
export const forceElectricSync = async (): Promise<void> => {
  try {
    console.log('🔄 Forzando sincronización con Electric SQL...')
    
    // Limpiar caché primero
    clearElectricCache()
    
    // Recargar la página para forzar una nueva sincronización
    window.location.reload()
  } catch (error) {
    console.error('❌ Error al forzar sincronización:', error)
  }
}

/**
 * Verifica si hay datos en caché de Electric SQL
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
