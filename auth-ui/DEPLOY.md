# Configuración para Deploy del Auth UI

## 🔐 Variables de Entorno Requeridas

Para que el deploy funcione correctamente, necesitas configurar estas variables en GitHub:

### 🔑 Secrets (Settings → Secrets and variables → Actions → Repository secrets)

```
NEXTAUTH_SECRET=tu-secret-muy-seguro-de-al-menos-32-caracteres
GOOGLE_CLIENT_ID=tu-google-client-id.apps.googleusercontent.com  
GOOGLE_CLIENT_SECRET=tu-google-client-secret
FIREBASE_SERVICE_ACCOUNT=contenido-completo-del-service-account-json
```

### 📋 Variables (Settings → Secrets and variables → Actions → Repository variables)

```
NEXTAUTH_URL=https://tu-dominio-auth.web.app
GOOGLE_PROJECT_ID=tu-proyecto-firebase-id
ALLOWED_DOMAINS=transportapp.com,empresa.com
GOOGLE_HD_DOMAIN=transportapp.com
```

## 🚀 Pasos para Deploy

1. **Configurar Firebase Hosting Target:**
   ```bash
   firebase target:apply hosting transport-auth-ui tu-proyecto-auth-id
   ```

2. **Push cambios al repositorio:**
   ```bash
   git add .
   git commit -m "Add auth-ui deploy configuration"
   git push
   ```

3. **El workflow se ejecutará automáticamente** cuando:
   - Haces push a `main` con cambios en `auth-ui/**`
   - Ejecutas manualmente desde GitHub Actions

## 🌐 URLs esperadas

- **Desarrollo**: http://localhost:3002
- **Producción**: https://tu-dominio-auth.web.app

## 🔧 Configuración de Google OAuth

Recuerda agregar las URLs de producción en Google Cloud Console:

**URIs de origen autorizados:**
```
https://tu-dominio-auth.web.app
```

**URIs de redirección autorizados:**
```
https://tu-dominio-auth.web.app/api/auth/callback/google
```

## 📝 Notas importantes

- El auth-ui se despliega como sitio estático (export)
- NextAuth.js funciona en modo edge runtime
- Las variables de entorno se inyectan en build time
- El sitio se almacena en Firebase Hosting
