# Usar Supabase SDK en Otro Subdominio

## Opción 1: Usar el Token del Fragment (Recomendado)

El código actual ya pasa el `access_token` de Supabase en el fragment. En tu otro subdominio:

### 1. Instalar Supabase SDK en el otro proyecto:
```bash
npm install @supabase/supabase-js
```

### 2. Crear cliente de Supabase:
```typescript
// lib/supabase.ts
import { createClient } from '@supabase/supabase-js'

export const supabase = createClient(
  'https://rbpdhapfcljecofrscnj.supabase.co',
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'
)
```

### 3. Procesar el fragment y establecer la sesión:
```typescript
// En tu app principal (ej: App.tsx o main.tsx)
useEffect(() => {
  const fragment = window.location.hash
  
  if (fragment.startsWith('#auth=')) {
    try {
      const encodedData = fragment.substring(6) // Remove '#auth='
      const authData = JSON.parse(atob(encodedData))
      
      if (authData.access_token && authData.provider === 'supabase') {
        // Establecer la sesión en Supabase
        supabase.auth.setSession({
          access_token: authData.access_token,
          refresh_token: authData.refresh_token,
        }).then(({ data, error }) => {
          if (error) {
            console.error('Error estableciendo sesión:', error)
          } else {
            console.log('✅ Sesión establecida:', data.session?.user?.email)
            // Limpiar fragment
            window.history.replaceState({}, '', window.location.pathname)
          }
        })
      }
    } catch (err) {
      console.error('Error procesando fragment:', err)
    }
  }
}, [])
```

### 4. Usar el SDK normalmente:
```typescript
// Cualquier componente
import { supabase } from '@/lib/supabase'

// Obtener usuario actual
const { data: { user } } = await supabase.auth.getUser()

// Hacer queries a la base de datos
const { data, error } = await supabase
  .from('tu_tabla')
  .select('*')
  .eq('user_id', user?.id)
```

## Opción 2: Cookies Compartidas (Para Subdominios)

Si ambos subdominios están bajo el mismo dominio (ej: `auth.micartapro.com` y `app.micartapro.com`):

### Configurar Supabase para usar cookies:
```typescript
export const supabase = createClient(
  'https://rbpdhapfcljecofrscnj.supabase.co',
  'tu-anon-key',
  {
    auth: {
      storage: typeof window !== 'undefined' 
        ? window.localStorage 
        : undefined,
      autoRefreshToken: true,
      persistSession: true,
      detectSessionInUrl: true,
      // Para compartir entre subdominios, usar cookies
      storageKey: 'sb-auth-token',
    },
  }
)
```

Luego configurar cookies con dominio compartido en el servidor.

## Recomendación

**Usa la Opción 1** (token en fragment) porque:
- ✅ Ya está implementado
- ✅ Funciona entre cualquier dominio/subdominio
- ✅ Más simple
- ✅ No requiere configuración de cookies

