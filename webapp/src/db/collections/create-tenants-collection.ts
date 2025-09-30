// Tipo para la estructura que devuelve Electric para tenants
export type ElectricTenantData = {
  id: string
  name: string
  country: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

// Factory para crear la colección inyectando el token y tenant_id
export const createTenantsCollection = (token: string, tenantId: string) => {
  const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=tenants&columns=id,name,country&where=id='${tenantId}'`
  
  return {
    id: 'tenants',
    url,
    headers: {
      'X-Access-Token': `Bearer ${token}`,
    },
    parser: {
      timestamptz: (iso: string) => new Date(iso),
      timestamp: (iso: string) => new Date(iso),
    },
    getKey: (tenant: ElectricTenantData) => tenant.id,
    
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
