import { useQuery } from '@tanstack/react-query'
import { findAccountByEmail, findTenantsByAccountId } from '../db/queries/accountQueries'
import type { ElectricAccountData } from '../db/collections/create-accounts-collection'
import type { ElectricTenantData } from '../db/collections/create-tenants-collection'

export interface AccountWithTenants {
  account: ElectricAccountData | null
  tenants: ElectricTenantData[]
  isLoading: boolean
  error: Error | null
}

export const useAccountData = (token: string, email: string): AccountWithTenants => {
  // Query para buscar la cuenta
  const accountQuery = useQuery({
    queryKey: ['account', email],
    queryFn: () => findAccountByEmail(token, email),
    enabled: !!token && !!email,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })

  // Query para buscar tenants (solo si existe la cuenta)
  const tenantsQuery = useQuery({
    queryKey: ['tenants', accountQuery.data?.id],
    queryFn: () => findTenantsByAccountId(token, accountQuery.data!.id),
    enabled: !!accountQuery.data?.id,
    staleTime: 5 * 60 * 1000, // 5 minutos
  })

  return {
    account: accountQuery.data || null,
    tenants: tenantsQuery.data || [],
    isLoading: accountQuery.isLoading || (!!accountQuery.data && tenantsQuery.isLoading),
    error: accountQuery.error || tenantsQuery.error || null
  }
}

// Hook para buscar cuenta específica
export const useAccountByEmail = (token: string, email: string) => {
  return useQuery({
    queryKey: ['account', email],
    queryFn: () => findAccountByEmail(token, email),
    enabled: !!token && !!email,
    staleTime: 5 * 60 * 1000,
  })
}

// Hook para buscar tenants de una cuenta específica
export const useTenantsByAccountId = (token: string, accountId: string) => {
  return useQuery({
    queryKey: ['tenants', accountId],
    queryFn: () => findTenantsByAccountId(token, accountId),
    enabled: !!token && !!accountId,
    staleTime: 5 * 60 * 1000,
  })
}