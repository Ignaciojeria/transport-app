import { supabase } from './supabase'
import { API_BASE_URL } from './config'

const STORAGE_BASE_URL = "https://storage.googleapis.com/micartapro-menus"
const SUPABASE_URL = 'https://rbpdhapfcljecofrscnj.supabase.co'

// Cache de clientes de Supabase autenticados para evitar múltiples instancias
const authenticatedClientsCache = new Map<string, any>()

/**
 * Obtiene o crea un cliente de Supabase autenticado con el token proporcionado
 * Reutiliza clientes existentes para evitar múltiples instancias
 * @param accessToken - Token de autenticación
 * @returns Cliente de Supabase autenticado
 */
export async function getAuthenticatedSupabaseClient(accessToken: string) {
  // Si ya existe un cliente para este token, reutilizarlo
  if (authenticatedClientsCache.has(accessToken)) {
    return authenticatedClientsCache.get(accessToken)
  }
  
  // Crear nuevo cliente solo si no existe
  const { createClient } = await import('@supabase/supabase-js')
  const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'
  
  // Usar un storage key único para evitar conflictos con el cliente principal
  // El storage key debe ser único por cliente para evitar el warning de múltiples instancias
  // Usamos un hash del token para crear una clave única pero estable
  const tokenHash = accessToken.substring(0, 16).replace(/[^a-zA-Z0-9]/g, '')
  const uniqueStorageKey = `sb-auth-${tokenHash}`
  
  // Crear un storage personalizado que use una clave única para aislar este cliente
  // del cliente principal de Supabase, evitando conflictos de storage
  const customStorage = {
    getItem: (key: string) => {
      // Usar una clave única para este cliente
      return localStorage.getItem(`${uniqueStorageKey}-${key}`)
    },
    setItem: (key: string, value: string) => {
      localStorage.setItem(`${uniqueStorageKey}-${key}`, value)
    },
    removeItem: (key: string) => {
      localStorage.removeItem(`${uniqueStorageKey}-${key}`)
    }
  }
  
  const client = createClient(SUPABASE_URL, supabaseAnonKey, {
    global: {
      headers: {
        Authorization: `Bearer ${accessToken}`
      }
    },
    auth: {
      autoRefreshToken: false, // Desactivar auto-refresh para evitar conflictos
      persistSession: false, // No persistir sesión para clientes autenticados manualmente
      detectSessionInUrl: false, // No detectar sesión en URL
      storage: customStorage // Storage personalizado con clave única para evitar conflictos
    }
  })
  
  // Guardar en cache (limitar a 10 clientes para evitar memory leak)
  if (authenticatedClientsCache.size >= 10) {
    // Eliminar el más antiguo (FIFO)
    const firstKey = authenticatedClientsCache.keys().next().value
    if (firstKey) {
      authenticatedClientsCache.delete(firstKey)
    }
  }
  
  authenticatedClientsCache.set(accessToken, client)
  return client
}

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
 * @param userId - ID del usuario
 * @param accessToken - Token de autenticación
 * @returns El menuID o null si no existe
 */
export async function getLatestMenuId(userId: string, accessToken: string): Promise<string | null> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    const { data, error } = await supabase
      .from('user_menus')
      .select('menu_id, created_at')
      .eq('user_id', userId)
      .order('created_at', { ascending: false })
      .limit(1)
      .single()

    if (error) {
      // PRIMERO verificar si es un error de permisos/RLS (406, PGRST301, etc.)
      // Estos errores tienen prioridad sobre "no encontrado" porque indican un problema de configuración
      const isPermissionError = 
        error.code === 'PGRST301' || 
        error.code === '42P01' || 
        error.code === '42501' ||
        error.message?.includes('406') || 
        error.message?.includes('permission') || 
        error.message?.includes('denied') || 
        error.message?.includes('does not exist') ||
        error.message?.includes('row-level security') ||
        error.details?.includes('406') ||
        error.hint?.includes('RLS')
      
      if (isPermissionError) {
        console.error('⚠️ Error de permisos o tabla no existe:')
        console.error('   - La tabla user_menus no existe, o')
        console.error('   - Las políticas RLS no están configuradas correctamente')
        console.error('📋 Solución: Ejecuta el archivo SUPABASE_SETUP_COMPLETE.sql en el SQL Editor de Supabase')
        console.error('   Este archivo crea la tabla user_menus y configura las políticas RLS necesarias')
        console.error('❌ Detalles del error:')
        console.error('   Código:', error.code)
        console.error('   Mensaje:', error.message)
        console.error('   Detalles:', error.details)
        console.error('   Hint:', error.hint)
        return null
      }
      
      // Si es PGRST116, significa que no se encontró ningún registro (normal si el usuario no tiene menú)
      if (error.code === 'PGRST116') {
        console.log('ℹ️ Usuario no tiene menú registrado aún')
        return null
      }
      
      // Otro tipo de error
      console.error('❌ Error obteniendo menuID:', error)
      console.error('   Código:', error.code)
      console.error('   Mensaje:', error.message)
      console.error('   Detalles:', error.details)
      console.error('   Hint:', error.hint)
      
      return null
    }

    return data?.menu_id || null
  } catch (err) {
    console.error('❌ Error en getLatestMenuId:', err)
    if (err instanceof Error) {
      console.error('   Mensaje:', err.message)
      console.error('   Stack:', err.stack)
    }
    return null
  }
}

/**
 * Genera la URL de la carta del restaurante usando slug
 * @param slug - Slug del menú
 * @param lang - Idioma opcional ('ES', 'PT', 'EN')
 * @param versionId - ID de la versión opcional. Si se proporciona, se usa esa versión específica (para interacciones)
 * @returns URL completa de la carta
 */
export function generateMenuUrlFromSlug(slug: string, lang?: string, versionId?: string): string {
  const baseUrl = typeof window !== 'undefined' 
    ? window.location.origin.includes('localhost')
      ? 'http://localhost:5173'
      : 'https://catalogo.micartapro.com'
    : 'https://catalogo.micartapro.com'
  
  let url = `${baseUrl}/m/${slug}`
  
  // Construir query parameters
  const params = new URLSearchParams()
  
  // Agregar version_id si se proporciona (para interacciones)
  if (versionId) {
    params.append('version_id', versionId)
  }
  
  // Agregar idioma si se proporciona
  if (lang && ['ES', 'PT', 'EN'].includes(lang)) {
    params.append('lang', lang)
  }
  
  // Agregar query string si hay parámetros
  const queryString = params.toString()
  if (queryString) {
    url += `?${queryString}`
  }
  
  return url
}

/**
 * Genera la URL del endpoint del backend para obtener el menú
 * @param menuId - ID del menú
 * @param versionId - ID de la versión opcional. Si se proporciona, se usa esa versión específica
 * @returns URL completa del endpoint del backend
 */
export function generateBackendMenuUrl(menuId: string, versionId?: string): string {
  let url = `${API_BASE_URL}/menu/${encodeURIComponent(menuId)}`
  
  // Construir query parameters
  const params = new URLSearchParams()
  
  // Agregar version_id si se proporciona
  if (versionId) {
    params.append('version_id', versionId)
  }
  
  // Agregar query string si hay parámetros
  const queryString = params.toString()
  if (queryString) {
    url += `?${queryString}`
  }
  
  return url
}

/**
 * Genera la URL de la carta del restaurante usando menu_id directamente
 * @param menuId - ID del menú
 * @param lang - Idioma opcional ('ES', 'PT', 'EN')
 * @param versionId - ID de la versión opcional. Si se proporciona, se usa esa versión específica (para interacciones)
 * @returns URL completa de la carta usando menu_id
 */
export function generateMenuUrlFromMenuId(menuId: string, lang?: string, versionId?: string): string {
  const baseUrl = typeof window !== 'undefined' 
    ? window.location.origin.includes('localhost')
      ? 'http://localhost:5173'
      : 'https://catalogo.micartapro.com'
    : 'https://catalogo.micartapro.com'
  
  let url = `${baseUrl}/m/${menuId}`
  
  // Construir query parameters
  const params = new URLSearchParams()
  
  // Agregar version_id si se proporciona (para interacciones)
  if (versionId) {
    params.append('version_id', versionId)
  }
  
  // Agregar idioma si se proporciona
  if (lang && ['ES', 'PT', 'EN'].includes(lang)) {
    params.append('lang', lang)
  }
  
  // Agregar query string si hay parámetros
  const queryString = params.toString()
  if (queryString) {
    url += `?${queryString}`
  }
  
  return url
}

/**
 * Genera la URL de la carta del restaurante (obtiene el slug automáticamente)
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @param lang - Idioma opcional ('ES', 'PT', 'EN')
 * @param versionId - ID de la versión opcional. Si se proporciona, se usa esa versión específica (para interacciones)
 * @returns Promise con la URL completa de la carta. Si no hay slug, usa menu_id directamente
 */
export async function generateMenuUrl(
  menuId: string, 
  accessToken: string, 
  lang?: string, 
  versionId?: string
): Promise<string | null> {
  try {
    // Intentar obtener el slug del menú
    const slug = await getMenuSlug(menuId, accessToken)
    
    if (slug) {
      // Si hay slug, generar URL usando el slug
      return generateMenuUrlFromSlug(slug, lang, versionId)
    }
    
    // Si no hay slug (normal en menús nuevos), usar menu_id directamente
    // No es un error, es esperado cuando se crea un menú por primera vez
    return generateMenuUrlFromMenuId(menuId, lang, versionId)
  } catch (error: any) {
    // Si hay un error (como 406 cuando no hay slug o políticas RLS), usar menu_id directamente
    // Esto es normal cuando se crea un menú por primera vez
    return generateMenuUrlFromMenuId(menuId, lang, versionId)
  }
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
 * Hace polling en Supabase hasta que la versión del menú exista
 * @param versionID - ID de la versión del menú
 * @param accessToken - Token de autenticación
 * @param maxAttempts - Número máximo de intentos (default: 120 - ~120 segundos)
 * @param intervalMs - Intervalo entre intentos en milisegundos (default: 1000)
 * @returns El contenido del menú cuando la versión existe
 */
export async function pollUntilVersionExists(
  versionID: string,
  accessToken: string,
  maxAttempts: number = 120,
  intervalMs: number = 1000
): Promise<any> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    let attempts = 0
    
    while (attempts < maxAttempts) {
      try {
        const { data, error } = await supabase
          .from('menu_versions')
          .select('content')
          .eq('id', versionID)
          .single()
        
        if (!error && data && data.content) {
          // Versión encontrada, retornar el contenido
          return data.content
        }
        
        // Si es error de "no encontrado", continuar polling
        if (error?.code === 'PGRST116') {
          await new Promise(resolve => setTimeout(resolve, intervalMs))
          attempts++
          continue
        }
        
        // Otro error, loguear y lanzar
        if (error) {
          console.error('Error de Supabase:', error)
          // Si es un error 406, probablemente es un problema de RLS
          // El 406 puede venir en el mensaje o en el código de error
          const errorMessage = error.message || ''
          const errorCode = error.code || ''
          const errorDetails = error.details || ''
          
          if (errorMessage.includes('406') || errorCode.includes('406') || errorDetails.includes('406')) {
            throw new Error(
              `❌ Error de permisos (RLS): La tabla menu_versions requiere políticas RLS.\n\n` +
              `📋 Solución: Ejecuta el SQL en RLS_POLICIES_MENU_VERSIONS.sql en el SQL Editor de Supabase.\n\n` +
              `Detalle técnico: ${errorMessage || errorDetails || errorCode}`
            )
          }
          
          // Si es un error de permisos genérico
          if (errorMessage.includes('permission') || errorMessage.includes('denied') || errorCode === 'PGRST301') {
            throw new Error(
              `❌ Error de permisos: Verifica que existan políticas RLS para menu_versions.\n\n` +
              `Ejecuta: RLS_POLICIES_MENU_VERSIONS.sql en Supabase SQL Editor.\n\n` +
              `Detalle: ${errorMessage}`
            )
          }
          
          throw new Error(`Error consultando versión: ${errorMessage || errorCode || 'Error desconocido'}`)
        }
        
        // Esperar antes del siguiente intento
        await new Promise(resolve => setTimeout(resolve, intervalMs))
        attempts++
      } catch (error: any) {
        // Si es un error de "no encontrado", continuar intentando
        if (error?.message?.includes('PGRST116') || error?.code === 'PGRST116') {
          await new Promise(resolve => setTimeout(resolve, intervalMs))
          attempts++
          continue
        }
        
        console.error(`Error en intento ${attempts + 1}:`, error)
        await new Promise(resolve => setTimeout(resolve, intervalMs))
        attempts++
      }
    }
    
    throw new Error(`Timeout: No se pudo obtener la versión ${versionID} después de ${maxAttempts} intentos`)
  } catch (error) {
    console.error('Error creando cliente de Supabase:', error)
    throw error
  }
}

/**
 * Hace polling al endpoint del backend para verificar que el menú esté listo
 * Útil cuando el usuario se registra por primera vez y el menú se está creando
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @param versionId - ID de la versión opcional para verificar
 * @param maxAttempts - Número máximo de intentos (default: 30)
 * @param intervalMs - Intervalo entre intentos en milisegundos (default: 2000)
 * @returns true si el menú está listo, false si no se encontró después de todos los intentos
 */
export async function pollUntilMenuExists(
  menuId: string,
  accessToken: string,
  versionId?: string,
  maxAttempts: number = 30,
  intervalMs: number = 2000
): Promise<boolean> {
  // Importar API_BASE_URL desde config
  const { API_BASE_URL } = await import('./config')
  let attempts = 0
  
  while (attempts < maxAttempts) {
    try {
      // Construir URL con versionId si se proporciona
      let url = `${API_BASE_URL}/menu/${menuId}`
      if (versionId) {
        url += `?version_id=${encodeURIComponent(versionId)}`
      }
      
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        },
        credentials: 'include'
      })
      
      if (response.ok) {
        // El menú está listo
        console.log('✅ Menú listo en el backend', versionId ? `(versión ${versionId})` : '')
        return true
      }
      
      // Si es 404, el menú aún no existe, continuar intentando
      if (response.status === 404) {
        console.log(`⏳ Intento ${attempts + 1}/${maxAttempts}: Menú aún no está listo (404)`)
        await new Promise(resolve => setTimeout(resolve, intervalMs))
        attempts++
        continue
      }
      
      // Otro error HTTP, loguear pero continuar intentando
      console.warn(`⚠️ Intento ${attempts + 1}/${maxAttempts}: Error ${response.status}`)
      await new Promise(resolve => setTimeout(resolve, intervalMs))
      attempts++
    } catch (error: any) {
      // Error de red, loguear pero continuar intentando
      console.error(`❌ Error en intento ${attempts + 1}:`, error)
      await new Promise(resolve => setTimeout(resolve, intervalMs))
      attempts++
    }
  }
  
  // No se encontró el menú después de todos los intentos
  console.warn(`⏱️ Timeout: El menú no estuvo listo después de ${maxAttempts} intentos`)
  return false
}

/**
 * Obtiene el slug activo para un menu_id
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @returns El slug activo o null si no existe
 */
export async function getMenuSlug(menuId: string, accessToken: string): Promise<string | null> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    const { data, error } = await supabase
      .from('menu_slugs')
      .select('slug')
      .eq('menu_id', menuId)
      .eq('is_active', true)
      .maybeSingle() // Usar maybeSingle() en lugar de single() para manejar cuando no hay slug

    if (error) {
      // Manejar errores comunes cuando no hay slug (menú nuevo)
      if (error.code === 'PGRST116' || error.code === '406' || error.message?.includes('406')) {
        // No se encontró ningún registro o error 406 (Not Acceptable) - es normal cuando no hay slug
        return null
      }
      console.error('Error obteniendo slug:', error)
      return null
    }

    // Si no hay data, significa que no existe slug para este menú (normal en menús nuevos)
    return data?.slug || null
  } catch (error: any) {
    // Manejar errores 406 específicamente (Not Acceptable - común cuando no hay políticas RLS o no hay slug)
    if (error?.code === '406' || error?.message?.includes('406') || error?.status === 406) {
      // Es normal que no haya slug cuando se crea un menú por primera vez
      return null
    }
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
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
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

/**
 * Obtiene todas las versiones de un menú ordenadas por número de versión (descendente)
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @returns Array de versiones con id, version_number, created_at, name, is_favorite
 */
export async function getMenuVersions(
  menuId: string,
  accessToken: string
): Promise<Array<{ id: string; version_number: number; created_at: string; name: string | null; is_favorite: boolean; content?: any }>> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    const { data, error } = await supabase
      .from('menu_versions')
      .select('id, version_number, created_at, name, is_favorite, content')
      .eq('menu_id', menuId)
      .order('version_number', { ascending: false })
    
    if (error) {
      console.error('Error obteniendo versiones del menú:', error)
      return []
    }
    
    return data || []
  } catch (error) {
    console.error('Error en getMenuVersions:', error)
    return []
  }
}

export interface MenuOrderRow {
  order_number: number
  event_payload: Record<string, unknown>
  event_type: string
  requested_time: string | null
  created_at: string | null
}

/** Fila de la proyección order_items_projection (un ítem por fila). */
export interface OrderItemProjectionRow {
  id: number
  menu_order_id: number
  order_number: number
  menu_id: string
  item_key: string
  item_name: string
  quantity: number
  unit: string
  station: string | null
  fulfillment: string
  status: string
  requested_time: string | null
  created_at: string
  updated_at: string
}

/** Ítem agrupado para la vista de cocina. */
export interface KitchenOrderItem {
  item_name: string
  quantity: number
  unit: string
  station: string | null
}

/** Orden ya agrupada para la vista de cocina (desde order_items_projection). */
export interface KitchenOrder {
  order_number: number
  menu_order_id: number
  requested_time: string | null
  created_at: string
  fulfillment: string
  items: KitchenOrderItem[]
}

/** Filtro por estación: ALL = todo; KITCHEN/BAR filtra por columna station en Supabase. */
export type StationFilter = 'ALL' | 'KITCHEN' | 'BAR'

/**
 * Obtiene las órdenes para la vista de cocina desde la proyección order_items_projection.
 * El filtro por estación se aplica en Supabase (columna station).
 * Agrupa ítems por orden (order_number, menu_order_id) y por item_name (sumando cantidades).
 * Orden: requested_time ASC, created_at ASC.
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @param stationFilter - Si es KITCHEN o BAR, filtra en la query por columna station
 * @returns Array de órdenes listas para mostrar en cocina
 */
export async function getKitchenOrdersFromProjection(
  menuId: string,
  accessToken: string,
  stationFilter: StationFilter = 'ALL'
): Promise<KitchenOrder[]> {
  try {
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    let query = supabase
      .from('order_items_projection')
      .select('order_number, menu_order_id, requested_time, created_at, fulfillment, item_name, quantity, unit, station')
      .eq('menu_id', menuId)
    if (stationFilter === 'KITCHEN' || stationFilter === 'BAR') {
      query = query.eq('station', stationFilter)
    }
    const { data: rows, error } = await query
      .order('requested_time', { ascending: true, nullsFirst: false })
      .order('created_at', { ascending: true, nullsFirst: false })

    if (error) {
      console.error('Error obteniendo proyección de órdenes:', error)
      return []
    }
    const items = (rows || []) as Array<{
      order_number: number
      menu_order_id: number
      requested_time: string | null
      created_at: string
      fulfillment: string
      item_name: string
      quantity: number
      unit: string
      station: string | null
    }>
    return groupProjectionItemsByOrder(items)
  } catch (error) {
    console.error('Error en getKitchenOrdersFromProjection:', error)
    return []
  }
}

function groupProjectionItemsByOrder(
  rows: Array<{
    order_number: number
    menu_order_id: number
    requested_time: string | null
    created_at: string
    fulfillment: string
    item_name: string
    quantity: number
    unit: string
    station: string | null
  }>
): KitchenOrder[] {
  const byOrder = new Map<string, KitchenOrder>()
  for (const r of rows) {
    const key = `${r.menu_order_id}:${r.order_number}`
    let order = byOrder.get(key)
    if (!order) {
      order = {
        order_number: r.order_number,
        menu_order_id: r.menu_order_id,
        requested_time: r.requested_time,
        created_at: r.created_at,
        fulfillment: r.fulfillment,
        items: []
      }
      byOrder.set(key, order)
    }
    const name = r.item_name?.trim() || '—'
    const existing = order.items.find(
      (i) => i.item_name === name && i.unit === r.unit && (i.station ?? '') === (r.station ?? '')
    )
    if (existing) {
      existing.quantity += r.quantity
    } else {
      order.items.push({
        item_name: name,
        quantity: r.quantity,
        unit: r.unit || 'EACH',
        station: r.station
      })
    }
  }
  return [...byOrder.values()].sort((a, b) => {
    const ta = a.requested_time ? new Date(a.requested_time).getTime() : 0
    const tb = b.requested_time ? new Date(b.requested_time).getTime() : 0
    if (ta !== tb) return ta - tb
    return new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  })
}

/**
 * Obtiene una orden completa desde menu_orders (para modal "Ver como hoja").
 * @param menuId - ID del menú
 * @param orderNumber - Número de orden
 * @param accessToken - Token de autenticación
 */
export async function getMenuOrderByNumber(
  menuId: string,
  orderNumber: number,
  accessToken: string
): Promise<MenuOrderRow | null> {
  try {
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    const { data, error } = await supabase
      .from('menu_orders')
      .select('order_number, event_payload, event_type, requested_time, created_at')
      .eq('menu_id', menuId)
      .eq('order_number', orderNumber)
      .maybeSingle()
    if (error || !data) return null
    return data as MenuOrderRow
  } catch {
    return null
  }
}

/**
 * Obtiene las órdenes del menú (DELIVERY o PICKUP) desde la tabla de eventos.
 * Usar getKitchenOrdersFromProjection para la vista de cocina.
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @returns Array de órdenes con order_number, event_payload, event_type, requested_time, created_at
 */
export async function getMenuOrders(
  menuId: string,
  accessToken: string
): Promise<MenuOrderRow[]> {
  try {
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    const { data, error } = await supabase
      .from('menu_orders')
      .select('order_number, event_payload, event_type, requested_time, created_at')
      .eq('menu_id', menuId)
      .order('requested_time', { ascending: true, nullsFirst: false })
      .order('created_at', { ascending: true, nullsFirst: false })

    if (error) {
      console.error('Error obteniendo órdenes del menú:', error)
      return []
    }
    return (data || []) as MenuOrderRow[]
  } catch (error) {
    console.error('Error en getMenuOrders:', error)
    return []
  }
}

/**
 * Suscripción Realtime a nuevos pedidos en menu_orders para un menu_id.
 * Cuando se inserta una fila, se llama onInsert (ej. recargar lista).
 * Requiere que la tabla esté en la publicación Realtime (ver sql/realtime_menu_orders.sql).
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @param onInsert - Callback al insertar una nueva orden
 * @returns Función para cancelar la suscripción
 */
export async function subscribeMenuOrdersRealtime(
  menuId: string,
  accessToken: string,
  onInsert: () => void
): Promise<() => void> {
  const supabase = await getAuthenticatedSupabaseClient(accessToken)
  const channel = supabase
    .channel(`menu_orders:${menuId}`)
    .on(
      'postgres_changes',
      {
        event: 'INSERT',
        schema: 'public',
        table: 'menu_orders',
        filter: `menu_id=eq.${menuId}`
      },
      () => {
        onInsert()
      }
    )
    .subscribe((status: string, err?: Error) => {
      if (import.meta.env.DEV) {
        console.log('[Cocina Realtime]', status, err ?? '')
      }
      if (status === 'CHANNEL_ERROR' && import.meta.env.DEV) {
        console.warn('[Cocina Realtime] Si no ves pedidos al instante, ejecuta en Supabase SQL: ALTER PUBLICATION supabase_realtime ADD TABLE public.menu_orders;')
      }
    })
  return () => {
    supabase.removeChannel(channel)
  }
}

/**
 * Actualiza el nombre de una versión del menú
 * @param versionId - ID de la versión
 * @param name - Nuevo nombre para la versión
 * @param accessToken - Token de autenticación
 * @returns true si se actualizó correctamente, false si hubo error
 */
export async function updateVersionName(
  versionId: string,
  name: string,
  accessToken: string
): Promise<boolean> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    console.log('🔄 Actualizando nombre de versión:', { versionId, name: name.trim() || null })
    
    const { data, error } = await supabase
      .from('menu_versions')
      .update({ name: name.trim() || null })
      .eq('id', versionId)
      .select('id, name')
    
    if (error) {
      console.error('❌ Error actualizando nombre de versión:', error)
      console.error('❌ Detalles del error:', {
        code: error.code,
        message: error.message,
        details: error.details,
        hint: error.hint
      })
      
      // Si es un error 406 o PGRST301, probablemente es un problema de RLS
      if (error.code === 'PGRST301' || error.message?.includes('406') || error.message?.includes('permission') || error.message?.includes('denied')) {
        console.error('⚠️ Error de permisos (RLS): La tabla menu_versions requiere una política RLS para UPDATE.')
        console.error('📋 Solución: Ejecuta el SQL para agregar la política "Users can update their own menu versions" en el SQL Editor de Supabase.')
      }
      
      return false
    }
    
    console.log('✅ Nombre de versión actualizado correctamente:', data)
    return true
  } catch (error) {
    console.error('❌ Error en updateVersionName:', error)
    return false
  }
}

/**
 * Actualiza el estado de favorito de una versión del menú
 * @param versionId - ID de la versión
 * @param isFavorite - Nuevo estado de favorito
 * @param accessToken - Token de autenticación
 * @returns true si se actualizó correctamente, false si hubo error
 */
export async function updateVersionFavorite(
  versionId: string,
  isFavorite: boolean,
  accessToken: string
): Promise<boolean> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    console.log('🔄 Actualizando favorito de versión:', { versionId, isFavorite })
    
    const { data, error } = await supabase
      .from('menu_versions')
      .update({ is_favorite: isFavorite })
      .eq('id', versionId)
      .select('id, is_favorite')
    
    if (error) {
      console.error('❌ Error actualizando favorito de versión:', error)
      console.error('❌ Detalles del error:', {
        code: error.code,
        message: error.message,
        details: error.details,
        hint: error.hint
      })
      
      // Si es un error 406 o PGRST301, probablemente es un problema de RLS
      if (error.code === 'PGRST301' || error.message?.includes('406') || error.message?.includes('permission') || error.message?.includes('denied')) {
        console.error('⚠️ Error de permisos (RLS): La tabla menu_versions requiere una política RLS para UPDATE.')
        console.error('📋 Solución: Ejecuta el SQL para agregar la política "Users can update their own menu versions" en el SQL Editor de Supabase.')
      }
      
      return false
    }
    
    console.log('✅ Favorito de versión actualizado correctamente:', data)
    return true
  } catch (error) {
    console.error('❌ Error en updateVersionFavorite:', error)
    return false
  }
}

/**
 * Obtiene el current_version_id de un menú
 * @param menuId - ID del menú
 * @param accessToken - Token de autenticación
 * @returns ID de la versión actual o null
 */
export async function getCurrentVersionId(
  menuId: string,
  accessToken: string
): Promise<string | null> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    const { data, error } = await supabase
      .from('menus')
      .select('current_version_id')
      .eq('id', menuId)
      .single()
    
    if (error) {
      console.error('Error obteniendo current_version_id:', error)
      return null
    }
    
    return data?.current_version_id || null
  } catch (error) {
    console.error('Error en getCurrentVersionId:', error)
    return null
  }
}

/**
 * Actualiza el current_version_id de un menú en Supabase
 * @param menuId - ID del menú
 * @param versionId - ID de la versión a activar
 * @param accessToken - Token de autenticación
 * @returns true si se actualizó correctamente, false si hubo error
 */
export async function updateCurrentVersionId(
  menuId: string,
  versionId: string,
  accessToken: string
): Promise<boolean> {
  try {
    // Usar cliente autenticado reutilizable
    const supabase = await getAuthenticatedSupabaseClient(accessToken)
    
    console.log('🔄 Actualizando current_version_id en tabla menus (Supabase):', { menuId, versionId })
    console.log('🔄 URL Supabase:', SUPABASE_URL)
    console.log('🔄 Tabla: menus')
    console.log('🔄 Operación: UPDATE current_version_id')
    
    // Verificar primero que el menú existe y pertenece al usuario
    const { data: menuData, error: checkError } = await supabase
      .from('menus')
      .select('id, user_id, current_version_id')
      .eq('id', menuId)
      .single()
    
    if (checkError) {
      console.error('❌ Error verificando menú antes de actualizar:', checkError)
      
      // Si es un error 406 o PGRST301, probablemente es un problema de RLS para SELECT
      if (checkError.code === 'PGRST301' || checkError.message?.includes('406') || checkError.message?.includes('permission') || checkError.message?.includes('denied')) {
        console.error('⚠️ Error de permisos (RLS): La tabla menus requiere una política RLS para SELECT.')
        console.error('📋 Solución: Ejecuta el SQL en RLS_POLICIES_MENUS_UPDATE.sql en el SQL Editor de Supabase.')
        console.error('📋 El SQL creará las políticas: "Users can view their own menus" (SELECT) y "Users can update their own menu current_version_id" (UPDATE)')
      }
      
      // Si es PGRST116, el menú no existe o no tiene permisos para verlo
      if (checkError.code === 'PGRST116') {
        console.error('⚠️ Menú no encontrado o sin permisos para verlo. Verifica:')
        console.error('   1. Que el menuId sea correcto:', menuId)
        console.error('   2. Que exista una política RLS para SELECT en la tabla menus')
        console.error('   3. Que el menú pertenezca al usuario autenticado')
      }
      
      return false
    }
    
    if (!menuData) {
      console.error('❌ Menú no encontrado:', menuId)
      return false
    }
    
    console.log('✅ Menú encontrado:', menuData)
    console.log('🔄 current_version_id actual:', menuData.current_version_id)
    console.log('🔄 current_version_id nuevo:', versionId)
    
    // Actualizar current_version_id en la tabla menus
    const { data, error } = await supabase
      .from('menus')
      .update({ current_version_id: versionId })
      .eq('id', menuId)
      .select('id, current_version_id')

    if (error) {
      console.error('❌ Error actualizando current_version_id en tabla menus:', error)
      console.error('❌ Detalles del error:', {
        code: error.code,
        message: error.message,
        details: error.details,
        hint: error.hint
      })
      
      // Si es un error 406 o PGRST301, probablemente es un problema de RLS
      if (error.code === 'PGRST301' || error.message?.includes('406') || error.message?.includes('permission') || error.message?.includes('denied')) {
        console.error('⚠️ Error de permisos (RLS): La tabla menus requiere una política RLS para UPDATE.')
        console.error('📋 Solución: Ejecuta el SQL en RLS_POLICIES_MENUS_UPDATE.sql en el SQL Editor de Supabase.')
        console.error('📋 El SQL creará la política: "Users can update their own menu current_version_id"')
      }
      
      return false
    }

    console.log('✅ current_version_id actualizado correctamente en tabla menus:', data)
    console.log('✅ Verificación: El menú ahora tiene current_version_id =', data[0]?.current_version_id)
    return true
  } catch (error) {
    console.error('❌ Error en updateCurrentVersionId:', error)
    return false
  }
}

/**
 * Obtiene los créditos del usuario
 * @param accessToken - Token de autenticación
 * @returns Promise con el balance de créditos y transacciones
 */
export async function getUserCredits(accessToken: string): Promise<{ balance: number; transactions: any[] } | null> {
  try {
    console.log('🔍 Obteniendo créditos del usuario')
    
    const response = await fetch(`${API_BASE_URL}/credits`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      console.error('❌ Error obteniendo créditos:', {
        status: response.status,
        statusText: response.statusText
      })
      return null
    }

    const data = await response.json()
    console.log('✅ Créditos obtenidos:', data.balance)
    return {
      balance: data.balance || 0,
      transactions: data.transactions || []
    }
  } catch (error: any) {
    console.error('❌ Error en getUserCredits:', error)
    return null
  }
}

/**
 * Verifica si el usuario tiene créditos suficientes
 * @param accessToken - Token de autenticación
 * @param requiredCredits - Créditos requeridos (por defecto 1)
 * @returns Promise<boolean> - true si tiene créditos suficientes, false en caso contrario
 */
export async function hasEnoughCredits(accessToken: string, requiredCredits: number = 1): Promise<boolean> {
  try {
    const credits = await getUserCredits(accessToken)
    if (!credits) {
      return false
    }
    return credits.balance >= requiredCredits
  } catch (error: any) {
    console.error('❌ Error en hasEnoughCredits:', error)
    return false
  }
}

/**
 * Verifica si el usuario tiene una suscripción activa
 * @param userId - ID del usuario
 * @param accessToken - Token de autenticación
 * @returns Promise<boolean> - true si tiene suscripción activa, false en caso contrario
 */
export async function hasActiveSubscription(userId: string, accessToken: string): Promise<boolean> {
  try {
    console.log('🔍 Verificando suscripción para usuario:', userId)
    
    // Llamar al endpoint del backend para verificar la suscripción
    const response = await fetch(`${API_BASE_URL}/check-subscription`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      console.error('❌ Error verificando suscripción:', {
        status: response.status,
        statusText: response.statusText
      })
      return false
    }

    const data = await response.json()
    const hasSubscription = data.has_active_subscription === true

    if (hasSubscription) {
      console.log('✅ Usuario tiene suscripción activa')
    } else {
      console.log('ℹ️ Usuario no tiene suscripción activa')
    }

    return hasSubscription
  } catch (error: any) {
    console.error('❌ Error en hasActiveSubscription:', error)
    return false
  }
}