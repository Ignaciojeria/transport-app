import { useState, useEffect } from 'react'
import { checkAccountAndGetTenants, type ElectricAccount, type ElectricTenant } from '../services/electricService'
import { useElectricLiveQuery } from './useElectricLiveQuery'

export type AuthFlowState = 'loading' | 'checking-account' | 'account-not-found' | 'loading-tenants' | 'tenants-loaded' | 'error'

export type AuthFlowResult = {
  state: AuthFlowState
  account: ElectricAccount | null
  tenants: ElectricTenant[]
  error: string | null
  retry: () => Promise<void>
}

export const useGoogleAuthFlow = (token: string, email: string) => {
  const [result, setResult] = useState<AuthFlowResult>({
    state: 'loading',
    account: null,
    tenants: [],
    error: null,
    retry: async () => {} // Placeholder inicial
  })

  // Usar LiveQuery para sincronización en tiempo real
  const { isLoading, error } = useElectricLiveQuery(
    () => checkAccountAndGetTenants(token, email),
    [token, email],
    {
      enabled: !!token && !!email,
      refetchInterval: 10000, // Refetch cada 10 segundos para mantener sincronización
      onSuccess: (data) => {
        console.log('🔄 Datos actualizados via LiveQuery:', data)
        
        if (!data) {
          // Account no existe, redirigir a creación de organización
          console.log('ℹ️ Cuenta no encontrada, permitiendo creación de organización')
          setResult(prev => ({ 
            ...prev, 
            state: 'account-not-found',
            account: null,
            tenants: []
          }))
        } else {
          // Account existe, mostrar tenants
          console.log('✅ Cuenta encontrada, mostrando organizaciones existentes:', data.tenants.length)
          setResult(prev => ({ 
            ...prev, 
            state: 'tenants-loaded',
            account: data.account,
            tenants: data.tenants
          }))
        }
      },
      onError: (error) => {
        console.error('❌ Error en LiveQuery:', error)
        setResult(prev => ({ 
          ...prev, 
          state: 'error',
          error: error.message
        }))
      }
    }
  )

  const executeAuthFlow = async () => {
    // Esta función ahora es solo para compatibilidad
    // La sincronización real se maneja con LiveQuery
    console.log('🔄 Ejecutando consulta manual (LiveQuery se encarga de la sincronización)')
  }

  // Actualizar la función retry en el estado
  useEffect(() => {
    setResult(prev => ({ ...prev, retry: executeAuthFlow }))
  }, [token, email])

  // Actualizar estado basado en LiveQuery
  useEffect(() => {
    if (isLoading) {
      setResult(prev => ({ ...prev, state: 'checking-account' }))
    } else if (error) {
      setResult(prev => ({ 
        ...prev, 
        state: 'error',
        error: error.message
      }))
    }
  }, [isLoading, error])

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