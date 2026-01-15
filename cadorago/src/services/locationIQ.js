/**
 * Servicio para integración con LocationIQ API
 */
const LOCATIONIQ_API_KEY = 'pk.6fb271168ebb42320bc9248737e83834';
const LOCATIONIQ_BASE_URL = 'https://api.locationiq.com/v1';
const LOCATIONIQ_STATIC_MAPS_URL = 'https://maps.locationiq.com/v3/staticmap';

/**
 * Geocodifica una dirección (convierte texto a coordenadas)
 * @param {string} address - Dirección a geocodificar
 * @returns {Promise<{lat: number, lon: number, display_name: string}>}
 */
export async function geocodeAddress(address) {
  if (!address || !address.trim()) {
    throw new Error('La dirección no puede estar vacía');
  }

  const url = `${LOCATIONIQ_BASE_URL}/search.php?key=${LOCATIONIQ_API_KEY}&q=${encodeURIComponent(address)}&format=json&limit=1`;
  
  try {
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`Error en geocodificación: ${response.status} ${response.statusText}`);
    }
    
    const data = await response.json();
    
    if (!data || data.length === 0) {
      throw new Error('No se encontró la dirección');
    }
    
    const result = data[0];
    return {
      lat: parseFloat(result.lat),
      lon: parseFloat(result.lon),
      display_name: result.display_name || address
    };
  } catch (error) {
    console.error('Error en geocodificación:', error);
    throw error;
  }
}

/**
 * Genera la URL de un mapa estático con un pin
 * @param {number} lat - Latitud
 * @param {number} lon - Longitud
 * @param {number} width - Ancho del mapa (default: 200)
 * @param {number} height - Alto del mapa (default: 150)
 * @returns {string} URL del mapa estático
 */
export function getStaticMapUrl(lat, lon, width = 200, height = 150, zoom = 18) {
  if (!lat || !lon || isNaN(lat) || isNaN(lon)) return '';
  
  // Formato: https://maps.locationiq.com/v3/staticmap?key=YOUR_API_KEY&center=lat,lon&zoom=18&size=widthxheight&markers=icon:small-red-cutout|lat,lon
  // Zoom 18 es más cercano que 16, ideal para ver direcciones específicas
  const url = `${LOCATIONIQ_STATIC_MAPS_URL}?key=${LOCATIONIQ_API_KEY}&center=${lat},${lon}&zoom=${zoom}&size=${width}x${height}&markers=icon:small-red-cutout|${lat},${lon}`;
  console.log('Generando URL del mapa:', url);
  return url;
}

/**
 * Autocompleta direcciones mientras el usuario escribe
 * @param {string} query - Texto de búsqueda
 * @returns {Promise<Array<{display_name: string, lat: number, lon: number, mapUrl: string}>>}
 */
export async function autocompleteAddress(query) {
  if (!query || query.length < 3) {
    return [];
  }

  const url = `${LOCATIONIQ_BASE_URL}/autocomplete.php?key=${LOCATIONIQ_API_KEY}&q=${encodeURIComponent(query)}&format=json&limit=5`;
  
  try {
    const response = await fetch(url);
    
    if (!response.ok) {
      return [];
    }
    
    const data = await response.json();
    
    if (!data || !Array.isArray(data)) {
      return [];
    }
    
    return data.map(item => {
      const lat = parseFloat(item.lat) || 0;
      const lon = parseFloat(item.lon) || 0;
      return {
        display_name: item.display_name || '',
        lat: lat,
        lon: lon,
        mapUrl: getStaticMapUrl(lat, lon, 200, 120)
      };
    });
  } catch (error) {
    console.error('Error en autocompletado:', error);
    return [];
  }
}
