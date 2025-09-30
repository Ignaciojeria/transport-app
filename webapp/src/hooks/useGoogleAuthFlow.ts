import { useState, useEffect } from 'react'
import { createAccountsCollection, type ElectricAccountData } from '../db/collections/create-accounts-collection'
import { createAccountTenantsCollection, type ElectricAccountTenantData } from '../db/collections/create-account-tenants-collection'
import { createTenantsCollection, type ElectricTenantData } from '../db/collections/create-tenants-collection'
import { getTenantIds } from '../db/collections/account-tenants-helpers'

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
        
        const accountsCollection = createAccountsCollection(token, email)
        
        // Simular la obtención de datos de la colección
        // En una implementación real, usarías el hook de Electric SQL
        const accountsData = await fetchAccountsData(accountsCollection)
        
        if (accountsData.length === 0) {
          // Account no existe, redirigir a creación de organización
          setResult(prev => ({ 
            ...prev, 
            state: 'account-not-found',
            account: null,
            tenants: []
          }))
          return
        }

        const account = accountsData[0]
        
        // Paso 2: Account existe, obtener tenants asociados
        setResult(prev => ({ 
          ...prev, 
          state: 'loading-tenants',
          account 
        }))

        const accountTenantsCollection = createAccountTenantsCollection(token, account.id)
        const accountTenantsData = await fetchAccountTenantsData(accountTenantsCollection)
        
        if (accountTenantsData.length === 0) {
          // No hay tenants asociados
          setResult(prev => ({ 
            ...prev, 
            state: 'tenants-loaded',
            tenants: []
          }))
          return
        }

        // Paso 3: Obtener detalles de cada tenant
        const tenantIds = getTenantIds(accountTenantsData)
        const tenantsPromises = tenantIds.map(tenantId => 
          fetchTenantData(createTenantsCollection(token, tenantId))
        )
        
        const tenantsData = await Promise.all(tenantsPromises)
        const tenants = tenantsData.filter(tenant => tenant !== null) as ElectricTenantData[]

        setResult(prev => ({ 
          ...prev, 
          state: 'tenants-loaded',
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

// Funciones auxiliares para simular la obtención de datos
// En una implementación real, estas serían reemplazadas por los hooks de Electric SQL
const fetchAccountsData = async (collection: { url: string; headers: Record<string, string> }): Promise<ElectricAccountData[]> => {
  // Simulación - en la realidad usarías el hook de Electric SQL
  try {
    const response = await fetch(collection.url, {
      headers: collection.headers
    })
    const data = await response.json()
    return data.rows || []
  } catch (error) {
    console.error('Error fetching accounts:', error)
    return []
  }
}

const fetchAccountTenantsData = async (collection: { url: string; headers: Record<string, string> }): Promise<ElectricAccountTenantData[]> => {
  // Simulación - en la realidad usarías el hook de Electric SQL
  try {
    const response = await fetch(collection.url, {
      headers: collection.headers
    })
    const data = await response.json()
    return data.rows || []
  } catch (error) {
    console.error('Error fetching account tenants:', error)
    return []
  }
}

const fetchTenantData = async (collection: { url: string; headers: Record<string, string> }): Promise<ElectricTenantData | null> => {
  // Simulación - en la realidad usarías el hook de Electric SQL
  try {
    const response = await fetch(collection.url, {
      headers: collection.headers
    })
    const data = await response.json()
    return data.rows?.[0] || null
  } catch (error) {
    console.error('Error fetching tenant:', error)
    return null
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
