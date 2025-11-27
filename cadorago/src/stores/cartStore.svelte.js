/**
 * Store del carrito de compras usando Svelte 5 runes
 */

class CartStore {
  constructor() {
    this.items = $state([]);
  }

  /**
   * Agrega un item al carrito
   * @param {Object} item - Item del menú a agregar
   * @param {Object} acompanamiento - Acompañamiento seleccionado (opcional)
   */
  addItem(item, acompanamiento = null) {
    // Si tiene acompañamientos, debe tener uno seleccionado
    if (item.acompanamientos && item.acompanamientos.length > 0 && !acompanamiento) {
      throw new Error('Debe seleccionar un acompañamiento');
    }
    
    // Determinar precio: usar precio del acompañamiento si existe, sino el precio del item
    const precio = acompanamiento ? acompanamiento.precio : (item.precio || 0);
    
    // Crear clave única: titulo + id del acompañamiento (si existe)
    const itemKey = acompanamiento 
      ? `${item.titulo}_${acompanamiento.id}` 
      : item.titulo;
    
    const existingItemIndex = this.items.findIndex(i => {
      const existingKey = i.acompanamientoId 
        ? `${i.titulo}_${i.acompanamientoId}` 
        : i.titulo;
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
        acompanamiento: acompanamiento ? acompanamiento.nombre : null,
        acompanamientoId: acompanamiento ? acompanamiento.id : null
      }];
    }
  }

  /**
   * Elimina un item del carrito
   * @param {string} titulo - Título del item a eliminar
   */
  removeItem(titulo) {
    this.items = this.items.filter(item => item.titulo !== titulo);
  }

  /**
   * Actualiza la cantidad de un item
   * @param {string} itemKey - Clave única del item (titulo o titulo_acompanamientoId)
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
        ? `${item.titulo}_${item.acompanamientoId}` 
        : item.titulo;
      if (currentKey === itemKey) {
        return { ...item, cantidad };
      }
      return item;
    });
  }
  
  /**
   * Elimina un item del carrito por clave única
   * @param {string} itemKey - Clave única del item
   */
  removeItemByKey(itemKey) {
    this.items = this.items.filter(item => {
      const currentKey = item.acompanamientoId 
        ? `${item.titulo}_${item.acompanamientoId}` 
        : item.titulo;
      return currentKey !== itemKey;
    });
  }

  /**
   * Calcula el total del carrito
   * @returns {number} Total del carrito
   */
  getTotal() {
    return this.items.reduce((total, item) => {
      return total + (item.precio * item.cantidad);
    }, 0);
  }

  /**
   * Obtiene la cantidad total de items
   * @returns {number} Cantidad total de items
   */
  getTotalItems() {
    return this.items.reduce((total, item) => total + item.cantidad, 0);
  }

  /**
   * Limpia el carrito
   */
  clear() {
    this.items = [];
  }

  /**
   * Genera el mensaje para WhatsApp con el pedido
   * @param {string} whatsappNumber - Número de WhatsApp
   * @param {string} nombreRetiro - Nombre de quien va a retirar
   * @param {string} horaRetiro - Hora de retiro
   * @returns {string} URL de WhatsApp con el mensaje
   */
  generateWhatsAppMessage(whatsappNumber, nombreRetiro = '', horaRetiro = '') {
    let message = "¡Hola! Me gustaría hacer el siguiente pedido:\n\n";
    
    this.items.forEach((item, index) => {
      message += `${index + 1}. ${item.titulo}`;
      if (item.acompanamiento) {
        message += ` (${item.acompanamiento})`;
      }
      if (item.cantidad > 1) {
        message += ` x${item.cantidad}`;
      }
      message += ` - $${item.precio.toLocaleString('es-CL')}`;
      if (item.cantidad > 1) {
        message += ` c/u (Total: $${(item.precio * item.cantidad).toLocaleString('es-CL')})`;
      }
      message += "\n";
    });
    
    message += `\n*Total: $${this.getTotal().toLocaleString('es-CL')}*\n\n`;
    
    if (nombreRetiro || horaRetiro) {
      message += "Información de retiro:\n";
      if (nombreRetiro) {
        message += `👤 Nombre: ${nombreRetiro}\n`;
      }
      if (horaRetiro) {
        message += `🕐 Hora de retiro: ${horaRetiro}\n`;
      }
      message += "\n";
    }
    
    message += "Gracias!";
    
    const encodedMessage = encodeURIComponent(message);
    const phoneNumber = whatsappNumber.replace(/[^0-9]/g, '');
    return `https://wa.me/${phoneNumber}?text=${encodedMessage}`;
  }
}

// Crear instancia única del store
export const cartStore = new CartStore();

