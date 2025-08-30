import type { Visit, Order, DeliveryUnit } from '../domain/route'

export interface ReportData {
  routeId: string
  routeDbId?: number
  routeLicense?: string
  visits: Visit[]
  localState: any
}

export interface DeliveryUnitData {
  visit: Visit
  visitIndex: number
  order: Order
  orderIndex: number
  unit: DeliveryUnit
  unitIndex: number
  status: 'delivered' | 'not-delivered' | undefined
}

export function generateReportData(
  visits: Visit[],
  localState: any,
  routeId: string
): DeliveryUnitData[] {
  const allUnits: DeliveryUnitData[] = []
  
  visits.forEach((visit: Visit, visitIndex: number) => {
    visit?.orders?.forEach((order: Order, orderIndex: number) => {
      order?.deliveryUnits?.forEach((unit: DeliveryUnit, unitIndex: number) => {
        const status = getDeliveryUnitStatus(localState, routeId, visitIndex, orderIndex, unitIndex)
        
        // Debug: verificar evidencia para no entregas
        if (status === 'not-delivered') {
          const evidence = getNonDeliveryEvidence(localState, routeId, visitIndex, orderIndex, unitIndex)
          // console.log(`🔍 Evidencia no entrega (${visitIndex},${orderIndex},${unitIndex}):`, evidence) // Comentado para reducir logs
        }
        
        allUnits.push({
          visit,
          visitIndex,
          order,
          orderIndex,
          unit,
          unitIndex,
          status
        })
      })
    })
  })
  
  return allUnits
}

export function getDeliveryUnitStatus(
  localState: any,
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): 'delivered' | 'not-delivered' | undefined {
  const key = `delivery:${routeId}:${visitIndex}-${orderIndex}-${unitIndex}`
  const data = localState?.[key]
  
  if (typeof data === 'string') {
    // Estado simple (formato anterior)
    return data as 'delivered' | 'not-delivered'
  } else if (data && typeof data === 'object' && data.status) {
    // Estado con evidencia (nuevo formato)
    return data.status as 'delivered' | 'not-delivered'
  }
  
  return undefined
}

export function getDeliveryEvidence(
  localState: any,
  routeId: string,
  visitIndex: number,
  unitIndex: number,
  orderIndex: number
): any {
  // Buscar en la clave evidence para entregas exitosas (nueva estructura del dominio)
  const evidenceKey = `evidence:${routeId}:${visitIndex}-${orderIndex}-${unitIndex}`
  const evidence = localState?.[evidenceKey]
  
  // Debug: ver qué datos estamos recibiendo
  console.log(`🔍 getDeliveryEvidence para ${evidenceKey}:`, evidence)
  
  if (evidence && typeof evidence === 'object') {
    // Ahora evidence ya es una entidad del dominio (convertida por useDeliveriesState)
    const recipientName = evidence.recipient?.fullName || ''
    const recipientRut = evidence.recipient?.nationalID || ''
    const takenAt = evidence.evidencePhotos?.[0]?.takenAt || evidence.delivery?.handledAt || Date.now()
    const photoDataUrl = evidence.evidencePhotos?.[0]?.url || ''
    
    console.log(`🔍 Datos extraídos:`, { recipientName, recipientRut, takenAt, photoDataUrl })
    
    const result = {
      recipientName,
      recipientRut,
      takenAt,
      photoDataUrl
    }
    
    console.log(`🔍 Resultado final:`, result)
    return result
  }
  
  // Fallback: buscar en delivery si no se encuentra en evidence
  const deliveryKey = `delivery:${routeId}:${visitIndex}-${orderIndex}-${unitIndex}`
  const deliveryData = localState?.[deliveryKey]
  
  if (deliveryData && typeof deliveryData === 'string' && deliveryData === 'delivered') {
    // Si solo hay estado 'delivered' sin evidencia detallada
    return {
      recipientName: '',
      recipientRut: '',
      takenAt: Date.now()
    }
  }
  
  return null
}

export function getNonDeliveryEvidence(
  localState: any,
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): any {
  // Buscar en la clave evidence para no entregas (nueva estructura del dominio)
  const evidenceKey = `evidence:${routeId}:${visitIndex}-${orderIndex}-${unitIndex}`
  const evidence = localState?.[evidenceKey]
  
  if (evidence && typeof evidence === 'object') {
    // Nueva estructura: evidence.delivery.failure y evidence.evidencePhotos
    const failure = evidence.delivery?.failure
    const evidencePhotos = evidence.evidencePhotos
    
    // Convertir evidencePhotos de formato Gun.js a array si es necesario
    let photosArray = []
    if (evidencePhotos) {
      if (Array.isArray(evidencePhotos)) {
        photosArray = evidencePhotos
      } else if (typeof evidencePhotos === 'object') {
        // Formato Gun.js: {"0": photo, "1": photo2, ...}
        photosArray = Object.values(evidencePhotos).filter(photo => photo !== null && photo !== undefined)
      }
    }
    
    // Extraer la primera foto como evidencia principal
    const mainPhoto = photosArray[0]
    
    return {
      reason: failure?.reason || '',
      observations: failure?.detail || '',
      takenAt: mainPhoto?.takenAt || evidence.delivery?.handledAt || Date.now(),
      photoDataUrl: mainPhoto?.url || ''
    }
  }
  
  // Fallback: buscar en delivery si no se encuentra en evidence
  const deliveryKey = `delivery:${routeId}:${visitIndex}-${orderIndex}-${unitIndex}`
  const deliveryData = localState?.[deliveryKey]
  
  if (deliveryData && typeof deliveryData === 'object' && deliveryData.status === 'not-delivered') {
    // Extraer información del dominio DeliveryFailure
    return {
      reason: deliveryData.failure?.reason || '',
      observations: deliveryData.failure?.detail || '',
      takenAt: deliveryData.timestamp || Date.now(),
      photoDataUrl: deliveryData.photoDataUrl || ''
    }
  }
  
  // Fallback: buscar en la clave antigua nd-evidence
  const ndEvidenceKey = `nd-evidence:${routeId}:${visitIndex}-${orderIndex}-${unitIndex}`
  const evidenceOld = localState?.[ndEvidenceKey]
  if (evidenceOld) {
    try {
      const parsed = typeof evidenceOld === 'string' ? JSON.parse(evidenceOld) : evidenceOld
      return parsed
    } catch {
      return null
    }
  }
  
  return null
}

// Función de debug para diagnosticar problemas de datos
export function debugDeliveryData(localState: any, routeId: string): void {
  console.log('🔍 Debug: Analizando estructura de datos para ruta:', routeId)
  
  // Buscar todas las claves relacionadas con la ruta
  const relevantKeys = Object.keys(localState).filter(key => key.includes(routeId))
  
  relevantKeys.forEach(key => {
    const data = localState[key]
    console.log(`📋 Clave: ${key}`, {
      tipo: typeof data,
      esObjeto: typeof data === 'object',
      propiedades: typeof data === 'object' ? Object.keys(data) : 'N/A',
      datos: data
    })
  })
}

export function generateCSVContent(units: DeliveryUnitData[], reportData: ReportData): string {
  const headers = [
    'ID_Ruta',
    'Patente_Vehiculo',
    'Secuencia_Visita',
    'Nombre_Cliente',
    'Direccion',
    'Telefono',
    'Email',
    'Referencia_Orden',
    'Unidad_Entrega',
    'Descripcion_Items',
    'Cantidad_Total',
    'Peso_Total',
    'Volumen_Total',
    'Estado_Entrega',
    'Nombre_Receptor',
    'Documento_Receptor',
    'Fecha_Gestion',
    'Motivo_No_Entrega',
    'Observaciones_No_Entrega',
    'Coordenadas_Lat',
    'Coordenadas_Lng'
  ]

  const rows = units.map(unit => {
    const { visit, visitIndex, order, orderIndex, unit: deliveryUnit, unitIndex, status } = unit
    const contactInfo = visit?.addressInfo?.contact || {}
    const coordinates = visit?.addressInfo?.coordinates
    const lat = coordinates?.latitude
    const lng = coordinates?.longitude
    
    const statusText = status === 'delivered' ? 'Entregado' : 
                      status === 'not-delivered' ? 'No Entregado' : 'Pendiente'
    
    // Calcular totales de la unidad
    const items = deliveryUnit?.items || []
    const totalQuantity = items.reduce((sum: number, item: any) => sum + (Number(item?.quantity) || 0), 0)
    const itemDescriptions = items.map((item: any) => item?.description || '').filter(Boolean).join('; ')
    
    // Obtener evidencia si existe
    let recipientName = ''
    let recipientDocument = ''
    let managementDate = ''
    let nonDeliveryReason = ''
    let nonDeliveryObservations = ''
    
    if (status === 'delivered') {
      const evidence = getDeliveryEvidence(reportData.localState, reportData.routeId, visitIndex, orderIndex, unitIndex)
      if (evidence) {
        recipientName = evidence?.recipientName || ''
        recipientDocument = evidence?.recipientRut || ''
        managementDate = evidence?.takenAt ? new Date(evidence.takenAt).toLocaleString('es-CL') : ''
      }
    } else if (status === 'not-delivered') {
      const evidence = getNonDeliveryEvidence(reportData.localState, reportData.routeId, visitIndex, orderIndex, unitIndex)
      if (evidence) {
        nonDeliveryReason = evidence?.reason || ''
        nonDeliveryObservations = evidence?.observations || ''
        managementDate = evidence?.takenAt ? new Date(evidence.takenAt).toLocaleString('es-CL') : ''
      }
    }
    
    return [
      (reportData.routeDbId || reportData.routeId)?.toString() || '',
      reportData.routeLicense || '',
      visit?.sequenceNumber?.toString() || (visitIndex + 1).toString(),
      contactInfo?.fullName || '',
      visit?.addressInfo?.addressLine1 || '',
      contactInfo?.phone || '',
      contactInfo?.email || '',
      order?.referenceID || '',
      `Unidad ${unitIndex + 1}`,
      itemDescriptions,
      totalQuantity.toString(),
      (deliveryUnit?.weight || '').toString(),
      (deliveryUnit?.volume || '').toString(),
      statusText,
      recipientName,
      recipientDocument,
      managementDate,
      nonDeliveryReason,
      nonDeliveryObservations,
      lat?.toString() || '',
      lng?.toString() || ''
    ]
  })

  return [
    headers.join(','),
    ...rows.map(row => row.map(field => `"${(field || '').toString().replace(/"/g, '""')}"`).join(','))
  ].join('\n')
}

export function generateExcelContent(units: DeliveryUnitData[], reportData: ReportData): string {
  const headers = [
    'ID_Ruta',
    'Patente_Vehiculo',
    'Secuencia_Visita',
    'Nombre_Cliente',
    'Direccion',
    'Telefono',
    'Email',
    'Referencia_Orden',
    'Unidad_Entrega',
    'Descripcion_Items',
    'Cantidad_Total',
    'Peso_Total',
    'Volumen_Total',
    'Estado_Entrega',
    'Nombre_Receptor',
    'Documento_Receptor',
    'Fecha_Gestion',
    'Motivo_No_Entrega',
    'Observaciones_No_Entrega',
    'Coordenadas_Lat',
    'Coordenadas_Lng'
  ]

  const rows = units.map(unit => {
    const { visit, visitIndex, order, orderIndex, unit: deliveryUnit, unitIndex, status } = unit
    const contactInfo = visit?.addressInfo?.contact || {}
    const coordinates = visit?.addressInfo?.coordinates
    const lat = coordinates?.latitude
    const lng = coordinates?.longitude
    
    const statusText = status === 'delivered' ? 'Entregado' : 
                      status === 'not-delivered' ? 'No Entregado' : 'Pendiente'
    
    // Calcular totales de la unidad
    const items = deliveryUnit?.items || []
    const totalQuantity = items.reduce((sum: number, item: any) => sum + (Number(item?.quantity) || 0), 0)
    const itemDescriptions = items.map((item: any) => item?.description || '').filter(Boolean).join('; ')
    
    // Obtener evidencia si existe
    let recipientName = ''
    let recipientDocument = ''
    let managementDate = ''
    let nonDeliveryReason = ''
    let nonDeliveryObservations = ''
    
    if (status === 'delivered') {
      const evidence = getDeliveryEvidence(reportData.localState, reportData.routeId, visitIndex, orderIndex, unitIndex)
      if (evidence) {
        recipientName = evidence?.recipientName || ''
        recipientDocument = evidence?.recipientRut || ''
        managementDate = evidence?.takenAt ? new Date(evidence.takenAt).toLocaleString('es-CL') : ''
      }
    } else if (status === 'not-delivered') {
      const evidence = getNonDeliveryEvidence(reportData.localState, reportData.routeId, visitIndex, orderIndex, unitIndex)
      if (evidence) {
        nonDeliveryReason = evidence?.reason || ''
        nonDeliveryObservations = evidence?.observations || ''
        managementDate = evidence?.takenAt ? new Date(evidence.takenAt).toLocaleString('es-CL') : ''
      }
    }
    
    return [
      (reportData.routeDbId || reportData.routeId)?.toString() || '',
      reportData.routeLicense || '',
      visit?.sequenceNumber?.toString() || (visitIndex + 1).toString(),
      contactInfo?.fullName || '',
      visit?.addressInfo?.addressLine1 || '',
      contactInfo?.phone || '',
      contactInfo?.email || '',
      order?.referenceID || '',
      `Unidad ${unitIndex + 1}`,
      itemDescriptions,
      totalQuantity.toString(),
      (deliveryUnit?.weight || '').toString(),
      (deliveryUnit?.volume || '').toString(),
      statusText,
      recipientName,
      recipientDocument,
      managementDate,
      nonDeliveryReason,
      nonDeliveryObservations,
      lat?.toString() || '',
      lng?.toString() || ''
    ]
  })

  return `
    <html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
    <head>
      <meta charset="utf-8"/>
      <!--[if gte mso 9]><xml><x:ExcelWorkbook><x:ExcelWorksheets><x:ExcelWorksheet><x:Name>Reporte Ruta</x:Name><x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions></x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]-->
    </head>
    <body>
      <table border="1">
        <thead>
          <tr>
            ${headers.map(header => `<th style="background-color: #4F46E5; color: white; font-weight: bold; padding: 8px;">${header}</th>`).join('')}
          </tr>
        </thead>
        <tbody>
          ${rows.map(row => `
            <tr>
              ${row.map(cell => `<td style="padding: 4px; border: 1px solid #ccc;">${cell || ''}</td>`).join('')}
            </tr>
          `).join('')}
        </tbody>
      </table>
    </body>
    </html>
  `
}

export function downloadFile(content: string, filename: string, mimeType: string) {
  const blob = new Blob([content], { type: mimeType })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  link.style.visibility = 'hidden'
  
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  
  // Vibración táctil si está disponible
  try { (navigator as any)?.vibrate?.(100) } catch {}
}

