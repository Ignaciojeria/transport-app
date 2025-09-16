import { CheckCircle, XCircle, Play, Package, User, MapPin, Users } from 'lucide-react'
import { IdentifierBadge } from './IdentifierBadge'
import { useLanguage } from '../hooks/useLanguage'

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
  onDeliverAll?: (visitIndex: number) => void
  onNonDeliverAll?: (visitIndex: number) => void
  openGroupedDelivery?: (visitIndex: number, group: any) => void
  openGroupedNonDelivery?: (visitIndex: number, group: any) => void
  selectedClient?: any
  hasMultipleClients?: boolean
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
  onClearSelection,
  onDeliverAll,
  onNonDeliverAll,
  openGroupedDelivery,
  openGroupedNonDelivery,
  selectedClient,
  hasMultipleClients
}: MapVisitCardProps) {
  const { t } = useLanguage()
  
  // Helper para verificar si un cliente tiene entregas parciales
  const hasClientPartialDeliveries = (clientName: string) => {
    const clientOrders = (visit.orders || []).filter((order: any) => 
      order.contact?.fullName === clientName
    )
    let totalClientUnits = 0
    let deliveredClientUnits = 0
    
    clientOrders.forEach((order: any) => {
      const orderIndex = (visit.orders || []).indexOf(order)
      ;(order.deliveryUnits || []).forEach((_unit: any, unitIndex: number) => {
        totalClientUnits++
        const status = getDeliveryUnitStatus(displayIdx, orderIndex, unitIndex)
        if (status === 'delivered') {
          deliveredClientUnits++
        }
      })
    })
    
    return deliveredClientUnits > 0 && deliveredClientUnits < totalClientUnits
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

  // Calcular estadísticas de la visita (filtradas por cliente seleccionado si aplica)
  const visitStats = (() => {
    let totalUnits = 0
    let pendingUnits = 0
    let deliveredUnits = 0
    let notDeliveredUnits = 0

    // Filtrar órdenes según el cliente seleccionado
    const ordersToCalculate = hasMultipleClients && selectedClient 
      ? (visit.orders || []).filter((order: any) => 
          order.contact?.fullName === selectedClient.clientName
        )
      : (visit.orders || [])

    ordersToCalculate.forEach((order: any) => {
      // Encontrar el índice real de la orden en la visita original
      const realOrderIndex = (visit.orders || []).findIndex((o: any) => o === order)
      
      order.deliveryUnits?.forEach((_unit: any, unitIndex: number) => {
        totalUnits++
        const status = getDeliveryUnitStatus(displayIdx, realOrderIndex, unitIndex)
        if (status === 'delivered') {
          deliveredUnits++
        } else if (status === 'not-delivered') {
          notDeliveredUnits++
        } else {
          pendingUnits++
        }
      })
    })

    return {
      totalUnits,
      pendingUnits,
      deliveredUnits,
      notDeliveredUnits,
      hasPendingUnits: pendingUnits > 0,
      hasDeliveredUnits: deliveredUnits > 0,
      isPartiallyDelivered: deliveredUnits > 0 && pendingUnits > 0
    }
  })()

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
            <span>{t.delivery.next} (#{visit?.sequenceNumber})</span>
          </button>
        </div>
      )}
      
      {/* Indicador de qué visita se está mostrando */}
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-gray-700">
          {isSelectedVisit ? t.visitCard.selectedVisit : t.visitCard.nextToDeliver}
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
              {hasMultipleClients && selectedClient ? (
                // Mostrar cliente seleccionado cuando hay múltiples clientes
                <div className="mb-1">
                  <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                    <User className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
                    <span className="truncate text-purple-700">{selectedClient.clientName}</span>
                  </h3>
                  <div className="text-xs text-gray-600 mb-1">
                    <span className="inline-block bg-purple-100 text-purple-700 px-2 py-1 rounded text-xs font-medium">
{t.visitCard.selectedClient}
                    </span>
                    {selectedClient.orderIndexes && selectedClient.orderIndexes.length > 1 && (
                      <span className="inline-block bg-gray-100 text-gray-600 px-2 py-1 rounded ml-1 text-xs">
                        {selectedClient.orderIndexes.length} {t.visitCard.orders}
                      </span>
                    )}
                  </div>
                </div>
              ) : hasMultipleClients ? (
                // Mostrar múltiples clientes cuando no hay selección específica
                <div className="mb-1">
                  <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                    <Users className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
                    <span className="text-indigo-700">{t.nextVisit.multipleClients}</span>
                  </h3>
                  <div className="text-xs text-gray-600 mb-1">
                    <span className="inline-block bg-indigo-100 text-indigo-700 px-2 py-1 rounded text-xs">
{t.nextVisit.selectClient}
                    </span>
                  </div>
                </div>
              ) : (
                // Mostrar cliente único
                <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                  <User className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
                  <span className="truncate">{visit.orders?.[0]?.contact?.fullName || 'Sin nombre'}</span>
                </h3>
              )}
              <p className="text-xs text-gray-600 flex items-start mb-2">
                <MapPin className="w-3 h-3 mr-1 mt-0.5 text-gray-500 flex-shrink-0" />
                <span className="line-clamp-2">{visit.addressInfo?.addressLine1}</span>
              </p>
            </div>
          </div>
        </div>
        <div className="p-4">
          {/* Botones de acción grupal - solo si hay 2 o más unidades pendientes */}
          {routeStarted && visitStats.hasPendingUnits && visitStats.pendingUnits >= 2 && onDeliverAll && onNonDeliverAll && (
            <div className="mb-4 p-3 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg border border-blue-200">
              <div className="flex items-center justify-between mb-2">
                <div>
                  <span className="text-sm font-medium text-gray-700">
                    {visitStats.isPartiallyDelivered ? t.delivery.actionsForRemaining : t.delivery.groupActions}
                  </span>
                  {visitStats.isPartiallyDelivered && (
                    <div className="text-xs text-gray-500 mt-1">
                      {visitStats.deliveredUnits} {t.status.delivered} • {visitStats.pendingUnits} {t.status.pendingUnits}
                    </div>
                  )}
                </div>
                <Package className="w-4 h-4 text-blue-600" />
              </div>
              <div className="flex space-x-2">
                <button
                  onClick={() => {
                    // Si hay múltiples clientes y uno seleccionado, crear grupo solo para ese cliente
                    if (hasMultipleClients && selectedClient) {
                      // Crear un grupo temporal solo con las unidades del cliente seleccionado
                      const clientUnits: any[] = []
                      
                      ;(visit.orders || [])
                        .filter((order: any) => order.contact?.fullName === selectedClient.clientName)
                        .forEach((order: any) => {
                          const orderIndex = (visit.orders || []).indexOf(order)
                          ;(order.deliveryUnits || []).forEach((unit: any, unitIndex: number) => {
                            const status = getDeliveryUnitStatus(displayIdx, orderIndex, unitIndex)
                            if (!status) { // Solo unidades pendientes
                              clientUnits.push({
                                unit,
                                uIdx: unitIndex,
                                orderIndex,
                                order
                              })
                            }
                          })
                        })
                      
                      if (clientUnits.length > 0) {
                        // Crear grupo temporal para el cliente seleccionado
                        const clientGroup = {
                          key: `client-${displayIdx}-${selectedClient.clientName}`,
                          addressInfo: visit.addressInfo,
                          units: clientUnits.map((unit) => ({
                            unit: unit.unit,
                            uIdx: unit.uIdx,
                            status: undefined,
                            visitIndex: displayIdx,
                            orderIndex: unit.orderIndex,
                            order: unit.order
                          })),
                          totalUnits: clientUnits.length,
                          pendingUnits: clientUnits.length
                        }
                        
                        // Usar openGroupedDelivery en lugar de onDeliverAll
                        if (openGroupedDelivery) {
                          openGroupedDelivery(displayIdx, clientGroup)
                        }
                      }
                    } else {
                      // Comportamiento original para visitas sin múltiples clientes
                      onDeliverAll && onDeliverAll(displayIdx)
                    }
                  }}
                  className="flex-1 bg-green-500 hover:bg-green-600 text-white px-3 py-2 rounded-lg font-medium flex items-center justify-center space-x-2 transition-colors text-sm"
                >
                  <CheckCircle className="w-4 h-4" />
                  <span>
                    {hasMultipleClients && selectedClient 
                      ? (hasClientPartialDeliveries(selectedClient.clientName) ? t.delivery.deliverRemaining : t.delivery.deliverAll)
                      : (visitStats.isPartiallyDelivered ? t.delivery.deliverRemaining : t.delivery.deliverAll)
                    }
                  </span>
                </button>
                <button
                  onClick={() => {
                    // Si hay múltiples clientes y uno seleccionado, crear grupo solo para ese cliente
                    if (hasMultipleClients && selectedClient) {
                      // Crear un grupo temporal solo con las unidades del cliente seleccionado
                      const clientUnits: any[] = []
                      
                      ;(visit.orders || [])
                        .filter((order: any) => order.contact?.fullName === selectedClient.clientName)
                        .forEach((order: any) => {
                          const orderIndex = (visit.orders || []).indexOf(order)
                          ;(order.deliveryUnits || []).forEach((unit: any, unitIndex: number) => {
                            const status = getDeliveryUnitStatus(displayIdx, orderIndex, unitIndex)
                            if (!status) { // Solo unidades pendientes
                              clientUnits.push({
                                unit,
                                uIdx: unitIndex,
                                orderIndex,
                                order
                              })
                            }
                          })
                        })
                      
                      if (clientUnits.length > 0) {
                        // Crear grupo temporal para el cliente seleccionado
                        const clientGroup = {
                          key: `client-nd-${displayIdx}-${selectedClient.clientName}`,
                          addressInfo: visit.addressInfo,
                          units: clientUnits.map((unit) => ({
                            unit: unit.unit,
                            uIdx: unit.uIdx,
                            status: undefined,
                            visitIndex: displayIdx,
                            orderIndex: unit.orderIndex,
                            order: unit.order
                          })),
                          totalUnits: clientUnits.length,
                          pendingUnits: clientUnits.length
                        }
                        
                        // Usar openGroupedNonDelivery en lugar de onNonDeliverAll
                        if (openGroupedNonDelivery) {
                          openGroupedNonDelivery(displayIdx, clientGroup)
                        }
                      }
                    } else {
                      // Comportamiento original para visitas sin múltiples clientes
                      onNonDeliverAll && onNonDeliverAll(displayIdx)
                    }
                  }}
                  className="flex-1 bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded-lg font-medium flex items-center justify-center space-x-2 transition-colors text-sm"
                >
                  <XCircle className="w-4 h-4" />
                  <span>
                    {hasMultipleClients && selectedClient 
                      ? (hasClientPartialDeliveries(selectedClient.clientName) ? t.delivery.notDeliverRemaining : t.delivery.notDeliverAll)
                      : (visitStats.isPartiallyDelivered ? t.delivery.notDeliverRemaining : t.delivery.notDeliverAll)
                    }
                  </span>
                </button>
              </div>
            </div>
          )}

          {/* Progreso de la visita - solo si hay múltiples unidades */}
          {visitStats.isPartiallyDelivered && visitStats.totalUnits > 1 && (
            <div className="mb-4 p-3 bg-gradient-to-r from-green-50 to-blue-50 rounded-lg border border-green-200">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-700">Progreso de la visita:</span>
                <div className="flex items-center space-x-2 text-xs text-gray-500">
                  <span className="flex items-center">
                    <div className="w-2 h-2 bg-green-500 rounded-full mr-1"></div>
                    {visitStats.deliveredUnits} entregadas
                  </span>
                  <span className="flex items-center">
                    <div className="w-2 h-2 bg-red-500 rounded-full mr-1"></div>
                    {visitStats.notDeliveredUnits} no entregadas
                  </span>
                  <span className="flex items-center">
                    <div className="w-2 h-2 bg-gray-400 rounded-full mr-1"></div>
                    {visitStats.pendingUnits} pendientes
                  </span>
                </div>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div 
                  className="bg-gradient-to-r from-green-500 to-blue-500 h-2 rounded-full transition-all duration-300"
                  style={{ 
                    width: `${((visitStats.deliveredUnits + visitStats.notDeliveredUnits) / visitStats.totalUnits) * 100}%` 
                  }}
                ></div>
              </div>
            </div>
          )}

          {/* Indicación para una sola unidad */}
          {routeStarted && visitStats.hasPendingUnits && visitStats.pendingUnits === 1 && (
            <div className="mb-4 p-3 bg-gradient-to-r from-yellow-50 to-orange-50 rounded-lg border border-yellow-200">
              <div className="flex items-center space-x-2">
                <Package className="w-4 h-4 text-orange-600" />
                <span className="text-sm text-gray-700">
                  Solo queda 1 unidad pendiente. Usa los botones individuales abajo.
                </span>
              </div>
            </div>
          )}

          {/* Solo mostrar órdenes si hay un cliente seleccionado o no hay múltiples clientes */}
          {(!hasMultipleClients || selectedClient) && (
            <>
              <h4 className="text-sm font-medium text-gray-800 mb-3 flex items-center">
                <Package size={18} />
                <span className="ml-2">{t.visitCard.deliveryUnits}:</span>
              </h4>
              {(() => {
                // Filtrar órdenes según el cliente seleccionado
                const ordersToShow = hasMultipleClients && selectedClient 
                  ? (visit.orders || []).filter((order: any) => 
                      order.contact?.fullName === selectedClient.clientName
                    )
                  : (visit.orders || [])
                
                return ordersToShow.map((order: any, orderIndex: number) => {
              // Encontrar el índice real de la orden en la visita original
              const realOrderIndex = (visit.orders || []).findIndex((o: any) => o === order)
              
              return (
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
                  status: getDeliveryUnitStatus(displayIdx, realOrderIndex, uIdx),
                }))
                .map(({ unit, uIdx, status }: { unit: any; uIdx: number; status: 'delivered' | 'not-delivered' | undefined }) => (
                  <div key={uIdx} className={`bg-gradient-to-r from-gray-50 to-blue-50 rounded-lg p-3 border ${getStatusColor(status).replace('bg-white ', '')}`}>
                    <div className="flex justify-between items-start mb-2">
                      <div className="flex-1 min-w-0">
                        {/* Identificadores prominentes */}
                        <div className="mb-2">
                          <IdentifierBadge 
                            lpn={unit.lpn} 
                            code={unit.code} 
                            size="sm"
                            className="mb-2"
                          />
                        </div>
                         <h5 className="text-sm font-medium text-gray-800 mb-2 truncate flex items-center">
                           <Package className="w-4 h-4 mr-2 text-gray-600" />
                           {t.visitCard.deliveryUnit} {uIdx + 1}
                         </h5>
                        {Array.isArray(unit.items) && unit.items.length > 0 && (
                          <div className="space-y-1 mb-2">
                            {unit.items.map((item: any, index: number) => (
                              <div key={index} className="flex items-center space-x-1">
                                <span className="w-1.5 h-1.5 bg-indigo-500 rounded-full"></span>
                                <span className="text-xs text-gray-700 truncate">{item.description}{item.quantity ? ` (${item.quantity})` : ''}</span>
                              </div>
                            ))}
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
                            <span>{t.delivery.changeToNotDelivered}</span>
                          </button>
                        ) : status === 'not-delivered' ? (
                          // Si está no entregado, mostrar solo opción de cambiar a entregado
                          <button
                            onClick={() => openDeliveryFor(displayIdx, orderIndex, uIdx)}
                            className="w-full flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-green-100 text-green-700 hover:bg-green-200"
                          >
                            <CheckCircle size={16} />
                            <span>{t.delivery.changeToDelivered}</span>
                          </button>
                        ) : (
                          // Si está pendiente, mostrar ambas opciones originales
                          <>
                            <button
                              onClick={() => openDeliveryFor(displayIdx, realOrderIndex, uIdx)}
                              className="flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-green-100 text-green-700 hover:bg-green-200"
                            >
                              <CheckCircle size={16} />
                              <span>{t.delivery.deliver}</span>
                            </button>
                            <button
                              onClick={() => openNonDeliveryFor(displayIdx, realOrderIndex, uIdx)}
                              className="flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-red-100 text-red-700 hover:bg-red-200"
                            >
                              <XCircle size={16} />
                              <span>{t.delivery.notDeliver}</span>
                            </button>
                          </>
                        )}
                      </div>
                    )}
                  </div>
                ))}
            </div>
              )
            })
          })()}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
