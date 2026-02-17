<script>
  import { trackingStore } from '../../stores/trackingStore.svelte.js';
  import { itemsStore } from '../../stores/cartStore.svelte.js';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { getSlugFromUrl } from '../../services/restaurantData.js';
  import { t } from '../../lib/useLanguage';

  const currentMenuId = $derived(
    restaurantDataStore.value?.id || (getSlugFromUrl()?.isMenuId ? getSlugFromUrl()?.value : null)
  );
  const allTrackings = $derived(Array.isArray($trackingStore) ? $trackingStore : []);
  /** Filtrar por menuId del menú actual; si no hay menuId, mostrar todos (compatibilidad) */
  const trackingsForMenu = $derived(
    currentMenuId
      ? allTrackings.filter((e) => (typeof e === 'object' ? e?.menuId : null) === currentMenuId)
      : allTrackings
  );
  /** Solo pedidos activos (no entregados) del menú actual */
  const activeTrackings = $derived(trackingsForMenu.filter((e) => !(typeof e === 'object' && e?.isDelivered === true)));
  const hasCartItems = $derived.by(() => {
    const list = Array.isArray($itemsStore) ? $itemsStore : [];
    const forMenu = currentMenuId ? list.filter((i) => (i?.menuId ?? null) === currentMenuId) : list;
    return forMenu.length > 0;
  });
  const showButton = $derived(activeTrackings.length > 0 && !hasCartItems);

  const firstId = $derived(activeTrackings[0] ? (typeof activeTrackings[0] === 'string' ? activeTrackings[0] : activeTrackings[0]?.id) : '');
  const isSingle = $derived(activeTrackings.length === 1);
  const buttonText = $derived(isSingle ? $t.cart.viewYourOrder : $t.cart.viewYourOrders);
  const menuQuery = $derived(currentMenuId ? `?m=${encodeURIComponent(currentMenuId)}` : '');
  const buttonHref = $derived(isSingle && firstId ? `/track/${encodeURIComponent(firstId)}${menuQuery}` : `/track${menuQuery}`);
</script>

{#if showButton}
  <a
    href={buttonHref}
    class="floating-tracking-btn"
    aria-label={buttonText}
    title={isSingle ? firstId : `${activeTrackings.length} pedidos`}
  >
    <span class="floating-tracking-btn-icon" aria-hidden="true">🧾</span>
    <span class="floating-tracking-btn-text">{buttonText}</span>
  </a>
{/if}

<style>
  .floating-tracking-btn {
    position: fixed;
    bottom: 1.25rem;
    right: 1.25rem;
    z-index: 60;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    background: #171717;
    color: white;
    font-size: 0.9375rem;
    font-weight: 600;
    text-decoration: none;
    border-radius: 9999px;
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.2), 0 2px 6px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .floating-tracking-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25), 0 4px 10px rgba(0, 0, 0, 0.15);
    background: #262626;
    color: white;
  }

  .floating-tracking-btn:active {
    transform: translateY(0);
  }

  .floating-tracking-btn-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.25rem;
    line-height: 1;
  }

  .floating-tracking-btn-text {
    white-space: nowrap;
  }

  @media (max-width: 640px) {
    .floating-tracking-btn {
      bottom: 1rem;
      right: 1rem;
      padding: 0.625rem 1rem;
      font-size: 0.875rem;
    }
  }
</style>
