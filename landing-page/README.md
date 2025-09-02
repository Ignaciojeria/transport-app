# TransportApp Landing Page

Una landing page moderna y responsive para TransportApp, construida con Next.js, TypeScript, Tailwind CSS y componentes de magic-ui.

## 🚀 Características

- **Diseño Responsive**: Optimizado para desktop, tablet y móvil
- **Animaciones Suaves**: Usando Framer Motion para transiciones elegantes
- **Componentes Magic UI**: Sistema de diseño consistente y moderno
- **SEO Optimizado**: Meta tags y estructura semántica
- **Performance**: Carga rápida y optimizada

## 📋 Secciones Incluidas

1. **Hero Section**: Propuesta de valor principal con CTA
2. **¿Cómo funciona?**: 3 pasos simples del proceso
3. **Vista del conductor**: Experiencia del usuario final
4. **Beneficios**: Ventajas para empresas
5. **Testimonios**: Casos de uso y feedback de clientes
6. **Call to Action**: Botones de conversión
7. **Footer**: Enlaces y información adicional

## 🛠️ Tecnologías

- **Next.js 14**: Framework de React con App Router
- **TypeScript**: Tipado estático
- **Tailwind CSS**: Estilos utilitarios
- **Framer Motion**: Animaciones
- **Lucide React**: Iconos
- **Magic UI**: Componentes de interfaz

## 🚀 Instalación y Ejecución

### Prerrequisitos

- Node.js 18+ 
- npm o yarn

### Pasos

1. **Instalar dependencias**:
   ```bash
   cd landing-page
   npm install
   ```

2. **Ejecutar en modo desarrollo**:
   ```bash
   npm run dev
   ```

3. **Abrir en el navegador**:
   ```
   http://localhost:3000
   ```

### Scripts Disponibles

- `npm run dev` - Servidor de desarrollo
- `npm run build` - Construir para producción
- `npm run start` - Servidor de producción
- `npm run lint` - Linter de código

## 🎨 Personalización

### Colores

Los colores principales están definidos en `tailwind.config.js`:

- **Primary**: Azul (#2563eb)
- **Secondary**: Gris (#f1f5f9)
- **Success**: Verde (#16a34a)
- **Warning**: Amarillo (#eab308)
- **Error**: Rojo (#dc2626)

### Componentes

Los componentes UI están en `components/ui/` y siguen el patrón de magic-ui:

- `Button` - Botones con variantes
- `Card` - Tarjetas de contenido

### Animaciones

Las animaciones están configuradas con Framer Motion:

- **Fade In**: Aparición suave
- **Slide Up**: Deslizamiento desde abajo
- **Stagger**: Animaciones escalonadas

## 📱 Responsive Design

La landing page está optimizada para:

- **Desktop**: 1024px+
- **Tablet**: 768px - 1023px
- **Mobile**: 320px - 767px

## 🔧 Estructura del Proyecto

```
landing-page/
├── app/
│   ├── globals.css          # Estilos globales
│   ├── layout.tsx           # Layout principal
│   └── page.tsx             # Página principal
├── components/
│   └── ui/                  # Componentes UI
├── lib/
│   └── utils.ts             # Utilidades
├── public/                  # Archivos estáticos
└── ...config files
```

## 🚀 Despliegue

### Vercel (Recomendado)

1. Conecta tu repositorio a Vercel
2. Configura las variables de entorno si es necesario
3. Despliega automáticamente

### Netlify

1. Conecta tu repositorio a Netlify
2. Configura el build command: `npm run build`
3. Configura el publish directory: `.next`

### Docker

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

## 📈 Optimizaciones

- **Imágenes**: Optimización automática con Next.js
- **CSS**: Purge automático de Tailwind
- **JavaScript**: Code splitting automático
- **SEO**: Meta tags y estructura semántica

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature
3. Commit tus cambios
4. Push a la rama
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT.
