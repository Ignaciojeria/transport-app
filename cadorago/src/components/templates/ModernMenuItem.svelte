<script>
  import { cartStore, itemsStore } from '../../stores/cartStore.svelte.js';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { getPriceFromPricing } from '../../services/menuData.js';
  import { getMultilingualText, getBaseText } from '../../lib/multilingual';
  import { getEffectiveCurrency, formatPrice } from '../../lib/currency';
  import MenuItemDescription from '../atoms/MenuItemDescription.svelte';
  import { t } from '../../lib/useLanguage';
  
  const { item } = $props();

  const currency = $derived(getEffectiveCurrency(restaurantDataStore.value));
  
  let itemImageError = $state(false);
  let showAcompanamientoModal = $state(false);
  let acompanamientoViewTransition = $state(false);
  let sideImageErrors = $state({});
  
  // Obtener textos multiidioma
  const itemTitle = $derived(getMultilingualText(item.title));
  const itemTitleBase = $derived(getBaseText(item.title));
  
  const hasAcompanamientos = $derived(item.sides && item.sides.length > 0);
  const pricingMode = $derived(item.pricing?.mode || 'UNIT');
  const needsQuantitySelector = $derived(
    pricingMode === 'WEIGHT' || 
    pricingMode === 'VOLUME' || 
    pricingMode === 'LENGTH' || 
    pricingMode === 'AREA'
  );
  
  const displayPrice = $derived(
    hasAcompanamientos 
      ? null 
      : (item.price || (item.pricing?.pricePerUnit || 0))
  );

  // Precio más bajo entre los sides (para mostrar "Desde CLP X" como en template hero)
  const minSidePrice = $derived.by(() => {
    if (!hasAcompanamientos || !item.sides || item.sides.length === 0) {
      return null;
    }
    const prices = item.sides.map(side => {
      return side.price ?? getPriceFromPricing(side.pricing);
    });
    return Math.min(...prices);
  });
  
  // Sincronizar items del carrito para reactividad
  let cartItemsList = $state([]);
  $effect(() => {
    const unsub = itemsStore.subscribe((v) => { cartItemsList = v ?? []; });
    return unsub;
  });

  const isInCart = $derived.by(() => {
    const matchingItems = cartItemsList.filter(i => {
      const iTitle = typeof i.title === 'string' ? i.title : i.title?.base || '';
      return iTitle === itemTitleBase;
    });
    return matchingItems.reduce((sum, i) => sum + (i.customQuantity || i.cantidad), 0) > 0;
  });

  const cartItems = $derived.by(() =>
    cartItemsList.filter(i => {
      const iTitle = typeof i.title === 'string' ? i.title : i.title?.base || '';
      return iTitle === itemTitleBase;
    })
  );
  const totalQuantity = $derived.by(() =>
    cartItems.reduce((sum, i) => sum + (i.customQuantity || i.cantidad), 0)
  );

  function handleAddToCart(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      showAcompanamientoModal = true;
      setTimeout(() => {
        acompanamientoViewTransition = true;
      }, 10);
      return;
    }
    if (needsQuantitySelector) {
      // TODO: Abrir modal de cantidad
      return;
    }
    cartStore.addItem(item);
  }

  function handleIncrement(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      showAcompanamientoModal = true;
      setTimeout(() => {
        acompanamientoViewTransition = true;
      }, 10);
      return;
    }
    if (needsQuantitySelector) {
      handleAddToCart(event);
      return;
    }
    cartStore.addItem(item);
  }

  function handleSelectAcompanamiento(acompanamiento, event) {
    if (event) {
      event.preventDefault();
      event.stopPropagation();
    }
    try {
      cartStore.addItem(item, acompanamiento);
      acompanamientoViewTransition = false;
      setTimeout(() => {
        showAcompanamientoModal = false;
      }, 300);
    } catch (err) {
      console.error('Error al agregar item:', err);
    }
  }

  function handleCancelAcompanamiento(event) {
    if (event) event.stopPropagation();
    acompanamientoViewTransition = false;
    setTimeout(() => {
      showAcompanamientoModal = false;
    }, 300);
  }

  function handleDecrement(event) {
    event.stopPropagation();
    if (cartItems.length === 0) return;
    const lastItem = cartItems[cartItems.length - 1];
    const lastTitleBase = typeof lastItem.title === 'string' ? lastItem.title : lastItem.title?.base || '';
    const itemKey = lastItem.acompanamientoId
      ? `${lastTitleBase}_${lastItem.acompanamientoId}`
      : lastTitleBase;
    if (lastItem.cantidad > 1) {
      cartStore.updateQuantity(itemKey, lastItem.cantidad - 1);
    } else {
      cartStore.removeItemByKey(itemKey);
    }
  }
  
  // Mapeo de foodAttributes con traducciones e iconos
  const foodAttributeMap = $derived(() => ({
    'GLUTEN': { icon: '🌾', label: $t.menu.allergens.gluten },
    'SEAFOOD': { icon: '🦐', label: $t.menu.allergens.seafood },
    'NUTS': { icon: '🥜', label: $t.menu.allergens.nuts },
    'DAIRY': { icon: '🧀', label: $t.menu.allergens.dairy },
    'EGGS': { icon: '🥚', label: $t.menu.allergens.egg },
    'SOY': { icon: '🫘', label: $t.menu.allergens.soy },
    'VEGAN': { icon: '🌱', label: $t.menu.allergens.vegan },
    'VEGETARIAN': { icon: '🥗', label: $t.menu.allergens.vegetarian },
    'SPICY': { icon: '🌶️', label: $t.menu.allergens.spicy },
    'ALCOHOL': { icon: '🍷', label: $t.menu.allergens.alcohol }
  }));
  
  // Obtener atributos alimentarios del campo foodAttributes
  const allergens = $derived(() => {
    if (!item.foodAttributes || !Array.isArray(item.foodAttributes)) {
      return [];
    }
    
    const map = foodAttributeMap();
    return item.foodAttributes
      .filter(attr => map[attr])
      .map(attr => ({
        icon: map[attr].icon,
        label: map[attr].label
      }));
  });
  
  // Badge opcional (ej: "MAIS VENDIDA", "ESPECIAL") - puede venir de item.badge o item.tags
  const badge = $derived(item.badge || item.tags?.find(t => ['mais vendida', 'especial', 'popular'].includes(t.toLowerCase())));
</script>

<div class="chef-card">
  <!-- Imagen del producto -->
  {#if item.photoUrl && !itemImageError}
    <img 
      src={item.photoUrl} 
      alt={itemTitle}
      class="chef-card-image"
      onerror={() => itemImageError = true}
    />
  {:else}
    <div class="chef-card-image bg-gray-200 flex items-center justify-center">
      <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
    </div>
  {/if}
  
  <div class="chef-card-content">
    <!-- Badge opcional -->
    {#if badge}
      <span class="chef-badge">{badge.toUpperCase()}</span>
    {/if}
    
    <!-- Título -->
    <h3 class="chef-card-title">{itemTitle}</h3>
    
    <!-- Descripción -->
    <MenuItemDescription description={item.description} className="chef-card-description" />
    
    <!-- Información de alérgenos -->
    {#if allergens().length > 0}
      <div class="allergen-info">
        {#each allergens() as allergen}
          <div class="allergen-badge">
            <span class="allergen-icon">{allergen.icon}</span>
            <span>{allergen.label}</span>
          </div>
        {/each}
      </div>
    {/if}
    
    <!-- Footer: Precio y controles de cantidad o botón Agregar -->
    <div class="chef-card-footer">
      {#if displayPrice}
        <span class="price">{formatPrice(displayPrice, currency)}</span>
      {:else if hasAcompanamientos && minSidePrice !== null}
        <span class="price text-sm">{$t.menu.fromPrice} {formatPrice(minSidePrice, currency)}</span>
      {:else if hasAcompanamientos}
        <span class="price text-sm">{$t.menu.viewOptions}</span>
      {:else}
        <span class="price text-sm">—</span>
      {/if}
      {#if isInCart}
        <div class="quantity-controls">
          <button
            type="button"
            class="qty-btn"
            onclick={handleDecrement}
            aria-label="Disminuir cantidad"
          >
            −
          </button>
          <span class="qty-value">{totalQuantity}</span>
          <button
            type="button"
            class="qty-btn"
            onclick={handleIncrement}
            aria-label="Aumentar cantidad"
          >
            +
          </button>
        </div>
      {:else}
        <button 
          class="cta-btn" 
          onclick={handleAddToCart}
        >
          {$t.menu.addToCart}
        </button>
      {/if}
    </div>
  </div>
</div>

<!-- Sheet de selección de agregado con transición deslizante desde la derecha -->
{#if showAcompanamientoModal}
  <div
    class="acompanamiento-sheet {acompanamientoViewTransition ? 'acompanamiento-sheet--visible' : ''}"
    role="dialog"
    aria-modal="true"
    aria-labelledby="acompanamiento-sheet-title"
    tabindex="-1"
    onkeydown={(e) => {
      if (e.key === 'Escape') handleCancelAcompanamiento(e);
    }}
  >
    <div class="acompanamiento-sheet-inner">
      <header class="acompanamiento-sheet-header">
        <button
          type="button"
          onclick={handleCancelAcompanamiento}
          class="acompanamiento-sheet-back"
          aria-label="Volver"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h3 id="acompanamiento-sheet-title" class="acompanamiento-sheet-title">
          {$t.menu.chooseOption}
        </h3>
        <div class="acompanamiento-sheet-spacer"></div>
      </header>

      <div class="acompanamiento-sheet-item-info">
        <h4 class="acompanamiento-sheet-item-title">{itemTitle}</h4>
        {#if itemDescription}
          <p class="acompanamiento-sheet-item-desc">{itemDescription}</p>
        {/if}
      </div>

      <div class="acompanamiento-sheet-options">
        {#each item.sides as side}
          <button
            type="button"
            onclick={(e) => handleSelectAcompanamiento(side, e)}
            class="acompanamiento-sheet-option"
          >
            <div class="acompanamiento-sheet-option-image">
              {#if side.photoUrl && !sideImageErrors[getBaseText(side.name)]}
                <img
                  src={side.photoUrl}
                  alt={getMultilingualText(side.name)}
                  onerror={() => {
                    sideImageErrors[getBaseText(side.name)] = true;
                  }}
                />
              {:else}
                <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              {/if}
            </div>
            <div class="acompanamiento-sheet-option-content">
              <span class="acompanamiento-sheet-option-name">{getMultilingualText(side.name)}</span>
              <span class="acompanamiento-sheet-option-price">
                {formatPrice(side.price ?? getPriceFromPricing(side.pricing), currency)}
              </span>
            </div>
          </button>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .chef-card {
    background: linear-gradient(135deg, var(--accent-yellow) 0%, #FFFBF0 100%);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: var(--shadow);
    transition: all 0.3s ease;
    border: 2px solid transparent;
    display: flex;
    flex-direction: column;
  }

  .chef-card:hover {
    transform: translateY(-8px);
    box-shadow: var(--shadow-hover);
    border-color: var(--primary);
  }

  .chef-card-image {
    width: 100%;
    height: 200px;
    object-fit: cover;
  }

  .chef-card-content {
    padding: 1.5rem;
    flex-grow: 1;
    display: flex;
    flex-direction: column;
  }

  .chef-badge {
    display: inline-block;
    background-color: var(--primary);
    color: var(--white);
    padding: 0.4rem 0.8rem;
    border-radius: 20px;
    font-size: 0.75rem;
    font-weight: 700;
    margin-bottom: 0.8rem;
    text-transform: uppercase;
    align-self: flex-start;
  }

  .chef-card-title {
    font-family: 'Playfair Display', serif;
    font-size: 1.5rem;
    font-weight: 800;
    margin-bottom: 0.8rem;
    color: var(--text-dark);
  }

  .chef-card-description {
    font-size: 0.95rem;
    color: var(--text-light);
    margin-bottom: 1rem;
    line-height: 1.5;
    flex-grow: 1;
  }

  .allergen-info {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
  }

  .allergen-badge {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.85rem;
    color: var(--text-dark);
    background-color: rgba(200, 90, 58, 0.1);
    padding: 0.4rem 0.8rem;
    border-radius: 6px;
  }

  .allergen-icon {
    font-size: 1.1rem;
  }

  .chef-card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 1rem;
    border-top: 1px solid rgba(200, 90, 58, 0.2);
  }

  .price {
    font-family: 'Playfair Display', serif;
    font-size: 1.5rem;
    font-weight: 800;
    color: var(--primary);
  }

  .cta-btn {
    background-color: var(--primary);
    color: var(--white);
    border: none;
    padding: 0.6rem 1.2rem;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.9rem;
    transition: all 0.3s ease;
    text-decoration: none;
    display: inline-block;
  }

  .cta-btn:hover:not(:disabled) {
    background-color: var(--primary-dark);
    transform: scale(1.05);
  }

  .quantity-controls {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    background-color: var(--white);
    border-radius: 6px;
    padding: 0.25rem;
    border: 1px solid var(--border-light);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  .qty-btn {
    width: 1.75rem;
    height: 1.75rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--secondary);
    border: none;
    border-radius: 4px;
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--text-dark);
    cursor: pointer;
    transition: background-color 0.2s;
    line-height: 1;
  }

  .qty-btn:hover {
    background-color: var(--primary);
    color: var(--white);
  }

  .qty-value {
    min-width: 1.5rem;
    text-align: center;
    font-weight: 700;
    font-size: 0.95rem;
    color: var(--text-dark);
  }

  /* Sheet de selección de agregado: transición deslizante desde la derecha */
  .acompanamiento-sheet {
    position: fixed;
    inset: 0;
    background-color: var(--white);
    z-index: 200;
    transform: translateX(100%);
    transition: transform 0.3s ease-in-out;
    overflow-y: auto;
  }

  .acompanamiento-sheet--visible {
    transform: translateX(0);
  }

  .acompanamiento-sheet-inner {
    max-width: 28rem;
    margin: 0 auto;
    padding: 1rem 1.25rem 2rem;
  }

  .acompanamiento-sheet-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding-bottom: 1rem;
    margin-bottom: 1rem;
    border-bottom: 1px solid var(--border-light);
    position: sticky;
    top: 0;
    background: var(--white);
    z-index: 1;
  }

  .acompanamiento-sheet-back {
    padding: 0.5rem;
    color: var(--text-light);
    background: none;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: color 0.2s, background-color 0.2s;
  }

  .acompanamiento-sheet-back:hover {
    color: var(--primary);
    background-color: var(--secondary);
  }

  .acompanamiento-sheet-title {
    font-family: 'Playfair Display', serif;
    font-size: 1.25rem;
    font-weight: 800;
    color: var(--text-dark);
    margin: 0;
  }

  .acompanamiento-sheet-spacer {
    width: 2.5rem;
  }

  .acompanamiento-sheet-item-info {
    margin-bottom: 1.25rem;
  }

  .acompanamiento-sheet-item-title {
    font-family: 'Playfair Display', serif;
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--text-dark);
    margin: 0 0 0.35rem 0;
  }

  .acompanamiento-sheet-item-desc {
    font-size: 0.9rem;
    color: var(--text-light);
    margin: 0;
    line-height: 1.4;
  }

  .acompanamiento-sheet-options {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .acompanamiento-sheet-option {
    display: flex;
    align-items: center;
    gap: 1rem;
    width: 100%;
    text-align: left;
    padding: 1rem;
    border: 2px solid var(--border-light);
    border-radius: 12px;
    background: var(--white);
    cursor: pointer;
    transition: border-color 0.2s, background-color 0.2s, transform 0.2s;
  }

  .acompanamiento-sheet-option:hover {
    border-color: var(--primary);
    background-color: rgba(200, 90, 58, 0.06);
  }

  .acompanamiento-sheet-option:active {
    transform: scale(0.98);
  }

  .acompanamiento-sheet-option-image {
    flex-shrink: 0;
    width: 4rem;
    height: 4rem;
    border-radius: 8px;
    background-color: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .acompanamiento-sheet-option-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .acompanamiento-sheet-option-content {
    flex: 1;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .acompanamiento-sheet-option-name {
    font-weight: 600;
    font-size: 1rem;
    color: var(--text-dark);
  }

  .acompanamiento-sheet-option-price {
    font-family: 'Playfair Display', serif;
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--primary);
  }

  @media (max-width: 768px) {
    .chef-card-title {
      font-size: 1.2rem;
    }

    .price {
      font-size: 1.2rem;
    }
  }
</style>
