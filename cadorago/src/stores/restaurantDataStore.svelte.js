/**
 * Store reactivo para los datos del restaurante usando Svelte 5 runes
 */
import { getUrlParams, fetchRestaurantData } from '../services/restaurantData.js';

class RestaurantDataStore {
  constructor() {
    this.data = $state(null);
    this.loading = $state(true);
    this.error = $state(null);
    
    // Cargar datos automáticamente al inicializar
    this.load();
  }
  
  async load() {
    this.loading = true;
    this.error = null;
    
    try {
      const { userID, menuID } = getUrlParams();
      
      if (!userID || !menuID) {
        throw new Error('Los parámetros userID y menuID son requeridos en la URL (ej: ?userID=xxx&menuID=yyy)');
      }
      
      this.data = await fetchRestaurantData(userID, menuID);
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

