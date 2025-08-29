# Componentes Modulares

Este directorio contiene componentes modulares extraídos del componente principal `App.tsx` para mejorar la mantenibilidad y reutilización del código.

## CameraCapture

Componente para capturar fotos usando la cámara del dispositivo.

### Props

- `onCapture: (photoDataUrl: string) => void` - Callback que se ejecuta cuando se captura una foto
- `title?: string` - Título del componente (por defecto: "Capturar foto")
- `buttonText?: string` - Texto del botón de activar cámara (por defecto: "Activar cámara")
- `className?: string` - Clases CSS adicionales

### Uso

```tsx
import { CameraCapture } from './components'

<CameraCapture
  onCapture={(photoDataUrl) => console.log('Foto capturada:', photoDataUrl)}
  title="Foto de evidencia"
  buttonText="Activar cámara"
/>
```

### Características

- Captura automática al tocar la pantalla
- Flash visual al capturar
- Vibración del dispositivo (si está disponible)
- Cámara trasera por defecto
- Botón para cambiar foto
- Limpieza automática de recursos de cámara

### Dependencias

- `react-webcam`
- `useState`, `useEffect`, `useRef` de React

## Sidebar

Componente del menú lateral con funcionalidades de reportes y estado de conexión.

### Props

- `isOpen: boolean` - Controla si el sidebar está abierto
- `onClose: () => void` - Callback para cerrar el sidebar
- `routeStarted: boolean` - Indica si la ruta está iniciada
- `onDownloadReport: () => void` - Callback para descargar reporte
- `syncInfo?: { deviceId: string } | null` - Información de sincronización
- `markerPosition?: { visitIndex: number; coordinates: [number, number]; timestamp: number; deviceId: string } | null` - Posición del marcador sincronizado

### Uso

```tsx
import { Sidebar } from './components'

<Sidebar
  isOpen={sidebarOpen}
  onClose={() => setSidebarOpen(false)}
  routeStarted={routeStarted}
  onDownloadReport={openDownloadModal}
  syncInfo={syncInfo}
  markerPosition={markerPosition}
/>
```

### Características

- Menú deslizable con overlay
- Botón de descarga de reporte (solo cuando la ruta está iniciada)
- Indicadores de estado de conexión:
  - Estado de internet
  - Estado de sincronización GunJS
  - Indicador de marcador sincronizado
- Diseño responsive con gradientes y animaciones

### Dependencias

- `lucide-react` para iconos
- Tailwind CSS para estilos

## DeliveryModal

Modal para capturar evidencia de entrega exitosa.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Callback para cerrar el modal
- `onSubmit: (evidence: { recipientName: string; recipientRut: string; photoDataUrl: string }) => void` - Callback para enviar la evidencia
- `submitting?: boolean` - Indica si se está enviando la evidencia

### Uso

```tsx
import { DeliveryModal } from './components'

<DeliveryModal
  isOpen={deliveryModal.open}
  onClose={closeDeliveryModal}
  onSubmit={handleDeliverySubmit}
  submitting={submittingEvidence}
/>
```

### Características

- Formulario para nombre y RUT del receptor
- Integración con CameraCapture para foto de evidencia
- Validación de campos obligatorios
- Botones de cancelar y confirmar entrega
- Estado de envío para deshabilitar botones

### Dependencias

- `CameraCapture` component
- `useState`, `useRef` de React

## NonDeliveryModal

Modal para capturar evidencia de no entrega.

### Props

- `isOpen: boolean` - Controla si el modal está abierto
- `onClose: () => void` - Callback para cerrar el modal
- `onSubmit: (evidence: { reason: string; observations: string; photoDataUrl: string }) => void` - Callback para enviar la evidencia
- `submitting?: boolean` - Indica si se está enviando la evidencia

### Uso

```tsx
import { NonDeliveryModal } from './components'

<NonDeliveryModal
  isOpen={ndModal.open}
  onClose={closeNdModal}
  onSubmit={handleNonDeliverySubmit}
  submitting={submittingEvidence}
/>
```

### Características

- Campo de motivo con lista de sugerencias filtrable
- Campo de observaciones opcional
- Integración con CameraCapture para foto de evidencia
- Validación de motivo obligatorio y foto
- Botones de cancelar y confirmar no entrega
- Estado de envío para deshabilitar botones

### Dependencias

- `CameraCapture` component
- `useState`, `useRef` de React
