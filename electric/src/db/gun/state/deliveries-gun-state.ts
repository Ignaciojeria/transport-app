import Gun from 'gun'
import 'gun/lib/radix'
import 'gun/lib/radisk'
import 'gun/lib/store'
import 'gun/lib/rindexed'
import { useEffect, useState } from 'react'
import type { 
  DeliveryUnit, 
  Recipient, 
  DeliveryLocation,
  DeliveryItem
} from '../../../domain/deliveries'
import type { 
  GunDeliveryEvidence, 
  GunDeliveryFailure
} from '../models/delivery-models'
import {
  mapDeliveryUnitToGun,
  mapGunToDeliveryUnit,
  mapDeliveryFailureToGun,
  mapGunToDeliveryFailure,
  createDeliveryUnitFromEvidence,
  createDeliveryUnitFromFailure
} from '../mappers/delivery-mappers'

// Configuración de Gun usando RAD (Radix Adaptive Database)
const gun = Gun({
  radisk: true, // Habilita RAD
  localStorage: false, // Deshabilitar localStorage para evitar cuotas
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
// Convierte internamente datos de Gun.js a entidades del dominio
export function useDeliveriesState() {
  const [state, setState] = useState<Record<string, any>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const unsubscribe = deliveriesData.map().on((value, key) => {
      if (value !== null && value !== undefined) {
        // Convertir datos de Gun.js al dominio ANTES de guardar en el estado
        const domainData = convertGunDataToDomain(key, value)
        setState(prev => ({ ...prev, [key]: domainData }))
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

// Función interna para convertir datos de Gun.js al dominio
function convertGunDataToDomain(key: string, value: any): any {
  // Si es una clave de evidence, convertir usando el mapper
  if (key.includes('evidence:')) {
    try {
      return mapGunToDeliveryUnit(value)
    } catch (error) {
      console.warn('Error convirtiendo evidence a dominio:', error)
      return value // Fallback a valor original
    }
  }
  
  // Si es una clave de nd-evidence, convertir usando el mapper
  if (key.includes('nd-evidence:')) {
    try {
      return mapGunToDeliveryFailure(value)
    } catch (error) {
      console.warn('Error convirtiendo nd-evidence a dominio:', error)
      return value // Fallback a valor original
    }
  }
  
  // Para otras claves, devolver el valor tal como está
  return value
}

// Mutadores usando los mappers y modelos internos

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
  
  if (status === 'not-delivered' && evidence) {
    // Para no entregas, crear entidad del dominio y mapear a Gun.js
    const deliveryUnit = createDeliveryUnitFromFailure(
      routeId, visitIndex, orderIndex, unitIndex, {
        reason: evidence.reason || '',
        observations: evidence.observations || '',
        photoDataUrl: evidence.photoDataUrl || ''
      }
    )
    
    const gunData = mapDeliveryFailureToGun(deliveryUnit)
    
    deliveriesData.get(key).put(gunData)
  } else if (status === 'delivered') {
    // Para entregas exitosas, solo guardar estado
    // La evidencia detallada ya se guardó en setDeliveryEvidence
    deliveriesData.get(key).put(status)
  } else {
    // Para otros estados, guardar estado simple
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

// Función que recibe entidades del dominio y las persiste en Gun.js
export function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  deliveryUnit: Partial<DeliveryUnit>
): void {
  const key = evidenceKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Mapear entidad del dominio a modelo interno de Gun.js
  const gunData = mapDeliveryUnitToGun(deliveryUnit)
  
  // Debug: ver qué estamos guardando
  console.log(`💾 setDeliveryEvidence - Datos a guardar:`, gunData)
  
  // Guardar evidencia detallada
  deliveriesData.get(key).put(gunData)
  
  // También actualizar el estado de delivery para mantener consistencia
  const deliveryStateKey = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  deliveriesData.get(deliveryStateKey).put('delivered')
  
  console.log(`✅ setDeliveryEvidence completado para ${key}`)
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
      if (value && typeof value === 'object') {
        // Mapear desde modelo interno de Gun.js a entidad del dominio
        const deliveryUnit = mapGunToDeliveryUnit(value as GunDeliveryEvidence)
        resolve(deliveryUnit)
      } else {
        resolve(undefined)
      }
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
  
  // Crear entidad del dominio
  const deliveryUnit: Partial<DeliveryUnit> = {
    delivery: {
      status: 'delivered',
      handledAt: new Date().toISOString(),
      location: location || { latitude: 0, longitude: 0 }
    },
    recipient: {
      fullName: recipient.fullName,
      nationalID: recipient.nationalID
    },
    evidencePhotos: [{
      takenAt: new Date().toISOString(),
      type: 'delivery',
      url: photoDataUrl
    }],
    items: items || [],
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`
  }
  
  // Mapear a modelo interno de Gun.js
  const gunData = mapDeliveryUnitToGun(deliveryUnit)
  
  deliveriesData.get(key).put(gunData)
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
  
  // Crear entidad del dominio usando el helper
  const deliveryUnit = createDeliveryUnitFromFailure(
    routeId, visitIndex, orderIndex, unitIndex, evidence
  )
  
  // Mapear a modelo interno de Gun.js
  const gunData = mapDeliveryFailureToGun(deliveryUnit)
  
  deliveriesData.get(key).put(gunData)
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
      if (value && typeof value === 'object') {
        // Mapear desde modelo interno de Gun.js a entidad del dominio
        const deliveryUnit = mapGunToDeliveryFailure(value as GunDeliveryFailure)
        resolve(deliveryUnit)
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
  const data = state[key]
  
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

// Función unificada que recibe entidades del dominio directamente
export function setDeliveryUnitByEntity(
  deliveryUnit: Partial<DeliveryUnit> & { 
    routeId: string
    visitIndex: number
    orderIndex: number
    unitIndex: number
  }
): void {
  const { routeId, visitIndex, orderIndex, unitIndex, ...unitData } = deliveryUnit
  
  // Usar el mapper apropiado según el estado
  let gunData: any
  
  if (unitData.delivery?.status === 'not-delivered') {
    // Para no entregas, usar el mapper de fallo
    gunData = mapDeliveryFailureToGun(unitData)
  } else {
    // Para entregas exitosas, usar el mapper normal
    gunData = mapDeliveryUnitToGun(unitData)
  }
  
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  
  // Guardar en un solo lugar usando la clave principal
  deliveriesData.get(key).put(gunData)
  
  console.log(`✅ setDeliveryUnitByEntity completado para ${key} con estado: ${unitData.delivery?.status}`)
}
