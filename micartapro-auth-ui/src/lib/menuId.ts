import { v7 as uuidv7 } from 'uuid'
import { supabase } from './supabase'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ||
  (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ? 'http://localhost:8082'
    : 'https://micartapro-backend-27303662337.us-central1.run.app')

// Cache para evitar llamadas simultáneas para el mismo usuario
const pendingRequests = new Map<string, Promise<string>>()

/**
 * Verifica si el usuario tiene al menos un menú.
 */
export async function hasUserMenus(userId: string): Promise<boolean> {
  const { data, error } = await supabase
    .from('user_menus')
    .select('menu_id')
    .eq('user_id', userId)
    .limit(1)
  if (error || !data || data.length === 0) return false
  return true
}

/**
 * Crea un nuevo menú con slug para un usuario (primera vez).
 * Inserta en user_menus, espera a que el menú exista (webhook), crea el slug.
 */
export async function createMenuWithSlug(
  userId: string,
  slug: string,
  accessToken: string
): Promise<string> {
  const newMenuId = uuidv7()
  const { data: inserted, error: insertError } = await supabase
    .from('user_menus')
    .insert({ user_id: userId, menu_id: newMenuId })
    .select('menu_id')
    .single()

  if (insertError || !inserted?.menu_id) {
    throw new Error(insertError?.message || 'Error al crear el menú')
  }

  const normalizedSlug = slug.trim().toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '')
  if (normalizedSlug) {
    for (let i = 0; i < 24; i++) {
      await new Promise((r) => setTimeout(r, 2500))
      const res = await fetch(`${API_BASE_URL}/menu/${encodeURIComponent(inserted.menu_id)}`, {
        headers: { Authorization: `Bearer ${accessToken}` },
      })
      if (res.ok) {
        const { data: existing } = await supabase.from('menu_slugs').select('id').eq('slug', normalizedSlug).maybeSingle()
        if (existing) throw new Error('SLUG_EXISTS')
        await supabase.from('menu_slugs').update({ is_active: false }).eq('menu_id', inserted.menu_id)
        const { error: slugErr } = await supabase.from('menu_slugs').insert({
          menu_id: inserted.menu_id,
          slug: normalizedSlug,
          is_active: true,
        })
        if (slugErr) throw new Error(slugErr.message)
        break
      }
      if (i === 23) console.warn('Timeout esperando menú para crear slug')
    }
  }
  return inserted.menu_id
}

/**
 * Obtiene el menuID más reciente del usuario. Si no tiene ninguno, retorna null
 * (el caller debe mostrar el formulario de slug y llamar createMenuWithSlug).
 */
export async function getOrCreateMenuId(userId: string): Promise<string | null> {
  // Si ya hay una petición en curso para este usuario, esperar a que termine
  if (pendingRequests.has(userId)) {
    console.log('⏳ Ya hay una petición en curso para este usuario, esperando...')
    return pendingRequests.get(userId)!
  }

  const requestPromise = getOrCreateMenuIdInternal(userId)
  pendingRequests.set(userId, requestPromise)
  try {
    return await requestPromise
  } finally {
    pendingRequests.delete(userId)
  }
}

async function getOrCreateMenuIdInternal(userId: string): Promise<string | null> {
  try {
    console.log('🔍 Buscando menuID para usuario:', userId)
    
    // 0. Verificar que tenemos una sesión activa
    const { data: { session }, error: sessionError } = await supabase.auth.getSession()
    if (sessionError || !session) {
      throw new Error('No hay sesión activa. Debes estar autenticado para crear un menuID.')
    }
    
    // Verificar que el userId coincide con el usuario autenticado
    if (session.user.id !== userId) {
      throw new Error(`El userId proporcionado (${userId}) no coincide con el usuario autenticado (${session.user.id})`)
    }
    
    console.log('✅ Sesión verificada:', {
      userId: session.user.id,
      email: session.user.email
    })
    
    // 1. Verificar si el usuario tiene AL MENOS un menuID
    // Buscamos cualquier menuID del usuario (puede tener múltiples)
    const { data: existingMenus, error: fetchError } = await supabase
      .from('user_menus')
      .select('menu_id, created_at')
      .eq('user_id', userId)
      .order('created_at', { ascending: false })
      .limit(1)

    // Manejar errores de fetch (solo errores críticos, no "no encontrado")
    if (fetchError) {
      // Si el error es que la tabla no existe, lanzar error
      if (fetchError.code === '42P01') {
        throw new Error('La tabla user_menus no existe. Ejecuta el SQL en SUPABASE_SETUP.md')
      }
      
      // Si el error es de RLS (permiso denegado), lanzar error
      if (fetchError.code === '42501' || fetchError.message?.includes('permission denied') || fetchError.message?.includes('RLS')) {
        throw new Error('Permiso denegado por RLS. Verifica las políticas de seguridad en Supabase.')
      }
      
      // Si el error es "no encontrado" (PGRST116), es normal, continuar para crear uno
      if (fetchError.code === 'PGRST116') {
        console.log('ℹ️ No se encontró ningún menuID para este usuario, se creará uno por defecto')
      } else {
        // Cualquier otro error, loguear pero continuar (intentaremos crear uno)
        console.warn('⚠️ Error al buscar menuID (continuando para crear uno por defecto):', fetchError.message)
      }
    }

    // Si el usuario ya tiene al menos un menuID, retornar el más reciente
    if (existingMenus && existingMenus.length > 0 && existingMenus[0].menu_id) {
      const latestMenuId = existingMenus[0].menu_id
      console.log('✅ El usuario ya tiene menús, usando el más reciente:', latestMenuId)
      return latestMenuId
    }

    // Usuario nuevo: retornar null para que el callback muestre el formulario de slug
    console.log('📝 Usuario nuevo sin menús, se debe mostrar formulario de slug')
    return null
  } catch (error) {
    console.error('❌ Error completo en getOrCreateMenuId:', error)
    throw error
  }
}

