export interface Carrier {
  name: string
  nationalID: string
}

export interface Driver {
  email: string
  nationalID: string
}

export interface RouteReference {
  id: number
  documentID: string
  referenceID: string
}

export interface VehicleStart {
  plate: string
}

export interface RouteStart {
  carrier: Carrier
  driver: Driver
  route: RouteReference
  startedAt: string
  vehicle: VehicleStart
}
