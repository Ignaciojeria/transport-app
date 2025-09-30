import type { ElectricTenantData } from './create-tenants-collection'

// Función helper para obtener solo el nombre del tenant
export const getTenantName = (tenant: ElectricTenantData): string => {
  return tenant.name
}

// Función helper para obtener solo el país del tenant
export const getTenantCountry = (tenant: ElectricTenantData): string => {
  return tenant.country
}

// Función helper para obtener el tenant por ID
export const getTenantById = (tenants: ElectricTenantData[], id: string): ElectricTenantData | undefined => {
  return tenants.find(tenant => tenant.id === id)
}

// Función helper para obtener tenants por país
export const getTenantsByCountry = (tenants: ElectricTenantData[], country: string): ElectricTenantData[] => {
  return tenants.filter(tenant => tenant.country === country)
}

// Función helper para buscar tenants por nombre (búsqueda parcial)
export const searchTenantsByName = (tenants: ElectricTenantData[], searchTerm: string): ElectricTenantData[] => {
  const term = searchTerm.toLowerCase()
  return tenants.filter(tenant => tenant.name.toLowerCase().includes(term))
}

// Función helper para obtener un mapa de países únicos
export const getUniqueCountries = (tenants: ElectricTenantData[]): string[] => {
  const countries = tenants.map(tenant => tenant.country)
  return [...new Set(countries)]
}

// Función helper para agrupar tenants por país
export const groupTenantsByCountry = (tenants: ElectricTenantData[]): Record<string, ElectricTenantData[]> => {
  return tenants.reduce((acc, tenant) => {
    if (!acc[tenant.country]) {
      acc[tenant.country] = []
    }
    acc[tenant.country].push(tenant)
    return acc
  }, {} as Record<string, ElectricTenantData[]>)
}

// Función helper para obtener información resumida del tenant
export const getTenantSummary = (tenant: ElectricTenantData): { id: string; name: string; country: string } => {
  return {
    id: tenant.id,
    name: tenant.name,
    country: tenant.country
  }
}
