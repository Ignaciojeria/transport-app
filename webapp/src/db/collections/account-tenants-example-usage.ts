import { createAccountTenantsCollection, type ElectricAccountTenantData } from './create-account-tenants-collection'
import { 
  getTenantIds, 
  getUniqueTenantIds, 
  hasTenant, 
  getAccountTenantsByTenantId,
  getAccountTenantByKeys,
  mapAccountToTenants
} from './account-tenants-helpers'

// Ejemplo de uso de la colección de account_tenants
export const exampleAccountTenantsUsage = () => {
  // Token de autenticación (debería venir de tu sistema de auth)
  const token = 'your-auth-token-here'
  const accountId = 'specific-account-id'
  
  // Crear la colección con token y account_id
  const accountTenantsCollection = createAccountTenantsCollection(token, accountId)
  
  return {
    accountTenantsCollection
  }
}

// Ejemplo de uso de las funciones helper
export const exampleHelperUsage = (accountTenants: ElectricAccountTenantData[]) => {
  // Obtener todos los tenant_ids
  const tenantIds = getTenantIds(accountTenants)
  console.log('Tenant IDs:', tenantIds)
  
  // Obtener tenant_ids únicos
  const uniqueTenantIds = getUniqueTenantIds(accountTenants)
  console.log('Unique Tenant IDs:', uniqueTenantIds)
  
  // Verificar si la cuenta tiene un tenant específico
  const hasSpecificTenant = hasTenant(accountTenants, 'tenant-123')
  console.log('Has tenant-123:', hasSpecificTenant)
  
  // Obtener account_tenants por tenant_id
  const tenantsForAccount = getAccountTenantsByTenantId(accountTenants, 'tenant-123')
  console.log('Account tenants for tenant-123:', tenantsForAccount)
  
  // Obtener account_tenant por account_id y tenant_id
  const specificAccountTenant = getAccountTenantByKeys(accountTenants, 'account-123', 'tenant-456')
  console.log('Specific account tenant:', specificAccountTenant)
  
  // Mapear account_id a tenant_ids
  const accountToTenantsMap = mapAccountToTenants(accountTenants)
  console.log('Account to Tenants mapping:', accountToTenantsMap)
  
  return {
    tenantIds,
    uniqueTenantIds,
    hasSpecificTenant,
    tenantsForAccount,
    specificAccountTenant,
    accountToTenantsMap
  }
}
