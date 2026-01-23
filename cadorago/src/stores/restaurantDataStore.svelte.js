/**
 * Store reactivo para los datos del restaurante usando Svelte 5 runes
 */
import { getSlugFromUrl, getVersionIdFromUrl, fetchRestaurantDataBySlug, fetchRestaurantDataById } from '../services/restaurantData.js';
import { adaptMenuData, DEFAULT_TEST_MENU } from '../services/menuData.js';

/**
 * Verifica si estamos en la ruta /test
 * @returns {boolean}
 */
function isTestRoute() {
  if (typeof window === 'undefined') return false;
  return window.location.pathname === '/test' || window.location.pathname.endsWith('/test');
}

class RestaurantDataStore {
  constructor() {
    this.data = $state(null);
    this.loading = $state(true);
    this.error = $state(null);
    
    // Cargar datos automáticamente al inicializar
    this.load();
    
    // Escuchar cambios en la ruta (para SPA)
    if (typeof window !== 'undefined') {
      // Usar popstate para detectar cambios de ruta
      window.addEventListener('popstate', () => {
        this.load();
      });
      
      // También escuchar cambios programáticos de ruta
      const originalPushState = history.pushState;
      const originalReplaceState = history.replaceState;
      
      history.pushState = (...args) => {
        originalPushState.apply(history, args);
        this.load();
      };
      
      history.replaceState = (...args) => {
        originalReplaceState.apply(history, args);
        this.load();
      };
    }
  }
  
  async load() {
    this.loading = true;
    this.error = null;
    
    try {
      // Verificar si estamos en la ruta /test
      if (isTestRoute()) {
        // Usar el menú por defecto para /test
        this.data = adaptMenuData(DEFAULT_TEST_MENU);
      } else {
        // Intentar obtener el slug o menuId desde la URL (ruta /m/{slug} o /m/{menuId})
        const urlData = getSlugFromUrl();
        
        if (!urlData || !urlData.value) {
          throw new Error('Se requiere un slug o menuId en la URL (ej: /m/mi-restaurante o /m/019be861-4f12-767f-a371-075d291277a8)');
        }
        
        // Obtener version_id opcional desde la URL (query param)
        // Si está presente, se usa esa versión específica (para interacciones)
        // Si no está presente, se usa la versión actual (visualización simple)
        const versionId = getVersionIdFromUrl();
        
        let rawData;
        
        // Si es un menuId (UUID), usar el endpoint por menuId
        if (urlData.isMenuId) {
          if (versionId) {
            console.log('🔍 Obteniendo menú por menuId con version_id específico (interacción):', { menuId: urlData.value, versionId });
          } else {
            console.log('🔍 Obteniendo menú por menuId versión actual (visualización simple):', urlData.value);
          }
          rawData = await fetchRestaurantDataById(urlData.value, versionId);
        } else {
          // Si es un slug, usar el endpoint por slug
          if (versionId) {
            console.log('🔍 Obteniendo menú por slug con version_id específico (interacción):', { slug: urlData.value, versionId });
          } else {
            console.log('🔍 Obteniendo menú por slug versión actual (visualización simple):', urlData.value);
          }
          rawData = await fetchRestaurantDataBySlug(urlData.value, versionId);
        }
        
        this.data = adaptMenuData(rawData);
      }
    } catch (err) {
      this.error = err.message || 'Error al cargar los datos del restaurante';
      console.error('Error en RestaurantDataStore:', err);
    } finally {
      this.loading = false;
    }
  }
  
  get value() {
    return this.data;
  }
}

// Crear instancia única del store
export const restaurantDataStore = new RestaurantDataStore();

