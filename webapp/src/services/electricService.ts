/**
 * Servicio para consultar Electric SQL usando el patrón del proyecto electric
 */

import { compareElectricVsDirect } from '../utils/directDbCheck'
import { syncElectricShapeComplete, parseAccountData, parseAccountTenantsData } from './electricSyncService'

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
 * Busca una cuenta por email en Electric SQL usando el patrón de colecciones
 * @param token - Token de autenticación
 * @param email - Email a buscar
 * @returns Cuenta encontrada o null si no existe
 */
export const findAccountByEmail = async (token: string, email: string): Promise<ElectricAccount | null> => {
  try {
    console.log('🔍 Buscando cuenta en Electric SQL para email:', email)
    
    // Usar sincronización incremental correcta
    const shapeId = `accounts_${email}`
    const baseUrl = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id,email&where=email='${email}'`
    
    const result = await syncElectricShapeComplete(
      shapeId,
      baseUrl,
      token,
      parseAccountData
    )

    if (result.error) {
      console.error('❌ Error en sincronización:', result.error)
      return null
    }

    if (!result.data) {
      console.log('ℹ️ No se encontró cuenta para el email:', email)
      return null
    }

    console.log('✅ Cuenta encontrada:', result.data)
    console.log('✅ Nuevo offset:', result.newOffset)
    
    // Comparar con verificación directa para detectar inconsistencias
    const comparison = await compareElectricVsDirect(email, result.data, token)
    
    // Si hay inconsistencia, no devolver los datos obsoletos
    if (!comparison.consistent) {
      console.warn('⚠️ Datos obsoletos detectados en Electric SQL, no devolviendo datos')
      console.warn('⚠️ Electric SQL tiene caché obsoleto, la cuenta no existe realmente')
      return null
    }
    
    return result.data
  } catch (error) {
    console.error('❌ Error al buscar cuenta en Electric SQL:', error)
    return null
  }
}

/**
 * Busca los tenants asociados a una cuenta
 * @param token - Token de autenticación
 * @param accountId - ID de la cuenta
 * @returns Lista de tenants asociados
 */
export const findTenantsByAccountId = async (token: string, accountId: string): Promise<ElectricTenant[]> => {
  try {
    console.log('🔍 Buscando tenants para account_id:', accountId)
    
    // Primero obtener las relaciones account_tenants usando sincronización incremental
    const accountTenantsShapeId = `account_tenants_${accountId}`
    const accountTenantsBaseUrl = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=account_tenants&columns=account_id,tenant_id&where=account_id='${accountId}'`
    
    const accountTenantsResult = await syncElectricShapeComplete(
      accountTenantsShapeId,
      accountTenantsBaseUrl,
      token,
      parseAccountTenantsData
    )

    if (accountTenantsResult.error) {
      console.error('❌ Error al sincronizar account_tenants:', accountTenantsResult.error)
      return []
    }

    const accountTenantItems = accountTenantsResult.data || []
    
    if (accountTenantItems.length === 0) {
      console.log('ℹ️ No hay tenants asociados a la cuenta')
      return []
    }

    // Obtener los detalles de cada tenant usando sincronización incremental
    const tenantIds = accountTenantItems.map(item => item.tenant_id)
    console.log('🔍 Tenant IDs a consultar:', tenantIds)
    
    const tenants: ElectricTenant[] = []
    
    for (const tenantId of tenantIds) {
      try {
        const tenantShapeId = `tenant_${tenantId}`
        const tenantBaseUrl = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=tenants&columns=id,name,country&where=id='${tenantId}'`
        
        const tenantResult = await syncElectricShapeComplete(
          tenantShapeId,
          tenantBaseUrl,
          token,
          parseAccountData // Usar el mismo parser para datos individuales
        )

        if (tenantResult.data) {
          tenants.push(tenantResult.data)
        }
      } catch (error) {
        console.error(`❌ Error al consultar tenant ${tenantId}:`, error)
      }
    }
    
    console.log('✅ Tenants encontrados:', tenants)
    return tenants
  } catch (error) {
    console.error('❌ Error al buscar tenants en Electric SQL:', error)
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
    console.error('❌ Error al verificar cuenta y obtener tenants:', error)
    return null
  }
}
