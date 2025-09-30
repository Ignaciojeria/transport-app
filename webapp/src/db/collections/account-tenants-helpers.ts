import type { ElectricAccountTenantData } from './create-account-tenants-collection'

// Función helper para obtener solo los tenant_ids de una cuenta
export const getTenantIds = (accountTenants: ElectricAccountTenantData[]): string[] => {
  return accountTenants.map(accountTenant => accountTenant.tenant_id)
}

// Función helper para obtener los tenant_ids únicos de una cuenta
export const getUniqueTenantIds = (accountTenants: ElectricAccountTenantData[]): string[] => {
  const tenantIds = getTenantIds(accountTenants)
  return [...new Set(tenantIds)]
}

// Función helper para verificar si una cuenta tiene un tenant específico
export const hasTenant = (accountTenants: ElectricAccountTenantData[], tenantId: string): boolean => {
  return accountTenants.some(accountTenant => accountTenant.tenant_id === tenantId)
}

// Función helper para obtener todos los account_tenants de un tenant específico
export const getAccountTenantsByTenantId = (
  accountTenants: ElectricAccountTenantData[], 
  tenantId: string
): ElectricAccountTenantData[] => {
  return accountTenants.filter(accountTenant => accountTenant.tenant_id === tenantId)
}

// Función helper para obtener la relación account-tenant por account_id y tenant_id
export const getAccountTenantByKeys = (
  accountTenants: ElectricAccountTenantData[], 
  accountId: string,
  tenantId: string
): ElectricAccountTenantData | undefined => {
  return accountTenants.find(accountTenant => 
    accountTenant.account_id === accountId && accountTenant.tenant_id === tenantId
  )
}

// Función helper para mapear account_id a tenant_ids
export const mapAccountToTenants = (accountTenants: ElectricAccountTenantData[]): Record<string, string[]> => {
  return accountTenants.reduce((acc, accountTenant) => {
    if (!acc[accountTenant.account_id]) {
      acc[accountTenant.account_id] = []
    }
    acc[accountTenant.account_id].push(accountTenant.tenant_id)
    return acc
  }, {} as Record<string, string[]>)
}
