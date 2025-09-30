/**
 * Extracción temprana de token antes de que React se monte
 * Esto evita que otros scripts limpien el fragment antes de que lo procesemos
 */

let earlyToken: string | null = null
let earlyEmail: string | null = null

/**
 * Extrae el token inmediatamente cuando se carga la página
 * Debe llamarse ANTES de que React se monte
 */
export const extractTokenEarly = (): { token: string | null; email: string | null } => {
  if (earlyToken !== null) {
    return { token: earlyToken, email: earlyEmail }
  }

  try {
    const currentUrl = window.location.href
    console.log('🔍 Extracción temprana - URL:', currentUrl)
    
    if (!currentUrl.includes('#')) {
      console.log('❌ No hay fragment en la URL')
      return { token: null, email: null }
    }

    const urlObj = new URL(currentUrl)
    const fragment = urlObj.hash
    console.log('🔍 Extracción temprana - Fragment:', fragment)

    if (!fragment) {
      return { token: null, email: null }
    }

    const params = new URLSearchParams(fragment.substring(1))
    const auth = params.get('auth')

    if (auth) {
      console.log('🔍 Extracción temprana - Auth encontrado')
      const authPayload = JSON.parse(atob(auth))
      console.log('🔍 Extracción temprana - Auth payload:', authPayload)

      if (authPayload.access_token) {
        earlyToken = authPayload.access_token
        earlyEmail = authPayload.user?.email || null
        
        console.log('✅ Extracción temprana - Token guardado:', earlyToken ? earlyToken.substring(0, 20) + '...' : 'null')
        console.log('✅ Extracción temprana - Email guardado:', earlyEmail)
        
        return { token: earlyToken, email: earlyEmail }
      }
    }

    return { token: null, email: null }
  } catch (error) {
    console.error('❌ Error en extracción temprana:', error)
    return { token: null, email: null }
  }
}

/**
 * Obtiene el token extraído tempranamente
 */
export const getEarlyToken = (): string | null => earlyToken

/**
 * Obtiene el email extraído tempranamente
 */
export const getEarlyEmail = (): string | null => earlyEmail
