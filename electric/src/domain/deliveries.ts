export interface DeliveryEvent {
  carrier: Carrier
  deliveryUnits: DeliveryUnit[]
  driver: Driver
  manualChange?: ManualChange
  route: RouteReference
  vehicle: Vehicle
}

export interface Carrier {
  name: string
  nationalID: string
}

export interface BusinessIdentifiers {
  commerce: string
  consumer: string
}

export interface DeliveryFailure {
  detail: string
  reason: string
  referenceID: string
}

export interface DeliveryLocation {
  latitude: number
  longitude: number
}

export interface Delivery {
  failure?: DeliveryFailure
  handledAt: string
  location: DeliveryLocation
  status: string
}

export interface EvidencePhoto {
  takenAt: string
  type: string
  url: string
}

export interface DeliveryItem {
  deliveredQuantity: number
  description: string
  quantity: number
  sku: string
}

export interface Recipient {
  fullName: string
  nationalID: string
}

export interface DeliveryUnit {
  businessIdentifiers: BusinessIdentifiers
  delivery: Delivery
  evidencePhotos: EvidencePhoto[]
  items: DeliveryItem[]
  lpn: string
  orderReferenceID: string
  recipient: Recipient
}

export interface Driver {
  email: string
  nationalID: string
}

export interface ManualChange {
  performedBy: string
  reason: string
}

export interface RouteReference {
  id: number
  referenceID: string
  sequenceNumber: number
  startedAt: string
}

export interface Vehicle {
  plate: string
}


