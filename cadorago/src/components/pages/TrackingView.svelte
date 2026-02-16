<script>
  import { onMount } from 'svelte';
  import { getOrderByTrackingId } from '../../services/trackingService.js';
  import { fetchRestaurantDataById } from '../../services/restaurantData.js';
  import { trackingStore } from '../../stores/trackingStore.svelte.js';
  import { getEffectiveCurrency, formatPrice } from '../../lib/currency';
  import { initLanguage, t } from '../../lib/useLanguage';

  let trackingCode = $state('');
  let order = $state(null);
  let loading = $state(false);
  let error = $state(null);
  let businessWhatsapp = $state(null);
  let menuData = $state(null);
  const currency = $derived(getEffectiveCurrency(menuData));

  const trackingFromUrl = $derived(() => {
    if (typeof window === 'undefined') return null;
    const match = window.location.pathname.match(/^\/track\/([^/]+)/);
    return match ? match[1] : null;
  });

  const savedTrackings = $derived(Array.isArray($trackingStore) ? $trackingStore : []);
  const hasTrackingInUrl = $derived(!!trackingFromUrl());

  /** Pedidos activos: no entregados ni cancelados */
  const activeOrders = $derived(
    savedTrackings.filter((e) => {
      const id = typeof e === 'string' ? e : e?.id;
      const summary = orderSummaries[id];
      if (!summary) return true;
      return summary.statusLabel !== 'Entregado' && summary.statusLabel !== 'Cancelado';
    })
  );
  /** Pedidos recientes: entregados o cancelados */
  const recentOrders = $derived(
    savedTrackings.filter((e) => {
      const id = typeof e === 'string' ? e : e?.id;
      const summary = orderSummaries[id];
      if (!summary) return false;
      return summary.statusLabel === 'Entregado' || summary.statusLabel === 'Cancelado';
    })
  );

  /** menuId para "Volver al menú" en lista: del primer pedido que tenga */
  const listMenuId = $derived(
    [...activeOrders, ...recentOrders].map((e) => {
      const id = typeof e === 'string' ? e : e?.id;
      return orderSummaries[id]?.menuId;
    }).find(Boolean) ?? order?.menuId
  );

  /** Cache de resumen por tracking: { orderNumber, statusLabel } para la lista */
  let orderSummaries = $state({});
  $effect(() => {
    const list = savedTrackings;
    if (!list.length || hasTrackingInUrl) return;
    const ids = list.map((e) => (typeof e === 'string' ? e : e?.id)).filter(Boolean);
    Promise.all(
      ids.map(async (id) => {
        try {
          const o = await getOrderByTrackingId(id);
          const statuses = o?.items ? [...new Set(o.items.map((i) => i.status))] : [];
          let statusLabel = 'Pedido confirmado';
          if (statuses.includes('CANCELLED') && statuses.length === 1) statusLabel = 'Cancelado';
          else if (statuses.includes('DELIVERED')) statusLabel = 'Entregado';
          else if (statuses.includes('DISPATCHED')) statusLabel = 'En camino';
          else if (statuses.includes('READY')) statusLabel = o?.fulfillment === 'DELIVERY' ? 'En camino' : 'Listo para retirar';
          else if (statuses.includes('IN_PROGRESS')) statusLabel = 'En preparación';
          return { id, orderNumber: o?.orderNumber ?? id, statusLabel, menuId: o?.menuId };
        } catch {
          return { id, orderNumber: id, statusLabel: 'Consultar' };
        }
      })
    ).then((results) => {
      orderSummaries = Object.fromEntries(results.map((r) => [r.id, { orderNumber: r.orderNumber, statusLabel: r.statusLabel, menuId: r.menuId }]));
      results.forEach((r) => {
        trackingStore.updateTracking(r.id, { isDelivered: r.statusLabel === 'Entregado' });
      });
    });
  });

  onMount(() => {
    initLanguage();
    const code = trackingFromUrl();
    if (code) {
      trackingCode = decodeURIComponent(code);
      fetchOrder(code);
    }
  });

  async function fetchOrder(code) {
    if (!code?.trim()) return;
    loading = true;
    error = null;
    order = null;
    businessWhatsapp = null;
    try {
      order = await getOrderByTrackingId(code.trim().toUpperCase());
      trackingStore.addTracking(order.trackingId, order.createdAt);
      const statuses = order?.items ? [...new Set(order.items.map((i) => i.status))] : [];
      const isDelivered = statuses.includes('DELIVERED');
      trackingStore.updateTracking(order.trackingId, { isDelivered });
      if (order?.menuId) {
        try {
          const menu = await fetchRestaurantDataById(order.menuId);
          menuData = menu;
          businessWhatsapp = menu?.businessInfo?.whatsapp || null;
        } catch {}
      }
    } catch (e) {
      error = e.message || 'Error al consultar el pedido';
    } finally {
      loading = false;
    }
  }

  let shareCopied = $state(false);

  async function shareTracking() {
    if (!order?.trackingId || typeof window === 'undefined') return;
    const url = `${window.location.origin}/track/${encodeURIComponent(order.trackingId)}`;
    const title = `Pedido #${order.orderNumber} - MiCartaPro`;
    const text = `Seguimiento del pedido #${order.orderNumber}`;
    try {
      if (navigator.share) {
        await navigator.share({ title, text, url });
      } else {
        await navigator.clipboard.writeText(url);
        shareCopied = true;
        setTimeout(() => { shareCopied = false; }, 2000);
      }
    } catch (err) {
      if (err?.name !== 'AbortError') {
        try {
          await navigator.clipboard.writeText(url);
          shareCopied = true;
          setTimeout(() => { shareCopied = false; }, 2000);
        } catch {}
      }
    }
  }

  function openWhatsAppContact() {
    if (!businessWhatsapp || !order) return;
    const total = order.items?.reduce((s, i) => s + (i.totalPrice || 0), 0) ?? 0;
    const statusLabel = overallStatus === 'PENDING' ? 'Pedido confirmado' : overallStatus === 'IN_PROGRESS' ? 'En preparación' : overallStatus === 'READY' ? (order.fulfillment === 'DELIVERY' ? 'En camino' : 'Listo para retirar') : overallStatus === 'DISPATCHED' ? 'En camino' : overallStatus === 'DELIVERED' ? 'Entregado' : overallStatus === 'CANCELLED' ? 'Cancelado' : 'Consultar';
    const lines = [
      'Hola, tengo una consulta sobre mi pedido:',
      '',
      `📋 *Pedido #${order.orderNumber}*`,
      `🔖 Código: ${order.trackingId}`,
      `📦 Tipo: ${order.fulfillment === 'DELIVERY' ? 'Envío a domicilio' : 'Retiro en local'}`,
      `📊 Estado: ${statusLabel}`,
      '',
      '*Detalle:*',
      ...(order.items || []).map((i) => `• ${i.itemName} ×${i.quantity}${i.unit && i.unit !== 'EACH' ? ` ${i.unit}` : ''} — $${(i.totalPrice || 0).toLocaleString('es-CL')}`),
      '',
      `*Total: $${total.toLocaleString('es-CL')}*`,
      '',
      `Fecha: ${order.createdAt ? new Date(order.createdAt).toLocaleString('es-CL') : ''}`,
    ];
    const msg = lines.join('\n');
    const phone = businessWhatsapp.replace(/[^0-9]/g, '');
    window.open(`https://wa.me/${phone}?text=${encodeURIComponent(msg)}`, '_blank');
  }

  function handleSubmit(e) {
    e.preventDefault();
    if (trackingCode.trim()) fetchOrder(trackingCode.trim());
  }

  const overallStatus = $derived.by(() => {
    if (!order?.items?.length) return null;
    const statuses = [...new Set(order.items.map((i) => i.status))];
    if (statuses.includes('CANCELLED') && statuses.length === 1) return 'CANCELLED';
    if (statuses.includes('DELIVERED')) return 'DELIVERED';
    if (statuses.includes('DISPATCHED')) return 'DISPATCHED';
    if (statuses.includes('READY')) return 'READY';
    if (statuses.includes('IN_PROGRESS')) return 'IN_PROGRESS';
    return 'PENDING';
  });

  // Timeline: pasos según fulfillment
  const timelineSteps = $derived.by(() => {
    const isDelivery = order?.fulfillment === 'DELIVERY';
    const s = overallStatus;
    return [
      { key: 'received', label: 'Pedido recibido', status: 'done' },
      { key: 'preparing', label: 'En preparación', status: s === 'PENDING' || s === 'IN_PROGRESS' ? 'current' : 'done' },
      { key: 'ready', label: isDelivery ? 'En camino' : 'Listo para retirar', status: s === 'READY' || s === 'DISPATCHED' ? 'current' : s === 'DELIVERED' ? 'done' : 'pending' },
      { key: 'delivered', label: 'Entregado', status: s === 'DELIVERED' ? 'done' : 'pending' },
    ];
  });
</script>

<div class="min-h-screen bg-slate-50 font-sans py-8 pb-24 px-4 sm:px-6 lg:px-8">
  <div class="max-w-lg mx-auto">
    <!-- Back link -->
    {#if hasTrackingInUrl}
      {#if order?.menuId}
        <a href="/m/{order.menuId}" class="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800 mb-6 transition-colors">
          ← Volver al menú
        </a>
      {:else}
        <a href="/track" class="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800 mb-6 transition-colors">
          ← Pedidos activos
        </a>
      {/if}
    {:else if listMenuId}
      <a href="/m/{listMenuId}" class="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800 mb-6 transition-colors">
        ← Volver al menú
      </a>
    {/if}

    <!-- Vista: lista de pedidos (solo cuando NO hay tracking en URL) -->
    {#if !hasTrackingInUrl}
      <h1 class="text-2xl font-bold text-slate-900 mb-2">Seguimiento de pedido</h1>
      <p class="text-slate-600 text-sm mb-6">Tus pedidos en curso y recientes</p>

      <!-- 1. Pedidos activos (máxima prioridad) -->
      {#if activeOrders.length > 0}
        <div class="mb-8">
          <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">Pedidos activos</h2>
          <div class="space-y-2">
            {#each activeOrders as entry}
              {@const id = typeof entry === 'string' ? entry : entry.id}
              {@const summary = orderSummaries[id]}
              <div class="flex items-center gap-2 p-4 bg-white rounded-xl border-2 border-slate-200 hover:border-slate-300 transition-all shadow-sm">
                <a href="/track/{encodeURIComponent(id)}" class="flex-1 flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3 min-w-0">
                  {#if summary}
                    <span class="font-semibold text-slate-900">Pedido #{summary.orderNumber}</span>
                    <span class="text-slate-600 text-sm font-medium">{summary.statusLabel}</span>
                  {:else}
                    <span class="font-mono font-semibold text-slate-900 truncate">{id}</span>
                    <span class="text-slate-400 text-sm">→ Cargando...</span>
                  {/if}
                </a>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- 2. Pedidos recientes (entregados / cancelados) -->
      {#if recentOrders.length > 0}
        <div class="mb-8">
          <h2 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">Pedidos recientes</h2>
          <div class="space-y-2">
            {#each recentOrders as entry}
              {@const id = typeof entry === 'string' ? entry : entry.id}
              {@const summary = orderSummaries[id]}
              <div class="flex items-center gap-2 p-4 bg-slate-50/80 rounded-xl border border-slate-200/60 hover:border-slate-300/80 transition-all">
                <a href="/track/{encodeURIComponent(id)}" class="flex-1 flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3 min-w-0">
                  {#if summary}
                    <span class="font-medium text-slate-700">Pedido #{summary.orderNumber}</span>
                    <span class="text-slate-500 text-sm">
                      {summary.statusLabel === 'Entregado' ? 'Entregado ✓' : summary.statusLabel}
                    </span>
                  {:else}
                    <span class="font-mono font-medium text-slate-700 truncate">{id}</span>
                    <span class="text-slate-400 text-sm">→ Cargando...</span>
                  {/if}
                </a>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- 3. Consultar por código (opcional, secundario) -->
      <div class="pt-4 border-t border-slate-200">
        <h2 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-3">Consultar pedido por código</h2>
        <form onsubmit={handleSubmit} class="flex gap-2">
          <input type="text" bind:value={trackingCode} placeholder="Ej: ABC12345" maxlength="12" disabled={loading}
            class="flex-1 px-4 py-2.5 rounded-lg border border-slate-200 bg-white focus:ring-2 focus:ring-slate-300 focus:border-slate-400 outline-none font-mono text-sm text-slate-900 placeholder:text-slate-400" />
          <button type="submit" disabled={loading || !trackingCode.trim()}
            class="px-5 py-2.5 bg-slate-200 hover:bg-slate-300 disabled:bg-slate-100 disabled:cursor-not-allowed text-slate-700 font-medium rounded-lg transition-colors text-sm">
            {loading ? 'Buscando...' : 'Consultar'}
          </button>
        </form>
      </div>
    {/if}

    <!-- Loading -->
    {#if hasTrackingInUrl && loading}
      <div class="flex flex-col items-center justify-center py-20">
        <div class="w-10 h-10 border-2 border-slate-200 border-t-slate-900 rounded-full animate-spin mb-4"></div>
        <p class="text-slate-500 text-sm">Cargando pedido...</p>
      </div>
    {/if}

    <!-- Error -->
    {#if error}
      <div class="p-4 bg-red-50 border border-red-100 rounded-xl text-red-700 text-sm">{error}</div>
    {/if}

    <!-- Vista: detalle del pedido (con timeline) -->
    {#if order && !loading}
      <div class="space-y-6">
        <!-- Header: Pedido #7, Estado, Retiro -->
        <div>
          <h1 class="text-2xl font-bold text-slate-900">Pedido #{order.orderNumber}</h1>
          <p class="text-slate-600 mt-1">{order.fulfillment === 'DELIVERY' ? 'Envío a domicilio' : 'Retiro en local'}</p>
          <p class="text-base font-semibold mt-2 {overallStatus === 'PENDING' ? 'text-slate-500' : overallStatus === 'IN_PROGRESS' ? 'text-blue-600' : overallStatus === 'READY' || overallStatus === 'DISPATCHED' ? 'text-emerald-600' : overallStatus === 'DELIVERED' ? 'text-slate-800' : 'text-slate-600'}">
            {overallStatus === 'PENDING' ? 'Pedido confirmado' : overallStatus === 'IN_PROGRESS' ? 'En preparación' : overallStatus === 'READY' ? 'Listo para retirar' : overallStatus === 'DISPATCHED' ? 'En camino' : overallStatus === 'DELIVERED' ? 'Entregado' : overallStatus === 'CANCELLED' ? 'Cancelado' : ''}
          </p>
        </div>

        <!-- Timeline -->
        <div class="bg-white rounded-xl border border-slate-200/80 p-6 shadow-sm">
          <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-4">Progreso</h2>
          <div class="relative">
            {#each timelineSteps as step, i}
              {@const done = step.status === 'done'}
              {@const current = step.status === 'current'}
              <div class="flex items-center gap-3 {i < timelineSteps.length - 1 ? 'mb-1' : ''}">
                <div class="flex flex-col items-center flex-shrink-0">
                  <div class="w-6 h-6 rounded-full flex items-center justify-center
                    {done ? 'bg-slate-900 text-white' : current ? 'bg-blue-500 text-white ring-4 ring-blue-100' : 'bg-slate-200 text-slate-400'}">
                    {#if done}
                      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" /></svg>
                    {:else if current}
                      <span class="w-1.5 h-1.5 rounded-full bg-white"></span>
                    {/if}
                  </div>
                  {#if i < timelineSteps.length - 1}
                    <div class="w-0.5 h-6 mt-1 {done ? 'bg-slate-900' : 'bg-slate-200'}"></div>
                  {/if}
                </div>
                <p class="font-medium text-sm {done ? 'text-slate-900' : current ? 'text-slate-900' : 'text-slate-400'}">{step.label}</p>
              </div>
            {/each}
          </div>
        </div>

        <!-- Detalle -->
        <div class="bg-white rounded-xl border border-slate-200/80 p-6 shadow-sm">
          <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-4">Detalle</h2>
          <ul class="space-y-3">
            {#each order.items as item}
              <li class="flex justify-between items-center py-2 border-b border-slate-100 last:border-0">
                <span class="font-medium text-slate-900">{item.itemName}</span>
                <span class="text-slate-500 text-sm">×{item.quantity} — {formatPrice(item.totalPrice || 0, currency)}</span>
              </li>
            {/each}
          </ul>
          <div class="flex justify-between items-center pt-4 mt-4 border-t-2 border-slate-200">
            <span class="font-bold text-slate-900">Total</span>
            <span class="font-bold text-lg text-slate-900">{formatPrice(order.items?.reduce((s, i) => s + (i.totalPrice || 0), 0) ?? 0, currency)}</span>
          </div>
        </div>

        <!-- Compartir seguimiento (acción secundaria) -->
        <button onclick={shareTracking}
          class="w-full flex items-center justify-center gap-2 px-4 py-3.5 bg-slate-900 hover:bg-slate-800 text-white font-semibold rounded-xl transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
          {shareCopied ? $t.tracking.shareCopied : $t.tracking.shareTracking}
        </button>

        <!-- Contactar con la tienda (acción terciaria) -->
        {#if businessWhatsapp}
          <button onclick={openWhatsAppContact}
            class="w-full flex items-center justify-center gap-2 px-4 py-3.5 border-2 border-slate-300 hover:border-slate-400 hover:bg-slate-50 text-slate-700 font-medium rounded-xl transition-colors">
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
            </svg>
            {$t.tracking.contactStore}
          </button>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  .animate-spin { animation: spin 0.8s linear infinite; }
</style>
