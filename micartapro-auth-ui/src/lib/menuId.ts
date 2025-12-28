import { v7 as uuidv7 } from 'uuid'
import { supabase } from './supabase'

/**
 * Obtiene el menuID del usuario actual, o crea uno nuevo si no existe.
 * Esta función es idempotente: si el usuario ya tiene un menuID, lo retorna.
 * Si no tiene uno, genera un nuevo UUID v7 y lo guarda en Supabase.
 * 
 * @param userId - El ID del usuario de Supabase (UUID)
 * @returns El menuID (UUID v7) del usuario
 * @throws Error si no se puede obtener o crear el menuID
 */
export async function getOrCreateMenuId(userId: string): Promise<string> {
  try {
    // 1. Intentar obtener el menuID existente
    const { data: existingMenu, error: fetchError } = await supabase
      .from('user_menus')
      .select('menu_id')
      .eq('user_id', userId)
      .maybeSingle()

    // Si hay un error que no sea "no encontrado", lanzarlo
    if (fetchError && fetchError.code !== 'PGRST116') {
      throw new Error(`Error al obtener menuID: ${fetchError.message}`)
    }

    // Si existe, retornarlo
    if (existingMenu && existingMenu.menu_id) {
      console.log('✅ MenuID existente encontrado:', existingMenu.menu_id)
      return existingMenu.menu_id
    }

    // 2. Si no existe, generar un nuevo UUID v7
    const newMenuId = uuidv7()
    console.log('🆕 Generando nuevo MenuID (UUID v7):', newMenuId)

    // 3. Intentar insertar (con upsert para evitar race conditions)
    const { data: insertedMenu, error: insertError } = await supabase
      .from('user_menus')
      .upsert(
        {
          user_id: userId,
          menu_id: newMenuId,
          created_at: new Date().toISOString(),
        },
        {
          onConflict: 'user_id', // Si ya existe, actualizar
        }
      )
      .select('menu_id')
      .single()

    if (insertError) {
      // Si falla por conflicto (otro proceso ya creó el menuID), intentar obtenerlo de nuevo
      if (insertError.code === '23505' || insertError.message?.includes('duplicate')) {
        console.log('⚠️ Race condition detectada, obteniendo menuID existente...')
        const { data: retryMenu, error: retryError } = await supabase
          .from('user_menus')
          .select('menu_id')
          .eq('user_id', userId)
          .single()

        if (retryMenu && retryMenu.menu_id) {
          console.log('✅ MenuID obtenido después de race condition:', retryMenu.menu_id)
          return retryMenu.menu_id
        }

        if (retryError) {
          throw new Error(`Error al obtener menuID después de race condition: ${retryError.message}`)
        }
      }

      throw new Error(`Error al crear menuID: ${insertError.message}`)
    }

    if (!insertedMenu || !insertedMenu.menu_id) {
      throw new Error('No se pudo crear el menuID')
    }

    console.log('✅ MenuID creado exitosamente:', insertedMenu.menu_id)
    return insertedMenu.menu_id
  } catch (error) {
    console.error('❌ Error en getOrCreateMenuId:', error)
    throw error
  }
}

