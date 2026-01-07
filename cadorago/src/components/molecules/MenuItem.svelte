<script>
  import MenuItemTitle from '../atoms/MenuItemTitle.svelte';
  import MenuItemDescription from '../atoms/MenuItemDescription.svelte';
  import Price from '../atoms/Price.svelte';
  import { cartStore } from '../../stores/cartStore.svelte.js';
  
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
  let selectedAcompanamiento = $state(null);
  
  const hasAcompanamientos = $derived(item.sides && item.sides.length > 0);
  const displayPrice = $derived(hasAcompanamientos ? null : (item.price || 0));
  
  // Verificar si el item está en el carrito (necesita verificar por clave única)
  const isInCart = $derived.by(() => {
    const items = cartStore.items;
    const matchingItems = items.filter(i => i.title === item.title);
    return matchingItems.reduce((sum, i) => sum + i.cantidad, 0) > 0;
  });
  
  const cartItems = $derived.by(() => cartStore.items.filter(i => i.title === item.title));
  const totalQuantity = $derived.by(() => cartItems.reduce((sum, i) => sum + i.cantidad, 0));
  
  function handleAddToCart(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      showAcompanamientoModal = true;
    } else {
      cartStore.addItem(item);
    }
  }
  
  function handleIncrement(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      showAcompanamientoModal = true;
    } else {
      cartStore.addItem(item);
    }
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
      showAcompanamientoModal = false;
      selectedAcompanamiento = null;
    } catch (error) {
      console.error('Error al agregar item:', error);
      alert('Error al agregar el item: ' + error.message);
    }
  }
  
  function handleCancelAcompanamiento(event) {
    if (event) {
      event.stopPropagation();
    }
    showAcompanamientoModal = false;
    selectedAcompanamiento = null;
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
      <div class="flex items-center gap-2 mb-2">
        {#if isInCart}
          <svg class="w-5 h-5 text-gray-700 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        {/if}
        <MenuItemTitle title={item.title} className="text-base sm:text-lg font-bold text-gray-900" />
      </div>
      {#if item.description}
        <MenuItemDescription description={item.description} className="text-sm text-gray-700 mb-2" />
      {/if}
      {#if hasAcompanamientos}
        <p class="text-xs sm:text-sm text-gray-600 mb-2">
          Selecciona un acompañamiento
        </p>
      {/if}
      {#if displayPrice}
        <Price price={displayPrice} />
      {/if}
      {#if isInCart && displayPrice && totalQuantity > 0}
        <p class="text-sm text-gray-600 mt-1">
          Total: <Price price={(item.price || 0) * totalQuantity} className="font-bold" />
        </p>
      {/if}
    </div>
    <!-- Imagen vacía a la derecha (estilo Uber Eats) -->
    <div class="flex-shrink-0 relative">
      <div class="w-24 h-24 sm:w-28 sm:h-28 rounded-lg bg-gray-100 flex items-center justify-center overflow-hidden">
        <svg class="w-12 h-12 sm:w-14 sm:h-14 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
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
          <span class="w-6 text-center font-semibold text-sm text-gray-700">
            {totalQuantity}
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

<!-- Modal de selección de acompañamiento -->
{#if showAcompanamientoModal}
  <div 
    class="fixed inset-0 bg-black/50 z-[100] flex items-center justify-center p-4" 
    role="dialog"
    aria-modal="true"
    aria-labelledby="acompanamiento-modal-title"
    tabindex="-1"
    onclick={handleCancelAcompanamiento}
    onkeydown={(e) => {
      if (e.key === 'Escape') {
        handleCancelAcompanamiento();
      }
    }}
  >
    <div 
      class="bg-white rounded-lg shadow-xl max-w-md w-full p-6 sm:p-8 max-h-[80vh] overflow-y-auto"
      role="document"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
    >
      <h3 id="acompanamiento-modal-title" class="text-xl sm:text-2xl font-bold text-gray-800 mb-4 sm:mb-6">
        Selecciona un acompañamiento
      </h3>
      <p class="text-sm sm:text-base text-gray-600 mb-4 sm:mb-6">
        {item.title}
      </p>
      
      <div class="space-y-3 sm:space-y-4" onclick={(e) => e.stopPropagation()}>
        {#each item.sides as acompanamiento}
          <button
            type="button"
            onclick={(e) => handleSelectAcompanamiento(acompanamiento, e)}
            class="w-full text-left p-4 rounded-lg border-2 border-gray-300 hover:border-gray-800 hover:bg-gray-50 transition-colors"
          >
            <div class="flex justify-between items-center">
              <span class="font-medium text-gray-800 text-sm sm:text-base">
                {acompanamiento.name}
              </span>
              <Price price={acompanamiento.price} className="text-base sm:text-lg" />
            </div>
          </button>
        {/each}
      </div>
      
      <div class="flex gap-3 sm:gap-4 mt-6 sm:mt-8">
        <button
          type="button"
          onclick={handleCancelAcompanamiento}
          class="flex-1 px-4 py-2 sm:py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm sm:text-base"
        >
          Cancelar
        </button>
      </div>
    </div>
  </div>
{/if}
