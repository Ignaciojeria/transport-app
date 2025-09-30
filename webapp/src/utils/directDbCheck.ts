/**
 * Utilidades para verificar directamente la base de datos
 * Esto nos ayuda a comparar con lo que devuelve Electric SQL
 */

/**
 * Verifica directamente si existe una cuenta en la base de datos
 * Usando el endpoint de registro para verificar si el email ya existe
 */
export const checkAccountDirectly = async (email: string): Promise<{
  exists: boolean
  message: string
  details?: any
}> => {
  try {
    console.log('🔍 Verificando cuenta directamente en la base de datos para:', email)
    
    // Usar el endpoint de registro para verificar si el email ya existe
    const response = await fetch('https://einar-main-f0820bc.d2.zuplo.dev/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email,
        organizationName: 'TEST_CHECK', // Nombre de prueba
        country: 'CL'
      })
    })

    const responseData = await response.json()
    console.log('🔍 Respuesta directa de la base de datos:', responseData)

    // Si el email ya existe, el servidor debería devolver un error
    if (!response.ok) {
      const errorMessage = responseData.message || responseData.error || 'Error desconocido'
      
      if (errorMessage.toLowerCase().includes('email') && 
          (errorMessage.toLowerCase().includes('exist') || 
           errorMessage.toLowerCase().includes('ya existe') ||
           errorMessage.toLowerCase().includes('already'))) {
        return {
          exists: true,
          message: 'Email ya existe en la base de datos',
          details: responseData
        }
      }
    }

    // Si la respuesta es exitosa, significa que el email no existía
    return {
      exists: false,
      message: 'Email no existe en la base de datos',
      details: responseData
    }

  } catch (error) {
    console.error('❌ Error al verificar cuenta directamente:', error)
    return {
      exists: false,
      message: 'Error al verificar cuenta directamente',
      details: error
    }
  }
}

/**
 * Compara los resultados de Electric SQL vs verificación directa
 */
export const compareElectricVsDirect = async (email: string, electricResult: any) => {
  console.log('🔄 Comparando Electric SQL vs verificación directa...')
  
  const directResult = await checkAccountDirectly(email)
  
  console.log('📊 Comparación de resultados:')
  console.log('  Electric SQL:', electricResult ? 'ENCONTRADO' : 'NO ENCONTRADO')
  console.log('  Verificación directa:', directResult.exists ? 'ENCONTRADO' : 'NO ENCONTRADO')
  
  if (electricResult && !directResult.exists) {
    console.warn('⚠️ INCONSISTENCIA: Electric SQL encuentra datos que no existen en la BD')
    console.warn('   Esto indica un problema de caché en Electric SQL')
  } else if (!electricResult && directResult.exists) {
    console.warn('⚠️ INCONSISTENCIA: La BD tiene datos que Electric SQL no encuentra')
    console.warn('   Esto indica un problema de sincronización')
  } else {
    console.log('✅ Los resultados son consistentes')
  }
  
  return {
    electric: electricResult ? 'FOUND' : 'NOT_FOUND',
    direct: directResult.exists ? 'FOUND' : 'NOT_FOUND',
    consistent: (!!electricResult) === directResult.exists,
    details: {
      electric: electricResult,
      direct: directResult
    }
  }
}
