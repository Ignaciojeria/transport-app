// Funciones utilitarias para la vista del mapa

// Helper para obtener [lat, lng] desde addressInfo
export const getLatLngFromAddressInfo = (addr: any): [number, number] | null => {
  const c = addr?.coordinates
  if (!c) return null
  
  if (Array.isArray(c?.point) && c.point.length >= 2) {
    return [c.point[1] as number, c.point[0] as number]
  }
  
  if (typeof c.latitude === 'number' && typeof c.longitude === 'number') {
    return [c.latitude as number, c.longitude as number]
  }
  
  return null
}

// Decodificador de polylines (Google Encoded Polyline Algorithm Format)
export const decodePolyline = (encoded: string): Array<[number, number]> => {
  let index = 0
  const len = encoded.length
  let lat = 0
  let lng = 0
  const coordinates: Array<[number, number]> = []
  
  while (index < len) {
    let b = 0
    let shift = 0
    let result = 0
    
    do {
      b = encoded.charCodeAt(index++) - 63
      result |= (b & 0x1f) << shift
      shift += 5
    } while (b >= 0x20)
    
    const dlat = (result & 1) ? ~(result >> 1) : (result >> 1)
    lat += dlat

    shift = 0
    result = 0
    
    do {
      b = encoded.charCodeAt(index++) - 63
      result |= (b & 0x1f) << shift
      shift += 5
    } while (b >= 0x20)
    
    const dlng = (result & 1) ? ~(result >> 1) : (result >> 1)
    lng += dlng

    coordinates.push([lat * 1e-5, lng * 1e-5])
  }
  
  return coordinates
}

// Helper para obtener gradiente complementario
export const getGradientColor = (baseColor: string): string => {
  const colorMap: Record<string, string> = {
    '#10B981': '#059669', // Verde claro -> Verde más oscuro
    '#EF4444': '#DC2626', // Rojo (no entregado) -> Rojo más oscuro
    '#1D4ED8': '#1E40AF', // Azul oscuro (parcial) -> Azul más oscuro
    '#6B7280': '#4B5563', // Gris -> Gris más oscuro
  }
  return colorMap[baseColor] || '#7C3AED'
}

// Función para obtener el estado de una visita completa
export const getVisitStatus = (
  visit: any,
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined,
  visitIndex: number
): 'completed' | 'not-delivered' | 'partial' | 'pending' => {
  if (!visit) return 'pending'
  
  const allUnits: Array<{ status: 'delivered' | 'not-delivered' | undefined }> = []
  
  // Recopilar el estado de todas las unidades de entrega de la visita
  ;(visit.orders || []).forEach((order: any, oIdx: number) => {
    ;(order.deliveryUnits || []).forEach((_unit: any, uIdx: number) => {
      const status = getDeliveryUnitStatus(visitIndex, oIdx, uIdx)
      allUnits.push({ status })
    })
  })
  
  if (allUnits.length === 0) return 'pending'
  
  const deliveredCount = allUnits.filter(u => u.status === 'delivered').length
  const notDeliveredCount = allUnits.filter(u => u.status === 'not-delivered').length
  const totalCount = allUnits.length
  const processedCount = deliveredCount + notDeliveredCount
  
  if (processedCount === 0) return 'pending'
  
  // Si todas las unidades están marcadas como no entregadas
  if (notDeliveredCount === totalCount) return 'not-delivered'
  
  // Si todas las unidades están procesadas (entregadas o no entregadas)
  if (processedCount === totalCount) {
    // Si hay al menos una entregada, considerarla completada exitosamente
    return deliveredCount > 0 ? 'completed' : 'not-delivered'
  }
  
  // Estado mixto: algunas procesadas, otras pendientes
  return 'partial'
}

// Obtener color del marcador según el estado de la visita
export const getVisitMarkerColor = (visitStatus: 'completed' | 'not-delivered' | 'partial' | 'pending'): string => {
  switch (visitStatus) {
    case 'completed':
      return '#10B981' // Verde
    case 'not-delivered':
      return '#EF4444' // Rojo
    case 'partial':
      return '#1D4ED8' // Azul oscuro
    case 'pending':
    default:
      return '#6B7280' // Gris
  }
}
