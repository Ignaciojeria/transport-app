// Exportar funciones del deliveries state
export {
  useDeliveriesState,
  setDeliveryStatus,
  getDeliveryStatusFromState,
  getNonDeliveryEvidenceFromState,
  setDeliveryEvidence,
  setSuccessfulDelivery,
  setFailedDelivery
} from './gun'

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
