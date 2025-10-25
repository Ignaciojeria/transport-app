# SisPro ERP - Grupo Gordillo

Sitio web para la gestión de productos del Grupo Gordillo con interfaz moderna y funcionalidad completa.

## Características

- **Sidenav interactiva** con navegación por módulos
- **Formulario de creación de productos** con campos:
  - `nombreProducto`: Nombre del producto
  - `descripcion`: Descripción del producto  
  - `precio`: Precio del producto
- **Tabla de productos** para visualizar los productos creados
- **Diseño responsive** con Tailwind CSS
- **Interfaz moderna** similar al diseño de referencia

## Instalación

1. Navegar a la carpeta del proyecto:
```bash
cd gordillo
```

2. Instalar dependencias:
```bash
npm install
```

3. Ejecutar en modo desarrollo:
```bash
npm run dev
```

4. Abrir en el navegador: `http://localhost:3000`

## Uso

1. **Navegación**: Usar el menú lateral (sidenav) para navegar entre diferentes módulos
2. **Crear Producto**: 
   - Seleccionar "Bodega" en el menú lateral
   - Hacer clic en "Gestión de Productos"
   - Completar el formulario con los datos del producto
   - Hacer clic en "CREAR"
3. **Ver Productos**: Los productos creados aparecerán automáticamente en la tabla debajo del formulario

## Estructura del Proyecto

```
gordillo/
├── src/
│   ├── components/
│   │   ├── Header.tsx          # Header con logo y menú de usuario
│   │   ├── Sidenav.tsx          # Navegación lateral
│   │   └── ProductForm.tsx      # Formulario de creación de productos
│   ├── App.tsx                  # Componente principal
│   ├── main.tsx                 # Punto de entrada
│   └── index.css                # Estilos globales
├── package.json                 # Dependencias del proyecto
├── tailwind.config.js           # Configuración de Tailwind CSS
├── tsconfig.json                # Configuración de TypeScript
└── vite.config.ts               # Configuración de Vite
```

## Tecnologías Utilizadas

- **React 18** con TypeScript
- **Tailwind CSS** para estilos
- **Vite** como bundler
- **Lucide React** para iconos

## Contrato de Producto

El formulario utiliza el siguiente contrato para la creación de productos:

```json
{
  "nombreProducto": "string",
  "descripcion": "string", 
  "precio": "number"
}
```

## Scripts Disponibles

- `npm run dev` - Ejecutar en modo desarrollo
- `npm run build` - Construir para producción
- `npm run preview` - Previsualizar build de producción
- `npm run lint` - Ejecutar linter
