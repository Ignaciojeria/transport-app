/**
 * Utilidades para verificar directamente la base de datos
 * Esto nos ayuda a comparar con lo que devuelve Electric SQL
 */

/**
 * Verifica directamente si existe una cuenta en la base de datos
 * Usando el endpoint de Electric SQL con un token válido
 */
export const checkAccountDirectly = async (email: string, token: string): Promise<{
  exists: boolean
  message: string
  details?: any
}> => {
  try {
    console.log('🔍 Verificando cuenta directamente en la base de datos para:', email)
    
    // Usar el endpoint de Electric SQL directamente con LiveQuery
    const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id,email&where=email='${email}'&live=true&offset=0_0`
    
    const response = await fetch(url, {
      headers: {
        'X-Access-Token': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      console.log('❌ Error en consulta directa:', response.status, response.statusText)
      return {
        exists: false,
        message: `Error en consulta directa: ${response.status}`,
        details: { status: response.status, statusText: response.statusText }
      }
    }

    const data = await response.json()
    console.log('🔍 Respuesta directa de Electric SQL:', data)

    // Verificar si hay datos reales (no solo objetos de control)
    if (Array.isArray(data) && data.length > 0) {
      const accountData = data.find(item => item.value && item.value.email)
      if (accountData) {
        return {
          exists: true,
          message: 'Email existe en la base de datos',
          details: accountData.value
        }
      }
    }

    return {
      exists: false,
      message: 'Email no existe en la base de datos',
      details: data
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
export const compareElectricVsDirect = async (email: string, electricResult: any, token: string) => {
  console.log('🔄 Comparando Electric SQL vs verificación directa...')
  
  const directResult = await checkAccountDirectly(email, token)
  
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
