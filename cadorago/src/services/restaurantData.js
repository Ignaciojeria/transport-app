/**
 * Servicio de datos del restaurante Cadorago
 */
const BASE_URL = "https://storage.googleapis.com/micartapro-menus";


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
 * Obtiene los parámetros userID y menuID desde la URL
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

