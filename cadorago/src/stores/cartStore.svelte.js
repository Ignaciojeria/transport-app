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
   */
  addItem(item) {
    const existingItemIndex = this.items.findIndex(i => i.titulo === item.titulo);
    
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
        cantidad: 1
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
   * @param {string} titulo - Título del item
   * @param {number} cantidad - Nueva cantidad
   */
  updateQuantity(titulo, cantidad) {
    if (cantidad <= 0) {
      this.removeItem(titulo);
      return;
    }
    
    // Crear nuevo array para forzar reactividad
    this.items = this.items.map(item => {
      if (item.titulo === titulo) {
        return { ...item, cantidad };
      }
      return item;
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
   * @returns {string} URL de WhatsApp con el mensaje
   */
  generateWhatsAppMessage(whatsappNumber) {
    let message = "¡Hola! Me gustaría hacer el siguiente pedido:\n\n";
    
    this.items.forEach((item, index) => {
      message += `${index + 1}. ${item.titulo}`;
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
    message += "Gracias!";
    
    const encodedMessage = encodeURIComponent(message);
    const phoneNumber = whatsappNumber.replace(/[^0-9]/g, '');
    return `https://wa.me/${phoneNumber}?text=${encodedMessage}`;
  }
}

// Crear instancia única del store
export const cartStore = new CartStore();

