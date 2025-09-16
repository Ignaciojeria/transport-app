// Utilidades para agrupar delivery units similares

export interface GroupableDeliveryUnit {
  unit: any
  uIdx: number
  status: 'delivered' | 'not-delivered' | undefined
  visitIndex: number
  orderIndex: number
  order: any
}

export interface DeliveryGroup {
  key: string // Clave única para el grupo (dirección + cliente)
  addressInfo: {
    addressLine1: string
    coordinates: any
  }
  // Información del primer cliente del grupo (para propósitos de display)
  primaryContact: {
    fullName: string
  }
  units: GroupableDeliveryUnit[]
  totalUnits: number
  pendingUnits: number
}

/**
 * Agrupa delivery units que pertenecen a la misma dirección
 * Esto permite agrupar diferentes clientes en la misma dirección
 */
export function groupDeliveryUnitsByLocation(
  visit: any,
  visitIndex: number,
  getDeliveryUnitStatus: (vIdx: number, oIdx: number, uIdx: number) => 'delivered' | 'not-delivered' | undefined
): DeliveryGroup[] {
  const groups: { [key: string]: DeliveryGroup } = {}
  
  // Obtener todas las delivery units de la visita
  const allUnits: GroupableDeliveryUnit[] = []
  
  visit.orders?.forEach((order: any, orderIndex: number) => {
    order.deliveryUnits?.forEach((unit: any, unitIndex: number) => {
      allUnits.push({
        unit,
        uIdx: unitIndex,
        status: getDeliveryUnitStatus(visitIndex, orderIndex, unitIndex),
        visitIndex,
        orderIndex,
        order
      })
    })
  })
  
  // Agrupar por dirección (sin importar el cliente)
  allUnits.forEach((deliveryUnit) => {
    const addressInfo = visit.addressInfo
    const address = addressInfo?.addressLine1 || 'Sin dirección'
    
    // Crear clave única para el grupo (solo dirección)
    const groupKey = address
    
    if (!groups[groupKey]) {
      groups[groupKey] = {
        key: groupKey,
        addressInfo: {
          addressLine1: addressInfo?.addressLine1,
          coordinates: addressInfo?.coordinates
        },
        primaryContact: {
          fullName: deliveryUnit.order?.contact?.fullName || 'Sin nombre'
        },
        units: [],
        totalUnits: 0,
        pendingUnits: 0
      }
    }
    
    groups[groupKey].units.push(deliveryUnit)
    groups[groupKey].totalUnits++
    
    if (!deliveryUnit.status) {
      groups[groupKey].pendingUnits++
    }
  })
  
  // Convertir a array y filtrar solo grupos con múltiples unidades
  return Object.values(groups).filter(group => group.totalUnits > 1)
}

/**
 * Verifica si un grupo tiene unidades pendientes
 */
export function hasGroupPendingUnits(group: DeliveryGroup): boolean {
  return group.pendingUnits > 0
}

/**
 * Obtiene las unidades pendientes de un grupo
 */
export function getGroupPendingUnits(group: DeliveryGroup): GroupableDeliveryUnit[] {
  return group.units.filter(unit => !unit.status)
}

/**
 * Obtiene las unidades entregadas de un grupo
 */
export function getGroupDeliveredUnits(group: DeliveryGroup): GroupableDeliveryUnit[] {
  return group.units.filter(unit => unit.status === 'delivered')
}

/**
 * Obtiene las unidades no entregadas de un grupo
 */
export function getGroupNotDeliveredUnits(group: DeliveryGroup): GroupableDeliveryUnit[] {
  return group.units.filter(unit => unit.status === 'not-delivered')
}

/**
 * Agrupa delivery units por dirección para la tarjeta de "Siguiente visita"
 * Esto permite mostrar múltiples clientes en la misma dirección
 */
export function groupDeliveryUnitsByAddressForNextVisit(
  visits: any[],
  getDeliveryUnitStatus: (vIdx: number, oIdx: number, uIdx: number) => 'delivered' | 'not-delivered' | undefined
): { [address: string]: { clients: string[], totalUnits: number, pendingUnits: number, visitIndex: number } } {
  const addressGroups: { [address: string]: { clients: string[], totalUnits: number, pendingUnits: number, visitIndex: number } } = {}
  
  visits.forEach((visit: any, visitIndex: number) => {
    const address = visit.addressInfo?.addressLine1 || 'Sin dirección'
    const clientName = visit.orders?.[0]?.contact?.fullName || 'Sin nombre'
    
    // Contar unidades totales y pendientes para esta visita
    let totalUnits = 0
    let pendingUnits = 0
    
    visit.orders?.forEach((order: any, orderIndex: number) => {
      order.deliveryUnits?.forEach((_unit: any, unitIndex: number) => {
        totalUnits++
        const status = getDeliveryUnitStatus(visitIndex, orderIndex, unitIndex)
        if (!status) {
          pendingUnits++
        }
      })
    })
    
    // Solo incluir si hay unidades pendientes
    if (pendingUnits > 0) {
      if (!addressGroups[address]) {
        addressGroups[address] = {
          clients: [],
          totalUnits: 0,
          pendingUnits: 0,
          visitIndex: visitIndex
        }
      }
      
      addressGroups[address].clients.push(clientName)
      addressGroups[address].totalUnits += totalUnits
      addressGroups[address].pendingUnits += pendingUnits
    }
  })
  
  return addressGroups
}
