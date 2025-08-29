// Exportar funciones del deliveries state
export {
  useDeliveriesState,
  setDeliveryStatus,
  getDeliveryStatusFromState,
  setDeliveryEvidence,
  setNonDeliveryEvidence
} from './deliveries-gun-state'

// Exportar funciones del route start state
export {
  setRouteStart,
  getRouteStart,
  isRouteStarted,
  getVehiclePlate,
  getDriverInfo,
  getCarrierInfo,
  clearRouteStart,
  useRouteStartedSync,
  routeStartedKey,
  setRouteStarted,
  setRouteLicense,
  getRouteLicenseFromState
} from './route-start-gun-state'

// Exportar hooks
export { useRouteStartSync } from '../hooks/useRouteStartSync'

// Exportar funciones de creación de colecciones
export { createRoutesCollection } from './create-routes-collection'
