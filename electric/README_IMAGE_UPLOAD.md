# Funcionalidad de Subida de Imágenes WebP

## Descripción

Esta funcionalidad permite capturar fotos de evidencia directamente en formato WebP y subirlas automáticamente al servidor usando las URLs de subida (`uploadUrl`) del contrato de rutas. Las URLs de las fotos se incluyen en los reportes generados.

## Características

- **Captura directa en WebP**: Las fotos se capturan directamente en formato WebP desde la cámara
- **Subida automática**: Las fotos se suben automáticamente al servidor después de ser capturadas
- **Integración con contratos**: Utiliza las URLs de subida y descarga del contrato de rutas
- **Reportes actualizados**: Los reportes CSV y Excel incluyen tanto las URLs de las fotos como las URLs de descarga del contrato
- **Procesamiento en tiempo real**: Indicadores visuales durante la captura y conversión
- **Doble URL en reportes**: Cada evidencia incluye tanto la foto local como la URL del contrato para acceso directo

## Componentes Modificados

### 1. `CameraCapture.tsx`
- **Captura directa en WebP**: Usa Canvas API para capturar frames del video y convertirlos a WebP
- **Configurable la calidad**: Control de calidad de imagen (0.1 a 1.0)
- **Indicadores de estado**: Muestra cuando se está procesando la imagen
- **Manejo de errores mejorado**: Mensajes específicos para diferentes tipos de errores
- **Logs de depuración**: Información detallada en consola para diagnóstico

### 2. `DeliveryModal.tsx`
- Integrado con la funcionalidad de subida de imágenes
- Muestra estado de subida en tiempo real
- Valida que la foto esté subida antes de permitir guardar

### 3. `NonDeliveryModal.tsx`
- Misma funcionalidad que DeliveryModal
- Campos mejorados para motivos de no entrega
- Validación de fotos subidas

### 4. `imageUpload.ts` (Nuevo)
- Utilidades para conversión a WebP (fallback)
- Funciones de subida al servidor
- Obtención de URLs de subida desde contratos
- Funciones de utilidad para diagnóstico

## Flujo de Trabajo

1. **Activación de Cámara**: El usuario activa la cámara desde el botón
2. **Captura de Foto**: Al tocar la pantalla, se captura el frame actual del video
3. **Conversión Directa a WebP**: El frame se convierte directamente a WebP usando Canvas API
4. **Subida Automática**: La foto se sube al servidor usando el `uploadUrl` del contrato
5. **Validación**: Solo se permite guardar la evidencia si la foto está subida
6. **Reporte**: Se incluyen dos URLs en los reportes:
   - **URL_Foto_Evidencia**: URL de la foto capturada localmente
   - **URL_Descarga_Contrato**: URL de descarga del contrato de ruta para acceso directo

## Estructura de Datos

### Contrato de Ruta
```typescript
interface Evidence {
  uploadUrl: string    // URL para subir la foto
  downloadUrl: string  // URL para descargar la foto
}
```

### Evidencia de Entrega
```typescript
interface EvidencePhoto {
  takenAt: string      // Fecha de captura
  type: string         // Tipo de evidencia
  url: string          // URL de la foto subida
}

interface Evidence {
  uploadUrl: string    // URL para subir la foto
  downloadUrl: string  // URL para descargar la foto desde el contrato
}
```

## Configuración

### Calidad de Imagen
```typescript
// En CameraCapture
quality={0.8} // Valor entre 0.1 y 1.0
```

### URLs de Subida y Descarga
Las URLs se obtienen automáticamente del contrato de ruta en este orden:
1. `deliveryUnit.evidences[0].uploadUrl` y `downloadUrl`
2. `order.evidences[0].uploadUrl` y `downloadUrl`
3. `visit.evidences[0].uploadUrl` y `downloadUrl`

**Nota**: Tanto `uploadUrl` como `downloadUrl` se extraen del mismo nivel del contrato para mantener consistencia.

## Reportes

### CSV
- Nueva columna: `URL_Foto_Evidencia` - URL de la foto capturada localmente
- Nueva columna: `URL_Descarga_Contrato` - URL de descarga del contrato de ruta
- Incluye ambas URLs para cada unidad de entrega

### Excel
- Mismas columnas: `URL_Foto_Evidencia` y `URL_Descarga_Contrato`
- Formato HTML con estilos para mejor visualización

## Manejo de Errores

- **Cámara no disponible**: Muestra mensaje de error específico
- **Error de conversión WebP**: Muestra detalle del error y permite reintentar
- **Error de subida**: Muestra detalle del error y permite reintentar
- **URL no encontrada**: Valida que exista `uploadUrl` en el contrato
- **Fallback**: Si no hay URL de subida, mantiene la funcionalidad local

## Dependencias

- `react-webcam`: Para acceso a la cámara
- Canvas API: Para captura y conversión a WebP
- Fetch API: Para subida de archivos

## Consideraciones Técnicas

- **Formato WebP**: Mejor compresión que JPEG, manteniendo calidad
- **Captura directa**: No hay conversión intermedia, mejor rendimiento
- **Subida PUT**: Las URLs de subida usan método HTTP PUT
- **Descarga GET**: Las URLs de descarga usan método HTTP GET para acceso directo
- **Validación**: Se requiere foto subida para guardar evidencia
- **Estado Reactivo**: UI actualizada en tiempo real durante captura y subida
- **Doble URL**: Los reportes incluyen tanto la foto local como la URL del contrato para máxima flexibilidad

## Logs de Depuración

La funcionalidad incluye logs detallados en la consola:

```javascript
// Al capturar
📸 Capturando imagen: {
  videoWidth: 1280,
  videoHeight: 720,
  canvasWidth: 1280,
  canvasHeight: 720,
  quality: 0.8
}

// Al completar
✅ Imagen WebP capturada exitosamente: {
  size: 45678,
  type: "image/webp",
  dataUrlLength: 123456
}
```

## Próximas Mejoras

- [ ] Soporte para múltiples fotos por evidencia
- [ ] Compresión adicional de imágenes
- [ ] Cache local de fotos subidas
- [ ] Sincronización offline
- [ ] Métricas de uso de almacenamiento
- [ ] Previsualización en tiempo real antes de capturar
