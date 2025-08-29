import Gun from 'gun'
import { useEffect, useState } from 'react'
import type { 
  DeliveryUnit, 
  Recipient, 
  DeliveryFailure,
  DeliveryLocation,
  DeliveryItem
} from '../domain/deliveries'

// Tipos para la interfaz de la aplicación (sin duplicar el dominio)
export interface DeliveryEvidence {
  recipient: Recipient
  photoDataUrl: string
  takenAt: number
  items?: DeliveryItem[]
  location?: DeliveryLocation
}

export interface NonDeliveryEvidence {
  reason: string
  observations?: string
  photoDataUrl: string
  takenAt: number
  location?: DeliveryLocation
  failure?: DeliveryFailure
}

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







// Exportar también la instancia de Gun por si necesitas funcionalidades avanzadas
export { gun, deliveriesData }
