import { useRef, useEffect, useState } from 'react'
import { CheckCircle } from 'lucide-react'
import { MapControls } from './MapControls'
import { MapVisitCard } from './MapVisitCard'
import { useLanguage } from '../hooks/useLanguage'
import { 
  getLatLngFromAddressInfo, 
  decodePolyline, 
  getGradientColor, 
  getVisitStatus, 
  getVisitMarkerColor 
} from './MapView.utils'
import type { Visit } from '../domain/route'

interface MapViewProps {
  routeId: string
  routeData: any
  visits: Visit[]
  routeStarted: boolean
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined
  getNextPendingVisitIndex: () => number | null
  getPositionedVisitIndex: () => number | null
  nextVisitIndex: number | null
  lastCenteredVisit: number | null
  markerPosition: any
  openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void
  openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void
  onDownloadReport: () => void
  setNextVisitIndex: (index: number | null) => void
  setLastCenteredVisit: (index: number | null) => void
  setMarkerPosition: (routeId: string, visitIndex: number, coordinates: [number, number]) => Promise<void>
  openNextNavigation: (provider: 'google' | 'waze' | 'geo') => void
  // Nuevas props para agrupación
  openGroupedDelivery?: (visitIndex: number, group: any) => void
  openGroupedNonDelivery?: (visitIndex: number, group: any) => void
  // Props para entregar todo
  onDeliverAll?: (visitIndex: number) => void
  onNonDeliverAll?: (visitIndex: number) => void
  // Props para selección de cliente
  selectedClientIndex?: number | null
  onClientSelect?: (clientIndex: number | null) => void
}

export function MapView({
  routeId,
  routeData,
  visits,
  routeStarted,
  getDeliveryUnitStatus,
  getNextPendingVisitIndex,
  getPositionedVisitIndex,
  nextVisitIndex,
  lastCenteredVisit,
  markerPosition,
  openDeliveryFor,
  openNonDeliveryFor,
  onDownloadReport,
  setNextVisitIndex,
  setLastCenteredVisit,
  setMarkerPosition,
  openNextNavigation,
  openGroupedDelivery,
  openGroupedNonDelivery,
  onDeliverAll,
  onNonDeliverAll,
  selectedClientIndex,
  onClientSelect
}: MapViewProps) {
  const { t } = useLanguage()
  const mapRef = useRef<HTMLDivElement | null>(null)
  const mapInstanceRef = useRef<any>(null)
  const [mapReady, setMapReady] = useState(false)
  const [forceUpdateCounter] = useState(0)
  
  // Estados para manejo de selección de clientes
  const [selectedClientName, setSelectedClientName] = useState<string | null>(null)
  const [wasManuallySelected, setWasManuallySelected] = useState(false)
  
  // Detectar si hay múltiples clientes en la misma dirección
  const getClientsAtSameLocation = () => {
    console.log('🔍 getClientsAtSameLocation DEBUG:', {
      wasManuallySelected,
      lastCenteredVisit,
      selectedClientIndex
    })
    
    // PRIORIDAD 1: Si fue selección manual desde modo lista, usar lastCenteredVisit
    if (wasManuallySelected && lastCenteredVisit !== null) {
      console.log('✅ Usando PRIORIDAD 1: selección manual, visitIndex:', lastCenteredVisit)
      const manuallySelectedVisit = visits[lastCenteredVisit]
      if (manuallySelectedVisit) {
        const orders = manuallySelectedVisit.orders || []
        console.log('📋 Visita seleccionada tiene', orders.length, 'órdenes')
        
        // Si la visita seleccionada manualmente tiene múltiples clientes, mostrarlos
        if (orders.length > 1) {
          console.log('👥 Múltiples clientes detectados en visita seleccionada')
          const clientMap = new Map()
          
          orders.forEach((order: any, orderIndex: number) => {
            const clientName = order.contact?.fullName || 'Sin nombre'
            const hasPendingUnits = (order.deliveryUnits || []).some((_unit: any, unitIndex: number) => 
              getDeliveryUnitStatus(lastCenteredVisit, orderIndex, unitIndex) === undefined
            )
            
            if (clientMap.has(clientName)) {
              const existing = clientMap.get(clientName)
              existing.hasPendingUnits = existing.hasPendingUnits || hasPendingUnits
              existing.orderIndexes.push(orderIndex)
            } else {
              clientMap.set(clientName, {
                index: lastCenteredVisit,
                orderIndex,
                orderIndexes: [orderIndex],
                clientName,
                hasPendingUnits
              })
            }
          })
          
          return Array.from(clientMap.values()).sort((a, b) => a.clientName.localeCompare(b.clientName))
        } else {
          console.log('👤 Un solo cliente detectado en visita seleccionada - no mostrar selector')
        }
        
        // Si la visita seleccionada manualmente tiene un solo cliente, no mostrar selector
        return []
      }
    }
    
    // PRIORIDAD 2: Si hay un cliente seleccionado programáticamente, usar su visita
    if (selectedClientIndex !== null) {
      const selectedVisit = visits[selectedClientIndex]
      if (!selectedVisit) return []
      
      // Si la visita seleccionada tiene múltiples órdenes/clientes, mostrarlos
      const orders = selectedVisit.orders || []
      if (orders.length > 1) {
        const clientMap = new Map()
        
        orders.forEach((order: any, orderIndex: number) => {
          const clientName = order.contact?.fullName || 'Sin nombre'
          const hasPendingUnits = (order.deliveryUnits || []).some((_unit: any, unitIndex: number) => 
            getDeliveryUnitStatus(selectedClientIndex, orderIndex, unitIndex) === undefined
          )
          
          if (clientMap.has(clientName)) {
            const existing = clientMap.get(clientName)
            existing.hasPendingUnits = existing.hasPendingUnits || hasPendingUnits
            existing.orderIndexes.push(orderIndex)
          } else {
            clientMap.set(clientName, {
              index: selectedClientIndex,
              orderIndex,
              orderIndexes: [orderIndex],
              clientName,
              hasPendingUnits
            })
          }
        })
        
        return Array.from(clientMap.values()).sort((a, b) => a.clientName.localeCompare(b.clientName))
      }
      
      // Si solo tiene una orden, buscar otras visitas en la misma dirección
      const selectedAddress = selectedVisit.addressInfo?.addressLine1
      if (!selectedAddress) return []
      
      const clientMap = new Map()
      
      visits
        .filter(visit => visit.addressInfo?.addressLine1 === selectedAddress)
        .forEach((visit, _, filteredVisits) => {
          const visitIndex = visits.indexOf(visit)
          
          ;(visit.orders || []).forEach((order: any, orderIndex: number) => {
            const clientName = order.contact?.fullName || 'Sin nombre'
            const hasPendingUnits = (order.deliveryUnits || []).some((_unit: any, unitIndex: number) => 
              getDeliveryUnitStatus(visitIndex, orderIndex, unitIndex) === undefined
            )
            
            if (clientMap.has(clientName)) {
              const existing = clientMap.get(clientName)
              existing.hasPendingUnits = existing.hasPendingUnits || hasPendingUnits
            } else {
              clientMap.set(clientName, {
                index: visitIndex,
                orderIndex,
                clientName,
                hasPendingUnits
              })
            }
          })
        })
      
      return Array.from(clientMap.values()).sort((a, b) => a.clientName.localeCompare(b.clientName))
    }
    
    // PRIORIDAD 3: Solo si no hay selección manual, buscar automáticamente visitas con múltiples clientes
    if (!wasManuallySelected) {
      for (let visitIndex = 0; visitIndex < visits.length; visitIndex++) {
        const visit = visits[visitIndex]
        const orders = visit.orders || []
        
        if (orders.length > 1) {
          const clientMap = new Map()
          
          orders.forEach((order: any, orderIndex: number) => {
            const clientName = order.contact?.fullName || 'Sin nombre'
            const hasPendingUnits = (order.deliveryUnits || []).some((_unit: any, unitIndex: number) => 
              getDeliveryUnitStatus(visitIndex, orderIndex, unitIndex) === undefined
            )
            
            if (clientMap.has(clientName)) {
              const existing = clientMap.get(clientName)
              existing.hasPendingUnits = existing.hasPendingUnits || hasPendingUnits
              existing.orderIndexes.push(orderIndex)
            } else {
              clientMap.set(clientName, {
                index: visitIndex,
                orderIndex,
                orderIndexes: [orderIndex],
                clientName,
                hasPendingUnits
              })
            }
          })
          
          return Array.from(clientMap.values()).sort((a, b) => a.clientName.localeCompare(b.clientName))
        }
      }
      
      // Si no hay visitas con múltiples clientes, buscar múltiples visitas en la misma dirección
      const addressGroups: { [key: string]: Map<string, any> } = {}
      
      visits.forEach((visit, index) => {
        const address = visit.addressInfo?.addressLine1
        if (address) {
          if (!addressGroups[address]) {
            addressGroups[address] = new Map()
          }
          
          ;(visit.orders || []).forEach((order: any, orderIndex: number) => {
            const clientName = order.contact?.fullName || 'Sin nombre'
            const hasPendingUnits = (order.deliveryUnits || []).some((_unit: any, unitIndex: number) => 
              getDeliveryUnitStatus(index, orderIndex, unitIndex) === undefined
            )
            
            if (addressGroups[address].has(clientName)) {
              const existing = addressGroups[address].get(clientName)
              existing.hasPendingUnits = existing.hasPendingUnits || hasPendingUnits
            } else {
              addressGroups[address].set(clientName, {
                index,
                orderIndex,
                clientName,
                hasPendingUnits
              })
            }
          })
        }
      })
      
      // Encontrar la primera dirección con múltiples clientes
      for (const [address, clientMap] of Object.entries(addressGroups)) {
        if (clientMap.size > 1) {
          return Array.from(clientMap.values()).sort((a, b) => a.clientName.localeCompare(b.clientName))
        }
      }
    }
    
    return []
  }
  
  const clientsAtSameLocation = getClientsAtSameLocation()
  const hasMultipleClients = clientsAtSameLocation.length > 1
  
  
  // Detectar cuando lastCenteredVisit cambia (selección manual desde modo lista)
  useEffect(() => {
    if (lastCenteredVisit !== null) {
      console.log('🔄 MapView: lastCenteredVisit cambió a:', lastCenteredVisit)
      setWasManuallySelected(true)
      // Limpiar selección de cliente cuando se selecciona manualmente una visita
      setSelectedClientName(null)
      console.log('✅ MapView: wasManuallySelected=true, selectedClientName=null')
    }
  }, [lastCenteredVisit])
  
  // Si hay múltiples clientes pero no hay uno seleccionado, seleccionar el primero automáticamente
  // SOLO si no fue una selección manual desde modo lista
  useEffect(() => {
    if (hasMultipleClients && selectedClientIndex === null && onClientSelect && !wasManuallySelected) {
      onClientSelect(clientsAtSameLocation[0].index)
    }
  }, [hasMultipleClients, selectedClientIndex, onClientSelect, clientsAtSameLocation, wasManuallySelected])
  
  // Cuando hay múltiples clientes, seleccionar el primero por defecto
  // SOLO si no fue una selección manual desde modo lista
  useEffect(() => {
    if (hasMultipleClients && !selectedClientName && !wasManuallySelected) {
      setSelectedClientName(clientsAtSameLocation[0]?.clientName || null)
    }
  }, [hasMultipleClients, selectedClientName, clientsAtSameLocation, wasManuallySelected])
  
  // Función para obtener el cliente seleccionado
  const getSelectedClient = () => {
    if (!selectedClientName || !hasMultipleClients) return null
    return clientsAtSameLocation.find(client => client.clientName === selectedClientName) || null
  }
  
  const selectedClient = getSelectedClient()
  
  // Debug temporal
  console.log('🔍 Debug selector de clientes:', {
    selectedClientIndex,
    clientsAtSameLocation,
    hasMultipleClients,
    totalVisits: visits.length,
    addresses: visits.map((v, i) => ({ index: i, address: v.addressInfo?.addressLine1, client: v.orders?.[0]?.contact?.fullName }))
  })
  // Estado para el pin de GPS del conductor
  const [gpsActive, setGpsActive] = useState(false)
  const gpsMarkerRef = useRef<any>(null)
  const gpsCircleRef = useRef<any>(null)
  
  // Mantener sincronizado el índice de "siguiente por entregar"
  const markersRef = useRef<any[]>([])

  // Función para hacer zoom al punto actualmente seleccionado/posicionado
  const zoomToCurrentlySelected = () => {
    const L = (window as any)?.L
    if (!L || !mapInstanceRef.current) return
    
    const currentSelectedIdx = getPositionedVisitIndex()
    if (typeof currentSelectedIdx !== 'number') {
      // Debug: logs removidos para limpiar la consola
      return
    }
    
    // Debug: logs removidos para limpiar la consola
    
    const visit = visits[currentSelectedIdx]
    const c = visit?.addressInfo?.coordinates
    const latlng = (typeof c?.latitude === 'number' && typeof c?.longitude === 'number')
      ? [c.latitude as number, c.longitude as number]
      : null
    
    if (latlng) {
      try { 
        mapInstanceRef.current.flyTo(latlng as any, 16, { duration: 0.6 })
        // Debug: logs removidos para limpiar la consola
      } catch (e) {
        console.error('❌ Error al hacer zoom:', e)
      }
    } else {
      // Debug: logs removidos para limpiar la consola
    }
  }

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

  // Icono en forma de marcador/pin de mapa para visita posicionada
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
    const L = (window as any)?.L
    if (!L || !mapInstanceRef.current) {
      return
    }

    // Limpiar marcadores existentes
    markersRef.current.forEach(marker => {
      try { mapInstanceRef.current.removeLayer(marker) } catch {}
    })
    markersRef.current = []

    // Recrear marcadores con estados actualizados
    visits.forEach((v: any, idx: number) => {
      const latlng = getLatLngFromAddressInfo(v?.addressInfo)
      if (latlng) {
        // Determinar si está posicionada usando función centralizada
        const positionedVisitIndex = getPositionedVisitIndex()
        const isCurrentlyPositioned = (positionedVisitIndex === idx)
        
        if (isCurrentlyPositioned) {
          // Debug: logs removidos para limpiar la consola
        }
        
        const visitStatus = getVisitStatus(v, getDeliveryUnitStatus, idx)
        const color = getVisitMarkerColor(visitStatus)
        const sequenceNumber = v?.sequenceNumber || (idx + 1)
        
        // Debug: logs removidos para limpiar la consola
        
        // Usar iconos optimizados
        const icon = isCurrentlyPositioned 
          ? createPositionedIcon(L, sequenceNumber, color)
          : createNumberedIcon(L, sequenceNumber, color)
        
        const marker = L.marker(latlng as any, { icon }).addTo(mapInstanceRef.current)
        
        // Agregar event listener para click en marcador
        marker.on('click', () => {
          // Debug: logs removidos para limpiar la consola
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
        const uniqueClients = Array.from(new Set(
          (v?.orders || []).map((order: any) => order.contact?.fullName).filter(Boolean)
        ))
        const visitInfo = uniqueClients.length > 1 
          ? `${uniqueClients.length} clientes: ${uniqueClients.join(', ')}`
          : uniqueClients[0] || `Visita ${sequenceNumber}`
        marker.bindTooltip(visitInfo, {
          permanent: false,
          direction: 'top',
          offset: [0, -20]
        })
        
        markersRef.current.push(marker)
      }
    })
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

    // Extraer waypoints desde startLocation y visitas
    const startLatLng = getLatLngFromAddressInfo(routeData?.vehicle?.startLocation?.addressInfo)
    const points: Array<[number, number]> = [
      ...(visits
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
      
      const visitStatus = getVisitStatus(visits[idx], getDeliveryUnitStatus, idx)
      const color = getVisitMarkerColor(visitStatus)
      const sequenceNumber = visits[idx]?.sequenceNumber || (idx + 1)
      
      // Usar forma diferente para visita posicionada
      const icon = isCurrentlyPositioned 
        ? createPositionedIcon(L, sequenceNumber, color)
        : createNumberedIcon(L, sequenceNumber, color)
      
      const marker = L.marker(latlng as any, { icon }).addTo(map)
      
      // Agregar event listener para click en marcador
      marker.on('click', () => {
        // Debug: logs removidos para limpiar la consola
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
      const visit = visits[idx]
      const uniqueClients = Array.from(new Set(
        (visit?.orders || []).map((order: any) => order.contact?.fullName).filter(Boolean)
      ))
      const visitInfo = uniqueClients.length > 1 
        ? `${uniqueClients.length} clientes: ${uniqueClients.join(', ')}`
        : uniqueClients[0] || `Visita ${sequenceNumber}`
      marker.bindTooltip(visitInfo, {
        permanent: false,
        direction: 'top',
        offset: [0, -20]
      })
      
      markersRef.current.push(marker)
    })

    // Ruta (polyline)
    const encoded = routeData?.geometry?.encoding === 'polyline' ? routeData?.geometry?.value : undefined
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
  }, [nextVisitIndex, JSON.stringify(visits.map((v: any) => ({ lat: v?.addressInfo?.coordinates?.latitude, lng: v?.addressInfo?.coordinates?.longitude })))])

  // UseEffect optimizado para actualizar solo marcadores cuando cambia el estado de entrega
  useEffect(() => {
    if (mapInstanceRef.current) {
      updateMapMarkers()
    }
  }, [forceUpdateCounter, nextVisitIndex, lastCenteredVisit, markerPosition])

  // Función para forzar actualización de marcadores (comentada por no uso)
  // const forceUpdateMarkers = () => {
  //   setForceUpdateCounter(prev => prev + 1)
  // }

  // Función para manejar el siguiente pendiente
  const handleNextPending = (nextPendingIdx: number) => {
    setNextVisitIndex(nextPendingIdx)
    setLastCenteredVisit(nextPendingIdx)
    
    // También sincronizar la posición del marcador
    const nextVisit = visits[nextPendingIdx]
    const c = nextVisit?.addressInfo?.coordinates
    const latlng = (typeof c?.latitude === 'number' && typeof c?.longitude === 'number')
      ? [c.latitude as number, c.longitude as number]
      : null
    if (latlng && latlng.length === 2) {
      setMarkerPosition(routeId, nextPendingIdx, latlng as [number, number])
    }
  }

  // Función para limpiar selección
  const handleClearSelection = () => {
    setNextVisitIndex(null)
    setLastCenteredVisit(null)
  }

  // Determinar qué visita mostrar en modo mapa
  const displayIdx = getPositionedVisitIndex()
  
  // Debug temporal para investigar el problema de displayIdx
  console.log('🗺️ MapView displayIdx DEBUG:', {
    displayIdx,
    lastCenteredVisit,
    nextVisitIndex,
    selectedClientIndex,
    visitSequenceAtDisplayIdx: visits[displayIdx]?.sequenceNumber
  })
  
  // Si no hay punto seleccionado/posicionado, mostrar mensaje de ruta completada
  if (typeof displayIdx !== 'number') {
    return (
      <div className="px-4 pt-4">
        <div className="relative">
          <div
            ref={mapRef}
            className="h-72 w-full rounded-xl overflow-hidden shadow-md bg-gray-100"
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
          <MapControls
            gpsActive={gpsActive}
            onGPSToggle={gpsActive ? stopGPS : startGPS}
            onZoomToSelected={zoomToCurrentlySelected}
            onNavigate={openNextNavigation}
          />
        </div>
        
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
              onClick={onDownloadReport}
              className="inline-flex items-center space-x-2 bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white px-6 py-3 rounded-lg font-medium transition-all duration-200 shadow-md hover:shadow-lg active:scale-95"
            >
              📊
              <span>Descargar Reporte</span>
            </button>
          </div>
        </div>
      </div>
    )
  }
  
  const visit = visits[displayIdx]
  // Es seleccionada si no es solo la automática (siguiente pendiente)
  const autoNext = getNextPendingVisitIndex()
  const isSelectedVisit = displayIdx !== autoNext || nextVisitIndex !== null || lastCenteredVisit !== null || (markerPosition && (Date.now() - markerPosition.timestamp) < 30000)
  
  // Verificar si la visita actual ya está procesada
  const visitStatus = getVisitStatus(visit, getDeliveryUnitStatus, displayIdx)
  
  // Verificar si hay otras visitas en la misma dirección que aún no han sido procesadas
  const currentAddress = visit.addressInfo?.addressLine1 || 'Sin dirección'
  const visitsAtSameAddress = visits.filter(v => 
    v.addressInfo?.addressLine1 === currentAddress && v !== visit
  )
  
  // Una visita se considera procesada solo si:
  // 1. TODOS los clientes de la visita actual han sido gestionados Y
  // 2. No hay otras visitas en la misma dirección con unidades pendientes
  
  // Verificar si todos los clientes de la visita actual han sido gestionados
  const allClientsProcessed = (() => {
    if (!hasMultipleClients || !selectedClient) {
      // Si no hay múltiples clientes o no hay cliente seleccionado, usar lógica original
      return visitStatus === 'completed' || visitStatus === 'not-delivered' || visitStatus === 'partial'
    }
    
    // Si hay múltiples clientes, verificar que TODOS los clientes hayan sido gestionados
    const uniqueClients = Array.from(new Set(
      (visit.orders || []).map((order: any) => order.contact?.fullName).filter(Boolean)
    ))
    
    return uniqueClients.every(clientName => {
      // Verificar si todas las unidades de este cliente han sido gestionadas
      const clientOrders = (visit.orders || []).filter((order: any) => 
        order.contact?.fullName === clientName
      )
      
      return clientOrders.every((order: any) => {
        const orderIndex = (visit.orders || []).indexOf(order)
        return (order.deliveryUnits || []).every((_unit: any, unitIndex: number) => {
          const status = getDeliveryUnitStatus(displayIdx, orderIndex, unitIndex)
          return status === 'delivered' || status === 'not-delivered'
        })
      })
    })
  })()
  
  const hasOtherVisitsAtSameAddress = visitsAtSameAddress.length > 0
  const otherVisitsProcessed = hasOtherVisitsAtSameAddress ? 
    visitsAtSameAddress.every(otherVisit => {
      const otherVisitStatus = getVisitStatus(otherVisit, getDeliveryUnitStatus, visits.indexOf(otherVisit))
      return otherVisitStatus === 'completed' || otherVisitStatus === 'not-delivered' || otherVisitStatus === 'partial'
    }) : true
  
  const isProcessed = allClientsProcessed && otherVisitsProcessed
  
  const nextPendingIdx = getNextPendingVisitIndex()
  const hasNextPending = typeof nextPendingIdx === 'number' && nextPendingIdx !== displayIdx

  return (
    <div className="px-4 pt-4">
      <div className="relative">
        <div
          ref={mapRef}
          className="h-72 w-full rounded-xl overflow-hidden shadow-md bg-gray-100"
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
        <MapControls
          gpsActive={gpsActive}
          onGPSToggle={gpsActive ? stopGPS : startGPS}
          onZoomToSelected={zoomToCurrentlySelected}
          onNavigate={openNextNavigation}
        />
      </div>
      
      {/* Selector de cliente cuando hay múltiples clientes en la misma dirección */}
      {hasMultipleClients && onClientSelect && (
        <div className="mb-4 p-4 bg-gradient-to-r from-purple-50 to-indigo-50 rounded-xl border border-purple-200 shadow-lg">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center space-x-2">
              <div className="w-3 h-3 bg-purple-500 rounded-full"></div>
              <span className="text-sm font-semibold text-gray-800">
🏢 {t.mapView.multipleClientsAtLocation}
              </span>
            </div>
            <span className="text-xs bg-purple-100 text-purple-700 px-2 py-1 rounded-full font-medium">
{clientsAtSameLocation.length} {t.nextVisit.clients}
            </span>
          </div>
          
          <div className="grid grid-cols-1 gap-2">
            {clientsAtSameLocation.map((client, clientIdx) => (
              <button
                key={`${client.index}-${client.orderIndex || 0}`}
                onClick={() => setSelectedClientName(client.clientName)}
                className={`p-4 rounded-lg border-2 transition-all duration-200 text-left ${
                  selectedClientName === client.clientName
                    ? 'border-purple-500 bg-purple-100 shadow-md transform scale-[1.02]'
                    : 'border-gray-200 bg-white hover:border-purple-300 hover:bg-purple-50 hover:shadow-md'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className={`w-4 h-4 rounded-full ${
                      client.hasPendingUnits ? 'bg-orange-400' : 'bg-green-500'
                    }`}></div>
                    <div className="flex flex-col">
                      <span className="font-semibold text-gray-800 text-base">
                        {client.clientName}
                      </span>
                      {client.orderIndexes && client.orderIndexes.length > 1 ? (
                        <span className="text-xs text-gray-500">
{client.orderIndexes.length} {t.visitCard.orders}
                        </span>
                      ) : (
                        client.orderIndex !== undefined && (
                          <span className="text-xs text-gray-500">
{t.visitCard.order} #{client.orderIndex + 1}
                          </span>
                        )
                      )}
                    </div>
                  </div>
                  <div className={`text-xs px-2 py-1 rounded-full font-medium ${
                    client.hasPendingUnits 
                      ? 'bg-orange-100 text-orange-700' 
                      : 'bg-green-100 text-green-700'
                  }`}>
{client.hasPendingUnits ? t.status.pending : t.status.completed}
                  </div>
                </div>
              </button>
            ))}
          </div>
          
          <div className="mt-3 text-xs text-gray-600 text-center bg-white/50 rounded-lg p-2">
👆 {t.mapView.selectClientToDeliver}
          </div>
        </div>
      )}
      
      <MapVisitCard
        visit={visit}
        displayIdx={displayIdx}
        isSelectedVisit={isSelectedVisit}
        isProcessed={isProcessed}
        hasNextPending={hasNextPending}
        nextPendingIdx={nextPendingIdx}
        routeStarted={routeStarted}
        getDeliveryUnitStatus={getDeliveryUnitStatus}
        openDeliveryFor={openDeliveryFor}
        openNonDeliveryFor={openNonDeliveryFor}
        onNextPending={handleNextPending}
        onClearSelection={handleClearSelection}
        onDeliverAll={onDeliverAll}
        onNonDeliverAll={onNonDeliverAll}
        openGroupedDelivery={openGroupedDelivery}
        openGroupedNonDelivery={openGroupedNonDelivery}
        selectedClient={selectedClient}
        hasMultipleClients={hasMultipleClients}
      />
    </div>
  )
}
