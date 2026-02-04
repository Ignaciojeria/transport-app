<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getKitchenOrdersFromProjection, subscribeMenuOrdersRealtime, getMenuOrderByNumber, groupOrderItemsForDisplay, type KitchenOrder, type MenuOrderRow, type StationFilter } from '../menuUtils'
  import { startPreparation, markReady, dispatchOrder, cancelOrder } from '../orderApi'
  import { getActiveJourney, createJourney } from '../journeyApi'
  import { t as tStore } from '../useLanguage'
  import { playNewOrderSound, ensureAudioUnlocked } from '../utils/newOrderSound'

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
  let kitchenMode = $state(false)
  /** Filtro por estación: se aplica en Supabase (columna station). */
  let stationFilter = $state<StationFilter>('ALL')
  /** Tab operativo en Cocina/Bar: Pendientes | En preparación | Listos (solo cuando stationFilter es KITCHEN o BAR). */
  let operationalTab = $state<'pending' | 'preparing' | 'done'>('pending')
  /** Tab en vista Entrega: Por entregar | Entregado | Cancelado. */
  let deliveryTab = $state<'pending' | 'delivered' | 'cancelled'>('pending')
  /** Acción en curso: clave `${aggregateId}-${action}` para deshabilitar botones. */
  let actionInProgress = $state<Set<string>>(new Set())
  /** Mensaje de error de la última acción (ej. fallo al llamar al backend). */
  let actionError = $state<string | null>(null)
  /** Vista QR: muestra códigos Cocina/Barra en lugar de la lista de órdenes. */
  let showQRView = $state(false)
  /** Vista órdenes: vertical (tabs + lista) o kanban (3 columnas). Aplica a Cocina/Bar y a Entrega. */
  let ordersViewMode = $state<'vertical' | 'kanban'>('vertical')
  let realtimeUnsubscribe = $state<(() => void) | null>(null)
  /** Modal QR ampliado: KITCHEN | BAR | ALL | null */
  let qrEnlarged = $state<'KITCHEN' | 'BAR' | 'ALL' | null>(null)
  /** Modal cancelar: orden seleccionada, motivo y comentario opcional. */
  let orderToCancel = $state<KitchenOrder | null>(null)
  let cancelReason = $state<string>('')
  let cancelComment = $state<string>('')
  let cancelInProgress = $state(false)
  let previousOrderNumbers = new Set<number>()
  let initialLoadDone = false
  let activeJourney = $state<{ id: string } | null>(null)
  let createJourneyInProgress = $state(false)
  let createJourneyError = $state<string | null>(null)

  const CANCEL_REASON_KEYS = ['outOfStock', 'orderError', 'customerLeft', 'paymentIssue', 'other'] as const

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const session = $derived(authState.session)
  const t = $derived($tStore)

  /** Devuelve el menuId cargado (para suscribir Realtime justo después). */
  async function loadOrders(): Promise<string | null> {
    if (!userId || !session?.access_token) {
      error = t.orders?.noSession ?? 'No Hay Sesión Activa'
      loading = false
      return null
    }

    try {
      loading = true
      error = null
      createJourneyError = null
      const currentMenuId = await getLatestMenuId(userId, session.access_token)
      if (!currentMenuId) {
        error = t.orders?.noMenu ?? 'No Se Encontró Un Menú'
        loading = false
        return null
      }
      menuId = currentMenuId
      const journey = await getActiveJourney(currentMenuId, session.access_token)
      activeJourney = journey
      const newOrders = journey
        ? await getKitchenOrdersFromProjection(currentMenuId, session.access_token, stationFilter, journey.id)
        : []
      orders = newOrders
      const newIds = new Set(newOrders.map((o) => Number(o.order_number)))
      if (initialLoadDone) {
        const addedIds = [...newIds].filter((id) => !previousOrderNumbers.has(id))
        const justArrived =
          addedIds.length > 0 &&
          newOrders.some((o) => {
            if (!addedIds.includes(Number(o.order_number))) return false
            const created = new Date(o.created_at).getTime()
            return Date.now() - created < 90_000
          })
        if (justArrived) playNewOrderSound()
      }
      newIds.forEach((id) => previousOrderNumbers.add(id))
      initialLoadDone = true
      return currentMenuId
    } catch (err: unknown) {
      console.error('Error cargando órdenes:', err)
      error = err instanceof Error ? err.message : 'Error al cargar las órdenes'
      return null
    } finally {
      loading = false
    }
  }

  async function openJourney() {
    if (!session?.access_token || !menuId) return
    createJourneyInProgress = true
    createJourneyError = null
    try {
      await createJourney(menuId, session.access_token, 'USER', t.jornada?.openJourneyReason ?? 'Apertura manual')
      await loadOrders()
    } catch (e) {
      console.error('Error creando jornada:', e)
      createJourneyError = t.jornada?.errorCreatingJourney ?? 'Error al abrir la jornada. Intenta de nuevo.'
    } finally {
      createJourneyInProgress = false
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

  /** Total de unidades de ítems. */
  function getItemCount(items: KitchenOrder['items']): number {
    return items.reduce((s, i) => s + i.quantity, 0)
  }

  /** Órdenes mostradas = las que vienen de Supabase ya filtradas por station. Excluimos órdenes totalmente canceladas. */
  const displayedOrders = $derived(
    orders.filter((o) => {
      const allCancelled = o.items.every((i) => i.status === 'CANCELLED')
      return !allCancelled
    })
  )

  /** En Cocina/Bar: órdenes particionadas por estado (pending / preparing / done) según status en proyección. */
  const ordersByTab = $derived.by(() => {
    if (stationFilter !== 'KITCHEN' && stationFilter !== 'BAR') return { pending: [] as KitchenOrder[], preparing: [] as KitchenOrder[], done: [] as KitchenOrder[] }
    const pending: KitchenOrder[] = []
    const preparing: KitchenOrder[] = []
    const done: KitchenOrder[] = []
    for (const o of displayedOrders) {
      const st = getOrderStatusFromItems(o, stationFilter)
      if (st === 'pending') pending.push(o)
      else if (st === 'preparing') preparing.push(o)
      else done.push(o)
    }
    return { pending, preparing, done }
  })

  /** En vista Entrega: órdenes por tab (por entregar | entregado | cancelado). Usa lista completa para incluir canceladas. */
  const ordersByDeliveryTab = $derived.by(() => {
    if (stationFilter !== 'ALL') return { pending: [] as KitchenOrder[], delivered: [] as KitchenOrder[], cancelled: [] as KitchenOrder[] }
    const pending = orders.filter((o) => !isOrderFullyDelivered(o) && !isOrderFullyCancelled(o))
    const delivered = orders.filter((o) => isOrderFullyDelivered(o))
    const cancelled = orders.filter((o) => isOrderFullyCancelled(o))
    return { pending, delivered, cancelled }
  })

  /** Órdenes a mostrar: en Entrega = las del tab activo (por entregar/entregado/cancelado); en Cocina/Bar = las del tab operativo. */
  const ordersToShow = $derived(
    stationFilter === 'ALL' ? ordersByDeliveryTab[deliveryTab] : ordersByTab[operationalTab]
  )

  async function markOrderAsDelivered(order: KitchenOrder) {
    if (!menuId || !session?.access_token) return
    const key = `${order.aggregate_id}-dispatch`
    if (actionInProgress.has(key)) return
    actionError = null
    actionInProgress = new Set(actionInProgress).add(key)
    try {
      await dispatchOrder(menuId, order.aggregate_id, session.access_token)
      await loadOrders()
    } catch (err) {
      actionError = err instanceof Error ? err.message : 'Error al marcar como entregado'
    } finally {
      const next = new Set(actionInProgress)
      next.delete(key)
      actionInProgress = next
    }
  }

  async function setStationFilterAndReload(filter: StationFilter) {
    stationFilter = filter
    if (filter === 'KITCHEN' || filter === 'BAR') operationalTab = 'pending'
    if (filter === 'ALL') deliveryTab = 'pending'
    if (menuId && session?.access_token && activeJourney) {
      loading = true
      try {
        orders = await getKitchenOrdersFromProjection(menuId, session.access_token, filter, activeJourney.id)
        previousOrderNumbers = new Set(orders.map((o) => Number(o.order_number)))
      } finally {
        loading = false
      }
    }
  }

  /** Estado por estación derivado de los ítems de la orden (proyección). Ignora ítems CANCELLED. */
  function getOrderStatusFromItems(order: KitchenOrder, station: string): 'pending' | 'preparing' | 'done' {
    const stationItems = order.items.filter((i) => (i.station ?? 'KITCHEN') === station && i.status !== 'CANCELLED')
    if (stationItems.length === 0) return 'done'
    if (stationItems.some((i) => i.status === 'PENDING')) return 'pending'
    if (stationItems.some((i) => i.status === 'IN_PROGRESS')) return 'preparing'
    return 'done'
  }

  /** Orden terminada: todos los ítems activos en estado terminal (DISPATCHED = retiro, DELIVERED = despacho). */
  function isOrderFullyDelivered(order: KitchenOrder): boolean {
    const active = order.items.filter((i) => i.status !== 'CANCELLED')
    return active.length > 0 && active.every((i) => i.status === 'DISPATCHED' || i.status === 'DELIVERED')
  }

  function isOrderFullyCancelled(order: KitchenOrder): boolean {
    return order.items.length > 0 && order.items.every((i) => i.status === 'CANCELLED')
  }

  async function cycleOrderStatus(order: KitchenOrder, station: string) {
    if (!menuId || !session?.access_token) return
    const current = getOrderStatusFromItems(order, station)
    const key = `${order.aggregate_id}-${current}-${station}`
    if (actionInProgress.has(key)) return
    actionError = null
    actionInProgress = new Set(actionInProgress).add(key)
    try {
      if (current === 'pending') {
        await startPreparation(menuId, order.aggregate_id, session.access_token, station, [])
      } else if (current === 'preparing') {
        const itemKeys = order.items
          .filter((i) => (i.station ?? 'KITCHEN') === station && i.status === 'IN_PROGRESS')
          .map((i) => i.item_key)
        if (itemKeys.length > 0) {
          await markReady(menuId, order.aggregate_id, session.access_token, station, itemKeys)
        }
      }
      await loadOrders()
    } catch (err) {
      actionError = err instanceof Error ? err.message : 'Error al actualizar estado'
    } finally {
      const next = new Set(actionInProgress)
      next.delete(key)
      actionInProgress = next
    }
  }

  /** Estaciones que tienen ítems en esta orden (KITCHEN, BAR). */
  function getStationsInOrder(order: KitchenOrder): string[] {
    const stations = new Set<string>()
    for (const i of order.items) {
      if (i.station) stations.add(i.station)
    }
    return stations.size > 0 ? [...stations] : ['KITCHEN']
  }

  /** Indica si Cocina y Barra marcaron la orden como lista (todas estaciones en READY, DISPATCHED o DELIVERED). */
  function isOrderReadyForDelivery(order: KitchenOrder): boolean {
    const stations = getStationsInOrder(order)
    return stations.every((st) => getOrderStatusFromItems(order, st) === 'done')
  }

  /** En Caja: estado resumido para mostrar (Pendiente / En preparación). */
  function getCajaOrderStatusLabel(order: KitchenOrder): 'pending' | 'preparing' {
    const stations = getStationsInOrder(order)
    const statuses = stations.map((st) => getOrderStatusFromItems(order, st))
    if (statuses.every((s) => s === 'pending')) return 'pending'
    return 'preparing'
  }

  function isActionInProgress(order: KitchenOrder, station: string, action: 'pending' | 'preparing'): boolean {
    const key = `${order.aggregate_id}-${action}-${station}`
    return actionInProgress.has(key)
  }

  function isDispatchInProgress(order: KitchenOrder): boolean {
    return actionInProgress.has(`${order.aggregate_id}-dispatch`)
  }

  function openCancelModal(order: KitchenOrder) {
    orderToCancel = order
    cancelReason = ''
    cancelComment = ''
  }

  function closeCancelModal() {
    if (!cancelInProgress) {
      orderToCancel = null
      cancelReason = ''
      cancelComment = ''
    }
  }

  function getCancelReasonLabel(key: string): string {
    return t.orders?.cancelReasons?.[key] ?? key
  }

  async function confirmCancelOrder() {
    if (!orderToCancel || !menuId || !session?.access_token || !cancelReason.trim()) return
    cancelInProgress = true
    actionError = null
    try {
      const reasonLabel = getCancelReasonLabel(cancelReason)
      const reasonText = cancelComment.trim()
        ? `${reasonLabel}: ${cancelComment.trim()}`
        : reasonLabel
      await cancelOrder(menuId, orderToCancel.aggregate_id, session.access_token, reasonText)
      orderToCancel = null
      cancelReason = ''
      cancelComment = ''
      await loadOrders()
    } catch (err) {
      actionError = err instanceof Error ? err.message : 'Error al cancelar la orden'
    } finally {
      cancelInProgress = false
    }
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

<div class="h-full flex flex-col bg-gray-50 kitchen-orders-root text-sm" class:kitchen-mode={kitchenMode}>
  <!-- Header -->
  <div class="flex-shrink-0 px-3 sm:px-4 py-2.5 border-b border-gray-200 bg-white" class:kitchen-mode-header={kitchenMode}>
    <div class="flex items-center justify-between gap-3">
      {#if !kitchenMode}
        <button
          type="button"
          onclick={onMenuClick}
          class="md:hidden p-1.5 -ml-1.5 rounded-lg hover:bg-gray-100 text-gray-600"
          aria-label="Abrir menú"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      {/if}
      <h1 class="text-lg sm:text-xl font-bold text-gray-800" class:text-2xl={kitchenMode}>
        {t.sidebar.orders}
      </h1>
      <div class="flex items-center gap-1.5">
        <button
          type="button"
          onclick={toggleKitchenMode}
          class="rounded-lg px-3 py-1.5 text-xs font-semibold {kitchenMode ? 'bg-amber-500 text-white hover:bg-amber-600' : 'bg-gray-200 text-gray-800 hover:bg-gray-300'}"
        >
          {kitchenMode ? (t.orders?.exitKitchenMode ?? 'Salir Modo Full') : (t.orders?.kitchenMode ?? 'Modo Full')}
        </button>
      </div>
    </div>
    {#if !kitchenMode}
      <p class="mt-0.5 text-xs text-gray-500">
        {t.orders?.subtitle ?? 'Ordenado por hora comprometida. Vista orientada a cocina.'}
      </p>
      <!-- Filtro por estación + opción QR: oculto en modo full -->
      <div class="mt-2 flex flex-wrap items-center gap-1.5">
        <button
          type="button"
          onclick={() => { ensureAudioUnlocked(); showQRView = false; setStationFilterAndReload('ALL'); }}
          class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {!showQRView && stationFilter === 'ALL' ? 'bg-gray-800 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
        >
          {t.orders?.filterAll ?? 'Entrega'}
        </button>
        <button
          type="button"
          onclick={() => { ensureAudioUnlocked(); showQRView = false; setStationFilterAndReload('KITCHEN'); }}
          class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {!showQRView && stationFilter === 'KITCHEN' ? 'bg-amber-600 text-white' : 'bg-amber-50 text-amber-800 hover:bg-amber-100 border border-amber-200'}"
        >
          {t.orders?.filterKitchen ?? 'Cocina'}
        </button>
        <button
          type="button"
          onclick={() => { ensureAudioUnlocked(); showQRView = false; setStationFilterAndReload('BAR'); }}
          class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {!showQRView && stationFilter === 'BAR' ? 'bg-blue-600 text-white' : 'bg-blue-50 text-blue-800 hover:bg-blue-100 border border-blue-200'}"
        >
          {t.orders?.filterBar ?? 'Barra'}
        </button>
        {#if menuId && session?.access_token}
          <button
            type="button"
            onclick={() => showQRView = true}
            class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {showQRView ? 'bg-gray-700 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200 border border-gray-300'}"
          >
            {t.orders?.showQR ?? 'QR'}
          </button>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Content -->
  <div class="flex-1 overflow-y-auto px-3 sm:px-4 py-3">
    {#if showQRView && menuId && session?.access_token}
      <!-- Vista QR: códigos Entrega, Cocina y Barra sin tapar filtros -->
      {@const baseUrl = typeof window !== 'undefined' ? window.location.origin : ''}
      {@const hashParams = session?.refresh_token ? `token=${encodeURIComponent(session.access_token)}&refresh_token=${encodeURIComponent(session.refresh_token)}` : `token=${encodeURIComponent(session.access_token)}`}
      {@const urlCaja = baseUrl ? `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=ALL#${hashParams}` : ''}
      {@const urlCocina = baseUrl ? `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=KITCHEN#${hashParams}` : ''}
      {@const urlBarra = baseUrl ? `${baseUrl}/?view=station&menu_id=${encodeURIComponent(menuId)}&station=BAR#${hashParams}` : ''}
      <div class="max-w-3xl">
        <p class="text-xs font-medium text-gray-700 mb-1">Acceso sin login (entrega, cocinero o barista escanea el código)</p>
        <p class="text-[11px] text-gray-500 mb-3">Haz clic en el código para agrandarlo.</p>
        <div class="flex flex-wrap gap-4">
          <div class="flex items-start gap-2 p-3 rounded-lg bg-gray-100 border border-gray-300">
            <button
              type="button"
              onclick={() => qrEnlarged = 'ALL'}
              class="flex-shrink-0 w-24 h-24 bg-white rounded border border-gray-300 overflow-hidden cursor-pointer hover:ring-2 hover:ring-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
              title="Clic para agrandar"
            >
              {#if urlCaja}
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=96x96&data=${encodeURIComponent(urlCaja)}`} alt="QR Entrega" class="w-full h-full object-contain pointer-events-none" />
              {/if}
            </button>
            <div class="min-w-0">
              <p class="text-xs font-semibold text-gray-900">{t.orders?.filterAll ?? 'Entrega'}</p>
              <p class="text-[11px] text-gray-700 mb-0.5">Ver todas las órdenes en tiempo real</p>
              <button type="button" onclick={() => urlCaja && navigator.clipboard.writeText(urlCaja).then(() => alert('Enlace copiado'))} class="text-[11px] text-gray-600 underline hover:no-underline">Copiar enlace</button>
            </div>
          </div>
          <div class="flex items-start gap-2 p-3 rounded-lg bg-amber-50 border border-amber-200">
            <button
              type="button"
              onclick={() => qrEnlarged = 'KITCHEN'}
              class="flex-shrink-0 w-24 h-24 bg-white rounded border border-amber-200 overflow-hidden cursor-pointer hover:ring-2 hover:ring-amber-400 focus:outline-none focus:ring-2 focus:ring-amber-500"
              title="Clic para agrandar"
            >
              {#if urlCocina}
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=96x96&data=${encodeURIComponent(urlCocina)}`} alt="QR Cocina" class="w-full h-full object-contain pointer-events-none" />
              {/if}
            </button>
            <div class="min-w-0">
              <p class="text-xs font-semibold text-amber-900">{t.orders?.filterKitchen ?? 'Cocina'}</p>
              <p class="text-[11px] text-amber-800 mb-0.5">Escanear → pedir nombre → ver pedidos en tiempo real</p>
              <button type="button" onclick={() => urlCocina && navigator.clipboard.writeText(urlCocina).then(() => alert('Enlace copiado'))} class="text-[11px] text-amber-700 underline hover:no-underline">Copiar enlace</button>
            </div>
          </div>
          <div class="flex items-start gap-2 p-3 rounded-lg bg-blue-50 border border-blue-200">
            <button
              type="button"
              onclick={() => qrEnlarged = 'BAR'}
              class="flex-shrink-0 w-24 h-24 bg-white rounded border border-blue-200 overflow-hidden cursor-pointer hover:ring-2 hover:ring-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              title="Clic para agrandar"
            >
              {#if urlBarra}
                <img src={`https://api.qrserver.com/v1/create-qr-code/?size=96x96&data=${encodeURIComponent(urlBarra)}`} alt="QR Barra" class="w-full h-full object-contain pointer-events-none" />
              {/if}
            </button>
            <div class="min-w-0">
              <p class="text-xs font-semibold text-blue-900">{t.orders?.filterBar ?? 'Barra'}</p>
              <p class="text-[11px] text-blue-800 mb-0.5">Escanear → pedir nombre → ver pedidos en tiempo real</p>
              <button type="button" onclick={() => urlBarra && navigator.clipboard.writeText(urlBarra).then(() => alert('Enlace copiado'))} class="text-[11px] text-blue-700 underline hover:no-underline">Copiar enlace</button>
            </div>
          </div>
        </div>
      </div>
    {:else if loading}
      <div class="flex items-center justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-2 border-blue-600 border-t-transparent"></div>
      </div>
      {:else if error}
      <div class="rounded-lg bg-red-50 border border-red-200 p-3 text-xs text-red-800">
        {error}
      </div>
      {:else if menuId && activeJourney === null}
      <section class="max-w-md mx-auto mt-8">
        <div class="rounded-xl bg-white border border-gray-200 shadow-sm p-8 text-center">
          <p class="text-gray-600 mb-6">
            {t.jornada?.noActiveJourney ?? 'No tienes una jornada abierta. Abre una para comenzar a registrar órdenes del día.'}
          </p>
          {#if createJourneyError}
            <p class="text-sm text-red-600 mb-4">{createJourneyError}</p>
          {/if}
          <button
            type="button"
            disabled={createJourneyInProgress}
            onclick={openJourney}
            class="w-full flex items-center justify-center gap-2 py-4 px-6 rounded-xl font-semibold text-white bg-blue-600 hover:bg-blue-700 focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-60 transition-colors"
          >
            {#if createJourneyInProgress}
              <span class="animate-spin inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full"></span>
              <span>{t.jornada?.openingJourney ?? 'Abriendo jornada...'}</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              <span>{t.jornada?.openJourney ?? 'Abrir jornada'}</span>
            {/if}
          </button>
        </div>
      </section>
      {:else if actionError}
      <div class="rounded-lg bg-amber-50 border border-amber-200 p-3 text-xs text-amber-800 mb-3">
        {actionError}
        <button type="button" onclick={() => (actionError = null)} class="ml-2 underline hover:no-underline">Cerrar</button>
      </div>
    {:else}
      <!-- Snippet: una card de orden (Cocina/Bar/Caja). compactStatus=true en Kanban; isDelivered=true en tab Entregado (Entrega). -->
      {#snippet orderCard(order: KitchenOrder, status: 'pending' | 'preparing' | 'done', isFirst: boolean, compactStatus: boolean = false, isDelivered: boolean = false, isCancelled: boolean = false)}
        {@const type = order.fulfillment}
        {@const itemCount = getItemCount(order.items)}
        {@const remainingMin = getRemainingMinutes(order.requested_time)}
        {@const timeColor = getRemainingTimeColor(remainingMin)}
        {@const isBarOrder = order.items.some((i) => i.station === 'BAR')}
        {@const useBarColor = stationFilter === 'BAR' || (stationFilter === 'ALL' && isBarOrder)}
        {@const readyForDelivery = isOrderReadyForDelivery(order)}
        {@const isDoneTab = status === 'done'}
        {@const barStatusForOrder = getOrderStatusFromItems(order, 'BAR')}
        {@const hasBarItems = order.items.some((i) => i.station === 'BAR')}
        {@const statusIcon = status === 'pending' ? '🟠' : status === 'preparing' ? '⏳' : '✓'}
        {@const statusTitle = status === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : status === 'preparing' ? (t.orders?.statusPreparing ?? 'En Preparación') : (t.orders?.statusDone ?? 'Listo')}
        <li class="bg-white rounded-lg border overflow-hidden kitchen-order-card order-card {isDoneTab ? 'order-card-done' : ''} {isFirst && !isDoneTab ? 'kitchen-order-first border-amber-400 shadow-md' : 'border-gray-200 shadow-sm'}">
          <div class="w-full px-3 py-2 sm:px-4 flex flex-wrap items-center gap-2 border-b border-gray-100 {isFirst && !isDoneTab ? 'sm:py-3' : ''} {compactStatus ? 'py-1.5 sm:py-2' : ''}">
            <span class="font-bold text-gray-900 tabular-nums {isFirst && !isDoneTab && !compactStatus ? 'text-2xl sm:text-3xl' : compactStatus ? 'text-lg sm:text-xl' : 'text-xl sm:text-2xl'}">#{order.order_number}</span>
            {#if compactStatus && stationFilter !== 'ALL'}
              <span class="text-sm opacity-90" aria-hidden="true" title={statusTitle}>{statusIcon}</span>
            {/if}
            {#if stationFilter === 'KITCHEN' && hasBarItems}
              <span class="text-sm" aria-hidden="true" title="{barStatusForOrder === 'done' ? (t.orders?.statusDone ?? 'Bar Listo') : (t.orders?.statusPreparing ?? 'Bar En Preparación')}">{barStatusForOrder === 'done' ? '🍺 ✔️' : '🍺 ⏳'}</span>
            {/if}
            <span class="font-semibold text-gray-700 {isFirst && !compactStatus ? 'text-base sm:text-lg' : compactStatus ? 'text-sm' : 'text-sm sm:text-base'}">
              {(t.orders?.forTime ?? 'Para')} {formatRequestedTime(order.requested_time)}
            </span>
            {#if remainingMin !== null}
              <span class="inline-flex items-center gap-0.5 rounded-full px-2 py-0.5 text-xs font-bold tabular-nums
                {timeColor === 'green' ? 'bg-green-100 text-green-800' : ''}
                {timeColor === 'yellow' ? 'bg-amber-200 text-amber-900' : ''}
                {timeColor === 'red' ? 'bg-red-100 text-red-800' : ''}">
                <span aria-hidden="true">{timeColor === 'green' ? '🟢' : timeColor === 'yellow' ? '🟡' : '🔴'}</span>
                {getRemainingTimeLabel(remainingMin)}
              </span>
            {/if}
            <span class="inline-flex items-center rounded-full font-medium text-xs {type === 'DELIVERY' ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'} {isFirst && !compactStatus ? 'px-2.5 py-1' : 'px-2 py-0.5'}">
              {getFulfillmentLabel(type)}
            </span>
            {#if stationFilter === 'ALL'}
              {@const kitchenSt = getOrderStatusFromItems(order, 'KITCHEN')}
              {@const barSt = getOrderStatusFromItems(order, 'BAR')}
              {@const orderHasBar = order.items.some((i) => i.station === 'BAR')}
              <div class="flex flex-wrap items-center gap-2 text-xs font-semibold">
                <span class="inline-flex items-center gap-0.5">{t.orders?.filterKitchen ?? 'Cocina'}: {kitchenSt === 'done' ? '✔️' : '⏳'}</span>
                <span class="inline-flex items-center gap-0.5">{t.orders?.filterBar ?? 'Barra'}: {orderHasBar ? (barSt === 'done' ? '✔️' : '⏳') : '—'}</span>
                <span class="inline-flex items-center gap-0.5 rounded-full border px-1.5 py-0.5 {readyForDelivery ? 'bg-green-50 text-green-800 border-green-200' : 'bg-amber-50 text-amber-900 border-amber-200'}">
                  {t.orders?.statusGeneralLabel ?? 'Estado General'}: {readyForDelivery ? (t.orders?.readyToDeliver ?? 'Listo Para Entregar') : (t.orders?.statusPreparing ?? 'En Preparación')}
                </span>
              </div>
            {:else if !compactStatus}
              <span class="inline-flex items-center gap-0.5 rounded-full font-bold border text-xs {isFirst ? 'px-2.5 py-1' : 'px-2 py-0.5'}
                {status === 'pending' ? 'bg-gray-100 text-gray-700 border-gray-200' : ''}
                {status === 'preparing' ? 'bg-amber-50 text-amber-900 border-amber-200' : ''}
                {status === 'done' ? 'bg-green-50 text-green-800 border-green-200' : ''}">
                {#if status === 'preparing'}<span aria-hidden="true">⏳</span>{/if}
                {status === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : status === 'preparing' ? (t.orders?.statusPreparing ?? 'En Preparación') : (t.orders?.statusDone ?? 'Listo')}
              </span>
            {/if}
          </div>
          <div class="px-3 py-2 sm:px-4 bg-amber-50/50 border-b border-amber-100 {isFirst && !compactStatus ? 'py-3' : ''} {compactStatus ? 'py-1.5 sm:py-2' : ''}">
            <div class="flex items-center justify-between gap-1.5 mb-1">
              <p class="font-semibold text-amber-800 uppercase tracking-wide text-[11px]">{t.orders?.itemsToPrepare ?? 'Qué Preparar'}</p>
              <span class="text-xs font-bold text-amber-800 tabular-nums">{(t.orders?.itemsCount ?? '{count} ítems').replace('{count}', String(itemCount))}</span>
            </div>
            <ul class="space-y-0.5 text-gray-900 {compactStatus ? 'text-sm' : isFirst ? 'text-base sm:text-lg' : 'text-sm sm:text-base'}">
              {#each groupOrderItemsForDisplay(order.items) as item}
                <li class="tabular-nums text-inherit">
                  <span class="font-bold text-amber-800">{item.quantity}×</span> <span class="font-normal">{item.item_name}</span>
                </li>
              {/each}
            </ul>
          </div>
          <div class="px-3 py-2 sm:px-4 border-t border-gray-100 {compactStatus ? 'py-1.5 sm:py-2' : ''}">
            {#if stationFilter === 'ALL'}
              {#if isCancelled}
                <div class="w-full py-2 px-3 rounded-lg text-xs font-bold bg-gray-200 text-gray-700 text-center">
                  ✕ {t.orders?.cancelled ?? 'Cancelado'}
                </div>
              {:else if isDelivered}
                <div class="w-full py-2 px-3 rounded-lg text-xs font-bold bg-green-100 text-green-800 text-center">
                  ✓ {type === 'DELIVERY' ? (t.orders?.delivered ?? 'Entregado') : (t.orders?.dispatched ?? 'Despachado')}
                </div>
              {:else}
                <div class="flex flex-col gap-1.5">
                  <button
                    type="button"
                    disabled={isDispatchInProgress(order)}
                    onclick={(e) => { e.stopPropagation(); markOrderAsDelivered(order); }}
                    class="w-full py-2 px-3 rounded-lg text-xs font-bold bg-green-600 hover:bg-green-700 text-white shadow transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {#if isDispatchInProgress(order)}
                      <span class="inline-block animate-spin mr-1">⏳</span>
                    {/if}
                    {type === 'DELIVERY' ? (t.orders?.dispatch ?? 'Despachar') : (t.orders?.deliver ?? 'Entregar')}
                  </button>
                  <button
                    type="button"
                    onclick={(e) => { e.stopPropagation(); openCancelModal(order); }}
                    class="w-full py-1.5 px-3 rounded-lg text-xs font-semibold bg-red-600 text-white hover:bg-red-700 transition-colors"
                  >
                    ✕ {t.orders?.cancelOrder ?? 'Cancelar Pedido'}
                  </button>
                </div>
              {/if}
            {:else}
              {#if status === 'pending'}
                <button
                  type="button"
                  disabled={isActionInProgress(order, stationFilter, 'pending')}
                  onclick={(e) => { e.stopPropagation(); cycleOrderStatus(order, stationFilter); }}
                  class="w-full py-2 px-3 rounded-lg text-xs font-bold text-white shadow transition-colors disabled:opacity-50 disabled:cursor-not-allowed {useBarColor ? 'bg-blue-600 hover:bg-blue-700' : 'bg-orange-500 hover:bg-orange-600'}"
                >
                  {#if isActionInProgress(order, stationFilter, 'pending')}
                    <span class="inline-block animate-spin mr-1">⏳</span>
                  {/if}
                  <span aria-hidden="true">🔥</span> {t.orders?.startPreparing ?? 'Iniciar Preparación'}
                </button>
              {:else if status === 'preparing'}
                <button
                  type="button"
                  disabled={isActionInProgress(order, stationFilter, 'preparing')}
                  onclick={(e) => { e.stopPropagation(); cycleOrderStatus(order, stationFilter); }}
                  class="w-full py-2 px-3 rounded-lg text-xs font-bold text-white shadow transition-colors disabled:opacity-50 disabled:cursor-not-allowed {useBarColor ? 'bg-blue-500 hover:bg-blue-600' : 'bg-amber-500 hover:bg-amber-600'}"
                >
                  {#if isActionInProgress(order, stationFilter, 'preparing')}
                    <span class="inline-block animate-spin mr-1">⏳</span>
                  {/if}
                  ✓ {t.orders?.markAsReady ?? 'Listo'}
                </button>
              {:else}
                <div class="w-full py-2 px-3 rounded-lg text-xs font-bold bg-green-100 text-green-800 text-center">
                  ✓ {t.orders?.statusDone ?? 'Listo'}
                </div>
              {/if}
            {/if}
          </div>
        </li>
      {/snippet}
      <!-- Tabs en vista Entrega: Pendiente | Listo | Cancelado (solo en vista vertical) -->
      {#if !showQRView && stationFilter === 'ALL' && ordersViewMode === 'vertical'}
        <div class="flex items-stretch mb-3 w-full rounded-lg overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="tablist" aria-label="{t.orders?.statusPending ?? 'Pendiente'}, {t.orders?.statusDone ?? 'Listo'}, {t.orders?.cancelled ?? 'Cancelado'}">
          <button
            type="button"
            role="tab"
            aria-selected={deliveryTab === 'pending'}
            onclick={() => (deliveryTab = 'pending')}
            class="flex-1 min-w-0 px-2 py-2.5 text-xs font-semibold transition-colors border-r border-gray-200 {deliveryTab === 'pending' ? 'bg-gray-800 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            📦 {t.orders?.statusPending ?? 'Pendiente'} {ordersByDeliveryTab.pending.length > 0 ? `(${ordersByDeliveryTab.pending.length})` : ''}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={deliveryTab === 'delivered'}
            onclick={() => (deliveryTab = 'delivered')}
            class="flex-1 min-w-0 px-2 py-2.5 text-xs font-semibold transition-colors border-r border-gray-200 {deliveryTab === 'delivered' ? 'bg-green-100 text-green-800 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            ✓ {t.orders?.statusDone ?? 'Listo'} {ordersByDeliveryTab.delivered.length > 0 ? `(${ordersByDeliveryTab.delivered.length})` : ''}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={deliveryTab === 'cancelled'}
            onclick={() => (deliveryTab = 'cancelled')}
            class="flex-1 min-w-0 px-2 py-2.5 text-xs font-semibold transition-colors {deliveryTab === 'cancelled' ? 'bg-gray-300 text-gray-800 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            ✕ {t.orders?.cancelled ?? 'Cancelado'} {ordersByDeliveryTab.cancelled.length > 0 ? `(${ordersByDeliveryTab.cancelled.length})` : ''}
          </button>
        </div>
      {/if}
      <!-- Toggle vista: Vertical (tabs + lista) vs 3 columnas (Kanban). Cocina/Bar y Entrega. -->
      {#if !showQRView && (stationFilter === 'KITCHEN' || stationFilter === 'BAR' || stationFilter === 'ALL')}
        <div class="flex items-center gap-1.5 mb-2">
          <span class="text-xs font-medium text-gray-600">{t.orders?.viewVertical ?? 'Vista'}:</span>
          <div class="flex rounded overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="group" aria-label="{t.orders?.viewVertical ?? 'Vertical'} / {t.orders?.viewThreeColumns ?? '3 Columnas'}">
            <button
              type="button"
              onclick={() => ordersViewMode = 'vertical'}
              class="px-2 py-1.5 text-xs font-semibold transition-colors {ordersViewMode === 'vertical' ? 'bg-white text-gray-900 shadow-sm border border-gray-200' : 'text-gray-600 hover:bg-gray-200'}"
            >
              {t.orders?.viewVertical ?? 'Vertical'}
            </button>
            <button
              type="button"
              onclick={() => ordersViewMode = 'kanban'}
              class="px-2 py-1.5 text-xs font-semibold transition-colors border-l border-gray-200 {ordersViewMode === 'kanban' ? 'bg-white text-gray-900 shadow-sm border border-gray-200' : 'text-gray-600 hover:bg-gray-200'}"
            >
              {t.orders?.viewThreeColumns ?? '3 Columnas'}
            </button>
          </div>
        </div>
      {/if}
      <!-- Vista vertical: tabs + una sola lista -->
      {#if ordersViewMode === 'vertical' && !showQRView && (stationFilter === 'KITCHEN' || stationFilter === 'BAR')}
        <div class="flex items-stretch mb-3 w-full rounded-lg overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="tablist" aria-label="{t.orders?.tabPending ?? 'Pendientes'}, {t.orders?.tabPreparing ?? 'En Preparación'}, {t.orders?.tabDone ?? 'Listos'}">
          <button
            type="button"
            role="tab"
            aria-selected={operationalTab === 'pending'}
            onclick={() => operationalTab = 'pending'}
            class="flex-1 min-w-0 px-2 py-2.5 text-xs font-semibold transition-colors border-r border-gray-200 {operationalTab === 'pending' ? 'bg-amber-200 text-amber-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            🟠 {t.orders?.tabPending ?? 'Pendientes'} {ordersByTab.pending.length > 0 ? `(${ordersByTab.pending.length})` : ''}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={operationalTab === 'preparing'}
            onclick={() => operationalTab = 'preparing'}
            class="flex-1 min-w-0 px-2 py-2.5 text-xs font-semibold transition-colors border-r border-gray-200 {operationalTab === 'preparing' ? 'bg-blue-100 text-blue-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            🔵 {t.orders?.tabPreparing ?? 'En Preparación'} {ordersByTab.preparing.length > 0 ? `(${ordersByTab.preparing.length})` : ''}
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={operationalTab === 'done'}
            onclick={() => operationalTab = 'done'}
            class="flex-1 min-w-0 px-2 py-2.5 text-xs font-semibold transition-colors {operationalTab === 'done' ? 'bg-green-100 text-green-800 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
          >
            🟢 {t.orders?.tabDone ?? 'Listos'} {ordersByTab.done.length > 0 ? `(${ordersByTab.done.length})` : ''}
          </button>
        </div>
      {/if}
      {#if stationFilter === 'ALL' && ordersViewMode === 'kanban'}
        <!-- Vista Kanban Entrega: 3 columnas (Pendiente | Listo | Cancelado) -->
        {@const deliveryKanbanColumns = [{ key: 'pending', label: t.orders?.statusPending ?? 'Pendiente', orders: ordersByDeliveryTab.pending, icon: '📦', bg: 'bg-gray-100 border-gray-200', headerBg: 'bg-gray-800 text-white' }, { key: 'delivered', label: t.orders?.statusDone ?? 'Listo', orders: ordersByDeliveryTab.delivered, icon: '✓', bg: 'bg-green-50/80 border-green-200', headerBg: 'bg-green-100 text-green-800' }, { key: 'cancelled', label: t.orders?.cancelled ?? 'Cancelado', orders: ordersByDeliveryTab.cancelled, icon: '✕', bg: 'bg-gray-50 border-gray-200', headerBg: 'bg-gray-300 text-gray-800' }]}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-3">
          {#each deliveryKanbanColumns as col}
            <div class="flex flex-col rounded-lg border {col.bg} overflow-hidden min-h-[160px]">
              <button
                type="button"
                onclick={() => { deliveryTab = col.key as 'pending' | 'delivered' | 'cancelled'; ordersViewMode = 'vertical'; }}
                class="flex items-center justify-center gap-1.5 w-full px-3 py-2 text-xs font-bold text-left transition-colors hover:opacity-90 {col.headerBg}"
                title="{t.orders?.viewVertical ?? 'Ver'} {col.label}"
              >
                <span aria-hidden="true">{col.icon}</span>
                <span>{col.label}</span>
                <span class="tabular-nums">({col.orders.length})</span>
              </button>
              <ul class="flex-1 overflow-y-auto p-2 space-y-2">
                {#each col.orders as order, index (order.order_number)}
                  {@const isDel = col.key === 'delivered'}
                  {@const isCan = col.key === 'cancelled'}
                  {@render orderCard(order, isOrderReadyForDelivery(order) ? 'done' : getCajaOrderStatusLabel(order), index === 0, true, isDel, isCan)}
                {/each}
              </ul>
            </div>
          {/each}
        </div>
      {:else if stationFilter !== 'ALL' && ordersViewMode === 'kanban'}
        <!-- Vista Kanban Cocina/Bar: 3 columnas (Pendientes | En preparación | Listos) -->
        {@const kanbanColumns = [{ key: 'pending', label: t.orders?.tabPending ?? 'Pendientes', orders: ordersByTab.pending, icon: '🟠', bg: 'bg-amber-50 border-amber-200', headerBg: 'bg-amber-200 text-amber-900' }, { key: 'preparing', label: t.orders?.tabPreparing ?? 'En Preparación', orders: ordersByTab.preparing, icon: '🔵', bg: 'bg-blue-50/80 border-blue-200', headerBg: 'bg-blue-100 text-blue-900' }, { key: 'done', label: t.orders?.tabDone ?? 'Listos', orders: ordersByTab.done, icon: '🟢', bg: 'bg-green-50/80 border-green-200', headerBg: 'bg-green-100 text-green-800' }]}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-3">
          {#each kanbanColumns as col}
            <div class="flex flex-col rounded-lg border {col.bg} overflow-hidden min-h-[160px]">
              <button
                type="button"
                onclick={() => { operationalTab = col.key as 'pending' | 'preparing' | 'done'; ordersViewMode = 'vertical'; }}
                class="flex items-center justify-center gap-1.5 w-full px-3 py-2 text-xs font-bold text-left transition-colors hover:opacity-90 {col.headerBg}"
                title="{t.orders?.viewVertical ?? 'Ver'} {col.label}"
              >
                <span aria-hidden="true">{col.icon}</span>
                <span>{col.label}</span>
                <span class="tabular-nums">({col.orders.length})</span>
              </button>
              <ul class="flex-1 overflow-y-auto p-2 space-y-2">
                {#each col.orders as order, index (order.order_number)}
                  {@render orderCard(order, col.key as 'pending' | 'preparing' | 'done', index === 0, true, false, false)}
                {/each}
              </ul>
            </div>
          {/each}
        </div>
      {:else if ordersToShow.length === 0}
        <div class="rounded-lg bg-gray-100 border border-gray-200 p-6 text-center text-xs text-gray-600">
          {#if stationFilter === 'ALL' && deliveryTab === 'delivered'}
            No hay órdenes entregadas.
          {:else if stationFilter === 'ALL' && deliveryTab === 'cancelled'}
            {t.orders?.emptyCancelled ?? 'No hay órdenes canceladas.'}
          {:else if stationFilter === 'ALL'}
            {t.orders?.empty ?? 'No hay órdenes aún.'}
          {:else}
            {t.orders?.emptyForStation ?? 'No hay órdenes para esta estación.'}
          {/if}
        </div>
      {:else}
      <ul class="space-y-3 kitchen-orders-list" class:kitchen-mode-list={kitchenMode}>
        {#each ordersToShow as order, index (order.order_number)}
          {@const cardStation = stationFilter === 'ALL' ? null : stationFilter}
          {@const status = cardStation !== null ? getOrderStatusFromItems(order, cardStation) : (isOrderReadyForDelivery(order) ? 'done' : getCajaOrderStatusLabel(order))}
          {@const isDel = stationFilter === 'ALL' && isOrderFullyDelivered(order)}
          {@const isCan = stationFilter === 'ALL' && isOrderFullyCancelled(order)}
          {@render orderCard(order, status, index === 0, false, isDel, isCan)}
        {/each}
      </ul>
      {/if}
    {/if}
  </div>

  <!-- Modal cancelar pedido (solo vista Entrega) -->
  {#if orderToCancel}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      role="dialog"
      aria-modal="true"
      aria-labelledby="cancel-modal-title"
      onclick={() => closeCancelModal()}
      onkeydown={(e) => e.key === 'Escape' && closeCancelModal()}
    >
      <div
        class="bg-white rounded-xl shadow-xl max-w-sm w-full p-4"
        role="document"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
      >
        <h2 id="cancel-modal-title" class="text-base font-bold text-gray-900 flex items-center gap-2">
          <span aria-hidden="true">⚠️</span>
          {t.orders?.cancelModalTitle ?? 'Cancelar Pedido'} #{orderToCancel.order_number}
        </h2>
        <p class="mt-1 text-xs text-gray-600">
          {t.orders?.cancelModalSubtitle ?? 'Esta acción no se puede deshacer.'}
        </p>
        <p class="mt-3 text-xs font-semibold text-gray-700">
          {t.orders?.cancelModalReasonLabel ?? 'Motivo (elige uno):'}
        </p>
        <div class="mt-1.5 space-y-1.5" role="radiogroup" aria-label={t.orders?.cancelModalReasonLabel ?? 'Motivo'}>
          {#each CANCEL_REASON_KEYS as key}
            <label class="flex items-center gap-2 cursor-pointer text-sm text-gray-800">
              <input
                type="radio"
                name="cancelReason"
                value={key}
                bind:group={cancelReason}
                class="rounded-full border-gray-300 text-red-600 focus:ring-red-500"
              />
              <span>{t.orders?.cancelReasons?.[key] ?? key}</span>
            </label>
          {/each}
        </div>
        <p class="mt-3 text-xs font-semibold text-gray-700">
          {t.orders?.cancelModalCommentLabel ?? 'Comentario (opcional):'}
        </p>
        <textarea
          bind:value={cancelComment}
          placeholder={t.orders?.cancelModalCommentPlaceholder ?? 'Ej: cliente no contestó...'}
          rows="2"
          class="mt-1 w-full rounded-lg border border-gray-300 px-2.5 py-1.5 text-xs placeholder-gray-400 focus:border-red-500 focus:ring-1 focus:ring-red-500"
        ></textarea>
        <div class="mt-4 flex gap-2">
          <button
            type="button"
            disabled={cancelInProgress}
            onclick={closeCancelModal}
            class="flex-1 py-2 rounded-lg text-xs font-semibold border border-gray-300 text-gray-700 hover:bg-gray-50 disabled:opacity-50"
          >
            {t.orders?.cancelModalBack ?? 'Volver'}
          </button>
          <button
            type="button"
            disabled={cancelInProgress || !cancelReason.trim()}
            onclick={confirmCancelOrder}
            class="flex-1 py-2 rounded-lg text-xs font-semibold bg-red-600 text-white hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {#if cancelInProgress}
              <span class="inline-block animate-spin mr-1">⏳</span>
            {/if}
            {t.orders?.cancelModalConfirm ?? 'Confirmar Cancelación'}
          </button>
        </div>
      </div>
    </div>
  {/if}

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
        class="bg-white rounded-lg shadow-xl p-4 max-w-sm w-full"
        role="document"
        tabindex="0"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
      >
        <h2 id="qr-enlarged-title" class="text-base font-bold text-gray-800 mb-3">
          {qrEnlarged === 'ALL' ? (t.orders?.filterAll ?? 'Entrega') : qrEnlarged === 'KITCHEN' ? (t.orders?.filterKitchen ?? 'Cocina') : (t.orders?.filterBar ?? 'Barra')}
        </h2>
        <div class="flex justify-center mb-3 bg-white rounded border border-gray-200 p-2">
          <img
            src={`https://api.qrserver.com/v1/create-qr-code/?size=280x280&data=${encodeURIComponent(enlargedUrl)}`}
            alt="QR ampliado"
            class="w-56 h-56 sm:w-64 sm:h-64 object-contain"
          />
        </div>
        <p class="text-[11px] text-gray-500 mb-2">Para un QR nuevo, actualiza la página o vuelve a Órdenes.</p>
        <div class="flex gap-1.5">
          <button
            type="button"
            onclick={() => enlargedUrl && navigator.clipboard.writeText(enlargedUrl).then(() => alert('Enlace copiado'))}
            class="flex-1 py-1.5 rounded-lg bg-gray-800 text-white text-xs font-medium hover:bg-gray-900"
          >
            Copiar enlace
          </button>
          <button type="button" onclick={() => qrEnlarged = null} class="py-1.5 px-3 rounded-lg border border-gray-300 text-gray-700 text-xs font-medium hover:bg-gray-50">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Modo cocina: fullscreen, alto contraste, tipografía optimizada */
  :global(.kitchen-mode) {
    background: #f5f5f5;
  }
  :global(.kitchen-mode .kitchen-mode-header) {
    padding: 0.5rem 0.75rem;
    border-bottom-width: 2px;
  }
  :global(.kitchen-mode .kitchen-orders-list) {
    padding: 0.25rem;
  }
  :global(.kitchen-mode .kitchen-order-card) {
    font-size: 0.9375rem;
  }
  :global(.kitchen-mode .kitchen-order-first) {
    font-size: 1rem;
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
        <h2 id="paper-order-title" class="text-lg font-semibold text-gray-800">{t.orders?.viewAsPaper ?? 'Ver Como Hoja'}</h2>
        <div class="flex items-center gap-2">
          <button type="button" onclick={printThermal} class="px-3 py-1.5 rounded-lg bg-gray-800 text-white text-sm font-medium hover:bg-gray-900">
            {t.orders?.printThermal ?? 'Imprimir En Térmica'}
          </button>
          <button type="button" onclick={closePaperView} class="p-2 rounded-lg hover:bg-gray-100 text-gray-600" aria-label="Cerrar">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
      </div>
      <div class="overflow-y-auto p-4 print:p-0 {thermalPrintMode ? 'ticket-thermal-preview' : ''}" style="background: linear-gradient(to bottom, #fafafa 0%, #fff 2rem);">
        <div id="ticket-print" class="max-w-sm mx-auto text-gray-900 text-sm">
          <div class="border-b-2 border-dashed border-gray-300 pb-2 mb-2">
            <p class="text-center text-[11px] text-gray-500 uppercase tracking-widest">Orden</p>
            <p class="text-center text-2xl font-bold text-gray-900 mt-0.5">#{paperOrder.order_number}</p>
            <p class="text-center text-base font-semibold mt-1.5 text-gray-700">{formatDetailDate(paperOrder.requested_time)}</p>
            <p class="text-center text-xs text-gray-600 mt-0.5">{fulfillment.type === 'DELIVERY' ? (t.orders?.delivery ?? 'Envío') : (t.orders?.pickup ?? 'Retiro')}</p>
          </div>
          <div class="space-y-2 text-sm">
            <div>
              <p class="font-semibold text-gray-700 mb-0.5 text-xs">Contacto</p>
              <p class="font-medium text-sm">{(contact.fullName || '').trim() || '—'}</p>
              <p class="text-sm">{contact.phone?.trim() || '—'}</p>
              {#if contact.email?.trim()}<p class="text-gray-600 text-xs">{contact.email}</p>{/if}
            </div>
            {#if fulfillment.address?.rawAddress}
              <div>
                <p class="font-semibold text-gray-700 mb-0.5 text-xs">Dirección</p>
                <p class="text-gray-800 text-sm">{fulfillment.address.rawAddress}</p>
                {#if fulfillment.address.deliveryDetails?.unit || fulfillment.address.deliveryDetails?.notes}
                  <p class="text-gray-600 text-xs mt-0.5">Depto/Unidad: {fulfillment.address.deliveryDetails?.unit ?? '—'} · Notas: {fulfillment.address.deliveryDetails?.notes || '—'}</p>
                {/if}
              </div>
            {/if}
            <div class="border-t border-dashed border-gray-300 pt-2">
              <p class="font-semibold text-gray-700 mb-1 text-xs">Pedido</p>
              <ul class="space-y-1">
                {#each items as item}
                  <li class="flex justify-between gap-2 items-baseline">
                    <span class="text-sm"><span class="font-bold">{item.quantity ?? 0}</span> × {item.productName ?? '—'}</span>
                    <span class="shrink-0 text-sm font-semibold">{formatCurrency(item.totalPrice ?? 0, totals.currency ?? 'CLP')}</span>
                  </li>
                {/each}
              </ul>
            </div>
            <div class="border-t-2 border-dashed border-gray-300 pt-2 flex justify-between items-baseline text-base font-bold">
              <span>Total</span>
              <span class="text-lg">{formatCurrency(totals.total ?? 0, totals.currency ?? 'CLP')}</span>
            </div>
          </div>
          <div class="mt-4 pt-2 border-t border-dashed border-gray-300 text-center text-xs text-gray-400">
            {formatDetailDate(p?.createdAt as string ?? null)}
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
