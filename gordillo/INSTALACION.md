# Instrucciones de Instalación - SisPro ERP

## Pasos para ejecutar el proyecto:

1. **Navegar a la carpeta del proyecto:**
```bash
cd gordillo
```

2. **Instalar dependencias:**
```bash
npm install
```

3. **Ejecutar en modo desarrollo:**
```bash
npm run dev
```

4. **Abrir en el navegador:**
El proyecto se abrirá automáticamente en `http://localhost:3000`

## Si Tailwind CSS no se aplica correctamente:

1. **Detener el servidor** (Ctrl+C)

2. **Reinstalar dependencias:**
```bash
rm -rf node_modules package-lock.json
npm install
```

3. **Ejecutar nuevamente:**
```bash
npm run dev
```

## Verificar que Tailwind esté funcionando:

- Los botones deben tener bordes redondeados y colores
- El sidenav debe tener fondo oscuro (gris-900)
- Los inputs deben tener estilos modernos
- Las transiciones deben funcionar al hacer hover

## Estructura de archivos importantes:

- `postcss.config.js` - Configuración de PostCSS para Tailwind
- `tailwind.config.js` - Configuración de Tailwind CSS
- `src/index.css` - Estilos globales con directivas de Tailwind
- `vite.config.ts` - Configuración de Vite

Si después de seguir estos pasos el problema persiste, verificar que:
1. El archivo `postcss.config.js` existe
2. Las dependencias `tailwindcss`, `autoprefixer` y `postcss` están instaladas
3. El archivo `src/index.css` contiene las directivas `@tailwind`
