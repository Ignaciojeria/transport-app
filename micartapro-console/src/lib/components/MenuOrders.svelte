<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getKitchenOrdersFromProjection, subscribeMenuOrdersRealtime, getMenuOrderByNumber, type KitchenOrder, type MenuOrderRow, type StationFilter } from '../menuUtils'
  import { t as tStore } from '../useLanguage'

  interface MenuOrdersProps {
    onMenuClick?: () => void
    onKitchenModeChange?: (active: boolean) => void
  }

  let { onMenuClick, onKitchenModeChange }: MenuOrdersProps = $props()

  let orders = $state<KitchenOrder[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let menuId = $state<string | null>(null)
  let paperOrder = $state<MenuOrderRow | null>(null)
  let thermalPrintMode = $state(false)
  /** Estado por (orden, estación): INICIAR/LISTO se aplican a BAR y COCINA de forma independiente. */
  let orderStatus = $state<Record<string, 'pending' | 'preparing' | 'done'>>({})
  let kitchenMode = $state(false)
  /** Filtro por estación: se aplica en Supabase (columna station). */
  let stationFilter = $state<StationFilter>('ALL')
  /** Tab operativo en Cocina/Bar: Pendientes | En preparación | Listos (solo cuando stationFilter es KITCHEN o BAR). */
  let operationalTab = $state<'pending' | 'preparing' | 'done'>('pending')
  /** Vista QR: muestra códigos Cocina/Barra en lugar de la lista de órdenes. */
  let showQRView = $state(false)
  /** Vista órdenes en Cocina/Bar: vertical (tabs + lista) o kanban (3 columnas). */
  let ordersViewMode = $state<'vertical' | 'kanban'>('vertical')
  let realtimeUnsubscribe = $state<(() => void) | null>(null)
  /** Modal QR ampliado: KITCHEN | BAR | ALL | null */
  let qrEnlarged = $state<'KITCHEN' | 'BAR' | 'ALL' | null>(null)

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
      orders = await getKitchenOrdersFromProjection(currentMenuId, session.access_token, stationFilter)
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

  /** fulfillment ya viene de la proyección (string: PICKUP | DELIVERY). */
  function getFulfillmentLabel(fulfillment: string): string {
    return fulfillment === 'DELIVERY' ? (t.orders?.delivery ?? 'Envío') : (t.orders?.pickup ?? 'Retiro')
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

  /** Total de unidades de ítems (la proyección ya agrupa por nombre en menuUtils). */
  function getItemCount(items: KitchenOrder['items']): number {
    return items.reduce((s, i) => s + i.quantity, 0)
  }

  /** Órdenes mostradas = las que vienen de Supabase ya filtradas por station. */
  const displayedOrders = $derived(orders)

  /** En Cocina/Bar: órdenes particionadas por estado (pending / preparing / done) para los tabs. */
  const ordersByTab = $derived.by(() => {
    if (stationFilter !== 'KITCHEN' && stationFilter !== 'BAR') return { pending: [] as KitchenOrder[], preparing: [] as KitchenOrder[], done: [] as KitchenOrder[] }
    const pending: KitchenOrder[] = []
    const preparing: KitchenOrder[] = []
    const done: KitchenOrder[] = []
    for (const o of displayedOrders) {
      const st = getOrderStatus(o.order_number, stationFilter)
      if (st === 'pending') pending.push(o)
      else if (st === 'preparing') preparing.push(o)
      else done.push(o)
    }
    return { pending, preparing, done }
  })

  /** Órdenes a mostrar: en Caja = todas; en Cocina/Bar = las del tab activo. */
  const ordersToShow = $derived(stationFilter === 'ALL' ? displayedOrders : ordersByTab[operationalTab])

  async function setStationFilterAndReload(filter: StationFilter) {
    stationFilter = filter
    if (filter === 'KITCHEN' || filter === 'BAR') operationalTab = 'pending'
    if (menuId && session?.access_token) {
      loading = true
      try {
        orders = await getKitchenOrdersFromProjection(menuId, session.access_token, filter)
      } finally {
        loading = false
      }
    }
  }

  const statusKey = (orderNumber: number, station: string) => `${orderNumber}-${station}`

  function getOrderStatus(orderNumber: number, station: string): 'pending' | 'preparing' | 'done' {
    return orderStatus[statusKey(orderNumber, station)] ?? 'pending'
  }

  function setOrderStatus(orderNumber: number, station: string, status: 'pending' | 'preparing' | 'done') {
    orderStatus = { ...orderStatus, [statusKey(orderNumber, station)]: status }
  }

  function cycleOrderStatus(orderNumber: number, station: string) {
    const current = getOrderStatus(orderNumber, station)
    if (current === 'pending') setOrderStatus(orderNumber, station, 'preparing')
    else if (current === 'preparing') setOrderStatus(orderNumber, station, 'done')
    else setOrderStatus(orderNumber, station, 'pending')
  }

  /** Estaciones que tienen ítems en esta orden (KITCHEN, BAR). */
  function getStationsInOrder(order: KitchenOrder): string[] {
    const stations = new Set<string>()
    for (const i of order.items) {
      if (i.station) stations.add(i.station)
    }
    return stations.size > 0 ? [...stations] : ['KITCHEN']
  }

  /** Indica si Cocina y Barra marcaron la orden como lista (solo informativo; la entrega no está bloqueada por estaciones). */
  function isOrderReadyForDelivery(order: KitchenOrder): boolean {
    const stations = getStationsInOrder(order)
    return stations.every((st) => getOrderStatus(order.order_number, st) === 'done')
  }

  /** En Caja: estado resumido para mostrar (Pendiente / En preparación). */
  function getCajaOrderStatusLabel(order: KitchenOrder): 'pending' | 'preparing' {
    const stations = getStationsInOrder(order)
    const statuses = stations.map((st) => getOrderStatus(order.order_number, st))
    if (statuses.every((s) => s === 'pending')) return 'pending'
    return 'preparing'
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

  async function openPaperView(order: KitchenOrder) {
    if (!menuId || !session?.access_token) return
    const full = await getMenuOrderByNumber(menuId, order.order_number, session.access_token)
    paperOrder = full
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
        {t.sidebar.orders}
      </h1>
      <div class="flex items-center gap-2">
        <button
          type="button"
          onclick={toggleKitchenMode}
          class="rounded-lg px-4 py-2 text-sm font-semibold {kitchenMode ? 'bg-amber-500 text-white hover:bg-amber-600' : 'bg-gray-200 text-gray-800 hover:bg-gray-300'}"
        >
          {kitchenMode ? (t.orders?.exitKitchenMode ?? 'Salir modo full') : (t.orders?.kitchenMode ?? 'Modo full')}
        </button>
      </div>
    </div>
    {#if !kitchenMode}
      <p class="mt-1 text-sm text-gray-500">
        {t.orders?.subtitle ?? 'Ordenado por hora comprometida. Vista orientada a cocina.'}
      </p>
      <!-- Filtro por estación + opción QR: oculto en modo full -->
      <div class="mt-3 flex flex-wrap items-center gap-2">
        <button
          type="button"
          onclick={() => { showQRView = false; setStationFilterAndReload('ALL'); }}
          class="rounded-lg px-4 py-2 text-sm font-semibold transition-colors {!showQRView && stationFilter === 'ALL' ? 'bg-gray-800 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
        >
          {t.orders?.filterAll ?? 'Entrega'}
        </button>
        <button
          type="button"
          onclick={() => { showQRView = false; setStationFilterAndReload('KITCHEN'); }}
          class="rounded-lg px-4 py-2 text-sm font-semibold transition-colors {!showQRView && stationFilter === 'KITCHEN' ? 'bg-amber-600 text-white' : 'bg-amber-50 text-amber-800 hover:bg-amber-100 border border-amber-200'}"
        >
          {t.orders?.filterKitchen ?? 'Cocina'}
        </button>
        <button
          type="button"
          onclick={() => { showQRView = false; setStationFilterAndReload('BAR'); }}
          class="rounded-lg px-4 py-2 text-sm font-semibold transition-colors {!showQRView && stationFilter === 'BAR' ? 'bg-blue-600 text-white' : 'bg-blue-50 text-blue-800 hover:bg-blue-100 border border-blue-200'}"
        >
          {t.orders?.filterBar ?? 'Barra'}
        </button>
        {#if menuId && session?.access_token}
          <button
            type="button"
            onclick={() => showQRView = true}
            class="rounded-lg px-4 py-2 text-sm font-semibold transition-colors {showQRView ? 'bg-gray-700 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200 border border-gray-300'}"
          >
            {t.orders?.showQR ?? 'QR'}
          </button>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Content -->
  <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
    {#if showQRView && menuId && session?.access_token}
      <!-- Vista QR: códigos Entrega, Cocina y Barra sin tapar filtros -->
      {@const baseUrl = typeof window !== 'undefined' ? window.location.origin : ''}
      {@const hashParams = session?.refresh_token ? `token=${encodeURIComponent(session.access_token)}&refresh_token=${encodeURIComponent(session.refresh_token)}` : `token=${encodeURIComponent(session.access_token)}`}
      {@const urlCaja = baseUrl ? `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=ALL#${hashParams}` : ''}
      {@const urlCocina = baseUrl ? `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=KITCHEN#${hashParams}` : ''}
      {@const urlBarra = baseUrl ? `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=BAR#${hashParams}` : ''}
      <div class="max-w-3xl">
        <p class="text-sm font-medium text-gray-700 mb-2">Acceso sin login (entrega, cocinero o barista escanea el código)</p>
        <p class="text-xs text-gray-500 mb-4">Haz clic en el código para agrandarlo.</p>
        <div class="flex flex-wrap gap-6">
          <div class="flex items-start gap-3 p-4 rounded-xl bg-gray-100 border border-gray-300">
            <button
              type="button"
              onclick={() => qrEnlarged = 'ALL'}
              class="flex-shrink-0 w-28 h-28 bg-white rounded-lg border border-gray-300 overflow-hidden cursor-pointer hover:ring-2 hover:ring-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
              title="Clic para agrandar"
            >
              {#if urlCaja}
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=112x112&data=${encodeURIComponent(urlCaja)}`} alt="QR Entrega" class="w-full h-full object-contain pointer-events-none" />
              {/if}
            </button>
            <div>
              <p class="font-semibold text-gray-900">{t.orders?.filterAll ?? 'Entrega'}</p>
              <p class="text-xs text-gray-700 mb-1">Ver todas las órdenes en tiempo real</p>
              <button type="button" onclick={() => urlCaja && navigator.clipboard.writeText(urlCaja).then(() => alert('Enlace copiado'))} class="text-xs text-gray-600 underline hover:no-underline">Copiar enlace</button>
            </div>
          </div>
          <div class="flex items-start gap-3 p-4 rounded-xl bg-amber-50 border border-amber-200">
            <button
              type="button"
              onclick={() => qrEnlarged = 'KITCHEN'}
              class="flex-shrink-0 w-28 h-28 bg-white rounded-lg border border-amber-200 overflow-hidden cursor-pointer hover:ring-2 hover:ring-amber-400 focus:outline-none focus:ring-2 focus:ring-amber-500"
              title="Clic para agrandar"
            >
              {#if urlCocina}
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=112x112&data=${encodeURIComponent(urlCocina)}`} alt="QR Cocina" class="w-full h-full object-contain pointer-events-none" />
              {/if}
            </button>
            <div>
              <p class="font-semibold text-amber-900">{t.orders?.filterKitchen ?? 'Cocina'}</p>
              <p class="text-xs text-amber-800 mb-1">Escanear → pedir nombre → ver pedidos en tiempo real</p>
              <button type="button" onclick={() => urlCocina && navigator.clipboard.writeText(urlCocina).then(() => alert('Enlace copiado'))} class="text-xs text-amber-700 underline hover:no-underline">Copiar enlace</button>
            </div>
          </div>
          <div class="flex items-start gap-3 p-4 rounded-xl bg-blue-50 border border-blue-200">
            <button
              type="button"
              onclick={() => qrEnlarged = 'BAR'}
              class="flex-shrink-0 w-28 h-28 bg-white rounded-lg border border-blue-200 overflow-hidden cursor-pointer hover:ring-2 hover:ring-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              title="Clic para agrandar"
            >
              {#if urlBarra}
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=112x112&data=${encodeURIComponent(urlBarra)}`} alt="QR Barra" class="w-full h-full object-contain pointer-events-none" />
              {/if}
            </button>
            <div>
              <p class="font-semibold text-blue-900">{t.orders?.filterBar ?? 'Barra'}</p>
              <p class="text-xs text-blue-800 mb-1">Escanear → pedir nombre → ver pedidos en tiempo real</p>
              <button type="button" onclick={() => urlBarra && navigator.clipboard.writeText(urlBarra).then(() => alert('Enlace copiado'))} class="text-xs text-blue-700 underline hover:no-underline">Copiar enlace</button>
            </div>
          </div>
        </div>
      </div>
    {:else if loading}
      <div class="flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-10 w-10 border-2 border-blue-600 border-t-transparent"></div>
      </div>
    {:else if error}
      <div class="rounded-lg bg-red-50 border border-red-200 p-4 text-red-800">
        {error}
      </div>
    {:else}
      <!-- Snippet: una card de orden (Cocina/Bar/Caja). compactStatus=true en Kanban: icono pequeño en cabecera en vez de badge. -->
      {#snippet orderCard(order: KitchenOrder, status: 'pending' | 'preparing' | 'done', isFirst: boolean, compactStatus: boolean = false)}
        {@const type = order.fulfillment}
        {@const itemCount = getItemCount(order.items)}
        {@const remainingMin = getRemainingMinutes(order.requested_time)}
        {@const timeColor = getRemainingTimeColor(remainingMin)}
        {@const isBarOrder = order.items.some((i) => i.station === 'BAR')}
        {@const useBarColor = stationFilter === 'BAR' || (stationFilter === 'ALL' && isBarOrder)}
        {@const readyForDelivery = isOrderReadyForDelivery(order)}
        {@const isDoneTab = status === 'done'}
        {@const barStatusForOrder = getOrderStatus(order.order_number, 'BAR')}
        {@const hasBarItems = order.items.some((i) => i.station === 'BAR')}
        {@const statusIcon = status === 'pending' ? '🟠' : status === 'preparing' ? '⏳' : '✓'}
        {@const statusTitle = status === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : status === 'preparing' ? (t.orders?.statusPreparing ?? 'En preparación') : (t.orders?.statusDone ?? 'Listo')}
        <li class="bg-white rounded-xl border-2 overflow-hidden kitchen-order-card order-card {isDoneTab ? 'order-card-done' : ''} {isFirst && !isDoneTab ? 'kitchen-order-first border-amber-400 shadow-lg' : 'border-gray-200 shadow-sm'}">
          <div class="w-full px-4 py-3 sm:px-5 flex flex-wrap items-center gap-4 border-b border-gray-100 {isFirst && !isDoneTab ? 'sm:py-6' : 'sm:py-4'} {compactStatus ? 'py-2 sm:py-3' : ''}">
            <span class="font-bold text-gray-900 tabular-nums {isFirst && !isDoneTab && !compactStatus ? 'text-4xl sm:text-5xl md:text-6xl' : compactStatus ? 'text-2xl sm:text-3xl' : 'text-3xl sm:text-4xl'}">#{order.order_number}</span>
            {#if compactStatus && stationFilter !== 'ALL'}
              <span class="text-base opacity-90" aria-hidden="true" title={statusTitle}>{statusIcon}</span>
            {/if}
            {#if stationFilter === 'KITCHEN' && hasBarItems}
              <span class="text-lg" aria-hidden="true" title="{barStatusForOrder === 'done' ? (t.orders?.statusDone ?? 'Bar listo') : (t.orders?.statusPreparing ?? 'Bar en preparación')}">{barStatusForOrder === 'done' ? '🍺 ✔️' : '🍺 ⏳'}</span>
            {/if}
            <span class="font-semibold text-gray-700 {isFirst && !compactStatus ? 'text-2xl sm:text-3xl md:text-4xl' : compactStatus ? 'text-base sm:text-lg' : 'text-xl sm:text-2xl'}">
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
            <span class="inline-flex items-center rounded-full font-medium {type === 'DELIVERY' ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'} {isFirst && !compactStatus ? 'px-4 py-2 text-base sm:text-lg' : 'px-3 py-1 text-sm'}">
              {getFulfillmentLabel(type)}
            </span>
            {#if stationFilter === 'ALL'}
              {@const kitchenSt = getOrderStatus(order.order_number, 'KITCHEN')}
              {@const barSt = getOrderStatus(order.order_number, 'BAR')}
              {@const orderHasBar = order.items.some((i) => i.station === 'BAR')}
              <div class="flex flex-wrap items-center gap-3 text-sm font-semibold">
                <span class="inline-flex items-center gap-1">{t.orders?.filterKitchen ?? 'Cocina'}: {kitchenSt === 'done' ? '✔️' : '⏳'}</span>
                <span class="inline-flex items-center gap-1">{t.orders?.filterBar ?? 'Barra'}: {orderHasBar ? (barSt === 'done' ? '✔️' : '⏳') : '—'}</span>
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-1 {readyForDelivery ? 'bg-green-50 text-green-800 border-green-200' : 'bg-amber-50 text-amber-900 border-amber-200'}">
                  {t.orders?.statusGeneralLabel ?? 'Estado general'}: {readyForDelivery ? (t.orders?.readyToDeliver ?? 'Listo para entregar') : (t.orders?.statusPreparing ?? 'En preparación')}
                </span>
              </div>
            {:else if !compactStatus}
              <span class="inline-flex items-center gap-1 rounded-full font-bold border {isFirst ? 'px-4 py-2 text-base sm:text-lg' : 'px-3 py-1 text-sm'}
                {status === 'pending' ? 'bg-gray-100 text-gray-700 border-gray-200' : ''}
                {status === 'preparing' ? 'bg-amber-50 text-amber-900 border-amber-200' : ''}
                {status === 'done' ? 'bg-green-50 text-green-800 border-green-200' : ''}">
                {#if status === 'preparing'}<span aria-hidden="true">⏳</span>{/if}
                {status === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : status === 'preparing' ? (t.orders?.statusPreparing ?? 'En preparación') : (t.orders?.statusDone ?? 'Listo')}
              </span>
            {/if}
          </div>
          <div class="px-4 py-3 sm:px-5 bg-amber-50/50 border-b border-amber-100 {isFirst && !compactStatus ? 'py-4 sm:py-5' : ''} {compactStatus ? 'py-2 sm:py-3' : ''}">
            <div class="flex items-center justify-between gap-2 mb-2">
              <p class="font-semibold text-amber-800 uppercase tracking-wide {isFirst && !compactStatus ? 'text-sm' : 'text-xs'}">{t.orders?.itemsToPrepare ?? 'Qué preparar'}</p>
              <span class="text-sm font-bold text-amber-800 tabular-nums">{(t.orders?.itemsCount ?? '{count} ítems').replace('{count}', String(itemCount))}</span>
            </div>
            <ul class="space-y-1 text-gray-900 {compactStatus ? 'text-base sm:text-lg' : isFirst ? 'text-xl sm:text-2xl md:text-3xl' : 'text-lg sm:text-xl'}">
              {#each order.items as item}
                <li class="tabular-nums">
                  <span class="font-bold text-amber-800">{item.quantity}×</span> <span class="font-normal">{item.item_name}</span>
                </li>
              {/each}
            </ul>
          </div>
          <div class="px-4 py-3 sm:px-5 border-t border-gray-100 {compactStatus ? 'py-2 sm:py-3' : ''}">
            {#if stationFilter === 'ALL'}
              <button
                type="button"
                onclick={(e) => { e.stopPropagation(); /* TODO: acción entregar */ }}
                class="w-full py-3 px-4 rounded-xl text-base font-bold bg-green-600 hover:bg-green-700 text-white shadow-md transition-colors"
              >
                {t.orders?.deliver ?? 'ENTREGAR'}
              </button>
            {:else}
              {#if status === 'pending'}
                <button
                  type="button"
                  onclick={(e) => { e.stopPropagation(); cycleOrderStatus(order.order_number, stationFilter); }}
                  class="w-full py-3 px-4 rounded-xl text-base font-bold text-white shadow-md transition-colors {useBarColor ? 'bg-blue-600 hover:bg-blue-700' : 'bg-orange-500 hover:bg-orange-600'}"
                >
                  <span aria-hidden="true">🔥</span> {t.orders?.startPreparing ?? 'Iniciar preparación'}
                </button>
              {:else if status === 'preparing'}
                <button
                  type="button"
                  onclick={(e) => { e.stopPropagation(); cycleOrderStatus(order.order_number, stationFilter); }}
                  class="w-full py-3 px-4 rounded-xl text-base font-bold text-white shadow-md transition-colors {useBarColor ? 'bg-blue-500 hover:bg-blue-600' : 'bg-amber-500 hover:bg-amber-600'}"
                >
                  ✓ {t.orders?.markAsReady ?? 'LISTO'}
                </button>
              {:else}
                <div class="w-full py-3 px-4 rounded-xl text-base font-bold bg-green-100 text-green-800 text-center">
                  ✓ {t.orders?.statusDone ?? 'LISTO'}
                </div>
              {/if}
            {/if}
          </div>
        </li>
      {/snippet}
      <!-- Toggle vista: Vertical (tabs + lista) vs 3 columnas (Kanban), solo en Cocina/Bar -->
      {#if !showQRView && (stationFilter === 'KITCHEN' || stationFilter === 'BAR')}
        <div class="flex items-center gap-2 mb-3">
          <span class="text-sm font-medium text-gray-600">{t.orders?.viewVertical ?? 'Vista'}:</span>
          <div class="flex rounded-lg overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="group" aria-label="{t.orders?.viewVertical ?? 'Vertical'} / {t.orders?.viewThreeColumns ?? '3 columnas'}">
            <button
              type="button"
              onclick={() => ordersViewMode = 'vertical'}
              class="px-3 py-2 text-sm font-semibold transition-colors {ordersViewMode === 'vertical' ? 'bg-white text-gray-900 shadow-sm border border-gray-200' : 'text-gray-600 hover:bg-gray-200'}"
            >
              {t.orders?.viewVertical ?? 'Vertical'}
            </button>
            <button
              type="button"
              onclick={() => ordersViewMode = 'kanban'}
              class="px-3 py-2 text-sm font-semibold transition-colors border-l border-gray-200 {ordersViewMode === 'kanban' ? 'bg-white text-gray-900 shadow-sm border border-gray-200' : 'text-gray-600 hover:bg-gray-200'}"
            >
              {t.orders?.viewThreeColumns ?? '3 columnas'}
            </button>
          </div>
        </div>
      {/if}
      <!-- Vista vertical: tabs + una sola lista -->
      {#if ordersViewMode === 'vertical' && !showQRView && (stationFilter === 'KITCHEN' || stationFilter === 'BAR')}
        <div class="flex items-stretch mb-4 w-full rounded-xl overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="tablist" aria-label="{t.orders?.tabPending ?? 'Pendientes'}, {t.orders?.tabPreparing ?? 'En preparación'}, {t.orders?.tabDone ?? 'Listos'}">
          <button
            type="button"
            role="tab"
            aria-selected={operationalTab === 'pending'}
            onclick={() => operationalTab = 'pending'}
            class="flex-1 min-w-0 px-3 py-3.5 text-base font-semibold transition-colors border-r border-gray-200 {operationalTab === 'pending' ? 'bg-amber-200 text-amber-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            🟠 {t.orders?.tabPending ?? 'Pendientes'} {ordersByTab.pending.length > 0 ? `(${ordersByTab.pending.length})` : ''}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={operationalTab === 'preparing'}
            onclick={() => operationalTab = 'preparing'}
            class="flex-1 min-w-0 px-3 py-3.5 text-base font-semibold transition-colors border-r border-gray-200 {operationalTab === 'preparing' ? 'bg-blue-100 text-blue-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            🔵 {t.orders?.tabPreparing ?? 'En preparación'} {ordersByTab.preparing.length > 0 ? `(${ordersByTab.preparing.length})` : ''}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={operationalTab === 'done'}
            onclick={() => operationalTab = 'done'}
            class="flex-1 min-w-0 px-3 py-3.5 text-base font-semibold transition-colors {operationalTab === 'done' ? 'bg-green-100 text-green-800 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            🟢 {t.orders?.tabDone ?? 'Listos'} {ordersByTab.done.length > 0 ? `(${ordersByTab.done.length})` : ''}
          </button>
        </div>
      {/if}
      {#if stationFilter !== 'ALL' && ordersViewMode === 'kanban'}
        <!-- Vista Kanban: 3 columnas clickeables (Pendientes | En preparación | Listos) -->
        {@const kanbanColumns = [{ key: 'pending', label: t.orders?.tabPending ?? 'Pendientes', orders: ordersByTab.pending, icon: '🟠', bg: 'bg-amber-50 border-amber-200', headerBg: 'bg-amber-200 text-amber-900' }, { key: 'preparing', label: t.orders?.tabPreparing ?? 'En preparación', orders: ordersByTab.preparing, icon: '🔵', bg: 'bg-blue-50/80 border-blue-200', headerBg: 'bg-blue-100 text-blue-900' }, { key: 'done', label: t.orders?.tabDone ?? 'Listos', orders: ordersByTab.done, icon: '🟢', bg: 'bg-green-50/80 border-green-200', headerBg: 'bg-green-100 text-green-800' }]}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
          {#each kanbanColumns as col}
            <div class="flex flex-col rounded-xl border-2 {col.bg} overflow-hidden min-h-[200px]">
              <button
                type="button"
                onclick={() => { operationalTab = col.key as 'pending' | 'preparing' | 'done'; ordersViewMode = 'vertical'; }}
                class="flex items-center justify-center gap-2 w-full px-4 py-3 font-bold text-left transition-colors hover:opacity-90 {col.headerBg}"
                title="{t.orders?.viewVertical ?? 'Ver'} {col.label}"
              >
                <span aria-hidden="true">{col.icon}</span>
                <span>{col.label}</span>
                <span class="tabular-nums">({col.orders.length})</span>
              </button>
              <ul class="flex-1 overflow-y-auto p-3 space-y-3">
                {#each col.orders as order, index (order.order_number)}
                  {@render orderCard(order, col.key as 'pending' | 'preparing' | 'done', index === 0)}
                {/each}
              </ul>
            </div>
          {/each}
        </div>
      {:else if ordersToShow.length === 0}
        <div class="rounded-lg bg-gray-100 border border-gray-200 p-8 text-center text-gray-600">
          {stationFilter === 'ALL' ? (t.orders?.empty ?? 'No hay órdenes aún.') : (t.orders?.emptyForStation ?? 'No hay órdenes para esta estación.')}
        </div>
      {:else}
      <ul class="space-y-5 kitchen-orders-list" class:kitchen-mode-list={kitchenMode}>
        {#each ordersToShow as order, index (order.order_number)}
          {@const cardStation = stationFilter === 'ALL' ? null : stationFilter}
          {@const status = cardStation !== null ? getOrderStatus(order.order_number, cardStation) : (isOrderReadyForDelivery(order) ? 'done' : getCajaOrderStatusLabel(order))}
          {@render orderCard(order, status, index === 0)}
        {/each}
      </ul>
      {/if}
    {/if}
  </div>

  <!-- Modal QR ampliado (clic en el código para agrandar) -->
  {#if qrEnlarged && menuId && session?.access_token && typeof window !== 'undefined'}
    {@const baseUrl = window.location.origin}
    {@const modalHashParams = session?.refresh_token ? `token=${encodeURIComponent(session.access_token)}&refresh_token=${encodeURIComponent(session.refresh_token)}` : `token=${encodeURIComponent(session.access_token)}`}
    {@const enlargedUrl = `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=${qrEnlarged}#${modalHashParams}`}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      role="dialog"
      aria-modal="true"
      aria-labelledby="qr-enlarged-title"
      tabindex="-1"
      onclick={() => qrEnlarged = null}
      onkeydown={(e) => e.key === 'Escape' && (qrEnlarged = null)}
    >
      <div
        class="bg-white rounded-xl shadow-xl p-6 max-w-sm w-full"
        role="document"
        tabindex="0"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
      >
        <h2 id="qr-enlarged-title" class="text-lg font-bold text-gray-800 mb-4">
          {qrEnlarged === 'ALL' ? (t.orders?.filterAll ?? 'Entrega') : qrEnlarged === 'KITCHEN' ? (t.orders?.filterKitchen ?? 'Cocina') : (t.orders?.filterBar ?? 'Barra')}
        </h2>
        <div class="flex justify-center mb-4 bg-white rounded-lg border border-gray-200 p-3">
          <img
            src={`https://api.qrserver.com/v1/create-qr-code/?size=280x280&data=${encodeURIComponent(enlargedUrl)}`}
            alt="QR ampliado"
            class="w-64 h-64 sm:w-72 sm:h-72 object-contain"
          />
        </div>
        <p class="text-xs text-gray-500 mb-3">Para un QR nuevo, actualiza la página o vuelve a Órdenes.</p>
        <div class="flex gap-2">
          <button
            type="button"
            onclick={() => enlargedUrl && navigator.clipboard.writeText(enlargedUrl).then(() => alert('Enlace copiado'))}
            class="flex-1 py-2 rounded-lg bg-gray-800 text-white text-sm font-medium hover:bg-gray-900"
          >
            Copiar enlace
          </button>
          <button type="button" onclick={() => qrEnlarged = null} class="py-2 px-4 rounded-lg border border-gray-300 text-gray-700 text-sm font-medium hover:bg-gray-50">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  {/if}
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
  /* Tab Listos: tarjetas atenuadas, no compiten con el foco principal */
  :global(.order-card-done) {
    opacity: 0.85;
    border-color: rgb(209 213 219);
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
