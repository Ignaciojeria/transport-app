import { useAccountData } from './useAccountData'
import type { ElectricAccountData } from '../db/collections/create-accounts-collection'
import type { ElectricTenantData } from '../db/collections/create-tenants-collection'

export type AuthFlowState = 'loading' | 'checking-account' | 'account-not-found' | 'loading-tenants' | 'tenants-loaded' | 'error'

export type AuthFlowResult = {
  state: AuthFlowState
  account: ElectricAccountData | null
  tenants: ElectricTenantData[]
  error: string | null
}

export const useGoogleAuthFlow = (token: string, email: string): AuthFlowResult => {
  const { account, tenants, isLoading, error } = useAccountData(token, email)
  
  // Determinar el estado basado en los datos
  if (isLoading) {
    return {
      state: 'loading',
      account: null,
      tenants: [],
      error: null
    }
  }
  
  if (error) {
    return {
      state: 'error',
      account: null,
      tenants: [],
      error: 'Error al cargar datos'
    }
  }
  
  if (!account) {
    return {
      state: 'account-not-found',
      account: null,
      tenants: [],
      error: null
    }
  }
  
  return {
    state: 'tenants-loaded',
    account,
    tenants,
    error: null
  }
}

// Hook para manejar la redirección basada en el estado
export const useAuthRedirect = (result: AuthFlowResult) => {
  const getRedirectPath = () => {
    switch (result.state) {
      case 'account-not-found':
        return '/create-organization'
      case 'tenants-loaded':
        return '/dashboard'
      case 'error':
        return '/error'
      default:
        return null
    }
  }

  const shouldRedirect = () => {
    return result.state === 'account-not-found' || result.state === 'tenants-loaded'
  }

  return {
    redirectPath: getRedirectPath(),
    shouldRedirect: shouldRedirect(),
    isLoading: result.state === 'loading' || result.state === 'checking-account' || result.state === 'loading-tenants'
  }
}