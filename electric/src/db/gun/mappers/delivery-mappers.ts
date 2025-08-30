import type { 
  DeliveryUnit
} from '../../../domain/deliveries'
import type { 
  GunDeliveryEvidence, 
  GunDeliveryFailure
} from '../models/delivery-models'

// Mappers para convertir entre dominio y modelos internos de Gun.js

export function mapDeliveryUnitToGun(
  deliveryUnit: Partial<DeliveryUnit>
): GunDeliveryEvidence {
  return {
    // Datos de delivery (planos)
    delivery_status: deliveryUnit.delivery?.status || 'delivered',
    delivery_handledAt: deliveryUnit.delivery?.handledAt || new Date().toISOString(),
    delivery_location_lat: deliveryUnit.delivery?.location?.latitude || 0,
    delivery_location_lng: deliveryUnit.delivery?.location?.longitude || 0,
    
    // Datos del receptor (planos)
    recipient_fullName: deliveryUnit.recipient?.fullName || '',
    recipient_nationalID: deliveryUnit.recipient?.nationalID || '',
    
    // Datos de evidencia (planos)
    evidence_0_takenAt: deliveryUnit.evidencePhotos?.[0]?.takenAt || new Date().toISOString(),
    evidence_0_type: deliveryUnit.evidencePhotos?.[0]?.type || 'delivery',
    evidence_0_url: deliveryUnit.evidencePhotos?.[0]?.url || '',
    
    // Referencia
    orderReferenceID: deliveryUnit.orderReferenceID || ''
  }
}

export function mapGunToDeliveryUnit(
  gunData: GunDeliveryEvidence
): Partial<DeliveryUnit> {
  return {
    delivery: {
      status: gunData.delivery_status as 'delivered' | 'not-delivered',
      handledAt: gunData.delivery_handledAt,
      location: {
        latitude: gunData.delivery_location_lat,
        longitude: gunData.delivery_location_lng
      }
    },
    recipient: {
      fullName: gunData.recipient_fullName,
      nationalID: gunData.recipient_nationalID
    },
    evidencePhotos: [{
      takenAt: gunData.evidence_0_takenAt,
      type: gunData.evidence_0_type as 'delivery' | 'non-delivery',
      url: gunData.evidence_0_url
    }],
    orderReferenceID: gunData.orderReferenceID
  }
}

export function mapDeliveryFailureToGun(
  deliveryUnit: Partial<DeliveryUnit>
): GunDeliveryFailure {
  return {
    // Datos de delivery (planos)
    delivery_status: deliveryUnit.delivery?.status || 'not-delivered',
    delivery_handledAt: deliveryUnit.delivery?.handledAt || new Date().toISOString(),
    delivery_location_lat: deliveryUnit.delivery?.location?.latitude || 0,
    delivery_location_lng: deliveryUnit.delivery?.location?.longitude || 0,
    
    // Datos de fallo (planos)
    failure_reason: deliveryUnit.delivery?.failure?.reason || '',
    failure_detail: deliveryUnit.delivery?.failure?.detail || '',
    failure_referenceID: deliveryUnit.delivery?.failure?.referenceID || '',
    
    // Datos de evidencia (planos)
    evidence_0_takenAt: deliveryUnit.evidencePhotos?.[0]?.takenAt || new Date().toISOString(),
    evidence_0_type: deliveryUnit.evidencePhotos?.[0]?.type || 'non-delivery',
    evidence_0_url: deliveryUnit.evidencePhotos?.[0]?.url || '',
    
    // Referencia
    orderReferenceID: deliveryUnit.orderReferenceID || ''
  }
}

export function mapGunToDeliveryFailure(
  gunData: GunDeliveryFailure
): Partial<DeliveryUnit> {
  return {
    delivery: {
      status: gunData.delivery_status as 'delivered' | 'not-delivered',
      handledAt: gunData.delivery_handledAt,
      location: {
        latitude: gunData.delivery_location_lat,
        longitude: gunData.delivery_location_lng
      },
      failure: {
        reason: gunData.failure_reason,
        detail: gunData.failure_detail,
        referenceID: gunData.failure_referenceID
      }
    },
    evidencePhotos: [{
      takenAt: gunData.evidence_0_takenAt,
      type: gunData.evidence_0_type as 'delivery' | 'non-delivery',
      url: gunData.evidence_0_url
    }],
    orderReferenceID: gunData.orderReferenceID
  }
}

// Helper para crear entidades del dominio desde datos simples
export function createDeliveryUnitFromEvidence(
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
): Partial<DeliveryUnit> {
  return {
    delivery: {
      status: 'delivered',
      handledAt: new Date().toISOString(),
      location: { latitude: 0, longitude: 0 }
    },
    recipient: {
      fullName: evidence.recipientName,
      nationalID: evidence.recipientRut
    },
    evidencePhotos: [{
      takenAt: new Date(evidence.takenAt).toISOString(),
      type: 'delivery',
      url: evidence.photoDataUrl,
    }],
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
}

export function createDeliveryUnitFromFailure(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: {
    reason: string
    observations: string
    photoDataUrl: string
  }
): Partial<DeliveryUnit> {
  return {
    delivery: {
      status: 'not-delivered',
      handledAt: new Date().toISOString(),
      location: { latitude: 0, longitude: 0 },
      failure: {
        reason: evidence.reason,
        detail: evidence.observations,
        referenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`
      }
    },
    evidencePhotos: [{
      takenAt: new Date().toISOString(),
      type: 'non-delivery',
      url: evidence.photoDataUrl,
    }],
    orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
  }
}
