# Colecciones de Electric SQL

Estas colecciones permiten interactuar con las tablas de Electric SQL usando el endpoint `/electric-me`.

## 1. Colección de Accounts

### Características
- **Endpoint**: `/electric-me/v1/shape?table=accounts&columns=id,email&where=email='{email}'`
- **Campos**: Obtiene `id` y `email` de la tabla `accounts`
- **Filtrado por email**: Requerido, filtra por `email` específico
- **Autenticación**: Usa Bearer token para autenticación

### Uso
```typescript
import { createAccountsCollection, type ElectricAccountData } from './db/collections'

// Crear colección con email específico
const accountsCollection = createAccountsCollection('your-token', 'user@example.com')
```

### Tipos
```typescript
type ElectricAccountData = {
  id: string
  email: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}
```

### Funciones Helper
- `getAccountIds(accounts)`: Extrae solo los IDs de un array de cuentas
- `getAccountEmails(accounts)`: Extrae solo los emails de un array de cuentas
- `getAccountByEmail(accounts, email)`: Obtiene cuenta por email
- `filterAccountsByReference(accounts, referenceId)`: Filtra cuentas por reference_id

## 2. Colección de Account Tenants

### Características
- **Endpoint**: `/electric-me/v1/shape?table=account_tenants&columns=account_id,tenant_id&where=account_id='{accountId}'`
- **Filtrado por account_id**: Requerido, filtra por `account_id` específico
- **Campos**: Obtiene solo `account_id` y `tenant_id`

### Uso
```typescript
import { createAccountTenantsCollection, type ElectricAccountTenantData } from './db/collections'

// Crear colección con account_id específico
const accountTenantsCollection = createAccountTenantsCollection('your-token', 'account-123')
```

### Tipos
```typescript
type ElectricAccountTenantData = {
  account_id: string
  tenant_id: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}
```

### Funciones Helper
- `getTenantIds(accountTenants)`: Extrae solo los tenant_ids
- `getUniqueTenantIds(accountTenants)`: Obtiene tenant_ids únicos
- `hasTenant(accountTenants, tenantId)`: Verifica si una cuenta tiene un tenant específico
- `getAccountTenantsByTenantId(accountTenants, tenantId)`: Obtiene account_tenants por tenant_id
- `getAccountTenantByKeys(accountTenants, accountId, tenantId)`: Obtiene account_tenant por account_id y tenant_id
- `mapAccountToTenants(accountTenants)`: Mapea account_id a tenant_ids

## 3. Colección de Tenants

### Características
- **Endpoint**: `/electric-me/v1/shape?table=tenants&columns=id,name,country&where=id='{tenantId}'`
- **Filtrado por id**: Requerido, filtra por `id` (UUID) específico
- **Campos**: Obtiene `id`, `name` y `country`

### Uso
```typescript
import { createTenantsCollection, type ElectricTenantData } from './db/collections'

// Crear colección con tenant_id específico
const tenantsCollection = createTenantsCollection('your-token', 'tenant-uuid-123')
```

### Tipos
```typescript
type ElectricTenantData = {
  id: string
  name: string
  country: string
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}
```

### Funciones Helper
- `getTenantName(tenant)`: Obtiene el nombre del tenant
- `getTenantCountry(tenant)`: Obtiene el país del tenant
- `getTenantById(tenants, id)`: Obtiene tenant por ID
- `getTenantsByCountry(tenants, country)`: Obtiene tenants por país
- `searchTenantsByName(tenants, searchTerm)`: Busca tenants por nombre (búsqueda parcial)
- `getUniqueCountries(tenants)`: Obtiene países únicos
- `groupTenantsByCountry(tenants)`: Agrupa tenants por país
- `getTenantSummary(tenant)`: Obtiene resumen del tenant (id, name, country)