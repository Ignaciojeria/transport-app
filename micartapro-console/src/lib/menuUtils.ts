import { supabase } from './supabase'

/**
 * Obtiene el menuID más reciente del usuario autenticado
 * @returns El menuID o null si no existe
 */
export async function getLatestMenuId(userId: string): Promise<string | null> {
  try {
    const { data, error } = await supabase
      .from('user_menus')
      .select('menu_id, created_at')
      .eq('user_id', userId)
      .order('created_at', { ascending: false })
      .limit(1)
      .single()

    if (error) {
      console.error('Error obteniendo menuID:', error)
      return null
    }

    return data?.menu_id || null
  } catch (err) {
    console.error('Error en getLatestMenuId:', err)
    return null
  }
}

/**
 * Genera la URL de la carta del restaurante
 * @param userId - ID del usuario
 * @param menuId - ID del menú
 * @returns URL completa de la carta
 */
export function generateMenuUrl(userId: string, menuId: string): string {
  const baseUrl = typeof window !== 'undefined' 
    ? window.location.origin.includes('localhost')
      ? 'http://localhost:5173'
      : 'https://cadorago.web.app'
    : 'https://cadorago.web.app'
  
  return `${baseUrl}/?userID=${userId}&menuID=${menuId}`
}

