/**
 * Servicio de tracking de órdenes
 * Llama al backend para obtener el estado de una orden por su código de seguimiento.
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ||
  (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ? 'http://localhost:8082'
    : 'https://micartapro-backend-27303662337.us-central1.run.app');

/**
 * Obtiene el estado de una orden por su trackingId
 * @param {string} trackingId - Código de seguimiento (ej: ABC12345)
 * @returns {Promise<Object>} Datos de la orden: { trackingId, aggregateId, orderNumber, menuId, fulfillment, items, requestedAt, createdAt }
 */
export async function getOrderByTrackingId(trackingId) {
  if (!trackingId || !trackingId.trim()) {
    throw new Error('El código de seguimiento es requerido');
  }

  const response = await fetch(
    `${API_BASE_URL}/api/tracking/${encodeURIComponent(trackingId.trim().toUpperCase())}`,
    {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    }
  );

  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('Pedido no encontrado. Verifica el código e intenta nuevamente.');
    }
    const text = await response.text();
    throw new Error(`Error al consultar el pedido: ${response.status} ${text}`);
  }

  return response.json();
}
