import { VisitCardDeliveryUnit } from './VisitCardDeliveryUnit'
import { groupDeliveryUnitsByLocation, hasGroupPendingUnits, getGroupPendingUnits, type DeliveryGroup } from './GroupedDeliveryUtils'
import { Package, Package2, Users, User } from 'lucide-react'
import { IdentifierBadge } from './IdentifierBadge'
import { OrderCard, DeliveryUnitCard } from './OrderCard'

interface VisitCardOrdersProps {
  visit: any
  visitIndex: number
  routeStarted: boolean
  onOpenDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenNonDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenGroupedDelivery?: (visitIndex: number, group: DeliveryGroup) => void
  onOpenGroupedNonDelivery?: (visitIndex: number, group: DeliveryGroup) => void
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, uIdx: number) => 'delivered' | 'not-delivered' | undefined
  shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean
  viewMode?: 'list' | 'map' // Nuevo prop para controlar agrupación
}

export function VisitCardOrders({
  visit,
  visitIndex,
  routeStarted,
  onOpenDelivery,
  onOpenNonDelivery,
  onOpenGroupedDelivery,
  onOpenGroupedNonDelivery,
  getDeliveryUnitStatus,
  shouldRenderByTab,
  viewMode = 'list'
}: VisitCardOrdersProps) {
  // Solo agrupar en modo mapa, no en modo lista
  const shouldGroup = viewMode === 'map'
  
  // Obtener grupos de delivery units agrupables solo si estamos en modo mapa
  const deliveryGroups = shouldGroup ? groupDeliveryUnitsByLocation(visit, visitIndex, getDeliveryUnitStatus) : []
  
  // Filtrar grupos que tienen unidades pendientes para el tab actual
  const relevantGroups = deliveryGroups.filter(group => {
    if (shouldRenderByTab(undefined)) { // Tab "en-ruta"
      return hasGroupPendingUnits(group)
    }
    return true // Para otros tabs, mostrar todos los grupos
  })

  return (
    <>
      {/* Mostrar grupos agrupados primero */}
      {relevantGroups.map((group) => {
        const pendingUnits = getGroupPendingUnits(group)
        const hasPendingForTab = pendingUnits.length > 0 && shouldRenderByTab(undefined)
        
        if (!hasPendingForTab && shouldRenderByTab(undefined)) return null
        
        return (
          <div key={group.key} className="mb-4">
            {/* Header del grupo - estilo similar al VisitCardHeader */}
            <div className="bg-white rounded-xl shadow-md border border-gray-100 mb-2">
              <div className="p-4">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center space-x-3">
                    <div className="bg-purple-100 p-2 rounded-lg">
                      <span className="text-purple-600 font-bold text-lg">1</span>
                    </div>
                    <div>
                      <h3 className="text-lg font-bold text-gray-900 flex items-center">
                        <Users className="w-5 h-5 mr-2 text-purple-600" />
                        {group.units?.[0]?.order?.contact?.fullName || 'Cliente'}
                      </h3>
                      <p className="text-sm text-gray-600 flex items-center">
                        <span className="text-gray-500 mr-1">📍</span>
                        <span className="line-clamp-2">{group.addressInfo.addressLine1}</span>
                      </p>
                    </div>
                  </div>
                  
                  {/* Botón de acción grupal */}
                  {routeStarted && hasPendingForTab && onOpenGroupedDelivery && onOpenGroupedNonDelivery && (
                    <div className="flex space-x-2">
                      <button
                        onClick={() => onOpenGroupedDelivery(visitIndex, group)}
                        className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg font-medium flex items-center space-x-2 transition-colors"
                      >
                        <Package className="w-4 h-4" />
                        <span>Entregar todo</span>
                      </button>
                      <button
                        onClick={() => onOpenGroupedNonDelivery(visitIndex, group)}
                        className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg font-medium flex items-center space-x-2 transition-colors"
                      >
                        <Package className="w-4 h-4" />
                        <span>No entregar</span>
                      </button>
                    </div>
                  )}
                </div>
                
                <div className="text-xs text-purple-600 bg-purple-50 px-2 py-1 rounded-full inline-block">
                  {group.totalUnits} unidades • {group.pendingUnits} pendientes
                </div>
              </div>
            </div>
            
            {/* Agrupar unidades por referenceID dentro del grupo */}
            <div className="space-y-3">
              {(() => {
                // Agrupar unidades por referenceID
                const unitsByReference = group.units
                  .filter(unit => shouldRenderByTab(unit.status))
                  .reduce((acc, unit) => {
                    const refId = unit.order?.referenceID || 'Sin referencia'
                    if (!acc[refId]) {
                      acc[refId] = []
                    }
                    acc[refId].push(unit)
                    return acc
                  }, {} as Record<string, any[]>)

                return Object.entries(unitsByReference).map(([referenceID, units]) => (
                  <div key={referenceID} className="space-y-2">
                    {/* Header de la orden */}
                    <div className="bg-gradient-to-r from-gray-50 to-blue-50 px-3 py-2 rounded-lg border border-gray-200">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <Package2 className="w-4 h-4 text-gray-600" />
                          <span className="font-medium text-gray-800">
                            Orden: {referenceID}
                          </span>
                        </div>
                        <span className="text-sm text-gray-500">
                          {units.length} unidad{units.length !== 1 ? 'es' : ''}
                        </span>
                      </div>
                    </div>

                    {/* Unidades de esta orden */}
                    <div className="space-y-2 ml-4">
                      {units.map((unit) => (
                        <DeliveryUnitCard
                          key={`${unit.orderIndex}-${unit.uIdx}`}
                          unit={unit.unit}
                          unitIndex={unit.uIdx}
                          visitIndex={visitIndex}
                          orderIndex={unit.orderIndex}
                          status={unit.status}
                          orderReferenceID={unit.order?.referenceID}
                        >
                          <div className="space-y-1 mb-3">
                            {unit.unit.items?.map((item: any, itemIndex: number) => (
                              <div key={itemIndex} className="text-sm text-gray-600">
                                • {item.description}
                              </div>
                            ))}
                            <div className="text-sm text-gray-600">
                              • {unit.unit.weight || 0}kg
                            </div>
                            <div className="text-sm text-gray-600">
                              • {unit.unit.volume || 0}m³
                            </div>
                          </div>
                          
                          <div className="flex items-center justify-between">
                            <div className="text-xs text-gray-500">
                              Cant.
                            </div>
                            <div className="text-lg font-bold text-purple-600">
                              {unit.unit.items?.reduce((sum: number, item: any) => sum + (item.quantity || 0), 0) || 0}
                            </div>
                          </div>
                          
                          {/* Botones individuales */}
                          {!unit.status && routeStarted && (
                            <div className="flex space-x-2 mt-3">
                              <button
                                onClick={() => onOpenDelivery(visitIndex, unit.orderIndex, unit.uIdx)}
                                className="flex-1 bg-green-100 hover:bg-green-200 text-gray-700 px-4 py-2 rounded-lg font-medium flex items-center justify-center space-x-2 transition-colors"
                              >
                                <span className="text-green-600">✓</span>
                                <span>Entregar</span>
                              </button>
                              <button
                                onClick={() => onOpenNonDelivery(visitIndex, unit.orderIndex, unit.uIdx)}
                                className="flex-1 bg-red-100 hover:bg-red-200 text-gray-700 px-4 py-2 rounded-lg font-medium flex items-center justify-center space-x-2 transition-colors"
                              >
                                <span className="text-red-600">✗</span>
                                <span>No entregado</span>
                              </button>
                            </div>
                          )}
                        </DeliveryUnitCard>
                      ))}
                    </div>
                  </div>
                ))
              })()}
            </div>
          </div>
        )
      })}
      
      {/* Mostrar unidades no agrupadas (que no pertenecen a ningún grupo) */}
      {(() => {
        // Procesar órdenes y agrupar por cliente
        const ordersByClient = new Map()
        
        visit.orders?.forEach((order: any, orderIndex: number) => {
          const orderUnits = (order.deliveryUnits || [])
            .map((unit: any, uIdx: number) => ({
              unit,
              uIdx,
              status: getDeliveryUnitStatus(visitIndex, orderIndex, uIdx),
            }))
            .filter((x: any) => shouldRenderByTab(x.status))
          
          // Verificar si esta unidad ya está en algún grupo (solo relevante en modo mapa)
          const isInGroup = shouldGroup && deliveryGroups.some(group => 
            group.units.some(groupUnit => 
              groupUnit.visitIndex === visitIndex && 
              groupUnit.orderIndex === orderIndex && 
              groupUnit.uIdx === orderUnits.findIndex((u: any) => u.uIdx === groupUnit.uIdx)
            )
          )
          
          if (isInGroup || orderUnits.length === 0) return
          
          const clientName = order.contact?.fullName || 'Sin nombre'
          
          if (!ordersByClient.has(clientName)) {
            ordersByClient.set(clientName, [])
          }
          
          ordersByClient.get(clientName).push({
            order,
            orderIndex,
            orderUnits
          })
        })
        
        // Convertir a array y ordenar por nombre de cliente
        const sortedClients = Array.from(ordersByClient.entries())
          .sort(([clientA], [clientB]) => clientA.localeCompare(clientB))
        
        return sortedClients.map(([clientName, clientOrders]) => (
          <div key={clientName} className="mb-6">
            {/* Header del cliente */}
            <div className="mb-4 p-4 bg-gradient-to-r from-indigo-50 to-blue-50 rounded-lg border border-indigo-200">
              <div className="flex items-center space-x-3">
                <User className="w-5 h-5 text-indigo-600" />
                <div className="flex-1">
                  <h3 className="text-sm font-bold text-gray-800">
                    {clientName}
                  </h3>
                  <span className="text-xs text-gray-600">
                    {clientOrders.length} {clientOrders.length === 1 ? 'orden' : 'órdenes'}
                  </span>
                </div>
                {clientOrders[0]?.order?.contact?.phone && (
                  <span className="text-xs text-gray-500 bg-white px-3 py-1 rounded-full border">
                    📞 {clientOrders[0].order.contact.phone}
                  </span>
                )}
              </div>
            </div>
            
            {/* Órdenes del cliente */}
            <div className="ml-4 space-y-4">
              {clientOrders.map(({ order, orderIndex, orderUnits }: any) => (
                <div key={orderIndex} className="border-l-2 border-indigo-200 pl-4">
                  {/* Información de la orden */}
                  <div className="mb-3 p-3 bg-white rounded-lg border border-gray-200 shadow-sm">
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center space-x-2">
                        <Package2 className="w-4 h-4 text-gray-600" />
                        <span className="text-sm font-medium text-gray-700">
                          Orden: {order.referenceID || `#${orderIndex + 1}`}
                        </span>
                      </div>
                      <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                        {orderUnits.length} {orderUnits.length === 1 ? 'unidad' : 'unidades'}
                      </span>
                    </div>
                    {order.instructions && (
                      <div className="text-xs text-gray-600 bg-blue-50 p-2 rounded border-l-2 border-blue-200">
                        <strong>Instrucciones:</strong> {order.instructions}
                      </div>
                    )}
                  </div>
                  
                  {/* Unidades de entrega */}
                  {orderUnits.map((x: any) => (
                    <div key={x.uIdx} className="mb-3">
                      <div className="mb-2">
                        <IdentifierBadge 
                          lpn={x.unit.lpn} 
                          code={x.unit.code} 
                          size="sm"
                        />
                      </div>
                      <VisitCardDeliveryUnit
                        unit={x.unit}
                        uIdx={x.uIdx}
                        status={x.status}
                        visitIndex={visitIndex}
                        orderIndex={orderIndex}
                        routeStarted={routeStarted}
                        orderReferenceID={order.referenceID}
                        onOpenDelivery={onOpenDelivery}
                        onOpenNonDelivery={onOpenNonDelivery}
                      />
                    </div>
                  ))}
                </div>
              ))}
            </div>
          </div>
        ))
      })()}
    </>
  )
}
