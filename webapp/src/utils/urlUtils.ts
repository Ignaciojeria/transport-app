/**
 * Utilidades para manejar URLs y extraer parámetros
 */

/**
 * Extrae el token JWT del fragment de la URL después de la redirección de Google OAuth
 * @param url - URL completa o fragment de la URL
 * @returns token JWT o null si no se encuentra
 */
export const extractTokenFromFragment = (url?: string): string | null => {
  try {
    const currentUrl = url || window.location.href
    const urlObj = new URL(currentUrl)
    
    // Buscar el token en el fragment (después del #)
    const fragment = urlObj.hash
    
    if (!fragment) {
      console.warn('No se encontró fragment en la URL')
      return null
    }
    
    // Parsear parámetros del fragment
    const params = new URLSearchParams(fragment.substring(1)) // Remover el #
    
    // Buscar el token en diferentes parámetros posibles
    const token = params.get('access_token') || 
                  params.get('id_token') || 
                  params.get('token') ||
                  params.get('jwt')
    
    if (!token) {
      console.warn('No se encontró token en el fragment de la URL')
      return null
    }
    
    return token
  } catch (error) {
    console.error('Error al extraer token del fragment:', error)
    return null
  }
}

/**
 * Extrae el email del fragment de la URL (como fallback)
 * @param url - URL completa o fragment de la URL
 * @returns email o null si no se encuentra
 */
export const extractEmailFromFragment = (url?: string): string | null => {
  try {
    const currentUrl = url || window.location.href
    const urlObj = new URL(currentUrl)
    const fragment = urlObj.hash
    
    if (!fragment) {
      return null
    }
    
    const params = new URLSearchParams(fragment.substring(1))
    return params.get('email')
  } catch (error) {
    console.error('Error al extraer email del fragment:', error)
    return null
  }
}

/**
 * Limpia la URL removiendo el fragment después de extraer el token
 * @param url - URL a limpiar
 * @returns URL limpia sin fragment
 */
export const cleanUrlAfterTokenExtraction = (url?: string): string => {
  try {
    const currentUrl = url || window.location.href
    const urlObj = new URL(currentUrl)
    urlObj.hash = ''
    return urlObj.toString()
  } catch (error) {
    console.error('Error al limpiar URL:', error)
    return url || window.location.href
  }
}
