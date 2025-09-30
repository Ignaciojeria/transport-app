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
 * Valida si un offset tiene el formato correcto para Electric SQL
 * Basado en el error "has invalid format", parece que Electric SQL es muy estricto
 */
const isValidOffset = (offset: string): boolean => {
  // Electric SQL acepta offsets en formato de string
  // Debe ser un número válido o '-1' para sincronización inicial
  if (offset === '-1') return true
  if (offset === '0') return true
  
  // Verificar si es un número válido simple (solo enteros positivos)
  const num = Number(offset)
  if (!isNaN(num) && num > 0 && Number.isInteger(num)) return true
  
  // Verificar formato con guión bajo (ej: 0_0, 123_456)
  // Pero solo si ambas partes son enteros
  if (offset.includes('_')) {
    const parts = offset.split('_')
    if (parts.length === 2) {
      const [first, second] = parts
      const firstNum = Number(first)
      const secondNum = Number(second)
      return !isNaN(firstNum) && !isNaN(secondNum) && 
             Number.isInteger(firstNum) && Number.isInteger(secondNum) &&
             firstNum >= 0 && secondNum >= 0
    }
  }
  
  return false
}

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
 * Patrón correcto: offset=-1 para inicial, luego offset+live=true para continuar
 */
export const syncElectricShape = async <T>(
  shapeId: string,
  url: string,
  token: string,
  parser: (data: any) => T,
  useLive: boolean = false
): Promise<ElectricSyncResult<T>> => {
  try {
    const syncState = getSyncState(shapeId)
    
    // Determinar el offset a usar
    let offset = syncState.isInitialized ? syncState.offset : '-1'
    
    // Validar formato del offset
    if (offset && offset !== '-1' && !isValidOffset(offset)) {
      console.warn(`⚠️ Offset inválido detectado: ${offset}, reseteando a -1`)
      offset = '-1'
      // Limpiar estado de sincronización para este shape
      clearSyncState(shapeId)
    } else if (offset && offset !== '-1') {
      console.log(`✅ Offset válido detectado: ${offset}`)
    }
    
    // Construir URL según el patrón correcto
    let syncUrl = `${url}&offset=${offset}`
    if (useLive && offset !== '-1') {
      // Solo usar live=true para sincronización continua (no inicial)
      syncUrl += '&live=true'
    }
    
    console.log(`🔄 Sincronizando shape ${shapeId} con offset: ${offset}${useLive ? ' (live)' : ''}`)
    
    const response = await fetch(syncUrl, {
      headers: {
        'X-Access-Token': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      // Si es error 400 y el offset es inválido, resetear a sincronización inicial
      if (response.status === 400 && offset !== '-1') {
        console.warn(`⚠️ Error 400 con offset ${offset}, reseteando a sincronización inicial`)
        clearSyncState(shapeId)
        
        // Reintentar con offset -1
        const retryUrl = `${url}&offset=-1`
        const retryResponse = await fetch(retryUrl, {
          headers: {
            'X-Access-Token': `Bearer ${token}`
          }
        })
        
        if (!retryResponse.ok) {
          throw new Error(`Error en sincronización inicial: ${retryResponse.status} ${retryResponse.statusText}`)
        }
        
        // Procesar respuesta de reintento
        const retryData = await retryResponse.json()
        const retryParsedData = parser(retryData)
        
        // Obtener nuevo offset de la respuesta de reintento
        const retryControlHeaders = retryData.find((item: any) => item.headers?.control === 'snapshot-end')
        const retryNewOffset = retryControlHeaders?.headers?.xmax ? String(retryControlHeaders.headers.xmax) : null
        
        // Actualizar estado con el nuevo offset
        const newSyncState: ElectricSyncState = {
          offset: retryNewOffset,
          isInitialized: true,
          lastSync: new Date()
        }
        saveSyncState(shapeId, newSyncState)
        
        return {
          data: retryParsedData,
          newOffset: retryNewOffset,
          hasMore: false,
          error: undefined
        }
      }
      
      throw new Error(`Error en sincronización: ${response.status} ${response.statusText}`)
    }

    const data = await response.json()
    console.log(`📊 Respuesta de sincronización para ${shapeId}:`, data)

    // Extraer el nuevo offset de la respuesta
    // Electric SQL devuelve el offset en los headers de control
    const controlHeaders = data.find((item: any) => item.headers?.control === 'snapshot-end')
    let newOffset = null
    
    // Buscar el offset en los headers de control
    if (controlHeaders?.headers?.xmax) {
      newOffset = String(controlHeaders.headers.xmax)
    }
    
    // También verificar si hay un electric-offset header en la respuesta HTTP
    const electricOffsetHeader = response.headers.get('electric-offset')
    if (electricOffsetHeader) {
      newOffset = electricOffsetHeader
      console.log(`📊 Offset encontrado en header: ${newOffset}`)
    }

    // Verificar mensajes de control
    const upToDateMessage = data.find((item: any) => item.headers?.control === 'up-to-date')
    const mustRefetchMessage = data.find((item: any) => item.headers?.control === 'must-refetch')
    
    if (mustRefetchMessage) {
      console.warn(`⚠️ Mensaje must-refetch recibido para ${shapeId}, limpiando estado`)
      clearSyncState(shapeId)
      // Retornar error para que el cliente re-sync desde cero
      return {
        data: null,
        newOffset: null,
        hasMore: false,
        error: 'must-refetch'
      }
    }
    
    // Parsear los datos
    const parsedData = parser(data)
    
    // Determinar si hay más datos (paginación)
    const hasMore = !upToDateMessage && newOffset !== null

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
      hasMore,
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

/**
 * Limpia todos los offsets inválidos y resetea a sincronización inicial
 */
export const clearInvalidOffsets = (): void => {
  try {
    const keys = Object.keys(localStorage)
    let cleanedCount = 0
    
    keys.forEach(key => {
      if (key.startsWith(SYNC_STATE_KEY)) {
        try {
          const stored = localStorage.getItem(key)
          if (stored) {
            const parsed = JSON.parse(stored)
            if (parsed.offset && !isValidOffset(parsed.offset)) {
              console.log(`🧹 Limpiando offset inválido: ${parsed.offset} en ${key}`)
              localStorage.removeItem(key)
              cleanedCount++
            }
          }
        } catch (error) {
          // Si no se puede parsear, limpiar también
          console.log(`🧹 Limpiando estado corrupto: ${key}`)
          localStorage.removeItem(key)
          cleanedCount++
        }
      }
    })
    
    console.log(`✅ Limpiados ${cleanedCount} offsets inválidos`)
  } catch (error) {
    console.error('Error al limpiar offsets inválidos:', error)
  }
}

/**
 * Función de prueba para validar formatos de offset
 * Útil para debug y verificar que los offsets se manejan correctamente
 */
export const testOffsetValidation = (): void => {
  const testOffsets = [
    '-1',      // Sincronización inicial
    '0',       // Offset cero
    '123',     // Número simple
    '0_0',     // Formato con guión bajo
    '123_456', // Formato con guión bajo
    '301131',  // El offset que causó el error
    'invalid', // Offset inválido
    '1_2_3',   // Formato inválido (más de 2 partes)
    'a_b',     // Formato inválido (no numérico)
  ]
  
  console.log('🧪 Probando validación de offsets:')
  testOffsets.forEach(offset => {
    const isValid = isValidOffset(offset)
    console.log(`  ${offset}: ${isValid ? '✅ VÁLIDO' : '❌ INVÁLIDO'}`)
  })
}

/**
 * Limpia offsets que han causado errores 400 en el pasado
 * Basado en patrones observados en los logs de error
 */
export const clearProblematicOffsets = (): void => {
  try {
    const keys = Object.keys(localStorage)
    let cleanedCount = 0
    
    keys.forEach(key => {
      if (key.startsWith(SYNC_STATE_KEY)) {
        try {
          const stored = localStorage.getItem(key)
          if (stored) {
            const parsed = JSON.parse(stored)
            if (parsed.offset) {
              // Limpiar offsets que sabemos que causan problemas
              // Basado en los errores reales observados
              const shouldClean = 
                parsed.offset === '301131' ||  // Error 400 observado
                parsed.offset === '0_0' ||    // Error 400 observado
                !isValidOffset(parsed.offset) // Cualquier offset inválido
              
              if (shouldClean) {
                console.log(`🧹 Limpiando offset problemático: ${parsed.offset} en ${key}`)
                localStorage.removeItem(key)
                cleanedCount++
              }
            }
          }
        } catch (error) {
          // Si no se puede parsear, limpiar también
          console.log(`🧹 Limpiando estado corrupto: ${key}`)
          localStorage.removeItem(key)
          cleanedCount++
        }
      }
    })
    
    console.log(`✅ Limpiados ${cleanedCount} offsets problemáticos`)
  } catch (error) {
    console.error('Error al limpiar offsets problemáticos:', error)
  }
}

/**
 * Sincronización inicial: offset=-1 (sin live)
 */
export const syncElectricShapeInitial = async <T>(
  shapeId: string,
  url: string,
  token: string,
  parser: (data: any) => T
): Promise<ElectricSyncResult<T>> => {
  console.log(`🔄 Sincronización inicial para shape ${shapeId}`)
  return syncElectricShape(shapeId, url, token, parser, false)
}

/**
 * Sincronización continua: offset=último + live=true
 */
export const syncElectricShapeLive = async <T>(
  shapeId: string,
  url: string,
  token: string,
  parser: (data: any) => T
): Promise<ElectricSyncResult<T>> => {
  console.log(`🔄 Sincronización continua (live) para shape ${shapeId}`)
  return syncElectricShape(shapeId, url, token, parser, true)
}

/**
 * Sincronización completa con paginación automática
 * Maneja la sincronización inicial y la paginación hasta obtener todos los datos
 */
export const syncElectricShapeComplete = async <T>(
  shapeId: string,
  url: string,
  token: string,
  parser: (data: any) => T
): Promise<ElectricSyncResult<T>> => {
  console.log(`🔄 Sincronización completa para shape ${shapeId}`)
  
  let allData: T | null = null
  let currentOffset: string | null = null
  let hasMore = true
  
  // Sincronización inicial
  const initialResult = await syncElectricShapeInitial(shapeId, url, token, parser)
  
  if (initialResult.error === 'must-refetch') {
    // Si hay must-refetch, limpiar todo y empezar de nuevo
    clearSyncState(shapeId)
    return await syncElectricShapeComplete(shapeId, url, token, parser)
  }
  
  if (initialResult.error) {
    return initialResult
  }
  
  allData = initialResult.data
  currentOffset = initialResult.newOffset
  hasMore = initialResult.hasMore
  
  // Continuar con paginación si es necesario
  while (hasMore && currentOffset) {
    console.log(`🔄 Continuando paginación con offset: ${currentOffset}`)
    
    const paginatedResult = await syncElectricShape(shapeId, url, token, parser, false)
    
    if (paginatedResult.error === 'must-refetch') {
      // Si hay must-refetch durante paginación, limpiar todo y empezar de nuevo
      clearSyncState(shapeId)
      return await syncElectricShapeComplete(shapeId, url, token, parser)
    }
    
    if (paginatedResult.error) {
      console.error(`❌ Error en paginación: ${paginatedResult.error}`)
      break
    }
    
    // Combinar datos (esto depende del tipo de datos)
    if (paginatedResult.data) {
      // Para arrays, concatenar
      if (Array.isArray(allData) && Array.isArray(paginatedResult.data)) {
        allData = [...allData, ...paginatedResult.data] as T
      }
      // Para objetos individuales, usar el último
      else {
        allData = paginatedResult.data
      }
    }
    
    currentOffset = paginatedResult.newOffset
    hasMore = paginatedResult.hasMore
  }
  
  console.log(`✅ Sincronización completa finalizada para ${shapeId}`)
  
  return {
    data: allData,
    newOffset: currentOffset,
    hasMore: false,
    error: undefined
  }
}
