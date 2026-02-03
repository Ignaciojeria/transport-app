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
 * Marca la orden como terminada: READY → DISPATCHED (retiro) o DELIVERED (despacho) según fulfillment.
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
