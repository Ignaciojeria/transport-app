import type { RouteStart } from '../../../domain/route-start'
import type { GunRouteStart } from '../models/route-start-models'

// Mappers para convertir entre dominio y modelos internos de Gun.js

export function mapRouteStartToGun(
  routeStart: RouteStart
): GunRouteStart {
  return {
    // Datos del carrier (planos)
    carrier_name: routeStart.carrier.name,
    carrier_nationalID: routeStart.carrier.nationalID,
    
    // Datos del driver (planos)
    driver_email: routeStart.driver.email,
    driver_nationalID: routeStart.driver.nationalID,
    
    // Datos de la ruta (planos)
    route_id: routeStart.route.id,
    route_referenceID: routeStart.route.referenceID,
    
    // Datos del vehículo (planos)
    vehicle_plate: routeStart.vehicle.plate,
    
    // Timestamp
    startedAt: routeStart.startedAt
  }
}

export function mapGunToRouteStart(
  gunData: GunRouteStart
): RouteStart {
  return {
    carrier: {
      name: gunData.carrier_name,
      nationalID: gunData.carrier_nationalID
    },
    driver: {
      email: gunData.driver_email,
      nationalID: gunData.driver_nationalID
    },
    route: {
      id: gunData.route_id,
      referenceID: gunData.route_referenceID
    },
    vehicle: {
      plate: gunData.vehicle_plate
    },
    startedAt: gunData.startedAt
  }
}
