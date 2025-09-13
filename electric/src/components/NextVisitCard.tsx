import { Play, User, MapPin, Users, ChevronDown, ChevronUp, Package } from 'lucide-react'
import { useState } from 'react'
import { IdentifierBadge } from './IdentifierBadge'

interface NextVisitCardProps {
  nextVisit: any
  nextIdx: number
  onCenterOnVisit: (visitIndex: number) => void
  addressGroups?: { [address: string]: { clients: string[], totalUnits: number, pendingUnits: number, visitIndex: number } }
  viewMode?: 'list' | 'map'
  allVisits?: any[]
}

export function NextVisitCard({ nextVisit, nextIdx, onCenterOnVisit, addressGroups, viewMode = 'list', allVisits = [] }: NextVisitCardProps) {
  const [isExpanded, setIsExpanded] = useState(false)
  const address = nextVisit.addressInfo?.addressLine1 || 'Sin dirección'
  const addressGroup = addressGroups?.[address]
  
  // Detectar múltiples clientes dentro de la visita actual
  const uniqueClientsInVisit = Array.from(new Set(
    (nextVisit.orders || []).map((order: any) => order.contact?.fullName).filter(Boolean)
  ))
  const hasMultipleClients = uniqueClientsInVisit.length > 1
  
  // Para compatibilidad con el acordeón expandible, usar addressGroup si existe
  const clientCount = hasMultipleClients ? uniqueClientsInVisit.length : 1
  
  // Debug log
  console.log('🔍 NextVisitCard DEBUG:', {
    sequenceNumber: nextVisit.sequenceNumber,
    address,
    uniqueClientsInVisit,
    hasMultipleClients,
    clientCount,
    ordersCount: nextVisit.orders?.length || 0
  })
  
  // Obtener todas las visitas que comparten la misma dirección
  const visitsAtSameAddress = allVisits.filter(visit => 
    visit.addressInfo?.addressLine1 === address
  )

  return (
    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl border-2 border-blue-200 p-4 mb-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-bold text-blue-800 flex items-center">
          <Play className="w-4 h-4 mr-2" />
          Siguiente visita
        </h3>
        <div className="flex items-center space-x-2">
          <span className="text-xs text-blue-600 bg-blue-100 px-2 py-1 rounded-full font-medium">
            #{nextVisit.sequenceNumber}
          </span>
          {viewMode === 'map' && hasMultipleClients && (
            <button
              onClick={() => setIsExpanded(!isExpanded)}
              className="w-6 h-6 bg-blue-100 hover:bg-blue-200 text-blue-600 rounded-full flex items-center justify-center transition-all duration-200"
              aria-label={isExpanded ? "Contraer detalles" : "Expandir detalles"}
            >
              {isExpanded ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
            </button>
          )}
        </div>
      </div>
      
      <div className="flex items-start space-x-3">
        <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-indigo-600 text-white rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0">
          {nextVisit.sequenceNumber}
        </div>
        <div className="flex-1 min-w-0">
          {/* Siempre mostrar la dirección como elemento principal */}
          <h4 className="text-sm font-bold text-gray-800 flex items-center mb-1">
            <MapPin className="w-3 h-3 mr-1 text-gray-500 flex-shrink-0" />
            <span className="line-clamp-2">{address}</span>
          </h4>
          
          {/* Mostrar información de clientes de forma discreta */}
          {hasMultipleClients ? (
            <div className="mb-2">
              <p className="text-xs text-gray-600 mb-1">
                <span className="text-blue-700 font-medium">{clientCount} clientes</span>
              </p>
              {addressGroup && (
                <div className="text-xs text-blue-600 bg-blue-100 px-2 py-1 rounded-full inline-block">
                  {addressGroup.totalUnits} unidades • {addressGroup.pendingUnits} pendientes
                </div>
              )}
            </div>
          ) : (
            <p className="text-xs text-gray-600">
              <span className="text-blue-700 font-medium">1 cliente</span>
            </p>
          )}
        </div>
        <button
          onClick={() => onCenterOnVisit(nextIdx)}
          className="w-8 h-8 bg-blue-50 hover:bg-blue-100 border border-blue-200 text-blue-600 rounded-lg flex items-center justify-center transition-all duration-200 hover:shadow-md active:scale-95 flex-shrink-0"
          aria-label={`Ver en mapa - Visita ${nextVisit.sequenceNumber}`}
          title="Ver en mapa"
        >
          <MapPin className="w-4 h-4" />
        </button>
      </div>

      {/* Acordeón expandible para modo mapa */}
      {viewMode === 'map' && hasMultipleClients && isExpanded && (
        <div className="mt-4 pt-4 border-t border-blue-200">
          <div className="space-y-3">
            {visitsAtSameAddress.map((visit, visitIndex) => (
              <div key={visitIndex} className="bg-white rounded-lg p-3 border border-gray-200">
                <div className="flex items-center justify-between mb-2">
                  <h5 className="text-sm font-semibold text-gray-800 flex items-center">
                    <User className="w-3 h-3 mr-1 text-gray-600" />
                    {visit.addressInfo?.contact?.fullName}
                  </h5>
                  <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                    #{visit.sequenceNumber}
                  </span>
                </div>
                
                <div className="space-y-2">
                  {visit.orders?.map((order: any, orderIndex: number) => (
                    <div key={orderIndex} className="bg-gray-50 rounded p-2">
                      <div className="flex items-center justify-between mb-1">
                        <span className="text-xs font-medium text-gray-700">
                          Orden: {order.referenceID}
                        </span>
                        <span className="text-xs text-gray-500">
                          {order.deliveryUnits?.length || 0} unidades
                        </span>
                      </div>
                      
                      <div className="space-y-1">
                        {order.deliveryUnits?.map((unit: any, unitIndex: number) => (
                          <div key={unitIndex} className="mb-2">
                            <IdentifierBadge 
                              lpn={unit.lpn} 
                              code={unit.code} 
                              size="sm"
                              className="mb-1"
                            />
                            <div className="text-xs text-gray-500">
                              {unit.items?.map((item: any) => item.description).join(', ')}
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
