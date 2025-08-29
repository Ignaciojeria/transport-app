import Gun from 'gun'
import { useEffect, useState } from 'react'
import type { 
  DeliveryUnit, 
  Recipient, 
  DeliveryFailure,
  DeliveryLocation,
  DeliveryItem
} from '../domain/deliveries'



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

// Función que crea una DeliveryUnit completa usando directamente la entidad de dominio
export function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  deliveryUnit: Partial<DeliveryUnit>
): void {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Asegurar que los campos requeridos estén presentes
  const completeDeliveryUnit: Partial<DeliveryUnit> = {
    ...deliveryUnit,
    orderReferenceID: deliveryUnit.orderReferenceID || `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(completeDeliveryUnit)
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
  
  const deliveryUnit: Partial<DeliveryUnit> = {
    recipient,
    evidencePhotos: [{
      takenAt: new Date().toISOString(),
      type: 'delivery',
      url: photoDataUrl,
    }],
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
  deliveryUnit: Partial<DeliveryUnit>
): void {
  const key = ndEvidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Asegurar que los campos requeridos estén presentes
  const completeDeliveryUnit: Partial<DeliveryUnit> = {
    ...deliveryUnit,
    orderReferenceID: deliveryUnit.orderReferenceID || `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
  
  deliveriesData.get(key).put(completeDeliveryUnit)
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
  return state[key] ?? undefined
}







// Exportar también la instancia de Gun por si necesitas funcionalidades avanzadas
export { gun, deliveriesData }
