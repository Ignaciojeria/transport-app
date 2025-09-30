import { useAccountByEmail } from './use-accounts-collection'
import { useTenantsByAccountId } from './use-tenants-collection'
import type { ElectricAccountData } from '../collections/create-accounts-collection'
import type { ElectricTenantData } from '../collections/create-tenants-collection'

export interface AccountWithTenants {
  account: ElectricAccountData | null
  tenants: ElectricTenantData[]
  isLoading: boolean
  error: boolean
}

// Hook principal que combina cuenta y tenants
export const useAccountData = (token: string, email: string): AccountWithTenants => {
  const account = useAccountByEmail(token, email)
  const tenants = useTenantsByAccountId(token, account?.id || '')
  
  return {
    account,
    tenants,
    isLoading: false, // Los hooks individuales manejan su propio loading
    error: false // Los hooks individuales manejan su propio error
  }
}
