/**
 * Store del carrito de compras usando Svelte 5 runes
 */
import { getPriceFromPricing } from '../services/menuData.js';

const STORAGE_KEY = 'cadorago_cart';

class CartStore {
  constructor() {
    // Cargar items desde localStorage al inicializar
    const savedItems = this.loadFromStorage();
    this.items = $state(savedItems);
  }

  /**
   * Carga los items del carrito desde localStorage
   * @returns {Array} Array de items del carrito
   */
  loadFromStorage() {
    try {
      if (typeof window === 'undefined' || !window.localStorage) {
        return [];
      }
      
      const saved = localStorage.getItem(STORAGE_KEY);
      if (!saved) {
        return [];
      }
      
      const parsed = JSON.parse(saved);
      // Validar que sea un array
      return Array.isArray(parsed) ? parsed : [];
    } catch (error) {
      console.warn('Error al cargar el carrito desde localStorage:', error);
      return [];
    }
  }

  /**
   * Guarda los items del carrito en localStorage
   */
  saveToStorage() {
    try {
      if (typeof window === 'undefined' || !window.localStorage) {
        return;
      }
      
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.items));
    } catch (error) {
      console.warn('Error al guardar el carrito en localStorage:', error);
    }
  }

  /**
   * Agrega un item al carrito
   * @param {Object} item - Item del menú a agregar
   * @param {Object} side - Side seleccionado (opcional)
   */
  addItem(item, side = null) {
    // Si tiene sides, debe tener uno seleccionado
    if (item.sides && item.sides.length > 0 && !side) {
      throw new Error('Debe seleccionar un acompañamiento');
    }
    
    // Determinar precio: usar precio del side si existe, sino el precio del item
    // El precio puede venir directamente (formato antiguo) o desde pricing (formato nuevo)
    let precio = 0;
    if (side) {
      precio = side.price || (side.pricing?.pricePerUnit || 0);
    } else {
      precio = item.price || (item.pricing?.pricePerUnit || 0);
    }
    
    // Crear clave única: title + name del side (si existe)
    const itemKey = side 
      ? `${item.title}_${side.name}` 
      : item.title;
    
    const existingItemIndex = this.items.findIndex(i => {
      const existingKey = i.acompanamientoId 
        ? `${i.title}_${i.acompanamientoId}` 
        : i.title;
      return existingKey === itemKey;
    });
    
    if (existingItemIndex !== -1) {
      // Crear nuevo array para forzar reactividad
      this.items = this.items.map((i, index) => {
        if (index === existingItemIndex) {
          return { ...i, cantidad: i.cantidad + 1 };
        }
        return i;
      });
    } else {
      this.items = [...this.items, {
        ...item,
        cantidad: 1,
        precio: precio,
        acompanamiento: side ? side.name : null,
        acompanamientoId: side ? side.name : null
      }];
    }
    
    // Guardar en localStorage después de modificar
    this.saveToStorage();
  }

  /**
   * Agrega un item al carrito con cantidad personalizada (para WEIGHT, VOLUME, etc.)
   * @param {Object} item - Item del menú a agregar
   * @param {number} quantity - Cantidad personalizada
   */
  addItemWithQuantity(item, quantity) {
    if (!item.pricing) {
      throw new Error('El item debe tener pricing para usar cantidad personalizada');
    }
    
    if (quantity <= 0) {
      throw new Error('La cantidad debe ser mayor a 0');
    }
    
    // Calcular precio según la cantidad y el modo de pricing
    const precio = getPriceFromPricing(item.pricing, quantity);
    
    // Crear clave única: title + cantidad (para permitir múltiples cantidades diferentes)
    const itemKey = `${item.title}_${quantity}`;
    
    // Verificar si ya existe un item con la misma cantidad
    const existingItemIndex = this.items.findIndex(i => {
      const existingKey = i.customQuantity 
        ? `${i.title}_${i.customQuantity}` 
        : i.title;
      return existingKey === itemKey;
    });
    
    if (existingItemIndex !== -1) {
      // Si existe, incrementar cantidad (pero esto es raro para items con cantidad personalizada)
      // Mejor agregar como nuevo item
      this.items = [...this.items, {
        ...item,
        cantidad: 1,
        customQuantity: quantity,
        precio: precio,
        pricing: item.pricing // Mantener el pricing original para recálculos
      }];
    } else {
      this.items = [...this.items, {
        ...item,
        cantidad: 1,
        customQuantity: quantity,
        precio: precio,
        pricing: item.pricing // Mantener el pricing original para recálculos
      }];
    }
    
    // Guardar en localStorage después de modificar
    this.saveToStorage();
  }

  /**
   * Elimina un item del carrito
   * @param {string} title - Título del item a eliminar
   */
  removeItem(title) {
    this.items = this.items.filter(item => item.title !== title);
    this.saveToStorage();
  }

  /**
   * Actualiza la cantidad de un item
   * @param {string} itemKey - Clave única del item (title o title_acompanamientoId)
   * @param {number} cantidad - Nueva cantidad
   */
  updateQuantity(itemKey, cantidad) {
    if (cantidad <= 0) {
      this.removeItemByKey(itemKey);
      return;
    }
    
    // Crear nuevo array para forzar reactividad
    this.items = this.items.map(item => {
      const currentKey = item.acompanamientoId 
        ? `${item.title}_${item.acompanamientoId}` 
        : item.title;
      if (currentKey === itemKey) {
        return { ...item, cantidad };
      }
      return item;
    });
    
    this.saveToStorage();
  }
  
  /**
   * Elimina un item del carrito por clave única
   * @param {string} itemKey - Clave única del item
   */
  removeItemByKey(itemKey) {
    this.items = this.items.filter(item => {
      const currentKey = item.acompanamientoId 
        ? `${item.title}_${item.acompanamientoId}` 
        : item.title;
      return currentKey !== itemKey;
    });
    this.saveToStorage();
  }

  /**
   * Calcula el total del carrito
   * @returns {number} Total del carrito
   */
  getTotal() {
    return this.items.reduce((total, item) => {
      // Si tiene cantidad personalizada y pricing, calcular precio dinámicamente
      if (item.customQuantity && item.pricing) {
        return total + getPriceFromPricing(item.pricing, item.customQuantity);
      }
      // Precio normal
      return total + (item.precio * item.cantidad);
    }, 0);
  }

  /**
   * Obtiene la cantidad total de items
   * @returns {number} Cantidad total de items
   */
  getTotalItems() {
    return this.items.reduce((total, item) => {
      // Para items con cantidad personalizada, sumar la cantidad personalizada
      if (item.customQuantity) {
        return total + item.customQuantity;
      }
      return total + item.cantidad;
    }, 0);
  }

  /**
   * Limpia el carrito
   */
  clear() {
    // Asignar nuevo array vacío - esto debería activar la reactividad en Svelte 5
    this.items = [];
    // Guardar en localStorage inmediatamente
    this.saveToStorage();
  }

  /**
   * Genera el mensaje para WhatsApp con el pedido
   * @param {string} whatsappNumber - Número de WhatsApp
   * @param {string} nombreRetiro - Nombre de quien va a retirar
   * @param {string} horaRetiro - Hora de retiro
   * @param {string} lang - Idioma ('ES', 'PT', 'EN')
   * @param {object} translations - Objeto con las traducciones de WhatsApp
   * @returns {string} URL de WhatsApp con el mensaje
   */
  generateWhatsAppMessage(whatsappNumber, nombreRetiro = '', horaRetiro = '', lang = 'ES', translations = null) {
    // Si no se pasan traducciones, usar valores por defecto en español
    const t = translations || {
      greeting: "¡Hola! Me gustaría hacer el siguiente pedido:\n\n",
      each: "c/u",
      itemTotal: "Total",
      orderTotal: "Total",
      pickupInfoLabel: "Información de retiro:\n",
      pickupNameLabel: "👤 Nombre:",
      pickupTimeLabel: "🕐 Hora de retiro:"
    };
    
    // Determinar locale para formateo de números
    const locale = lang === 'PT' ? 'pt-BR' : lang === 'EN' ? 'en-US' : 'es-CL';
    
    let message = t.greeting;
    
    this.items.forEach((item, index) => {
      message += `${index + 1}. ${item.title}`;
      if (item.acompanamiento) {
        message += ` (${item.acompanamiento})`;
      }
      
      // Manejar cantidad personalizada (WEIGHT, VOLUME, etc.)
      if (item.customQuantity) {
        const unitLabel = item.pricing?.unit === 'GRAM' ? 'g' : 
                         item.pricing?.unit === 'KILOGRAM' ? 'kg' :
                         item.pricing?.unit === 'MILLILITER' ? 'ml' :
                         item.pricing?.unit === 'LITER' ? 'L' :
                         item.pricing?.unit === 'METER' ? 'm' :
                         item.pricing?.unit === 'SQUARE_METER' ? 'm²' : '';
        message += ` (${item.customQuantity}${unitLabel})`;
        
        // Calcular precio para cantidad personalizada
        const itemPrice = getPriceFromPricing(item.pricing, item.customQuantity);
        message += ` - $${itemPrice.toLocaleString(locale)}`;
      } else {
        if (item.cantidad > 1) {
          message += ` x${item.cantidad}`;
        }
        message += ` - $${item.precio.toLocaleString(locale)}`;
        if (item.cantidad > 1) {
          message += ` ${t.each} (${t.itemTotal}: $${(item.precio * item.cantidad).toLocaleString(locale)})`;
        }
      }
      message += "\n";
    });
    
    message += `\n*${t.orderTotal}: $${this.getTotal().toLocaleString(locale)}*\n\n`;
    
    if (nombreRetiro || horaRetiro) {
      message += t.pickupInfoLabel;
      if (nombreRetiro) {
        message += `${t.pickupNameLabel} ${nombreRetiro}\n`;
      }
      if (horaRetiro) {
        message += `${t.pickupTimeLabel} ${horaRetiro}\n`;
      }
      message += "\n";
    }
    
    message += lang === 'EN' ? 'Thank you!' : lang === 'PT' ? 'Obrigado!' : 'Gracias!';
    
    const encodedMessage = encodeURIComponent(message);
    const phoneNumber = whatsappNumber.replace(/[^0-9]/g, '');
    return `https://wa.me/${phoneNumber}?text=${encodedMessage}`;
  }
}

// Crear instancia única del store
export const cartStore = new CartStore();

