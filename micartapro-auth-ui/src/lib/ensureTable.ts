import { supabase } from './supabase'

/**
 * Intenta crear la tabla user_menus si no existe.
 * NOTA: Esto requiere permisos de administrador en Supabase.
 * En producción, es mejor crear la tabla manualmente usando el SQL Editor.
 * 
 * @returns true si la tabla existe o se creó exitosamente, false en caso contrario
 */
export async function ensureUserMenusTable(): Promise<boolean> {
  try {
    // Intentar hacer una query simple para verificar si la tabla existe
    const { error: testError } = await supabase
      .from('user_menus')
      .select('user_id')
      .limit(0)

    // Si no hay error, la tabla existe
    if (!testError) {
      console.log('✅ Tabla user_menus ya existe')
      return true
    }

    // Si el error es que la tabla no existe, intentar crearla
    if (testError.code === '42P01' || testError.message?.includes('does not exist')) {
      console.warn('⚠️ Tabla user_menus no existe. Debes crearla manualmente usando el SQL en SUPABASE_SETUP.md')
      console.warn('⚠️ Ve a Supabase Dashboard > SQL Editor y ejecuta el SQL proporcionado')
      return false
    }

    // Otro tipo de error
    console.error('❌ Error verificando tabla:', testError)
    return false
  } catch (error) {
    console.error('❌ Error en ensureUserMenusTable:', error)
    return false
  }
}

