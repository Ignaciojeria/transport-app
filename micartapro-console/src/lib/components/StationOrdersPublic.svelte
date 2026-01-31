<script lang="ts">
  import { onMount } from 'svelte'
  import { getKitchenOrdersFromProjection, subscribeMenuOrdersRealtime, type KitchenOrder, type StationFilter } from '../menuUtils'
  import { t as tStore } from '../useLanguage'

  interface Props {
    menuId: string
    station: 'KITCHEN' | 'BAR'
  }
  const { menuId, station }: Props = $props()

  const STORAGE_TOKEN_KEY = $derived(`station_token_${menuId}_${station}`)
  const STORAGE_OPERATOR_KEY = $derived(`station_operator_${menuId}_${station}`)

  let orders = $state<KitchenOrder[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let token = $state<string | null>(null)
  let operatorName = $state('')
  let operatorSubmitted = $state(false)
  let realtimeUnsubscribe = $state<(() => void) | null>(null)
  let orderStatus = $state<Record<string, 'pending' | 'preparing' | 'done'>>({})
  const cleanupRef = { intervalId: null as ReturnType<typeof setInterval> | null, unsub: null as (() => void) | null }

  const t = $derived($tStore)
  const stationFilter = $derived(station) as StationFilter
  const displayedOrders = $derived(orders)
  const stationLabel = $derived(station === 'KITCHEN' ? (t.orders?.filterKitchen ?? 'Cocina') : (t.orders?.filterBar ?? 'Barra'))

  function getTokenFromHash(): string | null {
    if (typeof window === 'undefined') return null
    const hash = window.location.hash.slice(1)
    if (!hash) return null
    const params = new URLSearchParams(hash)
    return params.get('token')
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
      orders = await getKitchenOrdersFromProjection(menuId, accessToken, stationFilter)
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

  const statusKey = (orderNumber: number) => `${orderNumber}-${station}`

  function getOrderStatus(orderNumber: number): 'pending' | 'preparing' | 'done' {
    return orderStatus[statusKey(orderNumber)] ?? 'pending'
  }

  function cycleOrderStatus(orderNumber: number) {
    const current = getOrderStatus(orderNumber)
    const next = current === 'pending' ? 'preparing' : current === 'preparing' ? 'done' : 'pending'
    orderStatus = { ...orderStatus, [statusKey(orderNumber)]: next }
  }

  async function startOrdersAndRealtime(cleanupRef: { intervalId: ReturnType<typeof setInterval> | null; unsub: (() => void) | null }): Promise<void> {
    const accessToken = token || (typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(STORAGE_TOKEN_KEY) : null)
    if (!accessToken) return
    await loadOrders()
    const unsub = await subscribeMenuOrdersRealtime(menuId, accessToken, () => loadOrders())
    realtimeUnsubscribe = unsub
    cleanupRef.unsub = unsub
    cleanupRef.intervalId = setInterval(() => loadOrders(), 20_000)
  }

  function submitOperatorName() {
    const name = operatorName.trim()
    if (!name) return
    if (typeof sessionStorage !== 'undefined') sessionStorage.setItem(STORAGE_OPERATOR_KEY, name)
    operatorSubmitted = true
    startOrdersAndRealtime(cleanupRef)
  }

  onMount(() => {
    if (typeof window === 'undefined') return
    let fromHash = getTokenFromHash()
    if (fromHash) {
      sessionStorage.setItem(STORAGE_TOKEN_KEY, fromHash)
      token = fromHash
      clearHash()
    } else {
      token = sessionStorage.getItem(STORAGE_TOKEN_KEY)
    }

    const savedOperator = sessionStorage.getItem(STORAGE_OPERATOR_KEY)
    if (savedOperator) {
      operatorName = savedOperator
      operatorSubmitted = true
    }

    if (operatorSubmitted && token) {
      startOrdersAndRealtime(cleanupRef)
    }

    return () => {
      if (cleanupRef.intervalId) clearInterval(cleanupRef.intervalId)
      cleanupRef.unsub?.()
      realtimeUnsubscribe = null
    }
  })
</script>

<div class="min-h-screen flex flex-col bg-gray-50">
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
    <div class="flex-shrink-0 px-4 py-4 border-b border-gray-200 bg-white">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800">
        {stationLabel} — {operatorName || 'Operador'}
      </h1>
      <p class="text-sm text-gray-500 mt-1">Pedidos en tiempo real. Sin login.</p>
    </div>

    <!-- Lista de órdenes -->
    <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
      {#if loading && orders.length === 0}
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-10 w-10 border-2 border-amber-500 border-t-transparent"></div>
        </div>
      {:else if error}
        <div class="rounded-lg bg-red-50 border border-red-200 p-4 text-red-800">{error}</div>
      {:else if displayedOrders.length === 0}
        <div class="rounded-lg bg-gray-100 border border-gray-200 p-8 text-center text-gray-600">
          {t.orders?.emptyForStation ?? 'No hay órdenes para esta estación.'}
        </div>
      {:else}
        <ul class="space-y-5">
          {#each displayedOrders as order, index (order.order_number)}
            {@const type = order.fulfillment}
            {@const itemCount = getItemCount(order.items)}
            {@const isFirst = index === 0}
            {@const status = getOrderStatus(order.order_number)}
            {@const remainingMin = getRemainingMinutes(order.requested_time)}
            {@const timeColor = getRemainingTimeColor(remainingMin)}
            {@const useBarColor = station === 'BAR'}
            <li class="bg-white rounded-xl border-2 overflow-hidden {isFirst ? 'border-amber-400 shadow-lg' : 'border-gray-200'}">
              <div class="w-full px-4 py-3 sm:px-5 flex flex-wrap items-center gap-4 border-b border-gray-100">
                <span class="font-bold text-gray-900 tabular-nums {isFirst ? 'text-4xl' : 'text-3xl'}">#{order.order_number}</span>
                <span class="font-semibold text-gray-700">{t.orders?.forTime ?? 'Para'} {formatRequestedTime(order.requested_time)}</span>
                {#if remainingMin !== null}
                  <span class="inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-sm font-bold tabular-nums
                    {timeColor === 'green' ? 'bg-green-100 text-green-800' : ''}
                    {timeColor === 'yellow' ? 'bg-amber-200 text-amber-900' : ''}
                    {timeColor === 'red' ? 'bg-red-100 text-red-800' : ''}">
                    {getRemainingTimeLabel(remainingMin)}
                  </span>
                {/if}
                <span class="inline-flex items-center rounded-full font-medium {type === 'DELIVERY' ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'} px-3 py-1 text-sm">
                  {getFulfillmentLabel(type)}
                </span>
                <span class="inline-flex items-center gap-1 rounded-full font-bold border px-3 py-1 text-sm
                  {status === 'pending' ? 'bg-gray-100 text-gray-700 border-gray-200' : ''}
                  {status === 'preparing' ? 'bg-amber-50 text-amber-900 border-amber-200' : ''}
                  {status === 'done' ? 'bg-green-50 text-green-800 border-green-200' : ''}">
                  {status === 'pending' ? (t.orders?.statusPending ?? 'Pendiente') : status === 'preparing' ? (t.orders?.statusPreparing ?? 'En preparación') : (t.orders?.statusDone ?? 'Listo')}
                </span>
              </div>
              <div class="px-4 py-3 sm:px-5 bg-amber-50/50 border-b border-amber-100">
                <p class="font-semibold text-amber-800 uppercase tracking-wide text-xs mb-2">{t.orders?.itemsToPrepare ?? 'Qué preparar'}</p>
                <ul class="space-y-1 text-gray-900 text-lg">
                  {#each order.items as item}
                    <li><span class="font-bold text-amber-800">{item.quantity}×</span> {item.item_name}</li>
                  {/each}
                </ul>
              </div>
              <div class="px-4 py-3 sm:px-5 border-t border-gray-100">
                {#if status === 'pending'}
                  <button
                    type="button"
                    onclick={() => cycleOrderStatus(order.order_number)}
                    class="w-full py-3 px-4 rounded-xl text-base font-bold text-white shadow-md transition-colors {useBarColor ? 'bg-blue-600 hover:bg-blue-700' : 'bg-amber-500 hover:bg-amber-600'}"
                  >
                    {t.orders?.startPreparing ?? 'INICIAR'}
                  </button>
                {:else if status === 'preparing'}
                  <button
                    type="button"
                    onclick={() => cycleOrderStatus(order.order_number)}
                    class="w-full py-3 px-4 rounded-xl text-base font-bold text-white shadow-md {useBarColor ? 'bg-blue-500 hover:bg-blue-600' : 'bg-amber-500 hover:bg-amber-600'}"
                  >
                    ✓ {t.orders?.markAsReady ?? 'LISTO'}
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
  {/if}
</div>
