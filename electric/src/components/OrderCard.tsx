import { Package } from 'lucide-react'
import { IdentifierBadge } from './IdentifierBadge'

interface OrderCardProps {
  order: any
  orderIndex: number
  visitIndex: number
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined
  children?: React.ReactNode
  className?: string
}

export function OrderCard({ 
  order, 
  orderIndex, 
  visitIndex, 
  getDeliveryUnitStatus,
  children,
  className = ''
}: OrderCardProps) {
  const deliveryUnits = order?.deliveryUnits || []
  
  return (
    <div className={`bg-white rounded-xl shadow-md border border-gray-100 overflow-hidden ${className}`}>
      {/* Header de la orden */}
      <div className="bg-gradient-to-r from-gray-50 to-blue-50 px-4 py-3 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Package className="w-4 h-4 text-gray-600" />
            <span className="font-medium text-gray-800">
              Orden {orderIndex + 1}
            </span>
          </div>
          <span className="text-sm text-gray-500">
            {deliveryUnits.length} unidad{deliveryUnits.length !== 1 ? 'es' : ''}
          </span>
        </div>
        {order.referenceID && (
          <div className="mt-2">
            <span className="text-xs text-gray-600 font-mono bg-gray-100 px-2 py-1 rounded">
              {order.referenceID}
            </span>
          </div>
        )}
      </div>

      {/* Contenido de la orden */}
      <div className="p-4">
        {children}
      </div>
    </div>
  )
}

interface DeliveryUnitCardProps {
  unit: any
  unitIndex: number
  visitIndex: number
  orderIndex: number
  status: 'delivered' | 'not-delivered' | undefined
  orderReferenceID?: string
  children?: React.ReactNode
  className?: string
}

export function DeliveryUnitCard({ 
  unit, 
  unitIndex, 
  visitIndex, 
  orderIndex, 
  status,
  orderReferenceID,
  children,
  className = ''
}: DeliveryUnitCardProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'delivered':
        return 'bg-green-50 border-green-200'
      case 'not-delivered':
        return 'bg-red-50 border-red-200'
      default:
        return 'bg-gray-50 border-gray-200'
    }
  }

  return (
    <div className={`bg-gradient-to-r from-gray-50 to-blue-50 rounded-lg p-3 border ${getStatusColor(status || '')} ${className}`}>
      {/* Header de la unidad con identificadores prominentes */}
      <div className="mb-3">
        <IdentifierBadge 
          lpn={unit.lpn} 
          referenceID={orderReferenceID}
          size="sm"
          className="mb-2"
        />
        <h5 className="text-sm font-medium text-gray-800">
          Unidad de Entrega {unitIndex + 1}
        </h5>
      </div>

      {/* Contenido de la unidad */}
      {children}
    </div>
  )
}
