/**
 * Wrapper para el estado del driver que maneja automáticamente la cola de sincronización
 * cuando las operaciones se realizan offline
 */

import { syncQueue } from '../utils/sync-queue';
import type { DeliveryEvidence, NonDeliveryEvidence } from './driver-gun-state';
import type { RouteData } from '../utils/photo-upload';

// Re-export types
export type { DeliveryEvidence, NonDeliveryEvidence };

/**
 * Set route started status with offline sync support
 */
export async function setRouteStarted(routeId: string, started: boolean) {
  // Always update local state first (using existing local state system)
  const { setRouteStarted: setLocalRouteStarted } = await import('./driver-local-state');
  setLocalRouteStarted(routeId, started);

  // Queue for sync if offline, or sync immediately if online
  if (navigator.onLine) {
    // Try to sync immediately
    try {
      const { setRouteStarted: setGunRouteStarted } = await import('./driver-gun-state');
      setGunRouteStarted(routeId, started);
    } catch (error) {
      console.warn('Immediate sync failed, adding to queue:', error);
      queueRouteStartedSync(routeId, started);
    }
  } else {
    queueRouteStartedSync(routeId, started);
  }
}

function queueRouteStartedSync(routeId: string, started: boolean) {
  syncQueue.addToQueue({
    type: started ? 'route_start' : 'route_stop',
    routeId,
    data: { started },
    maxAttempts: 5,
    priority: 1 // High priority for route operations
  });
}

/**
 * Set route license with offline sync support
 */
export async function setRouteLicense(routeId: string, license: string) {
  // Note: No local state for license yet, but could be added
  
  if (navigator.onLine) {
    try {
      const { setRouteLicense: setGunRouteLicense } = await import('./driver-gun-state');
      setGunRouteLicense(routeId, license);
    } catch (error) {
      console.warn('Immediate license sync failed, adding to queue:', error);
      queueLicenseSync(routeId, license);
    }
  } else {
    queueLicenseSync(routeId, license);
  }
}

function queueLicenseSync(routeId: string, license: string) {
  syncQueue.addToQueue({
    type: 'license_set',
    routeId,
    data: { license },
    maxAttempts: 3,
    priority: 2
  });
}

/**
 * Set delivery status with offline sync support
 */
export async function setDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  status: 'delivered' | 'not-delivered'
) {
  // Always update local state first
  const { setDeliveryStatus: setLocalDeliveryStatus } = await import('./driver-local-state');
  setLocalDeliveryStatus(routeId, visitIndex, orderIndex, unitIndex, status);

  // Queue for sync if offline, or sync immediately if online
  if (navigator.onLine) {
    try {
      const { setDeliveryStatus: setGunDeliveryStatus } = await import('./driver-gun-state');
      setGunDeliveryStatus(routeId, visitIndex, orderIndex, unitIndex, status);
    } catch (error) {
      console.warn('Immediate delivery status sync failed, adding to queue:', error);
      queueDeliveryStatusSync(routeId, visitIndex, orderIndex, unitIndex, status);
    }
  } else {
    queueDeliveryStatusSync(routeId, visitIndex, orderIndex, unitIndex, status);
  }
}

function queueDeliveryStatusSync(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  status: 'delivered' | 'not-delivered'
) {
  syncQueue.addToQueue({
    type: 'delivery_status',
    routeId,
    data: { visitIndex, orderIndex, unitIndex, status },
    maxAttempts: 5,
    priority: 3
  });
}

/**
 * Set delivery evidence with offline sync support
 */
export async function setDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: DeliveryEvidence,
  routeData?: RouteData
) {
  // Always update local state first
  const { setDeliveryEvidence: setLocalDeliveryEvidence } = await import('./driver-local-state');
  setLocalDeliveryEvidence(routeId, visitIndex, orderIndex, unitIndex, evidence);

  // Queue for sync if offline, or sync immediately if online
  if (navigator.onLine) {
    try {
      const { setDeliveryEvidence: setGunDeliveryEvidence } = await import('./driver-gun-state');
      await setGunDeliveryEvidence(routeId, visitIndex, orderIndex, unitIndex, evidence, routeData);
    } catch (error) {
      console.warn('Immediate delivery evidence sync failed, adding to queue:', error);
      queueDeliveryEvidenceSync(routeId, visitIndex, orderIndex, unitIndex, evidence, routeData);
    }
  } else {
    queueDeliveryEvidenceSync(routeId, visitIndex, orderIndex, unitIndex, evidence, routeData);
  }
}

function queueDeliveryEvidenceSync(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: DeliveryEvidence,
  routeData?: RouteData
) {
  syncQueue.addToQueue({
    type: 'delivery_evidence',
    routeId,
    data: { visitIndex, orderIndex, unitIndex, evidence, routeData },
    maxAttempts: 5,
    priority: 4 // Lower priority than status changes
  });
}

/**
 * Set non-delivery evidence with offline sync support
 */
export async function setNonDeliveryEvidence(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: NonDeliveryEvidence,
  routeData?: RouteData
) {
  // Always update local state first
  const { setNonDeliveryEvidence: setLocalNonDeliveryEvidence } = await import('./driver-local-state');
  setLocalNonDeliveryEvidence(routeId, visitIndex, orderIndex, unitIndex, evidence);

  // Queue for sync if offline, or sync immediately if online
  if (navigator.onLine) {
    try {
      const { setNonDeliveryEvidence: setGunNonDeliveryEvidence } = await import('./driver-gun-state');
      await setGunNonDeliveryEvidence(routeId, visitIndex, orderIndex, unitIndex, evidence, routeData);
    } catch (error) {
      console.warn('Immediate non-delivery evidence sync failed, adding to queue:', error);
      queueNonDeliveryEvidenceSync(routeId, visitIndex, orderIndex, unitIndex, evidence, routeData);
    }
  } else {
    queueNonDeliveryEvidenceSync(routeId, visitIndex, orderIndex, unitIndex, evidence, routeData);
  }
}

function queueNonDeliveryEvidenceSync(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  evidence: NonDeliveryEvidence,
  routeData?: RouteData
) {
  syncQueue.addToQueue({
    type: 'non_delivery_evidence',
    routeId,
    data: { visitIndex, orderIndex, unitIndex, evidence, routeData },
    maxAttempts: 5,
    priority: 4 // Lower priority than status changes
  });
}

// Re-export read functions from local state (since we read locally first)
export {
  getDeliveryStatus,
  getDeliveryEvidence,
  getNonDeliveryEvidence,
  deliveryKey,
  evidenceKey,
  ndEvidenceKey,
  routeStartedKey
} from './driver-local-state';

// Re-export hooks and utility functions from gun state for reactive sync status
export {
  useDriverState,
  useGunData,
  useRouteStartedSync,
  getDeviceInfo,
  getAllRoutesSyncInfo,
  getDeliveryStatusFromState,
  getRouteLicenseFromState
} from './driver-gun-state';