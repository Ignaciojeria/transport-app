// Modelos internos específicos para Gun.js para RouteStart
// Estos modelos evitan referencias circulares y objetos anidados

export interface GunRouteStart {
  // Datos del carrier (planos)
  carrier_name: string
  carrier_nationalID: string
  
  // Datos del driver (planos)
  driver_email: string
  driver_nationalID: string
  
  // Datos de la ruta (planos)
  route_id: number
  route_referenceID: string
  
  // Datos del vehículo (planos)
  vehicle_plate: string
  
  // Timestamp
  startedAt: string
}
