import { useMemo } from 'react'
import { createAccountsCollection, type ElectricAccountData } from '../db/collections/create-accounts-collection'

export const useAccountByEmail = (token: string | null, email: string | null) => {
  const shapeResult = useMemo(() => {
    if (!token || !email) {
      return null
    }
    return createAccountsCollection(token, email)
  }, [token, email])

  const data = shapeResult?.data as ElectricAccountData[] | undefined
  const isLoading = shapeResult?.isLoading ?? false
  const error = shapeResult?.error

  // Determinar si la cuenta existe
  const accountExists = data && data.length > 0
  const accountId = accountExists ? data[0].id : null

  return {
    accountExists,
    accountId,
    isLoading,
    error,
    data
  }
}
