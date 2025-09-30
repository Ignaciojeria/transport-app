/**
 * Utilidades para forzar sincronización con Electric SQL
 */

/**
 * Fuerza la limpieza del caché de Electric SQL
 * Esto puede requerir múltiples intentos y diferentes estrategias
 */
export const forceElectricSync = async (token: string, email: string): Promise<boolean> => {
  console.log('🔄 Forzando sincronización completa con Electric SQL...')
  
  try {
    // Estrategia 1: Limpiar caché local
    console.log('🧹 Paso 1: Limpiando caché local...')
    const { clearElectricCache } = await import('./electricCacheUtils')
    clearElectricCache()
    
    // Estrategia 2: Múltiples consultas con diferentes parámetros
    console.log('🔄 Paso 2: Realizando consultas de sincronización...')
    const syncPromises = []
    
    for (let i = 0; i < 3; i++) {
      const timestamp = Date.now() + i
      const randomId = Math.random().toString(36).substring(7)
      const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id,email&where=email='${email}'&offset=-1&_sync=${timestamp}&_attempt=${i}&_random=${randomId}`
      
      syncPromises.push(
        fetch(url, {
          headers: {
            'X-Access-Token': `Bearer ${token}`,
            'Cache-Control': 'no-cache',
            'Pragma': 'no-cache'
          }
        }).then(response => response.json())
      )
    }
    
    await Promise.all(syncPromises)
    console.log('✅ Consultas de sincronización completadas')
    
    // Estrategia 3: Esperar un poco para que se propague
    console.log('⏳ Paso 3: Esperando propagación...')
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    return true
  } catch (error) {
    console.error('❌ Error al forzar sincronización:', error)
    return false
  }
}

/**
 * Verifica si Electric SQL está sincronizado con la base de datos
 */
export const isElectricSynced = async (email: string): Promise<{
  synced: boolean
  electricData: any
  directData: any
  message: string
}> => {
  try {
    console.log('🔍 Verificando sincronización de Electric SQL...')
    
    // Importar funciones dinámicamente para evitar dependencias circulares
    const { checkAccountDirectly } = await import('./directDbCheck')
    const { findAccountByEmail } = await import('../services/electricService')
    
    // Obtener token del localStorage
    const authData = localStorage.getItem('transport_auth')
    if (!authData) {
      return {
        synced: false,
        electricData: null,
        directData: null,
        message: 'No hay token de autenticación'
      }
    }
    
    const { access_token } = JSON.parse(authData)
    
    // Verificar ambos lados
    const [electricResult, directResult] = await Promise.all([
      findAccountByEmail(access_token, email),
      checkAccountDirectly(email, access_token)
    ])
    
    const synced = (!!electricResult) === directResult.exists
    
    return {
      synced,
      electricData: electricResult,
      directData: directResult,
      message: synced 
        ? 'Electric SQL está sincronizado' 
        : 'Electric SQL NO está sincronizado'
    }
  } catch (error) {
    console.error('❌ Error al verificar sincronización:', error)
    return {
      synced: false,
      electricData: null,
      directData: null,
      message: `Error: ${error instanceof Error ? error.message : 'Desconocido'}`
    }
  }
}
