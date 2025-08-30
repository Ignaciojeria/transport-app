import Gun from 'gun'
import { useEffect, useState } from 'react'
import type { 
  DeliveryUnit, 
  Recipient, 
  DeliveryFailure,
  DeliveryLocation,
  DeliveryItem,
  EvidencePhoto
} from '../domain/deliveries'



// Configuración de Gun para MVP con sincronización entre dispositivos
const gun = Gun({
  radisk: true,
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
        // Debug: log para claves de delivery
        if (key.includes('delivery:')) {
          // console.log(`🔄 Estado local actualizado - Clave: ${key}`, value) // Comentado para reducir logs
          // console.log(`🔄 Tipo de valor:`, typeof value) // Comentado para reducir logs
          if (typeof value === 'object') {
            // console.log(`🔄 Propiedades del objeto:`, Object.keys(value)) // Comentado para reducir logs
            if (value.failure) {
              // console.log(`🔄 DeliveryFailure encontrado:`, value.failure) // Comentado para reducir logs
            }
          }
        }
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
  status: 'delivered' | 'not-delivered',
  evidence?: {
    reason?: string
    observations?: string
    photoDataUrl?: string
  }
) {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // console.log(`🔧 setDeliveryStatus llamado:`, { routeId, visitIndex, orderIndex, unitIndex, status, evidence }) // Comentado para reducir logs
  
  if (status === 'not-delivered' && evidence) {
    // Para no entregas, usar el dominio DeliveryFailure
    const deliveryFailure: DeliveryFailure = {
      reason: evidence.reason || '',
      detail: evidence.observations || '',
      referenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`
    }
    
    const deliveryData = {
      status,
      failure: deliveryFailure,
      photoDataUrl: evidence.photoDataUrl,
      timestamp: Date.now(),
      deviceId: getDeviceId()
    }
    
    // console.log(`💾 Guardando datos de no entrega en clave: ${key}`) // Comentado para reducir logs
    // console.log(`💾 Datos guardados:`, deliveryData) // Comentado para reducir logs
    
    deliveriesData.get(key).put(deliveryData)
  } else {
    // Para entregas exitosas, solo guardar estado
    // console.log(`💾 Guardando estado simple: ${status} en clave: ${key}`) // Comentado para reducir logs
    deliveriesData.get(key).put(status)
  }
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

// Función que crea una DeliveryUnit completa usando directamente la entidad de dominio
export function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: {
    recipientName: string
    recipientRut: string
    photoDataUrl: string
    takenAt: number
  }
): void {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Crear Recipient usando el dominio correcto
  const recipient: Recipient = {
    fullName: evidence.recipientName,
    nationalID: evidence.recipientRut
  }
  
  // Crear evidencia de foto usando el dominio correcto
  const evidencePhoto: EvidencePhoto = {
    takenAt: new Date(evidence.takenAt).toISOString(),
    type: 'delivery',
    url: evidence.photoDataUrl,
  }
  
  // Crear estructura de DeliveryUnit usando el dominio correcto
  const deliveryUnit: Partial<DeliveryUnit> = {
    recipient,
    evidencePhotos: [evidencePhoto],
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(deliveryUnit)
}

export function getDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): Promise<Partial<DeliveryUnit> | undefined> {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  return new Promise((resolve) => {
    deliveriesData.get(key).once((value) => {
      resolve(value || undefined)
    })
  })
}

// Función específica para gestionar entrega exitosa
export function setSuccessfulDelivery(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  recipient: Recipient,
  photoDataUrl: string,
  items?: DeliveryItem[],
  location?: DeliveryLocation
): void {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Crear evidencia de foto usando el dominio correcto
  const evidencePhoto: EvidencePhoto = {
    takenAt: new Date().toISOString(),
    type: 'delivery',
    url: photoDataUrl,
  }
  
  const deliveryUnit: Partial<DeliveryUnit> = {
    recipient,
    evidencePhotos: [evidencePhoto],
    items: items || [],
    delivery: {
      status: 'delivered',
      handledAt: new Date().toISOString(),
      location: location || { latitude: 0, longitude: 0 },
    },
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(deliveryUnit)
}

// Función específica para gestionar no entrega
export function setFailedDelivery(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: {
    reason: string
    observations: string
    photoDataUrl: string
  }
): void {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Crear DeliveryFailure usando el dominio correcto
  const deliveryFailure: DeliveryFailure = {
    reason: evidence.reason,
    detail: evidence.observations,
    referenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`
  }
  
  // Crear evidencia de foto
  const evidencePhoto: EvidencePhoto = {
    takenAt: new Date().toISOString(),
    type: 'non-delivery',
    url: evidence.photoDataUrl,
  }
  
  // Crear estructura de DeliveryUnit usando el dominio correcto
  const deliveryUnit: Partial<DeliveryUnit> = {
    delivery: {
      status: 'not-delivered',
      handledAt: new Date().toISOString(),
      location: { latitude: 0, longitude: 0 },
      failure: deliveryFailure
    },
    evidencePhotos: [evidencePhoto],
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(deliveryUnit)
}

export function getNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): Promise<Partial<DeliveryUnit> | undefined> {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  return new Promise((resolve) => {
    deliveriesData.get(key).once((value) => {
      resolve(value || undefined)
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
  const data = state[key]
  
  // Debug: logs removidos para limpiar la consola
  
  if (typeof data === 'string') {
    // Estado simple (formato anterior)
    return data as 'delivered' | 'not-delivered'
  } else if (data && typeof data === 'object' && data.status) {
    // Estado con evidencia (nuevo formato)
    return data.status as 'delivered' | 'not-delivered'
  }
  
  return undefined
}

// Helper para obtener evidencia de no entrega desde el estado
export function getNonDeliveryEvidenceFromState(
  state: Record<string, any>,
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): {
  reason?: string
  observations?: string
  photoDataUrl?: string
  timestamp?: number
  deviceId?: string
} | null {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  const data = state[key]
  
  if (data && typeof data === 'object' && data.status === 'not-delivered') {
    return {
      reason: data.failure?.reason,
      observations: data.failure?.detail,
      photoDataUrl: data.photoDataUrl,
      timestamp: data.timestamp,
      deviceId: data.deviceId
    }
  }
  
  return null
}







// Exportar también la instancia de Gun por si necesitas funcionalidades avanzadas
export { gun, deliveriesData }
