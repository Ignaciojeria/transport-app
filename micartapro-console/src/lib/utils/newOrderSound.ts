/**
 * Sonido de aviso cuando llega un pedido nuevo.
 * Usa el archivo public/new_order.wav (debe estar en public/ para servirse como /new_order.wav).
 */

const NEW_ORDER_SOUND_URL = '/new_order.wav'

/** WAV mínimo silencioso para desbloquear el audio sin que suene la notificación. */
const SILENT_WAV_DATA_URL =
  'data:audio/wav;base64,UklGRiQAAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQAAAAA='

/**
 * Desbloquea el audio (llamar en un gesto del usuario: clic en filtro, etc.).
 * Reproduce un sonido silencioso para que el navegador permita luego el sonido de nueva orden.
 */
export function ensureAudioUnlocked(): void {
  if (typeof window === 'undefined') return
  try {
    const audio = new Audio(SILENT_WAV_DATA_URL)
    audio.volume = 0
    audio.play().catch(() => {})
  } catch {
    // Ignorar
  }
}

/**
 * Reproduce el sonido de pedido nuevo (archivo WAV desde public/new_order.wav).
 */
export function playNewOrderSound(): void {
  try {
    const audio = new Audio(NEW_ORDER_SOUND_URL)
    audio.volume = 0.7
    audio.play().catch(() => {})
  } catch {
    // Ignorar
  }
}
