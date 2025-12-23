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
  const cartItems = $derived(cartStore.items.filter(i => i.title === item.title));
  const totalQuantity = $derived(cartItems.reduce((sum, i) => sum + i.cantidad, 0));
  const isInCart = $derived(totalQuantity > 0);
  
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
  
  function handleItemClick() {
    if (!isInCart) {
      if (hasAcompanamientos) {
        showAcompanamientoModal = true;
      } else {
        cartStore.addItem(item);
      }
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
  onclick={handleItemClick}
  class={`bg-[#E8E4D9] rounded-lg p-6 sm:p-7 lg:p-8 cursor-pointer hover:bg-[#DDD9CE] transition-colors relative ${className}`}
  role="button"
  tabindex="0"
  onkeydown={(e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleItemClick();
    }
  }}
>
  <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4 sm:gap-6 lg:gap-8">
    <div class="flex-1 min-w-0">
      <MenuItemTitle title={item.title} />
      <MenuItemDescription description={item.description || ''} />
      {#if hasAcompanamientos}
        <p class="text-sm sm:text-base text-gray-500 mt-2 italic">
          Selecciona un acompañamiento
        </p>
      {/if}
    </div>
    <div class="flex items-center justify-between sm:justify-end gap-3 sm:gap-4 lg:gap-6 flex-shrink-0">
      {#if displayPrice}
        <div class="flex-shrink-0">
          <Price price={displayPrice} />
        </div>
      {/if}
      
      {#if isInCart}
        <!-- Controles de cantidad cuando está en el carrito -->
        <div class="flex items-center gap-1.5 sm:gap-2 lg:gap-3 flex-shrink-0">
          <button
            onclick={handleDecrement}
            class="w-8 h-8 sm:w-10 sm:h-10 lg:w-12 lg:h-12 xl:w-14 xl:h-14 bg-gray-800 hover:bg-gray-900 text-white rounded-full flex items-center justify-center transition-colors shadow-md hover:shadow-lg"
            aria-label="Disminuir cantidad"
          >
            <svg class="w-4 h-4 sm:w-5 sm:h-5 lg:w-6 lg:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
            </svg>
          </button>
          
          <span class="w-6 sm:w-8 lg:w-10 xl:w-12 text-center font-bold text-base sm:text-lg lg:text-xl xl:text-2xl text-gray-800">
            {totalQuantity}
          </span>
          
          <button
            onclick={handleIncrement}
            class="w-8 h-8 sm:w-10 sm:h-10 lg:w-12 lg:h-12 xl:w-14 xl:h-14 bg-gray-800 hover:bg-gray-900 text-white rounded-full flex items-center justify-center transition-colors shadow-md hover:shadow-lg"
            aria-label="Aumentar cantidad"
          >
            <svg class="w-4 h-4 sm:w-5 sm:h-5 lg:w-6 lg:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>
      {:else}
        <!-- Botón de agregar cuando no está en el carrito -->
        <button
          onclick={handleAddToCart}
          class="w-8 h-8 sm:w-10 sm:h-10 lg:w-12 lg:h-12 xl:w-14 xl:h-14 bg-gray-800 hover:bg-gray-900 text-white rounded-full flex items-center justify-center transition-colors flex-shrink-0 shadow-md hover:shadow-lg"
          aria-label="Agregar al carrito"
        >
          <svg class="w-4 h-4 sm:w-5 sm:h-5 lg:w-6 lg:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
