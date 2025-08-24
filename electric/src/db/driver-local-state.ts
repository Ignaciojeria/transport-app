import { z } from 'zod'
import { createCollection } from '@tanstack/react-db'
import { localStorageCollectionOptions } from '@tanstack/react-db'

// Esquema para estado local del driver, namespaciado por routeId
const DriverState = z.object({
  key: z.string(),
  value: z
    .union([
      z.literal('delivered'),
      z.literal('not-delivered'),
      z.literal('true'),
      z.literal('false'),
      z.null(),
    ])
    .nullable(),
})

export type DriverState = z.infer<typeof DriverState>

export const driverLocalState = createCollection(
  localStorageCollectionOptions<DriverState>({
    id: 'driver-local-state',
    storageKey: 'driver-local-state:v1',
    getKey: (item) => item.key,
  })
)

// Helpers para claves
export const routeStartedKey = (routeId: string) => `routeStarted:${routeId}`
export const deliveryKey = (routeId: string, vIdx: number, oIdx: number, uIdx: number) =>
  `delivery:${routeId}:${vIdx}-${oIdx}-${uIdx}`

// Mutadores
export function setRouteStarted(routeId: string, started: boolean) {
  const key = routeStartedKey(routeId)
  const existing = driverLocalState.get(key)
  if (existing) {
    driverLocalState.update(key, (d) => {
      d.value = started ? 'true' : 'false'
    })
  } else {
    driverLocalState.insert({ key, value: started ? 'true' : 'false' })
  }
}

export function setDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number,
  status: 'delivered' | 'not-delivered'
) {
  const key = deliveryKey(routeId, visitIndex, orderIndex, unitIndex)
  const existing = driverLocalState.get(key)
  if (existing) {
    driverLocalState.update(key, (d) => {
      d.value = status
    })
  } else {
    driverLocalState.insert({ key, value: status })
  }
}

export function getDeliveryStatus(
  routeId: string,
  visitIndex: number,
  orderIndex: number,
  unitIndex: number
): 'delivered' | 'not-delivered' | undefined {
  const item = driverLocalState.get(deliveryKey(routeId, visitIndex, orderIndex, unitIndex))
  return (item?.value as any) ?? undefined
}
