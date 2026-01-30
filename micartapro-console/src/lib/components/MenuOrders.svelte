<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getMenuOrders, subscribeMenuOrdersRealtime, type MenuOrderRow } from '../menuUtils'
  import { t as tStore } from '../useLanguage'

  interface MenuOrdersProps {
    onMenuClick?: () => void
    onKitchenModeChange?: (active: boolean) => void
  }

  let { onMenuClick, onKitchenModeChange }: MenuOrdersProps = $props()

  let orders = $state<MenuOrderRow[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let menuId = $state<string | null>(null)
  let paperOrder = $state<MenuOrderRow | null>(null)
  let thermalPrintMode = $state(false)
  let orderStatus = $state<Record<number, 'pending' | 'preparing' | 'done'>>({})
  let kitchenMode = $state(false)
  let realtimeUnsubscribe = $state<(() => void) | null>(null)

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const session = $derived(authState.session)
  const t = $derived($tStore)

  /** Devuelve el menuId cargado (para suscribir Realtime justo después). */
  async function loadOrders(): Promise<string | null> {
    if (!userId || !session?.access_token) {
      error = t.orders?.noSession ?? 'No hay sesión activa'
      loading = false
      return null
    }

    try {
      loading = true
      error = null
      const currentMenuId = await getLatestMenuId(userId, session.access_token)
      if (!currentMenuId) {
        error = t.orders?.noMenu ?? 'No se encontró un menú'
        loading = false
        return null
      }
      menuId = currentMenuId
      orders = await getMenuOrders(currentMenuId, session.access_token)
      return currentMenuId
    } catch (err: unknown) {
      console.error('Error cargando órdenes:', err)
      error = err instanceof Error ? err.message : 'Error al cargar las órdenes'
      return null
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

  function formatCurrency(amount: number, currency: string): string {
    if (currency === 'CLP') return `$${amount.toLocaleString('es-CL')}`
    return `${amount} ${currency}`
  }

  function formatDetailDate(iso: string | null): string {
    if (!iso) return '—'
    const d = new Date(iso)
    return d.toLocaleString('es-CL', { dateStyle: 'medium', timeStyle: 'short' })
  }

  /** Minutos hasta requested_time (negativo = atrasado). */
  function getRemainingMinutes(iso: string | null): number | null {
    if (!iso) return null
    const target = new Date(iso).getTime()
    const now = Date.now()
    return Math.round((target - now) / 60_000)
  }

  function getRemainingTimeLabel(minutes: number | null): string {
    if (minutes === null) return ''
    const abs = Math.abs(minutes)
    const isLate = minutes < 0
    const template = isLate ? (t.orders?.late ?? 'Atrasado {min} min') : (t.orders?.remainingIn ?? 'En {min} min')
    return template.replace('{min}', String(abs))
  }

  /** 'green' >15 min, 'yellow' 5–15, 'red' <5 o atrasado */
  function getRemainingTimeColor(minutes: number | null): 'green' | 'yellow' | 'red' | null {
    if (minutes === null) return null
    if (minutes < 0) return 'red'
    if (minutes <= 5) return 'red'
    if (minutes <= 15) return 'yellow'
    return 'green'
  }

  /** Agrupa ítems por productName, sumando cantidades. */
  function groupItems(items: Array<{ productName?: string; quantity?: number }>): Array<{ productName: string; quantity: number }> {
    const map = new Map<string, number>()
    for (const it of items) {
      const name = it.productName?.trim() || '—'
      map.set(name, (map.get(name) ?? 0) + (it.quantity ?? 0))
    }
    return [...map.entries()].map(([productName, quantity]) => ({ productName, quantity }))
  }

  function getOrderStatus(orderNumber: number): 'pending' | 'preparing' | 'done' {
    return orderStatus[orderNumber] ?? 'pending'
  }

  function setOrderStatus(orderNumber: number, status: 'pending' | 'preparing' | 'done') {
    orderStatus = { ...orderStatus, [orderNumber]: status }
  }

  function cycleOrderStatus(orderNumber: number) {
    const current = getOrderStatus(orderNumber)
    if (current === 'pending') setOrderStatus(orderNumber, 'preparing')
    else if (current === 'preparing') setOrderStatus(orderNumber, 'done')
    else setOrderStatus(orderNumber, 'pending')
  }

  async function toggleKitchenMode() {
    kitchenMode = !kitchenMode
    onKitchenModeChange?.(kitchenMode)
    if (kitchenMode) {
      try {
        await document.documentElement.requestFullscreen?.()
      } catch {
        // ignore if fullscreen not allowed
      }
    } else {
      try {
        await document.exitFullscreen?.()
      } catch {
        // ignore
      }
    }
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
    let cancelled = false
    let pollInterval: ReturnType<typeof setInterval> | null = null
    ;(async () => {
      const mid = await loadOrders()
      if (cancelled) return
      const token = session?.access_token
      if (mid && token) {
        realtimeUnsubscribe = await subscribeMenuOrdersRealtime(mid, token, () => {
          loadOrders()
        })
        // Polling cada 20 s por si Realtime no está habilitado en Supabase
        pollInterval = setInterval(() => {
          loadOrders()
        }, 20_000)
      }
    })()
    return () => {
      cancelled = true
      if (pollInterval) clearInterval(pollInterval)
      realtimeUnsubscribe?.()
      realtimeUnsubscribe = null
    }
  })
</script>

<div class="h-full flex flex-col bg-gray-50 kitchen-orders-root" class:kitchen-mode={kitchenMode}>
  <!-- Header -->
  <div class="flex-shrink-0 px-4 sm:px-6 py-4 border-b border-gray-200 bg-white" class:kitchen-mode-header={kitchenMode}>
    <div class="flex items-center justify-between gap-4">
      {#if !kitchenMode}
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
      {/if}
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800" class:text-3xl={kitchenMode}>
        {t.sidebar.kitchen}
      </h1>
      <div class="flex items-center gap-2">
        <button
          type="button"
          onclick={() => loadOrders()}
          class="rounded-lg px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 border border-gray-200"
          title={t.orders?.reload ?? 'Recargar'}
        >
          {t.orders?.reload ?? 'Recargar'}
        </button>
        <button
          type="button"
          onclick={toggleKitchenMode}
        class="rounded-lg px-4 py-2 text-sm font-semibold {kitchenMode ? 'bg-amber-500 text-white hover:bg-amber-600' : 'bg-gray-200 text-gray-800 hover:bg-gray-300'}"
      >
          {kitchenMode ? (t.orders?.exitKitchenMode ?? 'Salir modo cocina') : (t.orders?.kitchenMode ?? 'Modo cocina')}
        </button>
      </div>
    </div>
    {#if !kitchenMode}
      <p class="mt-1 text-sm text-gray-500">
        {t.orders?.subtitle ?? 'Ordenado por hora comprometida. Vista orientada a cocina.'}
      </p>
    {/if}
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
      <ul class="space-y-5 kitchen-orders-list" class:kitchen-mode-list={kitchenMode}>
        {#each orders as order, index (order.order_number)}
          {@const type = getFulfillmentType(order.event_payload)}
          {@const rawItems = (order.event_payload?.items as Array<{ productName?: string; quantity?: number }>) ?? []}
          {@const items = groupItems(rawItems)}
          {@const itemCount = rawItems.reduce((s, i) => s + (i.quantity ?? 0), 0)}
          {@const isFirst = index === 0}
          {@const status = getOrderStatus(order.order_number)}
          {@const remainingMin = getRemainingMinutes(order.requested_time)}
          {@const timeColor = getRemainingTimeColor(remainingMin)}
          <li class="bg-white rounded-xl border-2 overflow-hidden kitchen-order-card {isFirst ? 'kitchen-order-first border-amber-400 shadow-lg' : 'border-gray-200 shadow-sm'}">
            <!-- Cabecera cocina: número, hora, tipo, tiempo restante, estado (sin expandir) -->
            <div class="w-full px-4 py-3 sm:px-5 flex flex-wrap items-center gap-4 border-b border-gray-100 {isFirst ? 'sm:py-6' : 'sm:py-4'}">
              <span class="font-bold text-gray-900 tabular-nums {isFirst ? 'text-4xl sm:text-5xl md:text-6xl' : 'text-3xl sm:text-4xl'}">#{order.order_number}</span>
              <span class="font-semibold text-gray-700 {isFirst ? 'text-2xl sm:text-3xl md:text-4xl' : 'text-xl sm:text-2xl'}">
                {(t.orders?.forTime ?? 'Para')} {formatRequestedTime(order.requested_time)}
              </span>
              {#if remainingMin !== null}
                <span class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-sm font-bold tabular-nums
                  {timeColor === 'green' ? 'bg-green-100 text-green-800' : ''}
                  {timeColor === 'yellow' ? 'bg-amber-200 text-amber-900' : ''}
                  {timeColor === 'red' ? 'bg-red-100 text-red-800' : ''}">
                  <span aria-hidden="true">{timeColor === 'green' ? '🟢' : timeColor === 'yellow' ? '🟡' : '🔴'}</span>
                  {getRemainingTimeLabel(remainingMin)}
                </span>
              {/if}
              <span class="inline-flex items-center rounded-full font-medium {type === 'DELIVERY' ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'} {isFirst ? 'px-4 py-2 text-base sm:text-lg' : 'px-3 py-1 text-sm'}">
                {type === 'DELIVERY' ? (t.orders?.delivery ?? 'Envío') : (t.orders?.pickup ?? 'Retiro')}
              </span>
              <!-- Estado operativo: informativo (ámbar suave), no compite con el botón naranja -->
              <span class="inline-flex items-center gap-1 rounded-full font-bold border {isFirst ? 'px-4 py-2 text-base sm:text-lg' : 'px-3 py-1 text-sm'}
                {status === 'pending' ? 'bg-gray-100 text-gray-700 border-gray-200' : ''}
                {status === 'preparing' ? 'bg-amber-50 text-amber-900 border-amber-200' : ''}
                {status === 'done' ? 'bg-green-50 text-green-800 border-green-200' : ''}">
                {#if status === 'preparing'}<span aria-hidden="true">⏳</span>{/if}
                {status === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : status === 'preparing' ? (t.orders?.statusPreparing ?? 'En preparación') : (t.orders?.statusDone ?? 'Listo')}
              </span>
            </div>
            <!-- Qué preparar: listado vertical (un ítem por línea; cantidad en bold) -->
            <div class="px-4 py-3 sm:px-5 bg-amber-50/50 border-b border-amber-100 {isFirst ? 'py-4 sm:py-5' : ''}">
              <div class="flex items-center justify-between gap-2 mb-2">
                <p class="font-semibold text-amber-800 uppercase tracking-wide {isFirst ? 'text-sm' : 'text-xs'}">{t.orders?.itemsToPrepare ?? 'Qué preparar'}</p>
                <span class="text-sm font-bold text-amber-800 tabular-nums">{(t.orders?.itemsCount ?? '{count} ítems').replace('{count}', String(itemCount))}</span>
              </div>
              <ul class="space-y-1 text-gray-900 {isFirst ? 'text-xl sm:text-2xl md:text-3xl' : 'text-lg sm:text-xl'}">
                {#each items as item}
                  <li class="tabular-nums">
                    <span class="font-bold text-amber-800">{item.quantity}×</span> <span class="font-normal">{item.productName}</span>
                  </li>
                {/each}
              </ul>
            </div>
            <!-- Un solo CTA primario en la tarjeta (sin abrir el pedido) -->
            <div class="px-4 py-3 sm:px-5 border-t border-gray-100">
              {#if status === 'pending'}
                <button
                  type="button"
                  onclick={(e) => { e.stopPropagation(); cycleOrderStatus(order.order_number); }}
                  class="w-full py-3 px-4 rounded-xl text-base font-bold bg-orange-500 hover:bg-orange-600 text-white shadow-md transition-colors"
                >
                  🔥 {t.orders?.startPreparing ?? 'Iniciar preparación'}
                </button>
              {:else if status === 'preparing'}
                <button
                  type="button"
                  onclick={(e) => { e.stopPropagation(); cycleOrderStatus(order.order_number); }}
                  class="w-full py-3 px-4 rounded-xl text-base font-bold bg-amber-500 hover:bg-amber-600 text-white shadow-md transition-colors"
                >
                  ✓ {t.orders?.markAsReady ?? 'Marcar como listo'}
                </button>
              {:else}
                <div class="w-full py-3 px-4 rounded-xl text-base font-bold bg-green-100 text-green-800 text-center">
                  ✓ {t.orders?.statusDone ?? 'Listo'}
                </div>
              {/if}
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>

<style>
  /* Modo cocina: fullscreen, alto contraste, letra grande */
  :global(.kitchen-mode) {
    background: #f5f5f5;
  }
  :global(.kitchen-mode .kitchen-mode-header) {
    padding: 0.75rem 1rem;
    border-bottom-width: 2px;
  }
  :global(.kitchen-mode .kitchen-orders-list) {
    padding: 0.5rem;
  }
  :global(.kitchen-mode .kitchen-order-card) {
    font-size: 1.05rem;
  }
  :global(.kitchen-mode .kitchen-order-first) {
    font-size: 1.15rem;
  }
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
