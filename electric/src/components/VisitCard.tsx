import { VisitCardHeader } from './VisitCardHeader'
import { VisitCardOrders } from './VisitCardOrders'
import type { DeliveryGroup } from './GroupedDeliveryUtils'

interface VisitCardProps {
  visit: any
  visitIndex: number
  routeStarted: boolean
  onCenterOnVisit: (visitIndex: number) => void
  onOpenDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenNonDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenGroupedDelivery?: (visitIndex: number, group: DeliveryGroup) => void
  onOpenGroupedNonDelivery?: (visitIndex: number, group: DeliveryGroup) => void
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, uIdx: number) => 'delivered' | 'not-delivered' | undefined
  shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean
  viewMode?: 'list' | 'map' // Nuevo prop para controlar agrupación
}

export function VisitCard({
  visit,
  visitIndex,
  routeStarted,
  onCenterOnVisit,
  onOpenDelivery,
  onOpenNonDelivery,
  onOpenGroupedDelivery,
  onOpenGroupedNonDelivery,
  getDeliveryUnitStatus,
  shouldRenderByTab,
  viewMode = 'list'
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
    <div 
      className="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden border border-gray-100 active:scale-98 cursor-pointer"
      onClick={() => onCenterOnVisit(visitIndex)}
    >
                        <VisitCardHeader
                    visit={visit}
                    viewMode={viewMode}
                  />
      
      <div className="p-4">

                            <VisitCardOrders
                      visit={visit}
                      visitIndex={visitIndex}
                      routeStarted={routeStarted}
                      onOpenDelivery={onOpenDelivery}
                      onOpenNonDelivery={onOpenNonDelivery}
                      onOpenGroupedDelivery={onOpenGroupedDelivery}
                      onOpenGroupedNonDelivery={onOpenGroupedNonDelivery}
                      getDeliveryUnitStatus={getDeliveryUnitStatus}
                      shouldRenderByTab={shouldRenderByTab}
                      viewMode={viewMode}
                    />
      </div>
    </div>
  )
}
