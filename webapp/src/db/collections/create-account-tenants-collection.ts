import { createCollection } from '@tanstack/react-db'
import { electricCollectionOptions } from '@tanstack/electric-db-collection'

// Tipo para la estructura que devuelve Electric
export type ElectricAccountTenantData = {
  account_id: string
  tenant_id: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

// Factory para crear la colección inyectando el token
export const createAccountTenantsCollection = (token: string) => {
  const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=account_tenants&columns=account_id,tenant_id&offset=-1`
  
  return createCollection(
    electricCollectionOptions({
      id: 'account_tenants',
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
      getKey: (accountTenant: ElectricAccountTenantData) => `${accountTenant.account_id}-${accountTenant.tenant_id}`,
      
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