# MiCartaPro - Auth UI

Cáscara de autenticación minimalista para MiCartaPro con integración de Gmail/Google OAuth. Diseñado con el estilo visual de la landing page principal e inspirado en la interfaz limpia de Synadia.

## 🚀 Características

- ✅ **Autenticación con Google OAuth** - Solo botón de Gmail (minimalista)
- ✅ **Restricción por dominios** - Solo usuarios de dominios autorizados
- ✅ **Diseño responsive** - Optimizado para desktop y móvil
- ✅ **Partículas animadas** - Efectos visuales del estilo MiCartaPro
- ✅ **Manejo de errores** - Páginas específicas para diferentes tipos de error
- ✅ **Página de éxito** - Confirmación post-autenticación
- ✅ **Estilo minimalista** - Inspirado en Synadia con colores de MiCartaPro

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
   cd micartapro-auth-ui
   ```

2. **Instalar dependencias:**
   ```bash
   npm install
   ```

3. **Configurar variables de entorno:**
   ```bash
   cp env.example .env.local
   ```

4. **Editar el archivo `.env.local`** con tus credenciales reales

## 🚀 Desarrollo

```bash
# Iniciar servidor de desarrollo
npm run dev

# Construir para producción
npm run build

# Iniciar en modo producción
npm run start
```

La aplicación estará disponible en [http://localhost:3003](http://localhost:3003)

## 📱 Rutas Disponibles

- **`/`** - Página de login principal (minimalista)
- **`/auth/callback`** - Callback de OAuth de Google
- **`/auth/error`** - Página de errores de autenticación

## 📄 Licencia

© 2024 MiCartaPro. Todos los derechos reservados.

