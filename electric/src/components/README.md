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

- `isOpen: boolean` - Controla si el sidebar está abierto
- `onClose: () => void` - Función para cerrar el sidebar
- `onDownloadReport: () => void` - Función para abrir el modal de descarga
- `syncInfo?: { deviceId: string } | null` - Información de sincronización
- `routeStarted: boolean` - Estado de la ruta

### Uso

```tsx
import { Sidebar } from './components'

<Sidebar
  isOpen={sidebarOpen}
  onClose={() => setSidebarOpen(false)}
  onDownloadReport={() => setDownloadModal(true)}
  syncInfo={syncInfo}
  routeStarted={routeStarted}
/>
```

### Características

- Menú de navegación principal
- Indicador de estado de sincronización
- Botón de descarga de reportes
- Diseño responsive y accesible

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## DeliveryModal

Componente modal para confirmar entregas con evidencia fotográfica y datos del destinatario.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onSubmit: (evidence: DeliveryEvidence) => void` - Función que se ejecuta al confirmar
- `submitting: boolean` - Indica si se está procesando la entrega

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

- Captura de foto del destinatario
- Campos para nombre y RUT del destinatario
- Validación de datos requeridos
- Estado de carga durante el envío

### Dependencias

- `CameraCapture` - Para captura de fotos
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## NonDeliveryModal

Componente modal para registrar no entregas con motivo y evidencia fotográfica.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onSubmit: (evidence: NonDeliveryEvidence) => void` - Función que se ejecuta al confirmar
- `submitting: boolean` - Indica si se está procesando

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

- Selección de motivo de no entrega
- Campo de observaciones opcional
- Captura de foto de evidencia
- Validación de datos requeridos

### Dependencias

- `CameraCapture` - Para captura de fotos
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCard

Componente principal para mostrar una visita individual con sus órdenes y unidades de entrega.

### Props

- `visit: any` - Datos de la visita
- `visitIndex: number` - Índice de la visita
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de entrega
- `shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean` - Función para filtrar por pestaña
- `openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega
- `centerMapOnVisit: (visitIndex: number) => void` - Función para centrar mapa en la visita

### Uso

```tsx
import { VisitCard } from './components'

<VisitCard
  visit={visit}
  visitIndex={index}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
  shouldRenderByTab={shouldRenderByTab}
  openDeliveryFor={openDeliveryFor}
  openNonDeliveryFor={openNonDeliveryFor}
  centerMapOnVisit={centerMapOnVisit}
/>
```

### Características

- Renderizado condicional por pestaña activa
- Composición de sub-componentes especializados
- Manejo de estados de entrega
- Integración con funcionalidades del mapa

### Dependencias

- `VisitCardHeader` - Para el encabezado de la visita
- `VisitCardOrders` - Para la lista de órdenes
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## NextVisitCard

Componente especializado para mostrar la "Siguiente Disponible" en la pestaña "en-ruta".

### Props

- `nextVisit: any` - Datos de la siguiente visita disponible
- `nextVisitIndex: number` - Índice de la siguiente visita
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de entrega
- `openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega
- `centerMapOnVisit: (visitIndex: number) => void` - Función para centrar mapa en la visita

### Uso

```tsx
import { NextVisitCard } from './components'

<NextVisitCard
  nextVisit={nextVisit}
  nextVisitIndex={nextVisitIndex}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
  openDeliveryFor={openDeliveryFor}
  openNonDeliveryFor={openNonDeliveryFor}
  centerMapOnVisit={centerMapOnVisit}
/>
```

### Características

- Diseño destacado para la siguiente visita
- Indicador visual de prioridad
- Misma funcionalidad que VisitCard pero con estilo diferenciado

### Dependencias

- `VisitCardHeader` - Para el encabezado de la visita
- `VisitCardOrders` - Para la lista de órdenes
- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCardHeader

Componente para mostrar el encabezado de una visita con información básica.

### Props

- `visit: any` - Datos de la visita
- `visitIndex: number` - Índice de la visita
- `centerMapOnVisit: (visitIndex: number) => void` - Función para centrar mapa en la visita

### Uso

```tsx
import { VisitCardHeader } from './components'

<VisitCardHeader
  visit={visit}
  visitIndex={index}
  centerMapOnVisit={centerMapOnVisit}
/>
```

### Características

- Número de secuencia de la visita
- Nombre del contacto
- Dirección de la visita
- Botón para centrar mapa

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## VisitCardOrders

Componente para mostrar la lista de órdenes de una visita.

### Props

- `visit: any` - Datos de la visita
- `visitIndex: number` - Índice de la visita
- `getDeliveryUnitStatus: (visitIndex: number, orderIndex: number, unitIndex: number) => 'delivered' | 'not-delivered' | undefined` - Función para obtener el estado de entrega
- `shouldRenderByTab: (status?: 'delivered' | 'not-delivered') => boolean` - Función para filtrar por pestaña
- `openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega

### Uso

```tsx
import { VisitCardOrders } from './components'

<VisitCardOrders
  visit={visit}
  visitIndex={index}
  getDeliveryUnitStatus={getDeliveryUnitStatus}
  shouldRenderByTab={shouldRenderByTab}
  openDeliveryFor={openDeliveryFor}
  openNonDeliveryFor={openNonDeliveryFor}
/>
```

### Características

- Iteración sobre órdenes de la visita
- Filtrado por estado de entrega según pestaña
- Delegación de renderizado a VisitCardDeliveryUnit

### Dependencias

- `VisitCardDeliveryUnit` - Para unidades de entrega individuales
- Tailwind CSS - Para estilos

## VisitCardDeliveryUnit

Componente para mostrar una unidad de entrega individual con su estado y acciones.

### Props

- `unit: any` - Datos de la unidad de entrega
- `visitIndex: number` - Índice de la visita
- `orderIndex: number` - Índice de la orden
- `unitIndex: number` - Índice de la unidad
- `status: 'delivered' | 'not-delivered' | undefined` - Estado actual de la entrega
- `openDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de entrega
- `openNonDeliveryFor: (visitIndex: number, orderIndex: number, unitIndex: number) => void` - Función para abrir modal de no entrega

### Uso

```tsx
import { VisitCardDeliveryUnit } from './components'

<VisitCardDeliveryUnit
  unit={unit}
  visitIndex={visitIndex}
  orderIndex={orderIndex}
  unitIndex={unitIndex}
  status={status}
  openDeliveryFor={openDeliveryFor}
  openNonDeliveryFor={openNonDeliveryFor}
/>
```

### Características

- Información detallada de la unidad (descripción, peso, volumen, cantidad)
- Indicador visual del estado de entrega
- Botones de acción según el estado
- Colores diferenciados por estado

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos

## DownloadReportModal

Componente modal para descargar reportes de ruta en formato CSV o Excel.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Función para cerrar el modal
- `onDownloadReport: (format: 'csv' | 'excel') => void` - Función que se ejecuta al seleccionar un formato

### Uso

```tsx
import { DownloadReportModal } from './components'

<DownloadReportModal
  isOpen={downloadModal}
  onClose={closeDownloadModal}
  onDownloadReport={handleDownload}
/>
```

### Características

- Interfaz intuitiva para seleccionar formato de descarga
- Opciones para CSV y Excel
- Diseño responsive y accesible
- Feedback visual para cada opción

### Dependencias

- `lucide-react` - Para iconos
- Tailwind CSS - Para estilos
- Utilidades de reporte para generación de archivos

### Utilidades Asociadas

El componente utiliza funciones utilitarias para generar reportes:

- `generateReportData()` - Prepara los datos del reporte
- `generateCSVContent()` - Genera contenido CSV
- `generateExcelContent()` - Genera contenido Excel
- `downloadFile()` - Maneja la descarga del archivo

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

### Funcionalidades

- **Patente sugerida**: Muestra la patente asignada a la ruta si está disponible
- **Patente personalizada**: Permite ingresar una patente diferente
- **Validación**: Solo permite confirmar si hay una patente ingresada
- **UX mejorada**: Focus automático y soporte para tecla Enter

## VisitTabs

Componente para mostrar los tabs de navegación entre diferentes estados de visitas (En ruta, Entregados, No entregados).

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

- `lucide-react` - Para iconos (Play, CheckCircle, XCircle)
- Tailwind CSS - Para estilos y animaciones

### Funcionalidades

- **Navegación por tabs**: Cambio entre diferentes estados de visitas
- **Contadores en tiempo real**: Muestra cantidad de unidades en cada estado
- **Estados visuales**: Diferencia visual entre tab activo e inactivo
- **Iconos descriptivos**: Cada tab tiene un icono representativo
- **Posicionamiento sticky**: Los tabs permanecen visibles al hacer scroll
