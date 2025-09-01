# Componentes Modulares

Este directorio contiene componentes modulares extraídos del `App.tsx` principal para mejorar la organización y mantenibilidad del código.

## CameraCapture

Componente para capturar fotos usando la cámara del dispositivo.

### Props

- `onPhotoCapture: (photoDataUrl: string) => void` - Función que se ejecuta cuando se captura una foto

### Uso

```tsx
import { CameraCapture } from './components'

<CameraCapture onPhotoCapture={handlePhotoCapture} />
```

### Características

- Captura automática de fotos
- Confirmación automática de entrega al capturar
- Interfaz intuitiva y responsive
- Manejo de errores de cámara

### Dependencias

- `react-webcam` - Para acceso a la cámara
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## Sidebar

Componente de barra lateral que contiene el menú principal y funcionalidades adicionales.

### Props

- `onDownloadReport: () => void` - Función para abrir el modal de descarga de reporte
- `syncInfo?: { deviceId: string } | null` - Información de sincronización

### Uso

```tsx
import { Sidebar } from './components'

<Sidebar 
  onDownloadReport={openDownloadModal}
  syncInfo={syncInfo}
/>
```

### Características

- Menú de navegación principal
- Botón de descarga de reporte
- Indicador de estado de sincronización
- Diseño responsive

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## DeliveryModal

Componente modal para gestionar entregas exitosas con evidencia fotográfica.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onSubmit: (evidence: DeliveryEvidence) => void` - Función que se ejecuta al enviar la evidencia
- `submitting: boolean` - Indica si se está enviando la evidencia

### Uso

```tsx
import { DeliveryModal } from './components'

<DeliveryModal
  isOpen={evidenceModal.open}
  onClose={closeEvidenceModal}
  onSubmit={submitEvidence}
  submitting={submittingEvidence}
/>
```

### Características

- Captura de fotos con cámara
- Campos para razón y observaciones
- Validación de formulario
- Estado de carga durante envío

### Dependencias

- `CameraCapture` - Para captura de fotos
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## NonDeliveryModal

Componente modal para gestionar entregas no exitosas con evidencia fotográfica.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onSubmit: (evidence: NonDeliveryEvidence) => void` - Función que se ejecuta al enviar la evidencia
- `submitting: boolean` - Indica si se está enviando la evidencia

### Uso

```tsx
import { NonDeliveryModal } from './components'

<NonDeliveryModal
  isOpen={ndModal.open}
  onClose={closeNdModal}
  onSubmit={submitNonDelivery}
  submitting={submittingEvidence}
/>
```

### Características

- Captura de fotos con cámara
- Selección de razón de no entrega
- Campo para observaciones
- Validación de formulario

### Dependencias

- `CameraCapture` - Para captura de fotos
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCard

Componente principal para mostrar una visita individual con sus detalles y unidades de entrega.

### Props

- `visit: any` - Datos de la visita
- `visitIndex: number` - Índice de la visita en la lista
- `routeStarted: boolean` - Indica si la ruta ha comenzado
- `onCenterOnVisit: (visitIndex: number) => void` - Función para centrar el mapa en la visita
- `onOpenDelivery: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `onOpenNonDelivery: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de una unidad
- `shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean` - Función para determinar si mostrar la visita según el tab activo

### Uso

```tsx
import { VisitCard } from './components'

<VisitCard
  visit={visit}
  visitIndex={idx}
  routeStarted={routeStarted}
  onCenterOnVisit={() => {}}
  onOpenDelivery={openDeliveryFor}
  onOpenNonDelivery={openNonDeliveryFor}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
  shouldRenderByTab={shouldRenderByTab}
/>
```

### Características

- Header con información de contacto y dirección
- Lista de órdenes con unidades de entrega
- Botones de acción para gestionar entregas
- Filtrado por estado de entrega
- Integración con modales de entrega

### Dependencias

- `VisitCardHeader` - Para el encabezado de la visita
- `VisitCardOrders` - Para la lista de órdenes
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## NextVisitCard

Componente especializado para mostrar la "Siguiente Disponible" en la pestaña "En ruta".

### Props

- `nextVisit: any` - Datos de la siguiente visita disponible
- `nextIdx: number` - Índice de la siguiente visita
- `onCenterOnVisit: (visitIndex: number) => void` - Función para centrar el mapa en la visita

### Uso

```tsx
import { NextVisitCard } from './components'

<NextVisitCard
  nextVisit={nextVisit}
  nextIdx={nextIdx}
  onCenterOnVisit={() => {}}
/>
```

### Características

- Diseño destacado para la siguiente visita
- Botón para centrar en el mapa
- Información resumida de la visita

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCardHeader

Componente para mostrar el encabezado de una visita con información de contacto y dirección.

### Props

- `visit: any` - Datos de la visita
- `visitIndex: number` - Índice de la visita en la lista
- `onCenterOnVisit: (visitIndex: number) => void` - Función para centrar el mapa en la visita

### Uso

```tsx
import { VisitCardHeader } from './components'

<VisitCardHeader
  visit={visit}
  visitIndex={visitIndex}
  onCenterOnVisit={onCenterOnVisit}
/>
```

### Características

- Número de secuencia de la visita
- Nombre del contacto
- Dirección de la visita
- Botón para centrar en el mapa

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCardOrders

Componente para mostrar la lista de órdenes de una visita.

### Props

- `visit: any` - Datos de la visita
- `visitIndex: number` - Índice de la visita en la lista
- `routeStarted: boolean` - Indica si la ruta ha comenzado
- `onOpenDelivery: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `onOpenNonDelivery: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de una unidad

### Uso

```tsx
import { VisitCardOrders } from './components'

<VisitCardOrders
  visit={visit}
  visitIndex={visitIndex}
  routeStarted={routeStarted}
  onOpenDelivery={onOpenDelivery}
  onOpenNonDelivery={onOpenNonDelivery}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
/>
```

### Características

- Lista de órdenes con referencia ID
- Unidades de entrega con estado
- Botones de acción para gestionar entregas
- Integración con modales

### Dependencias

- `VisitCardDeliveryUnit` - Para unidades individuales
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCardDeliveryUnit

Componente para mostrar una unidad de entrega individual con su estado y acciones.

### Props

- `unit: any` - Datos de la unidad de entrega
- `unitIndex: number` - Índice de la unidad en la orden
- `orderIndex: number` - Índice de la orden en la visita
- `visitIndex: number` - Índice de la visita en la lista
- `status?: 'delivered' | 'not-delivered'` - Estado actual de la unidad
- `routeStarted: boolean` - Indica si la ruta ha comenzado
- `onOpenDelivery: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `onOpenNonDelivery: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega

### Uso

```tsx
import { VisitCardDeliveryUnit } from './components'

<VisitCardDeliveryUnit
  unit={unit}
  unitIndex={uIdx}
  orderIndex={orderIndex}
  visitIndex={visitIndex}
  status={status}
  routeStarted={routeStarted}
  onOpenDelivery={onOpenDelivery}
  onOpenNonDelivery={onOpenNonDelivery}
/>
```

### Características

- Información detallada de la unidad
- Estado visual con colores
- Botones de acción según el estado
- Manejo de peso, volumen y cantidad

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## DownloadReportModal

Componente modal para seleccionar el formato de descarga del reporte.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onDownloadReport: () => void` - Función que se ejecuta al confirmar la descarga

### Uso

```tsx
import { DownloadReportModal } from './components'

<DownloadReportModal
  isOpen={downloadModal}
  onClose={closeDownloadModal}
  onDownloadReport={downloadReport}
/>
```

### Características

- Selección entre formato CSV y Excel
- Botones de descarga
- Interfaz simple y clara

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## RouteStartModal

Componente modal para solicitar la patente del vehículo al iniciar una ruta.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onConfirm: (license: string) => void` - Función que se ejecuta al confirmar la patente
- `defaultLicense?: string` - Patente sugerida desde los datos de la ruta

### Uso

```tsx
import { RouteStartModal } from './components'

<RouteStartModal
  isOpen={routeStartModal}
  onClose={() => setRouteStartModal(false)}
  onConfirm={handleLicenseConfirm}
  defaultLicense={routeData?.vehicle?.plate}
/>
```

### Características

- Sugerencia de patente asignada a la ruta
- Opción para usar patente sugerida o ingresar otra
- Validación de entrada (máximo 8 caracteres)
- Conversión automática a mayúsculas
- Soporte para confirmación con Enter
- Focus automático en el input al abrir

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos
- React hooks (`useRef`, `useEffect`) - Para manejo del input y focus

## VisitTabs

Componente para los tabs de navegación entre estados de visitas.

### Props

- `activeTab: 'en-ruta' | 'entregados' | 'no-entregados'` - Tab activo actualmente
- `onTabChange: (tab: 'en-ruta' | 'entregados' | 'no-entregados') => void` - Función para cambiar de tab
- `inRouteUnits: number` - Cantidad de unidades en ruta
- `deliveredUnits: number` - Cantidad de unidades entregadas
- `notDeliveredUnits: number` - Cantidad de unidades no entregadas

### Uso

```tsx
import { VisitTabs } from './components'

<VisitTabs
  activeTab={activeTab}
  onTabChange={setActiveTab}
  inRouteUnits={inRouteUnits.length}
  deliveredUnits={deliveredUnits.length}
  notDeliveredUnits={notDeliveredUnits.length}
/>
```

### Características

- Tabs sticky con backdrop blur
- Contadores dinámicos para cada estado
- Iconos representativos para cada tab
- Transiciones suaves y estados hover
- Diseño responsive y accesible

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## MapView

Componente principal para la vista del mapa con toda la funcionalidad de navegación y gestión de visitas.

### Props

- `routeId: string` - ID de la ruta actual
- `routeData: any` - Datos completos de la ruta
- `visits: any[]` - Lista de visitas de la ruta
- `routeStarted: boolean` - Indica si la ruta ha comenzado
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de una unidad
- `getNextPendingVisitIndex: () => number | null` - Función para obtener el índice de la siguiente visita pendiente
- `getPositionedVisitIndex: () => number | null` - Función para obtener el índice de la visita posicionada
- `nextVisitIndex: number | null` - Índice de la siguiente visita
- `lastCenteredVisit: number | null` - Índice de la última visita centrada
- `markerPosition: any` - Información de posición del marcador
- `openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega
- `onDownloadReport: () => void` - Función para descargar reporte
- `setNextVisitIndex: (index: number | null) => void` - Función para establecer el índice de la siguiente visita
- `setLastCenteredVisit: (index: number | null) => void` - Función para establecer el índice de la última visita centrada
- `setMarkerPosition: (routeId: string, visitIndex: number, coordinates: [number, number]) => Promise<void>` - Función para establecer la posición del marcador
- `openNextNavigation: (provider: 'google' | 'waze' | 'geo') => void` - Función para abrir navegación externa

### Uso

```tsx
import { MapView } from './components'

<MapView
  routeId={routeId}
  routeData={routeData}
  visits={visits}
  routeStarted={routeStarted}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
  getNextPendingVisitIndex={getNextPendingVisitIndex}
  getPositionedVisitIndex={getPositionedVisitIndex}
  nextVisitIndex={nextVisitIndex}
  lastCenteredVisit={lastCenteredVisit}
  markerPosition={markerPosition}
  openDeliveryFor={openDeliveryFor}
  openNonDeliveryFor={openNonDeliveryFor}
  onDownloadReport={openDownloadModal}
  setNextVisitIndex={setNextVisitIndex}
  setLastCenteredVisit={setLastCenteredVisit}
  setMarkerPosition={setMarkerPosition}
  openNextNavigation={openNextNavigation}
/>
```

### Características

- Mapa interactivo con Leaflet
- Marcadores de visitas con estados visuales
- Controles de GPS y navegación
- Card de visita integrada
- Sincronización entre dispositivos
- Gestión completa de entregas desde el mapa

### Dependencias

- `MapControls` - Para controles del mapa
- `MapVisitCard` - Para mostrar la visita en modo mapa
- `MapView.utils` - Para funciones utilitarias
- Leaflet - Para el mapa
- Tailwind CSS - Para estilos

## MapControls

Componente para los controles flotantes del mapa.

### Props

- `gpsActive: boolean` - Indica si el GPS está activo
- `onGPSToggle: () => void` - Función para activar/desactivar GPS
- `onZoomToSelected: () => void` - Función para hacer zoom al punto seleccionado
- `onNavigate: (provider: 'google' | 'waze' | 'geo') => void` - Función para abrir navegación externa

### Uso

```tsx
import { MapControls } from './components'

<MapControls
  gpsActive={gpsActive}
  onGPSToggle={gpsActive ? stopGPS : startGPS}
  onZoomToSelected={zoomToCurrentlySelected}
  onNavigate={openNextNavigation}
/>
```

### Características

- Botón de GPS del conductor
- Botón de zoom al punto seleccionado
- Botones de navegación (Google Maps, Waze)
- Diseño flotante y responsive

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## MapVisitCard

Componente para mostrar la card de visita en modo mapa.

### Props

- `visit: any` - Datos de la visita
- `displayIdx: number` - Índice de la visita para mostrar
- `isSelectedVisit: boolean` - Indica si es la visita seleccionada
- `isProcessed: boolean` - Indica si la visita ya está procesada
- `hasNextPending: boolean` - Indica si hay siguiente visita pendiente
- `nextPendingIdx: number | null` - Índice de la siguiente visita pendiente
- `routeStarted: boolean` - Indica si la ruta ha comenzado
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de una unidad
- `openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega
- `onNextPending: (nextPendingIdx: number) => void` - Función para ir a la siguiente visita pendiente
- `onClearSelection: () => void` - Función para limpiar la selección

### Uso

```tsx
import { MapVisitCard } from './components'

<MapVisitCard
  visit={visit}
  displayIdx={displayIdx}
  isSelectedVisit={isSelectedVisit}
  isProcessed={isProcessed}
  hasNextPending={hasNextPending}
  nextPendingIdx={nextPendingIdx}
  routeStarted={routeStarted}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
  openDeliveryFor={openDeliveryFor}
  openNonDeliveryFor={openNonDeliveryFor}
  onNextPending={handleNextPending}
  onClearSelection={handleClearSelection}
/>
```

### Características

- Información completa de la visita
- Gestión de entregas desde el mapa
- Navegación a siguiente visita pendiente
- Estados visuales claros
- Integración con modales

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## MapView.utils

Archivo de utilidades para la vista del mapa con funciones auxiliares.

### Funciones Exportadas

- `getLatLngFromAddressInfo(addr: any): [number, number] | null` - Extrae coordenadas de addressInfo
- `decodePolyline(encoded: string): Array<[number, number]>` - Decodifica polylines de Google
- `getGradientColor(baseColor: string): string` - Obtiene color de gradiente complementario
- `getVisitStatus(visit: any, getDeliveryUnitStatus: Function, visitIndex: number): 'completed' | 'not-delivered' | 'partial' | 'pending'` - Determina el estado de una visita
- `getVisitMarkerColor(visitStatus: string): string` - Obtiene el color del marcador según el estado

### Uso

```tsx
import { 
  getLatLngFromAddressInfo, 
  decodePolyline, 
  getGradientColor, 
  getVisitStatus, 
  getVisitMarkerColor 
} from './MapView.utils'
```

### Características

- Funciones puras y reutilizables
- Manejo de coordenadas geográficas
- Decodificación de polylines
- Lógica de estados de visitas
- Colores de marcadores

### Dependencias

- Ninguna (funciones puras)
