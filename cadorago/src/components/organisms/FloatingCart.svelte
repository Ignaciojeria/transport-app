<script>
  import { cartStore } from '../../stores/cartStore.svelte.js';
  import Price from '../atoms/Price.svelte';
  import WhatsAppIcon from '../atoms/WhatsAppIcon.svelte';
  import { getRestaurantData } from '../../services/restaurantData.js';
  
  const restaurantData = getRestaurantData();
  
  // Valores derivados reactivos
  const items = $derived(cartStore.items);
  const total = $derived(cartStore.getTotal());
  const totalItems = $derived(cartStore.getTotalItems());
  
  let isExpanded = $state(false);
  let showOrderForm = $state(false);
  let showClearConfirm = $state(false);
  let nombreRetiro = $state('');
  let horaRetiro = $state('');
  
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
  
  function formatTimeInput(value) {
    // Remover todo excepto números
    let numbers = value.replace(/\D/g, '');
    
    // Limitar a 4 dígitos
    if (numbers.length > 4) {
      numbers = numbers.slice(0, 4);
    }
    
    // Formatear como HH:MM
    if (numbers.length >= 3) {
      return `${numbers.slice(0, 2)}:${numbers.slice(2)}`;
    } else if (numbers.length === 2) {
      return `${numbers}:`;
    }
    return numbers;
  }
  
  function validateTime(timeValue) {
    if (!timeValue) return false;
    const timeRegex = /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/;
    return timeRegex.test(timeValue);
  }
  
  function formatTimeForMessage(timeValue) {
    if (!timeValue) return '';
    const [hours, minutes] = timeValue.split(':');
    const hour = parseInt(hours);
    const min = minutes;
    const period = hour >= 12 ? 'PM' : 'AM';
    const displayHour = hour === 0 ? 12 : hour > 12 ? hour - 12 : hour;
    return `${displayHour}:${min} ${period}`;
  }
  
  function handleQuantityChange(cartItem, event) {
    const quantity = parseInt(event.currentTarget.value) || 0;
    const itemKey = cartItem.acompanamientoId 
      ? `${cartItem.title}_${cartItem.acompanamientoId}` 
      : cartItem.title;
    cartStore.updateQuantity(itemKey, quantity);
  }
  
  function handleRemoveItem(cartItem) {
    const itemKey = cartItem.acompanamientoId 
      ? `${cartItem.title}_${cartItem.acompanamientoId}` 
      : cartItem.title;
    cartStore.removeItemByKey(itemKey);
  }
  
  function handleSendOrderClick() {
    showOrderForm = true;
  }
  
  function handleCancelOrder() {
    showOrderForm = false;
    nombreRetiro = '';
    horaRetiro = '';
  }
  
  function handleTimeInput(event) {
    const value = event.currentTarget.value;
    horaRetiro = formatTimeInput(value);
  }
  
  function handleConfirmOrder() {
    if (!nombreRetiro.trim() || !horaRetiro.trim()) {
      alert('Por favor completa todos los campos');
      return;
    }
    
    if (!validateTime(horaRetiro)) {
      alert('Por favor ingresa una hora válida (formato: HH:MM, ejemplo: 14:30)');
      return;
    }
    
    const horaFormateada = formatTimeForMessage(horaRetiro);
    const horaConZona = `${horaFormateada} (${timeZoneName()})`;
    
    const url = cartStore.generateWhatsAppMessage(
      restaurantData.businessInfo.whatsapp,
      nombreRetiro.trim(),
      horaConZona
    );
    window.open(url, '_blank');
    showOrderForm = false;
    nombreRetiro = '';
    horaRetiro = '';
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
  <div class="fixed bottom-0 left-0 right-0 z-50 bg-[#E8E4D9] shadow-2xl border-t-2 border-gray-300">
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
                Tu pedido
              </p>
              <p class="text-xs sm:text-sm text-gray-600">
                {totalItems} {totalItems === 1 ? 'item' : 'items'}
              </p>
            </div>
          </button>
          
          <div class="text-right">
            <p class="text-xs sm:text-sm text-gray-600 mb-1">Total</p>
            <Price price={total} className="text-lg sm:text-xl lg:text-2xl font-bold" />
          </div>
        </div>
        
        <!-- Botones de acción -->
        <div class="flex items-center gap-2 sm:gap-3">
          <button
            onclick={handleClearCart}
            class="px-3 py-2 sm:px-4 sm:py-2.5 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors flex items-center gap-2"
            aria-label="Limpiar pedido"
          >
            <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            <span class="hidden sm:inline text-sm font-medium">Limpiar</span>
          </button>
          
          <button
            onclick={handleSendOrderClick}
            class="px-4 py-2 sm:px-6 sm:py-3 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors flex items-center gap-2 sm:gap-3 font-semibold text-sm sm:text-base"
            aria-label="Enviar pedido por WhatsApp"
          >
            <WhatsAppIcon className="w-5 h-5 sm:w-6 sm:h-6" />
            <span>Pedir por WhatsApp</span>
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
              <div class="bg-[#E8E4D9] rounded-lg p-4 sm:p-5">
                <div class="flex justify-between items-start gap-4 mb-3">
                  <div class="flex-1">
                    <h3 class="text-base sm:text-lg font-bold text-gray-800 mb-1">
                      {cartItem.title}
                      {#if cartItem.acompanamiento}
                        <span class="text-sm font-normal text-gray-600">({cartItem.acompanamiento})</span>
                      {/if}
                    </h3>
                    {#if cartItem.description}
                      <p class="text-xs sm:text-sm text-gray-600">
                        {cartItem.description}
                      </p>
                    {/if}
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
                    <label for="quantity-{cartItem.title}" class="text-xs sm:text-sm text-gray-700 font-medium">
                      Cantidad:
                    </label>
                    <input
                      id="quantity-{cartItem.title}"
                      type="number"
                      min="1"
                      value={cartItem.cantidad}
                      oninput={(e) => handleQuantityChange(cartItem, e)}
                      class="w-14 sm:w-16 px-2 py-1 border border-gray-300 rounded text-center text-xs sm:text-sm"
                    />
                  </div>
                  <div class="text-right">
                    <p class="text-sm sm:text-base text-gray-800 font-semibold">
                      ${(cartItem.precio * cartItem.cantidad).toLocaleString('es-CL')}
                    </p>
                    {#if cartItem.cantidad > 1}
                      <p class="text-xs text-gray-500">
                        ${cartItem.precio.toLocaleString('es-CL')} c/u
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

<!-- Modal de información de retiro -->
{#if showOrderForm}
  <div 
    class="fixed inset-0 bg-black/50 z-[60] flex items-center justify-center p-4" 
    role="dialog"
    aria-modal="true"
    aria-labelledby="order-form-title"
    tabindex="-1"
    onclick={handleCancelOrder}
    onkeydown={(e) => {
      if (e.key === 'Escape') {
        handleCancelOrder();
      }
    }}
  >
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6 sm:p-8">
      <h3 id="order-form-title" class="text-xl sm:text-2xl font-bold text-gray-800 mb-4 sm:mb-6">
        Información de Retiro
      </h3>
      
      <div class="space-y-4 sm:space-y-5" onclick={(e) => e.stopPropagation()}>
        <div>
          <label for="nombre-retiro" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
            Nombre de quien va a retirar *
          </label>
          <input
            id="nombre-retiro"
            type="text"
            bind:value={nombreRetiro}
            placeholder="Ej: Juan Pérez"
            class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
            required
          />
        </div>
        
        <div>
          <label for="hora-retiro" class="block text-sm sm:text-base font-medium text-gray-700 mb-2">
            Hora de retiro * <span class="text-xs text-gray-500 font-normal">({timeZoneName()})</span>
          </label>
          <input
            id="hora-retiro"
            type="text"
            value={horaRetiro}
            oninput={handleTimeInput}
            placeholder="Ej: 14:30"
            maxlength="5"
            class="w-full px-4 py-2 sm:py-3 border border-gray-300 rounded-lg text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-green-500"
            required
          />
          <p class="text-xs text-gray-500 mt-1">Formato: HH:MM (24 horas, ejemplo: 14:30 para 2:30 PM)</p>
        </div>
      </div>
      
      <div class="flex gap-3 sm:gap-4 mt-6 sm:mt-8" onclick={(e) => e.stopPropagation()}>
        <button
          onclick={handleCancelOrder}
          class="flex-1 px-4 py-2 sm:py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm sm:text-base"
        >
          Cancelar
        </button>
        <button
          onclick={handleConfirmOrder}
          class="flex-1 px-4 py-2 sm:py-3 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors font-semibold text-sm sm:text-base flex items-center justify-center gap-2"
        >
          <WhatsAppIcon className="w-5 h-5" />
          <span>Enviar Pedido</span>
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal de confirmación para limpiar carrito -->
{#if showClearConfirm}
  <div 
    class="fixed inset-0 bg-black/50 z-[60] flex items-center justify-center p-4" 
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
        ¿Vaciar carrito?
      </h3>
      
      <p class="text-gray-600 mb-6">
        ¿Estás seguro de que quieres vaciar el carrito? Esta acción no se puede deshacer.
      </p>
      
      <div class="flex gap-3 sm:gap-4">
        <button
          onclick={cancelClearCart}
          class="flex-1 px-4 py-2 sm:py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm sm:text-base"
        >
          Cancelar
        </button>
        <button
          onclick={confirmClearCart}
          class="flex-1 px-4 py-2 sm:py-3 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors font-semibold text-sm sm:text-base"
        >
          Vaciar Carrito
        </button>
      </div>
    </div>
  </div>
{/if}

