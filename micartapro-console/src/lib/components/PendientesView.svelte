<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { getPendingOrders, assignOrdersToJourney, cancelOrder, type PendingOrder, type PendingOrdersFilter } from '../orderApi'
  import { t as tStore } from '../useLanguage'

  interface PendientesViewProps {
    menuId: string
    accessToken: string
    onAssigned?: () => void
  }

  let { menuId, accessToken, onAssigned }: PendientesViewProps = $props()

  let orders = $state<PendingOrder[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let fromDate = $state<string>(getTodayLocal())
  let toDate = $state<string>(getTodayLocal())
  let selectedIds = $state<Set<number>>(new Set())
  let assignInProgress = $state(false)
  let cancelInProgress = $state(false)
  let actionError = $state<string | null>(null)
  let orderToCancel = $state<PendingOrder | null>(null)
  let cancelReason = $state('')
  let cancelComment = $state('')
  let orderToView = $state<PendingOrder | null>(null)

  const t = $derived($tStore)

  const CANCEL_REASON_KEYS = ['outOfStock', 'orderError', 'customerLeft', 'paymentIssue', 'other'] as const

  function getTodayLocal(): string {
    const d = new Date()
    return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0')
  }

  async function loadPendingOrders() {
    try {
      loading = true
      error = null
      const filter: PendingOrdersFilter = {}
      // UTC ISO-8601 explícito: rango de fechas en UTC (YYYY-MM-DD → 00:00:00Z y 23:59:59.999Z)
      if (fromDate) filter.fromDate = `${fromDate}T00:00:00.000Z`
      if (toDate) filter.toDate = `${toDate}T23:59:59.999Z`
      orders = await getPendingOrders(menuId, accessToken, filter)
    } catch (err) {
      console.error('Error cargando pendientes:', err)
      error = err instanceof Error ? err.message : 'Error al cargar órdenes pendientes'
      orders = []
    } finally {
      loading = false
    }
  }

  function formatDate(iso: string | null | undefined): string {
    if (!iso) return '—'
    const d = new Date(iso)
    return d.toLocaleString('es-CL', { dateStyle: 'short', timeStyle: 'short' })
  }

  /** "Atrasada 5 min" | "En 25 min" | "Ahora" | "17 Feb 15:00" — calculado en frontend para actualización en tiempo real */
  function formatTiempoDisplay(scheduledFor: string | null | undefined, createdAt: string): string {
    const now = new Date()
    let t: Date | null = null
    if (scheduledFor) {
      const parsed = new Date(scheduledFor)
      if (!isNaN(parsed.getTime())) t = parsed
    }
    if (!t) {
      const created = new Date(createdAt)
      if (!isNaN(created.getTime())) {
        return created.toLocaleDateString('es-CL', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
      }
      return '—'
    }
    const diffMin = Math.round((t.getTime() - now.getTime()) / 60_000)
    if (diffMin < -5) return `Atrasada ${-diffMin} min`
    if (diffMin <= 5) return 'Ahora'
    if (diffMin <= 60) return `En ${diffMin} min`
    return t.toLocaleDateString('es-CL', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
  }

  let tiempoTick = $state(0)
  let tiempoInterval: ReturnType<typeof setInterval> | null = null

  function formatCurrency(amount: number): string {
    return `$${amount.toLocaleString('es-CL')}`
  }

  function toggleSelect(id: number) {
    const next = new Set(selectedIds)
    if (next.has(id)) next.delete(id)
    else next.add(id)
    selectedIds = next
  }

  function toggleSelectAll() {
    if (selectedIds.size === orders.length) {
      selectedIds = new Set()
    } else {
      selectedIds = new Set(orders.map((o) => o.aggregateId))
    }
  }

  const allSelected = $derived(orders.length > 0 && selectedIds.size === orders.length)
  const someSelected = $derived(selectedIds.size > 0)

  async function handleAssign(idsOverride?: number[]) {
    const ids = idsOverride ?? [...selectedIds]
    if (ids.length === 0 || assignInProgress) return
    try {
      assignInProgress = true
      actionError = null
      await assignOrdersToJourney(menuId, ids, accessToken)
      if (!idsOverride) selectedIds = new Set()
      await loadPendingOrders()
      onAssigned?.()
    } catch (err) {
      actionError = err instanceof Error ? err.message : 'Error al asignar órdenes'
    } finally {
      assignInProgress = false
    }
  }

  function openViewModal(order: PendingOrder) {
    orderToView = order
  }

  function openCancelModal(order: PendingOrder) {
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

  async function confirmCancel() {
    if (!orderToCancel || !cancelReason.trim()) return
    cancelInProgress = true
    actionError = null
    try {
      const reasonText = cancelComment.trim()
        ? `${getCancelReasonLabel(cancelReason)}: ${cancelComment.trim()}`
        : getCancelReasonLabel(cancelReason)
      await cancelOrder(menuId, orderToCancel.aggregateId, accessToken, reasonText)
      closeCancelModal()
      await loadPendingOrders()
    } catch (err) {
      actionError = err instanceof Error ? err.message : 'Error al cancelar'
    } finally {
      cancelInProgress = false
    }
  }

  async function cancelSelected() {
    if (!someSelected || cancelInProgress) return
    const ids = [...selectedIds]
    cancelInProgress = true
    actionError = null
    try {
      for (const id of ids) {
        await cancelOrder(menuId, id, accessToken, 'Cancelación grupal')
      }
      selectedIds = new Set()
      await loadPendingOrders()
    } catch (err) {
      actionError = err instanceof Error ? err.message : 'Error al cancelar órdenes'
    } finally {
      cancelInProgress = false
    }
  }

  onMount(() => {
    loadPendingOrders()
    tiempoInterval = setInterval(() => {
      tiempoTick++
    }, 60_000)
  })

  onDestroy(() => {
    if (tiempoInterval) clearInterval(tiempoInterval)
  })
</script>

<div class="flex flex-col gap-4">
  <!-- Filtros -->
  <div class="flex flex-wrap items-center gap-3">
    <div class="flex items-center gap-2">
      <label class="text-xs font-medium text-gray-600">{t.orders?.pendingFrom ?? 'Desde'}</label>
      <input
        type="date"
        bind:value={fromDate}
        onchange={loadPendingOrders}
        class="rounded-lg border border-gray-300 px-2 py-1.5 text-sm"
      />
    </div>
    <div class="flex items-center gap-2">
      <label class="text-xs font-medium text-gray-600">{t.orders?.pendingTo ?? 'Hasta'}</label>
      <input
        type="date"
        bind:value={toDate}
        onchange={loadPendingOrders}
        class="rounded-lg border border-gray-300 px-2 py-1.5 text-sm"
      />
    </div>
    <button
      type="button"
      onclick={loadPendingOrders}
      class="rounded-lg bg-gray-100 px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-200"
    >
      {t.orders?.reload ?? 'Recargar'}
    </button>
  </div>

  {#if actionError}
    <div class="rounded-lg bg-amber-50 border border-amber-200 p-3 text-sm text-amber-800 flex items-center justify-between">
      <span>{actionError}</span>
      <button type="button" onclick={() => (actionError = null)} class="underline hover:no-underline">Cerrar</button>
    </div>
  {/if}

  {#if someSelected}
    <!-- Barra de acciones grupales -->
    <div class="rounded-lg bg-blue-50 border border-blue-200 p-3 flex flex-wrap items-center gap-3">
      <span class="font-semibold text-blue-900">
        {selectedIds.size} {t.orders?.pendingSelected ?? 'órdenes seleccionadas'}
      </span>
      <button
        type="button"
        disabled={assignInProgress}
        onclick={() => handleAssign()}
        class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
      >
        {#if assignInProgress}
          <span class="animate-spin inline-block mr-1">⏳</span>
        {/if}
        {t.orders?.pendingAssignToJourney ?? 'Asignar a jornada activa'}
      </button>
      <button
        type="button"
        disabled={cancelInProgress}
        onclick={cancelSelected}
        class="rounded-lg bg-red-600 px-4 py-2 text-sm font-semibold text-white hover:bg-red-700 disabled:opacity-50"
      >
        {t.orders?.pendingCancelOrders ?? 'Cancelar órdenes'}
      </button>
      <button
        type="button"
        onclick={() => (selectedIds = new Set())}
        class="rounded-lg bg-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-300"
      >
        {t.orders?.pendingClearSelection ?? 'Limpiar selección'}
      </button>
    </div>
  {/if}

  <!-- Tabla -->
  {#if loading}
    <div class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-2 border-blue-600 border-t-transparent"></div>
    </div>
  {:else if error}
    <div class="rounded-lg bg-red-50 border border-red-200 p-4 text-sm text-red-800">{error}</div>
  {:else if orders.length === 0}
    <p class="text-gray-500 text-center py-8">{t.orders?.pendingEmpty ?? 'No hay órdenes pendientes.'}</p>
  {:else}
    <div class="overflow-x-auto rounded-lg border border-gray-200 bg-white">
      <table class="min-w-full text-sm">
        <thead class="bg-gray-50 border-b border-gray-200">
          <tr>
            <th class="px-3 py-2 text-left">
              <input
                type="checkbox"
                checked={allSelected}
                onchange={toggleSelectAll}
                aria-label={t.orders?.pendingSelectAll ?? 'Seleccionar todas'}
                class="rounded"
              />
            </th>
            <th class="px-3 py-2 text-left font-semibold text-gray-700">Tracking</th>
            <th class="px-3 py-2 text-left font-semibold text-gray-700">Tiempo</th>
            <th class="px-3 py-2 text-left font-semibold text-gray-700">Creado</th>
            <th class="px-3 py-2 text-left font-semibold text-gray-700">Programado</th>
            <th class="px-3 py-2 text-right font-semibold text-gray-700">Total</th>
            <th class="px-3 py-2 text-left font-semibold text-gray-700">Estado</th>
            <th class="px-3 py-2 text-right font-semibold text-gray-700">Acciones</th>
          </tr>
        </thead>
        <tbody data-tick={tiempoTick}>
          {#each orders as order}
            {@const tiempo = formatTiempoDisplay(order.scheduledFor, order.createdAt)}
            <tr class="border-b border-gray-100 hover:bg-gray-50">
              <td class="px-3 py-2">
                <input
                  type="checkbox"
                  checked={selectedIds.has(order.aggregateId)}
                  onchange={() => toggleSelect(order.aggregateId)}
                  class="rounded"
                />
              </td>
              <td class="px-3 py-2 font-mono text-gray-900">{order.trackingId || order.aggregateId}</td>
              <td class="px-3 py-2 {tiempo.startsWith('Atrasada') ? 'text-red-600 font-medium' : 'text-gray-700'}">{tiempo}</td>
              <td class="px-3 py-2 text-gray-700">{formatDate(order.createdAt)}</td>
              <td class="px-3 py-2 text-gray-700">{formatDate(order.scheduledFor)}</td>
              <td class="px-3 py-2 text-right font-semibold tabular-nums">{formatCurrency(order.totalAmount)}</td>
              <td class="px-3 py-2">
                <span class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium bg-green-100 text-green-800">
                  {order.status}
                </span>
              </td>
              <td class="px-3 py-2">
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    onclick={() => (orderToView = order)}
                    class="rounded px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50"
                  >
                    {t.orders?.pendingViewDetails ?? 'Ver'}
                  </button>
                  <button
                    type="button"
                    onclick={() => handleAssign([order.aggregateId])}
                    class="rounded px-2 py-1 text-xs font-medium text-green-600 hover:bg-green-50"
                  >
                    {t.orders?.pendingAssign ?? 'Asignar'}
                  </button>
                  <button
                    type="button"
                    onclick={() => openCancelModal(order)}
                    class="rounded px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50"
                  >
                    {t.orders?.cancelOrder ?? 'Cancelar'}
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<!-- Modal Ver detalles -->
{#if orderToView}
  <div
    class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
    role="dialog"
    onclick={() => (orderToView = null)}
  >
    <div
      class="bg-white rounded-xl shadow-xl max-w-md w-full p-6"
      onclick={(e) => e.stopPropagation()}
    >
      <h3 class="text-lg font-bold text-gray-900 mb-4">Detalle de orden</h3>
      <dl class="space-y-2 text-sm mb-4">
        <div><dt class="font-medium text-gray-500">ID</dt><dd class="font-mono">{orderToView.aggregateId}</dd></div>
        <div><dt class="font-medium text-gray-500">Creado</dt><dd>{formatDate(orderToView.createdAt)}</dd></div>
        <div><dt class="font-medium text-gray-500">Programado</dt><dd>{formatDate(orderToView.scheduledFor)}</dd></div>
        <div><dt class="font-medium text-gray-500">Total</dt><dd class="font-semibold">{formatCurrency(orderToView.totalAmount)}</dd></div>
        <div><dt class="font-medium text-gray-500">Estado</dt><dd>{orderToView.status}</dd></div>
      </dl>
      <div class="border-t border-gray-200 pt-4">
        <h4 class="font-semibold text-gray-700 mb-2">Ítems</h4>
        {#if !orderToView.items || orderToView.items.length === 0}
          <p class="text-gray-500 text-sm">Sin ítems</p>
        {:else}
          <ul class="space-y-2 text-sm">
            {#each orderToView.items as item}
              <li class="flex justify-between items-start py-1 border-b border-gray-100 last:border-0">
                <span class="text-gray-800">{item.itemName} × {item.quantity} {item.unit}</span>
                <span class="font-medium tabular-nums">{formatCurrency(item.totalPrice)}</span>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
      <button
        type="button"
        onclick={() => (orderToView = null)}
        class="mt-4 w-full rounded-lg bg-gray-200 py-2 text-sm font-medium hover:bg-gray-300"
      >
        Cerrar
      </button>
    </div>
  </div>
{/if}

<!-- Modal Cancelar -->
{#if orderToCancel}
  <div
    class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
    role="dialog"
    onclick={closeCancelModal}
  >
    <div
      class="bg-white rounded-xl shadow-xl max-w-md w-full p-6"
      onclick={(e) => e.stopPropagation()}
    >
      <h3 class="text-lg font-bold text-gray-900">{t.orders?.cancelModalTitle ?? 'Cancelar Pedido'}</h3>
      <p class="text-sm text-gray-600 mt-1">{t.orders?.cancelModalSubtitle ?? 'Esta acción no se puede deshacer.'}</p>
      <p class="text-sm font-mono mt-2">#{orderToCancel.aggregateId}</p>
      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700">{t.orders?.cancelModalReasonLabel ?? 'Motivo'}</label>
        <div class="mt-1 space-y-1">
          {#each CANCEL_REASON_KEYS as key}
            <label class="flex items-center gap-2">
              <input type="radio" bind:group={cancelReason} value={key} class="rounded" />
              <span>{getCancelReasonLabel(key)}</span>
            </label>
          {/each}
        </div>
      </div>
      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700">{t.orders?.cancelModalCommentLabel ?? 'Comentario (opcional)'}</label>
        <input
          type="text"
          bind:value={cancelComment}
          placeholder={t.orders?.cancelModalCommentPlaceholder ?? 'Ej: Cliente no contestó...'}
          class="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
        />
      </div>
      <div class="mt-6 flex gap-2">
        <button
          type="button"
          onclick={closeCancelModal}
          disabled={cancelInProgress}
          class="flex-1 rounded-lg border border-gray-300 py-2 text-sm font-medium hover:bg-gray-50 disabled:opacity-50"
        >
          {t.orders?.cancelModalBack ?? 'Volver'}
        </button>
        <button
          type="button"
          onclick={confirmCancel}
          disabled={!cancelReason.trim() || cancelInProgress}
          class="flex-1 rounded-lg bg-red-600 py-2 text-sm font-semibold text-white hover:bg-red-700 disabled:opacity-50"
        >
          {#if cancelInProgress}⏳ {/if}{t.orders?.cancelModalConfirm ?? 'Confirmar Cancelación'}
        </button>
      </div>
    </div>
  </div>
{/if}
