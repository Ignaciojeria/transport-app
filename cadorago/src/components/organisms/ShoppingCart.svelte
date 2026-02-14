<script>
  import { cartStore, itemsStore } from '../../stores/cartStore.svelte.js';
  import Price from '../atoms/Price.svelte';
  import WhatsAppIcon from '../atoms/WhatsAppIcon.svelte';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { t, language } from '../../lib/useLanguage';
  import { getEffectiveCurrency, formatPrice } from '../../lib/currency';
  import { getPriceFromPricing } from '../../services/menuData.js';
  
  const { className = '' } = $props();
  
  const restaurantData = $derived(restaurantDataStore.value);
  const currency = $derived(getEffectiveCurrency(restaurantData));
  const currentLanguage = $derived($language);
  
  // $effect + subscribe: reactividad explícita (Svelte 5)
  let items = $state([]);
  $effect(() => {
    const unsub = itemsStore.subscribe((v) => { items = v ?? []; });
    return unsub;
  });
  const total = $derived.by(() => {
    return items.reduce((sum, item) => {
      if (item.customQuantity != null && item.pricing) {
        const p = getPriceFromPricing(item.pricing, item.customQuantity);
        return sum + (p * (item.cantidad || 1));
      }
      return sum + ((item.precio || 0) * (item.cantidad || 1));
    }, 0);
  });
  const totalItems = $derived.by(() => {
    return items.reduce((sum, item) => {
      if (item.customQuantity) return sum + item.customQuantity;
      return sum + (item.cantidad || 0);
    }, 0);
  });
  
  let showClearConfirm = $state(false);
  
  function handleQuantityChange(title, event) {
    const quantity = parseInt(event.currentTarget.value) || 0;
    cartStore.updateQuantity(title, quantity);
  }
  
  function handleRemoveItem(title) {
    cartStore.removeItem(title);
  }
  
  function handleSendOrder() {
    const url = cartStore.generateWhatsAppMessage(
      restaurantData?.businessInfo?.whatsapp || '',
      '',
      '',
      currentLanguage,
      $t.whatsapp
    );
    window.open(url, '_blank');
  }
  
  function handleClearCart() {
    showClearConfirm = true;
  }
  
  function confirmClearCart() {
    cartStore.clear();
    showClearConfirm = false;
  }
  
  function cancelClearCart() {
    showClearConfirm = false;
  }
</script>

<div class={`bg-white/90 backdrop-blur-sm rounded-lg shadow-lg p-6 sm:p-8 lg:p-10 ${className}`}>
  <div class="flex justify-between items-center mb-6 sm:mb-8">
    <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800">
      {$t.cart.shoppingCart}
    </h2>
    {#if items.length > 0}
      <button
        onclick={handleClearCart}
        class="text-sm text-red-600 hover:text-red-800 font-medium"
        aria-label={$t.cart.clearOrder}
      >
        {$t.cart.clear}
      </button>
    {/if}
  </div>
  
  {#if items.length === 0}
    <div class="text-center py-8 sm:py-12">
      <p class="text-lg sm:text-xl text-gray-600">
        {$t.cart.emptyCart}
      </p>
      <p class="text-sm sm:text-base text-gray-500 mt-2">
        {$t.cart.emptyCartMessage}
      </p>
    </div>
  {:else}
    <div class="space-y-4 sm:space-y-5 lg:space-y-6 mb-6 sm:mb-8">
      {#each items as cartItem}
        <div class="bg-gray-50 rounded-lg p-4 sm:p-5 lg:p-6 border border-gray-200">
          <div class="flex justify-between items-start gap-4 mb-3">
            <div class="flex-1">
              <h3 class="text-lg sm:text-xl lg:text-2xl font-bold text-gray-800 mb-1">
                {cartItem.title}
              </h3>
              {#if cartItem.description}
                <p class="text-sm sm:text-base text-gray-600">
                  {cartItem.description}
                </p>
              {/if}
            </div>
            <button
              onclick={() => handleRemoveItem(cartItem.title)}
              class="text-red-600 hover:text-red-800 text-xl font-bold"
              aria-label="Eliminar item"
            >
              ×
            </button>
          </div>
          
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3 sm:gap-4">
              <label for="quantity-{cartItem.title}" class="text-sm sm:text-base text-gray-700 font-medium">
                {$t.cart.quantity}
              </label>
              <input
                id="quantity-{cartItem.title}"
                type="number"
                min="1"
                value={cartItem.cantidad}
                oninput={(e) => handleQuantityChange(cartItem.title, e)}
                class="w-16 sm:w-20 px-2 py-1 border border-gray-300 rounded text-center text-sm sm:text-base"
              />
            </div>
            <div class="text-right">
              <p class="text-sm sm:text-base text-gray-600 mb-1">
                {formatPrice(cartItem.precio * cartItem.cantidad, currency)}
              </p>
              {#if cartItem.cantidad > 1}
                <p class="text-xs sm:text-sm text-gray-500">
                  {formatPrice(cartItem.precio, currency)} c/u
                </p>
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>
    
    <div class="border-t border-gray-300 pt-6 sm:pt-8">
      <div class="flex justify-between items-center mb-6 sm:mb-8">
        <span class="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-800">
          {$t.cart.total}:
        </span>
        <Price price={total} className="text-xl sm:text-2xl lg:text-3xl" />
      </div>
      
      <button
        onclick={handleSendOrder}
        class="w-full flex items-center justify-center gap-3 px-6 py-4 sm:px-8 sm:py-5 lg:px-10 lg:py-6 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors font-semibold text-base sm:text-lg lg:text-xl"
        aria-label="Enviar pedido por WhatsApp"
      >
        <WhatsAppIcon className="w-6 h-6 sm:w-7 sm:h-7 lg:w-8 lg:h-8" />
        <span>{$t.cart.sendOrderWhatsApp}</span>
      </button>
    </div>
  {/if}
</div>

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

