# Configuración para Deploy del MiCartaPro Auth UI

## 🔐 Variables de Entorno Requeridas

Para que el deploy funcione correctamente, necesitas configurar estas variables en GitHub:

### 🔑 Secrets (Settings → Secrets and variables → Actions → Repository secrets)

```
NEXTAUTH_SECRET=tu-secret-muy-seguro-de-al-menos-32-caracteres
GOOGLE_CLIENT_ID=27303662337-1icetdk7186gt37lh5ruq5hu0vq1r57t.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=tu-google-client-secret
FIREBASE_SERVICE_ACCOUNT=contenido-completo-del-service-account-json
```

### 📋 Variables (Settings → Secrets and variables → Actions → Repository variables)

```
NEXTAUTH_URL=https://micartapro-auth.web.app
GOOGLE_PROJECT_ID=tu-proyecto-firebase-id
ALLOWED_DOMAINS=micartapro.com,empresa.com
GOOGLE_HD_DOMAIN=micartapro.com
```

## 🚀 Pasos para Deploy

1. **Configurar Firebase Hosting Target:**
   ```bash
   firebase target:apply hosting micartapro-auth tu-proyecto-auth-id
   ```

2. **Build del proyecto:**
   ```bash
   cd micartapro-auth-ui
   npm install
   npm run build
   ```

3. **Deploy a Firebase:**
   ```bash
   firebase deploy --only hosting:micartapro-auth
   ```

## 🌐 URLs esperadas

- **Desarrollo**: http://localhost:3003
- **Producción**: https://micartapro-auth.web.app (o tu dominio personalizado)

## 🔧 Configuración de Google OAuth

Recuerda agregar las URLs de producción en Google Cloud Console:

**URIs de origen autorizados:**
```
https://micartapro-auth.web.app
http://localhost:3003
```

**URIs de redirección autorizados:**
```
https://micartapro-auth.web.app/auth/callback
http://localhost:3003/auth/callback
```

## 📝 Notas importantes

- El auth-ui se despliega como sitio estático (export)
- Las variables de entorno se inyectan en build time
- El sitio se almacena en Firebase Hosting
- El target de Firebase es: `micartapro-auth`

