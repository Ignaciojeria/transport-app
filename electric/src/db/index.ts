// Exportar funciones del driver state
export {
  useDriverState,
  useRouteStartedSync,
  routeStartedKey,
  setRouteStarted,
  setDeliveryStatus,
  getDeliveryStatusFromState,
  setDeliveryEvidence,
  setNonDeliveryEvidence,
  setRouteLicense,
  getRouteLicenseFromState
} from './driver-gun-state'

// Exportar funciones del route start state
export {
  setRouteStart,
  getRouteStart,
  isRouteStarted,
  getVehiclePlate,
  getDriverInfo,
  getCarrierInfo,
  clearRouteStart
} from './route-start-gun-state'

// Exportar hooks
export { useRouteStartSync } from '../hooks/useRouteStartSync'

// Exportar funciones de creación de colecciones
export { createRoutesCollection } from './create-routes-collection'
