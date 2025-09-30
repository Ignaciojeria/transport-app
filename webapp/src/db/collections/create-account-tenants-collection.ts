// Tipo para la estructura que devuelve Electric para account_tenants
export type ElectricAccountTenantData = {
  account_id: string
  tenant_id: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

// Factory para crear la colección inyectando el token y account_id
export const createAccountTenantsCollection = (token: string, accountId: string) => {
  const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=account_tenants&columns=account_id,tenant_id&where=account_id='${accountId}'`
  
  return {
    id: 'account_tenants',
    url,
    headers: {
      'X-Access-Token': `Bearer ${token}`,
    },
    parser: {
      timestamptz: (iso: string) => new Date(iso),
      timestamp: (iso: string) => new Date(iso),
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
  }
}
