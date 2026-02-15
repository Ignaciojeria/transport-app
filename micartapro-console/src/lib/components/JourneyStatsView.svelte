<script lang="ts">
  import { t as tStore } from '../useLanguage'
  import type { JourneyStats } from '../journeyApi'

  export interface JourneyInfo {
    id: string
    openedAt: string
    closedAt?: string
    reportXlsxUrl?: string
  }

  interface JourneyStatsViewProps {
    stats: JourneyStats | null
    statsLoading: boolean
    journeyInfo: JourneyInfo | null
    onBack?: () => void
    backLabel?: string
    onDownloadCSV?: () => void
  }

  let { stats, statsLoading, journeyInfo, onBack, backLabel, onDownloadCSV }: JourneyStatsViewProps = $props()

  const t = $derived($tStore)

  const PIE_COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16']

  function formatCurrency(n: number): string {
    if (n >= 1000) return `$${Math.round(n).toLocaleString()}`
    return `$${n.toFixed(0)}`
  }

  function formatWorkdayDate(isoDate: string): string {
    const d = new Date(isoDate)
    const months = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic']
    return `${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`
  }

  function formatTime(isoDate: string): string {
    const d = new Date(isoDate)
    return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  }

  const journeyDurationMinutes = $derived.by(() => {
    if (!journeyInfo?.openedAt || !journeyInfo?.closedAt) return null
    const start = new Date(journeyInfo.openedAt).getTime()
    const end = new Date(journeyInfo.closedAt).getTime()
    return Math.round((end - start) / 60000)
  })

  const averageTicket = $derived.by(() =>
    stats && stats.totalOrders > 0 ? stats.totalRevenue / stats.totalOrders : 0
  )

  const topProduct = $derived.by(() =>
    stats?.products?.length
      ? stats.products.reduce((a, b) => (a.totalRevenue >= b.totalRevenue ? a : b))
      : null
  )

  const itemsOrdered = $derived.by(() => stats?.itemsOrdered ?? 0)
</script>

{#if onBack}
  <div class="mb-6">
    <button
      type="button"
      onclick={onBack}
      class="inline-flex items-center gap-2 text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      {backLabel ?? (t.jornada?.back ?? 'Volver')}
    </button>
  </div>
{/if}

{#if statsLoading}
  <div class="flex justify-center py-16">
    <div class="animate-spin rounded-full h-12 w-12 border-2 border-blue-600 border-t-transparent"></div>
  </div>
{:else if stats}
  <div class="space-y-8 w-full max-w-full">
    <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between lg:gap-8">
      {#if journeyInfo}
        <div class="rounded-xl bg-white border border-gray-200 p-6 flex-1 min-w-0">
          <h2 class="text-lg font-bold text-gray-900 mb-1">{t.jornada?.workdayReport ?? 'Reporte de Jornada'}</h2>
          <p class="text-xl font-semibold text-gray-800">{formatWorkdayDate(journeyInfo.openedAt)}</p>
          <p class="text-sm text-gray-600 mt-1">
            {formatTime(journeyInfo.openedAt)}
            {#if journeyInfo.closedAt}
              – {formatTime(journeyInfo.closedAt)}
              {#if journeyDurationMinutes != null}
                · {t.jornada?.duration ?? 'Duración'}: {journeyDurationMinutes} min
              {/if}
            {:else}
              – {t.jornada?.active ?? 'En curso'}
            {/if}
          </p>
        </div>
      {/if}
      <div class="flex flex-wrap gap-2 shrink-0 {journeyInfo ? 'lg:mt-6' : ''}">
        {#if journeyInfo?.reportXlsxUrl}
          <a
            href={journeyInfo.reportXlsxUrl}
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white bg-green-600 hover:bg-green-700 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            {t.jornada?.downloadExcel ?? 'Descargar Excel'}
          </a>
        {/if}
        {#if stats.products.length > 0 && onDownloadCSV}
          <button
            type="button"
            onclick={onDownloadCSV}
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            {t.jornada?.downloadCSV ?? 'Descargar CSV'}
          </button>
        {/if}
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
      <div class="rounded-xl bg-white border-2 border-green-200 p-5 md:col-span-2 lg:col-span-2 md:row-span-2 flex flex-col justify-center">
        <p class="text-3xl md:text-4xl font-bold text-gray-900">{formatCurrency(stats.totalRevenue)}</p>
        <p class="text-sm font-medium text-gray-600 mt-1">{t.jornada?.revenueConcreted ?? 'Ventas concretadas'}</p>
        <p class="text-xs text-gray-500 mt-0.5">{(stats.ordersByStatus?.delivered ?? 0) + (stats.ordersByStatus?.dispatched ?? 0)} {t.jornada?.statsOrders ?? 'órdenes'}</p>
      </div>
      <div class="rounded-xl bg-white border border-green-100 p-4">
        <p class="text-xl font-bold text-green-800">{formatCurrency(stats.revenueByStatus?.delivered ?? 0)}</p>
        <p class="text-xs text-gray-500 mt-0.5">{t.jornada?.revenueDelivered ?? 'Entregadas'}</p>
        <p class="text-xs text-gray-400">{stats.ordersByStatus?.delivered ?? 0} {t.jornada?.statsOrders ?? 'órdenes'}</p>
      </div>
      <div class="rounded-xl bg-white border border-blue-100 p-4">
        <p class="text-xl font-bold text-blue-800">{formatCurrency(stats.revenueByStatus?.dispatched ?? 0)}</p>
        <p class="text-xs text-gray-500 mt-0.5">{t.jornada?.revenueDispatched ?? 'Despachadas'}</p>
        <p class="text-xs text-gray-400">{stats.ordersByStatus?.dispatched ?? 0} {t.jornada?.statsOrders ?? 'órdenes'}</p>
      </div>
      <div class="rounded-xl bg-white border border-amber-200 p-4">
        <p class="text-xl font-bold text-amber-800">{formatCurrency(stats.revenueByStatus?.pending ?? 0)}</p>
        <p class="text-xs text-gray-500 mt-0.5">{t.jornada?.revenuePending ?? 'Ventas pendientes'}</p>
        <p class="text-xs text-gray-400">{stats.ordersByStatus?.pending ?? 0} {t.jornada?.statsOrders ?? 'órdenes'}</p>
      </div>
      <div class="rounded-xl bg-white border border-gray-200 p-4">
        <p class="text-xl font-bold text-gray-600">{formatCurrency(stats.revenueByStatus?.cancelled ?? 0)}</p>
        <p class="text-xs text-gray-500 mt-0.5">{t.jornada?.revenueCancelled ?? 'Canceladas'}</p>
        <p class="text-xs text-gray-400">{stats.ordersByStatus?.cancelled ?? 0} {t.jornada?.statsOrders ?? 'órdenes'}</p>
      </div>
      <div class="rounded-xl bg-white border border-gray-200 p-4">
        <p class="text-xl font-bold text-gray-900">{itemsOrdered}</p>
        <p class="text-xs text-gray-500 mt-0.5">{t.jornada?.itemsOrdered ?? 'Ítems ordenados'}</p>
      </div>
      <div class="rounded-xl bg-white border border-gray-200 p-4 md:col-span-2 lg:col-span-2">
        <p class="text-xl font-bold text-gray-900">{formatCurrency(averageTicket)}</p>
        <p class="text-xs text-gray-500 mt-0.5">{t.jornada?.averageTicket ?? 'Ticket promedio'}</p>
      </div>
      {#if topProduct}
        <div class="rounded-xl bg-blue-50 border border-blue-200 p-4 md:col-span-2 lg:col-span-2">
          <p class="text-xs font-medium text-blue-600 mb-1">{t.jornada?.topProduct ?? 'Producto más vendido'}</p>
          <p class="text-sm font-semibold text-gray-900 truncate" title={topProduct.productName}>{topProduct.productName}</p>
          <p class="text-xs text-gray-600 mt-0.5">{topProduct.quantitySold} vendidos — {formatCurrency(topProduct.totalRevenue)}</p>
        </div>
      {/if}
    </div>

    {#if stats.products.length > 0}
      {@const maxRevenue = Math.max(...stats.products.map((p) => p.totalRevenue), 1)}
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
        <div class="rounded-xl bg-white border border-gray-200 p-6 min-w-0">
          <h3 class="text-base font-semibold text-gray-900 mb-4">{t.jornada?.chartByRevenue ?? 'Top productos por ventas'}</h3>
          <div class="space-y-3">
            {#each stats.products as p, i}
              <div class="flex items-center gap-3 w-full">
                <span class="w-32 md:w-48 lg:w-64 text-sm font-medium text-gray-900 truncate shrink-0" title={p.productName}>{p.productName}</span>
                <span class="text-sm font-semibold text-gray-700 shrink-0 w-20 text-right">{formatCurrency(p.totalRevenue)}</span>
                <div class="flex-1 min-w-0 h-6 bg-gray-100 rounded overflow-hidden">
                  <div
                    class="h-full rounded transition-all"
                    style="width: {(p.totalRevenue / maxRevenue) * 100}%; background-color: {PIE_COLORS[i % PIE_COLORS.length]}"
                    role="img"
                    aria-label="{p.productName}: {formatCurrency(p.totalRevenue)}"
                  ></div>
                </div>
              </div>
            {/each}
          </div>
        </div>
        <div class="rounded-xl bg-white border border-gray-200 overflow-hidden min-w-0">
          <h3 class="text-base font-semibold text-gray-900 p-4 pb-2">{t.jornada?.productsTable ?? 'Productos'}</h3>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-200 bg-gray-50">
                  <th class="text-left font-medium text-gray-700 px-4 py-3">{t.jornada?.productName ?? 'Producto'}</th>
                  <th class="text-right font-medium text-gray-700 px-4 py-3">{t.jornada?.quantity ?? 'Cant.'}</th>
                  <th class="text-right font-medium text-gray-700 px-4 py-3">{t.jornada?.revenue ?? 'Ventas'}</th>
                </tr>
              </thead>
              <tbody>
                {#each stats.products as p}
                  <tr class="border-b border-gray-100 last:border-b-0 hover:bg-gray-50/50">
                    <td class="px-4 py-3 font-medium text-gray-900">{p.productName}</td>
                    <td class="px-4 py-3 text-right text-gray-600">{p.quantitySold}</td>
                    <td class="px-4 py-3 text-right font-medium text-gray-900">{formatCurrency(p.totalRevenue)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    {:else}
      <p class="text-gray-500 text-center py-12">{t.jornada?.noStats ?? 'No hay datos de ventas para esta jornada.'}</p>
    {/if}
  </div>
{:else}
  <p class="text-gray-500 text-center py-12">{t.jornada?.errorLoadingStats ?? 'Error al cargar estadísticas.'}</p>
{/if}
