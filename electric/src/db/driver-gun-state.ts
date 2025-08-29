import { z } from 'zod'
import Gun from 'gun'
import { useEffect, useState } from 'react'
import type { 
  DeliveryUnit, 
  EvidencePhoto, 
  Recipient, 
  DeliveryFailure,
  DeliveryLocation,
  DeliveryItem
} from '../domain/deliveries'

// Esquemas de Zod basados en las entidades existentes
const DeliveryEvidence = z.object({
  recipient: z.object({
    fullName: z.string().min(1),
    nationalID: z.string().min(1),
  }),
  photoDataUrl: z.string().min(10),
  takenAt: z.number(),
  items: z.array(z.object({
    sku: z.string(),
    description: z.string(),
    quantity: z.number(),
    deliveredQuantity: z.number(),
  })).optional(),
})

const NonDeliveryEvidence = z.object({
  reason: z.string().min(1),
  observations: z.string().optional().default(''),
  photoDataUrl: z.string().min(10),
  takenAt: z.number(),
  failure: z.object({
    detail: z.string(),
    reason: z.string(),
    referenceID: z.string(),
  }).optional(),
})

export type DeliveryEvidence = z.infer<typeof DeliveryEvidence>
export type NonDeliveryEvidence = z.infer<typeof NonDeliveryEvidence>

// Configuración de Gun para MVP con sincronización entre dispositivos
const gun = Gun({
  localStorage: true,
  radisk: false,
  peers: [
    'https://peer.wallie.io/gun',
  ]
})

// Namespace para datos del driver
const driverData = gun.get('driver-state')

// Helpers para claves
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

// Hook reactivo para escuchar cambios en Gun
export function useGunData(key?: string) {
  const [data, setData] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!key) {
      setLoading(false)
      return
    }

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

// Mutadores usando las entidades de deliveries.ts
export function setRouteStarted(routeId: string, started: boolean) {
  const key = routeStartedKey(routeId)
  const timestamp = Date.now()
  const deviceId = getDeviceId()
  
  const value = {
    status: started ? 'true' : 'false',
    timestamp,
    deviceId,
    action: started ? 'route_started' : 'route_stopped'
  }
  
  driverData.get(key).put(value)
  driverData.get(`${key}_simple`).put(started ? 'true' : 'false')
}

export function setRouteLicense(routeId: string, license: string) {
  const key = routeLicenseKey(routeId)
  const timestamp = Date.now()
  const deviceId = getDeviceId()
  
  const value = {
    license: license,
    timestamp,
    deviceId,
    action: 'license_set'
  }
  
  driverData.get(key).put(value)
}

export function setDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  status: 'delivered' | 'not-delivered'
) {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  driverData.get(key).put(status)
}

export function getDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): 'delivered' | 'not-delivered' | undefined {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  let result: any = undefined
  
  driverData.get(key).once((value) => {
    result = value
  })
  
  return result ?? undefined
}

// Función mejorada que usa la entidad DeliveryUnit
export function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: DeliveryEvidence
) {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Convertir a formato compatible con DeliveryUnit
  const evidenceData = {
    ...evidence,
    // Convertir timestamp a ISO string para compatibilidad
    takenAt: new Date(evidence.takenAt).toISOString(),
    // Crear EvidencePhoto compatible
    evidencePhotos: [{
      takenAt: new Date(evidence.takenAt).toISOString(),
      type: 'delivery',
      url: evidence.photoDataUrl,
    }],
  }
  
  driverData.get(key).put(evidenceData)
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

// Función mejorada que usa la entidad DeliveryFailure
export function setNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: NonDeliveryEvidence
) {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Convertir a formato compatible con las entidades existentes
  const evidenceData = {
    ...evidence,
    takenAt: new Date(evidence.takenAt).toISOString(),
    // Crear DeliveryFailure si no existe
    failure: evidence.failure || {
      detail: evidence.observations || '',
      reason: evidence.reason,
      referenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
    },
  }
  
  driverData.get(key).put(evidenceData)
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

// Helper para obtener estado de delivery usando el estado reactivo
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
    
    const unsubscribe = driverData.get(key).on((data) => {
      if (data && typeof data === 'object') {
        setSyncData({
          isStarted: data.status === 'true',
          lastAction: data.action || 'unknown',
          deviceId: data.deviceId || 'unknown',
          timestamp: data.timestamp || 0,
          syncedDevices: [data.deviceId || 'unknown']
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
    
    driverData.map().once((data, key) => {
      if (key && key.includes('routeStarted:') && !key.includes('_simple')) {
        const routeId = key.replace('routeStarted:', '')
        routes[routeId] = data
      }
    })
    
    setTimeout(() => resolve(routes), 500)
  })
}

// Exportar también la instancia de Gun por si necesitas funcionalidades avanzadas
export { gun, driverData }
