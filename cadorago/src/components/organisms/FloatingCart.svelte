<script>
  import { cartStore, itemsStore } from '../../stores/cartStore.svelte.js';
  import { get } from 'svelte/store';
  import Price from '../atoms/Price.svelte';
  import OrderLoader from '../atoms/OrderLoader.svelte';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { trackingStore } from '../../stores/trackingStore.svelte.js';
  import { t, language } from '../../lib/useLanguage';
  import { getPriceFromPricing, getCostFromPricing } from '../../services/menuData.js';
  import { geocodeAddress, autocompleteAddress, getStaticMapUrl } from '../../services/locationIQ.js';
  import { getSlugFromUrl, getUrlParams } from '../../services/restaurantData.js';
  import { getBaseText, getMultilingualText } from '../../lib/multilingual';
  import MenuItemDescription from '../atoms/MenuItemDescription.svelte';
  import { getEffectiveCurrency, formatPrice } from '../../lib/currency';
  
  const restaurantData = $derived(restaurantDataStore.value);
  const currency = $derived(getEffectiveCurrency(restaurantData));
  const currentLanguage = $derived($language);
  
  // Leer $itemsStore en $effect para que Svelte 5 rastree la dependencia
  let items = $state([]);
  let total = $state(0);
  let totalItems = $state(0);
  $effect(() => {
    const list = $itemsStore ?? [];
    items = list;
    total = list.reduce((sum, item) => {
      if (item.customQuantity != null && item.pricing) {
        const p = getPriceFromPricing(item.pricing, item.customQuantity);
        return sum + (p * (item.cantidad || 1));
      }
      return sum + ((item.precio || 0) * (item.cantidad || 1));
    }, 0);
    totalItems = list.reduce((sum, item) => {
      if (item.customQuantity) return sum + item.customQuantity;
      return sum + (item.cantidad || 0);
    }, 0);
  });
  
  // Opciones de delivery disponibles
  const deliveryOptions = $derived(restaurantData?.deliveryOptions || []);
  const hasDelivery = $derived(deliveryOptions.some(opt => opt.type === 'DELIVERY'));
  const hasPickup = $derived(deliveryOptions.some(opt => opt.type === 'PICKUP'));
  
  let isExpanded = $state(false);
  let showOrderForm = $state(false);
  let showClearConfirm = $state(false);
  let showOrderLoader = $state(false);
  let lastTrackingId = $state(null);
  let orderViewTransition = $state(false); // Controla la transición de la vista
  let deliveryType = $state(null); // 'DELIVERY' | 'PICKUP' | null
  let deliveryStep = $state(1); // 1: dirección, 2: datos personales
  let nombreRetiro = $state('');
  let phoneContact = $state('');
  let deliveryAddress = $state('');
  let addressNumber = $state(''); // Número de casa/departamento
  let addressNotes = $state(''); // Indicaciones adicionales
  let addressSuggestions = $state([]);
  let showSuggestions = $state(false);
  let searchingAddress = $state(false);
  let selectedAddress = $state(null);
  
  // Obtener zona horaria
  const timeZone = $derived(Intl.DateTimeFormat().resolvedOptions().timeZone);
  const timeZoneName = $derived(() => {
    try {
      const formatter = new Intl.DateTimeFormat('es-CL', {
        timeZoneName: 'short',
        timeZone: timeZone
      });
      const parts = formatter.formatToParts(new Date());
      const tzPart = parts.find(part => part.type === 'timeZoneName');
      return tzPart ? tzPart.value : timeZone;
    } catch {
      return timeZone;
    }
  });
  
  function handleQuantityChange(cartItem, event) {
    const quantity = parseInt(event.currentTarget.value) || 0;
    const itemTitleBase = getBaseText(cartItem.title);
    const itemKey = cartItem.acompanamientoId 
      ? `${itemTitleBase}_${cartItem.acompanamientoId}` 
      : itemTitleBase;
    cartStore.updateQuantity(itemKey, quantity);
  }
  
  function handleRemoveItem(cartItem) {
    const itemTitleBase = getBaseText(cartItem.title);
    const itemKey = cartItem.acompanamientoId 
      ? `${itemTitleBase}_${cartItem.acompanamientoId}` 
      : itemTitleBase;
    cartStore.removeItemByKey(itemKey);
  }
  
  function handleSendOrderClick() {
    // Si solo hay una opción, seleccionarla automáticamente
    if (hasDelivery && !hasPickup) {
      deliveryType = 'DELIVERY';
    } else if (hasPickup && !hasDelivery) {
      deliveryType = 'PICKUP';
    }
    showOrderForm = true;
    // Activar transición después de un pequeño delay para suavidad
    setTimeout(() => {
      orderViewTransition = true;
    }, 10);
  }
  
  function handleCancelOrder() {
    // Primero ocultar la vista con transición
    orderViewTransition = false;
    // Luego limpiar después de la animación
    setTimeout(() => {
      showOrderForm = false;
      deliveryType = null;
      deliveryStep = 1;
      nombreRetiro = '';
      phoneContact = '';
      deliveryAddress = '';
      addressNumber = '';
      addressNotes = '';
      addressSuggestions = [];
      showSuggestions = false;
      selectedAddress = null;
    }, 300); // Duración de la transición
  }
  
  function handleNextStep() {
    if (deliveryStep === 1 && deliveryType === 'DELIVERY') {
      // Validar que haya una dirección seleccionada de las sugerencias
      if (!selectedAddress) {
        alert('Por favor selecciona una dirección de las sugerencias para ver el mapa y continuar');
        return;
      }
      deliveryStep = 2;
    }
  }
  
  function handleBackStep() {
    if (deliveryStep === 2) {
      deliveryStep = 1;
    }
  }
  
  // Autocompletado de direcciones
  let addressInputTimeout;
  let originalAddressInput = $state(''); // Guardar el texto original del usuario
  
  async function handleAddressInput(event) {
    const value = event.target.value;
    deliveryAddress = value;
    originalAddressInput = value; // Guardar el texto original
    selectedAddress = null;
    
    if (value.length < 3) {
      addressSuggestions = [];
      showSuggestions = false;
      return;
    }
    
    // Debounce para evitar demasiadas llamadas
    clearTimeout(addressInputTimeout);
    addressInputTimeout = setTimeout(async () => {
      searchingAddress = true;
      try {
        const suggestions = await autocompleteAddress(value);
        addressSuggestions = suggestions;
        showSuggestions = suggestions.length > 0;
      } catch (error) {
        console.error('Error en autocompletado:', error);
        addressSuggestions = [];
        showSuggestions = false;
      } finally {
        searchingAddress = false;
      }
    }, 300);
  }
  
  /**
   * Extrae números de una dirección (número de casa, departamento, etc.)
   * @param {string} address - Dirección a analizar
   * @returns {string} Números extraídos o cadena vacía
   */
  function extractAddressNumber(address) {
    if (!address) return '';
    
    // Buscar patrones comunes de números en direcciones chilenas:
    // - "inglaterra 59" -> "59"
    // - "inglaterra 59, la florida" -> "59"
    // - "calle 123, depto 45" -> "123, depto 45"
    // - "calle 123-45" -> "123-45"
    // - "calle 123 a" -> "123"
    
    // Buscar el primer número que aparece después del nombre de la calle
    // Patrón: palabra(s) + número + (opcional: letra, guión, coma, más números)
    const streetNumberPattern = /\b(\d+[a-z]?)\b/i;
    const match = address.match(streetNumberPattern);
    
    if (match) {
      const numberPart = match[1];
      const numberIndex = address.indexOf(numberPart);
      
      // Obtener todo lo que viene después del número hasta el final o hasta una coma que indique comuna
      const afterNumber = address.substring(numberIndex + numberPart.length).trim();
      
      // Si hay texto después del número que parece ser parte de la dirección (depto, casa, etc.)
      if (afterNumber) {
        const deptoPattern = /^[,]?\s*(depto|departamento|casa|block|torre|piso|oficina)\s*[:\s]*(\d+[a-z]?)/i;
        const deptoMatch = afterNumber.match(deptoPattern);
        if (deptoMatch) {
          return `${numberPart}, ${deptoMatch[1]} ${deptoMatch[2]}`;
        }
        
        // Si hay un guión seguido de más números (ej: "123-45")
        const hyphenPattern = /^[-]\s*(\d+[a-z]?)/i;
        const hyphenMatch = afterNumber.match(hyphenPattern);
        if (hyphenMatch) {
          return `${numberPart}-${hyphenMatch[1]}`;
        }
      }
      
      return numberPart;
    }
    
    return '';
  }
  
  /**
   * Combina la dirección normalizada de LocationIQ con el número extraído del input original
   * @param {string} normalizedAddress - Dirección normalizada de LocationIQ
   * @param {string} originalInput - Texto original que escribió el usuario
   * @returns {string} Dirección combinada
   */
  function combineAddressWithNumber(normalizedAddress, originalInput) {
    if (!normalizedAddress) return originalInput;
    if (!originalInput) return normalizedAddress;
    
    const extractedNumber = extractAddressNumber(originalInput);
    if (!extractedNumber) {
      return normalizedAddress;
    }
    
    // Verificar si la dirección normalizada ya contiene el número
    const normalizedLower = normalizedAddress.toLowerCase();
    const numberLower = extractedNumber.toLowerCase();
    
    // Extraer solo el número base (sin letras adicionales) para comparación
    const baseNumber = extractedNumber.match(/\d+/)?.[0] || '';
    
    // Si el número ya está en la dirección normalizada, no duplicarlo
    if (baseNumber && normalizedLower.includes(baseNumber)) {
      // Verificar si el número está en el contexto correcto (no como parte de un código postal)
      const numberIndex = normalizedLower.indexOf(baseNumber);
      const beforeNumber = normalizedLower.substring(Math.max(0, numberIndex - 10), numberIndex);
      const afterNumber = normalizedLower.substring(numberIndex + baseNumber.length, numberIndex + baseNumber.length + 10);
      
      // Si el número está precedido por palabras comunes de calles, asumimos que ya está incluido
      if (/calle|avenida|pasaje|plaza|boulevard|route|street|avenue/i.test(beforeNumber) || 
          /^[\s,]/i.test(afterNumber)) {
        return normalizedAddress;
      }
    }
    
    // Combinar: tomar la calle de la dirección normalizada y agregar el número
    // LocationIQ generalmente devuelve: "Calle, Comuna, Región, País"
    // Queremos: "Calle [número], Comuna, Región, País"
    
    // Buscar la primera coma para insertar el número antes de la comuna
    const firstCommaIndex = normalizedAddress.indexOf(',');
    if (firstCommaIndex > 0) {
      const streetName = normalizedAddress.substring(0, firstCommaIndex).trim();
      const rest = normalizedAddress.substring(firstCommaIndex);
      
      // Insertar el número después del nombre de la calle
      return `${streetName} ${extractedNumber}${rest}`;
    }
    
    // Si no hay coma, simplemente agregar el número al final
    return `${normalizedAddress} ${extractedNumber}`;
  }
  
  function selectAddress(suggestion) {
    // Combinar la dirección normalizada con el número del input original
    const combinedAddress = combineAddressWithNumber(suggestion.display_name, originalAddressInput);
    deliveryAddress = combinedAddress;
    
    // Actualizar el display_name en la sugerencia para mantener consistencia
    selectedAddress = {
      ...suggestion,
      display_name: combinedAddress
    };
    
    addressSuggestions = [];
    showSuggestions = false;
    // No avanzar automáticamente, dejar que el usuario vea el mapa
  }
  
  // Generar URL del mapa estático para la dirección seleccionada
  const selectedAddressMapUrl = $derived.by(() => {
    if (!selectedAddress) {
      console.log('No hay selectedAddress');
      return null;
    }
    if (!selectedAddress.lat || !selectedAddress.lon) {
      console.log('selectedAddress sin lat/lon:', selectedAddress);
      return null;
    }
    // Usar la función del servicio para generar la URL correctamente
    const url = getStaticMapUrl(selectedAddress.lat, selectedAddress.lon, 800, 400);
    console.log('URL del mapa generada:', url);
    return url;
  });
  
  function handleAddressBlur() {
    // Esperar un poco antes de ocultar sugerencias para permitir el click
    setTimeout(() => {
      showSuggestions = false;
    }, 200);
  }
  
  /**
   * Obtiene el menuId desde la URL o los datos del restaurante
   * @returns {string | null} ID del menú
   */
  function getMenuId() {
    // Intentar obtener desde los datos del restaurante (viene del backend)
    if (restaurantData?.id) {
      return restaurantData.id;
    }
    
    // Fallback: intentar obtener desde URL params (método legacy)
    const { menuID } = getUrlParams();
    if (menuID) {
      return menuID;
    }
    
    // Si se está usando slug o menuId, el id debería venir en restaurantData
    // Si no está, intentar obtener desde la URL (último recurso)
    const urlData = getSlugFromUrl();
    if (urlData && urlData.value) {
      // Si es un menuId, usarlo directamente
      if (urlData.isMenuId) {
        return urlData.value;
      }
      // Si es un slug, el id debería venir en restaurantData
      console.warn('No se pudo obtener el menuId. El backend debería devolver el id del menú.');
    }
    
    return null;
  }

  /**
   * Envía la orden a la API de micartapro-backend
   * @returns {Promise<{orderNumber: number, trackingId: string}>} Número de orden y código de seguimiento
   */
  async function sendOrderToAPI() {
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 
      (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
        ? 'http://localhost:8082'
        : 'https://micartapro-backend-27303662337.us-central1.run.app');
    
    const menuId = getMenuId();
    if (!menuId) {
      throw new Error('No se pudo obtener el ID del menú');
    }

    // Mapear items del carrito al formato de la API (orden: unit, quantity, unitPrice, totalPrice, unitCost, totalCost, pricingMode, productName)
    const items = get(itemsStore).map(item => {
      // Si tiene acompañamiento, usar pricing del side; si no, del item
      let pricing = item.pricing || {};
      if (item.acompanamientoId && item.sides?.length) {
        const side = item.sides.find(s => getBaseText(s?.name) === item.acompanamientoId);
        if (side?.pricing) pricing = side.pricing;
      }
      const pricingMode = pricing.mode === 'WEIGHT' || pricing.mode === 'VOLUME' || pricing.mode === 'LENGTH' || pricing.mode === 'AREA'
        ? pricing.mode
        : 'UNIT';

      const unit = pricing.unit || 'EACH';
      const quantity = item.customQuantity || item.cantidad;

      let unitPrice = 0;
      let totalPrice = 0;
      let unitCost = 0;
      let totalCost = 0;

      if (item.customQuantity != null && item.pricing) {
        unitPrice = pricing.pricePerUnit || 0;
        const precioPorUnidad = getPriceFromPricing(item.pricing, item.customQuantity);
        totalPrice = precioPorUnidad * (item.cantidad || 1);
        unitCost = pricing.costPerUnit ?? 0;
        totalCost = getCostFromPricing(item.pricing, item.customQuantity) * (item.cantidad || 1);
      } else {
        unitPrice = item.precio || 0;
        totalPrice = (item.precio || 0) * (item.cantidad || 1);
        unitCost = pricing.costPerUnit ?? 0;
        totalCost = getCostFromPricing(pricing, item.cantidad || 1);
      }

      return {
        unit,
        quantity,
        unitPrice,
        totalPrice,
        unitCost,
        totalCost,
        pricingMode,
        productName: getBaseText(item.title) + (item.acompanamiento ? ` (${item.acompanamiento})` : ''),
        station: item.station ?? ''
      };
    });

    // Calcular totales
    const subtotal = cartStore.getTotal();
    const deliveryFee = deliveryType === 'DELIVERY' ? 2000 : 0;
    const total = subtotal + deliveryFee;
    const totalCost = items.reduce((sum, it) => sum + (it.totalCost || 0), 0);

    // Fecha/hora solicitada en ISO string (PICKUP: ahora; DELIVERY: ahora + 1h)
    const now = new Date();
    const requestedTime = deliveryType === 'DELIVERY'
      ? (() => { const t = new Date(now); t.setHours(t.getHours() + 1); return t.toISOString(); })()
      : now.toISOString();

    // fulfillment: type, requestedTime, address (solo DELIVERY), contact
    const fulfillment = {
      type: deliveryType,
      requestedTime,
      address: {
        rawAddress: '',
        coordinates: { latitude: 0, longitude: 0 },
        deliveryDetails: { unit: '', notes: '' }
      },
      contact: {
        fullName: nombreRetiro.trim(),
        phone: (phoneContact || '').trim(),
        email: ''
      }
    };

    if (deliveryType === 'DELIVERY' && selectedAddress) {
      const rawAddress = selectedAddress.display_name || deliveryAddress;
      fulfillment.address = {
        rawAddress,
        coordinates: {
          latitude: parseFloat(selectedAddress.lat) || 0,
          longitude: parseFloat(selectedAddress.lon) || 0
        },
        deliveryDetails: {
          unit: (addressNumber || '').trim(),
          notes: (addressNotes || '').trim()
        }
      };
    }

    const orderData = {
      createdAt: '', // el backend lo sobrescribe siempre
      items,
      totals: {
        total,
        currency: currency,
        subtotal,
        deliveryFee,
        totalCost
      },
      fulfillment,
      businessInfo: {
        whatsapp: restaurantData?.businessInfo?.whatsapp || '',
        businessName: restaurantData?.businessInfo?.businessName || 'cadorago'
      }
    };

    const response = await fetch(`${API_BASE_URL}/menu/${menuId}/orders`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(orderData)
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Error al crear la orden: ${response.status} ${errorText}`);
    }

    const result = await response.json();
    return { orderNumber: result.orderNumber, trackingId: result.trackingId };
  }
  
  async function handleConfirmOrder() {
    // Validar tipo de entrega
    if (!deliveryType) {
      alert($t.cart.selectDeliveryType);
      return;
    }
    
    // Validar nombre y teléfono
    if (!nombreRetiro.trim()) {
      alert($t.cart.completeFields);
      return;
    }
    if (!phoneContact.trim()) {
      alert($t.cart.completeFields);
      return;
    }
    
    // Validar según el tipo
    if (deliveryType === 'PICKUP') {
      // Sin hora de retiro: solo nombre y teléfono
    } else if (deliveryType === 'DELIVERY') {
      // Validar que estemos en el paso 2
      if (deliveryStep === 1) {
        handleNextStep();
        return;
      }
      
      if (!deliveryAddress.trim()) {
        alert($t.cart.completeFields);
        return;
      }
      
      // Si no hay dirección seleccionada, intentar geocodificar
      if (!selectedAddress) {
        try {
          searchingAddress = true;
          selectedAddress = await geocodeAddress(deliveryAddress);
        } catch (error) {
          alert($t.cart.addressNotFound);
          searchingAddress = false;
          return;
        } finally {
          searchingAddress = false;
        }
      }
    }
    
    const horaConZona = null; // Sin hora de retiro en cadorago
    
    // Construir dirección completa para DELIVERY
    let fullDeliveryAddress = null;
    if (deliveryType === 'DELIVERY' && selectedAddress) {
      fullDeliveryAddress = selectedAddress.display_name || deliveryAddress;
      if (addressNumber.trim()) {
        fullDeliveryAddress += `, ${addressNumber.trim()}`;
      }
      if (addressNotes.trim()) {
        // Normalizar el emoji para asegurar codificación consistente
        const notesEmoji = '📝'.normalize('NFC');
        fullDeliveryAddress += `\n${notesEmoji} Indicaciones: ${addressNotes.trim()}`;
      }
    }
    
    // Mostrar loader
    showOrderForm = false;
    showOrderLoader = true;
    lastTrackingId = null;
    
    let orderResult = null;
    try {
      // Enviar orden a la API
      orderResult = await sendOrderToAPI();
    } catch (error) {
      console.error('Error al enviar orden a la API:', error);
      alert('Error al procesar la orden. Por favor, intenta nuevamente.');
      showOrderLoader = false;
      showOrderForm = true;
      return;
    }
    
    lastTrackingId = orderResult.trackingId || null;
    const menuId = getMenuId();
    trackingStore.addTracking(orderResult.trackingId || '', new Date().toISOString(), menuId);

    // Limpiar carrito y formulario antes de redirigir
    cartStore.clear();
    showOrderLoader = false;
    deliveryType = null;
    deliveryStep = 1;
    nombreRetiro = '';
    phoneContact = '';
    deliveryAddress = '';
    addressNumber = '';
    addressNotes = '';
    addressSuggestions = [];
    showSuggestions = false;
    selectedAddress = null;
    isExpanded = false;

    // Redirigir a la vista de tracking (con menuId para filtrar por menú)
    const m = menuId ? `?m=${encodeURIComponent(menuId)}` : '';
    window.location.href = `/track/${encodeURIComponent(orderResult.trackingId || '')}${m}`;
  }
  
  function handleClearCart() {
    showClearConfirm = true;
  }
  
  function confirmClearCart() {
    cartStore.clear();
    isExpanded = false;
    showClearConfirm = false;
  }
  
  function cancelClearCart() {
    showClearConfirm = false;
  }
  
  function toggleExpanded() {
    isExpanded = !isExpanded;
  }
</script>

{#if items.length > 0}
  <div class="fixed bottom-0 left-0 right-0 z-50 bg-white shadow-2xl border-t-2 border-gray-300">
    <!-- Resumen compacto (siempre visible cuando hay items) -->
    <div class="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-12 py-4 sm:py-5">
      <div class="flex items-center justify-between gap-4">
        <!-- Información del pedido -->
        <div class="flex items-center gap-4 sm:gap-6 flex-1">
          <button
            onclick={toggleExpanded}
            class="flex items-center gap-2 sm:gap-3 text-gray-800 hover:text-gray-900 transition-colors"
            aria-label={isExpanded ? "Contraer carrito" : "Expandir carrito"}
          >
            <svg 
              class="w-5 h-5 sm:w-6 sm:h-6 transition-transform {isExpanded ? 'rotate-180' : ''}" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
            <div class="text-left">
              <p class="text-sm sm:text-base font-semibold text-gray-800">
                {$t.cart.yourOrder}
              </p>
              <p class="text-xs sm:text-sm text-gray-600">
                {totalItems} {totalItems === 1 ? $t.cart.item : $t.cart.items}
              </p>
            </div>
          </button>
          
          <div class="text-right">
            <p class="text-xs sm:text-sm text-gray-600 mb-1">{$t.cart.total}</p>
            {#key items.length + total}
              <span class="font-bold text-gray-900 text-lg sm:text-xl lg:text-2xl">{formatPrice(total, currency)}</span>
            {/key}
          </div>
        </div>
        
        <!-- Botones de acción -->
        <div class="flex items-center gap-2 sm:gap-3">
          <button
            onclick={handleClearCart}
            class="px-3 py-2 sm:px-4 sm:py-2.5 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors flex items-center gap-2"
            aria-label={$t.cart.clearOrder}
          >
            <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            <span class="hidden sm:inline text-sm font-medium">{$t.cart.clearOrder}</span>
          </button>
          
          <button
            onclick={handleSendOrderClick}
            class="px-4 py-2 sm:px-6 sm:py-3 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors flex items-center gap-2 sm:gap-3 font-semibold text-sm sm:text-base"
            aria-label="Realizar pedido"
          >
            <span>{$t.cart.placeOrder}</span>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Carrito expandido -->
    {#if isExpanded}
      <div class="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-12 pb-4 sm:pb-5 border-t border-gray-300 bg-white/95">
        <div class="pt-4 sm:pt-6 max-h-[60vh] overflow-y-auto">
          <div class="space-y-3 sm:space-y-4">
            {#each items as cartItem}
              {@const unitLabel = cartItem.customQuantity && cartItem.pricing
                ? (cartItem.pricing.unit === 'GRAM' ? 'g' : 
                   cartItem.pricing.unit === 'KILOGRAM' ? 'kg' :
                   cartItem.pricing.unit === 'MILLILITER' ? 'ml' :
                   cartItem.pricing.unit === 'LITER' ? 'L' :
                   cartItem.pricing.unit === 'METER' ? 'm' :
                   cartItem.pricing.unit === 'SQUARE_METER' ? 'm²' : '')
                : ''}
              {@const needsDecimals = cartItem.customQuantity && cartItem.pricing
                ? (cartItem.pricing.unit === 'KILOGRAM' || 
                   cartItem.pricing.unit === 'LITER' ||
                   cartItem.pricing.unit === 'METER' ||
                   cartItem.pricing.unit === 'SQUARE_METER')
                : false}
              {@const itemPrice = cartItem.customQuantity && cartItem.pricing 
                ? getPriceFromPricing(cartItem.pricing, cartItem.customQuantity)
                : (cartItem.precio * cartItem.cantidad)}
              <div class="bg-gray-50 rounded-lg p-4 sm:p-5 border border-gray-200">
                <div class="flex justify-between items-start gap-4 mb-3">
                  <div class="flex-1">
                    <h3 class="text-base sm:text-lg font-bold text-gray-800 mb-1">
                      {getMultilingualText(cartItem.title)}
                      {#if cartItem.acompanamiento}
                        <span class="text-sm font-normal text-gray-600">({cartItem.acompanamiento})</span>
                      {/if}
                    </h3>
                    <MenuItemDescription description={cartItem.description} className="text-xs sm:text-sm text-gray-600" />
                  </div>
                  <button
                    onclick={() => handleRemoveItem(cartItem)}
                    class="text-red-600 hover:text-red-800 text-lg font-bold flex-shrink-0"
                    aria-label="Eliminar item"
                  >
                    ×
                  </button>
                </div>
                
                <div class="flex justify-between items-center">
                  <div class="flex items-center gap-2 sm:gap-3">
                    <label for="quantity-{getBaseText(cartItem.title)}" class="text-xs sm:text-sm text-gray-700 font-medium">
                      {$t.cart.quantity}
                    </label>
                    {#if cartItem.customQuantity && cartItem.pricing}
                      <span class="text-sm sm:text-base font-semibold text-gray-800">
                        {needsDecimals ? cartItem.customQuantity.toFixed(1) : cartItem.customQuantity} {unitLabel}
                      </span>
                    {:else}
                      <input
                        id="quantity-{getBaseText(cartItem.title)}"
                        type="number"
                        min="1"
                        value={cartItem.cantidad}
                        oninput={(e) => handleQuantityChange(cartItem, e)}
                        class="w-14 sm:w-16 px-2 py-1 border border-gray-300 rounded text-center text-xs sm:text-sm"
                      />
                    {/if}
                  </div>
                  <div class="text-right">
                    <p class="text-sm sm:text-base text-gray-800 font-semibold">
                      {formatPrice(itemPrice, currency)}
                    </p>
                    {#if !cartItem.customQuantity && cartItem.cantidad > 1}
                      <p class="text-xs text-gray-500">
                        {formatPrice(cartItem.precio, currency)} c/u
                      </p>
                    {/if}
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>
    {/if}
  </div>
{/if}

<!-- Modal de información de entrega/retiro -->
{#if showOrderForm}
  <div 
    class="fixed inset-0 bg-white z-[100] transition-transform duration-300 ease-in-out {orderViewTransition ? 'translate-x-0' : 'translate-x-full'}"
    role="dialog"
    aria-modal="true"
    aria-labelledby="order-form-title"
    tabindex="-1"
    onkeydown={(e) => {
      if (e.key === 'Escape') {
        handleCancelOrder();
      }
    }}
  >
    <div class="h-full w-full overflow-y-auto">
      <div class="max-w-2xl mx-auto p-4 sm:p-6 md:p-8">
        <!-- Header con botón de cerrar -->
        <div class="flex items-center justify-between mb-6 sm:mb-8 sticky top-0 bg-white z-10 pb-4 border-b border-gray-200 -mx-4 sm:-mx-6 md:-mx-8 px-4 sm:px-6 md:px-8">
          <button
            onclick={handleCancelOrder}
            class="p-2 text-gray-500 hover:text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
            aria-label="Cerrar"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h3 id="order-form-title" class="text-xl sm:text-2xl font-bold text-gray-800 flex-1 text-center">
            {deliveryType === 'DELIVERY' 
              ? (deliveryStep === 1 ? $t.cart.step1Title : $t.cart.step2Title)
              : deliveryType === 'PICKUP' 
              ? $t.cart.pickupInfo 
              : 'Información del Pedido'}
          </h3>
          <div class="w-10"></div> <!-- Spacer para centrar el título -->
        </div>
      
      <!-- Indicador de pasos para DELIVERY -->
      {#if deliveryType === 'DELIVERY'}
        <div class="flex items-center justify-center mb-6">
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-2">
              <div class="w-9 h-9 rounded-full flex items-center justify-center text-sm font-semibold transition-colors {deliveryStep >= 1 ? 'bg-green-500 text-white shadow-md' : 'bg-gray-200 text-gray-600'}">
                1
              </div>
              <span class="text-sm font-medium {deliveryStep >= 1 ? 'text-green-600' : 'text-gray-500'}">{$t.cart.step1Label}</span>
            </div>
            <div class="w-16 h-1 rounded-full transition-colors {deliveryStep >= 2 ? 'bg-green-500' : 'bg-gray-200'}"></div>
            <div class="flex items-center gap-2">
              <div class="w-9 h-9 rounded-full flex items-center justify-center text-sm font-semibold transition-colors {deliveryStep >= 2 ? 'bg-green-500 text-white shadow-md' : 'bg-gray-200 text-gray-600'}">
                2
              </div>
              <span class="text-sm font-medium {deliveryStep >= 2 ? 'text-green-600' : 'text-gray-500'}">{$t.cart.step2Label}</span>
            </div>
          </div>
        </div>
      {/if}
      
      <div class="space-y-4 sm:space-y-5" onclick={(e) => e.stopPropagation()}>
        <!-- Selección de tipo de entrega (si hay ambas opciones) -->
        {#if hasDelivery && hasPickup && !deliveryType}
          <div>
            <label class="block text-sm sm:text-base font-medium text-gray-700 mb-3">
              {$t.cart.deliveryType}
            </label>
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                onclick={() => deliveryType = 'DELIVERY'}
                class="px-4 py-3 border-2 border-gray-300 rounded-lg hover:border-green-500 hover:bg-green-50 transition-colors text-sm sm:text-base font-medium {deliveryType === 'DELIVERY' ? 'border-green-500 bg-green-50' : ''}"
              >
                📦 {$t.cart.delivery}
              </button>
              <button
                type="button"
                onclick={() => deliveryType = 'PICKUP'}
                class="px-4 py-3 border-2 border-gray-300 rounded-lg hover:border-green-500 hover:bg-green-50 transition-colors text-sm sm:text-base font-medium {deliveryType === 'PICKUP' ? 'border-green-500 bg-green-50' : ''}"
              >
                🏪 {$t.cart.pickup}
              </button>
            </div>
          </div>
        {/if}
        
        <!-- PASO 1: Dirección de entrega (solo DELIVERY) -->
        {#if deliveryType === 'DELIVERY' && deliveryStep === 1}
          <div class="space-y-4">
            <div class="relative">
              <label for="delivery-address" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                {$t.cart.deliveryAddress}
              </label>
              <input
                id="delivery-address"
                type="text"
                value={deliveryAddress}
                oninput={handleAddressInput}
                onblur={handleAddressBlur}
                onfocus={() => showSuggestions = addressSuggestions.length > 0}
                placeholder={$t.cart.deliveryAddressPlaceholder}
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
                required
              />
              {#if searchingAddress}
                <p class="text-xs text-gray-500 mt-1">{$t.cart.searchingAddress}</p>
              {/if}
              {#if showSuggestions && addressSuggestions.length > 0}
                <div class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-96 overflow-y-auto">
                  {#each addressSuggestions as suggestion}
                    <button
                      type="button"
                      onclick={() => selectAddress(suggestion)}
                      class="w-full text-left hover:bg-gray-50 border-b border-gray-100 last:border-b-0 transition-colors"
                    >
                      <div class="flex items-start gap-3 p-3">
                        {#if suggestion.mapUrl}
                          <div class="flex-shrink-0 w-24 h-20 rounded overflow-hidden border border-gray-200">
                            <img 
                              src={suggestion.mapUrl} 
                              alt="Mapa de ubicación"
                              class="w-full h-full object-cover"
                              loading="lazy"
                            />
                          </div>
                        {/if}
                        <div class="flex-1 min-w-0">
                          <p class="text-sm text-gray-800 font-medium leading-tight">
                            {suggestion.display_name}
                          </p>
                        </div>
                      </div>
                    </button>
                  {/each}
                </div>
              {/if}
            </div>
            
            <!-- Mapa con pin cuando hay dirección seleccionada -->
            {#if selectedAddress}
              <div class="mt-4 sm:mt-6">
                <label class="block text-sm sm:text-base font-medium text-gray-700 mb-3">
                  {$t.cart.confirmAddress}
                </label>
                {#if selectedAddressMapUrl}
                  <div class="w-full rounded-lg sm:rounded-lg overflow-hidden border-2 border-green-500 shadow-lg bg-gray-100">
                    <img 
                      src={selectedAddressMapUrl} 
                      alt="Mapa de ubicación seleccionada"
                      class="w-full h-auto"
                      onerror={(e) => {
                        console.error('Error cargando mapa. URL:', selectedAddressMapUrl);
                        console.error('selectedAddress:', selectedAddress);
                        if (e.target && e.target instanceof HTMLImageElement) {
                          e.target.style.display = 'none';
                        }
                      }}
                      onload={() => {
                        console.log('✅ Mapa cargado correctamente');
                      }}
                    />
                  </div>
                {:else}
                  <div class="w-full h-64 sm:h-64 bg-gray-200 rounded-lg border-2 border-gray-300 flex items-center justify-center">
                    <p class="text-gray-500">Generando mapa...</p>
                  </div>
                {/if}
                <div class="mt-3 p-3 bg-gray-50 rounded-lg border border-gray-200">
                  <p class="text-sm sm:text-base text-gray-800 font-medium">
                    📍 {selectedAddress.display_name}
                  </p>
                </div>
              </div>
            {/if}
          </div>
        {/if}
        
        <!-- PASO 2: Información de contacto (solo DELIVERY paso 2) -->
        {#if deliveryType === 'DELIVERY' && deliveryStep === 2}
          <div class="space-y-4">
            <div>
              <label class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                Dirección confirmada
              </label>
              <div class="p-3 bg-gray-50 rounded-lg border border-gray-200">
                <p class="text-sm text-gray-800 font-medium">{selectedAddress?.display_name || deliveryAddress}</p>
              </div>
            </div>
            
            <div>
              <label for="nombre-retiro" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                Nombre para la entrega *
              </label>
              <input
                id="nombre-retiro"
                type="text"
                bind:value={nombreRetiro}
                placeholder={$t.cart.nameFormatExample}
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
                required
              />
            </div>
            
            <div>
              <label for="phone-contact-delivery" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                {$t.cart.phone} *
              </label>
              <input
                id="phone-contact-delivery"
                type="tel"
                bind:value={phoneContact}
                placeholder={$t.cart.phonePlaceholder}
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
                required
              />
            </div>
            
            <div>
              <label for="address-number" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                {$t.cart.addressNumber}
              </label>
              <input
                id="address-number"
                type="text"
                bind:value={addressNumber}
                placeholder={$t.cart.addressNumberPlaceholder}
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>
            
            <div>
              <label for="address-notes" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                {$t.cart.addressNotes}
              </label>
              <textarea
                id="address-notes"
                bind:value={addressNotes}
                placeholder={$t.cart.addressNotesPlaceholder}
                rows="3"
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500 resize-none"
              ></textarea>
            </div>
          </div>
        {/if}
        
        <!-- Formulario PICKUP (sin pasos) -->
        {#if deliveryType === 'PICKUP'}
          <div class="space-y-4">
            <div>
              <label for="nombre-retiro-pickup" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                {$t.cart.pickupName}
              </label>
              <input
                id="nombre-retiro-pickup"
                type="text"
                bind:value={nombreRetiro}
                placeholder={$t.cart.nameFormatExample}
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
                required
              />
            </div>
            
            <div>
              <label for="phone-contact-pickup" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
                {$t.cart.phone} *
              </label>
              <input
                id="phone-contact-pickup"
                type="tel"
                bind:value={phoneContact}
                placeholder={$t.cart.phonePlaceholder}
                class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
                required
              />
            </div>
            
          </div>
        {/if}
      </div>
      
      <div class="flex gap-3 sm:gap-4 mt-6 sm:mt-8" onclick={(e) => e.stopPropagation()}>
        {#if deliveryType === 'DELIVERY' && deliveryStep === 2}
          <button
            onclick={handleBackStep}
            class="flex-1 px-4 py-2 sm:py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm sm:text-base"
          >
            {$t.cart.back}
          </button>
        {:else}
          <button
            onclick={handleCancelOrder}
            class="flex-1 px-4 py-2 sm:py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm sm:text-base"
          >
            {$t.cart.cancel}
          </button>
        {/if}
        
        {#if deliveryType === 'DELIVERY' && deliveryStep === 1}
          <button
            onclick={handleNextStep}
            disabled={!selectedAddress}
            class="flex-1 px-4 py-2 sm:py-3 bg-green-500 hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed text-white rounded-lg transition-colors font-semibold text-sm sm:text-base"
          >
            {$t.cart.next}
          </button>
        {:else}
          <button
            onclick={handleConfirmOrder}
            disabled={searchingAddress}
            class="flex-1 px-4 py-2 sm:py-3 bg-green-500 hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed text-white rounded-lg transition-colors font-semibold text-sm sm:text-base flex items-center justify-center gap-2"
          >
            <span>{$t.cart.placeOrder}</span>
          </button>
        {/if}
      </div>
      </div>
    </div>
  </div>
{/if}

<!-- Loader de preparación de pedido -->
{#if showOrderLoader}
  <OrderLoader
    message={$t.cart.preparingOrder}
    redirectingMessage={$t.cart.redirectingWhatsApp}
    trackingId={lastTrackingId}
  />
{/if}

<!-- Modal de confirmación para limpiar carrito -->
{#if showClearConfirm}
  <div 
    class="fixed inset-0 bg-black/50 z-[100] flex items-center justify-center p-4" 
    role="dialog"
    aria-modal="true"
    aria-labelledby="clear-confirm-title"
    tabindex="-1"
    onclick={cancelClearCart}
    onkeydown={(e) => {
      if (e.key === 'Escape') {
        cancelClearCart();
      }
    }}
  >
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6 sm:p-8" onclick={(e) => e.stopPropagation()}>
      <h3 id="clear-confirm-title" class="text-xl sm:text-2xl font-bold text-gray-800 mb-4">
        {$t.cart.confirmClear}
      </h3>
      
      <p class="text-gray-600 mb-6">
        {$t.cart.confirmClearMessage}
      </p>
      
      <div class="flex gap-3 sm:gap-4">
        <button
          onclick={cancelClearCart}
          class="flex-1 px-4 py-2 sm:py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm sm:text-base"
        >
          {$t.cart.no}
        </button>
        <button
          onclick={confirmClearCart}
          class="flex-1 px-4 py-2 sm:py-3 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors font-semibold text-sm sm:text-base"
        >
          {$t.cart.yes}
        </button>
      </div>
    </div>
  </div>
{/if}
