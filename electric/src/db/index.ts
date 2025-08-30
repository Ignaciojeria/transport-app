// Exportar funciones del deliveries state
export {
  useDeliveriesState,
  setDeliveryStatus,
  getDeliveryStatusFromState,
  getNonDeliveryEvidenceFromState,
  setDeliveryEvidence,
  setSuccessfulDelivery,
  setFailedDelivery,
  setDeliveryUnitByEntity
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
  getRouteLicenseFromState,
  useRouteStartDomain
} from './gun/state/route-start-gun-state'

// Exportar funciones de creación de colecciones
export { createRoutesCollection } from './create-routes-collection'
