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

// Esquemas de Zod basados en las entidades del dominio
const DeliveryEvidenceSchema = z.object({
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
  location: z.object({
    latitude: z.number(),
    longitude: z.number(),
  }).optional(),
})

const NonDeliveryEvidenceSchema = z.object({
  reason: z.string().min(1),
  observations: z.string().optional().default(''),
  photoDataUrl: z.string().min(10),
  takenAt: z.number(),
  location: z.object({
    latitude: z.number(),
    longitude: z.number(),
  }).optional(),
  failure: z.object({
    detail: z.string(),
    reason: z.string(),
    referenceID: z.string(),
  }).optional(),
})

const RouteStartedSchema = z.object({
  status: z.union([z.literal('true'), z.literal('false')]),
  timestamp: z.number(),
  deviceId: z.string(),
  action: z.enum(['route_started', 'route_stopped']),
  routeId: z.string(),
})

const RouteLicenseSchema = z.object({
  license: z.string(),
  timestamp: z.number(),
  deviceId: z.string(),
  action: z.literal('license_set'),
  routeId: z.string(),
})

export type DeliveryEvidence = z.infer<typeof DeliveryEvidenceSchema>
export type NonDeliveryEvidence = z.infer<typeof NonDeliveryEvidenceSchema>
export type RouteStarted = z.infer<typeof RouteStartedSchema>
export type RouteLicense = z.infer<typeof RouteLicenseSchema>

// Configuración de Gun para MVP con sincronización entre dispositivos
const gun = Gun({
  localStorage: true,
  radisk: false,
  peers: [
    'https://peer.wallie.io/gun',
  ]
})

// Namespace para datos de deliveries
const deliveriesData = gun.get('deliveries-state')

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

    const ref = deliveriesData.get(key)
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

// Hook para escuchar todos los cambios del estado de deliveries
export function useDeliveriesState() {
  const [state, setState] = useState<Record<string, any>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const unsubscribe = deliveriesData.map().on((value, key) => {
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
export function setRouteStarted(routeId: string, started: boolean): void {
  const key = routeStartedKey(routeId)
  const timestamp = Date.now()
  const deviceId = getDeviceId()
  
  const value: RouteStarted = {
    status: started ? 'true' : 'false',
    timestamp,
    deviceId,
    action: started ? 'route_started' : 'route_stopped',
    routeId
  }
  
  deliveriesData.get(key).put(value)
  // Mantener versión simple para compatibilidad
  deliveriesData.get(`${key}_simple`).put(started ? 'true' : 'false')
}

export function setRouteLicense(routeId: string, license: string): void {
  const key = routeLicenseKey(routeId)
  const timestamp = Date.now()
  const deviceId = getDeviceId()
  
  const value: RouteLicense = {
    license,
    timestamp,
    deviceId,
    action: 'license_set',
    routeId
  }
  
  deliveriesData.get(key).put(value)
}

export function setDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  status: 'delivered' | 'not-delivered'
) {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  deliveriesData.get(key).put(status)
}

export function getDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): Promise<'delivered' | 'not-delivered' | undefined> {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  
  return new Promise((resolve) => {
    deliveriesData.get(key).once((value) => {
      resolve(value ?? undefined)
    })
  })
}

// Función mejorada que crea una DeliveryUnit completa
export function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: DeliveryEvidence
): void {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Crear DeliveryUnit compatible con el dominio
  const deliveryUnit: Partial<DeliveryUnit> = {
    recipient: evidence.recipient,
    evidencePhotos: [{
      takenAt: new Date(evidence.takenAt).toISOString(),
      type: 'delivery',
      url: evidence.photoDataUrl,
    }],
    items: evidence.items || [],
    delivery: {
      status: 'delivered',
      handledAt: new Date(evidence.takenAt).toISOString(),
      location: evidence.location || { latitude: 0, longitude: 0 },
    },
    // Campos requeridos que podrían necesitar ser proporcionados
    businessIdentifiers: {
      commerce: '',
      consumer: '',
    },
    lpn: '',
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(deliveryUnit)
}

export function getDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): Promise<DeliveryEvidence | undefined> {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  return new Promise((resolve) => {
    deliveriesData.get(key).once((value) => {
      if (value) {
        try {
          // Convertir DeliveryUnit de vuelta a DeliveryEvidence
          const evidence: DeliveryEvidence = {
            recipient: value.recipient,
            photoDataUrl: value.evidencePhotos?.[0]?.url || '',
            takenAt: new Date(value.delivery?.handledAt || Date.now()).getTime(),
            items: value.items,
            location: value.delivery?.location,
          }
          resolve(evidence)
        } catch (error) {
          console.error('Error parsing delivery evidence:', error)
          resolve(undefined)
        }
      } else {
        resolve(undefined)
      }
    })
  })
}

// Función mejorada que crea una DeliveryUnit con failure para no entrega
export function setNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: NonDeliveryEvidence
): void {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Crear DeliveryUnit con failure para no entrega
  const deliveryUnit: Partial<DeliveryUnit> = {
    recipient: {
      fullName: 'N/A',
      nationalID: 'N/A',
    },
    evidencePhotos: [{
      takenAt: new Date(evidence.takenAt).toISOString(),
      type: 'non-delivery',
      url: evidence.photoDataUrl,
    }],
    items: [],
    delivery: {
      status: 'not-delivered',
      handledAt: new Date(evidence.takenAt).toISOString(),
      location: evidence.location || { latitude: 0, longitude: 0 },
      failure: evidence.failure || {
        detail: evidence.observations || '',
        reason: evidence.reason,
        referenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
      },
    },
    businessIdentifiers: {
      commerce: '',
      consumer: '',
    },
    lpn: '',
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(deliveryUnit)
}

export function getNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): Promise<NonDeliveryEvidence | undefined> {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  return new Promise((resolve) => {
    deliveriesData.get(key).once((value) => {
      if (value) {
        try {
          // Convertir DeliveryUnit de vuelta a NonDeliveryEvidence
          const evidence: NonDeliveryEvidence = {
            reason: value.delivery?.failure?.reason || '',
            observations: value.delivery?.failure?.detail || '',
            photoDataUrl: value.evidencePhotos?.[0]?.url || '',
            takenAt: new Date(value.delivery?.handledAt || Date.now()).getTime(),
            location: value.delivery?.location,
            failure: value.delivery?.failure,
          }
          resolve(evidence)
        } catch (error) {
          console.error('Error parsing non-delivery evidence:', error)
          resolve(undefined)
        }
      } else {
        resolve(undefined)
      }
    })
  })
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
    
    const unsubscribe = deliveriesData.get(key).on((data) => {
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
export function getAllRoutesSyncInfo(): Promise<Record<string, RouteStarted>> {
  return new Promise((resolve) => {
    const routes: Record<string, RouteStarted> = {}
    
    deliveriesData.map().once((data, key) => {
      if (key && key.includes('routeStarted:') && !key.includes('_simple')) {
        const routeId = key.replace('routeStarted:', '')
        if (data && typeof data === 'object') {
          routes[routeId] = data as RouteStarted
        }
      }
    })
    
    setTimeout(() => resolve(routes), 500)
  })
}

// Exportar también la instancia de Gun por si necesitas funcionalidades avanzadas
export { gun, deliveriesData }
