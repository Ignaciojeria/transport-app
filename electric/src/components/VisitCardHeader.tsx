import { User, MapPin, Users } from 'lucide-react'

interface VisitCardHeaderProps {
  visit: any
  visitIndex: number
  onCenterOnVisit: (visitIndex: number) => void
  viewMode?: 'list' | 'map'
}

export function VisitCardHeader({ visit, visitIndex, onCenterOnVisit, viewMode = 'list' }: VisitCardHeaderProps) {
  // Obtener todos los clientes únicos de la visita
  const uniqueClients = Array.from(new Set(
    (visit.orders || []).map((order: any) => order.contact?.fullName).filter(Boolean)
  ))
  
  const hasMultipleClients = uniqueClients.length > 1
  
  return (
    <div className="p-4 border-b border-gray-100">
      <div className="flex items-start space-x-3">
        <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 text-white rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0">
          {visit.sequenceNumber}
        </div>
        <div className="flex-1 min-w-0">
          {viewMode === 'list' ? (
            // En modo lista: solo mostrar la dirección en negrita arriba
            <>
              <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                <MapPin className="w-3 h-3 mr-1 text-gray-500 flex-shrink-0" />
                <span className="line-clamp-2">{visit.addressInfo?.addressLine1}</span>
              </h3>
              {hasMultipleClients && (
                <p className="text-xs text-gray-600 mb-2">
                  <span className="text-indigo-700 font-medium">{uniqueClients.length} clientes</span>
                </p>
              )}
            </>
          ) : (
            // En modo mapa: comportamiento original
            <>
              {hasMultipleClients ? (
                <div className="mb-1">
                  <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                    <Users className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
                    <span className="text-indigo-700">{uniqueClients.length} clientes</span>
                  </h3>
                  <div className="text-xs text-gray-600 mb-1">
                    {uniqueClients.map((client, index) => (
                      <span key={index} className="inline-block bg-gray-100 text-gray-700 px-2 py-1 rounded mr-1 mb-1">
                        {client}
                      </span>
                    ))}
                  </div>
                </div>
              ) : (
                <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
                  <User className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
                  <span className="truncate">{uniqueClients[0] || 'Sin nombre'}</span>
                </h3>
              )}
              <p className="text-xs text-gray-600 flex items-start mb-2">
                <MapPin className="w-3 h-3 mr-1 mt-0.5 text-gray-500 flex-shrink-0" />
                <span className="line-clamp-2">{visit.addressInfo?.addressLine1}</span>
              </p>
            </>
          )}
        </div>
        <button
          onClick={() => {
            console.log('📍 PIN CLICKED! visitIndex:', visitIndex, 'sequenceNumber:', visit.sequenceNumber)
            onCenterOnVisit(visitIndex)
          }}
          className="w-8 h-8 bg-blue-50 hover:bg-blue-100 border border-blue-200 text-blue-600 rounded-lg flex items-center justify-center transition-all duration-200 hover:shadow-md active:scale-95 flex-shrink-0"
          aria-label={`Ver en mapa - Visita ${visit.sequenceNumber}`}
          title="Ver en mapa"
        >
          <MapPin className="w-4 h-4" />
        </button>
      </div>
    </div>
  )
}
