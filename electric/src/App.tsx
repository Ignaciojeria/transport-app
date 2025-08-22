import { useLiveQuery } from '@tanstack/react-db'
import { useParams } from '@tanstack/react-router'
import { createRoutesCollection } from './db/create-routes-collection'
import { useMemo, useState } from 'react'
import { CheckCircle, XCircle, Play, Package, Phone, User, MapPin } from 'lucide-react'


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
          <DeliveryRouteView routeData={raw} />
        ) : (
          <pre>{JSON.stringify(data, (_key, value) => (typeof value === 'bigint' ? value.toString() : value), 2)}</pre>
        )
      })()}
    </div>
  )
}

type DeliveryRouteRaw = {
  vehicle?: { plate?: string }
  visits?: Array<any>
}

function DeliveryRouteView({ routeData }: { routeData: DeliveryRouteRaw }) {
  const [routeStarted, setRouteStarted] = useState(false)
  const [deliveryStates, setDeliveryStates] = useState<Record<string, 'delivered' | 'not-delivered' | undefined>>({})
  const [activeTab, setActiveTab] = useState<'en-ruta' | 'entregados' | 'no-entregados'>('en-ruta')

  const handleStartRoute = () => {
    setRouteStarted(true)
  }

  const markDeliveryUnit = (
    visitIndex: number,
    orderIndex: number,
    unitIndex: number,
    status: 'delivered' | 'not-delivered'
  ) => {
    const key = `${visitIndex}-${orderIndex}-${unitIndex}`
    setDeliveryStates((prev) => ({ ...prev, [key]: status }))
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

  const visits = routeData?.visits ?? []
  
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
    <div className="min-h-screen bg-gray-50 pb-40">
      <div className="bg-blue-600 text-white p-4 shadow-lg">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold">Ruta de Entrega</h1>
            <p className="text-blue-100">Vehículo: {routeData?.vehicle?.plate ?? '—'}</p>
          </div>
          {!routeStarted ? (
            <button
              onClick={handleStartRoute}
              className="flex items-center space-x-2 bg-green-500 hover:bg-green-600 px-4 py-2 rounded-lg font-medium transition-colors"
            >
              <Play size={20} />
              <span>Iniciar Ruta</span>
            </button>
          ) : (
            <div className="flex items-center space-x-2 text-green-200">
              <CheckCircle size={20} />
              <span>Ruta Iniciada</span>
            </div>
          )}
        </div>
      </div>

      {/* Tabs sticky: En ruta | Entregados | No entregados */}
      <div className="sticky top-0 z-20 bg-white border-b">
        <div className="flex divide-x">
          <button
            onClick={() => setActiveTab('en-ruta')}
            className={`flex-1 px-4 py-3 text-sm font-medium transition-colors ${
              activeTab === 'en-ruta' ? 'text-blue-700 bg-blue-50' : 'text-gray-700 hover:bg-gray-50'
            }`}
          >
            En ruta ({inRouteUnits.length})
          </button>
          <button
            onClick={() => setActiveTab('entregados')}
            className={`flex-1 px-4 py-3 text-sm font-medium transition-colors ${
              activeTab === 'entregados' ? 'text-green-700 bg-green-50' : 'text-gray-700 hover:bg-gray-50'
            }`}
          >
            Entregados ({deliveredUnits.length})
          </button>
          <button
            onClick={() => setActiveTab('no-entregados')}
            className={`flex-1 px-4 py-3 text-sm font-medium transition-colors ${
              activeTab === 'no-entregados' ? 'text-red-700 bg-red-50' : 'text-gray-700 hover:bg-gray-50'
            }`}
          >
            No entregados ({notDeliveredUnits.length})
          </button>
        </div>
      </div>

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
          <div key={visitIndex} className="bg-white rounded-lg shadow-md overflow-hidden">
            <div className="bg-gray-100 p-4 border-b">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="bg-blue-100 p-2 rounded-full">
                    <span className="text-blue-600 font-bold">{visit.sequenceNumber}</span>
                  </div>
                  <div>
                    <h3 className="font-semibold text-lg flex items-center space-x-2">
                      <User size={18} />
                      <span>{visit.addressInfo?.contact?.fullName}</span>
                    </h3>
                    {visit.addressInfo?.contact?.phone && (
                      <p className="text-gray-600 flex items-center space-x-2">
                        <Phone size={16} />
                        <span>{visit.addressInfo.contact.phone}</span>
                      </p>
                    )}
                  </div>
                </div>
              </div>
              <div className="mt-2 flex items-center space-x-2 text-gray-600">
                <MapPin size={16} />
                <span>{visit.addressInfo?.addressLine1}</span>
              </div>
            </div>

            <div className="p-4">
              <h4 className="font-medium text-gray-800 mb-3 flex items-center space-x-2">
                <Package size={18} />
                <span>Unidades de Entrega:</span>
              </h4>

              {visit.orders?.map((order: any, orderIndex: number) => (
                <div key={orderIndex} className="mb-4">
                  <p className="text-sm text-gray-600 mb-2">Orden: {order.referenceID}</p>
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
                        className={`border rounded-lg p-3 mb-3 transition-all ${getStatusColor(status)}`}
                      >
                        <div className="space-y-2">
                          <div className="font-medium">Unidad de Entrega {uIdx + 1}</div>
                          <div className="text-sm space-y-1">
                            {unit.items?.map((item: any, itemIndex: number) => (
                              <div key={itemIndex} className="flex justify-between">
                                <span>• {item.description}</span>
                                <span className="font-medium">Cantidad: {item.quantity}</span>
                              </div>
                            ))}
                          </div>
                          <div className="text-xs text-gray-500">
                            Peso: {unit.weight}kg | Volumen: {unit.volume}m³
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

      {/* Secciones inferiores: En ruta, Entregados, No entregados */}
      <div className="p-4 space-y-6">
        <div>
          <h4 className="text-sm font-semibold text-gray-700 mb-2">En ruta ({inRouteUnits.length})</h4>
          <div className="flex flex-wrap gap-2">
            {inRouteUnits.length === 0 ? (
              <span className="text-xs text-gray-500">Sin unidades pendientes.</span>
            ) : (
              inRouteUnits.map((x, idx) => (
                <span
                  key={`inroute-${idx}`}
                  className="text-xs border border-gray-300 bg-white text-gray-700 rounded-full px-2 py-1"
                  title={`${x.visit?.addressInfo?.contact?.fullName || ''}`}
                >
                  {(x.visit?.sequenceNumber ?? '—')} • {(x.order?.referenceID ?? '—')} • U{(x.uIdx + 1)}
                </span>
              ))
            )}
          </div>
        </div>

        <div>
          <h4 className="text-sm font-semibold text-green-700 mb-2">Entregados ({deliveredUnits.length})</h4>
          <div className="flex flex-wrap gap-2">
            {deliveredUnits.length === 0 ? (
              <span className="text-xs text-green-700/60">Sin unidades entregadas.</span>
            ) : (
              deliveredUnits.map((x, idx) => (
                <span
                  key={`deliv-${idx}`}
                  className="text-xs border border-green-200 bg-green-50 text-green-700 rounded-full px-2 py-1"
                  title={`${x.visit?.addressInfo?.contact?.fullName || ''}`}
                >
                  {(x.visit?.sequenceNumber ?? '—')} • {(x.order?.referenceID ?? '—')} • U{(x.uIdx + 1)}
                </span>
              ))
            )}
          </div>
        </div>

        <div className="mb-8">
          <h4 className="text-sm font-semibold text-red-700 mb-2">No entregados ({notDeliveredUnits.length})</h4>
          <div className="flex flex-wrap gap-2">
            {notDeliveredUnits.length === 0 ? (
              <span className="text-xs text-red-700/60">Sin unidades no entregadas.</span>
            ) : (
              notDeliveredUnits.map((x, idx) => (
                <span
                  key={`notdeliv-${idx}`}
                  className="text-xs border border-red-200 bg-red-50 text-red-700 rounded-full px-2 py-1"
                  title={`${x.visit?.addressInfo?.contact?.fullName || ''}`}
                >
                  {(x.visit?.sequenceNumber ?? '—')} • {(x.order?.referenceID ?? '—')} • U{(x.uIdx + 1)}
                </span>
              ))
            )}
          </div>
        </div>
      </div>

      {routeStarted && (
        <div className="fixed bottom-0 left-0 right-0 bg-white border-t p-4 shadow-lg">
          <div className="flex justify-between items-center">
            <div className="text-sm text-gray-600">Progreso de la ruta</div>
            <div className="flex space-x-4 text-sm">
              <span className="text-green-600 font-medium">
                ✓ {Object.values(deliveryStates).filter((s) => s === 'delivered').length} Entregadas
              </span>
              <span className="text-red-600 font-medium">
                ✗ {Object.values(deliveryStates).filter((s) => s === 'not-delivered').length} No Entregadas
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}