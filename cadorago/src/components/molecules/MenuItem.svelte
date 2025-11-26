<script>
  import MenuItemTitle from '../atoms/MenuItemTitle.svelte';
  import MenuItemDescription from '../atoms/MenuItemDescription.svelte';
  import Price from '../atoms/Price.svelte';
  import { cartStore } from '../../stores/cartStore.svelte.js';
  
  const { 
    item = {
      titulo: '',
      descripción: '',
      precio: 0
    },
    className = ''
  } = $props();
  
  // Verificar si el item está en el carrito
  const cartItem = $derived(cartStore.items.find(i => i.titulo === item.titulo));
  const quantity = $derived(cartItem?.cantidad || 0);
  const isInCart = $derived(quantity > 0);
  
  function handleAddToCart(event) {
    event.stopPropagation();
    cartStore.addItem(item);
  }
  
  function handleIncrement(event) {
    event.stopPropagation();
    cartStore.addItem(item);
  }
  
  function handleDecrement(event) {
    event.stopPropagation();
    if (quantity > 1) {
      cartStore.updateQuantity(item.titulo, quantity - 1);
    } else {
      cartStore.removeItem(item.titulo);
    }
  }
  
  function handleItemClick() {
    if (!isInCart) {
      cartStore.addItem(item);
    }
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
      <MenuItemTitle title={item.titulo} />
      <MenuItemDescription description={item.descripción} />
    </div>
    <div class="flex items-center justify-between sm:justify-end gap-3 sm:gap-4 lg:gap-6 flex-shrink-0">
      <div class="flex-shrink-0">
        <Price price={item.precio} />
      </div>
      
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
            {quantity}
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

