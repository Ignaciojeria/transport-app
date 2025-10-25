# Deploy de Gordillo a Firebase Hosting

## Configuración de Firebase

El proyecto gordillo está configurado para deploy automático a Firebase Hosting con el sitio llamado `gordillo`.

### Archivos de configuración:

1. **`firebase.json`** - Configuración de Firebase Hosting
   - Site: `gordillo`
   - Public directory: `dist`
   - Rewrites para SPA (Single Page Application)
   - Headers de cache para assets estáticos

2. **`.github/workflows/deploy-gordillo.yml`** - Workflow de GitHub Actions
   - Se ejecuta en push a main/master cuando hay cambios en `gordillo/**`
   - Build automático con Node.js 20
   - Deploy automático a Firebase Hosting

### Variables de entorno requeridas:

En GitHub Secrets y Variables:
- `FIREBASE_SERVICE_ACCOUNT` - Service account de Firebase
- `GOOGLE_PROJECT_ID` - ID del proyecto de Google Cloud/Firebase

### Comandos locales:

```bash
# Instalar dependencias
npm install

# Build para producción
npm run build

# Preview local del build
npm run preview

# Deploy manual (requiere Firebase CLI)
firebase deploy --project gordillo
```

### Estructura del build:

```
gordillo/
├── dist/                 # Archivos generados por npm run build
│   ├── index.html
│   ├── assets/
│   └── ...
├── firebase.json         # Configuración de Firebase
└── package.json         # Scripts de build
```

### URLs de acceso:

- **Desarrollo local**: `http://localhost:3000`
- **Firebase Hosting**: `https://gordillo.web.app` (o dominio personalizado)

### Notas importantes:

1. El workflow se ejecuta automáticamente en cada push a main/master
2. El sitio se llama `gordillo` en Firebase Hosting
3. Los assets estáticos tienen cache de 1 año
4. Todas las rutas se redirigen a `index.html` para SPA
