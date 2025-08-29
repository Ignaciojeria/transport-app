import { Play, User, MapPin } from 'lucide-react'

interface NextVisitCardProps {
  nextVisit: any
  nextIdx: number
  onCenterOnVisit: (visitIndex: number) => void
}

export function NextVisitCard({ nextVisit, nextIdx, onCenterOnVisit }: NextVisitCardProps) {
  return (
    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl border-2 border-blue-200 p-4 mb-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-bold text-blue-800 flex items-center">
          <Play className="w-4 h-4 mr-2" />
          Siguiente Disponible
        </h3>
        <span className="text-xs text-blue-600 bg-blue-100 px-2 py-1 rounded-full font-medium">
          #{nextVisit.sequenceNumber}
        </span>
      </div>
      <div className="flex items-start space-x-3">
        <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-indigo-600 text-white rounded-lg flex items-center justify-center font-bold text-sm shadow-md flex-shrink-0">
          {nextVisit.sequenceNumber}
        </div>
        <div className="flex-1 min-w-0">
          <h4 className="text-sm font-bold text-gray-800 flex items-center mb-1">
            <User className="w-3 h-3 mr-1 text-gray-600 flex-shrink-0" />
            <span className="truncate">{nextVisit.addressInfo?.contact?.fullName}</span>
          </h4>
          <p className="text-xs text-gray-600 flex items-start">
            <MapPin className="w-3 h-3 mr-1 mt-0.5 text-gray-500 flex-shrink-0" />
            <span className="line-clamp-2">{nextVisit.addressInfo?.addressLine1}</span>
          </p>
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
    </div>
  )
}
