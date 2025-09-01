import { CheckCircle, XCircle } from 'lucide-react'

interface VisitCardDeliveryUnitProps {
  unit: any
  uIdx: number
  status: string
  visitIndex: number
  orderIndex: number
  routeStarted: boolean
  onOpenDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenNonDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
}

export function VisitCardDeliveryUnit({
  unit,
  uIdx,
  status,
  visitIndex,
  orderIndex,
  routeStarted,
  onOpenDelivery,
  onOpenNonDelivery
}: VisitCardDeliveryUnitProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'delivered':
        return 'bg-green-50 border-green-200'
      case 'not-delivered':
        return 'bg-red-50 border-red-200'
      default:
        return 'bg-white border-gray-200'
    }
  }

  return (
    <div
      className={`bg-gradient-to-r from-gray-50 to-blue-50 rounded-lg p-3 border ${getStatusColor(status).replace('bg-white ', '')}`}
    >
      <div className="flex justify-between items-start mb-2">
        <div className="flex-1 min-w-0">
          <h5 className="text-sm font-medium text-gray-800 mb-2 truncate">
            Unidad de Entrega {uIdx + 1}
          </h5>
          {Array.isArray(unit.items) && unit.items.length > 0 && (
            <div className="flex items-center space-x-1 mb-2">
              <span className="w-1.5 h-1.5 bg-indigo-500 rounded-full"></span>
              <span className="text-xs text-gray-700 truncate">
                {unit.items[0]?.description}
              </span>
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
          {status === 'delivered' && (
            <div className="mt-2 inline-flex items-center text-[10px] px-2 py-0.5 rounded-full bg-green-100 text-green-700 border border-green-200">
              <CheckCircle className="w-3 h-3 mr-1" /> Evidencia registrada
            </div>
          )}
        </div>
        <div className="text-right ml-3">
          <span className="text-xs text-gray-500 block">Cant.</span>
          <span className="text-xl font-bold text-indigo-600">
            {(unit.items || []).reduce((a: number, it: any) => a + (Number(it?.quantity) || 0), 0)}
          </span>
        </div>
      </div>

      {routeStarted && (
        <div className="flex space-x-2 mt-3">
          {status === 'delivered' ? (
            // Si está entregado, mostrar solo opción de cambiar a no entregado
            <button
              onClick={() => onOpenNonDelivery(visitIndex, orderIndex, uIdx)}
              className="w-full flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-red-100 text-red-700 hover:bg-red-200"
            >
              <XCircle size={16} />
              <span>Cambiar a no entregado</span>
            </button>
          ) : status === 'not-delivered' ? (
            // Si está no entregado, mostrar solo opción de cambiar a entregado
            <button
              onClick={() => onOpenDelivery(visitIndex, orderIndex, uIdx)}
              className="w-full flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-green-100 text-green-700 hover:bg-green-200"
            >
              <CheckCircle size={16} />
              <span>Cambiar a entregado</span>
            </button>
          ) : (
            // Si está pendiente, mostrar ambas opciones originales
            <>
              <button
                onClick={() => onOpenDelivery(visitIndex, orderIndex, uIdx)}
                className="flex-1 flex items-center justify-center space-x-2 py-2 px-3 rounded-md font-medium transition-colors bg-green-100 text-green-700 hover:bg-green-200"
              >
                <CheckCircle size={16} />
                <span>entregar</span>
              </button>
              <button
                onClick={() => onOpenNonDelivery(visitIndex, orderIndex, uIdx)}
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
  )
}
