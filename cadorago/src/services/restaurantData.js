/**
 * Servicio de datos del restaurante Cadorago
 * 
 * NOTA: Este servicio ahora usa exclusivamente la API del backend.
 * El método legacy que usaba Google Cloud Storage ha sido eliminado.
 */

// URL del backend de micartapro
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 
  (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ? 'http://localhost:8082'
    : 'https://micartapro-backend-27303662337.us-central1.run.app');

/**
 * Verifica si una cadena es un UUID válido
 * @param {string} str - Cadena a verificar
 * @returns {boolean} true si es un UUID válido
 */
function isUUID(str) {
  if (!str) return false;
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return uuidRegex.test(str);
}

/**
 * Obtiene el slug o menuId desde la URL (ruta /m/{slug} o /m/{menuId})
 * @returns {{value: string | null, isMenuId: boolean}} Valor de la URL y si es un menuId
 */
export function getSlugFromUrl() {
  if (typeof window === 'undefined') {
    return null;
  }
  
  const pathname = window.location.pathname;
  // Buscar patrón /m/{slug} o /m/{menuId}
  const match = pathname.match(/^\/m\/([^\/]+)/);
  if (match && match[1]) {
    const value = match[1];
    return {
      value: value,
      isMenuId: isUUID(value)
    };
  }
  
  return null;
}

/**
 * Obtiene el version_id desde la URL (query params)
 * @returns {string | null} version_id si está presente en la URL
 */
export function getVersionIdFromUrl() {
  if (typeof window === 'undefined') {
    return null;
  }
  
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get('version_id');
}

/**
 * Obtiene los parámetros userID y menuID desde la URL (query params)
 * @deprecated Este método es legacy y ya no se usa. Se mantiene solo por compatibilidad.
 * @returns {{userID: string | null, menuID: string | null}}
 */
export function getUrlParams() {
  if (typeof window === 'undefined') {
    return { userID: null, menuID: null };
  }
  
  const urlParams = new URLSearchParams(window.location.search);
  return {
    userID: urlParams.get('userID'),
    menuID: urlParams.get('menuID')
  };
}

/**
 * Obtiene los datos del restaurante desde el backend usando el menuId
 * @param {string} menuId - ID del menú (UUID)
 * @param {string} [versionId] - ID de la versión opcional. Si se proporciona, se usa esa versión específica.
 *                                Si no se proporciona, se usa la versión actual (current_version_id)
 * @returns {Promise<Object>} Datos del restaurante
 */
export async function fetchRestaurantDataById(menuId, versionId = null) {
  try {
    // Construir URL con version_id opcional como query parameter
    let apiUrl = `${API_BASE_URL}/menu/${encodeURIComponent(menuId)}`;
    
    if (versionId) {
      apiUrl += `?version_id=${encodeURIComponent(versionId)}`;
      console.log('📦 Obteniendo menú desde backend por menuId con version_id específico:', apiUrl);
    } else {
      console.log('📦 Obteniendo menú desde backend por menuId (versión actual):', apiUrl);
    }
    
    const response = await fetch(apiUrl, {
      cache: 'no-store', // Evitar cache del navegador
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error(`Menú no encontrado para el menuId: ${menuId}${versionId ? ` con version_id: ${versionId}` : ''}`);
      }
      throw new Error(`Error al obtener datos del restaurante: ${response.status} ${response.statusText}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error al cargar datos del restaurante desde backend:', error);
    throw error;
  }
}

/**
 * Obtiene los datos del restaurante desde el backend usando el slug
 * @param {string} slug - Slug del menú
 * @param {string} [versionId] - ID de la versión opcional. Si se proporciona, se usa esa versión específica.
 *                                Si no se proporciona, se usa la versión actual (current_version_id)
 * @returns {Promise<Object>} Datos del restaurante
 */
export async function fetchRestaurantDataBySlug(slug, versionId = null) {
  try {
    // Construir URL con version_id opcional como query parameter
    let apiUrl = `${API_BASE_URL}/menu/slug/${encodeURIComponent(slug)}`;
    
    if (versionId) {
      apiUrl += `?version_id=${encodeURIComponent(versionId)}`;
      console.log('📦 Obteniendo menú desde backend con version_id específico:', apiUrl);
    } else {
      console.log('📦 Obteniendo menú desde backend (versión actual):', apiUrl);
    }
    
    const response = await fetch(apiUrl, {
      cache: 'no-store', // Evitar cache del navegador
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error(`Menú no encontrado para el slug: ${slug}${versionId ? ` con version_id: ${versionId}` : ''}`);
      }
      throw new Error(`Error al obtener datos del restaurante: ${response.status} ${response.statusText}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error al cargar datos del restaurante desde backend:', error);
    throw error;
  }
}
