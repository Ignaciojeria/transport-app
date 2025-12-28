import { supabase } from './supabase'

/**
 * Función de diagnóstico para verificar la configuración de Supabase
 * Úsala en la consola del navegador para debuggear problemas con menuID
 */
export async function diagnoseMenuIdIssue() {
  console.log('🔍 Iniciando diagnóstico de menuID...\n')
  
  // 1. Verificar sesión
  console.log('1️⃣ Verificando sesión...')
  const { data: { session }, error: sessionError } = await supabase.auth.getSession()
  
  if (sessionError) {
    console.error('❌ Error obteniendo sesión:', sessionError)
    return
  }
  
  if (!session || !session.user) {
    console.error('❌ No hay sesión activa')
    return
  }
  
  console.log('✅ Sesión activa:', {
    userId: session.user.id,
    email: session.user.email
  })
  
  // 2. Verificar que la tabla existe
  console.log('\n2️⃣ Verificando que la tabla user_menus existe...')
  const { data: tableTest, error: tableError } = await supabase
    .from('user_menus')
    .select('user_id')
    .limit(0)
  
  if (tableError) {
    if (tableError.code === '42P01') {
      console.error('❌ La tabla user_menus NO existe')
      console.error('   → Ejecuta el SQL en SUPABASE_SETUP.md')
    } else {
      console.error('❌ Error accediendo a la tabla:', tableError)
    }
    return
  }
  
  console.log('✅ La tabla user_menus existe')
  
  // 3. Verificar permisos de lectura
  console.log('\n3️⃣ Verificando permisos de lectura...')
  const { data: readTest, error: readError } = await supabase
    .from('user_menus')
    .select('*')
    .eq('user_id', session.user.id)
    .maybeSingle()
  
  if (readError) {
    if (readError.code === '42501' || readError.message?.includes('permission denied')) {
      console.error('❌ Permiso denegado para leer (problema de RLS)')
      console.error('   → Verifica las políticas RLS en Supabase')
      console.error('   → Asegúrate de que la política "Users can view their own menu" esté activa')
    } else {
      console.error('❌ Error al leer:', readError)
    }
  } else {
    console.log('✅ Permisos de lectura OK')
    if (readTest) {
      console.log('   → MenuID existente:', readTest.menu_id)
    } else {
      console.log('   → No hay menuID para este usuario')
    }
  }
  
  // 4. Verificar permisos de escritura
  console.log('\n4️⃣ Verificando permisos de escritura...')
  const testMenuId = `test-${Date.now()}`
  const { data: writeTest, error: writeError } = await supabase
    .from('user_menus')
    .insert({
      user_id: session.user.id,
      menu_id: testMenuId,
    })
    .select()
    .single()
  
  if (writeError) {
    if (writeError.code === '42501' || writeError.message?.includes('permission denied')) {
      console.error('❌ Permiso denegado para escribir (problema de RLS)')
      console.error('   → Verifica las políticas RLS en Supabase')
      console.error('   → Asegúrate de que la política "Users can insert their own menu" esté activa')
    } else if (writeError.code === '23505' || writeError.message?.includes('duplicate')) {
      console.log('⚠️ MenuID de prueba ya existe (esto es normal si ya hay un registro)')
      // Intentar eliminar el registro de prueba
      await supabase
        .from('user_menus')
        .delete()
        .eq('menu_id', testMenuId)
    } else {
      console.error('❌ Error al escribir:', writeError)
    }
  } else {
    console.log('✅ Permisos de escritura OK')
    // Eliminar el registro de prueba
    await supabase
      .from('user_menus')
      .delete()
      .eq('menu_id', testMenuId)
    console.log('   → Registro de prueba eliminado')
  }
  
  console.log('\n✅ Diagnóstico completado')
}

// Exportar también para uso en consola del navegador
if (typeof window !== 'undefined') {
  (window as any).diagnoseMenuIdIssue = diagnoseMenuIdIssue
}

