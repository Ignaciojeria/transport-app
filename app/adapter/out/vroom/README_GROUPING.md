# Agrupación de Referencias por Coordenadas

## Descripción

Esta funcionalidad permite agrupar referencias (`referenceID`) que tengan las mismas coordenadas en el polilyne de VROOM, evitando duplicados y mejorando la visualización de las rutas optimizadas.

## Problema Resuelto

Cuando tienes múltiples `referenceID` distintos pero con las mismas coordenadas (o coordenadas muy similares), el sistema ahora los agrupa en un solo punto del polilyne, mostrando todas las referencias asociadas a esa ubicación.

## Funciones Principales

### `groupReferencesByCoordinates(steps []Step, req optimization.FleetOptimization) map[string][]string`

Agrupa todas las referencias por coordenadas redondeadas a 6 decimales para evitar diferencias mínimas.

**Parámetros:**
- `steps`: Lista de pasos de la ruta de VROOM
- `req`: Solicitud de optimización original

**Retorna:**
- Mapa donde la clave es la coordenada (formato: "lat,lon") y el valor es la lista de referencias

### `getGroupedReferencesForStep(step Step, req optimization.FleetOptimization, allSteps []Step) []string`

Obtiene las referencias agrupadas para un step específico, considerando todas las referencias que comparten las mismas coordenadas.

**Parámetros:**
- `step`: Step específico para el cual obtener referencias
- `req`: Solicitud de optimización original
- `allSteps`: Todos los steps de la ruta

**Retorna:**
- Lista de referencias agrupadas para el step

### `removeDuplicateReferences(refs []string) []string`

Elimina referencias duplicadas de una lista.

### `removeDuplicateOrders(orders []domain.Order) []domain.Order`

Elimina órdenes duplicadas basándose en el `ReferenceID`.

## Ejemplo de Uso

```go
// Supongamos que tienes dos referencias con coordenadas similares:
// REF001: [-70.123456, -33.123456]
// REF002: [-70.123457, -33.123457]

// El sistema agrupará ambas referencias en el mismo punto del polilyne
// y mostrará: [REF001, REF002] en esa ubicación
```

## Implementación en las Funciones de Exportación

### ExportToPolylineJSON
- Usa `getGroupedReferencesForStep` para obtener referencias agrupadas
- Cada punto del polilyne muestra todas las referencias asociadas a esa coordenada

### ExportToGeoJSON
- Usa `getGroupedReferencesForStep` para obtener referencias agrupadas
- Los features de GeoJSON incluyen todas las referencias agrupadas en la propiedad `order_refs`

### Map (Conversión al Dominio)
- Agrupa órdenes por coordenadas usando `removeDuplicateOrders`
- Evita duplicados en el modelo de dominio

## Configuración

### Precisión de Coordenadas
La precisión se puede ajustar modificando el parámetro en `roundCoordinate`:

```go
// Actualmente configurado a 6 decimales
lat := roundCoordinate(step.Location[1], 6) // lat
lon := roundCoordinate(step.Location[0], 6) // lon
```

### Umbral de Agrupación
Para coordenadas muy cercanas pero no exactamente iguales, el sistema las redondea a 6 decimales, lo que significa que coordenadas que difieren en menos de 0.000001 grados se considerarán iguales.

## Beneficios

1. **Visualización Mejorada**: Evita puntos duplicados en el mapa
2. **Información Consolidada**: Muestra todas las referencias relacionadas en un solo punto
3. **Rendimiento**: Reduce la cantidad de puntos a renderizar
4. **Claridad**: Facilita la comprensión de las rutas optimizadas

## Casos de Uso

- **Múltiples entregas en el mismo edificio**: Diferentes referencias para el mismo destino
- **Coordenadas ligeramente diferentes**: Debido a precisiones de GPS o geocodificación
- **Agrupación lógica**: Referencias que deben tratarse como una sola parada

## Testing

Ejecuta las pruebas para verificar la funcionalidad:

```bash
go test ./app/adapter/out/vroom/model -v -run TestGroupReferencesByCoordinates
go test ./app/adapter/out/vroom/model -v -run TestGetGroupedReferencesForStep
``` 