import { useLiveQuery } from '@tanstack/react-db'
import { useMemo } from 'react'
import { createTenantsCollection, type ElectricTenantData } from '../collections/create-tenants-collection'
import { createAccountTenantsCollection } from '../collections/create-account-tenants-collection'

// Hook personalizado que combina las collections con useLiveQuery
export const useTenantsCollection = (token: string) => {
  const tenantsCollection = useMemo(() => createTenantsCollection(token), [token])
  const accountTenantsCollection = useMemo(() => createAccountTenantsCollection(token), [token])
  
  const tenantsQuery = useLiveQuery((queryBuilder: any) => 
    queryBuilder.from({ tenant: tenantsCollection })
  )
  
  const accountTenantsQuery = useLiveQuery((queryBuilder: any) => 
    queryBuilder.from({ accountTenant: accountTenantsCollection })
  )
  
  return {
    tenantsCollection,
    accountTenantsCollection,
    tenantsQuery,
    accountTenantsQuery,
    tenants: tenantsQuery.data || [],
    accountTenants: accountTenantsQuery.data || [],
    isLoading: tenantsQuery.isLoading || accountTenantsQuery.isLoading,
    error: tenantsQuery.isError || accountTenantsQuery.isError,
    // Métodos de las collections para mutaciones
    insertTenant: tenantsCollection.insert,
    updateTenant: tenantsCollection.update,
    deleteTenant: tenantsCollection.delete,
    insertAccountTenant: accountTenantsCollection.insert,
    updateAccountTenant: accountTenantsCollection.update,
    deleteAccountTenant: accountTenantsCollection.delete,
  }
}

// Hook para buscar tenants de una cuenta específica
export const useTenantsByAccountId = (token: string, accountId: string): ElectricTenantData[] => {
  const { tenants, accountTenants, error } = useTenantsCollection(token)
  
  const associatedTenants = useMemo(() => {
    if (!accountId || !Array.isArray(accountTenants) || !Array.isArray(tenants)) {
      return []
    }
    
    // Encontrar relaciones account_tenants para esta cuenta
    const accountTenantRelations = accountTenants.filter((at: any) => at.account_id === accountId)
    const tenantIds = accountTenantRelations.map((at: any) => at.tenant_id)
    
    // Filtrar tenants que están asociados a esta cuenta
    return tenants.filter((tenant: any) => tenantIds.includes(tenant.id)) as ElectricTenantData[]
  }, [accountId, accountTenants, tenants])
  
  if (error) {
    console.error('Error cargando tenants:', error)
    return []
  }
  
  return associatedTenants
}
