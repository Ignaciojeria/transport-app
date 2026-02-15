<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId } from '../menuUtils'
  import { getJourneys, getJourneyStats, type JourneyListItem, type JourneyStats } from '../journeyApi'
  import JourneyStatsView from './JourneyStatsView.svelte'
  import { t as tStore } from '../useLanguage'

  interface ReportesProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: ReportesProps = $props()

  let menuId = $state<string | null>(null)
  let loading = $state(true)
  let error = $state<string | null>(null)
  let journeys = $state<JourneyListItem[]>([])
  let view = $state<'list' | 'stats'>('list')
  let statsJourney = $state<JourneyListItem | null>(null)
  let stats = $state<JourneyStats | null>(null)
  let statsLoading = $state(false)
  const REPORT_PAGE_SIZE = 10
  let reportPage = $state(1)

  const session = $derived(authState.session)
  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const t = $derived($tStore)

  /** Lista de jornadas cerradas. */
  const closedJourneysList = $derived(journeys.filter((j) => j.status === 'CLOSED'))
  const reportTotalPages = $derived(Math.max(1, Math.ceil(closedJourneysList.length / REPORT_PAGE_SIZE)))
  const paginatedReports = $derived(
    closedJourneysList.slice((reportPage - 1) * REPORT_PAGE_SIZE, reportPage * REPORT_PAGE_SIZE)
  )

  /** Formato DD-MM HH:mm para jornadas. */
  function formatJourneyDate(isoDate: string): string {
    const d = new Date(isoDate)
    const day = String(d.getDate()).padStart(2, '0')
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const h = String(d.getHours()).padStart(2, '0')
    const min = String(d.getMinutes()).padStart(2, '0')
    return `${day}/${m} ${h}:${min}`
  }

  function goToReportPage(p: number) {
    reportPage = Math.max(1, Math.min(p, reportTotalPages))
  }

  async function load() {
    if (!userId || !session?.access_token) {
      error = t.jornada?.noSession ?? 'No Hay Sesión Activa'
      loading = false
      return
    }
    try {
      loading = true
      error = null
      const mid = await getLatestMenuId(userId, session.access_token)
      if (!mid) {
        error = t.jornada?.noMenu ?? 'No Se Encontró Un Menú'
        loading = false
        return
      }
      menuId = mid
      const journeysList = await getJourneys(mid, session.access_token)
      journeys = journeysList
      reportPage = 1
    } catch (e) {
      console.error('Error cargando reportes:', e)
      error = t.jornada?.errorLoading ?? 'Error Al Cargar Los Datos.'
    } finally {
      loading = false
    }
  }

  async function openStatsView(j: JourneyListItem) {
    if (!session?.access_token || !menuId) return
    statsJourney = j
    view = 'stats'
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

  function goBackToList() {
    view = 'list'
    statsJourney = null
    stats = null
  }

  function downloadStatsCSV() {
    if (!stats?.products?.length || !statsJourney) return
    const headers = [t.jornada?.productName ?? 'Producto', t.jornada?.quantity ?? 'Cant.', t.jornada?.revenue ?? 'Ventas', t.jornada?.cost ?? 'Costo', t.jornada?.profit ?? 'Ganancias']
    const rows = stats.products.map((p) => {
      const cost = p.totalCost ?? 0
      const margin = (p.totalRevenue ?? 0) - cost
      return [p.productName, String(p.quantitySold), String(p.totalRevenue), String(cost), String(margin)]
    })
    const csv = [headers.join(','), ...rows.map((r) => r.map((c) => `"${String(c).replace(/"/g, '""')}"`).join(','))].join('\n')
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `jornada-${statsJourney.id}.csv`
    a.click()
    URL.revokeObjectURL(url)
  }

  onMount(() => {
    load()
  })
</script>

<div class="h-full flex flex-col bg-gray-50">
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
    <h1 class="text-lg font-semibold text-gray-900">{t.jornada?.reports ?? 'Reportes de Jornadas'}</h1>
  </header>

  <div class="flex-1 min-h-0 overflow-y-auto p-4 md:p-6">
    {#if view === 'stats'}
      <JourneyStatsView
        stats={stats}
        statsLoading={statsLoading}
        journeyInfo={statsJourney ? { id: statsJourney.id, openedAt: statsJourney.openedAt, closedAt: statsJourney.closedAt, reportXlsxUrl: statsJourney.reportXlsxUrl } : null}
        onBack={goBackToList}
        backLabel={t.jornada?.back ?? 'Volver a reportes'}
        onDownloadCSV={downloadStatsCSV}
      />
    {:else if loading}
      <div class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-10 w-10 border-2 border-blue-600 border-t-transparent"></div>
      </div>
    {:else if error}
      <div class="rounded-xl bg-amber-50 border border-amber-200 p-4 text-amber-800">
        <p class="font-medium">{error}</p>
      </div>
    {:else if closedJourneysList.length === 0}
      <p class="text-gray-500 text-center py-12">{t.jornada?.reports ?? 'Reportes de Jornadas'}: no hay jornadas cerradas.</p>
    {:else}
      <!-- Lista de reportes -->
      <section>
        <h2 class="text-sm font-semibold text-gray-700 mb-3">{t.jornada?.reports ?? 'Reportes de Jornadas'}</h2>
        <div class="space-y-2">
          {#each paginatedReports as j (j.id)}
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
                  onclick={() => openStatsView(j)}
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
        {#if reportTotalPages > 1}
          <div class="flex items-center justify-between mt-4 px-2">
            <button
              type="button"
              disabled={reportPage <= 1}
              onclick={() => goToReportPage(reportPage - 1)}
              class="px-3 py-2 rounded-lg text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              ← {t.jornada?.prev ?? 'Anterior'}
            </button>
            <span class="text-sm text-gray-600">
              {t.jornada?.page ?? 'Página'} {reportPage} / {reportTotalPages}
            </span>
            <button
              type="button"
              disabled={reportPage >= reportTotalPages}
              onclick={() => goToReportPage(reportPage + 1)}
              class="px-3 py-2 rounded-lg text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {t.jornada?.next ?? 'Siguiente'} →
            </button>
          </div>
        {/if}
      </section>
    {/if}
  </div>
</div>
