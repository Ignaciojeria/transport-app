# Funcionalidad de Catálogo de Productos

## Descripción

Se ha implementado una nueva funcionalidad en la aplicación web que permite gestionar productos dentro de las organizaciones. Al hacer clic en una organización desde la lista, se muestra una sidenav con un menú que incluye "Catalog" como primer elemento.

## Componentes Creados

### 1. Tipos TypeScript (`src/types/product.ts`)
Define las interfaces para el contrato de productos según las especificaciones:
- `Product`: Estructura completa del producto
- `PaymentInfo`: Información de pagos
- `StockInfo`: Información de stock
- `PriceInfo`: Información de precios
- `LogisticsInfo`: Información logística
- `CreateProductRequest`: Contrato para crear productos

### 2. SideNav (`src/components/SideNav.tsx`)
Componente de navegación lateral que incluye:
- Menú con opciones: Catalog, Usuarios, Analytics, Reportes, Configuración
- Información de la organización seleccionada
- Diseño responsive (se oculta en móviles)
- Botón de cierre

### 3. Catalog (`src/components/Catalog.tsx`)
Componente principal para gestionar productos:
- Lista de productos con búsqueda y filtros
- Botón para crear nuevos productos
- Tarjetas de productos con información completa
- Estadísticas del catálogo
- Funciones de editar y eliminar productos

### 4. CreateProduct (`src/components/CreateProduct.tsx`)
Formulario completo para crear productos con:
- Información básica (nombre, descripción, imagen)
- Configuración de pagos (moneda, métodos, proveedor)
- Gestión de stock (unidades, peso, volumen)
- Configuración de precios (fijo, por peso, por volumen)
- Información logística (dimensiones, horarios, costos)
- Validación de campos requeridos

### 5. OrganizationDashboard (`src/components/OrganizationDashboard.tsx`)
Componente principal que maneja la navegación:
- Alterna entre lista de organizaciones y vista con sidenav
- Gestiona el estado de la sidenav y menú activo
- Renderiza el contenido según la opción seleccionada

## Funcionalidades Implementadas

### ✅ Navegación
- Al hacer clic en una organización se muestra la sidenav
- Menú lateral con "Catalog" como primer elemento
- Navegación entre diferentes secciones
- Botón para volver a la lista de organizaciones

### ✅ Gestión de Productos
- Crear productos con el contrato especificado
- Listar productos existentes
- Buscar productos por nombre, descripción o ID
- Eliminar productos
- Estadísticas del catálogo

### ✅ Formulario de Creación
- Todos los campos del contrato implementados
- Validación de campos requeridos
- Interfaz intuitiva y responsive
- Manejo de arrays para horarios y costos

### ✅ Diseño Responsive
- Sidenav que se adapta a diferentes tamaños de pantalla
- Formularios optimizados para móviles
- Tarjetas de productos con diseño flexible

## Uso

1. **Acceder a una organización**: Haz clic en cualquier organización de la lista
2. **Navegar al catálogo**: El menú "Catalog" se selecciona automáticamente
3. **Crear producto**: Haz clic en "Nuevo Producto" y completa el formulario
4. **Gestionar productos**: Usa las opciones de búsqueda, edición y eliminación

## Contrato de Producto

El formulario implementa exactamente el contrato especificado:

```json
{
  "referenceID": "customerReference",
  "name": "Palta Hass",
  "description": "Palta fresca",
  "image": "https://example.com/palta.jpg",
  "payment": {
    "currency": "CLP",
    "methods": ["credit_card", "debit_card", "transfer"],
    "provider": "Transbank"
  },
  "stock": {
    "fixed": { "availableUnits": 0 },
    "weight": { "availableWeight": 250 },
    "volume": { "availableVolume": 0 }
  },
  "price": {
    "fixedPrice": 250000,
    "weight": { "unitSize": 1, "pricePerUnit": 0 },
    "volume": { "unitSize": 1, "pricePerUnit": 0 }
  },
  "logistics": {
    "dimensions": { "height": 100, "length": 100, "width": 100 },
    "availabilityTime": [
      {
        "timeRange": { "from": "09:00", "to": "22:00" },
        "daysOfWeek": ["mon", "tue", "wed", "thu", "fri"]
      }
    ],
    "costs": [
      {
        "condition": "prime",
        "type": "fixed",
        "value": 50,
        "timeRange": { "from": "09:00", "to": "18:00" }
      }
    ]
  }
}
```

## Próximos Pasos

- Implementar persistencia de datos (API backend)
- Agregar funcionalidad de edición de productos
- Implementar las otras secciones del menú (Usuarios, Analytics, etc.)
- Agregar validaciones más avanzadas
- Implementar subida de imágenes
