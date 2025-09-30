/**
 * Servicio para manejar la sincronización incremental correcta con Electric SQL
 * Basado en la documentación oficial de Electric SQL HTTP API
 */

export interface ElectricSyncState {
  offset: string | null
  isInitialized: boolean
  lastSync: Date | null
}

export interface ElectricSyncResult<T> {
  data: T | null
  newOffset: string | null
  hasMore: boolean
  error?: string
}

/**
 * Almacena el estado de sincronización en localStorage
 */
const SYNC_STATE_KEY = 'electric_sync_state'

/**
 * Obtiene el estado de sincronización guardado
 */
export const getSyncState = (shapeId: string): ElectricSyncState => {
  try {
    const stored = localStorage.getItem(`${SYNC_STATE_KEY}_${shapeId}`)
    if (stored) {
      const parsed = JSON.parse(stored)
      return {
        offset: parsed.offset,
        isInitialized: parsed.isInitialized || false,
        lastSync: parsed.lastSync ? new Date(parsed.lastSync) : null
      }
    }
  } catch (error) {
    console.error('Error al obtener estado de sincronización:', error)
  }
  
  return {
    offset: null,
    isInitialized: false,
    lastSync: null
  }
}

/**
 * Guarda el estado de sincronización
 */
export const saveSyncState = (shapeId: string, state: ElectricSyncState): void => {
  try {
    localStorage.setItem(`${SYNC_STATE_KEY}_${shapeId}`, JSON.stringify({
      ...state,
      lastSync: state.lastSync?.toISOString() || null
    }))
  } catch (error) {
    console.error('Error al guardar estado de sincronización:', error)
  }
}

/**
 * Realiza una consulta incremental a Electric SQL
 */
export const syncElectricShape = async <T>(
  shapeId: string,
  url: string,
  token: string,
  parser: (data: any) => T
): Promise<ElectricSyncResult<T>> => {
  try {
    const syncState = getSyncState(shapeId)
    
    // Determinar el offset a usar
    const offset = syncState.isInitialized ? syncState.offset : '-1'
    const syncUrl = `${url}&offset=${offset}`
    
    console.log(`🔄 Sincronizando shape ${shapeId} con offset: ${offset}`)
    
    const response = await fetch(syncUrl, {
      headers: {
        'X-Access-Token': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error(`Error en sincronización: ${response.status} ${response.statusText}`)
    }

    const data = await response.json()
    console.log(`📊 Respuesta de sincronización para ${shapeId}:`, data)

    // Extraer el nuevo offset de la respuesta
    // Electric SQL devuelve el offset en los headers de control
    const controlHeaders = data.find((item: any) => item.headers?.control === 'snapshot-end')
    const newOffset = controlHeaders?.headers?.xmax || null

    // Parsear los datos
    const parsedData = parser(data)

    // Actualizar estado de sincronización
    const newSyncState: ElectricSyncState = {
      offset: newOffset,
      isInitialized: true,
      lastSync: new Date()
    }
    saveSyncState(shapeId, newSyncState)

    return {
      data: parsedData,
      newOffset,
      hasMore: false, // Electric SQL maneja esto internamente
      error: undefined
    }

  } catch (error) {
    console.error(`❌ Error en sincronización de ${shapeId}:`, error)
    return {
      data: null,
      newOffset: null,
      hasMore: false,
      error: error instanceof Error ? error.message : 'Error desconocido'
    }
  }
}

/**
 * Parser específico para datos de cuentas
 */
export const parseAccountData = (data: any[]): any | null => {
  if (!Array.isArray(data) || data.length === 0) {
    return null
  }

  // Buscar el primer objeto que tenga value (no los de control)
  const accountData = data.find(item => item.value && item.value.email)
  return accountData ? accountData.value : null
}

/**
 * Parser específico para datos de tenants
 */
export const parseTenantsData = (data: any[]): any[] => {
  if (!Array.isArray(data) || data.length === 0) {
    return []
  }

  // Filtrar solo objetos con datos reales
  return data
    .filter(item => item.value && item.value.id)
    .map(item => item.value)
}

/**
 * Parser específico para datos de account_tenants
 */
export const parseAccountTenantsData = (data: any[]): any[] => {
  if (!Array.isArray(data) || data.length === 0) {
    return []
  }

  // Filtrar solo objetos con datos reales
  return data
    .filter(item => item.value && item.value.tenant_id)
    .map(item => item.value)
}

/**
 * Limpia el estado de sincronización para un shape específico
 */
export const clearSyncState = (shapeId: string): void => {
  try {
    localStorage.removeItem(`${SYNC_STATE_KEY}_${shapeId}`)
    console.log(`🧹 Estado de sincronización limpiado para ${shapeId}`)
  } catch (error) {
    console.error('Error al limpiar estado de sincronización:', error)
  }
}

/**
 * Limpia todos los estados de sincronización
 */
export const clearAllSyncStates = (): void => {
  try {
    const keys = Object.keys(localStorage)
    keys.forEach(key => {
      if (key.startsWith(SYNC_STATE_KEY)) {
        localStorage.removeItem(key)
      }
    })
    console.log('🧹 Todos los estados de sincronización limpiados')
  } catch (error) {
    console.error('Error al limpiar todos los estados:', error)
  }
}
