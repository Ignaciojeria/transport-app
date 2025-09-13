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
      <div className="flex items-center justify-between">
        {/* Número de secuencia */}
        <div className="w-10 h-10 bg-indigo-600 text-white rounded-lg flex items-center justify-center font-bold text-lg flex-shrink-0">
          {visit.sequenceNumber}
        </div>

        {/* Información del cliente */}
        <div className="flex-1 mx-4">
          {viewMode === 'list' ? (
            // En modo lista: solo dirección prominente
            <>
              <h3 className="text-base font-bold text-gray-800 flex items-center">
                <MapPin className="w-4 h-4 mr-2 text-gray-600" />
                {visit.addressInfo?.addressLine1}
              </h3>
              {hasMultipleClients && (
                <p className="text-xs text-gray-500 mt-1">
                  {uniqueClients.length} clientes
                </p>
              )}
            </>
          ) : (
            // En modo mapa: comportamiento original
            <>
              {hasMultipleClients ? (
                <div>
                  <h3 className="text-base font-medium text-gray-800 flex items-center mb-1">
                    <Users className="w-4 h-4 mr-2 text-gray-600" />
                    {uniqueClients.length} clientes
                  </h3>
                  <p className="text-sm text-gray-600 flex items-center">
                    <MapPin className="w-4 h-4 mr-2 text-gray-500" />
                    {visit.addressInfo?.addressLine1}
                  </p>
                </div>
              ) : (
                <div>
                  <h3 className="text-base font-medium text-gray-800 flex items-center mb-1">
                    <User className="w-4 h-4 mr-2 text-gray-600" />
                    {uniqueClients[0] || 'Sin nombre'}
                  </h3>
                  <p className="text-sm text-gray-600 flex items-center">
                    <MapPin className="w-4 h-4 mr-2 text-gray-500" />
                    {visit.addressInfo?.addressLine1}
                  </p>
                </div>
              )}
            </>
          )}
        </div>

        {/* Botón del pin */}
        <button
          onClick={() => {
            console.log('📍 PIN CLICKED! visitIndex:', visitIndex, 'sequenceNumber:', visit.sequenceNumber)
            onCenterOnVisit(visitIndex)
          }}
          className="w-10 h-10 bg-blue-50 hover:bg-blue-100 border border-blue-200 text-blue-600 rounded-lg flex items-center justify-center transition-all duration-200 hover:shadow-md active:scale-95 flex-shrink-0"
          aria-label={`Ver en mapa - Visita ${visit.sequenceNumber}`}
          title="Ver en mapa"
        >
          <MapPin className="w-5 h-5" />
        </button>
      </div>
    </div>
  )
}
