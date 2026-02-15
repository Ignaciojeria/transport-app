/**
 * Cliente para los endpoints de transición de estado de órdenes del backend.
 * Todas las fechas se envían en UTC (ej: 2026-02-03T18:25:09.506000000Z).
 */

import { API_BASE_URL } from './config'

/** Fecha actual en UTC (ej: 2026-02-03T18:25:09.506000000Z). */
function getUtcIsoString(): string {
  const iso = new Date().toISOString()
  return iso.replace(/\.(\d{3})Z$/, (_, ms) => `.${ms}000000Z`)
}

export type OrderStatusTransitionError = { message: string; status?: number }

async function orderFetch(
  path: string,
  accessToken: string,
  body: Record<string, unknown>
): Promise<void> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${accessToken}`
    },
    body: JSON.stringify(body)
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
}

/**
 * Marca que Cocina/Barra está preparando (PENDING → IN_PROGRESS).
 * @param menuId - ID del menú
 * @param aggregateId - ID del agregado de la orden
 * @param station - KITCHEN | BAR
 * @param itemKeys - Si se omite o está vacío, aplica a todos los ítems de esa estación
 */
export async function startPreparation(
  menuId: string,
  aggregateId: number,
  accessToken: string,
  station: string,
  itemKeys: string[] = []
): Promise<void> {
  await orderFetch(
    `/menu/${encodeURIComponent(menuId)}/orders/${aggregateId}/start-preparation`,
    accessToken,
    { station, itemKeys, updatedAt: getUtcIsoString() }
  )
}

/**
 * Marca ítems como listos (IN_PROGRESS → READY).
 * @param menuId - ID del menú
 * @param aggregateId - ID del agregado de la orden
 * @param station - KITCHEN | BAR
 * @param itemKeys - item_key de los ítems que pasan a READY (requerido)
 */
export async function markReady(
  menuId: string,
  aggregateId: number,
  accessToken: string,
  station: string,
  itemKeys: string[]
): Promise<void> {
  if (itemKeys.length === 0) return
  await orderFetch(
    `/menu/${encodeURIComponent(menuId)}/orders/${aggregateId}/mark-ready`,
    accessToken,
    { station, itemKeys, updatedAt: getUtcIsoString() }
  )
}

/**
 * Marca la orden como terminada: READY → DELIVERED (PICKUP/retiro) o DISPATCHED (DELIVERY/despacho) según fulfillment.
 */
export async function dispatchOrder(
  menuId: string,
  aggregateId: number,
  accessToken: string
): Promise<void> {
  await orderFetch(
    `/menu/${encodeURIComponent(menuId)}/orders/${aggregateId}/dispatch`,
    accessToken,
    { updatedAt: getUtcIsoString() }
  )
}

/**
 * Cancela la orden (PENDING/IN_PROGRESS/READY → CANCELLED).
 */
export async function cancelOrder(
  menuId: string,
  aggregateId: number,
  accessToken: string,
  reason?: string
): Promise<void> {
  await orderFetch(
    `/menu/${encodeURIComponent(menuId)}/orders/${aggregateId}/cancel`,
    accessToken,
    { ...(reason != null && reason !== '' ? { reason } : {}), updatedAt: getUtcIsoString() }
  )
}

/** Ítem de una orden (para detalle). */
export interface OrderItem {
  itemName: string
  quantity: number
  unit: string
  totalPrice: number
  station?: string | null
}

/**
 * Obtiene los ítems de una orden por aggregate_id.
 */
export async function getOrderItems(
  menuId: string,
  aggregateId: number,
  accessToken: string
): Promise<OrderItem[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/orders/${aggregateId}/items`,
    { headers: { Authorization: `Bearer ${accessToken}` } }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  return res.json()
}

/** Orden pendiente (sin jornada asignada), con ítems agrupados por agregado. */
export interface PendingOrder {
  aggregateId: number
  trackingId: string
  createdAt: string
  scheduledFor?: string | null
  totalAmount: number
  status: string
  items: OrderItem[]
}

/** Filtro opcional para órdenes pendientes. */
export interface PendingOrdersFilter {
  fromDate?: string // RFC3339
  toDate?: string // RFC3339
}

/**
 * Obtiene órdenes pendientes (journey_id IS NULL) del menú.
 */
export async function getPendingOrders(
  menuId: string,
  accessToken: string,
  filter?: PendingOrdersFilter
): Promise<PendingOrder[]> {
  const params = new URLSearchParams()
  if (filter?.fromDate) params.set('fromDate', filter.fromDate)
  if (filter?.toDate) params.set('toDate', filter.toDate)
  const qs = params.toString()
  const url = `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/pending-orders${qs ? `?${qs}` : ''}`
  const res = await fetch(url, {
    headers: { Authorization: `Bearer ${accessToken}` }
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  return res.json()
}

/** Resultado de asignar una orden a la jornada. */
export interface AssignOrderResult {
  aggregateId: number
  orderNumber?: number | null
  assigned: boolean
}

/**
 * Asigna órdenes pendientes a la jornada activa.
 */
export async function assignOrdersToJourney(
  menuId: string,
  aggregateIds: number[],
  accessToken: string
): Promise<AssignOrderResult[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/journeys/assign-orders`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      },
      body: JSON.stringify({ aggregateIds })
    }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  const data = (await res.json()) as { results: AssignOrderResult[] }
  return data.results
}
