/* eslint-disable @typescript-eslint/no-explicit-any */
import { useLiveQuery } from '@tanstack/react-db'
import { useParams } from '@tanstack/react-router'
import { createRoutesCollection } from './db/create-routes-collection'
import { 
  useDriverState, 
  useRouteStartedSync,
  routeStartedKey, 
  setRouteStarted as setRouteStartedLocal, 
  setDeliveryStatus, 
  getDeliveryStatusFromState, 
  setDeliveryEvidence, 
  setNonDeliveryEvidence,
  setRouteLicense,
  getRouteLicenseFromState
} from './db/driver-gun-state'
import { useMemo, useState, useEffect, useRef } from 'react'
import { CheckCircle, XCircle, Play, Package, User, MapPin, Crosshair, Menu, Truck, Route, Map } from 'lucide-react'
import { Sidebar, DeliveryModal, NonDeliveryModal, VisitCard, NextVisitCard, DownloadReportModal, RouteStartModal, VisitTabs } from './components'
import { 
  generateReportData, 
  generateCSVContent, 
  generateExcelContent, 
  downloadFile,
  type ReportData 
} from './components/DownloadReportModal.utils'


// Componente para rutas específicas del driver
export function RouteComponent() {
  // Obtener el routeId de los parámetros de la ruta usando TanStack Router
  const { routeId } = useParams({ from: '/driver/routes/$routeId' })
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || ''
  const routes = useMemo(() => createRoutesCollection(token, routeId), [token, routeId])
  const { data } = useLiveQuery((query) => query.from({ route: routes }))

  return (
    <div>
      {/* Renderizar UI si hay datos */}
      {(() => {
        const d: any = data as any
        const route = Array.isArray(d?.route) ? d.route[0] : d?.route ?? (Array.isArray(d) ? d[0] : d)
        const raw = route?.raw
        const routeDbId = route?.id // ID de la base de datos
        return raw ? (
          <DeliveryRouteView routeId={routeId} routeData={raw} routeDbId={routeDbId} />
        ) : (
          <pre>{JSON.stringify(data, (_key, value) => (typeof value === 'bigint' ? value.toString() : value), 2)}</pre>
        )
      })()}
    </div>
  )
}

type DeliveryRouteRaw = {
  vehicle?: { plate?: string; startLocation?: { addressInfo?: any } }
  visits?: Array<any>
  geometry?: { encoding?: string; type?: string; value?: string }
}

function DeliveryRouteView({ routeId, routeData, routeDbId }: { routeId: string; routeData: DeliveryRouteRaw; routeDbId?: number }) {
  const [activeTab, setActiveTab] = useState<'en-ruta' | 'entregados' | 'no-entregados'>('en-ruta')
  const [viewMode, setViewMode] = useState<'list' | 'map'>('list')
  // fullscreen deshabilitado para evitar cambios por clic en el mapa
  const [nextVisitIndex, setNextVisitIndex] = useState<number | null>(null)
  const [lastCenteredVisit, setLastCenteredVisit] = useState<number | null>(null) // Recordar última visita centrada
  const mapRef = useRef<HTMLDivElement | null>(null)
  const mapInstanceRef = useRef<any>(null)
  const [mapReady, setMapReady] = useState(false)
  const [forceUpdateCounter, setForceUpdateCounter] = useState(0)

  // Estado para el pin de GPS del conductor
  const [gpsActive, setGpsActive] = useState(false)
  const gpsMarkerRef = useRef<any>(null)
  const gpsCircleRef = useRef<any>(null)

  // Modal de evidencia
  const [evidenceModal, setEvidenceModal] = useState<{ open: boolean; vIdx: number | null; oIdx: number | null; uIdx: number | null }>({ open: false, vIdx: null, oIdx: null, uIdx: null })
  const [submittingEvidence, setSubmittingEvidence] = useState(false)
  const [ndModal, setNdModal] = useState<{ open: boolean; vIdx: number | null; oIdx: number | null; uIdx: number | null }>({ open: false, vIdx: null, oIdx: null, uIdx: null })



  // Modal de descarga de reporte
  const [downloadModal, setDownloadModal] = useState(false)
  
  // Estado del sidebar
  const [sidebarOpen, setSidebarOpen] = useState(false)

  // Estado local reactivo via GunJS
  const { data: localState } = useDriverState()
  
  // Debug: Log cuando cambia el estado local (comentado en producción)
  // useEffect(() => {
  //   console.log('🔄 localState cambió:', localState?.s ? Object.keys(localState.s).filter(k => k.includes('delivery:')) : 'no state')
  // }, [localState])
  // Información de sincronización entre dispositivos
  const syncInfo = useRouteStartedSync(routeId)

  // Función para sincronizar posición del marcador entre dispositivos
  const setMarkerPosition = async (routeId: string, visitIndex: number, coordinates: [number, number]) => {
    try {
      const { driverData } = await import('./db/driver-gun-state')
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
      console.log('📍 Sincronizando posición de marcador:', data)
      driverData.get(key).put(JSON.stringify(data))
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
    import('./db/driver-gun-state').then(({ driverData }) => {
      const key = `marker_position:${routeId}`
      const unsubscribe = driverData.get(key).on((data) => {
        if (data && typeof data === 'string') {
          try {
            const parsed = JSON.parse(data)
            console.log('📡 Recibida posición de marcador sincronizada:', parsed)
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

  const handleStartRoute = () => {
    setRouteStartModal(true)
  }

  const handleLicenseConfirm = (license: string) => {
    if (!license.trim()) {
      return
    }
    
    // Guardar la patente ingresada en GunJS para sincronización
    setRouteLicense(routeId, license.trim())
    
    // Iniciar la ruta con la patente ingresada (no necesita coincidir)
    setRouteStartedLocal(routeId, true)
    setRouteStartModal(false)
  }



  // Nota: setDeliveryStatus se usa directamente en cada flujo

  const getDeliveryUnitStatus = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    const status = getDeliveryStatusFromState(localState?.s || {}, routeId, visitIndex, orderIndex, unitIndex)
    // Debug ocasional para verificar lecturas de estado
    if (visitIndex === 7 && orderIndex === 0 && unitIndex === 0) {
      console.log(`🔍 getDeliveryUnitStatus(7,0,0): ${status}`)
    }
    return status
  }

  const openDeliveryFor = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    setEvidenceModal({ open: true, vIdx: visitIndex, oIdx: orderIndex, uIdx: unitIndex })
  }

  const openNonDeliveryFor = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    setNdModal({ open: true, vIdx: visitIndex, oIdx: orderIndex, uIdx: unitIndex })
  }

  const closeNdModal = () => {
    setNdModal({ open: false, vIdx: null, oIdx: null, uIdx: null })
  }

  const closeEvidenceModal = () => {
    setEvidenceModal({ open: false, vIdx: null, oIdx: null, uIdx: null })
    setSubmittingEvidence(false)
  }

  



  const submitEvidence = async (evidence: { recipientName: string; recipientRut: string; photoDataUrl: string }) => {
    if (!evidenceModal.open || evidenceModal.vIdx === null || evidenceModal.oIdx === null || evidenceModal.uIdx === null) return
    try {
      setSubmittingEvidence(true)
      
      console.log('💾 Guardando evidencia de entrega para:', { routeId, vIdx: evidenceModal.vIdx, oIdx: evidenceModal.oIdx, uIdx: evidenceModal.uIdx })
      
      setDeliveryEvidence(routeId, evidenceModal.vIdx, evidenceModal.oIdx, evidenceModal.uIdx, {
        recipientName: evidence.recipientName,
        recipientRut: evidence.recipientRut,
        photoDataUrl: evidence.photoDataUrl,
        takenAt: Date.now(),
      } as any)
      console.log('📦 Estableciendo estado de entrega a "delivered"')
      setDeliveryStatus(routeId, evidenceModal.vIdx, evidenceModal.oIdx, evidenceModal.uIdx, 'delivered')
      closeEvidenceModal()
      // Actualizar marcadores manteniendo control manual
      advanceToNextAfterDelivery()
    } finally {
      setSubmittingEvidence(false)
    }
  }

  const submitNonDelivery = async (evidence: { reason: string; observations: string; photoDataUrl: string }) => {
    if (!ndModal.open || ndModal.vIdx === null || ndModal.oIdx === null || ndModal.uIdx === null) return
    try {
      setSubmittingEvidence(true)
      
      console.log('💾 Guardando evidencia de no entrega para:', { routeId, vIdx: ndModal.vIdx, oIdx: ndModal.oIdx, uIdx: ndModal.uIdx })
      
      setNonDeliveryEvidence(routeId, ndModal.vIdx, ndModal.oIdx, ndModal.uIdx, {
        reason: evidence.reason,
        observations: evidence.observations,
        photoDataUrl: evidence.photoDataUrl,
        takenAt: Date.now(),
      } as any)
      console.log('📦 Estableciendo estado de entrega a "not-delivered"')
      setDeliveryStatus(routeId, ndModal.vIdx, ndModal.oIdx, ndModal.uIdx, 'not-delivered')
      closeNdModal()
      // Actualizar marcadores manteniendo control manual
      advanceToNextAfterDelivery()
    } finally {
      setSubmittingEvidence(false)
    }
  }



  const getStatusColor = (status?: 'delivered' | 'not-delivered') => {
    switch (status) {
      case 'delivered':
        return 'text-green-600 bg-green-50 border-green-200'
      case 'not-delivered':
        return 'text-red-600 bg-red-50 border-red-200'
      default:
        return 'text-gray-600 bg-white border-gray-200'
    }
  }



  const visits = routeData?.visits ?? []
  
  // Obtener el estado de una visita completa
  const getVisitStatus = (visitIndex: number): 'completed' | 'not-delivered' | 'partial' | 'pending' => {
    const visit = (visits as any)[visitIndex]
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

  // Obtener color del marcador según el estado de la visita (sin considerar posicionamiento)
  const getVisitMarkerColor = (visitIndex: number): string => {
    const status = getVisitStatus(visitIndex)
    switch (status) {
      case 'completed':
        // Completamente entregado (verde)
        return '#10B981'
      case 'not-delivered':
        // Completamente no entregado (rojo)
        return '#EF4444'
      case 'partial':
        // Parcialmente entregado (azul más oscuro)
        return '#1D4ED8'
      case 'pending':
      default:
        // Pendiente (gris por defecto)
        return '#6B7280'
    }
  }

  // Función centralizada para determinar qué marcador debe estar posicionado
  const getPositionedVisitIndex = (): number | null => {
    // Siempre obtener la siguiente pendiente real
    const nextPending = getNextPendingVisitIndex()
    
    console.log('🔍 getPositionedVisitIndex - Estados:', {
      markerPosition: markerPosition?.visitIndex,
      markerPositionAge: markerPosition ? Date.now() - markerPosition.timestamp : null,
      nextVisitIndex,
      lastCenteredVisit,
      nextPending
    })
    
    // Priorizar estado sincronizado si es reciente (últimos 30 segundos)
    // CAMBIO: Permitir sincronización de cualquier visita, no solo pendientes
    if (markerPosition && (Date.now() - markerPosition.timestamp) < 30000) {
      console.log('📍 Usando posición sincronizada (cualquier estado):', markerPosition.visitIndex)
      return markerPosition.visitIndex
    }
    
    // Priorizar selección manual de cualquier visita
    // CAMBIO: Permitir selección manual de entregados/no entregados también
    if (nextVisitIndex !== null) {
      console.log('📍 Usando selección manual (cualquier estado):', nextVisitIndex)
      return nextVisitIndex
    }
    
    // Usar última visita centrada de cualquier estado
    // CAMBIO: Permitir usar última centrada aunque esté entregada/no entregada
    if (lastCenteredVisit !== null) {
      console.log('📍 Usando última centrada (cualquier estado):', lastCenteredVisit)
      return lastCenteredVisit
    }
    
    // Fallback: siguiente pendiente automática
    console.log('📍 Usando siguiente pendiente automática:', nextPending)
    if (nextPending !== null) {
      console.log(`✅ Retornando visita ${nextPending} como posicionada`)
    } else {
      console.log(`❌ No hay visitas pendientes para posicionar`)
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
      
      // Debug: Log estado de cada visita
      if (vIdx <= 5) { 
        console.log(`🔍 Visita ${vIdx}: hasPending=${hasPending}, seq=${visit?.sequenceNumber || vIdx + 1}`)
      }
      
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

  // Mantener sincronizado el índice de "siguiente por entregar"
  const markersRef = useRef<any[]>([])
  
  useEffect(() => {
    setNextVisitIndex(getNextPendingVisitIndex())
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(localState), JSON.stringify((visits || []).map((v: any) => v?.orders?.length))])

    // Helper para obtener gradiente complementario
    const getGradientColor = (baseColor: string): string => {
      const colorMap: Record<string, string> = {
        '#10B981': '#059669', // Verde claro -> Verde más oscuro
        '#EF4444': '#DC2626', // Rojo (no entregado) -> Rojo más oscuro
        '#1D4ED8': '#1E40AF', // Azul oscuro (parcial) -> Azul más oscuro
        '#6B7280': '#4B5563', // Gris -> Gris más oscuro
      }
      return colorMap[baseColor] || '#7C3AED'
    }

    // Icono circular normal para visitas
  const createNumberedIcon = (L: any, number: number, color = '#4F46E5') => {
      const gradientColor = getGradientColor(color)
      
      return L.divIcon({
        html: `
          <div style="
            background: linear-gradient(135deg, ${color}, ${gradientColor});
            color: white;
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            font-size: 12px;
            box-shadow: 0 4px 8px rgba(0,0,0,0.25);
            border: 2px solid white;
          ">${number}</div>
        `,
        className: 'custom-div-icon',
        iconSize: [32, 32],
        iconAnchor: [16, 16],
      })
    }

    // Icono en forma de marcador/pin de mapa para visita posicionada (CSS puro)
  const createPositionedIcon = (L: any, number: number, color = '#4F46E5') => {
      const gradientColor = getGradientColor(color)
      
      return L.divIcon({
        html: `
          <div style="
            position: relative;
            width: 36px;
            height: 46px;
            display: flex;
            align-items: flex-start;
            justify-content: center;
          ">
            <!-- Círculo superior del pin -->
            <div style="
              width: 36px;
              height: 36px;
              border-radius: 50% 50% 50% 0;
              transform: rotate(-45deg);
              background: linear-gradient(135deg, ${color}, ${gradientColor});
              border: 3px solid white;
              box-shadow: 0 4px 8px rgba(0,0,0,0.3);
              position: relative;
            "></div>
            <!-- Número centrado -->
            <div style="
              position: absolute;
              top: 8px;
              left: 50%;
              transform: translateX(-50%);
              color: white;
              font-weight: 700;
              font-size: 12px;
              z-index: 100;
              text-shadow: 0 1px 2px rgba(0,0,0,0.7);
              pointer-events: none;
            ">${number}</div>
          </div>
        `,
        className: 'custom-div-icon positioned',
        iconSize: [36, 46],
        iconAnchor: [18, 42], // Ancla en la punta del pin
      })
    }

  // Función optimizada para actualizar solo los marcadores sin recrear el mapa
  const updateMapMarkers = () => {
    // console.log('🚀 updateMapMarkers ejecutándose')
    const L = (window as any)?.L
    if (!L || !mapInstanceRef.current) {
      // console.log('❌ updateMapMarkers abortado - no hay L o mapInstance')
      return
    }

    // console.log('🧹 Limpiando marcadores existentes:', markersRef.current.length)
    // Limpiar marcadores existentes
    markersRef.current.forEach(marker => {
      try { mapInstanceRef.current.removeLayer(marker) } catch {}
    })
    markersRef.current = []

    // Helper: obtener [lat, lng] desde addressInfo
    const getLatLngFromAddressInfo = (addr: any): [number, number] | null => {
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

    // Recrear marcadores con estados actualizados
    // console.log('🔄 Recreando marcadores para', (visits || []).length, 'visitas')
    ;(visits || []).forEach((v: any, idx: number) => {
      const latlng = getLatLngFromAddressInfo(v?.addressInfo)
      if (latlng) {
        // Determinar si está posicionada usando función centralizada
        const positionedVisitIndex = getPositionedVisitIndex()
        const isCurrentlyPositioned = (positionedVisitIndex === idx)
        
        if (isCurrentlyPositioned) {
          console.log(`📍 Marcador ${idx} posicionado (único)`)
        }
        
        const color = getVisitMarkerColor(idx)
        const sequenceNumber = v?.sequenceNumber || (idx + 1)
        
        // Debug para identificar problema de colores
        const visitStatus = getVisitStatus(idx)
        console.log(`🎨 Marcador ${idx}: status=${visitStatus}, color=${color}, positioned=${isCurrentlyPositioned}`)
        
        // Usar iconos optimizados
        const icon = isCurrentlyPositioned 
          ? createPositionedIcon(L, sequenceNumber, color)
          : createNumberedIcon(L, sequenceNumber, color)
        
        const marker = L.marker(latlng as any, { icon }).addTo(mapInstanceRef.current)
        
        // Agregar event listener para click en marcador
        marker.on('click', () => {
          console.log(`🖱️ Click en marcador ${idx}`)
          // Vibración táctil si está disponible
          try { (navigator as any)?.vibrate?.(30) } catch {}
          
          // Sincronizar posición con otros dispositivos
          setMarkerPosition(routeId, idx, latlng)
          
          // Actualizar estado local para cambiar al marcador clickeado
          setNextVisitIndex(idx)
          setLastCenteredVisit(idx)
          
          // Centrar el mapa en la nueva posición con una transición suave
          try { 
            mapInstanceRef.current.flyTo(latlng as any, 16, { duration: 0.4 }) 
          } catch {}
        })
        
        // Agregar tooltip con información de la visita
        const visitInfo = v?.addressInfo?.contact?.fullName || `Visita ${sequenceNumber}`
        marker.bindTooltip(visitInfo, {
          permanent: false,
          direction: 'top',
          offset: [0, -20]
        })
        
        markersRef.current.push(marker)
      }
    })
    // console.log('✅ updateMapMarkers completado -', markersRef.current.length, 'marcadores creados')
  }

  // Inicialización dinámica de Leaflet y render del mapa con visitas
  const initializeLeafletMap = () => {
    if (typeof window === 'undefined') return
    const L = (window as any).L
    if (!L || !mapRef.current) return
    if (mapInstanceRef.current) {
      try { mapInstanceRef.current.remove() } catch {}
      mapInstanceRef.current = null
    }

    // Helper: obtener [lat, lng] desde addressInfo (acepta point [lng,lat] o {latitude,longitude})
    const getLatLngFromAddressInfo = (addr: any): [number, number] | null => {
      const c = addr?.coordinates
      if (!c) return null
      if (Array.isArray(c?.point) && c.point.length >= 2 && typeof c.point[0] === 'number' && typeof c.point[1] === 'number') {
        return [c.point[1] as number, c.point[0] as number]
      }
      if (typeof c.latitude === 'number' && typeof c.longitude === 'number') {
        return [c.latitude as number, c.longitude as number]
      }
      return null
    }

    // Extraer waypoints desde startLocation y visitas
    const startLatLng = getLatLngFromAddressInfo(routeData?.vehicle?.startLocation?.addressInfo)
    const points: Array<[number, number]> = [
      ...((visits || [])
        .map((v: any) => getLatLngFromAddressInfo(v?.addressInfo))
        .filter((p: any): p is [number, number] => Array.isArray(p))),
    ]
    const nextIdx = getNextPendingVisitIndex()

    // Determinar el centro inicial: última visita centrada, siguiente pendiente, o primera visita
    const centerIdx = lastCenteredVisit !== null ? lastCenteredVisit : 
                     (typeof nextIdx === 'number' ? nextIdx : 0)
    const defaultCenter: [number, number] = points[centerIdx] ?? points[0] ?? [-33.45, -70.66] // Santiago fallback
    const map = L.map(mapRef.current).setView(defaultCenter, points.length ? 16 : 12)
    map.attributionControl.setPrefix(false)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
    }).addTo(map)

    // Limpiar marcadores existentes del ref antes de crear nuevos
    markersRef.current = []

    // Marcador de inicio (opcional)
    if (startLatLng) {
      const startMarker = L.marker(startLatLng as any, { icon: createNumberedIcon(L, 0, '#10B981') }).addTo(map)
      markersRef.current.push(startMarker)
    }

    // Marcadores de visitas con colores según estado
    points.forEach((latlng, idx) => {
      // Determinar si esta visita está actualmente posicionada usando función centralizada
      const positionedVisitIndex = getPositionedVisitIndex()
      const isCurrentlyPositioned = (positionedVisitIndex === idx)
      
      const color = getVisitMarkerColor(idx)
      const sequenceNumber = (visits as any)[idx]?.sequenceNumber || (idx + 1)
      
      // Usar forma diferente para visita posicionada
      const icon = isCurrentlyPositioned 
        ? createPositionedIcon(L, sequenceNumber, color)
        : createNumberedIcon(L, sequenceNumber, color)
      
      const marker = L.marker(latlng as any, { icon }).addTo(map)
      
      // Agregar event listener para click en marcador
      marker.on('click', () => {
        console.log(`🖱️ Click en marcador inicial ${idx}`)
        // Vibración táctil si está disponible
        try { (navigator as any)?.vibrate?.(30) } catch {}
        
        // Sincronizar posición con otros dispositivos
        setMarkerPosition(routeId, idx, latlng)
        
        // Actualizar estado local para cambiar al marcador clickeado
        setNextVisitIndex(idx)
        setLastCenteredVisit(idx)
        
        // Centrar el mapa en la nueva posición con una transición suave
        try { 
          map.flyTo(latlng as any, 16, { duration: 0.4 }) 
        } catch {}
      })
      
      // Agregar tooltip con información de la visita
      const visit = (visits as any)[idx]
      const visitInfo = visit?.addressInfo?.contact?.fullName || `Visita ${sequenceNumber}`
      marker.bindTooltip(visitInfo, {
        permanent: false,
        direction: 'top',
        offset: [0, -20]
      })
      
      markersRef.current.push(marker)
    })

    // Ruta (polyline)
    // Decodificador de polylines (Google Encoded Polyline Algorithm Format)
    const decodePolyline = (encoded: string): Array<[number, number]> => {
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

    const encoded = (routeData as any)?.geometry?.encoding === 'polyline' ? (routeData as any)?.geometry?.value : undefined
    let routeLatLngs: Array<[number, number]> | null = null
    if (typeof encoded === 'string' && encoded.length > 0) {
      try {
        const decoded = decodePolyline(encoded)
        if (decoded.length >= 2) {
          routeLatLngs = decoded
        }
      } catch {}
    }

    const linePoints = routeLatLngs ?? (points.length >= 2 ? points : null)
    if (linePoints) {
      const line = L.polyline(linePoints as any, {
        color: '#4F46E5',
        weight: 4,
        opacity: 0.85,
        dashArray: '10,5',
      }).addTo(map)
      // Mantener la posición centrada si existe, si no, usar la siguiente pendiente o ajustar a la ruta
      if (lastCenteredVisit !== null && points[lastCenteredVisit]) {
        map.setView(points[lastCenteredVisit] as any, 16)
      } else if (typeof nextIdx === 'number' && points[nextIdx]) {
        map.setView(points[nextIdx] as any, 16)
      } else {
        map.fitBounds(line.getBounds(), { padding: [24, 24] })
      }
    } else if (points.length > 0 || startLatLng) {
      const group = L.featureGroup([
        ...points.map((p) => L.marker(p as any)),
        ...(startLatLng ? [L.marker(startLatLng as any)] : []),
      ])
      
      // Mantener la posición centrada si existe, si no, ajustar a todos los puntos
      if (lastCenteredVisit !== null && points[lastCenteredVisit]) {
        map.setView(points[lastCenteredVisit] as any, 16)
      } else {
        map.fitBounds(group.getBounds(), { padding: [24, 24] })
      }
    }

    mapInstanceRef.current = map
    setMapReady(true)
  }

  useEffect(() => {
    // Cargar Leaflet dinámicamente y luego inicializar
    if (typeof window === 'undefined') return
    if (viewMode !== 'map') return
    if (!(window as any).L) {
      setMapReady(false)
      const link = document.createElement('link')
      link.rel = 'stylesheet'
      link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
      document.head.appendChild(link)

      const script = document.createElement('script')
      script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
      script.onload = () => setTimeout(initializeLeafletMap, 50)
      document.body.appendChild(script)
    } else {
      setMapReady(false)
      setTimeout(initializeLeafletMap, 0)
    }

    return () => {
      if (mapInstanceRef.current) {
        try { mapInstanceRef.current.remove() } catch {}
        mapInstanceRef.current = null
      }
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [viewMode, nextVisitIndex, JSON.stringify((visits || []).map((v: any) => v?.addressInfo?.coordinates?.point))])

  // Optimización: Solo re-render cuando cambian datos esenciales del mapa
  const mapEssentialData = useMemo(() => {
    if (!localState?.s) return null
    // Solo incluir las claves de estado que afectan los marcadores del mapa
    // Usar el formato correcto de las claves: delivery:routeId:vIdx-oIdx-uIdx
    const essentialKeys = Object.keys(localState.s).filter(key => 
      key.startsWith(`delivery:${routeId}:`)
    )
    const result = essentialKeys.map(key => `${key}=${localState.s[key]}`).join(',')
    console.log('🗺️ mapEssentialData update:', { 
      essentialKeys: essentialKeys.length, 
      sampleKey: essentialKeys[0], 
      sampleValue: essentialKeys[0] ? localState.s[essentialKeys[0]] : null,
      result: result.substring(0, 100) + (result.length > 100 ? '...' : ''),
      timestamp: Date.now()
    })
    return result
  }, [localState?.s, routeId])

  useEffect(() => {
    console.log('🏗️ Main map useEffect disparado:', { viewMode, mapEssentialDataPreview: mapEssentialData?.substring(0, 50) })
    if (viewMode !== 'map') return
    setMapReady(false)
    setTimeout(initializeLeafletMap, 0)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [viewMode, mapEssentialData, lastCenteredVisit, nextVisitIndex])

  // UseEffect optimizado para actualizar solo marcadores cuando cambia el estado de entrega
  useEffect(() => {
    console.log('🔄 useEffect para updateMapMarkers disparado:', { 
      viewMode, 
      hasMapInstance: !!mapInstanceRef.current, 
      mapEssentialDataLength: mapEssentialData?.length || 0,
      nextVisitIndex,
      lastCenteredVisit
    })
    if (viewMode === 'map' && mapInstanceRef.current) {
      console.log('✅ Ejecutando updateMapMarkers desde useEffect')
      updateMapMarkers()
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mapEssentialData, nextVisitIndex, lastCenteredVisit, markerPosition, forceUpdateCounter])

  // UseEffect para reaccionar a cambios de posición sincronizada desde otros dispositivos
  useEffect(() => {
    if (markerPosition && viewMode === 'map' && mapInstanceRef.current) {
      const { visitIndex, coordinates, deviceId, timestamp } = markerPosition
      
      // Solo reaccionar si la posición viene de otro dispositivo y es reciente
      const currentDeviceId = syncInfo?.deviceId || 'unknown'
      const isFromOtherDevice = deviceId !== currentDeviceId
      const isRecent = (Date.now() - timestamp) < 10000 // 10 segundos
      
      console.log('📡 Evaluando posición sincronizada:', { 
        isFromOtherDevice, 
        isRecent, 
        currentDeviceId, 
        senderDeviceId: deviceId 
      })
      
      if (isFromOtherDevice && isRecent) {
        console.log('🔄 Aplicando posición sincronizada desde otro dispositivo')
        // Actualizar estado local para mostrar la posición sincronizada
        setNextVisitIndex(visitIndex)
        setLastCenteredVisit(visitIndex)
        
        // Centrar el mapa en la posición sincronizada con transición suave
        try { 
          mapInstanceRef.current.flyTo(coordinates as any, 16, { duration: 0.6 }) 
        } catch {}
        
        // Vibración suave para notificar cambio
        try { (navigator as any)?.vibrate?.(50) } catch {}
      }
    }
  }, [markerPosition, viewMode, syncInfo?.deviceId])



  // Funciones para manejar el GPS del conductor
  const startGPS = () => {
    if (!navigator.geolocation) {
      alert('El GPS no está disponible en este dispositivo')
      return
    }

    setGpsActive(true)
    
    const options = {
      enableHighAccuracy: true,
      timeout: 10000,
      maximumAge: 0
    }

    const success = (position: GeolocationPosition) => {
      const { latitude, longitude, accuracy } = position.coords
      
      // Actualizar marcador en el mapa si está disponible
      if (mapInstanceRef.current && (window as any)?.L) {
        updateGPSMarker([latitude, longitude], accuracy)
      }
    }

    const error = (err: GeolocationPositionError) => {
      console.error('Error GPS:', err)
      setGpsActive(false)
      let message = 'Error al obtener ubicación'
      switch (err.code) {
        case err.PERMISSION_DENIED:
          message = 'Permiso de ubicación denegado'
          break
        case err.POSITION_UNAVAILABLE:
          message = 'Información de ubicación no disponible'
          break
        case err.TIMEOUT:
          message = 'Tiempo de espera agotado'
          break
      }
      alert(message)
    }

    // Iniciar seguimiento continuo
    navigator.geolocation.watchPosition(success, error, options)
  }

  const stopGPS = () => {
    setGpsActive(false)
    
    // Remover marcador del mapa
    if (mapInstanceRef.current && gpsMarkerRef.current) {
      try {
        mapInstanceRef.current.removeLayer(gpsMarkerRef.current)
        gpsMarkerRef.current = null
      } catch {}
    }
    if (mapInstanceRef.current && gpsCircleRef.current) {
      try {
        mapInstanceRef.current.removeLayer(gpsCircleRef.current)
        gpsCircleRef.current = null
      } catch {}
    }
  }

  const updateGPSMarker = (latlng: [number, number], accuracy: number) => {
    const L = (window as any)?.L
    if (!L || !mapInstanceRef.current) return

    // Remover marcador anterior si existe
    if (gpsMarkerRef.current) {
      try { mapInstanceRef.current.removeLayer(gpsMarkerRef.current) } catch {}
    }
    if (gpsCircleRef.current) {
      try { mapInstanceRef.current.removeLayer(gpsCircleRef.current) } catch {}
    }

    // Crear marcador de GPS con icono personalizado
    const gpsIcon = L.divIcon({
      html: `
        <div style="
          position: relative;
          width: 24px;
          height: 24px;
          display: flex;
          align-items: center;
          justify-content: center;
        ">
          <!-- Círculo principal pulsante -->
          <div style="
            width: 24px;
            height: 24px;
            border-radius: 50%;
            background: #00D4AA;
            border: 3px solid white;
            box-shadow: 0 0 0 3px #00D4AA;
            animation: gps-pulse 2s infinite;
            position: relative;
          "></div>
          <!-- Punto central -->
          <div style="
            position: absolute;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: white;
            border: 2px solid #00D4AA;
          "></div>
          <!-- Indicador de dirección -->
          <div style="
            position: absolute;
            top: -2px;
            left: 50%;
            transform: translateX(-50%);
            width: 0;
            height: 0;
            border-left: 4px solid transparent;
            border-right: 4px solid transparent;
            border-bottom: 6px solid #00D4AA;
          "></div>
        </div>
        <style>
          @keyframes gps-pulse {
            0% { box-shadow: 0 0 0 0 rgba(0, 212, 170, 0.7); }
            70% { box-shadow: 0 0 0 10px rgba(0, 212, 170, 0); }
            100% { box-shadow: 0 0 0 0 rgba(0, 212, 170, 0); }
          }
        </style>
      `,
      className: 'gps-marker',
      iconSize: [24, 24],
      iconAnchor: [12, 12],
    })

    // Crear marcador
    gpsMarkerRef.current = L.marker(latlng as any, { icon: gpsIcon }).addTo(mapInstanceRef.current)
    
    // Crear círculo de precisión
    gpsCircleRef.current = L.circle(latlng as any, {
      radius: accuracy,
      color: '#00D4AA',
      fillColor: '#00D4AA',
      fillOpacity: 0.1,
      weight: 1,
      opacity: 0.6
    }).addTo(mapInstanceRef.current)

    // Tooltip con información del GPS
    gpsMarkerRef.current.bindTooltip(`
      <div class="text-center">
        <div class="font-bold text-green-700">Tu ubicación</div>
        <div class="text-xs text-gray-600">Precisión: ${Math.round(accuracy)}m</div>
        <div class="text-xs text-gray-500">${latlng[0].toFixed(6)}, ${latlng[1].toFixed(6)}</div>
      </div>
    `, {
      permanent: false,
      direction: 'top',
      offset: [0, -15]
    })
  }

  // Función para hacer zoom al punto actualmente seleccionado/posicionado (sin cambiar selección)
  const zoomToCurrentlySelected = () => {
    const L = (window as any)?.L
    if (!L || !mapInstanceRef.current) return
    
    // Obtener el índice del marcador actualmente posicionado
    const currentSelectedIdx = getPositionedVisitIndex()
    if (typeof currentSelectedIdx !== 'number') {
      console.log('📍 No hay marcador posicionado para hacer zoom')
      return
    }
    
    console.log(`🔍 Haciendo zoom al marcador actualmente seleccionado: ${currentSelectedIdx}`)
    
    // Obtener latlng de la visita seleccionada
    const visit = (visits as any)[currentSelectedIdx]
    const c = visit?.addressInfo?.coordinates
    const latlng = Array.isArray(c?.point)
      ? [c.point[1] as number, c.point[0] as number]
      : (typeof c?.latitude === 'number' && typeof c?.longitude === 'number'
          ? [c.latitude as number, c.longitude as number]
          : null)
    if (latlng) {
      try { 
        mapInstanceRef.current.flyTo(latlng as any, 16, { duration: 0.6 })
        console.log(`✅ Zoom realizado a visita ${currentSelectedIdx}`)
      } catch (e) {
        console.error('❌ Error al hacer zoom:', e)
      }
    } else {
      console.log('❌ No se pudo obtener coordenadas para la visita')
    }
  }



  const centerOnVisit = (visitIndex: number) => {
    // Obtener latlng de la visita específica
    const visit = (visits as any)[visitIndex]
    if (!visit) return
    
    const c = visit?.addressInfo?.coordinates
    const latlng = Array.isArray(c?.point)
      ? [c.point[1] as number, c.point[0] as number]
      : (typeof c?.latitude === 'number' && typeof c?.longitude === 'number'
          ? [c.latitude as number, c.longitude as number]
          : null)
    
    if (latlng && latlng.length === 2) {
      // Sincronizar la posición del marcador seleccionado entre dispositivos
      console.log('🎯 Usuario seleccionó "Ver en mapa" para visita', visitIndex)
      setMarkerPosition(routeId, visitIndex, latlng as [number, number])
      
      // Cambiar a modo mapa primero
      setViewMode('map')
      // Guardar la visita seleccionada para centrarse después
      setNextVisitIndex(visitIndex)
      // Guardar la última visita centrada
      setLastCenteredVisit(visitIndex)
      
      // Vibración táctil si está disponible
      try { (navigator as any)?.vibrate?.(50) } catch {}
      
      // Función para centrar el mapa cuando esté listo
      const attemptCenter = (attempts = 0) => {
        const L = (window as any)?.L
        if (L && mapInstanceRef.current) {
          try { 
            mapInstanceRef.current.flyTo(latlng as any, 16, { duration: 0.8 }) 
          } catch {}
        } else if (attempts < 10) {
          // Reintentar hasta que el mapa esté listo
          setTimeout(() => attemptCenter(attempts + 1), 200)
        }
      }
      
      // Iniciar el intento de centrado
      setTimeout(() => attemptCenter(), 100)
    }
  }

  const advanceToNextAfterDelivery = () => {
    // Actualizar marcadores sin mover la vista (funciona en cualquier modo)
    console.log('🔄 Actualizando después de gestionar entrega (sin mover mapa)')
    
    // Esperar un poco para que el estado se actualice después de la entrega
    setTimeout(() => {
      console.log('🧹 Manteniendo visita actual después de gestionar entrega...')
      
      // MANTENER la visita actual seleccionada para evitar salto automático
      // Obtener la visita que se está mostrando actualmente
      const currentDisplayedVisit = getPositionedVisitIndex()
      
      if (typeof currentDisplayedVisit === 'number') {
        // Forzar que se mantenga la visita actual como seleccionada
        console.log('🔒 Fijando visita actual como seleccionada para evitar salto automático:', currentDisplayedVisit)
        setNextVisitIndex(currentDisplayedVisit)
        setLastCenteredVisit(currentDisplayedVisit)
      }
      
      // Verificar cuál sería el siguiente punto pendiente para logs
      const nextPendingVisit = getNextPendingVisitIndex()
      console.log('📍 Siguiente punto pendiente disponible:', nextPendingVisit, '(pero no saltando automáticamente)')
      
      // Forzar múltiples actualizaciones para asegurar que el estado se refleje
      if (mapInstanceRef.current) {
        // Primera actualización inmediata
        setTimeout(() => {
          updateMapMarkers()
          console.log('🔄 Primera actualización de marcadores')
        }, 50)
        
        // Segunda actualización después de más tiempo para asegurar que GunJS se haya sincronizado
        setTimeout(() => {
          updateMapMarkers()
          console.log('🔄 Segunda actualización de marcadores (post-GunJS)')
        }, 300)
        
        // Tercera actualización como respaldo
        setTimeout(() => {
          updateMapMarkers()
          console.log('🔄 Tercera actualización de marcadores (respaldo)')
        }, 600)
        
        // Forzar re-render del useEffect después de que todo haya sido procesado
        setTimeout(() => {
          setForceUpdateCounter(prev => prev + 1)
          console.log('🔄 Forzando re-render con counter')
        }, 800)
      }
    }, 200) // Pausa más larga para asegurar que GunJS ha procesado el cambio
  }

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
      {viewMode === 'map' && (() => {
      // Siempre mostrar el mapa, incluso cuando todas las entregas estén gestionadas
      return (
        <div className="px-4 pt-4">
          <div className="relative">
            <div
              ref={mapRef}
              className={`h-72 w-full rounded-xl overflow-hidden shadow-md bg-gray-100`}
              style={{ zIndex: 1 }}
            >
              {!mapReady && (
                <div className="absolute inset-0 flex items-center justify-center bg-gradient-to-br from-blue-100 to-indigo-100">
                  <div className="text-center">
                    <div className="animate-spin rounded-full h-10 w-10 border-4 border-indigo-500 border-t-transparent mx-auto mb-3"></div>
                    <p className="text-indigo-600 text-sm font-medium">Cargando mapa…</p>
                  </div>
                </div>
              )}
            </div>
            {/* Controles flotantes del mapa */}
            <div className="absolute top-3 right-3 space-y-2" style={{ zIndex: 1000 }}>
              {/* Botón de GPS del conductor */}
              <button
                onClick={gpsActive ? stopGPS : startGPS}
                className={`w-10 h-10 rounded-lg shadow-lg flex items-center justify-center transition-all ${
                  gpsActive 
                    ? 'bg-green-500 text-white hover:bg-green-600' 
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                } hover:shadow-xl`}
                aria-label={gpsActive ? 'Desactivar GPS' : 'Activar GPS'}
                title={gpsActive ? 'Desactivar GPS' : 'Activar GPS del conductor'}
              >
                <div className={`w-5 h-5 ${gpsActive ? 'animate-pulse' : ''}`}>
                  <svg 
                    viewBox="0 0 24 24" 
                    fill="none" 
                    stroke="currentColor" 
                    strokeWidth="2" 
                    strokeLinecap="round" 
                    strokeLinejoin="round"
                    className="w-full h-full"
                  >
                    <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
                    <circle cx="12" cy="9" r="2.5"/>
                  </svg>
                </div>
              </button>
              
              <button
                onClick={zoomToCurrentlySelected}
                className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-gray-700 hover:bg-gray-50 hover:shadow-xl transition-all"
                aria-label="Zoom al punto seleccionado"
                title="Hacer zoom al punto seleccionado"
              >
                <Crosshair className="w-5 h-5" />
              </button>
              <div className="flex flex-col gap-2">
                <button
                  onClick={() => openNextNavigation('google')}
                  className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-blue-600 hover:bg-gray-50 hover:shadow-xl transition-all"
                  aria-label="Navegar con Google Maps"
                  title="Google Maps"
                >
                  G
                </button>
                <button
                  onClick={() => openNextNavigation('waze')}
                  className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-indigo-600 hover:bg-gray-50 hover:shadow-xl transition-all"
                  aria-label="Navegar con Waze"
                  title="Waze"
                >
                  W
                </button>
              </div>
            </div>
          </div>
        </div>
      )
      })()}

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
              onCenterOnVisit={centerOnVisit}
            />
          )
        })()}
        
        {visits.map((visit: any, visitIndex: number) => (
          <VisitCard
            key={visitIndex}
            visit={visit}
            visitIndex={visitIndex}
            routeStarted={routeStarted}
            onCenterOnVisit={centerOnVisit}
            onOpenDelivery={openDeliveryFor}
            onOpenNonDelivery={openNonDeliveryFor}
            getDeliveryUnitStatus={getDeliveryUnitStatus}
            shouldRenderByTab={shouldRenderByTab}
          />
        ))}
      </div>
      )}

      {/* En modo mapa: mostrar la visita seleccionada o la siguiente pendiente debajo del mapa */}
      {viewMode === 'map' && (() => {
        // Usar la misma lógica que el mapa para determinar qué visita mostrar
        const displayIdx = getPositionedVisitIndex()
        
        // Si no hay punto seleccionado/posicionado, mostrar mensaje de ruta completada
        if (typeof displayIdx !== 'number') {
          return (
            <div className="p-4 space-y-4">
              <div className="bg-gradient-to-r from-green-50 to-emerald-50 rounded-xl border-2 border-green-200 p-6 text-center">
                <div className="flex items-center justify-center mb-3">
                  <CheckCircle className="w-8 h-8 text-green-600" />
                </div>
                <h2 className="text-lg font-semibold text-gray-800 mb-2">¡Ruta Completada!</h2>
                <p className="text-sm text-gray-600 mb-4">Todas las entregas han sido gestionadas exitosamente.</p>
                <p className="text-xs text-gray-500 mb-4">El mapa muestra el estado final de todas las visitas.</p>
                
                {/* Botón de descarga CSV */}
                <button
                  onClick={openDownloadModal}
                  className="inline-flex items-center space-x-2 bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white px-6 py-3 rounded-lg font-medium transition-all duration-200 shadow-md hover:shadow-lg active:scale-95"
                >
                  📊
                  <span>Descargar Reporte</span>
                </button>
              </div>
            </div>
          )
        }
        
        const visit: any = (visits as any)[displayIdx]
        // Es seleccionada si no es solo la automática (siguiente pendiente)
        const autoNext = getNextPendingVisitIndex()
        const isSelectedVisit = displayIdx !== autoNext || nextVisitIndex !== null || lastCenteredVisit !== null || (markerPosition && (Date.now() - markerPosition.timestamp) < 30000)
        
        // Debug para modo mapa
        console.log('🗺️ Modo mapa - Determinando qué mostrar:', {
          displayIdx,
          autoNext,
          isSelectedVisit,
          nextVisitIndex,
          lastCenteredVisit,
          hasRecentMarkerPosition: markerPosition && (Date.now() - markerPosition.timestamp) < 30000
        })
        // Verificar si la visita actual ya está procesada
        const visitStatus = getVisitStatus(displayIdx)
        const isProcessed = visitStatus === 'completed' || visitStatus === 'not-delivered' || visitStatus === 'partial'
        const nextPendingIdx = getNextPendingVisitIndex()
        const hasNextPending = typeof nextPendingIdx === 'number' && nextPendingIdx !== displayIdx
        
        return (
          <div className="p-4 space-y-4">
            {/* Sección "Siguiente a Entregar" cuando la visita actual está procesada */}
            {isProcessed && hasNextPending && (
              <div className="bg-gradient-to-r from-green-50 to-blue-50 rounded-xl border-2 border-green-200 p-4 mb-4">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-sm font-bold text-green-800 flex items-center">
                    <CheckCircle className="w-4 h-4 mr-2" />
                    ¡Gestión Completada!
                  </h3>
                  <span className="text-xs text-green-600 bg-green-100 px-2 py-1 rounded-full font-medium">
                    ✓ Procesado
                  </span>
                </div>
                <button
                  onClick={() => {
                    // Ir específicamente al siguiente pendiente y marcarlo como seleccionado
                    setNextVisitIndex(nextPendingIdx)
                    setLastCenteredVisit(nextPendingIdx)
                    console.log('🎯 Usuario presionó "Siguiente a Entregar" - saltando a visita:', nextPendingIdx)
                    
                    // También sincronizar la posición del marcador
                    if (nextPendingIdx !== null) {
                      const nextVisit = (visits as any)[nextPendingIdx]
                      const c = nextVisit?.addressInfo?.coordinates
                      const latlng = Array.isArray(c?.point)
                        ? [c.point[1] as number, c.point[0] as number]
                        : (typeof c?.latitude === 'number' && typeof c?.longitude === 'number'
                            ? [c.latitude as number, c.longitude as number]
                            : null)
                      if (latlng && latlng.length === 2) {
                        setMarkerPosition(routeId, nextPendingIdx, latlng as [number, number])
                      }
                    }
                  }}
                  className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white py-3 px-4 rounded-lg font-medium flex items-center justify-center space-x-2 transition-all duration-200 shadow-md hover:shadow-lg active:scale-95"
                >
                  <Play className="w-4 h-4" />
                  <span>Siguiente a Entregar (#{(visits as any)[nextPendingIdx]?.sequenceNumber})</span>
                </button>
              </div>
            )}
            
            {/* Indicador de qué visita se está mostrando */}
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-medium text-gray-700">
                {isSelectedVisit ? 'Visita seleccionada' : 'Siguiente a entregar'}
              </h3>
              {isSelectedVisit && !isProcessed && (
                <button
                  onClick={() => {
                    // Limpiar todas las selecciones para volver al automático
                    setNextVisitIndex(null)
                    setLastCenteredVisit(null)
                    console.log('🔄 Usuario solicitó ver siguiente - limpiando selecciones')
                  }}
                  className="text-xs text-blue-600 hover:text-blue-800 font-medium"
                >
                  Ver siguiente
                </button>
              )}
            </div>
            
            <div className={`bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden border ${
              isSelectedVisit ? 'border-blue-300 ring-2 ring-blue-100' : 'border-gray-100'
            }`}>
              <div className="p-4 border-b border-gray-100">
                <div className="flex items-start space-x-3">
                  <div className={`w-8 h-8 rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0 text-white ${
                    isSelectedVisit 
                      ? 'bg-gradient-to-br from-indigo-500 to-purple-600' 
                      : 'bg-gradient-to-br from-indigo-500 to-purple-600'
                  }`}>
                    {visit.sequenceNumber}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                      <User className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
                      <span className="truncate">{visit.addressInfo?.contact?.fullName}</span>
                    </h3>
                    <p className="text-xs text-gray-600 flex items-start mb-2">
                      <MapPin className="w-3 h-3 mr-1 mt-0.5 text-gray-500 flex-shrink-0" />
                      <span className="line-clamp-2">{visit.addressInfo?.addressLine1}</span>
                    </p>
                  </div>
                </div>
              </div>
              <div className="p-4">
                <h4 className="text-sm font-medium text-gray-800 mb-3 flex items-center">
                  <Package size={18} />
                  <span className="ml-2">Unidades de Entrega:</span>
                </h4>
                {(visit.orders || []).map((order: any, orderIndex: number) => (
                  <div key={orderIndex} className="mb-4">
                    <div className="mb-2">
                      <span className="inline-block bg-gradient-to-r from-orange-400 to-red-500 text-white px-2 py-1 rounded-lg text-xs font-medium">
                        {order.referenceID}
                      </span>
                    </div>
                    {(order.deliveryUnits || [])
                      .map((unit: any, uIdx: number): { unit: any; uIdx: number; status: 'delivered' | 'not-delivered' | undefined } => ({
                        unit,
                        uIdx,
                        status: getDeliveryUnitStatus(displayIdx, orderIndex, uIdx),
                      }))
                      .map(({ unit, uIdx, status }: { unit: any; uIdx: number; status: 'delivered' | 'not-delivered' | undefined }) => (
                        <div key={uIdx} className={`bg-gradient-to-r from-gray-50 to-blue-50 rounded-lg p-3 border ${getStatusColor(status).replace('bg-white ', '')}`}>
                          <div className="flex justify-between items-start mb-2">
                            <div className="flex-1 min-w-0">
                              <h5 className="text-sm font-medium text-gray-800 mb-2 truncate">Unidad de Entrega {uIdx + 1}</h5>
                              {Array.isArray(unit.items) && unit.items.length > 0 && (
                                <div className="flex items-center space-x-1 mb-2">
                                  <span className="w-1.5 h-1.5 bg-indigo-500 rounded-full"></span>
                                  <span className="text-xs text-gray-700 truncate">{unit.items[0]?.description}</span>
                                </div>
                              )}
                              <div className="flex items-center space-x-3 text-xs text-gray-600">
                                <span className="flex items-center">
                                  <span className="w-1.5 h-1.5 bg-green-500 rounded-full mr-1"></span>
                                  {typeof unit.weight === 'number' ? `${unit.weight}kg` : unit.weight}
                                </span>
                                <span className="flex items-center">
                                  <span className="w-1.5 h-1.5 bg-blue-500 rounded-full mr-1"></span>
                                  {typeof unit.volume === 'number' ? `${unit.volume}m³` : unit.volume}
                                </span>
                              </div>
                            </div>
                            <div className="text-right ml-3">
                              <span className="text-xs text-gray-500 block">Cant.</span>
                              <span className="text-xl font-bold text-indigo-600">{(unit.items || []).reduce((a: number, it: any) => a + (Number(it?.quantity) || 0), 0)}</span>
                            </div>
                          </div>
                          {routeStarted && (
                            <div className="flex space-x-2 mt-3">
                              {/* Permitir cambios de estado en vista mapa siempre */}
                              {status === 'delivered' ? (
                                // Si está entregado, mostrar solo opción de cambiar a no entregado
                                <button
                                  onClick={() => openNonDeliveryFor(displayIdx, orderIndex, uIdx)}
                                  className="w-full flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-red-100 text-red-700 hover:bg-red-200"
                                >
                                  <XCircle size={16} />
                                  <span>Cambiar a no entregado</span>
                                </button>
                              ) : status === 'not-delivered' ? (
                                // Si está no entregado, mostrar solo opción de cambiar a entregado
                                <button
                                  onClick={() => openDeliveryFor(displayIdx, orderIndex, uIdx)}
                                  className="w-full flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-green-100 text-green-700 hover:bg-green-200"
                                >
                                  <CheckCircle size={16} />
                                  <span>Cambiar a entregado</span>
                                </button>
                              ) : (
                                // Si está pendiente, mostrar ambas opciones originales
                                <>
                                  <button
                                    onClick={() => openDeliveryFor(displayIdx, orderIndex, uIdx)}
                                    className="flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-green-100 text-green-700 hover:bg-green-200"
                                  >
                                    <CheckCircle size={16} />
                                    <span>entregar</span>
                                  </button>
                                  <button
                                    onClick={() => openNonDeliveryFor(displayIdx, orderIndex, uIdx)}
                                    className="flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-red-100 text-red-700 hover:bg-red-200"
                                  >
                                    <XCircle size={16} />
                                    <span>no entregado</span>
                                  </button>
                                </>
                              )}
                            </div>
                          )}
                        </div>
                      ))}
                  </div>
                ))}
              </div>
            </div>
          </div>
        )
      })()}
      

      

      {/* Barra inferior de progreso eliminada por redundancia con la barra superior */}
    {/* Modal de evidencia */}
    <DeliveryModal
      isOpen={evidenceModal.open}
      onClose={closeEvidenceModal}
      onSubmit={submitEvidence}
      submitting={submittingEvidence}
    />
    {/* Modal de No Entregado */}
    <NonDeliveryModal
      isOpen={ndModal.open}
      onClose={closeNdModal}
      onSubmit={submitNonDelivery}
      submitting={submittingEvidence}
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