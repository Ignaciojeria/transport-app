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
  
  function handleAddToCart(event) {
    event.stopPropagation();
    cartStore.addItem(item);
  }
  
  function handleItemClick() {
    cartStore.addItem(item);
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
  <div class="flex justify-between items-start gap-6 sm:gap-8">
    <div class="flex-1">
      <MenuItemTitle title={item.titulo} />
      <MenuItemDescription description={item.descripción} />
    </div>
    <div class="flex items-center gap-4 sm:gap-6">
      <div class="flex-shrink-0">
        <Price price={item.precio} />
      </div>
      <button
        onclick={handleAddToCart}
        class="w-10 h-10 sm:w-12 sm:h-12 lg:w-14 lg:h-14 bg-gray-800 hover:bg-gray-900 text-white rounded-full flex items-center justify-center transition-colors flex-shrink-0 shadow-md hover:shadow-lg"
        aria-label="Agregar al carrito"
      >
        <svg class="w-5 h-5 sm:w-6 sm:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
      </button>
    </div>
  </div>
</div>

