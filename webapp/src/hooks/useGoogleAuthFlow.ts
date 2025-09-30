import { useState, useEffect } from 'react'
import { findAccountByEmail, findTenantsByAccountId } from '../db/queries/accountQueries'
import type { ElectricAccountData } from '../db/collections/create-accounts-collection'
import type { ElectricTenantData } from '../db/collections/create-tenants-collection'

export type AuthFlowState = 'loading' | 'checking-account' | 'account-not-found' | 'loading-tenants' | 'tenants-loaded' | 'error'

export type AuthFlowResult = {
  state: AuthFlowState
  account: ElectricAccountData | null
  tenants: ElectricTenantData[]
  error: string | null
}

export const useGoogleAuthFlow = (token: string, email: string) => {
  const [result, setResult] = useState<AuthFlowResult>({
    state: 'loading',
    account: null,
    tenants: [],
    error: null
  })

  useEffect(() => {
    if (!token || !email) return

    const executeAuthFlow = async () => {
      try {
        // Paso 1: Verificar si el account existe
        setResult(prev => ({ ...prev, state: 'checking-account' }))
        
        // Buscar la cuenta
        const account = await findAccountByEmail(token, email)
        
        if (!account) {
          // Account no existe, redirigir a creación de organización
          setResult(prev => ({ 
            ...prev, 
            state: 'account-not-found',
            account: null,
            tenants: []
          }))
          return
        }

        // Account existe, buscar tenants
        setResult(prev => ({ ...prev, state: 'loading-tenants' }))
        const tenants = await findTenantsByAccountId(token, account.id)

        // Mostrar tenants
        setResult(prev => ({ 
          ...prev, 
          state: 'tenants-loaded',
          account,
          tenants
        }))

      } catch (error) {
        console.error('Error en el flujo de autenticación:', error)
        setResult(prev => ({ 
          ...prev, 
          state: 'error',
          error: error instanceof Error ? error.message : 'Error desconocido'
        }))
      }
    }

    executeAuthFlow()
  }, [token, email])

  return result
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
