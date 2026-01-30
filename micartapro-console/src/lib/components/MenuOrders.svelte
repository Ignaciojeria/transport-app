<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getMenuOrders, type MenuOrderRow } from '../menuUtils'
  import { t as tStore } from '../useLanguage'

  interface MenuOrdersProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: MenuOrdersProps = $props()

  let orders = $state<MenuOrderRow[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let menuId = $state<string | null>(null)
  let expandedOrderNumber = $state<number | null>(null)
  let paperOrder = $state<MenuOrderRow | null>(null)
  let thermalPrintMode = $state(false)

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const session = $derived(authState.session)
  const t = $derived($tStore)

  async function loadOrders() {
    if (!userId || !session?.access_token) {
      error = t.orders?.noSession ?? 'No hay sesión activa'
      loading = false
      return
    }

    try {
      loading = true
      error = null
      const currentMenuId = await getLatestMenuId(userId, session.access_token)
      if (!currentMenuId) {
        error = t.orders?.noMenu ?? 'No se encontró un menú'
        loading = false
        return
      }
      menuId = currentMenuId
      orders = await getMenuOrders(currentMenuId, session.access_token)
    } catch (err: unknown) {
      console.error('Error cargando órdenes:', err)
      error = err instanceof Error ? err.message : 'Error al cargar las órdenes'
    } finally {
      loading = false
    }
  }

  function formatRequestedTime(iso: string | null): string {
    if (!iso) return '—'
    const d = new Date(iso)
    return d.toLocaleString('es-CL', {
      dateStyle: 'short',
      timeStyle: 'short'
    })
  }

  function getFulfillmentType(payload: Record<string, unknown>): string {
    const type = (payload?.fulfillment as Record<string, unknown>)?.type
    return typeof type === 'string' ? type : '—'
  }

  function togglePayload(orderNumber: number) {
    expandedOrderNumber = expandedOrderNumber === orderNumber ? null : orderNumber
  }

  function formatCurrency(amount: number, currency: string): string {
    if (currency === 'CLP') return `$${amount.toLocaleString('es-CL')}`
    return `${amount} ${currency}`
  }

  function formatDetailDate(iso: string | null): string {
    if (!iso) return '—'
    const d = new Date(iso)
    return d.toLocaleString('es-CL', { dateStyle: 'medium', timeStyle: 'short' })
  }

  function openPaperView(order: MenuOrderRow) {
    paperOrder = order
  }

  function closePaperView() {
    paperOrder = null
  }

  function printThermal() {
    thermalPrintMode = true
    requestAnimationFrame(() => {
      window.print()
      thermalPrintMode = false
    })
  }

  onMount(() => {
    loadOrders()
  })
</script>

<div class="h-full flex flex-col bg-gray-50">
  <!-- Header -->
  <div class="flex-shrink-0 px-4 sm:px-6 py-4 border-b border-gray-200 bg-white">
    <div class="flex items-center justify-between gap-4">
      <button
        type="button"
        onclick={onMenuClick}
        class="md:hidden p-2 -ml-2 rounded-lg hover:bg-gray-100 text-gray-600"
        aria-label="Abrir menú"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800">
        {t.sidebar.orders}
      </h1>
    </div>
    <p class="mt-1 text-sm text-gray-500">
      {t.orders?.subtitle ?? 'Ordenado por número de orden y hora solicitada para planificar entregas o preparación.'}
    </p>
  </div>

  <!-- Content -->
  <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
    {#if loading}
      <div class="flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-10 w-10 border-2 border-blue-600 border-t-transparent"></div>
      </div>
    {:else if error}
      <div class="rounded-lg bg-red-50 border border-red-200 p-4 text-red-800">
        {error}
      </div>
    {:else if orders.length === 0}
      <div class="rounded-lg bg-gray-100 border border-gray-200 p-8 text-center text-gray-600">
        {t.orders?.empty ?? 'No hay órdenes aún.'}
      </div>
    {:else}
      <ul class="space-y-4">
        {#each orders as order (order.order_number)}
          {@const type = getFulfillmentType(order.event_payload)}
          <li class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
            <button
              type="button"
              class="w-full text-left px-4 py-4 sm:px-5 sm:py-4 flex items-center justify-between gap-3 hover:bg-gray-50 transition-colors"
              onclick={() => togglePayload(order.order_number)}
            >
              <div class="flex items-center gap-3 min-w-0">
                <span class="flex-shrink-0 font-semibold text-gray-800">#{order.order_number}</span>
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {type === 'DELIVERY' ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'}">
                  {type === 'DELIVERY' ? (t.orders?.delivery ?? 'Envío') : (t.orders?.pickup ?? 'Retiro')}
                </span>
                <span class="text-sm text-gray-500 truncate">
                  {formatRequestedTime(order.requested_time)}
                </span>
              </div>
              <svg
                class="w-5 h-5 text-gray-400 flex-shrink-0 transition-transform {expandedOrderNumber === order.order_number ? 'rotate-180' : ''}"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            {#if expandedOrderNumber === order.order_number}
              {@const p = order.event_payload}
              {@const items = (p?.items as Array<{ productName?: string; quantity?: number; unitPrice?: number; totalPrice?: number; unit?: string; pricingMode?: string }>) ?? []}
              {@const totals = (p?.totals as { subtotal?: number; deliveryFee?: number; total?: number; currency?: string }) ?? {}}
              {@const fulfillment = (p?.fulfillment as { type?: string; requestedTime?: string; address?: { rawAddress?: string; coordinates?: { latitude?: number; longitude?: number }; deliveryDetails?: { unit?: string; notes?: string } }; contact?: { fullName?: string; phone?: string; email?: string } }) ?? {}}
              {@const contact = fulfillment.contact ?? {}}
              {@const createdAt = p?.createdAt as string | undefined}
              <div class="border-t border-gray-200 bg-gray-50 px-4 py-4 sm:px-5 space-y-4">
                <!-- Items -->
                <section class="rounded-lg bg-white border border-gray-200 overflow-hidden">
                  <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-wide px-3 py-2 bg-gray-100 border-b border-gray-200">Items</h4>
                  <ul class="divide-y divide-gray-100">
                    {#each items as item}
                      <li class="px-3 py-2.5 flex justify-between items-baseline gap-2">
                        <span class="text-sm text-gray-800 font-medium">{item.productName ?? '—'}</span>
                        <span class="text-sm text-gray-600 shrink-0">{item.quantity ?? 0} × {formatCurrency(item.unitPrice ?? 0, totals.currency ?? 'CLP')} = {formatCurrency(item.totalPrice ?? 0, totals.currency ?? 'CLP')}</span>
                      </li>
                    {/each}
                  </ul>
                </section>
                <!-- Totals -->
                <section class="rounded-lg bg-white border border-gray-200 overflow-hidden">
                  <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-wide px-3 py-2 bg-gray-100 border-b border-gray-200">Totales</h4>
                  <dl class="px-3 py-2.5 space-y-1 text-sm">
                    <div class="flex justify-between"><dt class="text-gray-500">Subtotal</dt><dd class="text-gray-800 font-medium">{formatCurrency(totals.subtotal ?? 0, totals.currency ?? 'CLP')}</dd></div>
                    <div class="flex justify-between"><dt class="text-gray-500">Envío</dt><dd class="text-gray-800 font-medium">{formatCurrency(totals.deliveryFee ?? 0, totals.currency ?? 'CLP')}</dd></div>
                    <div class="flex justify-between border-t border-gray-100 pt-2 mt-2"><dt class="text-gray-700 font-semibold">Total</dt><dd class="text-gray-900 font-bold">{formatCurrency(totals.total ?? 0, totals.currency ?? 'CLP')}</dd></div>
                  </dl>
                </section>
                <!-- Contacto (destacado) -->
                <section class="rounded-lg bg-white border-2 border-blue-200 overflow-hidden">
                  <h4 class="text-xs font-semibold text-blue-800 uppercase tracking-wide px-3 py-2.5 bg-blue-50 border-b border-blue-200">Contacto</h4>
                  <dl class="px-3 py-3 space-y-2.5 text-sm">
                    <div><dt class="text-gray-500 text-xs mb-0.5">Nombre</dt><dd class="text-gray-900 font-semibold">{(contact.fullName || '').trim() || '—'}</dd></div>
                    <div><dt class="text-gray-500 text-xs mb-0.5">Teléfono</dt><dd class="text-gray-800">{contact.phone?.trim() || '—'}</dd></div>
                    <div><dt class="text-gray-500 text-xs mb-0.5">Email</dt><dd class="text-gray-800">{contact.email?.trim() || '—'}</dd></div>
                  </dl>
                </section>
                <!-- Entrega / Retiro -->
                <section class="rounded-lg bg-white border border-gray-200 overflow-hidden">
                  <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-wide px-3 py-2 bg-gray-100 border-b border-gray-200">Entrega / Retiro</h4>
                  <dl class="px-3 py-2.5 space-y-2 text-sm">
                    <div class="flex justify-between gap-2"><dt class="text-gray-500 shrink-0">Tipo</dt><dd class="text-gray-800 font-medium">{fulfillment.type ?? '—'}</dd></div>
                    <div class="flex justify-between gap-2"><dt class="text-gray-500 shrink-0">Hora solicitada</dt><dd class="text-gray-800">{formatDetailDate(fulfillment.requestedTime ?? null)}</dd></div>
                    {#if fulfillment.address?.rawAddress}
                      <div><dt class="text-gray-500 text-xs mb-0.5">Dirección</dt><dd class="text-gray-800">{fulfillment.address.rawAddress}</dd></div>
                      {#if fulfillment.address.deliveryDetails?.unit || fulfillment.address.deliveryDetails?.notes}
                        <div class="text-gray-600 text-xs">Depto/Unidad: {fulfillment.address.deliveryDetails?.unit ?? '—'} · Notas: {fulfillment.address.deliveryDetails?.notes || '—'}</div>
                      {/if}
                    {/if}
                  </dl>
                </section>
                {#if createdAt}
                  <p class="text-xs text-gray-400">Creado: {formatDetailDate(createdAt)}</p>
                {/if}
                <div class="flex justify-end pt-2">
                  <button
                    type="button"
                    onclick={() => openPaperView(order)}
                    class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-white border border-gray-300 text-gray-700 text-sm font-medium hover:bg-gray-50 shadow-sm print:hidden"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    {t.orders?.viewAsPaper ?? 'Ver como hoja'}
                  </button>
                </div>
              </div>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>

<style>
  /* Vista previa térmica en pantalla (ancho 80mm) */
  :global(.ticket-thermal-preview #ticket-print) {
    max-width: 80mm;
    margin-left: auto;
    margin-right: auto;
  }
  /* Impresión: solo el ticket, optimizado para térmica 80mm */
  @media print {
    :global(body) * {
      visibility: hidden;
    }
    :global(#ticket-print),
    :global(#ticket-print) * {
      visibility: visible;
    }
    :global(#ticket-print) {
      position: fixed !important;
      left: 0 !important;
      top: 0 !important;
      width: 80mm !important;
      max-width: 80mm !important;
      margin: 0 !important;
      padding: 4mm !important;
      color: black !important;
      background: white !important;
      font-family: monospace, 'Courier New', sans-serif !important;
      box-shadow: none !important;
      border: none !important;
    }
    :global(#ticket-print) * {
      color: black !important;
      background: transparent !important;
      border-color: black !important;
    }
  }
</style>

<!-- Modal vista papel -->
{#if paperOrder}
  {@const p = paperOrder.event_payload}
  {@const items = (p?.items as Array<{ productName?: string; quantity?: number; unitPrice?: number; totalPrice?: number }>) ?? []}
  {@const totals = (p?.totals as { subtotal?: number; deliveryFee?: number; total?: number; currency?: string }) ?? {}}
  {@const fulfillment = (p?.fulfillment as { type?: string; requestedTime?: string; address?: { rawAddress?: string; deliveryDetails?: { unit?: string; notes?: string } }; contact?: { fullName?: string; phone?: string; email?: string } }) ?? {}}
  {@const contact = fulfillment.contact ?? {}}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 print:bg-transparent print:p-0"
    role="dialog"
    aria-modal="true"
    aria-labelledby="paper-order-title"
  >
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-hidden flex flex-col print:max-h-none print:shadow-none print:rounded-none">
      <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200 print:hidden">
        <h2 id="paper-order-title" class="text-lg font-semibold text-gray-800">{t.orders?.viewAsPaper ?? 'Ver como hoja'}</h2>
        <div class="flex items-center gap-2">
          <button type="button" onclick={printThermal} class="px-3 py-1.5 rounded-lg bg-gray-800 text-white text-sm font-medium hover:bg-gray-900">
            {t.orders?.printThermal ?? 'Imprimir en térmica'}
          </button>
          <button type="button" onclick={closePaperView} class="p-2 rounded-lg hover:bg-gray-100 text-gray-600" aria-label="Cerrar">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
      </div>
      <div class="overflow-y-auto p-6 print:p-0 {thermalPrintMode ? 'ticket-thermal-preview' : ''}" style="background: linear-gradient(to bottom, #fafafa 0%, #fff 2rem);">
        <div id="ticket-print" class="max-w-sm mx-auto text-gray-900">
          <div class="border-b-2 border-dashed border-gray-300 pb-3 mb-3">
            <p class="text-center text-xs text-gray-500 uppercase tracking-widest">Orden</p>
            <p class="text-center text-3xl font-bold text-gray-900 mt-0.5">#{paperOrder.order_number}</p>
            <p class="text-center text-xl font-semibold mt-2 text-gray-700">{formatDetailDate(paperOrder.requested_time)}</p>
            <p class="text-center text-base text-gray-600 mt-0.5">{fulfillment.type === 'DELIVERY' ? (t.orders?.delivery ?? 'Envío') : (t.orders?.pickup ?? 'Retiro')}</p>
          </div>
          <div class="space-y-3 text-base">
            <div>
              <p class="font-semibold text-gray-700 mb-1 text-sm">Contacto</p>
              <p class="font-medium text-lg">{(contact.fullName || '').trim() || '—'}</p>
              <p class="text-lg">{contact.phone?.trim() || '—'}</p>
              {#if contact.email?.trim()}<p class="text-gray-600 text-sm">{contact.email}</p>{/if}
            </div>
            {#if fulfillment.address?.rawAddress}
              <div>
                <p class="font-semibold text-gray-700 mb-1 text-sm">Dirección</p>
                <p class="text-gray-800">{fulfillment.address.rawAddress}</p>
                {#if fulfillment.address.deliveryDetails?.unit || fulfillment.address.deliveryDetails?.notes}
                  <p class="text-gray-600 text-sm mt-0.5">Depto/Unidad: {fulfillment.address.deliveryDetails?.unit ?? '—'} · Notas: {fulfillment.address.deliveryDetails?.notes || '—'}</p>
                {/if}
              </div>
            {/if}
            <div class="border-t border-dashed border-gray-300 pt-3">
              <p class="font-semibold text-gray-700 mb-2 text-sm">Pedido</p>
              <ul class="space-y-2">
                {#each items as item}
                  <li class="flex justify-between gap-2 items-baseline">
                    <span class="text-base"><span class="font-bold text-lg">{item.quantity ?? 0}</span> × {item.productName ?? '—'}</span>
                    <span class="shrink-0 text-lg font-semibold">{formatCurrency(item.totalPrice ?? 0, totals.currency ?? 'CLP')}</span>
                  </li>
                {/each}
              </ul>
            </div>
            <div class="border-t-2 border-dashed border-gray-300 pt-3 flex justify-between items-baseline text-xl font-bold">
              <span>Total</span>
              <span class="text-2xl">{formatCurrency(totals.total ?? 0, totals.currency ?? 'CLP')}</span>
            </div>
          </div>
          <div class="mt-6 pt-3 border-t border-dashed border-gray-300 text-center text-sm text-gray-400">
            {formatDetailDate(p?.createdAt as string ?? null)}
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
