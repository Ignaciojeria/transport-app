/* eslint-disable @typescript-eslint/no-explicit-any */
import { useLiveQuery, createCollection } from '@tanstack/react-db'
import { queryCollectionOptions } from '@tanstack/query-db-collection'
import { QueryClient } from '@tanstack/query-core'
import { useParams } from '@tanstack/react-router'
import { createRoutesCollection } from './db/create-routes-collection'
import { useMemo, useState, useEffect, useRef } from 'react'
import { CheckCircle, XCircle, Play, Package, User, MapPin, Crosshair } from 'lucide-react'


// Helpers mínimos para IndexedDB (persistencia offline real)
const IDB_DB_NAME = 'transport-app'
const IDB_STORE = 'driver-progress'
function openIDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(IDB_DB_NAME, 1)
    req.onupgradeneeded = () => {
      const db = req.result
      if (!db.objectStoreNames.contains(IDB_STORE)) db.createObjectStore(IDB_STORE)
    }
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
}
async function idbGet<T = unknown>(key: string): Promise<T | undefined> {
  const db = await openIDB()
  return new Promise((resolve, reject) => {
    const tx = db.transaction(IDB_STORE, 'readonly')
    const store = tx.objectStore(IDB_STORE)
    const r = store.get(key)
    r.onsuccess = () => resolve(r.result as T | undefined)
    r.onerror = () => reject(r.error)
  })
}
async function idbSet<T = unknown>(key: string, value: T): Promise<void> {
  const db = await openIDB()
  return new Promise((resolve, reject) => {
    const tx = db.transaction(IDB_STORE, 'readwrite')
    const store = tx.objectStore(IDB_STORE)
    const r = store.put(value as any, key)
    r.onsuccess = () => resolve()
    r.onerror = () => reject(r.error)
  })
}
async function idbDel(key: string): Promise<void> {
  const db = await openIDB()
  return new Promise((resolve, reject) => {
    const tx = db.transaction(IDB_STORE, 'readwrite')
    const store = tx.objectStore(IDB_STORE)
    const r = store.delete(key)
    r.onsuccess = () => resolve()
    r.onerror = () => reject(r.error)
  })
}

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
        const raw = Array.isArray(d?.route)
          ? d.route[0]?.raw
          : d?.route?.raw ?? (Array.isArray(d) ? d[0]?.raw : d?.raw)
        return raw ? (
          <DeliveryRouteView routeId={routeId} routeData={raw} />
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

function DeliveryRouteView({ routeId, routeData }: { routeId: string; routeData: DeliveryRouteRaw }) {
  const [routeStarted, setRouteStarted] = useState(false)
  const [deliveryStates, setDeliveryStates] = useState<Record<string, 'delivered' | 'not-delivered' | undefined>>({})
  const [activeTab, setActiveTab] = useState<'en-ruta' | 'entregados' | 'no-entregados'>('en-ruta')
  const [viewMode, setViewMode] = useState<'list' | 'map'>('list')
  // fullscreen por tap en el mapa (sin botón explícito)
  const [mapFullscreen, setMapFullscreen] = useState(false)
  const [nextVisitIndex, setNextVisitIndex] = useState<number | null>(null)
  const mapRef = useRef<HTMLDivElement | null>(null)
  const mapInstanceRef = useRef<any>(null)
  const [mapReady, setMapReady] = useState(false)

  // TanStack DB: colección local de progreso (persistimos en localStorage por simplicidad)
  const queryClientRef = useRef<QueryClient | null>(null)
  if (!queryClientRef.current) queryClientRef.current = new QueryClient()
  const progressCollection = useMemo(() => {
    const storageKey = `driver_progress:${routeId}`
    return createCollection(
      queryCollectionOptions<{ id: string; routeStarted: boolean; deliveryStates: Record<string, 'delivered' | 'not-delivered'>; updatedAt?: string }>({
        id: `driver_progress:${routeId}`,
        queryKey: ['driver_progress', routeId],
        queryClient: queryClientRef.current!,
        getKey: (item) => item.id,
        // Carga desde IndexedDB
        queryFn: async () => {
          try {
            const stored = await idbGet(storageKey)
            if (!stored) return []
            if (stored && typeof stored === 'object') return [stored as any]
          } catch {}
          return []
        },
        // Persistencia en IndexedDB sin refetch automático
        async onInsert({ transaction }) {
          for (const m of transaction.mutations) {
            try { await idbSet(storageKey, m.modified) } catch {}
          }
          return { refetch: false }
        },
        async onUpdate({ transaction }) {
          for (const m of transaction.mutations) {
            const current = (await idbGet(storageKey)) || {}
            const next = { ...current as any, ...(m.modified ?? {}), ...(m.changes ?? {}), updatedAt: new Date().toISOString() }
            try { await idbSet(storageKey, next) } catch {}
          }
          return { refetch: false }
        },
        async onDelete({ transaction }) {
          for (const _m of transaction.mutations) {
            try { await idbDel(storageKey) } catch {}
          }
          return { refetch: false }
        },
      })
    )
  }, [routeId])

  const { data: progressData } = useLiveQuery((query) => query.from({ progress: progressCollection }))
  const progressItem: any = useMemo(() => {
    const d: any = progressData as any
    return Array.isArray(d?.progress) ? d.progress[0] : d?.progress ?? (Array.isArray(d) ? d[0] : d)
  }, [progressData])

  // Hidratar estado local desde la colección
  useEffect(() => {
    if (!progressItem) return
    if (typeof progressItem.routeStarted === 'boolean') setRouteStarted(progressItem.routeStarted)
    if (progressItem.deliveryStates && typeof progressItem.deliveryStates === 'object') {
      setDeliveryStates(progressItem.deliveryStates as Record<string, 'delivered' | 'not-delivered' | undefined>)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(progressItem)])

  const persistProgress = (partial: { routeStarted?: boolean; deliveryStates?: Record<string, 'delivered' | 'not-delivered' | undefined> }) => {
    const next = {
      id: routeId,
      routeStarted: partial.routeStarted ?? routeStarted,
      deliveryStates: (partial.deliveryStates as any) ?? deliveryStates,
      updatedAt: new Date().toISOString(),
    }
    try {
      if (progressItem && progressItem.id) {
        ;(progressCollection as any).update(routeId, next)
      } else {
        ;(progressCollection as any).insert(next)
      }
    } catch {}
  }

  const handleStartRoute = () => {
    setRouteStarted(true)
    persistProgress({ routeStarted: true })
  }

  const markDeliveryUnit = (
    visitIndex: number,
    orderIndex: number,
    unitIndex: number,
    status: 'delivered' | 'not-delivered'
  ) => {
    const key = `${visitIndex}-${orderIndex}-${unitIndex}`
    setDeliveryStates((prev) => {
      const updated = { ...prev, [key]: status }
      persistProgress({ deliveryStates: updated })
      return updated
    })
  }

  const getDeliveryUnitStatus = (visitIndex: number, orderIndex: number, unitIndex: number) => {
    const key = `${visitIndex}-${orderIndex}-${unitIndex}`
    return deliveryStates[key]
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

  const getTabIcon = (tabId: 'en-ruta' | 'entregados' | 'no-entregados') => {
    switch (tabId) {
      case 'en-ruta':
        return <Play className="w-3 h-3" />
      case 'entregados':
        return <CheckCircle className="w-3 h-3" />
      case 'no-entregados':
        return <XCircle className="w-3 h-3" />
      default:
        return null
    }
  }

  const visits = routeData?.visits ?? []
  
  // Siguiente visita pendiente (primera con unidades sin estado)
  const getNextPendingVisitIndex = (): number | null => {
    for (let vIdx = 0; vIdx < (visits || []).length; vIdx++) {
      const visit: any = (visits as any)[vIdx]
      const hasPending = (visit?.orders || []).some((order: any, oIdx: number) =>
        (order?.deliveryUnits || []).some((_u: any, uIdx: number) => getDeliveryUnitStatus(vIdx, oIdx, uIdx) === undefined)
      )
      if (hasPending) return vIdx
    }
    return null
  }

  // Mantener sincronizado el índice de "siguiente por entregar"
  useEffect(() => {
    setNextVisitIndex(getNextPendingVisitIndex())
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(deliveryStates), JSON.stringify((visits || []).map((v: any) => v?.orders?.length))])
  
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

    const defaultCenter: [number, number] = points[0] ?? [-33.45, -70.66] // Santiago fallback
    const map = L.map(mapRef.current).setView(defaultCenter, points.length ? 14 : 12)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
    }).addTo(map)

    // Icono numerado para visitas
    const createNumberedIcon = (number: number, color = '#4F46E5') =>
      L.divIcon({
        html: `\n          <div style="\n            background: linear-gradient(135deg, ${color}, #7C3AED);\n            color: white;\n            width: 28px;\n            height: 28px;\n            border-radius: 50%;\n            display: flex;\n            align-items: center;\n            justify-content: center;\n            font-weight: 700;\n            font-size: 12px;\n            box-shadow: 0 4px 8px rgba(0,0,0,0.25);\n            border: 2px solid white;\n          ">${number}</div>\n        `,
        className: 'custom-div-icon',
        iconSize: [28, 28],
        iconAnchor: [14, 14],
      })

    // Marcador de inicio (opcional)
    if (startLatLng) {
      L.marker(startLatLng as any, { icon: createNumberedIcon(0, '#10B981') }).addTo(map)
    }

    // Marcadores de visitas
    points.forEach((latlng, idx) => {
      const isNext = typeof nextIdx === 'number' && idx === nextIdx
      L.marker(latlng as any, { icon: createNumberedIcon(idx + 1, isNext ? '#EF4444' : '#4F46E5') }).addTo(map)
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
      // Si hay siguiente, centramos ahí; si no, ajustamos a la ruta
      if (typeof nextIdx === 'number' && points[nextIdx]) {
        map.setView(points[nextIdx] as any, 16)
      } else {
        map.fitBounds(line.getBounds(), { padding: [24, 24] })
      }
    } else if (points.length > 0 || startLatLng) {
      const group = L.featureGroup([
        ...points.map((p) => L.marker(p as any)),
        ...(startLatLng ? [L.marker(startLatLng as any)] : []),
      ])
      map.fitBounds(group.getBounds(), { padding: [24, 24] })
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

  // Re-render del mapa cuando cambian los estados de entrega para recalcular "siguiente"
  useEffect(() => {
    if (viewMode !== 'map') return
    setMapReady(false)
    setTimeout(initializeLeafletMap, 0)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [viewMode, JSON.stringify(deliveryStates)])

  const centerOnNext = () => {
    const L = (window as any)?.L
    if (!L || !mapInstanceRef.current) return
    const nextIdx = getNextPendingVisitIndex()
    if (typeof nextIdx !== 'number') return
    // Obtener latlng de la visita
    const visit = (visits as any)[nextIdx]
    const c = visit?.addressInfo?.coordinates
    const latlng = Array.isArray(c?.point)
      ? [c.point[1] as number, c.point[0] as number]
      : (typeof c?.latitude === 'number' && typeof c?.longitude === 'number'
          ? [c.latitude as number, c.longitude as number]
          : null)
    if (latlng) {
      try { mapInstanceRef.current.flyTo(latlng as any, 16, { duration: 0.6 }) } catch {}
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
  type MappedUnit = {
    unit: any
    uIdx: number
    status: 'delivered' | 'not-delivered' | undefined
  }
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
      <div className="bg-gradient-to-r from-indigo-600 to-purple-600 text-white p-4 shadow-lg">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center space-x-3">
            <div className="bg-white/10 backdrop-blur-sm rounded-lg p-2">
              <Package className="w-5 h-5" />
            </div>
            <div>
              <h1 className="text-lg font-bold">Ruta de Entrega</h1>
              <p className="text-indigo-100 text-sm flex items-center">
                <Package className="w-3 h-3 mr-1" />
                {routeData?.vehicle?.plate ?? '—'}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setViewMode((m) => (m === 'list' ? 'map' : 'list'))}
              className="bg-white/10 hover:bg-white/20 text-white px-3 py-2 rounded-lg font-medium transition-all duration-200 text-sm active:scale-95"
              aria-label="Alternar mapa/lista"
            >
              {viewMode === 'list' ? 'Mapa' : 'Lista'}
            </button>
          {!routeStarted ? (
            <button
              onClick={handleStartRoute}
              className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg font-medium flex items-center space-x-2 transition-all duration-200 text-sm active:scale-95"
            >
              <Play className="w-4 h-4" />
              <span>Iniciar</span>
            </button>
          ) : (
            <div className="flex items-center space-x-2 text-green-200">
              <CheckCircle size={20} />
              <span>Ruta Iniciada</span>
            </div>
          )}
          </div>
        </div>
      </div>

      {/* Vista de mapa (solo cuando viewMode === 'map') */}
      {viewMode === 'map' && (() => {
      const nextIdxForMap = getNextPendingVisitIndex()
      if (nextIdxForMap === null) {
        return (
          <div className="px-4 pt-10">
            <div className="bg-white rounded-xl shadow-md border border-gray-100 p-6 text-center">
              <div className="flex items-center justify-center mb-2">
                <CheckCircle className="w-6 h-6 text-green-600" />
              </div>
              <h2 className="text-base font-semibold text-gray-800">Ruta terminada</h2>
              <p className="text-sm text-gray-600 mt-1">Todas las entregas fueron gestionadas.</p>
            </div>
          </div>
        )
      }
      return (
        <div className="px-4 pt-4">
          <div className="relative">
            <div
              ref={mapRef}
              className={`${mapFullscreen ? 'h-[70vh]' : 'h-72'} w-full rounded-xl overflow-hidden shadow-md bg-gray-100`}
              style={{ zIndex: 1 }}
              onClick={() => setMapFullscreen((f) => !f)}
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
              <button
                onClick={centerOnNext}
                className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-gray-700 hover:bg-gray-50 hover:shadow-xl transition-all"
                aria-label="Centrar en siguiente"
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
      <div className="sticky top-0 z-20 bg-white/80 backdrop-blur border-b">
        <div className="flex">
          <button
            onClick={() => setActiveTab('en-ruta')}
            className={`flex-1 py-3 px-2 text-center text-xs font-medium transition-all duration-200 border-b-2 ${
              activeTab === 'en-ruta'
                ? 'bg-gradient-to-r from-blue-50 to-indigo-50 border-indigo-500 text-indigo-700'
                : 'bg-gray-50 border-transparent text-gray-500 hover:bg-gray-100'
            }`}
          >
            <div className="flex flex-col items-center space-y-1">
              <div className="flex items-center space-x-1">
                {getTabIcon('en-ruta')}
                <span className="truncate">En ruta</span>
              </div>
              <span className={`${
                activeTab === 'en-ruta' ? 'bg-indigo-200 text-indigo-800' : 'bg-gray-200 text-gray-600'
              } px-2 py-0.5 rounded-full text-xs font-bold`}>
                ({inRouteUnits.length})
              </span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('entregados')}
            className={`flex-1 py-3 px-2 text-center text-xs font-medium transition-all duration-200 border-b-2 ${
              activeTab === 'entregados'
                ? 'bg-gradient-to-r from-blue-50 to-indigo-50 border-indigo-500 text-indigo-700'
                : 'bg-gray-50 border-transparent text-gray-500 hover:bg-gray-100'
            }`}
          >
            <div className="flex flex-col items-center space-y-1">
              <div className="flex items-center space-x-1">
                {getTabIcon('entregados')}
                <span className="truncate">Entregados</span>
              </div>
              <span className={`${
                activeTab === 'entregados' ? 'bg-indigo-200 text-indigo-800' : 'bg-gray-200 text-gray-600'
              } px-2 py-0.5 rounded-full text-xs font-bold`}>
                ({deliveredUnits.length})
              </span>
            </div>
          </button>
          <button
            onClick={() => setActiveTab('no-entregados')}
            className={`flex-1 py-3 px-2 text-center text-xs font-medium transition-all duration-200 border-b-2 ${
              activeTab === 'no-entregados'
                ? 'bg-gradient-to-r from-blue-50 to-indigo-50 border-indigo-500 text-indigo-700'
                : 'bg-gray-50 border-transparent text-gray-500 hover:bg-gray-100'
            }`}
          >
            <div className="flex flex-col items-center space-y-1">
              <div className="flex items-center space-x-1">
                {getTabIcon('no-entregados')}
                <span className="truncate">No entregados</span>
              </div>
              <span className={`${
                activeTab === 'no-entregados' ? 'bg-indigo-200 text-indigo-800' : 'bg-gray-200 text-gray-600'
              } px-2 py-0.5 rounded-full text-xs font-bold`}>
                ({notDeliveredUnits.length})
              </span>
            </div>
          </button>
        </div>
      </div>
      )}

      {viewMode === 'list' && (
      <div className="p-4 space-y-4">
        {visits.map((visit: any, visitIndex: number) => {
          const matchesForTab: number = (visit?.orders || []).reduce(
            (acc: number, order: any, orderIndex: number) => {
              const countInOrder = (order?.deliveryUnits || []).reduce(
                (a: number, _unit: any, uIdx: number) =>
                  a + (shouldRenderByTab(getDeliveryUnitStatus(visitIndex, orderIndex, uIdx)) ? 1 : 0),
                0
              )
              return acc + countInOrder
            },
            0
          )
          if (matchesForTab === 0) return null
          return (
          <div key={visitIndex} className="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden border border-gray-100 active:scale-98">
            <div className="p-4 border-b border-gray-100">
              <div className="flex items-start space-x-3">
                <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 text-white rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0">
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

              {visit.orders?.map((order: any, orderIndex: number) => (
                <div key={orderIndex} className="mb-4">
                  <div className="mb-2">
                    <span className="inline-block bg-gradient-to-r from-orange-400 to-red-500 text-white px-2 py-1 rounded-lg text-xs font-medium">
                      {order.referenceID}
                    </span>
                  </div>
                  {(order.deliveryUnits || [])
                    .map((unit: any, uIdx: number): MappedUnit => ({
                      unit,
                      uIdx,
                      status: getDeliveryUnitStatus(visitIndex, orderIndex, uIdx),
                    }))
                    .filter((x: MappedUnit) => shouldRenderByTab(x.status))
                    .map((x: MappedUnit) => (
                      // extraemos para legibilidad
                      (({ unit, uIdx, status }: MappedUnit) => (
                      <div
                        key={uIdx}
                        className={`bg-gradient-to-r from-gray-50 to-blue-50 rounded-lg p-3 border ${getStatusColor(status).replace('bg-white ', '')}`}
                      >
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
                            <button
                              onClick={() => markDeliveryUnit(visitIndex, orderIndex, uIdx, 'delivered')}
                              className={`flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors ${
                                status === 'delivered'
                                  ? 'bg-green-600 text-white'
                                  : 'bg-green-100 text-green-700 hover:bg-green-200'
                              }`}
                            >
                              <CheckCircle size={16} />
                              <span>entregado</span>
                            </button>
                            <button
                              onClick={() => markDeliveryUnit(visitIndex, orderIndex, uIdx, 'not-delivered')}
                              className={`flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors ${
                                status === 'not-delivered'
                                  ? 'bg-red-600 text-white'
                                  : 'bg-red-100 text-red-700 hover:bg-red-200'
                              }`}
                            >
                              <XCircle size={16} />
                              <span>no entregado</span>
                            </button>
                          </div>
                        )}
                      </div>
                      ))(x)
                    ))}
                </div>
              ))}
            </div>
          </div>
          )
        })}
      </div>
      )}

      {/* En modo mapa: mostrar sólo la siguiente visita debajo del mapa (si no está en pantalla completa) */}
      {viewMode === 'map' && !mapFullscreen && (() => {
        const nextIdx = getNextPendingVisitIndex()
        if (typeof nextIdx !== 'number') return null
        const visit: any = (visits as any)[nextIdx]
        return (
          <div className="p-4 space-y-4">
            <div className="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden border border-gray-100">
              <div className="p-4 border-b border-gray-100">
                <div className="flex items-start space-x-3">
                  <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 text-white rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0">
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
                        status: getDeliveryUnitStatus(nextIdx, orderIndex, uIdx),
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
                              <button
                                onClick={() => markDeliveryUnit(nextIdx, orderIndex, uIdx, 'delivered')}
                                className={`flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors ${
                                  status === 'delivered' ? 'bg-green-600 text-white' : 'bg-green-100 text-green-700 hover:bg-green-200'
                                }`}
                              >
                                <CheckCircle size={16} />
                                <span>entregado</span>
                              </button>
                              <button
                                onClick={() => markDeliveryUnit(nextIdx, orderIndex, uIdx, 'not-delivered')}
                                className={`flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors ${
                                  status === 'not-delivered' ? 'bg-red-600 text-white' : 'bg-red-100 text-red-700 hover:bg-red-200'
                                }`}
                              >
                                <XCircle size={16} />
                                <span>no entregado</span>
                              </button>
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
    </div>
  )
}