<script lang="ts">
  import { onMount } from 'svelte'
  import { getKitchenOrdersFromProjection, subscribeMenuOrdersRealtime, refreshSupabaseToken, groupOrderItemsForDisplay, type KitchenOrder, type KitchenOrderItem, type StationFilter } from '../menuUtils'
  import { getActiveJourney } from '../journeyApi'
  import { dispatchOrder, cancelOrder, startPreparation, markReady } from '../orderApi'
  import { t as tStore } from '../useLanguage'
  import { playNewOrderSound, ensureAudioUnlocked } from '../utils/newOrderSound'

  interface Props {
    menuId: string
    station: 'KITCHEN' | 'BAR' | 'ALL'
  }
  const { menuId, station }: Props = $props()

  const STORAGE_TOKEN_KEY = $derived(`station_token_${menuId}_${station}`)
  const STORAGE_REFRESH_KEY = $derived(`station_refresh_${menuId}_${station}`)
  const STORAGE_OPERATOR_KEY = $derived(`station_operator_${menuId}_${station}`)

  /** Renovar token cada 50 min para que el cocinero no tenga que escanear de nuevo (~1 h de vida del access_token). */
  const REFRESH_INTERVAL_MS = 50 * 60 * 1000

  let orders = $state<KitchenOrder[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let token = $state<string | null>(null)
  let operatorName = $state('')
  let operatorSubmitted = $state(false)
  let realtimeUnsubscribe = $state<(() => void) | null>(null)
  /** Tab operativo en Cocina/Bar: Pendientes | En preparación | Listos (solo cuando station es KITCHEN o BAR). */
  let operationalTab = $state<'pending' | 'preparing' | 'done'>('pending')
  /** Tab en vista Entrega (station ALL): Pendiente | Listo | Cancelado (mismo patrón que consola de órdenes). */
  let deliveryTab = $state<'pending' | 'delivered' | 'cancelled'>('pending')
  /** Vista órdenes: vertical (tabs + lista) o kanban (3 columnas). */
  let ordersViewMode = $state<'vertical' | 'kanban'>('vertical')
  /** Modo full (pantalla completa) en vista pública. */
  let fullMode = $state(false)
  /** Escala de fuentes y tarjetas: pequeño / normal / grande (persistido en sessionStorage). */
  type FontScaleLevel = 'xxsmall' | 'xsmall' | 'small' | 'medium' | 'large'
  let fontScaleLevel = $state<FontScaleLevel>('medium')
  const STORAGE_SCALE_KEY = $derived(`station_scale_${menuId}_${station}`)
  const STORAGE_VIEW_MODE_KEY = $derived(`station_view_mode_${menuId}_${station}`)
  let dispatchInProgress = $state<Set<number>>(new Set())
  const cleanupRef = { intervalId: null as ReturnType<typeof setInterval> | null, refreshIntervalId: null as ReturnType<typeof setInterval> | null, unsub: null as (() => void) | null }
  /** Para detectar pedidos nuevos y reproducir sonido (mismo criterio que MenuOrders). */
  let previousOrderNumbers = new Set<number>()
  let initialLoadDone = false
  /** Modal cancelar pedido (solo vista Entrega). */
  let orderToCancel = $state<KitchenOrder | null>(null)
  let cancelReason = $state<string>('')
  let cancelComment = $state<string>('')
  let cancelInProgress = $state(false)
  const CANCEL_REASON_KEYS = ['outOfStock', 'orderError', 'customerLeft', 'paymentIssue', 'other'] as const

  const t = $derived($tStore)
  const stationFilter = $derived(station) as StationFilter
  const displayedOrders = $derived(orders)
  const stationLabel = $derived(station === 'KITCHEN' ? (t.orders?.filterKitchen ?? 'Cocina') : station === 'BAR' ? (t.orders?.filterBar ?? 'Barra') : (t.orders?.filterAll ?? 'Entrega'))

  /** Estado por estación derivado de los ítems de la orden (proyección). Igual criterio que consola MenuOrders. */
  function getOrderStatusFromItems(order: KitchenOrder, stationKey: string): 'pending' | 'preparing' | 'done' {
    const stationItems = order.items.filter((i) => (i.station ?? 'KITCHEN') === stationKey && i.status !== 'CANCELLED')
    if (stationItems.length === 0) return 'done'
    if (stationItems.some((i) => i.status === 'PENDING')) return 'pending'
    if (stationItems.some((i) => i.status === 'IN_PROGRESS')) return 'preparing'
    return 'done'
  }

  /** En Cocina/Bar: órdenes particionadas por estado para los tabs (desde proyección, no estado local). */
  const ordersByTab = $derived.by(() => {
    if (station !== 'KITCHEN' && station !== 'BAR') return { pending: [] as KitchenOrder[], preparing: [] as KitchenOrder[], done: [] as KitchenOrder[] }
    const stKey = station as 'KITCHEN' | 'BAR'
    const pending: KitchenOrder[] = []
    const preparing: KitchenOrder[] = []
    const done: KitchenOrder[] = []
    for (const o of displayedOrders) {
      const st = getOrderStatusFromItems(o, stKey)
      if (st === 'pending') pending.push(o)
      else if (st === 'preparing') preparing.push(o)
      else done.push(o)
    }
    return { pending, preparing, done }
  })

  /** En vista Entrega (ALL): órdenes por tab Pendiente | Listo | Cancelado (mismo patrón que MenuOrders). */
  const ordersByDeliveryTab = $derived.by(() => {
    if (station !== 'ALL') return { pending: [] as KitchenOrder[], delivered: [] as KitchenOrder[], cancelled: [] as KitchenOrder[] }
    const pending = displayedOrders.filter((o) => !isOrderFullyDelivered(o) && !isOrderFullyCancelled(o))
    const delivered = displayedOrders.filter((o) => isOrderFullyDelivered(o))
    const cancelled = displayedOrders.filter((o) => isOrderFullyCancelled(o))
    return { pending, delivered, cancelled }
  })

  /** Órdenes a mostrar: en Entrega (ALL) = las del tab activo (Pendiente/Listo/Cancelado); en Cocina/Bar = las del tab operativo. */
  const ordersToShow = $derived(station === 'ALL' ? ordersByDeliveryTab[deliveryTab] : ordersByTab[operationalTab])

  function getTokensFromHash(): { token: string | null; refresh_token: string | null } {
    if (typeof window === 'undefined') return { token: null, refresh_token: null }
    const hash = window.location.hash.slice(1)
    if (!hash) return { token: null, refresh_token: null }
    const params = new URLSearchParams(hash)
    return { token: params.get('token'), refresh_token: params.get('refresh_token') }
  }

  function clearHash() {
    if (typeof window === 'undefined') return
    const url = new URL(window.location.href)
    url.hash = ''
    window.history.replaceState(null, '', url.toString())
  }

  async function loadOrders(): Promise<void> {
    const accessToken = token || sessionStorage.getItem(STORAGE_TOKEN_KEY)
    if (!accessToken) {
      error = 'Enlace inválido o expirado. Escanee de nuevo el código.'
      loading = false
      return
    }
    try {
      loading = true
      error = null
      const journey = await getActiveJourney(menuId, accessToken)
      const newOrders = journey
        ? await getKitchenOrdersFromProjection(menuId, accessToken, stationFilter, journey.id)
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
        if (justArrived) {
          ensureAudioUnlocked()
          playNewOrderSound()
        }
      }
      newIds.forEach((id) => previousOrderNumbers.add(id))
      initialLoadDone = true
    } catch (err) {
      console.error('Error cargando órdenes:', err)
      error = err instanceof Error ? err.message : 'Error al cargar las órdenes'
      orders = []
    } finally {
      loading = false
    }
  }

  function formatRequestedTime(iso: string | null): string {
    if (!iso) return '—'
    return new Date(iso).toLocaleString('es-CL', { dateStyle: 'short', timeStyle: 'short' })
  }

  function getFulfillmentLabel(fulfillment: string): string {
    return fulfillment === 'DELIVERY' ? (t.orders?.delivery ?? 'Envío') : (t.orders?.pickup ?? 'Retiro')
  }

  function getRemainingMinutes(iso: string | null): number | null {
    if (!iso) return null
    return Math.round((new Date(iso).getTime() - Date.now()) / 60_000)
  }

  function getRemainingTimeLabel(minutes: number | null): string {
    if (minutes === null) return ''
    const abs = Math.abs(minutes)
    const isLate = minutes < 0
    const template = isLate ? (t.orders?.late ?? 'Atrasado {min} min') : (t.orders?.remainingIn ?? 'En {min} min')
    return template.replace('{min}', String(abs))
  }

  function getRemainingTimeColor(minutes: number | null): 'green' | 'yellow' | 'red' | null {
    if (minutes === null) return null
    if (minutes < 0) return 'red'
    if (minutes <= 5) return 'red'
    if (minutes <= 15) return 'yellow'
    return 'green'
  }

  function getItemCount(items: KitchenOrder['items']): number {
    return items.reduce((s, i) => s + i.quantity, 0)
  }

  /** Acción en curso en Cocina/Bar (start-preparation / mark-ready) para evitar doble clic. */
  let kitchenActionInProgress = $state<Set<string>>(new Set())

  function isKitchenActionInProgress(order: KitchenOrder, st: 'KITCHEN' | 'BAR'): boolean {
    const current = getOrderStatusFromItems(order, st)
    return kitchenActionInProgress.has(`${order.aggregate_id}-${current}-${st}`)
  }

  async function cycleOrderStatus(order: KitchenOrder, st: 'KITCHEN' | 'BAR') {
    const accessToken = token || (typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(STORAGE_TOKEN_KEY) : null)
    if (!menuId || !accessToken) return
    const current = getOrderStatusFromItems(order, st)
    const key = `${order.aggregate_id}-${current}-${st}`
    if (kitchenActionInProgress.has(key)) return
    kitchenActionInProgress = new Set(kitchenActionInProgress).add(key)
    try {
      if (current === 'pending') {
        await startPreparation(menuId, order.aggregate_id, accessToken, st, [])
      } else if (current === 'preparing') {
        const itemKeys = order.items
          .filter((i) => (i.station ?? 'KITCHEN') === st && i.status === 'IN_PROGRESS')
          .map((i) => i.item_key)
        if (itemKeys.length > 0) {
          await markReady(menuId, order.aggregate_id, accessToken, st, itemKeys)
        }
      }
      await loadOrders()
    } catch (err) {
      console.error('Error al actualizar estado en cocina:', err)
      error = err instanceof Error ? err.message : 'Error al actualizar estado'
    } finally {
      const next = new Set(kitchenActionInProgress)
      next.delete(key)
      kitchenActionInProgress = next
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

  /** Indica si Cocina y Barra marcaron la orden como lista (derivado de ítems de la proyección). */
  function isOrderReadyForDelivery(order: KitchenOrder): boolean {
    const stations = getStationsInOrder(order)
    return stations.every((st) => getOrderStatusFromItems(order, st) === 'done')
  }

  /** Orden terminada: todos los ítems activos en DELIVERED (PICKUP) o DISPATCHED (DELIVERY). */
  function isOrderFullyDelivered(order: KitchenOrder): boolean {
    const active = order.items.filter((i) => i.status !== 'CANCELLED')
    return active.length > 0 && active.every((i) => i.status === 'DISPATCHED' || i.status === 'DELIVERED')
  }

  /** Orden cancelada: todos los ítems en CANCELLED. */
  function isOrderFullyCancelled(order: KitchenOrder): boolean {
    return order.items.length > 0 && order.items.every((i) => i.status === 'CANCELLED')
  }

  async function markOrderAsDelivered(order: KitchenOrder) {
    const accessToken = token || (typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(STORAGE_TOKEN_KEY) : null)
    if (!menuId || !accessToken) return
    if (dispatchInProgress.has(order.aggregate_id)) return
    dispatchInProgress = new Set(dispatchInProgress).add(order.aggregate_id)
    try {
      await dispatchOrder(menuId, order.aggregate_id, accessToken)
      await loadOrders()
    } catch (err) {
      console.error('Error al marcar como entregado:', err)
    } finally {
      const next = new Set(dispatchInProgress)
      next.delete(order.aggregate_id)
      dispatchInProgress = next
    }
  }

  /** Estado resumido para mostrar en vista Entrega (Pendiente / En preparación). */
  function getCajaOrderStatusLabel(order: KitchenOrder): 'pending' | 'preparing' {
    const stations = getStationsInOrder(order)
    const statuses = stations.map((st) => getOrderStatusFromItems(order, st))
    if (statuses.every((s) => s === 'pending')) return 'pending'
    return 'preparing'
  }

  async function toggleFullMode() {
    try {
      if (fullMode) {
        await document.exitFullscreen?.()
        fullMode = false
      } else {
        await document.documentElement.requestFullscreen?.()
        fullMode = true
      }
    } catch (_) {
      fullMode = !!document.fullscreenElement
    }
  }

  async function startOrdersAndRealtime(cleanupRef: { intervalId: ReturnType<typeof setInterval> | null; refreshIntervalId: ReturnType<typeof setInterval> | null; unsub: (() => void) | null }): Promise<void> {
    const accessToken = token || (typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(STORAGE_TOKEN_KEY) : null)
    if (!accessToken) return
    await loadOrders()
    const unsub = await subscribeMenuOrdersRealtime(menuId, accessToken, () => loadOrders())
    realtimeUnsubscribe = unsub
    cleanupRef.unsub = unsub
    cleanupRef.intervalId = setInterval(() => loadOrders(), 20_000)
  }

  /** Renueva el access_token con el refresh_token y vuelve a suscribir realtime/polling. */
  async function doRefreshToken(): Promise<void> {
    const refreshTokenValue = typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(STORAGE_REFRESH_KEY) : null
    if (!refreshTokenValue) return
    const data = await refreshSupabaseToken(refreshTokenValue)
    if (!data) {
      error = 'Sesión expirada. Escanee de nuevo el código desde la consola.'
      if (typeof sessionStorage !== 'undefined') {
        sessionStorage.removeItem(STORAGE_TOKEN_KEY)
        sessionStorage.removeItem(STORAGE_REFRESH_KEY)
      }
      token = null
      return
    }
    token = data.access_token
    if (typeof sessionStorage !== 'undefined') {
      sessionStorage.setItem(STORAGE_TOKEN_KEY, data.access_token)
      if (data.refresh_token) sessionStorage.setItem(STORAGE_REFRESH_KEY, data.refresh_token)
    }
    error = null
    // Re-suscribir realtime y polling con el nuevo token
    cleanupRef.intervalId && clearInterval(cleanupRef.intervalId)
    cleanupRef.unsub?.()
    cleanupRef.intervalId = null
    cleanupRef.unsub = null
    await startOrdersAndRealtime(cleanupRef)
  }

  function submitOperatorName() {
    const name = operatorName.trim()
    if (!name) return
    ensureAudioUnlocked()
    if (typeof localStorage !== 'undefined') localStorage.setItem(STORAGE_OPERATOR_KEY, name)
    operatorSubmitted = true
    startOrdersAndRealtime(cleanupRef)
  }

  function setFontScaleLevel(level: FontScaleLevel) {
    fontScaleLevel = level
    if (typeof localStorage !== 'undefined') localStorage.setItem(STORAGE_SCALE_KEY, level)
  }

  function decreaseFontScale() {
    if (fontScaleLevel === 'large') setFontScaleLevel('medium')
    else if (fontScaleLevel === 'medium') setFontScaleLevel('small')
    else if (fontScaleLevel === 'small') setFontScaleLevel('xsmall')
    else if (fontScaleLevel === 'xsmall') setFontScaleLevel('xxsmall')
  }

  function increaseFontScale() {
    if (fontScaleLevel === 'xxsmall') setFontScaleLevel('xsmall')
    else if (fontScaleLevel === 'xsmall') setFontScaleLevel('small')
    else if (fontScaleLevel === 'small') setFontScaleLevel('medium')
    else if (fontScaleLevel === 'medium') setFontScaleLevel('large')
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
    const accessToken = token || (typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(STORAGE_TOKEN_KEY) : null)
    if (!orderToCancel || !menuId || !accessToken || !cancelReason.trim()) return
    cancelInProgress = true
    try {
      const reasonLabel = getCancelReasonLabel(cancelReason)
      const reasonText = cancelComment.trim()
        ? `${reasonLabel}: ${cancelComment.trim()}`
        : reasonLabel
      await cancelOrder(menuId, orderToCancel.aggregate_id, accessToken, reasonText)
      orderToCancel = null
      cancelReason = ''
      cancelComment = ''
      await loadOrders()
    } catch (err) {
      console.error('Error al cancelar la orden:', err)
      error = err instanceof Error ? err.message : 'Error al cancelar la orden'
    } finally {
      cancelInProgress = false
    }
  }

  onMount(() => {
    if (typeof window === 'undefined') return
    const { token: fromHash, refresh_token: refreshFromHash } = getTokensFromHash()
    if (fromHash) {
      sessionStorage.setItem(STORAGE_TOKEN_KEY, fromHash)
      token = fromHash
      if (refreshFromHash) {
        sessionStorage.setItem(STORAGE_REFRESH_KEY, refreshFromHash)
      }
      clearHash()
    } else {
      token = sessionStorage.getItem(STORAGE_TOKEN_KEY)
    }

    const savedOperator = localStorage.getItem(STORAGE_OPERATOR_KEY) ?? sessionStorage.getItem(STORAGE_OPERATOR_KEY)
    if (savedOperator) {
      operatorName = savedOperator
      operatorSubmitted = true
    }

    if (operatorSubmitted && token) {
      startOrdersAndRealtime(cleanupRef)
      // Renovar token cada ~50 min para que dure todo el turno sin volver a escanear
      if (sessionStorage.getItem(STORAGE_REFRESH_KEY)) {
        cleanupRef.refreshIntervalId = setInterval(() => doRefreshToken(), REFRESH_INTERVAL_MS)
      }
    }

    const savedScale = localStorage.getItem(STORAGE_SCALE_KEY)
    if (savedScale === 'xxsmall' || savedScale === 'xsmall' || savedScale === 'small' || savedScale === 'medium' || savedScale === 'large') {
      fontScaleLevel = savedScale
    }
    const savedViewMode = localStorage.getItem(STORAGE_VIEW_MODE_KEY)
    if (savedViewMode === 'vertical' || savedViewMode === 'kanban') {
      ordersViewMode = savedViewMode
    }

    const onFullscreenChange = () => {
      fullMode = !!document.fullscreenElement
    }
    document.addEventListener('fullscreenchange', onFullscreenChange)

    return () => {
      document.removeEventListener('fullscreenchange', onFullscreenChange)
      if (cleanupRef.intervalId) clearInterval(cleanupRef.intervalId)
      if (cleanupRef.refreshIntervalId) clearInterval(cleanupRef.refreshIntervalId)
      cleanupRef.unsub?.()
      realtimeUnsubscribe = null
    }
  })
</script>

<div class="min-h-screen flex flex-col bg-gray-50" class:kitchen-mode={fullMode}>
  <!-- Sin token -->
  {#if typeof window !== 'undefined' && !token && !sessionStorage.getItem(STORAGE_TOKEN_KEY)}
    <div class="flex-1 flex items-center justify-center p-6">
      <div class="text-center max-w-md">
        <p class="text-red-600 font-medium">Enlace inválido o expirado.</p>
        <p class="text-gray-600 mt-2 text-sm">Escanee de nuevo el código de {stationLabel} desde la consola del dueño del menú.</p>
      </div>
    </div>
  {:else if !operatorSubmitted}
    <!-- Pedir nombre del operador -->
    <div class="flex-1 flex items-center justify-center p-6">
      <div class="w-full max-w-sm bg-white rounded-xl shadow-lg border border-gray-200 p-6">
        <h1 class="text-xl font-bold text-gray-800 mb-2">Vista de {stationLabel}</h1>
        <p class="text-gray-600 text-sm mb-4">Ingrese su nombre para gestionar los pedidos entrantes.</p>
        <input
          type="text"
          bind:value={operatorName}
          placeholder="Ej. Juan, María..."
          class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
          onkeydown={(e) => e.key === 'Enter' && submitOperatorName()}
        />
        <button
          type="button"
          onclick={submitOperatorName}
          disabled={!operatorName.trim()}
          class="mt-4 w-full py-3 rounded-lg font-semibold text-white bg-amber-500 hover:bg-amber-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Entrar
        </button>
      </div>
    </div>
  {:else}
    <!-- Header -->
    <div class="flex-shrink-0 px-4 py-4 border-b border-gray-200 bg-white flex flex-wrap items-center justify-between gap-2">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-gray-800">
          {ordersViewMode === 'kanban' ? `${t.sidebar?.orders ?? 'Kanban'} - ${stationLabel}` : stationLabel} — {operatorName || 'Operador'}
        </h1>
        {#if !fullMode}
          <p class="text-sm text-gray-500 mt-1">Pedidos en tiempo real. Sin login.</p>
        {/if}
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <!-- Tamaño fuentes y tarjetas: A− / A+ -->
        <div class="flex rounded overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="group" aria-label="Tamaño de texto y tarjetas">
          <button
            type="button"
            onclick={decreaseFontScale}
            disabled={fontScaleLevel === 'xxsmall'}
            class="px-3 py-2 text-sm font-bold text-gray-700 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            title="Reducir tamaño"
          >
            A−
          </button>
          <button
            type="button"
            onclick={increaseFontScale}
            disabled={fontScaleLevel === 'large'}
            class="px-3 py-2 text-sm font-bold text-gray-700 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors border-l border-gray-200"
            title="Aumentar tamaño"
          >
            A+
          </button>
        </div>
        <!-- Toggle vista: Vertical / Kanban (3 columnas) -->
        <div class="flex rounded overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="group" aria-label="{t.orders?.viewVertical ?? 'Vertical'} / {t.orders?.viewThreeColumns ?? '3 Columnas'}">
          <button
            type="button"
            onclick={() => { ensureAudioUnlocked(); ordersViewMode = 'vertical'; if (typeof localStorage !== 'undefined') localStorage.setItem(STORAGE_VIEW_MODE_KEY, 'vertical'); }}
            class="px-3 py-2 text-sm font-semibold transition-colors {ordersViewMode === 'vertical' ? 'bg-white text-gray-900 shadow-sm border border-gray-200' : 'text-gray-600 hover:bg-gray-200'}"
          >
            {t.orders?.viewVertical ?? 'Vertical'}
          </button>
          <button
            type="button"
            onclick={() => { ensureAudioUnlocked(); ordersViewMode = 'kanban'; if (typeof localStorage !== 'undefined') localStorage.setItem(STORAGE_VIEW_MODE_KEY, 'kanban'); }}
            class="px-3 py-2 text-sm font-semibold transition-colors border-l border-gray-200 {ordersViewMode === 'kanban' ? 'bg-white text-gray-900 shadow-sm border border-gray-200' : 'text-gray-600 hover:bg-gray-200'}"
          >
            {t.orders?.viewThreeColumns ?? '3 Columnas'}
          </button>
        </div>
        <button
          type="button"
          onclick={() => { ensureAudioUnlocked(); toggleFullMode(); }}
          class="rounded-lg px-4 py-2 text-sm font-semibold shrink-0 {fullMode ? 'bg-amber-500 text-white hover:bg-amber-600' : 'bg-gray-200 text-gray-800 hover:bg-gray-300'}"
        >
          {fullMode ? (t.orders?.exitKitchenMode ?? 'Salir modo full') : (t.orders?.kitchenMode ?? 'Modo full')}
        </button>
      </div>
    </div>

    <!-- Lista de órdenes (escala aplicada con transform para que cards y fuentes cambien) -->
    <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
      <div class="station-scale-wrapper" data-scale={fontScaleLevel}>
      {#if loading && orders.length === 0}
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-10 w-10 border-2 border-amber-500 border-t-transparent"></div>
        </div>
      {:else if error}
        <div class="rounded-lg bg-red-50 border border-red-200 p-4 text-red-800">{error}</div>
      {:else}
        {#snippet orderCard(order: KitchenOrder, isFirst: boolean, columnKey: 'pending' | 'preparing' | 'done' | 'delivered' | 'cancelled')}
          {@const type = order.fulfillment}
          {@const itemCount = getItemCount(order.items)}
          {@const remainingMin = getRemainingMinutes(order.requested_time)}
          {@const timeColor = getRemainingTimeColor(remainingMin)}
          {@const useBarColor = station === 'BAR' || (station === 'ALL' && order.items.some((i: KitchenOrderItem) => i.station === 'BAR'))}
          {@const isDoneTab = (station === 'ALL' && (columnKey === 'delivered' || columnKey === 'cancelled')) || ((station === 'KITCHEN' || station === 'BAR') && columnKey === 'done')}
          {@const cardStatus = station === 'ALL' ? null : (columnKey as 'pending' | 'preparing' | 'done')}
          {@const kitchenSt = station === 'ALL' ? getOrderStatusFromItems(order, 'KITCHEN') : null}
          {@const barSt = station === 'ALL' ? getOrderStatusFromItems(order, 'BAR') : null}
          {@const orderHasBar = station === 'ALL' ? order.items.some((i) => i.station === 'BAR') : false}
          {@const readyForDelivery = station === 'ALL' ? isOrderReadyForDelivery(order) : false}
          <li class="bg-white rounded-xl border-2 overflow-hidden station-order-card {isDoneTab ? 'order-card-done opacity-80 border-gray-300' : ''} {isFirst && !isDoneTab ? 'border-amber-400 shadow-lg' : 'border-gray-200'}">
            <div class="w-full px-4 py-3 sm:px-5 flex flex-wrap items-center gap-4 border-b border-gray-100 {isFirst && !isDoneTab ? 'sm:py-6' : 'sm:py-4'}">
              <span class="font-bold text-gray-900 tabular-nums {isFirst && !isDoneTab ? 'text-4xl sm:text-5xl md:text-6xl' : 'text-3xl sm:text-4xl'}">#{order.order_number}</span>
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
                {getFulfillmentLabel(type)}
              </span>
              {#if station === 'ALL'}
                <div class="flex flex-wrap items-center gap-3 text-sm font-semibold">
                  <span class="inline-flex items-center gap-1">{t.orders?.filterKitchen ?? 'Cocina'}: {kitchenSt === 'done' ? '✔️' : '⏳'}</span>
                  <span class="inline-flex items-center gap-1">{t.orders?.filterBar ?? 'Barra'}: {orderHasBar ? (barSt === 'done' ? '✔️' : '⏳') : '—'}</span>
                  <span class="inline-flex items-center gap-1 rounded-full border px-2 py-1 {readyForDelivery ? 'bg-green-50 text-green-800 border-green-200' : 'bg-amber-50 text-amber-900 border-amber-200'}">
                    {t.orders?.statusGeneralLabel ?? 'Estado General'}: {readyForDelivery ? (t.orders?.readyToDeliver ?? 'Listo Para Entregar') : (t.orders?.statusPreparing ?? 'En Preparación')}
                  </span>
                </div>
              {:else}
                <span class="inline-flex items-center gap-1 rounded-full font-bold border {isFirst ? 'px-4 py-2 text-base sm:text-lg' : 'px-3 py-1 text-sm'}
                  {cardStatus === 'pending' ? 'bg-gray-100 text-gray-700 border-gray-200' : ''}
                  {cardStatus === 'preparing' ? 'bg-amber-50 text-amber-900 border-amber-200' : ''}
                  {cardStatus === 'done' ? 'bg-green-50 text-green-800 border-green-200' : ''}">
                  {#if cardStatus === 'preparing'}<span aria-hidden="true">⏳</span>{/if}
                  {cardStatus === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : cardStatus === 'preparing' ? (t.orders?.statusPreparing ?? 'En Preparación') : (t.orders?.statusDone ?? 'Listo')}
                </span>
              {/if}
            </div>
            <div class="px-4 py-3 sm:px-5 bg-amber-50/50 border-b border-amber-100 {isFirst ? 'py-4 sm:py-5' : ''}">
              <div class="flex items-center justify-between gap-2 mb-2">
                <p class="font-semibold text-amber-800 uppercase tracking-wide {isFirst ? 'text-sm' : 'text-xs'}">{t.orders?.itemsToPrepare ?? 'Qué Preparar'}</p>
                <span class="text-sm font-bold text-amber-800 tabular-nums">{(t.orders?.itemsCount ?? '{count} Ítems').replace('{count}', String(itemCount))}</span>
              </div>
              <ul class="space-y-1 text-gray-900 {isFirst ? 'text-xl sm:text-2xl md:text-3xl' : 'text-lg sm:text-xl'}">
                {#each groupOrderItemsForDisplay(order.items) as item}
                  <li class="tabular-nums">
                    <span class="font-bold text-amber-800">{item.quantity}×</span> <span class="font-normal">{item.item_name}</span>
                  </li>
                {/each}
              </ul>
            </div>
            <div class="px-4 py-3 sm:px-5 border-t border-gray-100">
              {#if station === 'ALL'}
                {#if isOrderFullyCancelled(order)}
                  <div class="w-full py-3 px-4 rounded-xl text-base font-bold bg-gray-200 text-gray-700 text-center">
                    ✕ {t.orders?.cancelled ?? 'Cancelado'}
                  </div>
                {:else if isOrderFullyDelivered(order)}
                  <div class="w-full py-3 px-4 rounded-xl text-base font-bold bg-green-100 text-green-800 text-center">
                    ✓ {type === 'PICKUP' ? (t.orders?.delivered ?? 'Entregado') : (t.orders?.dispatched ?? 'Despachado')}
                  </div>
                {:else}
                  <div class="flex flex-col gap-2">
                    <button
                      type="button"
                      disabled={dispatchInProgress.has(order.aggregate_id)}
                      onclick={(e) => { e.stopPropagation(); markOrderAsDelivered(order); }}
                      class="w-full py-3 px-4 rounded-xl text-base font-bold bg-green-600 hover:bg-green-700 text-white shadow-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {#if dispatchInProgress.has(order.aggregate_id)}
                        <span class="inline-block animate-spin mr-1">⏳</span>
                      {/if}
                      {type === 'PICKUP' ? (t.orders?.deliver ?? 'Entregar') : (t.orders?.dispatch ?? 'Despachar')}
                    </button>
                    <button
                      type="button"
                      onclick={(e) => { e.stopPropagation(); openCancelModal(order); }}
                      class="w-full py-2.5 px-4 rounded-xl text-base font-semibold bg-red-600 hover:bg-red-700 text-white shadow-md transition-colors"
                    >
                      ✕ {t.orders?.cancelOrder ?? 'Cancelar Pedido'}
                    </button>
                  </div>
                {/if}
              {:else}
                {#if cardStatus === 'pending'}
                  <button
                    type="button"
                    disabled={isKitchenActionInProgress(order, station)}
                    onclick={() => cycleOrderStatus(order, station)}
                    class="w-full py-3 px-4 rounded-xl text-base font-bold text-white shadow-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed {useBarColor ? 'bg-blue-600 hover:bg-blue-700' : 'bg-orange-500 hover:bg-orange-600'}"
                  >
                    {#if isKitchenActionInProgress(order, station)}
                      <span class="inline-block animate-spin mr-1">⏳</span>
                    {:else}
                      <span aria-hidden="true">🔥</span>
                    {/if}
                    {t.orders?.startPreparing ?? 'Iniciar Preparación'}
                  </button>
                {:else if cardStatus === 'preparing'}
                  <button
                    type="button"
                    disabled={isKitchenActionInProgress(order, station)}
                    onclick={() => cycleOrderStatus(order, station)}
                    class="w-full py-3 px-4 rounded-xl text-base font-bold text-white shadow-md disabled:opacity-50 disabled:cursor-not-allowed {useBarColor ? 'bg-blue-500 hover:bg-blue-600' : 'bg-amber-500 hover:bg-amber-600'}"
                  >
                    {#if isKitchenActionInProgress(order, station)}
                      <span class="inline-block animate-spin mr-1">⏳</span>
                    {/if}
                    ✓ {t.orders?.markAsReady ?? 'LISTO'}
                  </button>
                {:else}
                  <div class="w-full py-3 px-4 rounded-xl text-base font-bold bg-green-100 text-green-800 text-center">
                    ✓ {t.orders?.statusDone ?? 'Listo'}
                  </div>
                {/if}
              {/if}
            </div>
          </li>
        {/snippet}
        <!-- Tabs solo en vista vertical -->
        <!-- Tabs vista Entrega (ALL): Pendiente | Listo | Cancelado -->
        {#if station === 'ALL' && ordersViewMode === 'vertical'}
          <div class="flex items-stretch mb-4 w-full rounded-xl overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="tablist" aria-label="{t.orders?.statusPending ?? 'Pendiente'}, {t.orders?.statusDone ?? 'Listo'}, {t.orders?.cancelled ?? 'Cancelado'}">
            <button
              type="button"
              role="tab"
              aria-selected={deliveryTab === 'pending'}
              onclick={() => (deliveryTab = 'pending')}
              class="flex-1 min-w-0 px-3 py-3.5 text-base font-semibold transition-colors border-r border-gray-200 {deliveryTab === 'pending' ? 'bg-gray-800 text-white shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
            >
              📦 {t.orders?.statusPending ?? 'Pendiente'} {ordersByDeliveryTab.pending.length > 0 ? `(${ordersByDeliveryTab.pending.length})` : ''}
            </button>
            <button
              type="button"
              role="tab"
              aria-selected={deliveryTab === 'delivered'}
              onclick={() => (deliveryTab = 'delivered')}
              class="flex-1 min-w-0 px-3 py-3.5 text-base font-semibold transition-colors border-r border-gray-200 {deliveryTab === 'delivered' ? 'bg-green-100 text-green-800 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
            >
              ✓ {t.orders?.statusDone ?? 'Listo'} {ordersByDeliveryTab.delivered.length > 0 ? `(${ordersByDeliveryTab.delivered.length})` : ''}
            </button>
            <button
              type="button"
              role="tab"
              aria-selected={deliveryTab === 'cancelled'}
              onclick={() => (deliveryTab = 'cancelled')}
              class="flex-1 min-w-0 px-3 py-3.5 text-base font-semibold transition-colors {deliveryTab === 'cancelled' ? 'bg-gray-300 text-gray-800 shadow-sm' : 'text-gray-600 hover:bg-gray-200'}"
            >
              ✕ {t.orders?.cancelled ?? 'Cancelado'} {ordersByDeliveryTab.cancelled.length > 0 ? `(${ordersByDeliveryTab.cancelled.length})` : ''}
            </button>
          </div>
        {/if}
        <!-- Segmented control: Cocina y Bar, solo en vista vertical -->
        {#if (station === 'KITCHEN' || station === 'BAR') && ordersViewMode === 'vertical'}
          <div class="flex items-stretch mb-4 w-full rounded-xl overflow-hidden border border-gray-200 bg-gray-100 shadow-inner" role="tablist" aria-label="{t.orders?.tabPending ?? 'Pendientes'}, {t.orders?.tabPreparing ?? 'En Preparación'}, {t.orders?.tabDone ?? 'Listos'}">
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
              🔵 {t.orders?.tabPreparing ?? 'En Preparación'} {ordersByTab.preparing.length > 0 ? `(${ordersByTab.preparing.length})` : ''}
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
        {#if station === 'ALL' && ordersViewMode === 'kanban'}
          <!-- Kanban Entrega: Pendiente | Listo | Cancelado -->
          {@const deliveryCols = [{ key: 'pending' as const, label: t.orders?.statusPending ?? 'Pendiente', orders: ordersByDeliveryTab.pending, icon: '📦', bg: 'bg-gray-100 border-gray-200', headerBg: 'bg-gray-800 text-white' }, { key: 'delivered' as const, label: t.orders?.statusDone ?? 'Listo', orders: ordersByDeliveryTab.delivered, icon: '✓', bg: 'bg-green-50/80 border-green-200', headerBg: 'bg-green-100 text-green-800' }, { key: 'cancelled' as const, label: t.orders?.cancelled ?? 'Cancelado', orders: ordersByDeliveryTab.cancelled, icon: '✕', bg: 'bg-gray-50 border-gray-200', headerBg: 'bg-gray-300 text-gray-800' }]}
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            {#each deliveryCols as col}
              <div class="flex flex-col rounded-xl border {col.bg} overflow-hidden min-h-[200px]">
                <div class="flex items-center justify-center gap-2 w-full px-4 py-3 text-sm font-bold {col.headerBg}">
                  <span aria-hidden="true">{col.icon}</span>
                  <span>{col.label}</span>
                  <span class="tabular-nums">({col.orders.length})</span>
                </div>
                <ul class="flex-1 overflow-y-auto p-3 space-y-4">
                  {#each col.orders as order, index (order.order_number)}
                    {@render orderCard(order, index === 0, col.key)}
                  {/each}
                </ul>
              </div>
            {/each}
          </div>
        {:else if (station === 'KITCHEN' || station === 'BAR') && ordersViewMode === 'kanban'}
          <!-- Kanban Cocina/Bar: Pendientes | En preparación | Listos -->
          {@const opCols = [{ key: 'pending' as const, label: t.orders?.tabPending ?? 'Pendientes', orders: ordersByTab.pending, icon: '🟠', bg: 'bg-amber-50 border-amber-200', headerBg: 'bg-amber-200 text-amber-900' }, { key: 'preparing' as const, label: t.orders?.tabPreparing ?? 'En Preparación', orders: ordersByTab.preparing, icon: '🔵', bg: 'bg-blue-50/80 border-blue-200', headerBg: 'bg-blue-100 text-blue-900' }, { key: 'done' as const, label: t.orders?.tabDone ?? 'Listos', orders: ordersByTab.done, icon: '🟢', bg: 'bg-green-50/80 border-green-200', headerBg: 'bg-green-100 text-green-800' }]}
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            {#each opCols as col}
              <div class="flex flex-col rounded-xl border {col.bg} overflow-hidden min-h-[200px]">
                <div class="flex items-center justify-center gap-2 w-full px-4 py-3 text-sm font-bold {col.headerBg}">
                  <span aria-hidden="true">{col.icon}</span>
                  <span>{col.label}</span>
                  <span class="tabular-nums">({col.orders.length})</span>
                </div>
                <ul class="flex-1 overflow-y-auto p-3 space-y-4">
                  {#each col.orders as order, index (order.order_number)}
                    {@render orderCard(order, index === 0, col.key)}
                  {/each}
                </ul>
              </div>
            {/each}
          </div>
        {:else if ordersToShow.length === 0}
          <div class="rounded-lg bg-gray-100 border border-gray-200 p-8 text-center text-gray-600">
            {#if station === 'ALL' && deliveryTab === 'delivered'}
              {t.orders?.empty ?? 'No Hay Órdenes Aún.'}
            {:else if station === 'ALL' && deliveryTab === 'cancelled'}
              {t.orders?.emptyCancelled ?? 'No Hay Órdenes Canceladas.'}
            {:else if station === 'ALL'}
              {t.orders?.empty ?? 'No Hay Órdenes Aún.'}
            {:else}
              {t.orders?.emptyForStation ?? 'No Hay Órdenes Para Esta Estación.'}
            {/if}
          </div>
        {:else}
          <ul class="space-y-5">
            {#each ordersToShow as order, index (order.order_number)}
              {@render orderCard(order, index === 0, station === 'ALL' ? deliveryTab : operationalTab)}
            {/each}
          </ul>
        {/if}
      {/if}
      </div>
    </div>
  {/if}

  <!-- Modal cancelar pedido (solo vista Entrega) -->
  {#if orderToCancel}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      role="dialog"
      aria-modal="true"
      aria-labelledby="cancel-modal-title"
      tabindex="-1"
      onclick={(e) => e.target === e.currentTarget && closeCancelModal()}
      onkeydown={(e) => e.key === 'Escape' && closeCancelModal()}
    >
      <div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-4">
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
</div>

<style>
  /* Modo full: mismo aspecto que en admin (cocina) */
  :global(.kitchen-mode) {
    background: #f5f5f5;
  }

  /* Escala de cards y fuentes con transform (Tailwind usa rem respecto a :root; scale() sí escala todo) */
  .station-scale-wrapper[data-scale='xxsmall'] {
    transform: scale(0.65);
    transform-origin: top left;
    width: 153.846%; /* 1/0.65 */
  }
  .station-scale-wrapper[data-scale='xsmall'] {
    transform: scale(0.75);
    transform-origin: top left;
    width: 133.333%; /* 1/0.75 */
  }
  .station-scale-wrapper[data-scale='small'] {
    transform: scale(0.875);
    transform-origin: top left;
    width: 114.286%; /* 1/0.875 → tras escalar el ancho visual es 100% */
  }
  .station-scale-wrapper[data-scale='medium'] {
    transform: scale(1);
    transform-origin: top left;
    width: 100%;
  }
  .station-scale-wrapper[data-scale='large'] {
    transform: scale(1.15);
    transform-origin: top left;
    width: 86.956%; /* 1/1.15 */
  }
</style>
