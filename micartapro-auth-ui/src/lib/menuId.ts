import { v7 as uuidv7 } from 'uuid'
import { supabase } from './supabase'

// Cache para evitar llamadas simultáneas para el mismo usuario
const pendingRequests = new Map<string, Promise<string>>()

/**
 * Verifica si el usuario tiene al menos un menuID, y si no tiene ninguno, crea uno por defecto.
 * Esta función es idempotente: si el usuario ya tiene al menos un menuID, no crea uno nuevo.
 * Si no tiene ninguno, genera un nuevo UUID v7 y lo guarda en Supabase.
 * 
 * NOTA: El usuario puede tener múltiples menús, pero en el callback solo se crea uno por defecto
 * si no tiene ninguno.
 * 
 * Evita race conditions usando un cache de peticiones pendientes.
 * 
 * @param userId - El ID del usuario de Supabase (UUID)
 * @returns El menuID (UUID v7) creado o el más reciente si ya existe alguno
 * @throws Error si no se puede obtener o crear el menuID
 */
export async function getOrCreateMenuId(userId: string): Promise<string> {
  // Si ya hay una petición en curso para este usuario, esperar a que termine
  if (pendingRequests.has(userId)) {
    console.log('⏳ Ya hay una petición en curso para este usuario, esperando...')
    return pendingRequests.get(userId)!
  }

  // Crear la petición y guardarla en el cache
  const requestPromise = getOrCreateMenuIdInternal(userId)
  pendingRequests.set(userId, requestPromise)

  try {
    const result = await requestPromise
    return result
  } finally {
    // Limpiar el cache después de completar
    pendingRequests.delete(userId)
  }
}

async function getOrCreateMenuIdInternal(userId: string): Promise<string> {
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

    // Si el usuario ya tiene al menos un menuID, no crear uno nuevo
    if (existingMenus && existingMenus.length > 0 && existingMenus[0].menu_id) {
      const latestMenuId = existingMenus[0].menu_id
      console.log('✅ El usuario ya tiene al menos un menuID, usando el más reciente:', latestMenuId)
      return latestMenuId
    }

    // Si llegamos aquí, el usuario NO tiene ningún menuID, así que creamos uno por defecto
    console.log('📝 El usuario no tiene ningún menuID, creando uno por defecto...')

    // 2. Si el usuario no tiene ningún menuID, generar uno nuevo por defecto
    const newMenuId = uuidv7()
    console.log('🆕 Generando nuevo MenuID por defecto (UUID v7):', newMenuId)

    // 3. Insertar el nuevo menuID (INSERT simple, no UPSERT)
    // No usamos onConflict porque el usuario puede tener múltiples menús
    const { data: insertedMenu, error: insertError } = await supabase
      .from('user_menus')
      .insert({
        user_id: userId,
        menu_id: newMenuId,
        created_at: new Date().toISOString(),
      })
      .select('menu_id')
      .single()

    if (insertError) {
      console.error('❌ Error al insertar menuID:', {
        code: insertError.code,
        message: insertError.message,
        details: insertError.details,
        hint: insertError.hint
      })
      
      // Si falla por duplicado (race condition), intentar obtener el menuID que se creó
      if (insertError.code === '23505' || insertError.message?.includes('duplicate') || insertError.message?.includes('unique')) {
        console.log('⚠️ MenuID duplicado detectado (race condition), obteniendo el existente...')
        const { data: retryMenus, error: retryError } = await supabase
          .from('user_menus')
          .select('menu_id')
          .eq('user_id', userId)
          .order('created_at', { ascending: false })
          .limit(1)

        if (retryMenus && retryMenus.length > 0 && retryMenus[0].menu_id) {
          console.log('✅ MenuID obtenido después de race condition:', retryMenus[0].menu_id)
          return retryMenus[0].menu_id
        }

        if (retryError) {
          throw new Error(`Error al obtener menuID después de race condition: ${retryError.message}`)
        }
      }
      
      // Si es error de RLS
      if (insertError.code === '42501' || insertError.message?.includes('permission denied') || insertError.message?.includes('RLS')) {
        throw new Error('Permiso denegado por RLS al insertar. Verifica que la política "Users can insert their own menu" esté correctamente configurada.')
      }

      throw new Error(`Error al crear menuID: ${insertError.message} (Código: ${insertError.code})`)
    }

    if (!insertedMenu || !insertedMenu.menu_id) {
      throw new Error('No se pudo crear el menuID: respuesta vacía')
    }

    console.log('✅ MenuID por defecto creado exitosamente:', insertedMenu.menu_id)
    return insertedMenu.menu_id
  } catch (error) {
    console.error('❌ Error completo en getOrCreateMenuId:', error)
    throw error
  }
}

