# Sistema de Deploy Automático para Tenants

Este repositorio contiene la configuración para el deploy automático de sitios web de tenants usando Firebase Hosting y GitHub Actions.

## Estructura del Repositorio

```
├── .github/
│   └── workflows/
│       └── deploy.yml          # Workflow de deploy automático
├── config/
│   └── tenant.json             # Configuración del tenant
├── data/                       # Datos específicos del tenant
├── assets/                     # Recursos estáticos
├── index.html                  # Página principal del sitio
├── firebase.json               # Configuración de Firebase Hosting
├── .gitignore                  # Archivos a ignorar en Git
└── README.md                   # Este archivo
```

## Cómo Funciona

### 1. Creación Automática
- Cuando se crea un nuevo tenant, se genera automáticamente este repositorio
- Se crean todos los archivos de template necesarios
- Se configura el workflow de GitHub Actions

### 2. Deploy Automático
- Cada vez que se hace push a `main` o `master`, se ejecuta el deploy
- El sitio se despliega automáticamente a Firebase Hosting
- Se asigna un DNS específico para el tenant: `{{.TenantID}}-site`

### 3. Configuración de Firebase
- **Proyecto**: `transport-app-agents`
- **Target**: `{{.TenantID}}-site`
- **URL**: `https://{{.TenantID}}-site.web.app`

## Variables de Entorno Requeridas

Para que el deploy funcione, el repositorio necesita tener configuradas estas variables en GitHub:

### Secrets (Settings → Secrets and variables → Actions)
- `FIREBASE_SERVICE_ACCOUNT`: JSON del service account de Firebase
- `GITHUB_TOKEN`: Token de GitHub (se genera automáticamente)

### Variables (Settings → Secrets and variables → Actions)
- `GOOGLE_PROJECT_ID`: ID del proyecto de Google Cloud (`transport-app-agents`)

## Personalización

### Modificar el Sitio
1. Edita `index.html` para cambiar el contenido
2. Agrega archivos en `assets/` para recursos estáticos
3. Modifica `config/tenant.json` para cambiar la configuración
4. Haz commit y push para desplegar automáticamente

### Agregar Funcionalidades
1. Crea nuevos archivos HTML/CSS/JS según necesites
2. Actualiza `firebase.json` si necesitas cambiar la configuración de hosting
3. El workflow se ejecutará automáticamente en cada push

## Troubleshooting

### Deploy Fallido
1. Verifica que las variables de entorno estén configuradas
2. Revisa los logs del workflow en la pestaña "Actions"
3. Asegúrate de que el service account tenga permisos de Firebase Hosting

### Sitio No Se Actualiza
1. Verifica que el push se haya hecho a la rama `main` o `master`
2. Revisa que los archivos modificados estén en las rutas correctas
3. Espera unos minutos para que Firebase procese el deploy

## Contacto

Para soporte técnico, contacta al equipo de Transport App.
