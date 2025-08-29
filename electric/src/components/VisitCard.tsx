import { User, MapPin, Package, CheckCircle, XCircle } from 'lucide-react'
import { VisitCardHeader } from './VisitCardHeader'
import { VisitCardOrders } from './VisitCardOrders'

interface VisitCardProps {
  visit: any
  visitIndex: number
  routeStarted: boolean
  onCenterOnVisit: (visitIndex: number) => void
  onOpenDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenNonDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, uIdx: number) => 'delivered' | 'not-delivered' | undefined
  shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean
}

export function VisitCard({
  visit,
  visitIndex,
  routeStarted,
  onCenterOnVisit,
  onOpenDelivery,
  onOpenNonDelivery,
  getDeliveryUnitStatus,
  shouldRenderByTab
}: VisitCardProps) {
  // Solo mostrar si tiene elementos pendientes para el tab actual
  const matchesForTab: number = (visit?.orders || []).reduce(
    (acc: number, order: any, orderIndex: number) => {
      const countInOrder = (order?.deliveryUnits || []).reduce(
        (a: number, _unit: any, uIdx: number) =>
          a + (shouldRenderByTab(getDeliveryUnitStatus(visitIndex, orderIndex, uIdx)) ? 1 : 0),
        0
      )
      return acc + countInOrder
    },
    0
  )

  if (matchesForTab === 0) return null

  return (
    <div className="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden border border-gray-100 active:scale-98">
                        <VisitCardHeader
                    visit={visit}
                    visitIndex={visitIndex}
                    onCenterOnVisit={onCenterOnVisit}
                  />
      
      <div className="p-4">
        <h4 className="text-sm font-medium text-gray-800 mb-3 flex items-center">
          <Package size={18} />
          <span className="ml-2">Unidades de Entrega:</span>
        </h4>

                            <VisitCardOrders
                      visit={visit}
                      visitIndex={visitIndex}
                      routeStarted={routeStarted}
                      onOpenDelivery={onOpenDelivery}
                      onOpenNonDelivery={onOpenNonDelivery}
                      getDeliveryUnitStatus={getDeliveryUnitStatus}
                      shouldRenderByTab={shouldRenderByTab}
                    />
      </div>
    </div>
  )
}
