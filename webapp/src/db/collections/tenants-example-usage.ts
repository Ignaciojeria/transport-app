import { createTenantsCollection, type ElectricTenantData } from './create-tenants-collection'
import { 
  getTenantName, 
  getTenantCountry, 
  getTenantById, 
  getTenantsByCountry,
  searchTenantsByName,
  getUniqueCountries,
  groupTenantsByCountry,
  getTenantSummary
} from './tenants-helpers'

// Ejemplo de uso de la colección de tenants
export const exampleTenantsUsage = () => {
  // Token de autenticación (debería venir de tu sistema de auth)
  const token = 'your-auth-token-here'
  
  // Crear la colección con token
  const tenantsCollection = createTenantsCollection(token)
  
  return {
    tenantsCollection
  }
}

// Ejemplo de uso de las funciones helper
export const exampleTenantsHelperUsage = (tenants: ElectricTenantData[]) => {
  // Obtener nombre del primer tenant
  const firstName = tenants.length > 0 ? getTenantName(tenants[0]) : null
  console.log('First tenant name:', firstName)
  
  // Obtener país del primer tenant
  const firstCountry = tenants.length > 0 ? getTenantCountry(tenants[0]) : null
  console.log('First tenant country:', firstCountry)
  
  // Obtener tenant por ID
  const specificTenant = getTenantById(tenants, 'tenant-uuid-123')
  console.log('Specific tenant:', specificTenant)
  
  // Obtener tenants por país
  const tenantsInChile = getTenantsByCountry(tenants, 'Chile')
  console.log('Tenants in Chile:', tenantsInChile)
  
  // Buscar tenants por nombre
  const searchResults = searchTenantsByName(tenants, 'Transport')
  console.log('Tenants with "Transport" in name:', searchResults)
  
  // Obtener países únicos
  const uniqueCountries = getUniqueCountries(tenants)
  console.log('Unique countries:', uniqueCountries)
  
  // Agrupar tenants por país
  const tenantsByCountry = groupTenantsByCountry(tenants)
  console.log('Tenants grouped by country:', tenantsByCountry)
  
  // Obtener resumen del tenant
  const tenantSummary = tenants.length > 0 ? getTenantSummary(tenants[0]) : null
  console.log('Tenant summary:', tenantSummary)
  
  return {
    firstName,
    firstCountry,
    specificTenant,
    tenantsInChile,
    searchResults,
    uniqueCountries,
    tenantsByCountry,
    tenantSummary
  }
}
