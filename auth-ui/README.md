# TransportApp - Auth UI

Cáscara de autenticación minimalista para TransportApp con integración de Gmail/Google OAuth. Diseñado con el estilo visual de la landing page principal e inspirado en la interfaz limpia de Synadia.

## 🚀 Características

- ✅ **Autenticación con Google OAuth** - Solo botón de Gmail (minimalista)
- ✅ **Restricción por dominios** - Solo usuarios de dominios autorizados
- ✅ **Diseño responsive** - Optimizado para desktop y móvil
- ✅ **Partículas animadas** - Efectos visuales del estilo TransportApp
- ✅ **Manejo de errores** - Páginas específicas para diferentes tipos de error
- ✅ **Página de éxito** - Confirmación post-autenticación
- ✅ **Estilo minimalista** - Inspirado en Synadia con colores de TransportApp

## 🛠️ Tecnologías

- **Next.js 14** - Framework React con App Router
- **NextAuth.js** - Sistema de autenticación completo
- **TailwindCSS** - Framework de estilos utilitarios
- **Framer Motion** - Animaciones fluidas
- **TypeScript** - Tipado estático
- **Lucide React** - Iconos modernos

## 📦 Instalación

1. **Navegar al directorio del proyecto:**
   ```bash
   cd auth-ui
   ```

2. **Instalar dependencias:**
   ```bash
   npm install
   # o
   yarn install
   # o
   pnpm install
   ```

3. **Configurar variables de entorno:**
   ```bash
   cp env.example .env.local
   ```

4. **Editar el archivo `.env.local`** con tus credenciales reales:
   ```env
   NEXTAUTH_URL=http://localhost:3001
   NEXTAUTH_SECRET=tu-secret-muy-seguro-de-al-menos-32-caracteres
   
   # Credenciales de Google OAuth
   GOOGLE_CLIENT_ID=tu-client-id.apps.googleusercontent.com
   GOOGLE_CLIENT_SECRET=tu-client-secret
   
   # Opcional: Restringir dominios
   ALLOWED_DOMAINS=transportapp.com,empresa.com
   GOOGLE_HD_DOMAIN=transportapp.com
   ```

## 🔧 Configuración de Google OAuth

### 1. Crear proyecto en Google Cloud Console

1. Ir a [Google Cloud Console](https://console.cloud.google.com)
2. Crear un nuevo proyecto o seleccionar uno existente
3. Habilitar la **Google+ API** y **Gmail API**

### 2. Configurar OAuth 2.0

1. Ir a **Credenciales** > **Crear credenciales** > **ID de cliente OAuth 2.0**
2. Tipo de aplicación: **Aplicación web**
3. Nombre: `TransportApp Auth`
4. URIs de origen autorizados:
   ```
   http://localhost:3001
   https://tu-dominio.com
   ```
5. URIs de redirección autorizados:
   ```
   http://localhost:3001/api/auth/callback/google
   https://tu-dominio.com/api/auth/callback/google
   ```

### 3. Configurar pantalla de consentimiento OAuth

1. Ir a **Pantalla de consentimiento OAuth**
2. Tipo de usuario: **Interno** (para uso corporativo)
3. Completar información de la aplicación:
   - Nombre: "TransportApp"
   - Email de soporte: tu-email@empresa.com
   - Dominio autorizado: empresa.com

## 🚀 Desarrollo

```bash
# Iniciar servidor de desarrollo
npm run dev

# Construir para producción
npm run build

# Iniciar en modo producción
npm run start
```

La aplicación estará disponible en [http://localhost:3001](http://localhost:3001)

## 📱 Rutas Disponibles

- **`/`** - Página de login principal (minimalista)
- **`/success`** - Página de éxito post-autenticación
- **`/auth/error`** - Página de errores de autenticación
- **`/api/auth/*`** - Endpoints de NextAuth.js

## 🎨 Estilo y Diseño

El diseño está inspirado en:
- **Colores**: Paleta de azules de TransportApp (blue-600, blue-700)
- **Tipografía**: Inter font family
- **Efectos**: Partículas animadas y degradados de la landing page
- **Layout**: Panel dividido similar a OpenObserve
- **Iconos**: Lucide React con temática de transporte

## 🔐 Seguridad

### Restricciones Implementadas

1. **Dominios autorizados**: Solo emails de dominios específicos
2. **Google Workspace**: Restricción HD para organizaciones
3. **Tokens seguros**: NextAuth.js maneja tokens JWT automáticamente
4. **Variables de entorno**: Credenciales nunca expuestas al cliente

### Testing

La aplicación está configurada solo para autenticación con Google OAuth. 
No hay credenciales de demo ya que se enfoca en la integración corporativa.

## 🐛 Solución de Problemas

### Error: "Configuration"
- Verificar que todas las variables de entorno estén configuradas
- Verificar que NEXTAUTH_SECRET tenga al menos 32 caracteres

### Error: "AccessDenied" 
- El dominio del email no está en ALLOWED_DOMAINS
- El usuario no tiene acceso al Google Workspace configurado

### Error: "CredentialsSignin"
- Email o contraseña incorrectos para login con credenciales
- Verificar la implementación en `route.ts`

## 🚀 Despliegue

### Variables de entorno de producción:
```env
NEXTAUTH_URL=https://auth.transportapp.com
NEXTAUTH_SECRET=secret-de-produccion-super-seguro
GOOGLE_CLIENT_ID=production-client-id
GOOGLE_CLIENT_SECRET=production-client-secret
ALLOWED_DOMAINS=transportapp.com
```

### Plataformas recomendadas:
- **Vercel** (configuración automática)
- **Netlify** 
- **AWS Amplify**
- **Railway**

## 📄 Licencia

© 2024 TransportApp. Todos los derechos reservados.

---

## 🔗 Enlaces Útiles

- [NextAuth.js Documentation](https://next-auth.js.org/)
- [Google OAuth 2.0 Setup](https://developers.google.com/identity/protocols/oauth2)
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [Framer Motion Documentation](https://www.framer.com/motion/)
