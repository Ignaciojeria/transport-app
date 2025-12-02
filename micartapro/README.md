# MiCartaPro Landing Page

Landing page para MiCartaPro - Tu Menú Digital, Sin Complicaciones.

## Características

- 🎨 Diseño moderno y responsivo
- 📱 Optimizado para todos los dispositivos
- 🧪 Demo interactiva embebida de Cadorago
- 💬 Integración con WhatsApp para cotizaciones
- ⚡ Construido con Next.js 14 y TypeScript

## Instalación

```bash
npm install
```

## Desarrollo

```bash
npm run dev
```

Abre [http://localhost:3000](http://localhost:3000) en tu navegador.

## Build

```bash
npm run build
```

El build se generará en la carpeta `out/`.

## Despliegue

Este proyecto está configurado para desplegarse en Firebase Hosting. El archivo `firebase.json` está configurado para usar la carpeta `out/` como directorio público.

## Estructura del Proyecto

```
micartapro/
├── app/
│   ├── globals.css      # Estilos globales
│   ├── layout.tsx        # Layout principal
│   └── page.tsx          # Página principal
├── components/
│   ├── ui/               # Componentes UI reutilizables
│   └── DemoEmbed.tsx     # Componente de demo embebida
├── lib/
│   ├── utils.ts          # Utilidades
│   └── whatsapp.ts       # Funciones de WhatsApp
└── public/               # Archivos estáticos
```

## Configuración

- **WhatsApp**: El número de WhatsApp para cotizaciones está configurado en `lib/whatsapp.ts` (+56957857558)
- **Demo URL**: La URL de la demo está configurada en `components/DemoEmbed.tsx` (https://cadorago.web.app/)

