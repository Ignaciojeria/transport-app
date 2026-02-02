<script>
  import { cartStore } from '../../stores/cartStore.svelte.js';
  import { getPriceFromPricing } from '../../services/menuData.js';
  import { t } from '../../lib/useLanguage';
  
  const { item } = $props();
  
  let itemImageError = $state(false);
  
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
  
  const isInCart = $derived.by(() => {
    const items = cartStore.items;
    const matchingItems = items.filter(i => i.title === item.title);
    return matchingItems.reduce((sum, i) => sum + (i.customQuantity || i.cantidad), 0) > 0;
  });
  
  function handleAddToCart(event) {
    event.stopPropagation();
    if (hasAcompanamientos) {
      // TODO: Abrir modal de acompañamientos
      return;
    }
    if (needsQuantitySelector) {
      // TODO: Abrir modal de cantidad
      return;
    }
    
    cartStore.addItem({
      title: item.title,
      precio: displayPrice,
      cantidad: 1,
      photoUrl: item.photoUrl,
      pricing: item.pricing
    });
  }
  
  // Mapeo de alérgenos comunes con traducciones
  const allergenMap = $derived(() => ({
    'gluten': { icon: '🌾', label: $t.menu.allergens.gluten, key: 'gluten' },
    'glúten': { icon: '🌾', label: $t.menu.allergens.gluten, key: 'gluten' },
    'gluten': { icon: '🌾', label: $t.menu.allergens.gluten, key: 'gluten' },
    'ovo': { icon: '🥚', label: $t.menu.allergens.egg, key: 'egg' },
    'egg': { icon: '🥚', label: $t.menu.allergens.egg, key: 'egg' },
    'huevo': { icon: '🥚', label: $t.menu.allergens.egg, key: 'egg' },
    'lactose': { icon: '🧀', label: $t.menu.allergens.lactose, key: 'lactose' },
    'lácteos': { icon: '🧀', label: $t.menu.allergens.lactose, key: 'lactose' },
    'lactosa': { icon: '🧀', label: $t.menu.allergens.lactose, key: 'lactose' },
    'frutos do mar': { icon: '🦐', label: $t.menu.allergens.seafood, key: 'seafood' },
    'seafood': { icon: '🦐', label: $t.menu.allergens.seafood, key: 'seafood' },
    'mariscos': { icon: '🦐', label: $t.menu.allergens.seafood, key: 'seafood' }
  }));
  
  // Extraer información de alérgenos de la descripción o de un campo específico
  const allergens = $derived(() => {
    const found = [];
    const foundKeys = new Set(); // Evitar duplicados
    const desc = (item.description || '').toLowerCase();
    const title = (item.title || '').toLowerCase();
    const text = `${desc} ${title}`;
    const map = allergenMap();
    
    for (const [key, value] of Object.entries(map)) {
      if (text.includes(key) && !foundKeys.has(value.key)) {
        found.push(value);
        foundKeys.add(value.key);
      }
    }
    return found;
  });
  
  // Badge opcional (ej: "MAIS VENDIDA", "ESPECIAL") - puede venir de item.badge o item.tags
  const badge = $derived(item.badge || item.tags?.find(t => ['mais vendida', 'especial', 'popular'].includes(t.toLowerCase())));
</script>

<div class="chef-card">
  <!-- Imagen del producto -->
  {#if item.photoUrl && !itemImageError}
    <img 
      src={item.photoUrl} 
      alt={item.title}
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
    <h3 class="chef-card-title">{item.title}</h3>
    
    <!-- Descripción -->
    {#if item.description}
      <p class="chef-card-description">{item.description}</p>
    {/if}
    
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
    
    <!-- Footer: Precio y Botón -->
    <div class="chef-card-footer">
      {#if displayPrice}
        <span class="price">CLP {displayPrice.toLocaleString('es-CL')}</span>
      {:else if hasAcompanamientos}
        <span class="price text-sm">{$t.menu.viewOptions}</span>
      {:else}
        <span class="price text-sm">—</span>
      {/if}
      <button 
        class="cta-btn" 
        onclick={handleAddToCart}
        disabled={isInCart}
      >
        {isInCart ? $t.menu.addedToCart : $t.menu.addToCart}
      </button>
    </div>
  </div>
</div>

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

  .cta-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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
