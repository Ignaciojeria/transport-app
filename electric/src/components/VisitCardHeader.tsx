import { User, MapPin } from 'lucide-react'

interface VisitCardHeaderProps {
  visit: any
  visitIndex: number
  onCenterOnVisit: (visitIndex: number) => void
}

export function VisitCardHeader({ visit, visitIndex, onCenterOnVisit }: VisitCardHeaderProps) {
  return (
    <div className="p-4 border-b border-gray-100">
      <div className="flex items-start space-x-3">
        <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 text-white rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0">
          {visit.sequenceNumber}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-sm font-bold text-gray-800 flex items-center mb-1">
            <User className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
            <span className="truncate">{visit.orders?.[0]?.contact?.fullName || 'Sin nombre'}</span>
          </h3>
          <p className="text-xs text-gray-600 flex items-start mb-2">
            <MapPin className="w-3 h-3 mr-1 mt-0.5 text-gray-500 flex-shrink-0" />
            <span className="line-clamp-2">{visit.addressInfo?.addressLine1}</span>
          </p>
        </div>
        <button
          onClick={() => onCenterOnVisit(visitIndex)}
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
