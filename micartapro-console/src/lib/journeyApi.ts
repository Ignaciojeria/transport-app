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
