<script>
  import MenuItemTitle from '../atoms/MenuItemTitle.svelte';
  import MenuItemDescription from '../atoms/MenuItemDescription.svelte';
  import Price from '../atoms/Price.svelte';
  import QuantitySelector from './QuantitySelector.svelte';
  import { cartStore } from '../../stores/cartStore.svelte.js';
  import { getPriceFromPricing, getPricingLimits } from '../../services/menuData.js';
  
  const { 
    item = {
      title: '',
      description: '',
      price: 0,
      sides: null
    },
    className = ''
  } = $props();
  
  let showAcompanamientoModal = $state(false);
  let showQuantityModal = $state(false);
  let selectedAcompanamiento = $state(null);
  let acompanamientoViewTransition = $state(false); // Controla la transición de la vista
  let itemImageError = $state(false); // Controla si la imagen del item falló al cargar
  let sideImageErrors = $state({}); // Controla errores de carga de imágenes de sides
  
  const hasAcompanamientos = $derived(item.sides && item.sides.length > 0);

  // Mostrar etiqueta Cocina/Barra solo cuando la URL tiene ?station=true (ej. preview desde consola)
  const showStation = $derived(typeof window !== 'undefined' && new URLSearchParams(window.location?.search || '').get('station') === 'true');
  
  // Detectar el modo de pricing
  const pricingMode = $derived(item.pricing?.mode || 'UNIT');
  const needsQuantitySelector = $derived(
    pricingMode === 'WEIGHT' || 
    pricingMode === 'VOLUME' || 
    pricingMode === 'LENGTH' || 
    pricingMode === 'AREA'
  );
  
  // Calcular el precio más económico de los acompañamientos
  const minSidePrice = $derived.by(() => {
    if (!hasAcompanamientos || !item.sides || item.sides.length === 0) {
      return null;
    }
    
    const prices = item.sides.map(side => {
      return side.price || (side.pricing?.pricePerUnit || 0);
    });
    
    return Math.min(...prices);
  });
  
  // El precio puede venir directamente (formato antiguo) o desde pricing (formato nuevo)
  const displayPrice = $derived(
    hasAcompanamientos 
      ? null 
      : (item.price || (item.pricing?.pricePerUnit || 0))
  );
  
  // Verificar si el item está en el carrito (necesita verificar por clave única)
  const isInCart = $derived.by(() => {
    const items = cartStore.items;
    const matchingItems = items.filter(i => i.title === item.title);
    return matchingItems.reduce((sum, i) => sum + (i.customQuantity || i.cantidad), 0) > 0;
  });
  
  const cartItems = $derived.by(() => cartStore.items.filter(i => i.title === item.title));
  const totalQuantity = $derived.by(() => {
    return cartItems.reduce((sum, i) => {
      // Para items con cantidad personalizada, sumar la cantidad
      if (i.customQuantity) {
        return sum + i.customQuantity;
      }
      return sum + i.cantidad;
    }, 0);
  });
  
  // Calcular precio total para items en el carrito
  const totalPrice = $derived.by(() => {
    return cartItems.reduce((sum, cartItem) => {
      if (cartItem.customQuantity && cartItem.pricing) {
        // Calcular precio según cantidad personalizada
        return sum + getPriceFromPricing(cartItem.pricing, cartItem.customQuantity);
      }
      // Precio normal
      return sum + (cartItem.precio * cartItem.cantidad);
    }, 0);
  });
  
  function handleAddToCart(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      showAcompanamientoModal = true;
      // Activar transición después de un pequeño delay para suavidad
      setTimeout(() => {
        acompanamientoViewTransition = true;
      }, 10);
    } else if (needsQuantitySelector) {
      showQuantityModal = true;
    } else {
      cartStore.addItem(item);
    }
  }
  
  function handleIncrement(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      showAcompanamientoModal = true;
      // Activar transición después de un pequeño delay para suavidad
      setTimeout(() => {
        acompanamientoViewTransition = true;
      }, 10);
    } else if (needsQuantitySelector) {
      showQuantityModal = true;
    } else {
      cartStore.addItem(item);
    }
  }
  
  function handleQuantityConfirm(quantity) {
    cartStore.addItemWithQuantity(item, quantity);
    showQuantityModal = false;
  }
  
  function handleQuantityCancel() {
    showQuantityModal = false;
  }
  
  function handleDecrement(event) {
    event.stopPropagation();
    // Si tiene múltiples variantes en el carrito, eliminar la última
    if (cartItems.length > 0) {
      const lastItem = cartItems[cartItems.length - 1];
      const itemKey = lastItem.acompanamientoId 
        ? `${lastItem.title}_${lastItem.acompanamientoId}` 
        : lastItem.title;
      
      if (lastItem.cantidad > 1) {
        cartStore.updateQuantity(itemKey, lastItem.cantidad - 1);
      } else {
        cartStore.removeItemByKey(itemKey);
      }
    }
  }
  
  function handleItemClick(event) {
    // Evitar que se active si se hace clic en los botones de cantidad
    if (event.target.closest('button')) {
      return;
    }
    
    // Si el item no está en el carrito, agregarlo
    if (!isInCart) {
      handleAddToCart(event);
    }
  }
  
  function handleSelectAcompanamiento(acompanamiento, event) {
    if (event) {
      event.preventDefault();
      event.stopPropagation();
    }
    
    // Agregar directamente al carrito al seleccionar el acompañamiento
    try {
      cartStore.addItem(item, acompanamiento);
      // Primero ocultar la vista con transición
      acompanamientoViewTransition = false;
      // Luego limpiar después de la animación
      setTimeout(() => {
        showAcompanamientoModal = false;
        selectedAcompanamiento = null;
      }, 300); // Duración de la transición
    } catch (error) {
      console.error('Error al agregar item:', error);
      alert('Error al agregar el item: ' + error.message);
    }
  }
  
  function handleCancelAcompanamiento(event) {
    if (event) {
      event.stopPropagation();
    }
    // Primero ocultar la vista con transición
    acompanamientoViewTransition = false;
    // Luego limpiar después de la animación
    setTimeout(() => {
      showAcompanamientoModal = false;
      selectedAcompanamiento = null;
    }, 300); // Duración de la transición
  }
</script>

<div 
  class={`rounded-lg p-4 sm:p-5 mb-4 sm:mb-5 transition-all cursor-pointer hover:shadow-md ${isInCart ? 'bg-gray-100 border-2 border-gray-400' : 'bg-gray-50 border border-gray-200 hover:border-gray-300'} ${className}`}
  onclick={handleItemClick}
  role="button"
  tabindex="0"
  onkeydown={(e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleItemClick(e);
    }
  }}
>
  <div class="flex items-start gap-4">
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2 mb-2 flex-wrap">
        {#if isInCart}
          <svg class="w-5 h-5 text-gray-700 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        {/if}
        <MenuItemTitle title={item.title} className="text-base sm:text-lg font-bold text-gray-900" />
        {#if showStation && item.station === 'KITCHEN'}
          <span class="text-xs font-medium px-2 py-0.5 rounded bg-amber-100 text-amber-800" title="Cocina">Cocina</span>
        {:else if showStation && item.station === 'BAR'}
          <span class="text-xs font-medium px-2 py-0.5 rounded bg-blue-100 text-blue-800" title="Barra">Barra</span>
        {/if}
      </div>
      {#if item.description}
        <MenuItemDescription description={item.description} className="text-sm text-gray-700 mb-2" />
      {/if}
      {#if hasAcompanamientos}
        <p class="text-xs sm:text-sm text-gray-600 mb-2">
          Selecciona una opción
        </p>
        {#if minSidePrice !== null}
          <p class="text-sm sm:text-base text-gray-800 font-semibold">
            Desde <Price price={minSidePrice} />
          </p>
        {/if}
      {/if}
      {#if displayPrice && !needsQuantitySelector}
        <Price price={displayPrice} />
      {:else if needsQuantitySelector}
        <p class="text-xs sm:text-sm text-gray-600 mb-2">
          Precio por {item.pricing?.unit === 'GRAM' ? 'gramo' : item.pricing?.unit === 'KILOGRAM' ? 'kg' : item.pricing?.unit?.toLowerCase() || 'unidad'}
        </p>
        {#if item.pricing}
          <Price price={item.pricing.pricePerUnit} />
          <span class="text-xs text-gray-500"> / {item.pricing.baseUnit || 1} {item.pricing.unit === 'GRAM' ? 'g' : item.pricing.unit === 'KILOGRAM' ? 'kg' : item.pricing.unit?.toLowerCase() || ''}</span>
        {/if}
      {/if}
      {#if isInCart && totalQuantity > 0}
        <p class="text-sm text-gray-600 mt-1">
          {#if needsQuantitySelector}
            <span class="text-xs text-gray-500">{totalQuantity} {item.pricing?.unit === 'GRAM' ? 'g' : item.pricing?.unit === 'KILOGRAM' ? 'kg' : item.pricing?.unit?.toLowerCase() || ''} - </span>
          {/if}
          Total: <Price price={totalPrice} className="font-bold" />
        </p>
      {/if}
    </div>
    <!-- Imagen del item a la derecha (estilo Uber Eats) -->
    <div class="flex-shrink-0 relative">
      <div class="w-24 h-24 sm:w-28 sm:h-28 rounded-lg bg-gray-100 flex items-center justify-center overflow-hidden">
        {#if item.photoUrl && !itemImageError}
          <img 
            src={item.photoUrl} 
            alt={item.title}
            class="w-full h-full object-cover"
            onerror={() => {
              itemImageError = true
            }}
          />
        {:else}
          <svg class="w-12 h-12 sm:w-14 sm:h-14 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        {/if}
      </div>
      <!-- Controles de cantidad sobre la imagen (esquina inferior derecha) -->
      {#if isInCart}
        <div class="absolute -bottom-2 -right-2 flex items-center gap-1 bg-white rounded-full shadow-lg border border-gray-200 p-1">
          <button
            onclick={handleDecrement}
            class="w-7 h-7 rounded-full bg-gray-100 hover:bg-gray-200 flex items-center justify-center transition-colors"
            aria-label="Disminuir cantidad"
          >
            <svg class="w-4 h-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
            </svg>
          </button>
          <span class="w-6 text-center font-semibold text-sm text-gray-700" title={needsQuantitySelector ? `${totalQuantity} ${item.pricing?.unit === 'GRAM' ? 'g' : item.pricing?.unit === 'KILOGRAM' ? 'kg' : ''}` : ''}>
            {needsQuantitySelector ? '✏️' : totalQuantity}
          </span>
          <button
            onclick={handleIncrement}
            class="w-7 h-7 rounded-full bg-gray-100 hover:bg-gray-200 flex items-center justify-center transition-colors"
            aria-label="Aumentar cantidad"
          >
            <svg class="w-4 h-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>
      {:else}
        <button
          onclick={handleAddToCart}
          class="absolute -bottom-2 -right-2 w-8 h-8 bg-white hover:bg-gray-50 rounded-full shadow-lg border border-gray-300 flex items-center justify-center transition-colors"
          aria-label="Agregar al carrito"
        >
          <svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </button>
      {/if}
    </div>
  </div>
</div>

<!-- Modal de selección de cantidad (WEIGHT/VOLUME/etc) -->
{#if showQuantityModal}
  <div 
    class="fixed inset-0 bg-black/50 z-[100] flex items-center justify-center p-4" 
    role="dialog"
    aria-modal="true"
    aria-labelledby="quantity-modal-title"
    tabindex="-1"
    onclick={handleQuantityCancel}
    onkeydown={(e) => {
      if (e.key === 'Escape') {
        handleQuantityCancel();
      }
    }}
  >
    <div 
      class="bg-white rounded-lg shadow-xl max-w-md w-full p-6 sm:p-8 max-h-[80vh] overflow-y-auto"
      role="document"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
    >
      <h3 id="quantity-modal-title" class="text-xl sm:text-2xl font-bold text-gray-800 mb-2">
        {item.title}
      </h3>
      {#if item.description}
        <p class="text-sm text-gray-600 mb-4 sm:mb-6">
          {item.description}
        </p>
      {/if}
      
      {#if item.pricing}
        {@const limits = getPricingLimits(item.pricing)}
        {@const unitLabel = item.pricing.unit === 'GRAM' ? 'g' : item.pricing.unit === 'KILOGRAM' ? 'kg' : item.pricing.unit === 'MILLILITER' ? 'ml' : item.pricing.unit === 'LITER' ? 'L' : item.pricing.unit === 'METER' ? 'm' : item.pricing.unit === 'SQUARE_METER' ? 'm²' : ''}
        <QuantitySelector
          pricing={item.pricing}
          min={limits.min}
          max={limits.max}
          step={limits.step}
          value={limits.min}
          unit={unitLabel}
          onConfirm={handleQuantityConfirm}
          onCancel={handleQuantityCancel}
        />
      {/if}
    </div>
  </div>
{/if}

<!-- Vista de selección de acompañamiento con view transition -->
{#if showAcompanamientoModal}
  <div 
    class="fixed inset-0 bg-white z-[100] transition-transform duration-300 ease-in-out {acompanamientoViewTransition ? 'translate-x-0' : 'translate-x-full'}"
    role="dialog"
    aria-modal="true"
    aria-labelledby="acompanamiento-modal-title"
    tabindex="-1"
    onkeydown={(e) => {
      if (e.key === 'Escape') {
        handleCancelAcompanamiento();
      }
    }}
  >
    <div class="h-full w-full overflow-y-auto">
      <div class="max-w-2xl mx-auto p-4 sm:p-6 md:p-8">
        <!-- Header con botón de cerrar -->
        <div class="flex items-center justify-between mb-6 sm:mb-8 sticky top-0 bg-white z-10 pb-4 border-b border-gray-200 -mx-4 sm:-mx-6 md:-mx-8 px-4 sm:px-6 md:px-8">
          <button
            onclick={handleCancelAcompanamiento}
            class="p-2 text-gray-500 hover:text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
            aria-label="Cerrar"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h3 id="acompanamiento-modal-title" class="text-xl sm:text-2xl font-bold text-gray-800 flex-1 text-center">
            Selecciona una opción
          </h3>
          <div class="w-10"></div> <!-- Spacer para centrar el título -->
        </div>
        
        <!-- Información del item -->
        <div class="mb-6">
          <h4 class="text-lg sm:text-xl font-semibold text-gray-800 mb-2">
            {item.title}
          </h4>
          {#if item.description}
            <p class="text-sm sm:text-base text-gray-600">
              {item.description}
            </p>
          {/if}
        </div>
        
        <!-- Lista de acompañamientos -->
        <div class="space-y-3 sm:space-y-4">
          {#each item.sides as acompanamiento}
            <button
              type="button"
              onclick={(e) => handleSelectAcompanamiento(acompanamiento, e)}
              class="w-full text-left p-4 sm:p-5 rounded-lg border-2 border-gray-300 hover:border-green-500 hover:bg-green-50 transition-all duration-200 active:scale-[0.98]"
            >
              <div class="flex items-center gap-4">
                <!-- Imagen del acompañamiento (estilo Uber Eats) -->
                <div class="flex-shrink-0">
                  <div class="w-20 h-20 sm:w-24 sm:h-24 rounded-lg bg-gray-100 flex items-center justify-center overflow-hidden">
                    {#if acompanamiento.photoUrl && !sideImageErrors[acompanamiento.name]}
                      <img 
                        src={acompanamiento.photoUrl} 
                        alt={acompanamiento.name}
                        class="w-full h-full object-cover"
                        onerror={() => {
                          sideImageErrors[acompanamiento.name] = true
                        }}
                      />
                    {:else}
                      <svg class="w-10 h-10 sm:w-12 sm:h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                    {/if}
                  </div>
                </div>
                <!-- Información del acompañamiento -->
                <div class="flex-1 flex justify-between items-center gap-2 flex-wrap">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-medium text-gray-800 text-base sm:text-lg">
                      {acompanamiento.name}
                    </span>
                    {#if showStation && acompanamiento.station === 'KITCHEN'}
                      <span class="text-xs font-medium px-2 py-0.5 rounded bg-amber-100 text-amber-800">Cocina</span>
                    {:else if showStation && acompanamiento.station === 'BAR'}
                      <span class="text-xs font-medium px-2 py-0.5 rounded bg-blue-100 text-blue-800">Barra</span>
                    {/if}
                  </div>
                  <Price price={acompanamiento.price} className="text-lg sm:text-xl font-semibold" />
                </div>
              </div>
            </button>
          {/each}
        </div>
      </div>
    </div>
  </div>
{/if}
