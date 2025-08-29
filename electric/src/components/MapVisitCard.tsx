import { CheckCircle, XCircle, Play, Package, User, MapPin } from 'lucide-react'

interface MapVisitCardProps {
  visit: any
  displayIdx: number
  isSelectedVisit: boolean
  isProcessed: boolean
  hasNextPending: boolean
  nextPendingIdx: number | null
  routeStarted: boolean
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined
  openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void
  openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void
  onNextPending: (nextPendingIdx: number) => void
  onClearSelection: () => void
}

export function MapVisitCard({
  visit,
  displayIdx,
  isSelectedVisit,
  isProcessed,
  hasNextPending,
  nextPendingIdx,
  routeStarted,
  getDeliveryUnitStatus,
  openDeliveryFor,
  openNonDeliveryFor,
  onNextPending,
  onClearSelection
}: MapVisitCardProps) {
  
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
            onClick={() => nextPendingIdx !== null && onNextPending(nextPendingIdx)}
            className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white py-3 px-4 rounded-lg font-medium flex items-center justify-center space-x-2 transition-all duration-200 shadow-md hover:shadow-lg active:scale-95"
          >
            <Play className="w-4 h-4" />
            <span>Siguiente a Entregar (#{visit?.sequenceNumber})</span>
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
            onClick={onClearSelection}
            className="text-xs text-blue-600 hover:text-blue-800 font-medium"
          >
            Ver siguiente
          </button>
        )}
      </div>
      
      {/* Card de la visita */}
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
}
