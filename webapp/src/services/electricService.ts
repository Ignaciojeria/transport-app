/**
 * Servicio para consultar Electric SQL usando el patrón del proyecto electric
 */

export interface ElectricAccount {
  id: string
  email: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

export interface ElectricTenant {
  id: string
  name: string
  country: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

export interface ElectricAccountTenant {
  account_id: string
  tenant_id: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

/**
 * Busca una cuenta por email en Electric SQL
 * @param token - Token de autenticación
 * @param email - Email a buscar
 * @returns Cuenta encontrada o null si no existe
 */
export const findAccountByEmail = async (token: string, email: string): Promise<ElectricAccount | null> => {
  try {
    console.log('🔍 Buscando cuenta para email:', email)
    
    const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id,email&where=email='${email}'&offset=-1`
    
    const response = await fetch(url, {
      headers: {
        'X-Access-Token': `Bearer ${token}`,
      }
    })

    if (!response.ok) {
      console.error('❌ Error al consultar Electric SQL:', response.status)
      return null
    }

    const data = await response.json()
    console.log('🔍 Respuesta de Electric SQL:', data)
    
    // Electric SQL maneja el caché automáticamente
    if (data.rows && data.rows.length > 0) {
      console.log('✅ Cuenta encontrada:', data.rows[0])
      return data.rows[0]
    }
    
    console.log('ℹ️ No se encontró cuenta para el email:', email)
    return null
  } catch (error) {
    console.error('❌ Error al buscar cuenta:', error)
    return null
  }
}

/**
 * Busca los tenants (organizaciones) asociados a una cuenta
 * Lógica: account -> account_tenants -> tenants
 * @param token - Token de autenticación
 * @param accountId - ID de la cuenta
 * @returns Lista de tenants (organizaciones) asociados
 */
export const findTenantsByAccountId = async (token: string, accountId: string): Promise<ElectricTenant[]> => {
  try {
    console.log('🔍 Buscando organizaciones para account_id:', accountId)
    
    // Paso 1: Buscar relaciones en account_tenants
    const accountTenantsUrl = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=account_tenants&columns=account_id,tenant_id&where=account_id='${accountId}'&offset=-1`
    
    const accountTenantsResponse = await fetch(accountTenantsUrl, {
      headers: {
        'X-Access-Token': `Bearer ${token}`,
      }
    })

    if (!accountTenantsResponse.ok) {
      console.error('❌ Error al consultar account_tenants:', accountTenantsResponse.status)
      return []
    }

    const accountTenantsData = await accountTenantsResponse.json()
    console.log('🔍 Account tenants encontrados:', accountTenantsData)
    
    if (!accountTenantsData.rows || accountTenantsData.rows.length === 0) {
      console.log('ℹ️ No hay tenants asociados a la cuenta')
      return []
    }

    // Paso 2: Obtener tenant_ids
    const tenantIds = accountTenantsData.rows.map((at: ElectricAccountTenant) => at.tenant_id)
    console.log('🔍 Tenant IDs a consultar:', tenantIds)
    
    // Paso 3: Buscar detalles de cada tenant
    const tenants: ElectricTenant[] = []
    
    for (const tenantId of tenantIds) {
      try {
        const tenantUrl = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=tenants&columns=id,name,country&where=id='${tenantId}'&offset=-1`
        
        const tenantResponse = await fetch(tenantUrl, {
          headers: {
            'X-Access-Token': `Bearer ${token}`,
          }
        })

        if (tenantResponse.ok) {
          const tenantData = await tenantResponse.json()
          if (tenantData.rows && tenantData.rows.length > 0) {
            tenants.push(tenantData.rows[0])
          }
        }
      } catch (error) {
        console.error(`❌ Error al consultar tenant ${tenantId}:`, error)
      }
    }
    
    console.log('✅ Organizaciones encontradas:', tenants)
    return tenants
  } catch (error) {
    console.error('❌ Error al buscar tenants:', error)
    return []
  }
}

/**
 * Verifica si una cuenta existe y obtiene sus tenants
 * @param token - Token de autenticación
 * @param email - Email a verificar
 * @returns Objeto con la cuenta y sus tenants, o null si no existe
 */
export const checkAccountAndGetTenants = async (token: string, email: string): Promise<{
  account: ElectricAccount
  tenants: ElectricTenant[]
} | null> => {
  try {
    console.log('🔍 Verificando cuenta y obteniendo tenants para:', email)
    
    // Buscar la cuenta
    const account = await findAccountByEmail(token, email)
    if (!account) {
      console.log('ℹ️ Cuenta no encontrada, usuario puede crear organización')
      return null
    }
    
    // Buscar los tenants asociados
    const tenants = await findTenantsByAccountId(token, account.id)
    console.log('✅ Cuenta encontrada con tenants:', { account, tenants })
    
    return { account, tenants }
  } catch (error) {
    console.error('❌ Error al verificar cuenta:', error)
    return null
  }
}
