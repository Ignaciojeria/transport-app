/**
 * Extrae el email de un JWT token
 * @param token - JWT token en formato string
 * @returns email extraído del token o null si no se puede extraer
 */
export const extractEmailFromJWT = (token: string): string | null => {
  try {
    console.log('🔍 Extrayendo email del JWT:', token.substring(0, 20) + '...')
    
    // Decodificar el JWT (solo la parte del payload)
    const parts = token.split('.')
    console.log('🔍 Partes del JWT:', parts.length)
    
    if (parts.length !== 3) {
      console.error('❌ Token JWT inválido: no tiene 3 partes')
      return null
    }

    // Decodificar el payload (parte del medio)
    const payload = parts[1]
    console.log('🔍 Payload (base64):', payload)
    
    // Agregar padding si es necesario para base64
    const paddedPayload = payload + '='.repeat((4 - payload.length % 4) % 4)
    console.log('🔍 Payload con padding:', paddedPayload)
    
    // Decodificar de base64
    const decodedPayload = atob(paddedPayload)
    console.log('🔍 Payload decodificado:', decodedPayload)
    
    // Parsear el JSON
    const payloadObj = JSON.parse(decodedPayload)
    console.log('🔍 Payload parseado:', payloadObj)
    
    // Extraer el email
    const email = payloadObj.email || payloadObj.sub || payloadObj.preferred_username
    console.log('🔍 Email extraído:', email)
    
    if (!email) {
      console.error('❌ No se encontró email en el token JWT')
      console.log('🔍 Campos disponibles:', Object.keys(payloadObj))
      return null
    }
    
    console.log('✅ Email encontrado:', email)
    return email
  } catch (error) {
    console.error('❌ Error al extraer email del JWT:', error)
    return null
  }
}

/**
 * Valida si un token JWT está expirado
 * @param token - JWT token en formato string
 * @returns true si está expirado, false si no
 */
export const isJWTExpired = (token: string): boolean => {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return true

    const payload = parts[1]
    const paddedPayload = payload + '='.repeat((4 - payload.length % 4) % 4)
    const decodedPayload = atob(paddedPayload)
    const payloadObj = JSON.parse(decodedPayload)
    
    const exp = payloadObj.exp
    if (!exp) return true
    
    const currentTime = Math.floor(Date.now() / 1000)
    return currentTime >= exp
  } catch (error) {
    console.error('Error al validar expiración del JWT:', error)
    return true
  }
}
