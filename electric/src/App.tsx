/* eslint-disable @typescript-eslint/no-explicit-any */
import { useParams } from '@tanstack/react-router'
import { useRoutes } from './db'
import { 
  useDeliveriesState, 
  useRouteStartedSync,
  routeStartedKey, 
  setRouteStarted as setRouteStartedLocal, 
  getDeliveryStatusFromState, 
  setDeliveryEvidence, 
  setRouteLicense,
  getRouteLicenseFromState,
  setRouteStart,
  setDeliveryUnitByEntity,
} from './db'
import { useState, useEffect } from 'react'
import { Play, Menu, Truck, Route, Map } from 'lucide-react'
import { Sidebar, DeliveryModal, NonDeliveryModal, VisitCard, NextVisitCard, DownloadReportModal, RouteStartModal, VisitTabs, MapView } from './components'
import { 
  generateReportData, 
  generateCSVContent, 
  generateExcelContent, 
  downloadFile,
  debugDeliveryData,
  type ReportData 
} from './components/DownloadReportModal.utils'
import type { Route as RouteType } from './domain/route'
import type { RouteStart } from './domain/route-start'
import type { DeliveryUnit, DeliveryEvent } from './domain/deliveries'


// Componente para rutas específicas del driver
export function RouteComponent() {
  // Obtener el routeId de los parámetros de la ruta usando TanStack Router
  const { routeId } = useParams({ from: '/driver/routes/$routeId' })
  
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || ''
  
  const routes = useRoutes(token, routeId)

  return (
    <div>
      {/* Renderizar UI si hay datos */}
      {routes.length > 0 ? (
        <DeliveryRouteView 
          routeId={routeId} 
          routeData={routes[0]} 
          routeDbId={routes[0].electricId} 
        />
      ) : (
        <div>Cargando rutas...</div>
      )}
    </div>
  )
}

function DeliveryRouteView({ routeId, routeData, routeDbId }: { routeId: string; routeData: RouteType; routeDbId?: string }) {
  const [activeTab, setActiveTab] = useState<'en-ruta' | 'entregados' | 'no-entregados'>('en-ruta')
  const [viewMode, setViewMode] = useState<'list' | 'map'>('list')
  // fullscreen deshabilitado para evitar cambios por clic en el mapa
  const [nextVisitIndex, setNextVisitIndex] = useState<number | null>(null)
  const [lastCenteredVisit, setLastCenteredVisit] = useState<number | null>(null) // Recordar última visita centrada


  // Modal de evidencia
  const [evidenceModal, setEvidenceModal] = useState<{ open: boolean; vIdx: number | null; oIdx: number | null; uIdx: number | null }>({ open: false, vIdx: null, oIdx: null, uIdx: null })
  const [submittingEvidence, setSubmittingEvidence] = useState(false)
  const [ndModal, setNdModal] = useState<{ open: boolean; vIdx: number | null; oIdx: number | null; uIdx: number | null }>({ open: false, vIdx: null, oIdx: null, uIdx: null })



  // Modal de descarga de reporte
  const [downloadModal, setDownloadModal] = useState(false)
  
  // Estado del sidebar
  const [sidebarOpen, setSidebarOpen] = useState(false)

  // Estado local reactivo via GunJS
  const { data: localState } = useDeliveriesState()
  
  // Debug: Log cuando cambia el estado local (comentado en producción)
  // useEffect(() => {
  //   console.log('🔄 localState cambió:', localState?.s ? Object.keys(localState.s).filter(k => k.includes('delivery:')) : 'no state')
  // }, [localState])
  // Información de sincronización entre dispositivos
  const syncInfo = useRouteStartedSync(routeId)

  // Función para sincronizar posición del marcador entre dispositivos
  const setMarkerPosition = async (routeId: string, visitIndex: number, coordinates: [number, number]) => {
    try {
      const { deliveriesData } = await import('./db/gun')
      const key = `marker_position:${routeId}`
      const deviceId = (() => {
        try {
          return syncInfo?.deviceId || `device_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
        } catch {
          return `device_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
        }
      })()
      const data = {
        visitIndex,
        coordinates,
        timestamp: Date.now(),
        deviceId
      }
      // Debug: logs removidos para limpiar la consola
              deliveriesData.get(key).put(JSON.stringify(data))
    } catch (error) {
      console.error('Error sincronizando posición de marcador:', error)
    }
  }

  // Hook para escuchar cambios de posición del marcador
  const [markerPosition, setMarkerPosition_] = useState<{
    visitIndex: number
    coordinates: [number, number]
    timestamp: number
    deviceId: string
  } | null>(null)

  useEffect(() => {
    import('./db/gun').then(({ deliveriesData }) => {
      const key = `marker_position:${routeId}`
      const unsubscribe = deliveriesData.get(key).on((data) => {
        if (data && typeof data === 'string') {
          try {
            const parsed = JSON.parse(data)
            // Debug: logs removidos para limpiar la consola
            setMarkerPosition_(parsed)
          } catch {}
        }
      })

      return () => {
        if (unsubscribe && typeof unsubscribe.off === 'function') {
          unsubscribe.off()
        }
      }
    })
  }, [routeId])

  const routeStarted = (localState?.s?.[`${routeStartedKey(routeId)}_simple`] === 'true')

  const [routeStartModal, setRouteStartModal] = useState(false)
  const [initialDeliveryEvent, setInitialDeliveryEvent] = useState<DeliveryEvent | undefined>(undefined)
  const [initialNonDeliveryEvent, setInitialNonDeliveryEvent] = useState<DeliveryEvent | undefined>(undefined)

  const handleStartRoute = () => {
    setRouteStartModal(true)
  }

  const handleLicenseConfirm = async (license: string) => {
    if (!license.trim()) {
      return
    }
    
    try {
      // Crear la entidad RouteStart
      const routeStart: RouteStart = {
        carrier: {
          name: '', // Por ahora vacío como mencionaste
          nationalID: ''
        },
        driver: {
          email: 'driver@example.com', // Por ahora hardcodeado
          nationalID: '12345678-9' // Por ahora hardcodeado
        },
        route: {
          id: parseInt(routeDbId || routeId) || 0,
          documentID: routeId,
          referenceID: routeId
        },
        startedAt: new Date().toISOString(),
        vehicle: {
          plate: license.trim()
        }
      }
      
      // Guardar en la nueva colección RouteStart
      await setRouteStart(routeId, routeStart)
      
      // También mantener compatibilidad con el sistema anterior
      setRouteLicense(routeId, license.trim())
      setRouteStartedLocal(routeId, true)
      
      setRouteStartModal(false)
      console.log('🚀 Ruta iniciada con RouteStart:', routeStart)
      
    } catch (error) {
      console.error('Error iniciando ruta:', error)
      alert('Error al iniciar la ruta. Por favor intenta nuevamente.')
    }
  }



  // Nota: setDeliveryStatus se usa directamente en cada flujo

  const getDeliveryUnitStatus = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    const status = getDeliveryStatusFromState(localState?.s || {}, routeId, visitIndex, orderIndex, unitIndex)
    
    // Debug: logs removidos para limpiar la consola
    
    return status
  }

  const openDeliveryFor = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    // Crear un DeliveryEvent inicial simplificado
    const initialDeliveryEvent: DeliveryEvent = {
      carrier: {
        name: '',
        nationalID: ''
      },
      deliveryUnits: [{
        businessIdentifiers: {
          commerce: '',
          consumer: ''
        },
        delivery: {
          status: 'pending',
          handledAt: new Date().toISOString(),
          location: { latitude: 0, longitude: 0 }
        },
        evidencePhotos: [],
        items: [],
        lpn: '',
        orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
        recipient: {
          fullName: '',
          nationalID: ''
        }
      }],
      driver: {
        email: '',
        nationalID: ''
      },
      route: {
        id: routeData?.id || 0,
        documentID: routeData?.documentID || '',
        referenceID: routeData?.referenceID || '',
        sequenceNumber: 0,
        startedAt: new Date().toISOString()
      },
      vehicle: {
        plate: routeData?.vehicle?.plate || ''
      }
    }
    
    setEvidenceModal({ open: true, vIdx: visitIndex, oIdx: orderIndex, uIdx: unitIndex })
    setInitialDeliveryEvent(initialDeliveryEvent)
  }

  const openNonDeliveryFor = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    // Crear un DeliveryEvent inicial para no entrega
    const initialNonDeliveryEvent: DeliveryEvent = {
      carrier: {
        name: '',
        nationalID: ''
      },
      deliveryUnits: [{
        businessIdentifiers: {
          commerce: '',
          consumer: ''
        },
        delivery: {
          status: 'pending',
          handledAt: new Date().toISOString(),
          location: { latitude: 0, longitude: 0 }
        },
        evidencePhotos: [],
        items: [],
        lpn: '',
        orderReferenceID: `${routeId}-${visitIndex}-${orderIndex}-${unitIndex}`,
        recipient: {
          fullName: '',
          nationalID: ''
        }
      }],
      driver: {
        email: '',
        nationalID: ''
      },
      route: {
        id: routeData?.id || 0,
        documentID: routeData?.documentID || '',
        referenceID: routeData?.referenceID || '',
        sequenceNumber: 0,
        startedAt: new Date().toISOString()
      },
      vehicle: {
        plate: routeData?.vehicle?.plate || ''
      }
    }
    
    setNdModal({ open: true, vIdx: visitIndex, oIdx: orderIndex, uIdx: unitIndex })
    setInitialNonDeliveryEvent(initialNonDeliveryEvent)
  }

  const closeNdModal = () => {
    setNdModal({ open: false, vIdx: null, oIdx: null, uIdx: null })
    setInitialNonDeliveryEvent(undefined)
  }

  const closeEvidenceModal = () => {
    setEvidenceModal({ open: false, vIdx: null, oIdx: null, uIdx: null })
    setSubmittingEvidence(false)
    setInitialDeliveryEvent(undefined)
  }

  



  const submitEvidence = async (deliveryEvent: DeliveryEvent) => {
    if (!evidenceModal.open || evidenceModal.vIdx === null || evidenceModal.oIdx === null || evidenceModal.uIdx === null) return
    try {
      setSubmittingEvidence(true)
      
      console.log('💾 Guardando evidencia de entrega para:', { routeId, vIdx: evidenceModal.vIdx, oIdx: evidenceModal.oIdx, uIdx: evidenceModal.uIdx })
      
      // Extraer datos del DeliveryEvent
      const recipientName = deliveryEvent.deliveryUnits[0]?.recipient?.fullName || ''
      const recipientRut = deliveryEvent.deliveryUnits[0]?.recipient?.nationalID || ''
      const photoDataUrl = deliveryEvent.deliveryUnits[0]?.evidencePhotos[0]?.url || ''
      
      // Crear entidad del dominio
      const deliveryUnit: Partial<DeliveryUnit> = {
        delivery: {
          status: 'delivered',
          handledAt: new Date().toISOString(),
          location: { latitude: 0, longitude: 0 }
        },
        recipient: {
          fullName: recipientName,
          nationalID: recipientRut
        },
        evidencePhotos: [{
          takenAt: new Date().toISOString(),
          type: 'delivery',
          url: photoDataUrl,
        }],
        orderReferenceID: `${routeId}-${evidenceModal.vIdx}-${evidenceModal.oIdx}-${evidenceModal.uIdx}`,
      }
      
      // Pasar entidad del dominio a setDeliveryEvidence
      setDeliveryEvidence(routeId, evidenceModal.vIdx, evidenceModal.oIdx, evidenceModal.uIdx, deliveryUnit)
      console.log('📦 Estado de entrega establecido a "delivered" por setDeliveryEvidence')
      closeEvidenceModal()
    } finally {
      setSubmittingEvidence(false)
    }
  }

  const submitNonDelivery = async (deliveryEvent: DeliveryEvent) => {
    if (!ndModal.open || ndModal.vIdx === null || ndModal.oIdx === null || ndModal.uIdx === null) return
    try {
      setSubmittingEvidence(true)
      
      console.log('💾 Guardando evidencia de no entrega para:', { routeId, vIdx: ndModal.vIdx, oIdx: ndModal.oIdx, uIdx: ndModal.uIdx })
      
      // Extraer datos del DeliveryEvent
      const failure = deliveryEvent.deliveryUnits[0]?.delivery?.failure
      const photoDataUrl = deliveryEvent.deliveryUnits[0]?.evidencePhotos[0]?.url
      
      if (!failure || !photoDataUrl) {
        console.error('❌ Datos incompletos en DeliveryEvent')
        return
      }
      
      // Crear la entidad del dominio completa
      const deliveryUnit: Partial<DeliveryUnit> & {
        routeId: string
        visitIndex: number
        orderIndex: number
        unitIndex: number
      } = {
        routeId,
        visitIndex: ndModal.vIdx,
        orderIndex: ndModal.oIdx,
        unitIndex: ndModal.uIdx,
        delivery: {
          status: 'not-delivered',
          handledAt: new Date().toISOString(),
          location: { latitude: 0, longitude: 0 }, // TODO: obtener ubicación real
          failure: {
            reason: failure.reason,
            detail: failure.detail,
            referenceID: `${routeId}-${ndModal.vIdx}-${ndModal.oIdx}-${ndModal.uIdx}`
          }
        },
        evidencePhotos: [{
          takenAt: new Date().toISOString(),
          type: 'non-delivery',
          url: photoDataUrl,
        }],
        orderReferenceID: `${routeId}-${ndModal.vIdx}-${ndModal.oIdx}-${ndModal.uIdx}`,
      }
      
      // Usar la función unificada que recibe la entidad del dominio
      setDeliveryUnitByEntity(deliveryUnit)
      
      closeNdModal()
    } finally {
      setSubmittingEvidence(false)
    }
  }






  const visits = routeData?.visits ?? []

  // Función centralizada para determinar qué marcador debe estar posicionado
  const getPositionedVisitIndex = (): number | null => {
    // Siempre obtener la siguiente pendiente real
    const nextPending = getNextPendingVisitIndex()
    
    // Debug: logs removidos para limpiar la consola
    
    // PRIORIDAD 1: Selección manual desde botón de mapa (lastCenteredVisit)
    // Esta debe tener prioridad absoluta cuando se establece manualmente
    if (lastCenteredVisit !== null) {
              // Debug: logs removidos para limpiar la consola
      return lastCenteredVisit
    }
    
    // PRIORIDAD 2: Estado sincronizado si es reciente (últimos 30 segundos)
    if (markerPosition && (Date.now() - markerPosition.timestamp) < 30000) {
              // Debug: logs removidos para limpiar la consola
      return markerPosition.visitIndex
    }
    
    // PRIORIDAD 3: Selección manual de cualquier visita
    if (nextVisitIndex !== null) {
              // Debug: logs removidos para limpiar la consola
      return nextVisitIndex
    }
    
    // PRIORIDAD 4: Fallback - siguiente pendiente automática
            // Debug: logs removidos para limpiar la consola
    if (nextPending !== null) {
              // Debug: logs removidos para limpiar la consola
    } else {
              // Debug: logs removidos para limpiar la consola
    }
    return nextPending
  }



  // Siguiente visita pendiente (busca la siguiente en orden secuencial, no solo la primera)
  const getNextPendingVisitIndex = (): number | null => {
    if (!visits || visits.length === 0) return null
    
    // Crear array de visitas con su estado de pendiente y número de secuencia
    const visitStatus = (visits as any[]).map((visit: any, vIdx: number) => {
      const hasPending = (visit?.orders || []).some((order: any, oIdx: number) =>
        (order?.deliveryUnits || []).some((_u: any, uIdx: number) => getDeliveryUnitStatus(vIdx, oIdx, uIdx) === undefined)
      )
      
      // Debug: logs removidos para limpiar la consola
      
      return {
        index: vIdx,
        sequenceNumber: visit?.sequenceNumber || vIdx + 1,
        hasPending
      }
    })
    
    // Filtrar solo las que tienen elementos pendientes
    const pendingVisits = visitStatus.filter(v => v.hasPending)
    
    if (pendingVisits.length === 0) return null
    
    // Si hay una visita seleccionada manualmente y tiene pendientes, mantenerla
    if (nextVisitIndex !== null) {
      const selectedVisit = visitStatus[nextVisitIndex]
      if (selectedVisit?.hasPending) {
        return nextVisitIndex
      }
    }
    
    // Si hay una última visita centrada, buscar la siguiente después de esa
    if (lastCenteredVisit !== null) {
      const lastCenteredSequence = visitStatus[lastCenteredVisit]?.sequenceNumber
      if (lastCenteredSequence) {
        // Buscar la siguiente visita pendiente después de la última centrada
        const nextPending = pendingVisits
          .filter(v => v.sequenceNumber > lastCenteredSequence)
          .sort((a, b) => a.sequenceNumber - b.sequenceNumber)[0]
        
        if (nextPending) return nextPending.index
      }
    }
    
    // Si no hay contexto previo, devolver la primera pendiente por orden de secuencia
    return pendingVisits.sort((a, b) => a.sequenceNumber - b.sequenceNumber)[0].index
  }

  useEffect(() => {
    setNextVisitIndex(getNextPendingVisitIndex())
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(localState), JSON.stringify((visits || []).map((v: any) => v?.orders?.length))])

  // Función para abrir modal de descarga
  const openDownloadModal = () => {
    setDownloadModal(true)
  }

  // Función para cerrar modal de descarga
  const closeDownloadModal = () => {
    setDownloadModal(false)
  }

  // Función para generar y descargar reporte en formato especificado
  const downloadReport = (format: 'csv' | 'excel') => {
    // Debug: Si hay problemas con los datos, descomenta esta línea
    debugDeliveryData(localState?.s || {}, routeId)
    try {
      const routeLicense = getRouteLicenseFromState(localState?.s || {}, routeId)
      
      // Generar datos del reporte usando las utilidades
      const reportData: ReportData = {
        routeId,
        routeDbId,
        routeLicense: routeLicense || routeData?.vehicle?.plate,
        visits: visits || [],
        localState: localState?.s || {}
      }
      
      const units = generateReportData(visits || [], localState?.s || {}, routeId)
      
      // Generar contenido según el formato
      let content: string
      let filename: string
      let mimeType: string
      
      const now = new Date()
      const timestamp = now.toISOString().slice(0, 19).replace(/[:.]/g, '-')
      
      if (format === 'excel') {
        content = generateExcelContent(units, reportData)
        filename = `Reporte_Ruta_${routeId}_${timestamp}.xls`
        mimeType = 'application/vnd.ms-excel;charset=utf-8;'
      } else {
        content = generateCSVContent(units, reportData)
        filename = `Reporte_Ruta_${routeId}_${timestamp}.csv`
        mimeType = 'text/csv;charset=utf-8;'
      }
      
      // Descargar archivo
      downloadFile(content, filename, mimeType)
      
      console.log(`📊 Reporte ${format.toUpperCase()} descargado:`, filename)
      
      // Cerrar modal
      closeDownloadModal()
      
    } catch (error) {
      console.error('Error generando reporte:', error)
      alert('Error al generar el reporte. Por favor intenta nuevamente.')
    }
  }

  const openNextNavigation = (provider: 'google' | 'waze' | 'geo' = 'google') => {
    const nextIdx = getNextPendingVisitIndex()
    if (typeof nextIdx !== 'number') return
    const visit = (visits as any)[nextIdx]
    const c = visit?.addressInfo?.coordinates
    const name = visit?.addressInfo?.contact?.fullName || 'Destino'
    const address = visit?.addressInfo?.addressLine1
    const latlng = Array.isArray(c?.point)
      ? [c.point[1] as number, c.point[0] as number]
      : (typeof c?.latitude === 'number' && typeof c?.longitude === 'number'
          ? [c.latitude as number, c.longitude as number]
          : null)
    let url = ''
    if (provider === 'waze' && latlng) {
      url = `https://waze.com/ul?ll=${latlng[0]},${latlng[1]}&navigate=yes`
    } else if (provider === 'geo' && latlng) {
      const label = encodeURIComponent(name)
      url = `geo:${latlng[0]},${latlng[1]}?q=${latlng[0]},${latlng[1]}(${label})`
    } else {
      // Google Maps por defecto
      if (latlng) {
        url = `https://www.google.com/maps/dir/?api=1&destination=${latlng[0]},${latlng[1]}&travelmode=driving`
      } else if (typeof address === 'string' && address.length > 0) {
        url = `https://www.google.com/maps/dir/?api=1&destination=${encodeURIComponent(address)}&travelmode=driving`
      }
    }
    if (url) {
      try { window.open(url, '_blank', 'noopener,noreferrer') } catch {}
    }
  }
  
  // Construir una lista plana de unidades de entrega para agrupar por estado
  const allUnits: Array<any> = (visits || []).flatMap((visit: any, vIdx: number) =>
    (visit?.orders || []).flatMap((order: any, oIdx: number) =>
      (order?.deliveryUnits || []).map((unit: any, uIdx: number) => ({
        visit,
        order,
        unit,
        vIdx,
        oIdx,
        uIdx,
        status: getDeliveryUnitStatus(vIdx, oIdx, uIdx),
      }))
    )
  )

  const inRouteUnits = allUnits.filter((u) => !u.status)
  const deliveredUnits = allUnits.filter((u) => u.status === 'delivered')
  const notDeliveredUnits = allUnits.filter((u) => u.status === 'not-delivered')

  const shouldRenderByTab = (status?: 'delivered' | 'not-delivered') => {
    if (activeTab === 'entregados') return status === 'delivered'
    if (activeTab === 'no-entregados') return status === 'not-delivered'
    return !status // en-ruta
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 pb-8">
      {/* Sidebar */}
      <Sidebar
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        routeStarted={routeStarted}
        onDownloadReport={openDownloadModal}
        syncInfo={syncInfo}
        markerPosition={markerPosition}
      />
      <div className="bg-gradient-to-r from-indigo-600 to-purple-600 text-white p-4 shadow-lg">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center space-x-3">
            <button 
              onClick={() => setSidebarOpen(true)}
              className="bg-white/10 backdrop-blur-sm rounded-lg p-2 hover:bg-white/20 transition-colors duration-200"
              aria-label="Abrir menú"
            >
              <Menu className="w-5 h-5" />
            </button>
            <div>
              <h1 className="text-lg font-bold flex items-center">
                <Route className="w-5 h-5 mr-2" />
                ID RUTA: {routeDbId || routeId}
              </h1>
              <p className="text-indigo-100 text-sm flex items-center">
                <Truck className="w-3 h-3 mr-1" />
                PATENTE: 
                <span className="bg-white/20 text-white px-2 py-1 rounded-lg ml-2 font-mono text-xs">
                  {getRouteLicenseFromState(localState?.s || {}, routeId) || (routeData?.vehicle?.plate ?? '—')}
                </span>
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {/* Mostrar botón mapa solo cuando la ruta esté iniciada */}
            {routeStarted && (
              <button
                onClick={() => setViewMode((m) => (m === 'list' ? 'map' : 'list'))}
                className="bg-white/10 hover:bg-white/20 text-white px-3 py-2 rounded-lg font-medium transition-all duration-200 text-sm active:scale-95 flex items-center space-x-2"
                aria-label="Alternar mapa/lista"
              >
                <Map className="w-4 h-4" />
                <span>{viewMode === 'list' ? 'Mapa' : 'Lista'}</span>
              </button>
            )}
            {!routeStarted ? (
              <button
                onClick={handleStartRoute}
                className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg font-medium flex items-center space-x-2 transition-all duration-200 text-sm active:scale-95"
              >
                <Play className="w-4 h-4" />
                <span>Iniciar</span>
              </button>
            ) : null}
          </div>
        </div>
      </div>

      {/* Vista de mapa (solo cuando viewMode === 'map') */}
      {viewMode === 'map' && (
        <MapView
          routeId={routeId}
          routeData={routeData}
          visits={visits}
          routeStarted={routeStarted}
          getDeliveryUnitStatus={getDeliveryUnitStatus}
          getNextPendingVisitIndex={getNextPendingVisitIndex}
          getPositionedVisitIndex={getPositionedVisitIndex}
          nextVisitIndex={nextVisitIndex}
          lastCenteredVisit={lastCenteredVisit}
          markerPosition={markerPosition}
          openDeliveryFor={openDeliveryFor}
          openNonDeliveryFor={openNonDeliveryFor}
          onDownloadReport={openDownloadModal}
          setNextVisitIndex={setNextVisitIndex}
          setLastCenteredVisit={setLastCenteredVisit}
          setMarkerPosition={setMarkerPosition}
          openNextNavigation={openNextNavigation}
        />
      )}

      {/* Tabs sticky: En ruta | Entregados | No entregados (ocultas en modo mapa) */}
      {viewMode === 'list' && (
        <VisitTabs
          activeTab={activeTab}
          onTabChange={setActiveTab}
          inRouteUnits={inRouteUnits.length}
          deliveredUnits={deliveredUnits.length}
          notDeliveredUnits={notDeliveredUnits.length}
        />
      )}

      {viewMode === 'list' && (
      <div className="p-4 space-y-4">
        {/* Sección "Siguiente a Entregar" - solo en pestaña "En ruta" */}
        {activeTab === 'en-ruta' && (() => {
          const nextIdx = getPositionedVisitIndex()
          if (typeof nextIdx !== 'number') return null
          const nextVisit: any = (visits as any)[nextIdx]
          if (!nextVisit) return null
          
          // Solo mostrar si tiene elementos pendientes para el tab actual
          const pendingForTab = (nextVisit?.orders || []).reduce(
            (acc: number, order: any, orderIndex: number) => {
              const countInOrder = (order?.deliveryUnits || []).reduce(
                (a: number, _unit: any, uIdx: number) =>
                  a + (shouldRenderByTab(getDeliveryUnitStatus(nextIdx, orderIndex, uIdx)) ? 1 : 0),
                0
              )
              return acc + countInOrder
            },
            0
          )
          
          if (pendingForTab === 0) return null
          
          return (
            <NextVisitCard
              nextVisit={nextVisit}
              nextIdx={nextIdx}
              onCenterOnVisit={(visitIndex: number) => {
                setViewMode('map')
                setLastCenteredVisit(visitIndex)
                setNextVisitIndex(null) // Limpiar selección automática para dar prioridad a la manual
              }}
            />
          )
        })()}
        
        {visits.map((visit: any, visitIndex: number) => (
          <VisitCard
            key={visitIndex}
            visit={visit}
            visitIndex={visitIndex}
            routeStarted={routeStarted}
            onCenterOnVisit={(visitIndex: number) => {
              setViewMode('map')
              setLastCenteredVisit(visitIndex)
              setNextVisitIndex(null) // Limpiar selección automática para dar prioridad a la manual
            }}
            onOpenDelivery={openDeliveryFor}
            onOpenNonDelivery={openNonDeliveryFor}
            getDeliveryUnitStatus={getDeliveryUnitStatus}
            shouldRenderByTab={shouldRenderByTab}
          />
        ))}
      </div>
      )}

      {/* En modo mapa: mostrar la visita seleccionada o la siguiente pendiente debajo del mapa */}

      

      {/* Barra inferior de progreso eliminada por redundancia con la barra superior */}
    {/* Modal de evidencia */}
    <DeliveryModal
      isOpen={evidenceModal.open}
      onClose={closeEvidenceModal}
      onSubmit={submitEvidence}
      initialDeliveryEvent={initialDeliveryEvent}
      submitting={submittingEvidence}
      routeData={routeData}
      visitIndex={evidenceModal.vIdx || undefined}
      orderIndex={evidenceModal.oIdx || undefined}
      unitIndex={evidenceModal.uIdx || undefined}
    />
    {/* Modal de No Entregado */}
    <NonDeliveryModal
      isOpen={ndModal.open}
      onClose={closeNdModal}
      onSubmit={submitNonDelivery}
      initialDeliveryEvent={initialNonDeliveryEvent}
      submitting={submittingEvidence}
      routeData={routeData}
      visitIndex={ndModal.vIdx || undefined}
      orderIndex={ndModal.oIdx || undefined}
      unitIndex={ndModal.uIdx || undefined}
    />


    {/* Modal de inicio de ruta */}
    <RouteStartModal
      isOpen={routeStartModal}
      onClose={() => setRouteStartModal(false)}
      onConfirm={handleLicenseConfirm}
      defaultLicense={routeData?.vehicle?.plate}
    />

    {/* Modal de descarga de reporte */}
    <DownloadReportModal
      isOpen={downloadModal}
      onClose={closeDownloadModal}
      onDownloadReport={downloadReport}
    />
    </div>
  )
}