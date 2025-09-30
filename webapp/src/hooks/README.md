# Hooks de Autenticación

Este directorio contiene hooks personalizados para manejar el flujo de autenticación con Google y la gestión de organizaciones.

## useGoogleAuthFlow

Hook principal que maneja todo el flujo de autenticación después del login con Google.

### Características
- Verifica si existe una cuenta con el email de Google
- Si no existe, indica que debe crear una organización
- Si existe, carga todos los tenants asociados a la cuenta
- Maneja estados de loading y errores

### Uso
```typescript
import { useGoogleAuthFlow } from '../hooks/useGoogleAuthFlow'

const MyComponent = () => {
  const authResult = useGoogleAuthFlow(token, email)
  
  // authResult contiene:
  // - state: 'loading' | 'checking-account' | 'account-not-found' | 'loading-tenants' | 'tenants-loaded' | 'error'
  // - account: ElectricAccountData | null
  // - tenants: ElectricTenantData[]
  // - error: string | null
}
```

## useAuthRedirect

Hook auxiliar que determina a dónde redirigir basado en el estado de autenticación.

### Uso
```typescript
import { useAuthRedirect } from '../hooks/useGoogleAuthFlow'

const MyComponent = () => {
  const authResult = useGoogleAuthFlow(token, email)
  const { redirectPath, shouldRedirect, isLoading } = useAuthRedirect(authResult)
  
  // redirectPath: '/create-organization' | '/dashboard' | '/error' | null
  // shouldRedirect: boolean
  // isLoading: boolean
}
```

## Flujo Completo

1. **Usuario inicia sesión con Google** → Obtiene token y email
2. **useGoogleAuthFlow verifica la cuenta**:
   - Busca en la tabla `accounts` por email
   - Si no existe → `state: 'account-not-found'`
   - Si existe → continúa al paso 3
3. **Carga tenants asociados**:
   - Busca en `account_tenants` por `account_id`
   - Para cada `tenant_id`, obtiene detalles de `tenants`
   - Retorna `state: 'tenants-loaded'` con la lista completa

## Componentes Relacionados

- `GoogleAuthHandler`: Componente principal que usa los hooks
- `TenantsList`: Muestra la lista de organizaciones disponibles
- `CreateOrganization`: Formulario para crear nueva organización
- `LoadingSpinner`: Componente de loading reutilizable

## Ejemplo de Implementación

```typescript
import GoogleAuthHandler from '../components/GoogleAuthHandler'

const App = () => {
  const [token, setToken] = useState('')
  const [email, setEmail] = useState('')
  
  // Lógica de Google OAuth aquí...
  
  return (
    <GoogleAuthHandler 
      token={token}
      email={email}
      onError={(error) => console.error('Auth error:', error)}
    />
  )
}
```
