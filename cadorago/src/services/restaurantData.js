/**
 * Servicio de datos del restaurante Cadorago
 */
const BASE_URL = "https://storage.googleapis.com/micartapro-menus";

// URL del backend de micartapro
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 
  (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ? 'http://localhost:8082'
    : 'https://micartapro-backend-27303662337.us-central1.run.app');


/**
 * Obtiene el nombre del archivo desde latest.json
 * @param {string} userID - ID del usuario
 * @param {string} menuID - ID del menú
 * @returns {Promise<string>} Nombre del archivo JSON
 */
async function getLatestFilename(userID, menuID) {
  // Agregar timestamp para evitar cache y asegurar obtener la versión más reciente
  const timestamp = Date.now();
  const latestUrl = `${BASE_URL}/${userID}/menus/${menuID}/latest.json?t=${timestamp}`;
  const response = await fetch(latestUrl, {
    cache: 'no-store' // Evitar cache del navegador
  });
  
  if (!response.ok) {
    throw new Error(`Error al obtener latest.json: ${response.status} ${response.statusText}`);
  }
  
  const data = await response.json();
  return data.filename;
}

/**
 * Obtiene los datos del restaurante desde Google Cloud Storage
 * @param {string} userID - ID del usuario
 * @param {string} menuID - ID del menú
 * @returns {Promise<Object>} Datos del restaurante
 */
export async function fetchRestaurantData(userID, menuID) {
  try {
    // Primero obtenemos el nombre del archivo desde latest.json
    const filename = await getLatestFilename(userID, menuID);
    
    // Luego obtenemos los datos del restaurante desde el archivo referenciado
    // Agregar timestamp para evitar cache y asegurar obtener la versión más reciente
    const timestamp = Date.now();
    const dataUrl = `${BASE_URL}/${userID}/menus/${menuID}/${filename}?t=${timestamp}`;
    
    // Mostrar la URL en la consola
    console.log('📦 URL del storage:', dataUrl);
    
    const response = await fetch(dataUrl, {
      cache: 'no-store' // Evitar cache del navegador
    });
    
    if (!response.ok) {
      throw new Error(`Error al obtener datos del restaurante: ${response.status} ${response.statusText}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error al cargar datos del restaurante:', error);
    throw error;
  }
}

/**
 * Obtiene el slug desde la URL (ruta /m/{slug})
 * @returns {string | null} Slug del menú
 */
export function getSlugFromUrl() {
  if (typeof window === 'undefined') {
    return null;
  }
  
  const pathname = window.location.pathname;
  // Buscar patrón /m/{slug}
  const match = pathname.match(/^\/m\/([^\/]+)/);
  if (match && match[1]) {
    return match[1];
  }
  
  return null;
}

/**
 * Obtiene los parámetros userID y menuID desde la URL (query params)
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
 * Obtiene los datos del restaurante desde el backend usando el slug
 * @param {string} slug - Slug del menú
 * @returns {Promise<Object>} Datos del restaurante
 */
export async function fetchRestaurantDataBySlug(slug) {
  try {
    const apiUrl = `${API_BASE_URL}/menu/slug/${encodeURIComponent(slug)}`;
    
    console.log('📦 Obteniendo menú desde backend:', apiUrl);
    
    const response = await fetch(apiUrl, {
      cache: 'no-store', // Evitar cache del navegador
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error(`Menú no encontrado para el slug: ${slug}`);
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
