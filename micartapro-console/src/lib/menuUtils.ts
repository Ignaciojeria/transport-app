import { supabase } from './supabase'

const STORAGE_BASE_URL = "https://storage.googleapis.com/micartapro-menus"

/**
 * Obtiene los datos del restaurante desde Google Cloud Storage
 * @param userId - ID del usuario
 * @param menuId - ID del menú
 * @returns Promise con los datos del restaurante o null si hay error
 */
export async function fetchRestaurantData(userId: string, menuId: string): Promise<any | null> {
  try {
    // Obtener latest.json
    const latest = await getLatestJson(userId, menuId)
    if (!latest || !latest.filename) {
      return null
    }
    
    // Obtener los datos del restaurante
    const menuUrl = `${STORAGE_BASE_URL}/${userId}/menus/${menuId}/${latest.filename}`
    const response = await fetch(menuUrl, {
      cache: 'no-store'
    })
    
    if (!response.ok) {
      return null
    }
    
    return await response.json()
  } catch (error) {
    console.error('Error obteniendo datos del restaurante:', error)
    return null
  }
}

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
 * @param lang - Idioma opcional ('ES', 'PT', 'EN')
 * @returns URL completa de la carta
 */
export function generateMenuUrl(userId: string, menuId: string, lang?: string): string {
  const baseUrl = typeof window !== 'undefined' 
    ? window.location.origin.includes('localhost')
      ? 'http://localhost:5173'
      : 'https://cadorago.web.app'
    : 'https://cadorago.web.app'
  
  let url = `${baseUrl}/?userID=${userId}&menuID=${menuId}`
  
  // Agregar parámetro de idioma si se proporciona
  if (lang && ['ES', 'PT', 'EN'].includes(lang)) {
    url += `&lang=${lang}`
  }
  
  return url
}

/**
 * Obtiene el contenido de latest.json desde Google Cloud Storage
 * @param userId - ID del usuario
 * @param menuId - ID del menú
 * @returns El objeto latest.json con el campo filename
 */
async function getLatestJson(userId: string, menuId: string): Promise<{ filename: string } | null> {
  try {
    const latestUrl = `${STORAGE_BASE_URL}/${userId}/menus/${menuId}/latest.json`
    const response = await fetch(latestUrl, {
      cache: 'no-store' // Evitar caché para obtener siempre la versión más reciente
    })
    
    if (!response.ok) {
      if (response.status === 404) {
        return null
      }
      throw new Error(`Error al obtener latest.json: ${response.status} ${response.statusText}`)
    }
    
    const data = await response.json()
    return data
  } catch (error) {
    console.error('Error obteniendo latest.json:', error)
    throw error
  }
}

/**
 * Hace polling de latest.json hasta que el idempotencyKey coincida
 * @param userId - ID del usuario
 * @param menuId - ID del menú
 * @param idempotencyKey - Clave de idempotencia a esperar
 * @param maxAttempts - Número máximo de intentos (default: 60)
 * @param intervalMs - Intervalo entre intentos en milisegundos (default: 2000)
 * @returns El menú actualizado cuando el idempotencyKey coincide
 */
export async function pollUntilMenuUpdated(
  userId: string,
  menuId: string,
  idempotencyKey: string,
  maxAttempts: number = 60,
  intervalMs: number = 2000
): Promise<any> {
  let attempts = 0
  
  while (attempts < maxAttempts) {
    try {
      const latest = await getLatestJson(userId, menuId)
      
      if (latest && latest.filename) {
        // El filename en latest.json es "idempotencyKey.json", necesitamos extraer solo el idempotencyKey
        const filenameWithoutExt = latest.filename.replace('.json', '')
        
        if (filenameWithoutExt === idempotencyKey) {
          // Coincide! Obtener el menú actualizado
          const menuUrl = `${STORAGE_BASE_URL}/${userId}/menus/${menuId}/${latest.filename}`
          const response = await fetch(menuUrl, {
            cache: 'no-store'
          })
          
          if (!response.ok) {
            throw new Error(`Error al obtener menú actualizado: ${response.status} ${response.statusText}`)
          }
          
          const menuData = await response.json()
          return menuData
        }
      }
      
      // Esperar antes del siguiente intento
      await new Promise(resolve => setTimeout(resolve, intervalMs))
      attempts++
    } catch (error) {
      console.error(`Error en intento ${attempts + 1}:`, error)
      // Continuar intentando a menos que sea un error crítico
      await new Promise(resolve => setTimeout(resolve, intervalMs))
      attempts++
    }
  }
  
  throw new Error(`Timeout: No se pudo obtener el menú actualizado después de ${maxAttempts} intentos`)
}

/**
 * Hace polling para verificar que el menú exista en GCS (verificando latest.json)
 * Útil cuando el usuario se registra por primera vez y el menú se está creando
 * @param userId - ID del usuario
 * @param menuId - ID del menú
 * @param maxAttempts - Número máximo de intentos (default: 30)
 * @param intervalMs - Intervalo entre intentos en milisegundos (default: 2000)
 * @returns true si el menú existe, false si no se encontró después de todos los intentos
 */
export async function pollUntilMenuExists(
  userId: string,
  menuId: string,
  maxAttempts: number = 30,
  intervalMs: number = 2000
): Promise<boolean> {
  let attempts = 0
  
  while (attempts < maxAttempts) {
    try {
      const latest = await getLatestJson(userId, menuId)
      
      if (latest && latest.filename) {
        // El menú existe
        return true
      }
      
      // Esperar antes del siguiente intento
      await new Promise(resolve => setTimeout(resolve, intervalMs))
      attempts++
    } catch (error: any) {
      // Si es 404, el menú aún no existe, continuar intentando
      if (error.message?.includes('404')) {
        await new Promise(resolve => setTimeout(resolve, intervalMs))
        attempts++
        continue
      }
      
      // Otro error, loguear pero continuar intentando
      console.error(`Error en intento ${attempts + 1}:`, error)
      await new Promise(resolve => setTimeout(resolve, intervalMs))
      attempts++
    }
  }
  
  // No se encontró el menú después de todos los intentos
  return false
}


export interface Entitlement {
  v: number
  user_id: string
  plan: string
  status: string
  access: boolean
  starts_at: string
  ends_at: string
}

export async function fetchEntitlement(userId: string): Promise<Entitlement | null> {
  try {
    // Agregar timestamp para evitar cache y asegurar obtener la versión más reciente
    const timestamp = Date.now()
    const entitlementUrl = STORAGE_BASE_URL + '/' + userId + '/entitlement.json?t=' + timestamp
    const response = await fetch(entitlementUrl, {
      cache: 'no-store' // Evitar cache del navegador
    })
    
    if (!response.ok) {
      if (response.status === 404) {
        return null
      }
      console.error('Error al obtener entitlement.json: ' + response.status + ' ' + response.statusText)
      return null
    }
    
    const data = await response.json()
    return data
  } catch (error) {
    console.error('Error obteniendo entitlement.json:', error)
    return null
  }
}

export function calculateTrialDaysRemaining(endsAt: string | null | undefined): number | null {
  if (!endsAt) {
    return null
  }
  
  try {
    const endDate = new Date(endsAt)
    const now = new Date()
    
    // Calcular la diferencia en milisegundos
    const diffMs = endDate.getTime() - now.getTime()
    
    // Convertir a días (redondear hacia arriba para mostrar días completos)
    const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
    
    // Si ya expiró, retornar 0
    return Math.max(0, diffDays)
  } catch (error) {
    console.error('Error calculando días restantes:', error)
    return null
  }
}

/**
 * Obtiene el slug activo para un menu_id
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @returns El slug activo o null si no existe
 */
export async function getMenuSlug(menuId: string, accessToken: string): Promise<string | null> {
  try {
    const { createClient } = await import('@supabase/supabase-js')
    const supabaseUrl = 'https://rbpdhapfcljecofrscnj.supabase.co'
    const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'
    
    // Crear cliente autenticado con el token
    const supabase = createClient(supabaseUrl, supabaseAnonKey, {
      global: {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      }
    })
    
    const { data, error } = await supabase
      .from('menu_slugs')
      .select('slug')
      .eq('menu_id', menuId)
      .eq('is_active', true)
      .single()

    if (error) {
      if (error.code === 'PGRST116') {
        // No se encontró ningún registro
        return null
      }
      console.error('Error obteniendo slug:', error)
      return null
    }

    return data?.slug || null
  } catch (error) {
    console.error('Error en getMenuSlug:', error)
    return null
  }
}

/**
 * Crea un nuevo slug para un menu_id
 * @param menuId - ID del menú
 * @param slug - Slug a crear
 * @param accessToken - Token de autenticación
 * @returns El slug creado o null si hubo error
 */
export async function createMenuSlug(menuId: string, slug: string, accessToken: string): Promise<string | null> {
  try {
    const { createClient } = await import('@supabase/supabase-js')
    const supabaseUrl = 'https://rbpdhapfcljecofrscnj.supabase.co'
    const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'
    
    // Crear cliente autenticado con el token
    const supabase = createClient(supabaseUrl, supabaseAnonKey, {
      global: {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      }
    })
    
    // Verificar si el slug ya existe
    const { data: existing, error: checkError } = await supabase
      .from('menu_slugs')
      .select('id, menu_id')
      .eq('slug', slug)
      .maybeSingle()

    if (checkError && checkError.code !== 'PGRST116') {
      console.error('Error verificando slug:', checkError)
      return null
    }

    if (existing) {
      // El slug ya existe
      if (existing.menu_id === menuId) {
        // Es el mismo menú, activar el slug si no está activo
        const { error: updateError } = await supabase
          .from('menu_slugs')
          .update({ is_active: true })
          .eq('slug', slug)
          .eq('menu_id', menuId)
        
        if (updateError) {
          console.error('Error activando slug:', updateError)
          return null
        }
        return slug
      } else {
        // El slug existe para otro menú
        throw new Error('SLUG_EXISTS')
      }
    }

    // Desactivar otros slugs del mismo menú
    await supabase
      .from('menu_slugs')
      .update({ is_active: false })
      .eq('menu_id', menuId)

    // Crear el nuevo slug
    const { data, error } = await supabase
      .from('menu_slugs')
      .insert({
        menu_id: menuId,
        slug: slug,
        is_active: true
      })
      .select('slug')
      .single()

    if (error) {
      console.error('Error creando slug:', error)
      return null
    }

    return data?.slug || null
  } catch (error) {
    if (error instanceof Error && error.message === 'SLUG_EXISTS') {
      throw error
    }
    console.error('Error en createMenuSlug:', error)
    return null
  }
}

/**
 * Genera una URL usando el slug
 * @param slug - Slug del menú
 * @param lang - Idioma opcional
 * @returns URL completa con el slug
 */
export function generateSlugUrl(slug: string, lang?: string): string {
  let url = `https://catalogo.micartapro.com/m/${slug}`
  
  // Agregar parámetro de idioma si se proporciona
  if (lang && ['ES', 'PT', 'EN'].includes(lang)) {
    url += `?lang=${lang}`
  }
  
  return url
}