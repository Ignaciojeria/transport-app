/**
 * Redimensiona una imagen a un tamaño máximo manteniendo la proporción
 * @param file - Archivo de imagen
 * @param maxWidth - Ancho máximo en píxeles (por defecto 400px para menús)
 * @param maxHeight - Alto máximo en píxeles (por defecto 400px)
 * @param quality - Calidad de compresión (0-1, por defecto 0.8)
 * @returns Promise con el Blob de la imagen redimensionada
 */
export async function resizeImage(
  file: File,
  maxWidth: number = 400,
  maxHeight: number = 400,
  quality: number = 0.8
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    
    reader.onload = (e) => {
      const img = new Image()
      
      img.onload = () => {
        // Calcular nuevas dimensiones manteniendo la proporción
        let width = img.width
        let height = img.height
        
        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width = width * ratio
          height = height * ratio
        }
        
        // Crear canvas para redimensionar
        const canvas = document.createElement('canvas')
        canvas.width = width
        canvas.height = height
        
        const ctx = canvas.getContext('2d')
        if (!ctx) {
          reject(new Error('No se pudo obtener el contexto del canvas'))
          return
        }
        
        // Dibujar imagen redimensionada
        ctx.drawImage(img, 0, 0, width, height)
        
        // Convertir a Blob
        canvas.toBlob(
          (blob) => {
            if (blob) {
              resolve(blob)
            } else {
              reject(new Error('Error al convertir la imagen'))
            }
          },
          'image/jpeg',
          quality
        )
      }
      
      img.onerror = () => {
        reject(new Error('Error al cargar la imagen'))
      }
      
      if (e.target?.result) {
        img.src = e.target.result as string
      }
    }
    
    reader.onerror = () => {
      reject(new Error('Error al leer el archivo'))
    }
    
    reader.readAsDataURL(file)
  })
}
