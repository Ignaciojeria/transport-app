import { createAccountsCollection, type ElectricAccountData } from './create-accounts-collection'

// Ejemplo de uso de la colección de accounts
export const exampleUsage = () => {
  // Token de autenticación (debería venir de tu sistema de auth)
  const token = 'your-auth-token-here'
  const email = 'user@example.com'
  
  // Crear la colección con token y email
  const accountsCollection = createAccountsCollection(token, email)
  
  return {
    accountsCollection
  }
}

// Función helper para obtener solo los IDs de las cuentas
export const getAccountIds = (accounts: ElectricAccountData[]): string[] => {
  return accounts.map(account => account.id)
}

// Función helper para obtener solo los emails de las cuentas
export const getAccountEmails = (accounts: ElectricAccountData[]): string[] => {
  return accounts.map(account => account.email)
}

// Función helper para filtrar cuentas por reference_id
export const filterAccountsByReference = (accounts: ElectricAccountData[], referenceId: string): ElectricAccountData[] => {
  return accounts.filter(account => account.reference_id === referenceId)
}

// Función helper para obtener cuenta por email
export const getAccountByEmail = (accounts: ElectricAccountData[], email: string): ElectricAccountData | undefined => {
  return accounts.find(account => account.email === email)
}
