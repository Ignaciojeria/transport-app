# Vista de Desarrollo

Esta vista te permite probar y desarrollar componentes de manera interactiva. Actualmente incluye un componente de mapa que puede recibir y mostrar lineStrings (rutas).

## Componentes Disponibles

### Map Component

El componente `Map` es reutilizable y puede ser usado en cualquier parte de tu aplicación.

#### Props

- `lineString: number[][]` - Array de coordenadas [lat, lng]
- `center: [number, number]` - Centro del mapa (default: [40.4168, -3.7038])
- `zoom: number` - Nivel de zoom (default: 6)
- `height: string` - Altura del mapa (default: "400px")
- `showMarkers: boolean` - Mostrar marcadores en los puntos (default: true)
- `lineColor: string` - Color de la línea (default: "red")
- `lineWeight: number` - Grosor de la línea (default: 3)
- `lineOpacity: number` - Opacidad de la línea (default: 0.7)

#### Uso Básico

```svelte
<script>
  import Map from '$lib/components/Map.svelte';
  
  const route = [
    [40.4168, -3.7038], // Madrid
    [41.3851, 2.1734]   // Barcelona
  ];
</script>

<Map lineString={route} />
```

#### Uso Avanzado

```svelte
<script>
  import Map from '$lib/components/Map.svelte';
  
  const route = [
    [40.4168, -3.7038], // Madrid
    [41.3851, 2.1734]   // Barcelona
  ];
</script>

<Map 
  lineString={route}
  height="600px"
  lineColor="blue"
  lineWeight={5}
  lineOpacity={0.8}
  showMarkers={false}
  center={[40.9, -2.2]}
  zoom={7}
/>
```

## Funcionalidades de la Vista de Desarrollo

### Controles Interactivos

1. **Agregar Punto Aleatorio** - Añade un punto aleatorio a la ruta actual
2. **Restaurar Ruta Original** - Vuelve a la ruta de ejemplo inicial
3. **Mostrar/Ocultar Marcadores** - Alterna la visibilidad de los marcadores
4. **Cambiar Color** - Cambia el color de la línea (Rojo, Azul, Verde, Naranja)
5. **Cargar Rutas de Ejemplo**:
   - Madrid → Barcelona
   - Ruta Costera (Bilbao → A Coruña → Santiago → Sevilla)
   - Ruta Aleatoria (5 puntos aleatorios)

### Formato de LineString

El lineString debe ser un array de arrays, donde cada sub-array contiene dos números:
- Primer número: Latitud
- Segundo número: Longitud

```javascript
const lineString = [
  [40.4168, -3.7038], // Madrid
  [41.3851, 2.1734],  // Barcelona
  [37.3891, -5.9845]  // Sevilla
];
```

## Agregando Nuevos Componentes

Para agregar nuevos componentes a esta vista de desarrollo:

1. Crea tu componente en `$lib/components/`
2. Importa y usa el componente en `dev/+page.svelte`
3. Agrega controles interactivos para probar diferentes configuraciones
4. Documenta el uso del componente

## Dependencias

- **Leaflet** - Librería de mapas
- **@types/leaflet** - Tipos TypeScript para Leaflet
- **Tailwind CSS** - Framework de estilos

## Notas Técnicas

- El mapa usa OpenStreetMap como proveedor de tiles
- Los marcadores muestran las coordenadas al hacer clic
- El mapa se ajusta automáticamente para mostrar toda la ruta
- El componente es reactivo y se actualiza cuando cambian las props 

# Página de Desarrollo - Visualización de Rutas

Esta página permite visualizar las rutas optimizadas generadas por el sistema de optimización de flota.

## Funcionalidades

### 🚀 Carga Automática de Polylines (Nuevo)

La página ahora soporta carga automática de archivos polyline numerados:

- **Archivos esperados**: `polyline_001.json`, `polyline_002.json`, `polyline_003.json`, etc.
- **Ubicación**: `/ui/static/dev/` (accesible via `/dev/` en el navegador)
- **Carga automática**: El mapa busca automáticamente hasta 20 archivos numerados
- **Colores únicos**: Cada polyline tiene un color diferente de la paleta
- **Marcadores detallados**: Información completa de cada paso de la ruta

#### Estructura esperada de los archivos JSON:

```json
{
  "routes": [
    {
      "vehicle": 1,
      "cost": 1500,
      "duration": 3600,
      "route": [[lat1, lng1], [lat2, lng2], ...],
      "steps": [
        {
          "step_type": "start",
          "step_number": 0,
          "arrival": 0,
          "location": [lat, lng],
          "vehicle": 1,
          "order_refs": ["ORD001", "ORD002"]
        },
        {
          "step_type": "pickup",
          "step_number": 1,
          "arrival": 300,
          "location": [lat, lng],
          "vehicle": 1,
          "order_refs": ["ORD001"]
        }
      ]
    }
  ]
}
```

### 📁 Carga Manual (Legacy)

También soporta la carga manual desde un archivo único:

- **Archivo**: `/dev/polyline.json`
- **Formato**: Array de rutas con coordenadas y marcadores
- **Colores**: Asignados por vehículo

## Controles de la Interfaz

### Selector de Modo
- **Carga Automática**: Busca archivos numerados automáticamente
- **Carga Manual**: Carga desde el archivo polyline.json tradicional

### Información Visual
- **Rutas**: Líneas de colores en el mapa
- **Marcadores**: Puntos con información detallada
- **Popups**: Información completa al hacer clic en marcadores

## Tipos de Marcadores

- **▶ Inicio**: Verde - Punto de partida del vehículo
- **⏹️ Fin**: Rojo - Punto de llegada del vehículo  
- **📦 Recogida**: Color de la ruta - Punto de recogida
- **Números**: Color de la ruta - Puntos de entrega numerados

## Información en Popups

### Carga Automática
- Polyline y número de ruta
- Número de vehículo
- Tipo de paso y número
- Tiempo de llegada
- Referencias de órdenes asociadas

### Carga Manual
- Número de ruta
- Número de vehículo
- Tipo de paso
- Tiempo de llegada
- Referencias de órdenes

## Configuración del Mapa

- **Centro**: Santiago, Chile (-33.52245, -70.575)
- **Zoom**: 12
- **Grosor de línea**: 5px
- **Opacidad**: 0.7
- **Marcadores**: Habilitados por defecto

## Desarrollo

### Generación de Archivos Polyline

Los archivos se generan automáticamente en el backend cuando se ejecuta una optimización:

```go
// En app/adapter/out/vroom/vrp.go
polylineFilename := fmt.Sprintf("ui/static/dev/polyline_%03d.json", optimizationIndex+1)
individualVroomResponse.ExportToPolylineJSON(polylineFilename, fleetOptimization)
```

### Estructura de Archivos

```
ui/static/dev/
├── polyline_001.json     # Primera optimización
├── polyline_002.json     # Segunda optimización
├── polyline_003.json     # Tercera optimización
└── polyline.json         # Archivo legacy (carga manual)
```

## Troubleshooting

### No se ven las rutas
1. Verificar que los archivos JSON existen en `/ui/static/dev/`
2. Revisar la consola del navegador para errores
3. Confirmar que el formato JSON es correcto

### Marcadores no aparecen
1. Verificar que `showMarkers` está habilitado
2. Confirmar que los steps tienen coordenadas válidas
3. Revisar que el campo `location` existe en cada step

### Colores no se aplican
1. Verificar que el array `routeColors` está definido
2. Confirmar que cada polyline tiene un índice válido
3. Revisar que no hay conflictos con estilos CSS 