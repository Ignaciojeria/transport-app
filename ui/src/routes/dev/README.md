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