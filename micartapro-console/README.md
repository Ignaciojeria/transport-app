# MiCartaPro Console

Panel de control para gestionar menús digitales, plantillas y códigos QR.

## 🚀 Inicio Rápido

### Instalar dependencias
```bash
npm install
```

### Desarrollo
```bash
npm run dev
```

La aplicación se ejecutará en `http://localhost:5174`

### Build para producción
```bash
npm run build
```

## 🔐 Autenticación

Este proyecto se conecta con Supabase y recibe la autenticación desde `micartapro-auth-ui`.

### Flujo de Autenticación

1. Usuario inicia sesión en `micartapro-auth-ui` (puerto 3003)
2. Después de autenticarse, es redirigido a este proyecto con un fragment `#auth=...`
3. El proyecto procesa el fragment y establece la sesión de Supabase
4. El usuario puede usar todas las funcionalidades del SDK de Supabase

### URLs de Redirección

- **Desarrollo local**: `http://localhost:5174`
- **Producción**: `https://console.micartapro.com`

## 📦 Estructura

```
micartapro-console/
├── src/
│   ├── lib/
│   │   ├── supabase.js    # Cliente de Supabase
│   │   └── auth.js        # Lógica de autenticación
│   ├── App.svelte         # Componente principal
│   ├── main.js            # Punto de entrada
│   └── app.css            # Estilos globales
├── index.html
└── package.json
```

## 🛠️ Tecnologías

- **Svelte 5**: Framework frontend
- **Vite**: Build tool
- **Supabase**: Backend y autenticación
- **TailwindCSS**: Estilos

