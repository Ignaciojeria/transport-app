<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import {
    getUserMenus,
    createNewMenu,
    setActiveMenuId,
    getLatestMenuId,
    type UserMenu
  } from '../menuUtils'
  import { t as tStore } from '../useLanguage'

  interface NegociosViewProps {
    onMenuClick?: () => void
    onBusinessSelected?: () => void
  }

  let { onMenuClick, onBusinessSelected }: NegociosViewProps = $props()

  let menus = $state<UserMenu[]>([])
  let activeMenuId = $state<string | null>(null)
  let loading = $state(true)
  let error = $state<string | null>(null)
  let creating = $state(false)
  let showCreateModal = $state(false)
  let newMenuSlug = $state('')

  const session = $derived(authState.session)
  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const t = $derived($tStore)

  function formatDate(iso: string): string {
    try {
      const d = new Date(iso)
      return d.toLocaleDateString(undefined, {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    } catch {
      return iso
    }
  }

  async function load() {
    if (!userId || !session?.access_token) {
      loading = false
      return
    }
    try {
      loading = true
      error = null
      const [menusData, latestId] = await Promise.all([
        getUserMenus(userId, session.access_token),
        getLatestMenuId(userId, session.access_token)
      ])
      menus = menusData
      activeMenuId = latestId
    } catch (e) {
      console.error('Error cargando negocios:', e)
      error = t.negocios?.errorLoading ?? 'Error al cargar los negocios.'
    } finally {
      loading = false
    }
  }

  function openCreateModal() {
    newMenuSlug = ''
    showCreateModal = true
  }

  function closeCreateModal() {
    showCreateModal = false
    newMenuSlug = ''
  }

  async function handleCreateNew() {
    if (!userId || !session?.access_token) return
    creating = true
    error = null
    try {
      const slug = newMenuSlug.trim() || undefined
      const newId = await createNewMenu(userId, session.access_token, slug)
      if (newId) {
        closeCreateModal()
        activeMenuId = newId
        await load()
        onMenuClick?.()
        onBusinessSelected?.()
      } else {
        error = t.negocios?.errorCreating ?? 'Error al crear el nuevo menú.'
      }
    } catch (e) {
      console.error('Error creando menú:', e)
      error = t.negocios?.errorCreating ?? 'Error al crear el nuevo menú.'
    } finally {
      creating = false
    }
  }

  async function handleSelect(menu: UserMenu) {
    if (!userId || !session?.access_token || activeMenuId === menu.menuId) return
    const ok = await setActiveMenuId(userId, menu.menuId, session.access_token)
    if (ok) {
      activeMenuId = menu.menuId
      onMenuClick?.()
      onBusinessSelected?.()
    }
  }

  onMount(() => {
    load()

    // Demo Remotion: seleccionar negocio por slug al recibir postMessage
    const msgHandler = (e: MessageEvent) => {
      const allowed = ['http://localhost:', 'http://127.0.0.1:']
      if (!allowed.some(o => e.origin?.startsWith(o))) return
      const slug = e.data?.clickSelectBusiness
      if (slug && menus.length > 0) {
        const menu = menus.find((m) => m.slug === slug)
        if (menu && menu.menuId !== activeMenuId) {
          handleSelect(menu)
        }
      }
    }
    window.addEventListener('message', msgHandler)
    return () => window.removeEventListener('message', msgHandler)
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
    <h1 class="text-lg font-semibold text-gray-900">{t.negocios?.title ?? 'Negocios'}</h1>
  </header>

  <div class="flex-1 min-h-0 overflow-y-auto p-4 md:p-6">
    {#if loading}
      <div class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-10 w-10 border-2 border-blue-600 border-t-transparent"></div>
      </div>
    {:else if error}
      <div class="rounded-xl bg-amber-50 border border-amber-200 p-4 text-amber-800">
        <p class="font-medium">{error}</p>
      </div>
    {:else}
      <p class="text-gray-600 text-sm mb-6">
        {t.negocios?.subtitle ?? 'Selecciona tu negocio activo o crea uno nuevo para trabajar.'}
      </p>

      <!-- Botón crear nuevo -->
      <button
        type="button"
        disabled={creating}
        onclick={openCreateModal}
        class="w-full flex items-center justify-center gap-2 py-4 px-6 rounded-xl font-semibold text-white bg-blue-600 hover:bg-blue-700 focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-60 transition-colors mb-8"
      >
        {#if creating}
          <span class="animate-spin inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full"></span>
          <span>{t.negocios?.creating ?? 'Creando...'}</span>
        {:else}
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>{t.negocios?.createNew ?? 'Crear nuevo menú'}</span>
        {/if}
      </button>

      <!-- Lista de negocios -->
      <h2 class="text-sm font-semibold text-gray-700 mb-3">{t.negocios?.yourBusinesses ?? 'Tus negocios'}</h2>
      {#if menus.length === 0}
        <p class="text-gray-500 text-sm py-4">{t.negocios?.noBusinesses ?? 'Aún no tienes negocios. Crea uno arriba.'}</p>
      {:else}
        <div class="space-y-2">
          {#each menus as menu}
            <button
              type="button"
              onclick={() => handleSelect(menu)}
              class="w-full flex items-center justify-between p-4 rounded-xl border transition-colors text-left {activeMenuId === menu.menuId
                ? 'border-blue-600 bg-blue-50 ring-2 ring-blue-600'
                : 'border-gray-200 bg-white hover:border-gray-300 hover:bg-gray-50'}"
            >
              <div class="flex items-center gap-3">
                <span
                  class="inline-flex w-3 h-3 rounded-full {activeMenuId === menu.menuId ? 'bg-blue-600' : 'bg-gray-300'}"
                  aria-hidden="true"
                ></span>
                <div>
                  <p class="font-medium text-gray-900">
                    {menu.slug || `${t.negocios?.business ?? 'Negocio'} ${menu.menuId.slice(0, 8)}...`}
                  </p>
                  <p class="text-xs text-gray-500">{formatDate(menu.createdAt)}</p>
                </div>
              </div>
              {#if activeMenuId === menu.menuId}
                <span class="text-xs font-medium text-blue-600">{t.negocios?.active ?? 'Activo'}</span>
              {:else}
                <span class="text-xs text-gray-400">{t.negocios?.select ?? 'Seleccionar'}</span>
              {/if}
            </button>
          {/each}
        </div>
      {/if}
    {/if}
  </div>

  <!-- Modal crear nuevo menú -->
  {#if showCreateModal}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
      role="dialog"
      aria-modal="true"
      aria-labelledby="new-menu-title"
    >
      <!-- Backdrop: clic para cerrar -->
      <button
        type="button"
        class="absolute inset-0 -z-10 cursor-default"
        aria-label={t.negocios?.cancel ?? 'Cerrar'}
        onclick={closeCreateModal}
      ></button>
      <div
        class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6 relative z-10"
        role="document"
      >
        <h2 id="new-menu-title" class="text-lg font-semibold text-gray-900 mb-4">
          {t.negocios?.newMenuTitle ?? 'Nuevo negocio'}
        </h2>
        <label for="new-menu-slug" class="block text-sm font-medium text-gray-700 mb-2">
          {t.negocios?.slugLabel ?? 'Identificador de tu negocio (slug)'}
        </label>
        <input
          id="new-menu-slug"
          type="text"
          bind:value={newMenuSlug}
          placeholder={t.negocios?.slugPlaceholder ?? 'ej: mi-restaurante'}
          class="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none mb-2"
          onkeydown={(e) => e.key === 'Enter' && handleCreateNew()}
        />
        <p class="text-xs text-gray-500 mb-6">
          {t.negocios?.slugHint ?? 'Solo letras minúsculas, números y guiones. Será la URL de tu carta.'}
        </p>
        <div class="flex gap-3 justify-end">
          <button
            type="button"
            onclick={closeCreateModal}
            class="px-4 py-2 rounded-xl font-medium text-gray-700 hover:bg-gray-100"
          >
            {t.negocios?.cancel ?? 'Cancelar'}
          </button>
          <button
            type="button"
            disabled={creating}
            onclick={handleCreateNew}
            class="px-4 py-2 rounded-xl font-semibold text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-60"
          >
            {#if creating}
              <span class="animate-spin inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full align-middle"></span>
              <span class="ml-1">{t.negocios?.creating ?? 'Creando...'}</span>
            {:else}
              {t.negocios?.create ?? 'Crear'}
            {/if}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
