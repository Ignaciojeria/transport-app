/**
 * Utilidades para manejo de imágenes y subida a servidor
 */

// Función para convertir imagen a formato WebP con manejo robusto de errores
export async function convertToWebP(imageDataUrl: string, quality: number = 0.8): Promise<Blob> {
  return new Promise((resolve, reject) => {
    try {
      // Validar que la imagen sea válida
      if (!imageDataUrl || typeof imageDataUrl !== 'string') {
        reject(new Error('URL de imagen inválida'))
        return
      }

      // Crear canvas y contexto
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      
      if (!ctx) {
        reject(new Error('No se pudo obtener el contexto del canvas'))
        return
      }
      
      // Crear imagen
      const img = new Image()
      
      // Configurar eventos de la imagen
      img.onload = () => {
        try {
          // Validar dimensiones de la imagen
          if (img.width === 0 || img.height === 0) {
            reject(new Error('La imagen no tiene dimensiones válidas'))
            return
          }
          
          // Configurar dimensiones del canvas
          canvas.width = img.width
          canvas.height = img.height
          
          // Limpiar canvas
          ctx.clearRect(0, 0, canvas.width, canvas.height)
          
          // Dibujar imagen en el canvas
          ctx.drawImage(img, 0, 0)
          
          // Intentar convertir a WebP
          canvas.toBlob(
            (blob) => {
              if (blob) {
                console.log('✅ Imagen convertida a WebP exitosamente:', {
                  size: blob.size,
                  type: blob.type,
                  dimensions: `${canvas.width}x${canvas.height}`
                })
                resolve(blob)
              } else {
                reject(new Error('No se pudo generar el blob WebP'))
              }
            },
            'image/webp',
            quality
          )
        } catch (error) {
          reject(new Error(`Error procesando imagen: ${error instanceof Error ? error.message : 'Error desconocido'}`))
        }
      }
      
      img.onerror = (error) => {
        console.error('❌ Error cargando imagen:', error)
        reject(new Error('No se pudo cargar la imagen desde la URL proporcionada'))
      }
      
      // Configurar crossOrigin para evitar problemas CORS
      img.crossOrigin = 'anonymous'
      
      // Cargar imagen
      img.src = imageDataUrl
      
    } catch (error) {
      reject(new Error(`Error inesperado: ${error instanceof Error ? error.message : 'Error desconocido'}`))
    }
  })
}

// Función para subir imagen usando uploadUrl del contrato
export async function uploadImageToServer(
  imageBlob: Blob, 
  uploadUrl: string, 
  contentType: string = 'image/webp'
): Promise<string> {
  try {
    // Validar parámetros
    if (!imageBlob || !(imageBlob instanceof Blob)) {
      throw new Error('Blob de imagen inválido')
    }
    
    if (!uploadUrl || typeof uploadUrl !== 'string') {
      throw new Error('URL de subida inválida')
    }
    
    console.log('📤 Iniciando subida de imagen:', {
      size: imageBlob.size,
      type: imageBlob.type,
      uploadUrl: uploadUrl.substring(0, 100) + '...'
    })
    
    const response = await fetch(uploadUrl, {
      method: 'PUT',
      body: imageBlob,
      headers: {
        'Content-Type': contentType,
      },
    })
    
    if (!response.ok) {
      const errorText = await response.text().catch(() => 'Sin detalles del error')
      throw new Error(`Error al subir imagen: ${response.status} ${response.statusText}. Detalles: ${errorText}`)
    }
    
    // Las URLs son pre-firmadas del contrato, no se deben modificar
    // Si no se proporciona downloadUrl, usar la uploadUrl sin parámetros como fallback
    const downloadUrl = uploadUrl.split('?')[0] // Fallback: remover parámetros de firma
    
    console.log('✅ Imagen subida exitosamente. URL de descarga:', downloadUrl)
    return downloadUrl
    
  } catch (error) {
    console.error('❌ Error subiendo imagen:', error)
    throw error
  }
}

// Función principal que convierte y sube la imagen
export async function processAndUploadImage(
  imageDataUrl: string, 
  uploadUrl: string, 
  downloadUrl?: string,
  quality: number = 0.8
): Promise<{ downloadUrl: string; blob: Blob }> {
  try {
    console.log('🔄 Iniciando procesamiento de imagen...')
    
    // Validar parámetros
    if (!imageDataUrl || typeof imageDataUrl !== 'string') {
      throw new Error('URL de imagen inválida')
    }
    
    if (!uploadUrl || typeof uploadUrl !== 'string') {
      throw new Error('URL de subida inválida')
    }
    
    if (quality < 0.1 || quality > 1.0) {
      console.warn('⚠️ Calidad fuera de rango (0.1-1.0), usando valor por defecto 0.8')
      quality = 0.8
    }
    
    console.log('🔄 Convirtiendo imagen a WebP...')
    const webpBlob = await convertToWebP(imageDataUrl, quality)
    
    console.log('📤 Subiendo imagen WebP al servidor...', {
      blobType: webpBlob.type,
      blobSize: webpBlob.size
    })
    await uploadImageToServer(webpBlob, uploadUrl, 'image/webp')
    
    // Usar la downloadUrl pre-firmada del contrato tal como viene
    const finalDownloadUrl = downloadUrl || uploadUrl.split('?')[0]
    
    console.log('✅ Imagen procesada y subida exitosamente:', finalDownloadUrl)
    return { downloadUrl: finalDownloadUrl, blob: webpBlob }
    
  } catch (error) {
    console.error('❌ Error procesando/subiendo imagen:', error)
    throw error
  }
}

// Función para obtener la URL de subida y descarga desde el contrato de ruta
export function getUploadUrlFromRoute(
  routeData: any, 
  visitIndex: number, 
  orderIndex: number, 
  unitIndex: number
): { uploadUrl: string | null; downloadUrl: string | null } {
  try {
    // Validar parámetros
    if (!routeData || typeof visitIndex !== 'number' || typeof orderIndex !== 'number' || typeof unitIndex !== 'number') {
      console.warn('⚠️ Parámetros inválidos para getUploadUrlFromRoute')
      return { uploadUrl: null, downloadUrl: null }
    }
    
    const visit = routeData?.visits?.[visitIndex]
    const order = visit?.orders?.[orderIndex]
    const deliveryUnit = order?.deliveryUnits?.[unitIndex]
    
    // Buscar en evidences del delivery unit
    if (deliveryUnit?.evidences && Array.isArray(deliveryUnit.evidences)) {
      const evidence = deliveryUnit.evidences[0] // Tomar la primera evidencia disponible
      if (evidence?.uploadUrl || evidence?.downloadUrl) {
        console.log('✅ URLs encontradas en deliveryUnit.evidences:', {
          uploadUrl: !!evidence.uploadUrl,
          downloadUrl: !!evidence.downloadUrl
        })
        return {
          uploadUrl: evidence.uploadUrl || null,
          downloadUrl: evidence.downloadUrl || null
        }
      }
    }
    
    // Fallback: buscar en el nivel de order
    if (order?.evidences && Array.isArray(order.evidences)) {
      const evidence = order.evidences[0]
      if (evidence?.uploadUrl || evidence?.downloadUrl) {
        console.log('✅ URLs encontradas en order.evidences:', {
          uploadUrl: !!evidence.uploadUrl,
          downloadUrl: !!evidence.downloadUrl
        })
        return {
          uploadUrl: evidence.uploadUrl || null,
          downloadUrl: evidence.downloadUrl || null
        }
      }
    }
    
    // Fallback: buscar en el nivel de visit
    if (visit?.evidences && Array.isArray(visit.evidences)) {
      const evidence = visit.evidences[0]
      if (evidence?.uploadUrl || evidence?.downloadUrl) {
        console.log('✅ URLs encontradas en visit.evidences:', {
          uploadUrl: !!evidence.uploadUrl,
          downloadUrl: !!evidence.downloadUrl
        })
        return {
          uploadUrl: evidence.uploadUrl || null,
          downloadUrl: evidence.downloadUrl || null
        }
      }
    }
    
    console.warn('⚠️ No se encontraron URLs en el contrato de ruta')
    return { uploadUrl: null, downloadUrl: null }
    
  } catch (error) {
    console.error('❌ Error obteniendo URLs:', error)
    return { uploadUrl: null, downloadUrl: null }
  }
}

// Función de utilidad para verificar soporte de WebP
export function isWebPSupported(): boolean {
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    if (!ctx) return false
    
    // Intentar crear un blob WebP
    let supported = false
    canvas.toBlob((blob) => {
      if (blob && blob.type === 'image/webp') {
        supported = true
      }
    }, 'image/webp', 0.8)
    
    return supported
  } catch {
    return false
  }
}

// Función de utilidad para obtener información de la imagen
export function getImageInfo(imageDataUrl: string): Promise<{ width: number; height: number; size: number }> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => {
      resolve({
        width: img.width,
        height: img.height,
        size: imageDataUrl.length // Tamaño aproximado de la data URL
      })
    }
    img.onerror = () => reject(new Error('No se pudo cargar la imagen'))
    img.src = imageDataUrl
  })
}
