/**
 * Store del carrito de compras
 * Usa writable store para garantizar reactividad en todos los componentes
 */
import { writable } from 'svelte/store';
import { getPriceFromPricing } from '../services/menuData.js';
import { getBaseText } from '../lib/multilingual';

const STORAGE_KEY = 'cadorago_cart';

function loadFromStorage() {
  try {
    if (typeof window === 'undefined' || !window.localStorage) return [];
    const saved = localStorage.getItem(STORAGE_KEY);
    if (!saved) return [];
    const parsed = JSON.parse(saved);
    return Array.isArray(parsed) ? parsed : [];
  } catch (error) {
    console.warn('Error al cargar el carrito:', error);
    return [];
  }
}

function saveToStorage(itemsToSave) {
  try {
    if (typeof window === 'undefined' || !window.localStorage) return;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(itemsToSave));
  } catch (error) {
    console.warn('Error al guardar el carrito:', error);
  }
}

// writable store: reactividad garantizada en Svelte
export const itemsStore = writable(loadFromStorage());

class CartStore {
  get items() {
    let val;
    itemsStore.subscribe((v) => { val = v; })();
    return val ?? [];
  }

  _setItems(value) {
    itemsStore.set(value);
  }

  /**
   * Agrega un item al carrito
   * @param {Object} item - Item del menú a agregar
   * @param {Object} side - Side seleccionado (opcional)
   * @param {Array<{descriptionId: string, optionId: string}>} descriptionSelections - Selecciones de description.selectables (opcional)
   * @param {string} menuId - ID del menú (asocia el item al menú actual, como el tracking)
   */
  addItem(item, side = null, descriptionSelections = [], menuId = null) {
    // Si tiene sides, debe tener uno seleccionado
    if (item.sides && item.sides.length > 0 && !side) {
      throw new Error('Debe seleccionar un acompañamiento');
    }
    
    // Determinar precio: usar precio del side si existe, sino el precio del item
    let precio = 0;
    if (side) {
      precio = side.price || (side.pricing?.pricePerUnit || 0);
    } else {
      precio = item.price || (item.pricing?.pricePerUnit || 0);
    }
    
    // Crear clave única: title + side + descriptionSelections
    const itemTitleBase = getBaseText(item.title);
    const sideNameBase = side ? getBaseText(side.name) : null;
    const descKey = descriptionSelections.length
      ? '_' + descriptionSelections.map(s => `${s.descriptionId}:${s.optionId}`).sort().join(',')
      : '';
    const baseKey = (side ? `${itemTitleBase}_${sideNameBase}` : itemTitleBase) + descKey;
    
    itemsStore.update((current) => {
      const existingItemIndex = current.findIndex(i => {
        if ((i.menuId ?? null) !== (menuId ?? null)) return false;
        const iTitleBase = getBaseText(i.title);
        const iSide = i.acompanamientoId ?? null;
        const iDesc = i.descriptionSelections?.length
          ? '_' + i.descriptionSelections.map(s => `${s.descriptionId}:${s.optionId}`).sort().join(',')
          : '';
        const existingKey = (iSide ? `${iTitleBase}_${iSide}` : iTitleBase) + iDesc;
        return existingKey === baseKey;
      });
      
      let next;
      if (existingItemIndex !== -1) {
        next = current.map((i, index) => {
          if (index === existingItemIndex) {
            return { ...i, cantidad: i.cantidad + 1 };
          }
          return i;
        });
      } else {
        const station = side?.station ?? item?.station ?? null;
        next = [...current, {
          ...item,
          cantidad: 1,
          precio: precio,
          acompanamiento: side ? getBaseText(side.name) : null,
          acompanamientoId: side ? getBaseText(side.name) : null,
          descriptionSelections: descriptionSelections.length ? descriptionSelections : undefined,
          station,
          menuId: menuId || undefined
        }];
      }
      saveToStorage(next);
      return next;
    });
  }

  /** Obtiene la clave única de un item del carrito (para updateQuantity/removeItemByKey) */
  getItemKey(cartItem) {
    const titleBase = getBaseText(cartItem.title);
    let baseKey;
    if (cartItem.customQuantity != null) {
      baseKey = `${titleBase}_${cartItem.customQuantity}`;
    } else {
      const side = cartItem.acompanamientoId ?? null;
      const desc = cartItem.descriptionSelections?.length
        ? '_' + cartItem.descriptionSelections.map(s => `${s.descriptionId}:${s.optionId}`).sort().join(',')
        : '';
      baseKey = (side ? `${titleBase}_${side}` : titleBase) + desc;
    }
    const menuId = cartItem.menuId ?? null;
    return menuId ? `${menuId}_${baseKey}` : baseKey;
  }

  /** Obtiene items filtrados por menuId (null = todos, para compatibilidad) */
  getItemsForMenu(menuId) {
    const list = this.items;
    if (!menuId) return list;
    return list.filter((i) => (i.menuId ?? null) === menuId);
  }

  /**
   * Agrega un item al carrito con cantidad personalizada (para WEIGHT, VOLUME, etc.)
   * @param {Object} item - Item del menú a agregar
   * @param {number} quantity - Cantidad personalizada
   * @param {string} menuId - ID del menú (asocia el item al menú actual)
   */
  addItemWithQuantity(item, quantity, menuId = null) {
    if (!item.pricing) {
      throw new Error('El item debe tener pricing para usar cantidad personalizada');
    }
    
    if (quantity <= 0) {
      throw new Error('La cantidad debe ser mayor a 0');
    }
    
    // Calcular precio según la cantidad y el modo de pricing
    const precio = getPriceFromPricing(item.pricing, quantity);
    
    // Crear clave única: title + cantidad (para permitir múltiples cantidades diferentes)
    const itemTitleBase = getBaseText(item.title);
    const baseKey = `${itemTitleBase}_${quantity}`;
    
    itemsStore.update((current) => {
      const existingItemIndex = current.findIndex(i => {
        if ((i.menuId ?? null) !== (menuId ?? null)) return false;
        const iTitleBase = getBaseText(i.title);
        const existingKey = i.customQuantity
          ? `${iTitleBase}_${i.customQuantity}`
          : iTitleBase;
        return existingKey === baseKey;
      });
      
      const newItem = {
        ...item,
        cantidad: 1,
        customQuantity: quantity,
        precio: precio,
        pricing: item.pricing,
        menuId: menuId || undefined
      };
      const next = existingItemIndex !== -1
        ? current.map((i, idx) => (idx === existingItemIndex ? { ...i, cantidad: i.cantidad + 1 } : i))
        : [...current, newItem];
      saveToStorage(next);
      return next;
    });
  }

  /**
   * Elimina un item del carrito
   * @param {string} title - Título del item a eliminar
   */
  removeItem(title) {
    const titleBase = getBaseText(title);
    itemsStore.update((current) => {
      const next = current.filter(item => getBaseText(item.title) !== titleBase);
      saveToStorage(next);
      return next;
    });
  }

  /**
   * Actualiza la cantidad de un item
   * @param {string} itemKey - Clave única del item (usar getItemKey)
   * @param {number} cantidad - Nueva cantidad
   */
  updateQuantity(itemKey, cantidad) {
    if (cantidad <= 0) {
      this.removeItemByKey(itemKey);
      return;
    }
    
    itemsStore.update((current) => {
      const next = current.map(item => {
        const currentKey = this.getItemKey(item);
        if (currentKey === itemKey) {
          return { ...item, cantidad };
        }
        return item;
      });
      saveToStorage(next);
      return next;
    });
  }
  
  /**
   * Elimina un item del carrito por clave única
   * @param {string} itemKey - Clave única del item (usar getItemKey)
   */
  removeItemByKey(itemKey) {
    itemsStore.update((current) => {
      const next = current.filter(item => this.getItemKey(item) !== itemKey);
      saveToStorage(next);
      return next;
    });
  }

  /**
   * Calcula el total del carrito (opcionalmente filtrado por menuId)
   * @param {string} menuId - ID del menú; si se pasa, solo suma items de ese menú
   * @returns {number} Total del carrito
   */
  getTotal(menuId = null) {
    const list = menuId ? this.getItemsForMenu(menuId) : this.items;
    return list.reduce((total, item) => {
      // Si tiene cantidad personalizada y pricing, calcular precio dinámicamente
      if (item.customQuantity != null && item.pricing) {
        const precioUnitario = getPriceFromPricing(item.pricing, item.customQuantity);
        return total + (precioUnitario * (item.cantidad || 1));
      }
      // Precio normal (items con sides o UNIT sin customQuantity)
      return total + ((item.precio || 0) * (item.cantidad || 1));
    }, 0);
  }

  /**
   * Obtiene la cantidad total de items (opcionalmente filtrado por menuId)
   * @param {string} menuId - ID del menú; si se pasa, solo cuenta items de ese menú
   * @returns {number} Cantidad total de items
   */
  getTotalItems(menuId = null) {
    const list = menuId ? this.getItemsForMenu(menuId) : this.items;
    return list.reduce((total, item) => {
      // Para items con cantidad personalizada, sumar la cantidad personalizada
      if (item.customQuantity) {
        return total + item.customQuantity;
      }
      return total + item.cantidad;
    }, 0);
  }

  /**
   * Limpia el carrito (opcionalmente solo items del menuId indicado)
   * @param {string} menuId - Si se pasa, solo elimina items de ese menú; si no, limpia todo
   */
  clear(menuId = null) {
    if (!menuId) {
      itemsStore.set([]);
      saveToStorage([]);
      return;
    }
    itemsStore.update((current) => {
      const next = current.filter((i) => (i.menuId ?? null) !== menuId);
      saveToStorage(next);
      return next;
    });
  }

  /**
   * Genera el mensaje para WhatsApp con el pedido
   * @param {string} whatsappNumber - Número de WhatsApp
   * @param {string} nombreRetiro - Nombre de quien va a retirar
   * @param {string} horaRetiro - Hora de retiro (solo para PICKUP)
   * @param {string} deliveryAddress - Dirección de entrega (solo para DELIVERY)
   * @param {string} lang - Idioma ('ES', 'PT', 'EN')
   * @param {object} translations - Objeto con las traducciones de WhatsApp
   * @param {number} orderNumber - Número de orden (opcional)
   * @returns {string} URL de WhatsApp con el mensaje
   */
  generateWhatsAppMessage(whatsappNumber, nombreRetiro = '', horaRetiro = '', deliveryAddress = null, lang = 'ES', translations = null, orderNumber = null, menuId = null) {
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
    const itemsToUse = menuId ? this.getItemsForMenu(menuId) : this.items;
    itemsToUse.forEach((item, index) => {
      const itemTitle = getBaseText(item.title);
      message += `${index + 1}. ${itemTitle}`;
      if (item.acompanamiento) {
        message += ` (${item.acompanamiento})`;
      }
      if (item.descriptionSelections?.length) {
        const prefParts = item.descriptionSelections.map(s => s.optionId).join(', ');
        message += ` [${prefParts}]`;
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
    
    message += `\n*${t.orderTotal}: $${this.getTotal(menuId).toLocaleString(locale)}*\n\n`;
    
    // Agregar número de orden si está disponible
    if (orderNumber !== null) {
      message += lang === 'EN' 
        ? `📋 Order Number: ${orderNumber}\n\n`
        : lang === 'PT'
        ? `📋 Número do Pedido: ${orderNumber}\n\n`
        : `📋 Número de Orden: ${orderNumber}\n\n`;
    }
    
    // Información de entrega o retiro
    if (deliveryAddress) {
      // DELIVERY
      message += lang === 'EN' ? '📦 Delivery Information:\n' : lang === 'PT' ? '📦 Informações de Entrega:\n' : '📦 Información de Envío:\n';
      if (nombreRetiro) {
        message += `${t.pickupNameLabel} ${nombreRetiro}\n`;
      }
      message += `📍 ${lang === 'EN' ? 'Address' : lang === 'PT' ? 'Endereço' : 'Dirección'}: ${deliveryAddress}\n`;
      message += "\n";
    } else if (nombreRetiro || horaRetiro) {
      // PICKUP
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
    
    // Codificar el mensaje correctamente para WhatsApp
    // Normalizar el mensaje para asegurar codificación UTF-8 consistente
    // Esto previene problemas intermitentes con emojis
    const normalizedMessage = message.normalize('NFC'); // Normalizar a NFC (Canonical Composition)
    const encodedMessage = encodeURIComponent(normalizedMessage);
    const phoneNumber = whatsappNumber.replace(/[^0-9]/g, '');
    return `https://wa.me/${phoneNumber}?text=${encodedMessage}`;
  }
}

// Crear instancia única del store
export const cartStore = new CartStore();

