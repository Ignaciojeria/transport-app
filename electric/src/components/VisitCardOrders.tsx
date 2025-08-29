import { VisitCardDeliveryUnit } from './VisitCardDeliveryUnit'

interface VisitCardOrdersProps {
  visit: any
  visitIndex: number
  routeStarted: boolean
  onOpenDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  onOpenNonDelivery: (visitIndex: number, orderIndex: number, uIdx: number) => void
  getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, uIdx: number) => 'delivered' | 'not-delivered' | undefined
  shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean
}

export function VisitCardOrders({
  visit,
  visitIndex,
  routeStarted,
  onOpenDelivery,
  onOpenNonDelivery,
  getDeliveryUnitStatus,
  shouldRenderByTab
}: VisitCardOrdersProps) {
  return (
    <>
      {visit.orders?.map((order: any, orderIndex: number) => (
        <div key={orderIndex} className="mb-4">
          <div className="mb-2">
            <span className="inline-block bg-gradient-to-r from-orange-400 to-red-500 text-white px-2 py-1 rounded-lg text-xs font-medium">
              {order.referenceID}
            </span>
          </div>
          {(order.deliveryUnits || [])
            .map((unit: any, uIdx: number) => ({
              unit,
              uIdx,
              status: getDeliveryUnitStatus(visitIndex, orderIndex, uIdx),
            }))
            .filter((x: any) => shouldRenderByTab(x.status))
            .map((x: any) => (
                                        <VisitCardDeliveryUnit
                            key={x.uIdx}
                            unit={x.unit}
                            uIdx={x.uIdx}
                            status={x.status}
                            visitIndex={visitIndex}
                            orderIndex={orderIndex}
                            routeStarted={routeStarted}
                            onOpenDelivery={onOpenDelivery}
                            onOpenNonDelivery={onOpenNonDelivery}
                          />
            ))}
        </div>
      ))}
    </>
  )
}
