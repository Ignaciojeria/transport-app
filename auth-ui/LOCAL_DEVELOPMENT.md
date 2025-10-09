# Configuración para Desarrollo Local

## 🔧 Redirección Automática

El auth-ui ahora detecta automáticamente si está ejecutándose en desarrollo local y ajusta la URL de redirección en consecuencia:

### Desarrollo Local
- **Auth UI**: `http://localhost:3002`
- **Redirección**: `http://localhost:5173/` (webapp)

### Producción
- **Auth UI**: `https://tu-dominio-auth.web.app`
- **Redirección**: `https://console.transport-app.com`

## 🚀 Cómo Funciona

El código detecta el entorno basándose en el `hostname`:

```typescript
const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
const baseUrl = isLocalDev ? 'http://localhost:5173' : 'https://console.transport-app.com'
```

## 📋 Pasos para Desarrollo

1. **Ejecutar Auth UI en localhost:**
   ```bash
   cd auth-ui
   npm run dev
   # Se ejecuta en http://localhost:3002
   ```

2. **Ejecutar Webapp en localhost:**
   ```bash
   cd webapp
   npm run dev
   # Se ejecuta en http://localhost:5173
   ```

3. **Probar el flujo completo:**
   - Ir a `http://localhost:3002`
   - Hacer clic en "Iniciar Sesión con Google"
   - Después de la autenticación, será redirigido automáticamente a `http://localhost:5173/`

## 🔍 URLs de Google OAuth

Asegúrate de que en Google Cloud Console tengas configuradas estas URLs:

### URIs de origen autorizados:
```
http://localhost:3002
https://tu-dominio-auth.web.app
```

### URIs de redirección autorizados:
```
http://localhost:3002/auth/callback
https://tu-dominio-auth.web.app/auth/callback
```

## ✅ Verificación

Para verificar que funciona correctamente:

1. Abre las herramientas de desarrollador (F12)
2. Ve a la consola
3. Inicia el proceso de autenticación
4. Verás en la consola: `🚀 Redirigiendo a: http://localhost:5173#auth=...`

Esto confirma que la detección de entorno está funcionando correctamente.
