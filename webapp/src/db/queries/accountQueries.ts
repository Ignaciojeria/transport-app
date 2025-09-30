import type { ElectricAccountData } from '../collections/create-accounts-collection'
import type { ElectricTenantData } from '../collections/create-tenants-collection'
import type { ElectricAccountTenantData } from '../collections/create-account-tenants-collection'

// Función para buscar cuenta por email usando la collection
export const findAccountByEmail = async (token: string, email: string): Promise<ElectricAccountData | null> => {
  try {
    console.log('🔍 Buscando cuenta para email:', email)
    
    // Consulta directa a Electric SQL
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
    
    // Electric SQL devuelve un array directo con objetos {key, value, headers}
    if (Array.isArray(data) && data.length > 0) {
      // Buscar el primer item que tenga value (no solo headers)
      const accountItem = data.find(item => item.value)
      if (accountItem && accountItem.value) {
        console.log('✅ Cuenta encontrada:', accountItem.value)
        return accountItem.value
      }
    }
    
    // Fallback: verificar si hay data.rows (formato alternativo)
    if (data.rows && data.rows.length > 0) {
      console.log('✅ Cuenta encontrada (formato rows):', data.rows[0])
      return data.rows[0]
    }
    
    console.log('ℹ️ No se encontró cuenta para el email:', email)
    return null
  } catch (error) {
    console.error('❌ Error al buscar cuenta:', error)
    return null
  }
}

// Función para buscar tenants por account_id usando las collections
export const findTenantsByAccountId = async (token: string, accountId: string): Promise<ElectricTenantData[]> => {
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
    
    // Electric SQL devuelve un array directo con objetos {key, value, headers}
    let accountTenants: ElectricAccountTenantData[] = []
    if (Array.isArray(accountTenantsData) && accountTenantsData.length > 0) {
      accountTenants = accountTenantsData
        .filter(item => item.value)
        .map(item => item.value)
    } else if (accountTenantsData.rows && accountTenantsData.rows.length > 0) {
      accountTenants = accountTenantsData.rows
    }
    
    if (accountTenants.length === 0) {
      console.log('ℹ️ No hay tenants asociados a la cuenta')
      return []
    }

    // Paso 2: Obtener tenant_ids
    const tenantIds = accountTenants.map((at: ElectricAccountTenantData) => at.tenant_id)
    console.log('🔍 Tenant IDs a consultar:', tenantIds)
    
    // Paso 3: Buscar detalles de cada tenant
    const tenants: ElectricTenantData[] = []
    
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
          
          // Electric SQL devuelve un array directo con objetos {key, value, headers}
          if (Array.isArray(tenantData) && tenantData.length > 0) {
            const tenantItem = tenantData.find(item => item.value)
            if (tenantItem && tenantItem.value) {
              tenants.push(tenantItem.value)
            }
          } else if (tenantData.rows && tenantData.rows.length > 0) {
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
