/**
 * API de jornadas: obtener jornada activa y abrir nueva jornada.
 */

import { API_BASE_URL } from './config'

export interface ActiveJourney {
  id: string
  menuId: string
  status: string
  openedAt: string
  openedBy: string
  reason?: string
}

export interface CreateJourneyResponse {
  id: string
  menuId: string
  status: string
  openedAt: string
  openedBy: string
  reason?: string
}

export interface JourneyListItem {
  id: string
  menuId: string
  status: string
  openedAt: string
  closedAt?: string
  reportPdfUrl?: string
  reportXlsxUrl?: string
}

/**
 * Obtiene la jornada activa (OPEN) del menú. Devuelve null si no hay (404).
 */
export async function getActiveJourney(
  menuId: string,
  accessToken: string
): Promise<ActiveJourney | null> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/active-journey`,
    {
      headers: { Authorization: `Bearer ${accessToken}` }
    }
  )
  if (res.status === 404) return null
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  const j = await res.json() as Record<string, unknown>
  // Backend puede devolver PascalCase (ID, MenuID) o camelCase
  return {
    id: (j.id ?? j.ID) as string,
    menuId: (j.menuId ?? j.MenuID) as string,
    status: (j.status ?? j.Status) as string,
    openedAt: (j.openedAt ?? j.OpenedAt) as string,
    openedBy: (j.openedBy ?? j.OpenedBy) as string,
    reason: (j.reason ?? j.OpenedReason) as string | undefined
  }
}

/**
 * Crea una nueva jornada OPEN para el menú. Falla con error si ya hay una abierta (409).
 */
export async function createJourney(
  menuId: string,
  accessToken: string,
  openedBy: 'USER' | 'SYSTEM' = 'USER',
  reason?: string
): Promise<CreateJourneyResponse> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/journeys`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`
      },
      body: JSON.stringify({ openedBy, reason: reason ?? 'Apertura manual' })
    }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  return res.json()
}

/**
 * Lista jornadas del menú (abiertas y cerradas), ordenadas por fecha desc.
 */
export async function getJourneys(
  menuId: string,
  accessToken: string
): Promise<JourneyListItem[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/journeys`,
    {
      headers: { Authorization: `Bearer ${accessToken}` }
    }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  const raw = (await res.json()) as Record<string, unknown>[]
  return raw.map((j) => ({
    id: (j.id ?? j.ID) as string,
    menuId: (j.menuId ?? j.MenuID) as string,
    status: (j.status ?? j.Status) as string,
    openedAt: (j.openedAt ?? j.OpenedAt) as string,
    closedAt: (j.closedAt ?? j.ClosedAt) as string | undefined,
    reportPdfUrl: (j.reportPdfUrl ?? j.ReportPDFURL) as string | undefined,
    reportXlsxUrl: (j.reportXlsxUrl ?? j.ReportXLSXURL) as string | undefined
  }))
}

export interface JourneyStatsProduct {
  productName: string
  quantitySold: number
  totalRevenue: number
  percentage: number
  percentageByQuantity: number
}

export interface JourneyStats {
  totalRevenue: number
  totalOrders: number
  itemsOrdered?: number
  products: JourneyStatsProduct[]
  revenueByStatus?: {
    delivered: number
    dispatched: number
    pending: number
    cancelled: number
  }
  ordersByStatus?: {
    delivered: number
    dispatched: number
    pending: number
    cancelled: number
  }
}

/**
 * Obtiene estadísticas de ventas de una jornada (productos más vendidos, revenue, etc.).
 */
export async function getJourneyStats(
  menuId: string,
  journeyId: string,
  accessToken: string
): Promise<JourneyStats> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/journeys/${encodeURIComponent(journeyId)}/stats`,
    {
      headers: { Authorization: `Bearer ${accessToken}` }
    }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
  return res.json()
}

/**
 * Cierra la jornada activa (OPEN) del menú. Falla con 404 si no hay jornada abierta.
 */
export async function closeJourney(
  menuId: string,
  accessToken: string
): Promise<void> {
  const res = await fetch(
    `${API_BASE_URL}/api/menus/${encodeURIComponent(menuId)}/journeys/close`,
    {
      method: 'POST',
      headers: { Authorization: `Bearer ${accessToken}` }
    }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Error ${res.status}`)
  }
}
