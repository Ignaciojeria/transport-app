# Servicios de Electric SQL

Este directorio contiene los servicios para interactuar con Electric SQL siguiendo el patrón del proyecto `/electric`.

## 🔧 Configuración de la API

### Parámetros Requeridos

Todas las consultas a Electric SQL requieren el parámetro `offset`:

- **`offset=-1`**: Para obtener todos los datos actuales (sincronización inicial)
- **`offset=0_0`**: Para modo en vivo (después de la sincronización inicial)

### Endpoint Base

```
https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape
```

### Headers Requeridos

```javascript
{
  'X-Access-Token': `Bearer ${token}`
}
```

## 📋 Ejemplos de Uso

### 1. Buscar Cuenta por Email

```javascript
const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=accounts&columns=id,email&where=email='${email}'&offset=-1`
```

### 2. Buscar Account Tenants

```javascript
const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=account_tenants&columns=account_id,tenant_id&where=account_id='${accountId}'&offset=-1`
```

### 3. Buscar Tenants

```javascript
const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=tenants&columns=id,name,country&where=id='${tenantId}'&offset=-1`
```

## 🚨 Errores Comunes

### Error: "offset can't be blank"

**Causa**: Falta el parámetro `offset` en la URL.

**Solución**: Agregar `&offset=-1` al final de la URL.

```javascript
// ❌ Incorrecto
const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=accounts&columns=id,email&where=email='${email}'`

// ✅ Correcto
const url = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=accounts&columns=id,email&where=email='${email}'&offset=-1`
```

## 🔄 Flujo de Sincronización

1. **Sincronización Inicial**: `offset=-1`
   - Obtiene todos los datos actuales
   - Usado para cargar datos por primera vez

2. **Modo En Vivo**: `offset=0_0&live=true`
   - Recibe actualizaciones en tiempo real
   - Requiere un `handle` de forma

## 📊 Estructura de Respuesta

Electric SQL devuelve un **array de objetos** con la siguiente estructura:

```typescript
[
  {
    headers: {
      operation: 'insert' | 'update' | 'delete',
      relation: ['schema', 'table'],
      // ... otros metadatos
    },
    key: "\"public\".\"table_name\"/\"record_id\"",
    value: {
      id: string,
      email: string,
      // ... otros campos de la tabla
    }
  },
  {
    headers: {
      control: 'snapshot-end',
      xip_list: [],
      xmax: string,
      xmin: string
    }
    // ... objeto de control (sin value)
  }
]
```

### 🔍 Procesamiento de Datos

```typescript
// Filtrar solo los objetos con datos (que tienen 'value')
const dataItems = response.filter(item => item.value && item.value.email)

// Obtener el primer registro de datos
const firstRecord = dataItems.find(item => item.value && item.value.email)
if (firstRecord) {
  const account = firstRecord.value
  // Usar account.id, account.email, etc.
}
```

## 🛠️ Servicios Disponibles

- **`findAccountByEmail()`**: Busca cuentas por email
- **`findTenantsByAccountId()`**: Obtiene tenants asociados a una cuenta
- **`checkAccountAndGetTenants()`**: Función principal que verifica cuenta y obtiene tenants

## 📚 Referencias

- [Electric SQL HTTP API Documentation](https://electric-sql.com/docs/api/http)
- [Proyecto Electric de Referencia](../electric/src/db/collections/)
