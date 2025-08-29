import { z } from 'zod'
import Gun from 'gun'
import { useEffect, useState } from 'react'

// Esquemas iguales que antes
const DeliveryEvidence = z.object({
  recipientName: z.string().min(1),
  recipientRut: z.string().min(1),
  photoDataUrl: z.string().min(10),
  takenAt: z.number(),
})

const NonDeliveryEvidence = z.object({
  reason: z.string().min(1),
  observations: z.string().optional().default(''),
  photoDataUrl: z.string().min(10),
  takenAt: z.number(),
})

export type DeliveryEvidence = z.infer<typeof DeliveryEvidence>
export type NonDeliveryEvidence = z.infer<typeof NonDeliveryEvidence>

// Configuración mejorada de Gun para persistencia local robusta
// Incluye múltiples capas de persistencia para garantizar que los datos offline se mantengan
const gun = Gun({
  radisk: true,  // Habilita Radix para persistencia local
  localStorage: true,  // Habilita localStorage como respaldo adicional
  peers: [
    'https://peer.wallie.io/gun',
  ],
  // Configuración adicional para mejor persistencia local
  axe: false,  // Desactiva el algoritmo de consenso para mejor rendimiento local
  multicast: false,  // Desactiva multicast para evitar conflictos en entornos locales
})

// Namespace para datos del driver
const driverData = gun.get('driver-state')

// Sistema de respaldo local adicional usando localStorage
class LocalBackup {
  private static readonly STORAGE_KEY = 'gun-driver-state-backup'
  
  static save(key: string, value: any): void {
    try {
      const backup = this.loadAll()
      backup[key] = {
        value,
        timestamp: Date.now(),
        version: '1.0'
      }
      localStorage.setItem(this.STORAGE_KEY, JSON.stringify(backup))
    } catch (error) {
      console.warn('Error saving to localStorage backup:', error)
    }
  }
  
  static load(key: string): any {
    try {
      const backup = this.loadAll()
      const item = backup[key]
      return item ? item.value : null
    } catch (error) {
      console.warn('Error loading from localStorage backup:', error)
      return null
    }
  }
  
  static loadAll(): Record<string, any> {
    try {
      const stored = localStorage.getItem(this.STORAGE_KEY)
      return stored ? JSON.parse(stored) : {}
    } catch (error) {
      console.warn('Error loading localStorage backup:', error)
      return {}
    }
  }
  
  static remove(key: string): void {
    try {
      const backup = this.loadAll()
      delete backup[key]
      localStorage.setItem(this.STORAGE_KEY, JSON.stringify(backup))
    } catch (error) {
      console.warn('Error removing from localStorage backup:', error)
    }
  }
  
  static clear(): void {
    try {
      localStorage.removeItem(this.STORAGE_KEY)
    } catch (error) {
      console.warn('Error clearing localStorage backup:', error)
    }
  }
}

// Helpers para claves (igual que antes)
export const routeStartedKey = (routeId: string) => `routeStarted:${routeId}`
export const routeLicenseKey = (routeId: string) => `routeLicense:${routeId}`
export const deliveryKey = (routeId: string, vIdx: number, oIdx: number, uIdx: number) =>
  `delivery:${routeId}:${vIdx}-${oIdx}-${uIdx}`
export const evidenceKey = (routeId: string, vIdx: number, oIdx: number, uIdx: number) =>
  `evidence:${routeId}:${vIdx}-${oIdx}-${uIdx}`
export const ndEvidenceKey = (routeId: string, vIdx: number, oIdx: number, uIdx: number) =>
  `nd-evidence:${routeId}:${vIdx}-${oIdx}-${uIdx}`

// Helper para generar un ID único del dispositivo
function getDeviceId(): string {
  let deviceId = localStorage.getItem('gun-device-id')
  if (!deviceId) {
    deviceId = `device-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    localStorage.setItem('gun-device-id', deviceId)
  }
  return deviceId
}

// Helper para obtener información del dispositivo actual
export function getDeviceInfo() {
  return {
    id: getDeviceId(),
    userAgent: navigator.userAgent,
    timestamp: Date.now(),
    online: navigator.onLine
  }
}

// Hook reactivo para escuchar cambios en Gun (similar a useLiveQuery)
export function useGunData(key?: string) {
  const [data, setData] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!key) {
      setLoading(false)
      return
    }

    // Escuchar cambios en tiempo real
    const ref = driverData.get(key)
    const unsubscribe = ref.on((value, _key) => {
      setData(value)
      setLoading(false)
    })

    return () => {
      if (unsubscribe && typeof unsubscribe.off === 'function') {
        unsubscribe.off()
      }
    }
  }, [key])

  return { data, loading }
}

// Hook para escuchar todos los cambios del estado del driver
export function useDriverState() {
  const [state, setState] = useState<Record<string, any>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Cargar datos del respaldo local al inicializar
    const loadLocalBackup = () => {
      try {
        const localData = LocalBackup.loadAll()
        const initialState: Record<string, any> = {}
        
        Object.entries(localData).forEach(([key, item]) => {
          if (item && item.value !== null && item.value !== undefined) {
            initialState[key] = item.value
          }
        })
        
        if (Object.keys(initialState).length > 0) {
          setState(initialState)
          console.log('Datos del respaldo local cargados:', initialState)
        }
      } catch (error) {
        console.warn('Error cargando respaldo local:', error)
      }
    }

    // Cargar respaldo local inmediatamente
    loadLocalBackup()

    // Escuchar todos los cambios en el namespace del driver
    const unsubscribe = driverData.map().on((value, key) => {
      if (value !== null && value !== undefined) {
        setState(prev => ({ ...prev, [key]: value }))
      } else {
        setState(prev => {
          const newState = { ...prev }
          delete newState[key]
          return newState
        })
      }
      setLoading(false)
    })

    return () => {
      if (unsubscribe && typeof unsubscribe.off === 'function') {
        unsubscribe.off()
      }
    }
  }, [])

  return { data: { s: state }, loading }
}

// Mutadores - API mejorada con respaldo local robusto
export function setRouteStarted(routeId: string, started: boolean) {
  const key = routeStartedKey(routeId)
  const timestamp = Date.now()
  const deviceId = getDeviceId()
  
  // Guardar con metadatos para sincronización
  const value = {
    status: started ? 'true' : 'false',
    timestamp,
    deviceId,
    action: started ? 'route_started' : 'route_stopped'
  }
  
  // Guardar en GunJS (con persistencia Radix)
  driverData.get(key).put(value)
  
  // También guardar en key simple para compatibilidad
  driverData.get(`${key}_simple`).put(started ? 'true' : 'false')
  
  // Respaldo local adicional para garantizar persistencia
  LocalBackup.save(key, value)
  LocalBackup.save(`${key}_simple`, started ? 'true' : 'false')
}

export function setRouteLicense(routeId: string, license: string) {
  const key = routeLicenseKey(routeId)
  const timestamp = Date.now()
  const deviceId = getDeviceId()
  
  // Guardar con metadatos para sincronización
  const value = {
    license: license,
    timestamp,
    deviceId,
    action: 'license_set'
  }
  
  // Guardar en GunJS (con persistencia Radix)
  driverData.get(key).put(value)
  
  // Respaldo local adicional para garantizar persistencia
  LocalBackup.save(key, value)
}

export function setDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  status: 'delivered' | 'not-delivered'
) {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Guardar en GunJS (con persistencia Radix)
  driverData.get(key).put(status)
  
  // Respaldo local adicional para garantizar persistencia
  LocalBackup.save(key, status)
}

export function getDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): 'delivered' | 'not-delivered' | undefined {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Primero intentar obtener del respaldo local (síncrono)
  const localResult = LocalBackup.load(key)
  if (localResult !== null) {
    return localResult
  }
  
  // Si no hay respaldo local, intentar obtener de GunJS
  let result: any = undefined
  
  // Gun es asíncrono, pero podemos intentar obtener del cache local
  driverData.get(key).once((value) => {
    result = value
  })
  
  return result ?? undefined
}

export function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: DeliveryEvidence
) {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Guardar en GunJS (con persistencia Radix)
  driverData.get(key).put(evidence)
  
  // Respaldo local adicional para garantizar persistencia
  LocalBackup.save(key, evidence)
}

export function getDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): DeliveryEvidence | undefined {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  let result: any = undefined
  
  driverData.get(key).once((value) => {
    result = value
  })
  
  return result ? DeliveryEvidence.parse(result) : undefined
}

export function setNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: NonDeliveryEvidence
) {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Guardar en GunJS (con persistencia Radix)
  driverData.get(key).put(evidence)
  
  // Respaldo local adicional para garantizar persistencia
  LocalBackup.save(key, evidence)
}

export function getNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): NonDeliveryEvidence | undefined {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  let result: any = undefined
  
  driverData.get(key).once((value) => {
    result = value
  })
  
  return result ? NonDeliveryEvidence.parse(result) : undefined
}

// Helper mejorado para obtener estado de delivery usando el estado reactivo
export function getDeliveryStatusFromState(
  state: Record<string, any>,
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): 'delivered' | 'not-delivered' | undefined {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  return state[key] ?? undefined
}

export function getRouteLicenseFromState(
  state: Record<string, any>,
  routeId: string
): string | undefined {
  const key = routeLicenseKey(routeId)
  const licenseData = state[key]
  return licenseData?.license ?? undefined
}

// Hook específico para monitorear sincronización de ruta iniciada
export function useRouteStartedSync(routeId: string) {
  const [syncData, setSyncData] = useState<{
    isStarted: boolean
    lastAction: string
    deviceId: string
    timestamp: number
    syncedDevices: string[]
  } | null>(null)

  useEffect(() => {
    const key = routeStartedKey(routeId)
    
    // Escuchar cambios con metadatos
    const unsubscribe = driverData.get(key).on((data) => {
      if (data && typeof data === 'object') {
        setSyncData({
          isStarted: data.status === 'true',
          lastAction: data.action || 'unknown',
          deviceId: data.deviceId || 'unknown',
          timestamp: data.timestamp || 0,
          syncedDevices: [data.deviceId || 'unknown'] // En un MVP simple, solo mostramos el último dispositivo
        })
      }
    })

    return () => {
      if (unsubscribe && typeof unsubscribe.off === 'function') {
        unsubscribe.off()
      }
    }
  }, [routeId])

  return syncData
}

// Función para obtener información de sincronización de todas las rutas
export function getAllRoutesSyncInfo() {
  return new Promise((resolve) => {
    const routes: Record<string, any> = {}
    
    // Escuchar todos los cambios relacionados con rutas
    driverData.map().once((data, key) => {
      if (key && key.includes('routeStarted:') && !key.includes('_simple')) {
        const routeId = key.replace('routeStarted:', '')
        routes[routeId] = data
      }
    })
    
    setTimeout(() => resolve(routes), 500) // Dar tiempo para recopilar datos
  })
}

// Función para sincronizar datos del respaldo local con GunJS
export function syncLocalBackupToGun() {
  try {
    const localData = LocalBackup.loadAll()
    let syncedCount = 0
    
    Object.entries(localData).forEach(([key, item]) => {
      if (item && item.value !== null && item.value !== undefined) {
        // Sincronizar con GunJS
        driverData.get(key).put(item.value)
        syncedCount++
      }
    })
    
    console.log(`Sincronizados ${syncedCount} elementos del respaldo local con GunJS`)
    return syncedCount
  } catch (error) {
    console.error('Error sincronizando respaldo local:', error)
    return 0
  }
}

// Función para limpiar respaldo local después de sincronización exitosa
export function clearLocalBackup() {
  LocalBackup.clear()
  console.log('Respaldo local limpiado después de sincronización')
}

// Exportar también la instancia de Gun por si necesitas funcionalidades avanzadas
export { gun, driverData }
