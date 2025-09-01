// Modelos internos específicos para Gun.js
// Estos modelos evitan referencias circulares y objetos anidados

export interface GunDeliveryEvidence {
  // Datos de delivery (planos)
  delivery_status: string
  delivery_handledAt: string
  delivery_location_lat: number
  delivery_location_lng: number
  
  // Datos del receptor (planos)
  recipient_fullName: string
  recipient_nationalID: string
  
  // Datos de evidencia (planos)
  evidence_0_takenAt: string
  evidence_0_type: string
  evidence_0_url: string
  
  // Referencia
  orderReferenceID: string
}

export interface GunDeliveryFailure {
  // Datos de delivery (planos)
  delivery_status: string
  delivery_handledAt: string
  delivery_location_lat: number
  delivery_location_lng: number
  
  // Datos de fallo (planos)
  failure_reason: string
  failure_detail: string
  failure_referenceID: string
  
  // Datos de evidencia (planos)
  evidence_0_takenAt: string
  evidence_0_type: string
  evidence_0_url: string
  
  // Referencia
  orderReferenceID: string
}

export interface GunDeliveryStatus {
  status: 'delivered' | 'not-delivered'
  timestamp?: number
  deviceId?: string
}

export interface GunRouteInfo {
  routeId: string
  startedAt: string
  vehiclePlate: string
  driverInfo: string
  carrierInfo: string
}
