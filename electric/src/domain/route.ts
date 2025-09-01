export interface Coordinates {
  latitude: number
  longitude: number
}

export interface Contact {
  email: string
  fullName: string
  nationalID: string
  phone: string
}

export interface PoliticalArea {
  adminAreaLevel1: string
  adminAreaLevel2: string
  adminAreaLevel3: string
  adminAreaLevel4: string
  code: string
}

export interface AddressInfo {
  addressLine1: string
  addressLine2: string
  contact: Contact
  coordinates: Coordinates
  politicalArea: PoliticalArea
  zipCode: string
}

export interface NodeInfo {
  referenceID: string
}

export interface Location {
  addressInfo: AddressInfo
  nodeInfo: NodeInfo
}

export interface VehicleCapacity {
  deliveryUnitsQuantity: number
  insurance: number
  volume: number
  weight: number
}

export interface TimeWindow {
  start: string
  end: string
}

export interface Vehicle {
  capacity: VehicleCapacity
  endLocation: Location
  plate: string
  skills: string[]
  startLocation: Location
  timeWindow: TimeWindow
}

export interface Geometry {
  encoding: string
  type: string
  value: string
}

export interface Evidence {
  downloadUrl: string
  uploadUrl: string
}

export interface DeliveryItem {
  description: string
  quantity: number
  sku: string
}

export interface DeliveryUnit {
  documentID: string
  evidences: Evidence[]
  items: DeliveryItem[]
  lpn: string
  price: number
  skills: string[]
  volume: number
  weight: number
}

export interface Order {
  deliveryUnits: DeliveryUnit[]
  documentID: string
  referenceID: string
}

export interface Visit {
  addressInfo: AddressInfo
  deliveryInstructions: string
  instructions: string
  nodeInfo: NodeInfo
  orders: Order[]
  sequenceNumber: number
  serviceTime: number
  timeWindow: TimeWindow
  type: string
  unassignedReason: string
}

export interface Route {
  id: number
  createdAt: string
  documentID: string
  geometry: Geometry
  planReferenceID: string
  referenceID: string
  vehicle: Vehicle
  visits: Visit[]
}