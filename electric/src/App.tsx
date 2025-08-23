/* eslint-disable @typescript-eslint/no-explicit-any */
import { useLiveQuery } from '@tanstack/react-db'
import { useParams } from '@tanstack/react-router'
import { createRoutesCollection } from './db/create-routes-collection'
import { useMemo, useState } from 'react'
import { CheckCircle, XCircle, Play, Package, User, MapPin } from 'lucide-react'


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

      {/* Tabs sticky: En ruta | Entregados | No entregados */}
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

      {/* Resumen de la Ruta */}
      <div className="p-4">
        <div className="mt-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-xl p-4 shadow-lg">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="text-sm font-bold mb-1">Resumen de la Ruta</h3>
              <p className="text-indigo-100 text-xs">Total entregas programadas</p>
            </div>
            <div className="text-right">
              <div className="text-2xl font-bold">{allUnits.length}</div>
              <div className="text-indigo-200 text-xs">entregas</div>
            </div>
          </div>
        </div>
      </div>

      {/* Barra inferior de progreso eliminada por redundancia con la barra superior */}
    </div>
  )
}