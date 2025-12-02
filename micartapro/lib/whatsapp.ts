/**
 * Genera un enlace de WhatsApp con un mensaje predefinido
 * @param phoneNumber - Número de teléfono (con código de país)
 * @param message - Mensaje a enviar
 * @returns URL de WhatsApp
 */
export function generateWhatsAppLink(phoneNumber: string, message: string): string {
  const cleanPhone = phoneNumber.replace(/[^0-9]/g, '')
  const encodedMessage = encodeURIComponent(message)
  return `https://wa.me/${cleanPhone}?text=${encodedMessage}`
}

/**
 * Abre WhatsApp con un mensaje de cotización
 */
export function openWhatsAppQuote(): void {
  const phoneNumber = '+56957857558'
  const message = 'Hola! Me gustaría cotizar MiCartaPro para mi restaurante.'
  const url = generateWhatsAppLink(phoneNumber, message)
  window.open(url, '_blank')
}

