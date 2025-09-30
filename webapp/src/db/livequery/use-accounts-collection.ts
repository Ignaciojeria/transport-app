import { useLiveQuery } from '@tanstack/react-db'
import { useMemo } from 'react'
import { createAccountsCollection, type ElectricAccountData } from '../collections/create-accounts-collection'

// Hook personalizado que combina la collection con useLiveQuery
export const useAccountsCollection = (token: string, email: string) => {
  const collection = useMemo(() => createAccountsCollection(token, email), [token, email])
  
  const query = useLiveQuery((queryBuilder: any) => 
    queryBuilder.from({ account: collection })
  )
  
  // Buscar la cuenta por email
  const account = useMemo(() => {
    if (Array.isArray(query.data) && query.data.length > 0) {
      const foundAccount = query.data.find((item: any) => item.email === email)
      return foundAccount as ElectricAccountData || null
    }
    return null
  }, [query.data, email])
  
  return {
    collection,
    query,
    account,
    isLoading: query.isLoading,
    error: query.isError,
    // Métodos de la collection para mutaciones
    insert: collection.insert,
    update: collection.update,
    delete: collection.delete,
  }
}

// Hook simplificado que solo devuelve la cuenta encontrada
export const useAccountByEmail = (token: string, email: string): ElectricAccountData | null => {
  const { account, error } = useAccountsCollection(token, email)
  
  if (error) {
    console.error('Error cargando cuenta:', error)
    return null
  }
  
  return account
}
