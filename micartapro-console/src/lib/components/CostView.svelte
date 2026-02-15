<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import {
    getLatestMenuId,
    getCurrentVersionId,
    getCurrentVersionContent,
    updateMenuVersionContent,
    type MenuContent
  } from '../menuUtils'
  import { t as tStore } from '../useLanguage'

  interface CostViewProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: CostViewProps = $props()

  let menuId = $state<string | null>(null)
  let versionId = $state<string | null>(null)
  let content = $state<MenuContent | null>(null)
  let loading = $state(true)
  let saving = $state(false)
  let error = $state<string | null>(null)
  let successMessage = $state<string | null>(null)

  const session = $derived(authState.session)
  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const t = $derived($tStore)

  /** Unidades permitidas por el agente (menu_create_request.go, pricing.go) */
  const ALLOWED_UNITS = [
    { value: 'EACH', label: 'Unidad' },
    { value: 'GRAM', label: 'Gramo' },
    { value: 'KILOGRAM', label: 'Kilogramo' },
    { value: 'MILLILITER', label: 'Mililitro' },
    { value: 'LITER', label: 'Litro' },
    { value: 'METER', label: 'Metro' },
    { value: 'SQUARE_METER', label: 'm²' },
  ] as const

  function getItemName(item: { title?: { base?: string } }): string {
    return item?.title?.base ?? 'Sin nombre'
  }

  function getSideName(side: { name?: { base?: string } }): string {
    return side?.name?.base ?? 'Sin nombre'
  }

  function getCost(pricing: { costPerUnit?: number } | undefined): number {
    return pricing?.costPerUnit ?? 0
  }

  function getPrice(pricing: { pricePerUnit?: number } | undefined): number {
    return pricing?.pricePerUnit ?? 0
  }

  function getUnit(pricing: { unit?: string } | undefined): string {
    return pricing?.unit ?? 'EACH'
  }

  function getBaseUnit(pricing: { baseUnit?: number } | undefined): number {
    return pricing?.baseUnit ?? 1
  }

  function setItemUnit(catIdx: number, itemIdx: number, value: string) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    if (!item) return
    if (!item.pricing) item.pricing = { pricePerUnit: 0, mode: 'UNIT', unit: 'EACH', baseUnit: 1 }
    item.pricing.unit = value
    if (value === 'EACH') item.pricing.baseUnit = 1
    content = { ...content }
  }

  function setItemBaseUnit(catIdx: number, itemIdx: number, value: number) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    if (!item) return
    if (!item.pricing) item.pricing = { pricePerUnit: 0, mode: 'UNIT', unit: 'EACH', baseUnit: 1 }
    item.pricing.baseUnit = value
    content = { ...content }
  }

  function setItemPrice(catIdx: number, itemIdx: number, value: number) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    if (!item) return
    if (!item.pricing) item.pricing = { pricePerUnit: 0, mode: 'UNIT', unit: 'EACH', baseUnit: 1 }
    item.pricing.pricePerUnit = value
    content = { ...content }
  }

  function setItemCost(catIdx: number, itemIdx: number, value: number) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    if (!item) return
    if (!item.pricing) item.pricing = { pricePerUnit: 0, mode: 'UNIT', unit: 'EACH', baseUnit: 1 }
    item.pricing.costPerUnit = value
    content = { ...content }
  }

  function setSideUnit(catIdx: number, itemIdx: number, sideIdx: number, value: string) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    const side = item?.sides?.[sideIdx]
    if (!side) return
    if (!side.pricing) side.pricing = { pricePerUnit: 0, unit: 'EACH', baseUnit: 1 }
    side.pricing.unit = value
    if (value === 'EACH') side.pricing.baseUnit = 1
    content = { ...content }
  }

  function setSideBaseUnit(catIdx: number, itemIdx: number, sideIdx: number, value: number) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    const side = item?.sides?.[sideIdx]
    if (!side) return
    if (!side.pricing) side.pricing = { pricePerUnit: 0 }
    side.pricing.baseUnit = value
    content = { ...content }
  }

  function setSideCost(catIdx: number, itemIdx: number, sideIdx: number, value: number) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    const side = item?.sides?.[sideIdx]
    if (!side) return
    if (!side.pricing) side.pricing = { pricePerUnit: 0 }
    side.pricing.costPerUnit = value
    content = { ...content }
  }

  function setSidePrice(catIdx: number, itemIdx: number, sideIdx: number, value: number) {
    if (!content?.menu) return
    const cat = content.menu[catIdx]
    const item = cat?.items?.[itemIdx]
    const side = item?.sides?.[sideIdx]
    if (!side) return
    if (!side.pricing) side.pricing = { pricePerUnit: 0 }
    side.pricing.pricePerUnit = value
    content = { ...content }
  }

  async function load() {
    if (!userId || !session?.access_token) {
      error = 'No hay sesión activa'
      loading = false
      return
    }
    try {
      loading = true
      error = null
      successMessage = null
      const mid = await getLatestMenuId(userId, session.access_token)
      if (!mid) {
        error = 'No se encontró un menú'
        loading = false
        return
      }
      menuId = mid
      const vid = await getCurrentVersionId(mid, session.access_token)
      versionId = vid
      const c = await getCurrentVersionContent(mid, session.access_token)
      content = c
    } catch (e) {
      console.error('Error cargando menú:', e)
      error = 'Error al cargar el menú'
    } finally {
      loading = false
    }
  }

  async function save() {
    if (!versionId || !content || !session?.access_token) return
    try {
      saving = true
      error = null
      successMessage = null
      const ok = await updateMenuVersionContent(versionId, content, session.access_token)
      if (ok) {
        successMessage = t.cost?.saved ?? 'Precio y costo guardados correctamente'
      } else {
        error = 'Error al guardar los costos'
      }
    } catch (e) {
      console.error('Error guardando:', e)
      error = 'Error al guardar los costos'
    } finally {
      saving = false
    }
  }

  onMount(() => load())
</script>

<div class="h-full overflow-y-auto bg-gray-50 p-4 md:p-6">
  <div class="max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold text-gray-900 mb-4">
      {t.cost?.title ?? 'Costos del menú'}
    </h1>
    <p class="text-gray-600 mb-6">
      {t.cost?.subtitle ?? 'Configura el costo de cada ítem y acompañamiento. Se usa para márgenes y reportes.'}
    </p>

    {#if loading}
      <div class="flex items-center justify-center py-12">
        <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600"></div>
      </div>
    {:else if error}
      <div class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700 mb-4">
        {error}
      </div>
      <button
        onclick={load}
        class="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded-lg text-sm font-medium"
      >
        Reintentar
      </button>
    {:else if !content?.menu?.length}
      <div class="bg-amber-50 border border-amber-200 rounded-lg p-4 text-amber-800">
        No hay categorías en el menú.
      </div>
    {:else}
      {#if successMessage}
        <div class="bg-green-50 border border-green-200 rounded-lg p-4 text-green-700 mb-4">
          {successMessage}
        </div>
      {/if}

      <div class="space-y-6 mb-8">
        {#each content.menu as category, catIdx}
          {@const catTitle = category.title?.base ?? 'Sin categoría'}
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <div class="px-4 py-3 bg-gray-100 border-b border-gray-200 font-semibold text-gray-800">
              {catTitle}
            </div>
            <div class="divide-y divide-gray-100">
              {#each category.items ?? [] as item, itemIdx}
                <div class="p-4">
                  <div class="font-medium text-gray-900 mb-2">{getItemName(item)}</div>
                  <div class="flex items-center gap-4 flex-wrap mb-2">
                    <label class="flex items-center gap-2">
                      <span class="text-sm text-gray-500">{t.cost?.unit ?? 'Unidad de venta'}:</span>
                      <select
                        value={getUnit(item.pricing)}
                        onchange={(e) => setItemUnit(catIdx, itemIdx, (e.target as HTMLSelectElement).value)}
                        class="px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white"
                      >
                        {#each ALLOWED_UNITS as u}
                          <option value={u.value}>{u.label}</option>
                        {/each}
                      </select>
                    </label>
                    {#if getUnit(item.pricing) !== 'EACH'}
                      <label class="flex items-center gap-2">
                        <span class="text-sm text-gray-500">{t.cost?.baseUnit ?? 'Base'}:</span>
                        <input
                          type="number"
                          min="0.001"
                          step="0.01"
                          value={getBaseUnit(item.pricing)}
                          oninput={(e) => setItemBaseUnit(catIdx, itemIdx, parseFloat((e.target as HTMLInputElement).value) || 1)}
                          class="w-20 px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        />
                      </label>
                    {/if}
                  </div>
                  <div class="flex items-center gap-4 flex-wrap">
                    <label class="flex items-center gap-2">
                      <span class="text-sm text-gray-500">{t.cost?.priceSale ?? 'Precio venta'}:</span>
                      <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={getPrice(item.pricing)}
                        oninput={(e) => setItemPrice(catIdx, itemIdx, parseFloat((e.target as HTMLInputElement).value) || 0)}
                        class="w-24 px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </label>
                    <label class="flex items-center gap-2">
                      <span class="text-sm text-gray-500">{t.cost?.cost ?? 'Costo'}:</span>
                      <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={getCost(item.pricing)}
                        oninput={(e) => setItemCost(catIdx, itemIdx, parseFloat((e.target as HTMLInputElement).value) || 0)}
                        class="w-24 px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </label>
                  </div>
                  {#if item.sides?.length}
                    <div class="ml-4 mt-2 space-y-2">
                      {#each item.sides as side, sideIdx}
                        <div class="text-sm">
                          <div class="text-gray-600 mb-1">↳ {getSideName(side)}</div>
                          <div class="flex items-center gap-4 flex-wrap ml-4 mb-2">
                            <label class="flex items-center gap-2">
                              <span class="text-gray-500">{t.cost?.unit ?? 'Unidad de venta'}:</span>
                              <select
                                value={getUnit(side.pricing)}
                                onchange={(e) => setSideUnit(catIdx, itemIdx, sideIdx, (e.target as HTMLSelectElement).value)}
                                class="px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white"
                              >
                                {#each ALLOWED_UNITS as u}
                                  <option value={u.value}>{u.label}</option>
                                {/each}
                              </select>
                            </label>
                            {#if getUnit(side.pricing) !== 'EACH'}
                              <label class="flex items-center gap-2">
                                <span class="text-gray-500">{t.cost?.baseUnit ?? 'Base'}:</span>
                                <input
                                  type="number"
                                  min="0.001"
                                  step="0.01"
                                  value={getBaseUnit(side.pricing)}
                                  oninput={(e) => setSideBaseUnit(catIdx, itemIdx, sideIdx, parseFloat((e.target as HTMLInputElement).value) || 1)}
                                  class="w-16 px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                />
                              </label>
                            {/if}
                          </div>
                          <div class="flex items-center gap-4 flex-wrap ml-4">
                            <label class="flex items-center gap-2">
                              <span class="text-gray-500">{t.cost?.priceSale ?? 'Precio venta'}:</span>
                              <input
                                type="number"
                                min="0"
                                step="0.01"
                                value={getPrice(side.pricing)}
                                oninput={(e) => setSidePrice(catIdx, itemIdx, sideIdx, parseFloat((e.target as HTMLInputElement).value) || 0)}
                                class="w-20 px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                              />
                            </label>
                            <label class="flex items-center gap-2">
                              <span class="text-gray-500">{t.cost?.cost ?? 'Costo'}:</span>
                              <input
                                type="number"
                                min="0"
                                step="0.01"
                                value={getCost(side.pricing)}
                                oninput={(e) => setSideCost(catIdx, itemIdx, sideIdx, parseFloat((e.target as HTMLInputElement).value) || 0)}
                                class="w-20 px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                              />
                            </label>
                          </div>
                        </div>
                      {/each}
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>

      <div class="flex gap-3">
        <button
          onclick={save}
          disabled={saving}
          class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
        >
          {saving ? 'Guardando...' : (t.cost?.save ?? 'Guardar costos')}
        </button>
        <button
          onclick={load}
          disabled={saving}
          class="px-4 py-2.5 bg-gray-200 hover:bg-gray-300 disabled:opacity-50 text-gray-700 rounded-lg font-medium"
        >
          {t.cost?.reload ?? 'Recargar'}
        </button>
      </div>
    {/if}
  </div>
</div>
