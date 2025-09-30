import { createCollection } from '@tanstack/react-db'
import { electricCollectionOptions } from '@tanstack/electric-db-collection'

// Tipo para la estructura que devuelve Electric
export type ElectricAccountData = {
  id: string
  email: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

// Factory para crear la colección inyectando el token y email
export const createAccountsCollection = (token: string, email: string) => {
  const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id,email&where=email='${email}'&offset=-1`
  
  return createCollection(
    electricCollectionOptions({
      id: 'accounts',
      shapeOptions: {
        url,
        headers: {
          'X-Access-Token': `Bearer ${token}`,
        },
        parser: {
          timestamptz: (iso: string) => new Date(iso),
          timestamp: (iso: string) => new Date(iso),
        },
      },
      getKey: (account: ElectricAccountData) => account.id,
      
      async onInsert() {
        return { txid: [Date.now()] }
      },
      
      async onUpdate() {
        return { txid: [Date.now()] }
      },
      
      async onDelete() {
        return { txid: [Date.now()] }
      },
    })
  )
}