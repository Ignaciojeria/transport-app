/**
 * Store reactivo para los datos del restaurante usando Svelte 5 runes
 */
import { getSlugFromUrl, getVersionIdFromUrl, fetchRestaurantDataBySlug } from '../services/restaurantData.js';
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
        // Intentar obtener el slug desde la URL (ruta /m/{slug})
        const slug = getSlugFromUrl();
        
        if (!slug) {
          throw new Error('Se requiere un slug en la URL (ej: /m/mi-restaurante)');
        }
        
        // Obtener version_id opcional desde la URL (query param)
        // Si está presente, se usa esa versión específica (para interacciones)
        // Si no está presente, se usa la versión actual (visualización simple)
        const versionId = getVersionIdFromUrl();
        
        if (versionId) {
          console.log('🔍 Obteniendo menú con version_id específico (interacción):', { slug, versionId });
        } else {
          console.log('🔍 Obteniendo menú versión actual (visualización simple):', slug);
        }
        
        // Usar siempre el endpoint del backend con el slug
        const rawData = await fetchRestaurantDataBySlug(slug, versionId);
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

