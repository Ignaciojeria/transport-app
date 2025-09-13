import { CheckCircle, XCircle, Package } from 'lucide-react'
import { IdentifierBadge } from './IdentifierBadge'

interface VisitCardDeliveryUnitProps {
  unit: any
  uIdx: number
  status: string
  visitIndex: number
  orderIndex: number
  routeStarted: boolean
  orderReferenceID?: string
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
  orderReferenceID,
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
    <div className="bg-white rounded-lg p-4 border border-gray-200">
      {/* Badge prominente */}
      <div className="mb-3">
        <span className="inline-block bg-orange-500 text-white px-3 py-1 rounded-lg text-sm font-medium">
          {orderReferenceID || `CODE-${uIdx + 1}`}
        </span>
      </div>

      <div className="flex justify-between items-start">
        <div className="flex-1">
          {/* Título de la unidad */}
          <h5 className="text-base font-medium text-gray-800 mb-3 flex items-center">
            <Package className="w-4 h-4 mr-2 text-gray-600" />
            Unidad de Entrega {uIdx + 1}
          </h5>

          {/* Información del producto */}
          <div className="space-y-1 text-sm text-gray-600">
            {Array.isArray(unit.items) && unit.items.length > 0 && (
              <div className="flex items-center">
                <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
                <span>{unit.items[0]?.description}</span>
              </div>
            )}
            <div className="flex items-center">
              <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
              <span>{typeof unit.weight === 'number' ? `${unit.weight}kg` : unit.weight}</span>
            </div>
            <div className="flex items-center">
              <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
              <span>{typeof unit.volume === 'number' ? `${unit.volume}m³` : unit.volume}</span>
            </div>
          </div>

          {status === 'delivered' && (
            <div className="mt-3 inline-flex items-center text-xs px-2 py-1 rounded-full bg-green-100 text-green-700 border border-green-200">
              <CheckCircle className="w-3 h-3 mr-1" /> Evidencia registrada
            </div>
          )}
        </div>

        {/* Cantidad en la esquina derecha */}
        <div className="text-right ml-4">
          <span className="text-xs text-gray-500 block">Cant.</span>
          <span className="text-2xl font-bold text-indigo-600">
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
