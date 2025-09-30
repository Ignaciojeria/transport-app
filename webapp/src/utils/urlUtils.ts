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
    console.log('🔍 URL completa:', currentUrl)
    
    const urlObj = new URL(currentUrl)
    console.log('🔍 URL object:', {
      origin: urlObj.origin,
      pathname: urlObj.pathname,
      search: urlObj.search,
      hash: urlObj.hash
    })
    
    // Buscar el token en el fragment (después del #)
    const fragment = urlObj.hash
    console.log('🔍 Fragment completo:', fragment)
    
    if (!fragment) {
      console.warn('❌ No se encontró fragment en la URL')
      return null
    }
    
    // Parsear parámetros del fragment
    const fragmentWithoutHash = fragment.substring(1) // Remover el #
    console.log('🔍 Fragment sin #:', fragmentWithoutHash)
    
    const params = new URLSearchParams(fragmentWithoutHash)
    console.log('🔍 Parámetros parseados:', Object.fromEntries(params.entries()))
    
    // Buscar el token en diferentes parámetros posibles
    const accessToken = params.get('access_token')
    const idToken = params.get('id_token')
    const token = params.get('token')
    const jwt = params.get('jwt')
    const auth = params.get('auth') // Nuevo: parámetro 'auth' que contiene el payload completo
    
    console.log('🔍 Tokens encontrados:', {
      access_token: accessToken ? `${accessToken.substring(0, 20)}...` : null,
      id_token: idToken ? `${idToken.substring(0, 20)}...` : null,
      token: token ? `${token.substring(0, 20)}...` : null,
      jwt: jwt ? `${jwt.substring(0, 20)}...` : null,
      auth: auth ? `${auth.substring(0, 20)}...` : null
    })
    
    // Si tenemos el parámetro 'auth', extraer el access_token de ahí
    if (auth) {
      try {
        console.log('🔍 Procesando parámetro auth...')
        const authPayload = JSON.parse(atob(auth))
        console.log('🔍 Auth payload decodificado:', authPayload)
        
        if (authPayload.access_token) {
          console.log('✅ Access token encontrado en auth payload')
          return authPayload.access_token
        }
      } catch (error) {
        console.error('❌ Error al decodificar auth payload:', error)
      }
    }
    
    const finalToken = accessToken || idToken || token || jwt
    
    if (!finalToken) {
      console.warn('❌ No se encontró token en el fragment de la URL')
      console.log('🔍 Parámetros disponibles:', Array.from(params.keys()))
      return null
    }
    
    console.log('✅ Token encontrado:', finalToken.substring(0, 20) + '...')
    return finalToken
  } catch (error) {
    console.error('❌ Error al extraer token del fragment:', error)
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
    
    // Buscar email directamente en parámetros
    const directEmail = params.get('email')
    if (directEmail) {
      return directEmail
    }
    
    // Buscar email en el payload de auth
    const auth = params.get('auth')
    if (auth) {
      try {
        const authPayload = JSON.parse(atob(auth))
        console.log('🔍 Extrayendo email del auth payload:', authPayload)
        
        if (authPayload.user && authPayload.user.email) {
          console.log('✅ Email encontrado en auth.user.email:', authPayload.user.email)
          return authPayload.user.email
        }
      } catch (error) {
        console.error('❌ Error al decodificar auth payload para email:', error)
      }
    }
    
    return null
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
