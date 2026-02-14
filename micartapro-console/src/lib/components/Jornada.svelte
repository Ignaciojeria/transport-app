<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getKitchenOrdersFromProjection, type KitchenOrder } from '../menuUtils'
  import { getActiveJourney, createJourney, closeJourney, getJourneys, getJourneyStats, type JourneyListItem, type JourneyStats } from '../journeyApi'
  import { t as tStore } from '../useLanguage'

  interface JornadaProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: JornadaProps = $props()

  let menuId = $state<string | null>(null)
  let orders = $state<KitchenOrder[]>([])
  let activeJourney = $state<{ id: string } | null>(null)
  let loading = $state(true)
  let error = $state<string | null>(null)
  let showCloseModal = $state(false)
  let closeInProgress = $state(false)
  let closeResult = $state<'idle' | 'success' | 'coming_soon'>('idle')
  let createJourneyInProgress = $state(false)
  let createJourneyError = $state<string | null>(null)
  let journeys = $state<JourneyListItem[]>([])
  let showStatsModal = $state(false)
  let statsJourney = $state<JourneyListItem | null>(null)
  let stats = $state<JourneyStats | null>(null)
  let statsLoading = $state(false)

  const session = $derived(authState.session)
  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const t = $derived($tStore)

  /** Orden totalmente entregada/despachada (estado terminal). */
  function isOrderFullyDelivered(order: KitchenOrder): boolean {
    const active = order.items.filter((i) => i.status !== 'CANCELLED')
    return active.length > 0 && active.every((i) => i.status === 'DISPATCHED' || i.status === 'DELIVERED')
  }

  /** Orden totalmente cancelada. */
  function isOrderFullyCancelled(order: KitchenOrder): boolean {
    return order.items.length > 0 && order.items.every((i) => i.status === 'CANCELLED')
  }

  /** Fecha de hoy en zona local YYYY-MM-DD. */
  function todayLocalDateString(): string {
    const d = new Date()
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
  }

  /** Formato DD-MM HH:mm para jornadas. */
  function formatJourneyDate(isoDate: string): string {
    const d = new Date(isoDate)
    const day = String(d.getDate()).padStart(2, '0')
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const h = String(d.getHours()).padStart(2, '0')
    const min = String(d.getMinutes()).padStart(2, '0')
    return `${day}/${m} ${h}:${min}`
  }

  /** Formato DD-MM-YYYY para mostrar en modal. */
  function formatDateDDMMYYYY(isoDate: string): string {
    const d = new Date(isoDate + 'T12:00:00')
    const day = String(d.getDate()).padStart(2, '0')
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const y = d.getFullYear()
    return `${day}-${m}-${y}`
  }

  /** Resumen rápido: todas las órdenes de la jornada activa (ya filtradas por journey_id al cargar). */
  const total = $derived(orders.length)
  const delivered = $derived(orders.filter((o) => isOrderFullyDelivered(o)).length)
  const cancelled = $derived(orders.filter((o) => isOrderFullyCancelled(o)).length)
  const pending = $derived(orders.filter((o) => !isOrderFullyDelivered(o) && !isOrderFullyCancelled(o)).length)

  /** Hora de apertura: mínima created_at de las órdenes de la jornada, formateada HH:mm. */
  const openedSince = $derived.by(() => {
    if (orders.length === 0) return null
    const minCreated = orders.reduce((acc, o) => (o.created_at < acc ? o.created_at : acc), orders[0].created_at)
    const d = new Date(minCreated)
    return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  })

  const dateLabel = $derived(formatDateDDMMYYYY(todayLocalDateString()))
  const modalMessage = $derived((t.jornada?.closeModalMessage ?? 'Estás por cerrar la jornada del {date}. Las órdenes pendientes quedarán marcadas como no entregadas. ¿Deseas continuar?').replace('{date}', dateLabel))

  async function load() {
    if (!userId || !session?.access_token) {
      error = t.jornada?.noSession ?? 'No Hay Sesión Activa'
      loading = false
      return
    }
    try {
      loading = true
      error = null
      createJourneyError = null
      const mid = await getLatestMenuId(userId, session.access_token)
      if (!mid) {
        error = t.jornada?.noMenu ?? 'No Se Encontró Un Menú'
        loading = false
        return
      }
      menuId = mid
      const [journey, journeysList] = await Promise.all([
        getActiveJourney(mid, session.access_token),
        getJourneys(mid, session.access_token)
      ])
      activeJourney = journey
      journeys = journeysList
      const list = journey
        ? await getKitchenOrdersFromProjection(mid, session.access_token, 'ALL', journey.id)
        : []
      orders = list
    } catch (e) {
      console.error('Error cargando jornada:', e)
      error = t.jornada?.errorLoading ?? 'Error Al Cargar Los Datos De La Jornada.'
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
      await load()
    } catch (e) {
      console.error('Error creando jornada:', e)
      createJourneyError = t.jornada?.errorCreatingJourney ?? 'Error al abrir la jornada. Intenta de nuevo.'
    } finally {
      createJourneyInProgress = false
    }
  }

  function openCloseModal() {
    closeResult = 'idle'
    showCloseModal = true
  }

  function cancelCloseModal() {
    if (!closeInProgress) showCloseModal = false
  }

  async function confirmCloseWorkday() {
    if (!session?.access_token || !menuId) return
    closeInProgress = true
    closeResult = 'idle'
    try {
      await closeJourney(menuId, session.access_token)
      closeResult = 'success'
      await load()
      setTimeout(() => {
        showCloseModal = false
        closeResult = 'idle'
      }, 1500)
    } catch {
      closeResult = 'coming_soon'
    } finally {
      closeInProgress = false
    }
  }

  async function openStatsModal(j: JourneyListItem) {
    if (!session?.access_token || !menuId) return
    statsJourney = j
    showStatsModal = true
    stats = null
    statsLoading = true
    try {
      stats = await getJourneyStats(menuId, j.id, session.access_token)
    } catch (e) {
      console.error('Error cargando estadísticas:', e)
      stats = null
    } finally {
      statsLoading = false
    }
  }

  function closeStatsModal() {
    showStatsModal = false
    statsJourney = null
    stats = null
  }

  /** Formato de moneda para mostrar revenue. */
  function formatCurrency(n: number): string {
    if (n >= 1000) return `$${Math.round(n).toLocaleString()}`
    return `$${n.toFixed(0)}`
  }

  /** Colores para el gráfico de torta (por índice). */
  const PIE_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16']

  /** Gradiente conic para el gráfico de torta por revenue. */
  const pieGradient = $derived.by(() => {
    if (!stats?.products?.length) return ''
    let acc = 0
    return stats.products
      .map((p, i) => {
        const start = acc
        acc += p.percentage
        return `${PIE_COLORS[i % PIE_COLORS.length]} ${start}% ${acc}%`
      })
      .join(', ')
  })

  /** Gradiente conic para el gráfico de torta por cantidad. */
  const pieGradientByQuantity = $derived.by(() => {
    if (!stats?.products?.length) return ''
    let acc = 0
    return stats.products
      .map((p, i) => {
        const start = acc
        acc += p.percentageByQuantity
        return `${PIE_COLORS[i % PIE_COLORS.length]} ${start}% ${acc}%`
      })
      .join(', ')
  })

  /** Ticket promedio (ventas / órdenes). */
  const averageTicket = $derived.by(() => {
    if (!stats || stats.totalOrders === 0) return 0
    return stats.totalRevenue / stats.totalOrders
  })

  /** Producto top por revenue. */
  const topByRevenue = $derived.by(() => {
    if (!stats?.products?.length) return null
    return stats.products.reduce((a, b) => (a.totalRevenue >= b.totalRevenue ? a : b))
  })

  /** Producto top por cantidad vendida. */
  const topByQuantity = $derived.by(() => {
    if (!stats?.products?.length) return null
    return stats.products.reduce((a, b) => (a.quantitySold >= b.quantitySold ? a : b))
  })

  onMount(() => {
    load()
  })
</script>

<div class="h-full flex flex-col bg-gray-50">
  <!-- Header con menú -->
  <header class="flex-shrink-0 flex items-center gap-3 px-4 py-3 bg-white border-b border-gray-200 md:px-6">
    <button
      type="button"
      onclick={onMenuClick}
      class="p-2 rounded-lg hover:bg-gray-100 md:hidden"
      aria-label={t.sidebar?.closeMenu ?? 'Menú'}
    >
      <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>
    <h1 class="text-lg font-semibold text-gray-900">{t.jornada?.title ?? 'Jornada'}</h1>
  </header>

  <div class="flex-1 overflow-y-auto p-4 md:p-6">
    {#if loading}
      <div class="flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-10 w-10 border-2 border-blue-600 border-t-transparent"></div>
      </div>
    {:else if error}
      <div class="rounded-xl bg-amber-50 border border-amber-200 p-4 text-amber-800">
        <p class="font-medium">{error}</p>
      </div>
    {:else if activeJourney === null && menuId}
      <!-- No hay jornada activa: CTA para abrir una -->
      {@const closedJourneysNoActive = journeys.filter((j) => j.status === 'CLOSED')}
      {#if closedJourneysNoActive.length > 0}
        <section class="mb-8">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">{t.jornada?.reports ?? 'Reportes de Jornadas'}</h2>
          <div class="space-y-2">
            {#each closedJourneysNoActive as j (j.id)}
              <div class="flex items-center justify-between rounded-xl bg-white border border-gray-200 p-4">
                <div class="text-sm">
                  <span class="font-medium text-gray-900">{formatJourneyDate(j.openedAt)}</span>
                  {#if j.closedAt}
                    <span class="text-gray-500"> – {formatJourneyDate(j.closedAt)}</span>
                  {/if}
                </div>
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    onclick={() => openStatsModal(j)}
                    class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                    {t.jornada?.stats ?? 'Estadísticas'}
                  </button>
                  {#if j.reportXlsxUrl}
                    <a
                      href={j.reportXlsxUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-white bg-green-600 hover:bg-green-700 transition-colors"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      {t.jornada?.downloadExcel ?? 'Descargar Excel'}
                    </a>
                  {:else}
                    <span class="text-xs text-gray-400">{t.jornada?.reportGenerating ?? 'Generando...'}</span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </section>
      {/if}
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
    {:else}
      <!-- Jornada activa: estado actual, resumen y cerrar -->
      <section class="mb-6">
        <div class="rounded-xl bg-white border border-gray-200 shadow-sm p-5">
          <div class="flex items-center gap-2 mb-3">
            <span class="inline-flex w-3 h-3 rounded-full bg-green-500" aria-hidden="true"></span>
            <h2 class="text-sm font-semibold text-gray-700">{t.jornada?.active ?? 'Jornada Activa'}</h2>
          </div>
          <dl class="grid grid-cols-1 gap-2 text-sm">
            <div>
              <dt class="text-gray-500">{t.jornada?.date ?? 'Fecha'}</dt>
              <dd class="font-medium text-gray-900">{dateLabel}</dd>
            </div>
            <div>
              <dt class="text-gray-500">{t.jornada?.openedSince ?? 'Abierta Desde'}</dt>
              <dd class="font-medium text-gray-900">{openedSince ?? '—'}</dd>
            </div>
          </dl>
        </div>
      </section>

      <section class="mb-8">
        <h2 class="text-sm font-semibold text-gray-700 mb-3">{t.jornada?.summary ?? 'Resumen Rápido'}</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
          <div class="rounded-xl bg-white border border-gray-200 p-4 text-center">
            <p class="text-2xl font-bold text-gray-900">{total}</p>
            <p class="text-xs text-gray-500">{t.jornada?.totalOrders ?? 'Órdenes Totales'}</p>
          </div>
          <div class="rounded-xl bg-white border border-gray-200 p-4 text-center">
            <p class="text-2xl font-bold text-green-600">{delivered}</p>
            <p class="text-xs text-gray-500">{t.jornada?.delivered ?? 'Entregadas'}</p>
          </div>
          <div class="rounded-xl bg-white border border-gray-200 p-4 text-center">
            <p class="text-2xl font-bold text-gray-600">{cancelled}</p>
            <p class="text-xs text-gray-500">{t.jornada?.cancelled ?? 'Canceladas'}</p>
          </div>
          <div class="rounded-xl bg-white border border-amber-200 bg-amber-50/50 p-4 text-center">
            <p class="text-2xl font-bold text-amber-700">{pending}</p>
            <p class="text-xs text-gray-500">{t.jornada?.pending ?? 'Pendientes'}</p>
            {#if pending > 0}
              <p class="text-[10px] text-amber-600 mt-0.5">⚠</p>
            {/if}
          </div>
        </div>
      </section>

      <section class="mt-8 pt-6 border-t border-gray-200">
        <button
          type="button"
          onclick={openCloseModal}
          class="w-full max-w-sm mx-auto flex items-center justify-center gap-2 py-4 px-6 rounded-xl font-semibold text-white bg-gray-800 hover:bg-gray-900 focus:ring-2 focus:ring-offset-2 focus:ring-gray-700 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          <span>{t.jornada?.closeWorkday ?? 'Cerrar Jornada'}</span>
        </button>
      </section>

      <!-- Reportes de jornadas cerradas -->
      {@const closedJourneys = journeys.filter((j) => j.status === 'CLOSED')}
      {#if closedJourneys.length > 0}
        <section class="mt-10 pt-6 border-t border-gray-200">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">{t.jornada?.reports ?? 'Reportes de Jornadas'}</h2>
          <div class="space-y-2">
            {#each closedJourneys as j (j.id)}
              <div class="flex items-center justify-between rounded-xl bg-white border border-gray-200 p-4">
                <div class="text-sm">
                  <span class="font-medium text-gray-900">{formatJourneyDate(j.openedAt)}</span>
                  {#if j.closedAt}
                    <span class="text-gray-500"> – {formatJourneyDate(j.closedAt)}</span>
                  {/if}
                </div>
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    onclick={() => openStatsModal(j)}
                    class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                    {t.jornada?.stats ?? 'Estadísticas'}
                  </button>
                  {#if j.reportXlsxUrl}
                    <a
                      href={j.reportXlsxUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium text-white bg-green-600 hover:bg-green-700 transition-colors"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      {t.jornada?.downloadExcel ?? 'Descargar Excel'}
                    </a>
                  {:else}
                    <span class="text-xs text-gray-400">{t.jornada?.reportGenerating ?? 'Generando...'}</span>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </section>
      {/if}
    {/if}
  </div>
</div>

<!-- Modal Cerrar Jornada -->
{#if showCloseModal}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
    role="dialog"
    aria-modal="true"
    aria-labelledby="jornada-close-title"
    tabindex="-1"
    onclick={(e) => e.target === e.currentTarget && !closeInProgress && cancelCloseModal()}
    onkeydown={(e) => e.key === 'Escape' && !closeInProgress && cancelCloseModal()}
  >
    <div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6">
      <h2 id="jornada-close-title" class="text-xl font-bold text-gray-900 mb-3">
        {t.jornada?.closeModalTitle ?? 'Cerrar Jornada'}
      </h2>
      <p class="text-gray-600 text-sm mb-6">{modalMessage}</p>

      {#if closeResult === 'success'}
        <p class="text-green-600 font-medium mb-4">{t.jornada?.success ?? 'Jornada Cerrada Correctamente.'}</p>
      {:else if closeResult === 'coming_soon'}
        <p class="text-amber-700 text-sm mb-4">{t.jornada?.comingSoon ?? 'El Cierre De Jornada Estará Disponible En Una Próxima Actualización.'}</p>
      {/if}

      <div class="flex gap-3 justify-end">
        <button
          type="button"
          disabled={closeInProgress}
          onclick={cancelCloseModal}
          class="px-4 py-2 rounded-lg font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 disabled:opacity-50"
        >
          {t.jornada?.closeModalCancel ?? 'Cancelar'}
        </button>
        <button
          type="button"
          disabled={closeInProgress}
          onclick={confirmCloseWorkday}
          class="px-4 py-2 rounded-lg font-medium text-white bg-gray-800 hover:bg-gray-900 disabled:opacity-50"
        >
          {closeInProgress ? (t.jornada?.closing ?? 'Cerrando...') : (t.jornada?.closeModalConfirm ?? 'Cerrar Jornada')}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal Estadísticas de Jornada -->
{#if showStatsModal}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
    role="dialog"
    aria-modal="true"
    aria-labelledby="jornada-stats-title"
    tabindex="-1"
    onclick={(e) => e.target === e.currentTarget && closeStatsModal()}
    onkeydown={(e) => e.key === 'Escape' && closeStatsModal()}
  >
    <div class="bg-white rounded-2xl shadow-xl max-w-lg w-full p-6 max-h-[90vh] overflow-y-auto">
      <div class="flex items-center justify-between mb-4">
        <h2 id="jornada-stats-title" class="text-xl font-bold text-gray-900">
          {t.jornada?.stats ?? 'Estadísticas'}
          {#if statsJourney}
            <span class="text-sm font-normal text-gray-500 ml-2">
              {formatJourneyDate(statsJourney.openedAt)}
            </span>
          {/if}
        </h2>
        <button
          type="button"
          onclick={closeStatsModal}
          class="p-2 rounded-lg hover:bg-gray-100 text-gray-500"
          aria-label="Cerrar"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      {#if statsLoading}
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-10 w-10 border-2 border-blue-600 border-t-transparent"></div>
        </div>
      {:else if stats}
        <div class="space-y-6">
          <!-- Métricas rápidas -->
          <div class="grid grid-cols-2 gap-3">
            <div class="rounded-xl bg-gray-50 p-4">
              <p class="text-2xl font-bold text-gray-900">{formatCurrency(stats.totalRevenue)}</p>
              <p class="text-xs text-gray-500">{t.jornada?.totalRevenue ?? 'Ventas totales'}</p>
            </div>
            <div class="rounded-xl bg-gray-50 p-4">
              <p class="text-2xl font-bold text-gray-900">{stats.totalOrders}</p>
              <p class="text-xs text-gray-500">{t.jornada?.statsOrders ?? 'Órdenes'}</p>
            </div>
            <div class="rounded-xl bg-gray-50 p-4">
              <p class="text-2xl font-bold text-gray-900">{formatCurrency(averageTicket)}</p>
              <p class="text-xs text-gray-500">{t.jornada?.averageTicket ?? 'Ticket promedio'}</p>
            </div>
            <div class="rounded-xl bg-blue-50 p-4 border border-blue-100">
              <p class="text-sm font-semibold text-gray-900 truncate" title={topByRevenue?.productName}>
                {topByRevenue?.productName ?? '—'}
              </p>
              <p class="text-xs text-blue-600">{t.jornada?.topByRevenue ?? 'Top ventas'}</p>
            </div>
            <div class="rounded-xl bg-emerald-50 p-4 border border-emerald-100">
              <p class="text-sm font-semibold text-gray-900 truncate" title={topByQuantity?.productName}>
                {topByQuantity?.productName ?? '—'}
              </p>
              <p class="text-xs text-emerald-600">{t.jornada?.topByQuantity ?? 'Top unidades'}</p>
            </div>
          </div>

          {#if stats.products.length > 0}
            <!-- Gráficos de torta: por ventas y por unidades -->
            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col items-center">
                <p class="text-xs font-medium text-gray-600 mb-2">{t.jornada?.chartByRevenue ?? 'Por ventas'}</p>
                <div
                  class="w-32 h-32 rounded-full shrink-0"
                  style="background: conic-gradient({pieGradient})"
                  role="img"
                  aria-label="Gráfico por ventas"
                ></div>
              </div>
              <div class="flex flex-col items-center">
                <p class="text-xs font-medium text-gray-600 mb-2">{t.jornada?.chartByQuantity ?? 'Por unidades'}</p>
                <div
                  class="w-32 h-32 rounded-full shrink-0"
                  style="background: conic-gradient({pieGradientByQuantity})"
                  role="img"
                  aria-label="Gráfico por unidades vendidas"
                ></div>
              </div>
            </div>

            <!-- Lista de productos -->
            <div>
              <h3 class="text-sm font-semibold text-gray-700 mb-2">{t.jornada?.topProducts ?? 'Productos más vendidos'}</h3>
              <ul class="space-y-2">
                {#each stats.products as p, i}
                  <li class="flex items-center justify-between rounded-lg bg-gray-50 px-3 py-2">
                    <span class="flex items-center gap-2">
                      <span
                        class="w-3 h-3 rounded-full shrink-0"
                        style="background-color: {PIE_COLORS[i % PIE_COLORS.length]}"
                      ></span>
                      <span class="font-medium text-gray-900">{p.productName}</span>
                    </span>
                    <span class="text-sm text-gray-600">
                      {p.quantitySold} ud · {formatCurrency(p.totalRevenue)} ({p.percentage.toFixed(0)}% ventas / {p.percentageByQuantity.toFixed(0)}% ud)
                    </span>
                  </li>
                {/each}
              </ul>
            </div>
          {:else}
            <p class="text-gray-500 text-center py-6">{t.jornada?.noStats ?? 'No hay datos de ventas para esta jornada.'}</p>
          {/if}
        </div>
      {:else}
        <p class="text-gray-500 text-center py-6">{t.jornada?.errorLoadingStats ?? 'Error al cargar estadísticas.'}</p>
      {/if}
    </div>
  </div>
{/if}
