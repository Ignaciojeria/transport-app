import { User, Users, MapPin } from 'lucide-react'
import { useLanguage } from '../hooks/useLanguage'

interface VisitCardHeaderProps {
  visit: any
  viewMode?: 'list' | 'map'
}

export function VisitCardHeader({ visit, viewMode = 'list' }: VisitCardHeaderProps) {
  const { t } = useLanguage()
  
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
              <h3 className="text-lg font-bold text-gray-800 flex items-center">
                <MapPin className="w-5 h-5 mr-2 text-gray-600" />
                {visit.addressInfo?.addressLine1}
              </h3>
              {hasMultipleClients && (
                <p className="text-sm text-gray-500 mt-1">
                  {uniqueClients.length} {t.nextVisit.clients}
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
                    {uniqueClients.length} {t.nextVisit.clients}
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
                    {uniqueClients.length > 0 ? String(uniqueClients[0]) : 'Sin nombre'}
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

      </div>
    </div>
  )
}
