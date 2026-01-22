<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getMenuVersions, getCurrentVersionId, updateCurrentVersionId, generateMenuUrlFromSlug, getMenuSlug, updateVersionName, updateVersionFavorite } from '../menuUtils'
  import { language } from '../useLanguage'

  interface MenuHistoryProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: MenuHistoryProps = $props()

  interface MenuVersion {
    id: string
    version_number: number
    created_at: string
    name: string | null
    is_favorite: boolean
    content?: any
  }

  let versions = $state<MenuVersion[]>([])
  let loading = $state(true)
  let error = $state<string | null>(null)
  let currentVersionId = $state<string | null>(null)
  let showPreview = $state(false)
  let previewVersionId = $state<string | null>(null)
  let previewUrl = $state<string | null>(null)
  let activatingVersionId = $state<string | null>(null)
  let menuSlug = $state<string | null>(null)
  let iframeKey = $state(0) // Key para forzar recarga del iframe
  let editingNameId = $state<string | null>(null) // ID de la versión cuyo nombre se está editando
  let editingName = $state<string>('') // Nombre temporal mientras se edita

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const session = $derived(authState.session)
  const currentLanguage = $derived($language)

  async function loadVersions() {
    if (!userId || !session?.access_token) {
      error = 'No hay sesión activa'
      loading = false
      return
    }

    try {
      loading = true
      error = null

      // Obtener menuId
      const menuId = await getLatestMenuId(userId)
      if (!menuId) {
        error = 'No se encontró un menú'
        loading = false
        return
      }

      // Obtener slug para generar URLs de preview
      const slug = await getMenuSlug(menuId, session.access_token)
      if (slug) {
        menuSlug = slug
      }

      // Obtener versiones y versión actual en paralelo
      const [versionsData, currentId] = await Promise.all([
        getMenuVersions(menuId, session.access_token),
        getCurrentVersionId(menuId, session.access_token)
      ])

      currentVersionId = currentId
      
      // Ordenar versiones: primero la activa, luego favoritas, luego el resto
      versions = versionsData.sort((a, b) => {
        // 1. La versión activa siempre primero
        if (a.id === currentId && b.id !== currentId) return -1
        if (b.id === currentId && a.id !== currentId) return 1
        
        // 2. Si ambas son activas o ninguna es activa, ordenar por favoritas
        if (a.is_favorite && !b.is_favorite) return -1
        if (!a.is_favorite && b.is_favorite) return 1
        
        // 3. Si ambas tienen el mismo estado de favorito, ordenar por número de versión (descendente)
        return b.version_number - a.version_number
      })
    } catch (err: any) {
      console.error('Error cargando versiones:', err)
      error = err.message || 'Error al cargar el historial'
    } finally {
      loading = false
    }
  }

  async function previewVersion(versionId: string) {
    if (!menuSlug) {
      alert('No se puede previsualizar: falta el slug del menú')
      return
    }

    previewVersionId = versionId
    previewUrl = generateMenuUrlFromSlug(menuSlug, currentLanguage, versionId)
    showPreview = true
    iframeKey++ // Forzar recarga del iframe
  }

  function closePreview() {
    showPreview = false
    previewUrl = null
    previewVersionId = null
  }

  async function activateVersion(versionId: string) {
    if (!userId || !session?.access_token) {
      alert('No hay sesión activa')
      return
    }

    if (!confirm('¿Deseas activar esta versión del menú?')) {
      return
    }

    try {
      activatingVersionId = versionId

      const menuId = await getLatestMenuId(userId)
      if (!menuId) {
        alert('No se encontró un menú')
        return
      }

      const success = await updateCurrentVersionId(menuId, versionId, session.access_token)
      
      if (success) {
        currentVersionId = versionId
        // Recargar versiones para actualizar el estado
        await loadVersions()
        // Cerrar preview y volver a la lista
        closePreview()
        alert('Versión activada exitosamente')
      } else {
        alert('Error al activar la versión')
      }
    } catch (err: any) {
      console.error('Error activando versión:', err)
      alert('Error al activar la versión: ' + err.message)
    } finally {
      activatingVersionId = null
    }
  }

  function formatDate(dateString: string): string {
    const date = new Date(dateString)
    return date.toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  function startEditingName(version: MenuVersion) {
    editingNameId = version.id
    editingName = version.name || ''
  }

  function cancelEditingName() {
    editingNameId = null
    editingName = ''
  }

  async function saveVersionName(versionId: string) {
    if (!session?.access_token) {
      alert('No hay sesión activa')
      return
    }

    try {
      console.log('💾 Guardando nombre de versión:', { versionId, name: editingName })
      const success = await updateVersionName(versionId, editingName, session.access_token)
      
      if (success) {
        console.log('✅ Nombre guardado exitosamente')
        // Actualizar el nombre en la lista local
        versions = versions.map(v => 
          v.id === versionId ? { ...v, name: editingName.trim() || null } : v
        )
        editingNameId = null
        editingName = ''
      } else {
        console.error('❌ Error al guardar el nombre')
        alert('Error al guardar el nombre. Revisa la consola para más detalles.')
      }
    } catch (err: any) {
      console.error('❌ Error guardando nombre:', err)
      alert('Error al guardar el nombre: ' + err.message)
    }
  }

  async function toggleFavorite(versionId: string, currentFavorite: boolean) {
    if (!session?.access_token) {
      alert('No hay sesión activa')
      return
    }

    try {
      const newFavoriteState = !currentFavorite
      console.log('💾 Actualizando favorito:', { versionId, newFavoriteState })
      const success = await updateVersionFavorite(versionId, newFavoriteState, session.access_token)
      
      if (success) {
        console.log('✅ Favorito actualizado exitosamente')
        // Actualizar el estado de favorito en la lista local
        versions = versions.map(v => 
          v.id === versionId ? { ...v, is_favorite: newFavoriteState } : v
        )
        // Reordenar después de cambiar favorito
        const currentId = currentVersionId
        versions = versions.sort((a, b) => {
          if (a.id === currentId && b.id !== currentId) return -1
          if (b.id === currentId && a.id !== currentId) return 1
          if (a.is_favorite && !b.is_favorite) return -1
          if (!a.is_favorite && b.is_favorite) return 1
          return b.version_number - a.version_number
        })
      } else {
        console.error('❌ Error al actualizar favorito')
        alert('Error al actualizar favorito. Revisa la consola para más detalles.')
      }
    } catch (err: any) {
      console.error('❌ Error actualizando favorito:', err)
      alert('Error al actualizar favorito: ' + err.message)
    }
  }

  onMount(() => {
    loadVersions()
  })
</script>

<div class="flex flex-col h-screen bg-white relative overflow-hidden">
  <!-- Vista de Lista (oculta cuando showPreview es true) -->
  <div 
    class="flex flex-col h-full transition-transform duration-300 ease-in-out {showPreview ? '-translate-x-full' : 'translate-x-0'}"
  >
    <!-- Header -->
    <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between">
      <button 
        onclick={onMenuClick}
        class="md:hidden p-2 hover:bg-gray-100 rounded-full transition-colors" 
        aria-label="Abrir menú"
      >
        <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <div class="hidden md:block w-9"></div>
      <h1 class="text-lg font-medium text-gray-900">Historial de Versiones</h1>
      <div class="w-9"></div>
    </header>

    <!-- Contenido -->
    <div class="flex-1 overflow-auto p-6">
    {#if loading}
      <div class="flex items-center justify-center h-full">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
          <p class="text-gray-600">Cargando versiones...</p>
        </div>
      </div>
    {:else if error}
      <div class="flex items-center justify-center h-full">
        <div class="text-center max-w-md">
          <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">Error</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <button
            onclick={loadVersions}
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            Reintentar
          </button>
        </div>
      </div>
    {:else if versions.length === 0}
      <div class="flex items-center justify-center h-full">
        <div class="text-center">
          <p class="text-gray-600">No hay versiones disponibles</p>
        </div>
      </div>
    {:else}
      <div class="max-w-4xl mx-auto space-y-4">
        {#each versions as version (version.id)}
          <div class="bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-start justify-between gap-3">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  {#if editingNameId === version.id}
                    <!-- Modo edición de nombre -->
                    <div class="flex-1 flex items-center gap-2">
                      <input
                        type="text"
                        bind:value={editingName}
                        placeholder="Nombre de la versión..."
                        class="flex-1 px-3 py-1.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                        onkeydown={(e) => {
                          if (e.key === 'Enter') {
                            saveVersionName(version.id)
                          } else if (e.key === 'Escape') {
                            cancelEditingName()
                          }
                        }}
                        autofocus
                      />
                      <button
                        onclick={() => saveVersionName(version.id)}
                        class="px-3 py-1.5 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm transition-colors"
                        title="Guardar"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                      </button>
                      <button
                        onclick={cancelEditingName}
                        class="px-3 py-1.5 bg-gray-300 hover:bg-gray-400 text-gray-700 rounded-lg text-sm transition-colors"
                        title="Cancelar"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  {:else}
                    <!-- Modo visualización -->
                    <h3 class="text-lg font-semibold text-gray-900">
                      {version.name || `Versión ${version.version_number}`}
                    </h3>
                    {#if version.id === currentVersionId}
                      <span class="px-2 py-1 bg-green-100 text-green-800 text-xs font-medium rounded-full">
                        Activa
                      </span>
                    {/if}
                    <!-- Botón de favorito -->
                    <button
                      onclick={() => toggleFavorite(version.id, version.is_favorite)}
                      class="p-1.5 hover:bg-gray-100 rounded-full transition-colors"
                      title={version.is_favorite ? 'Quitar de favoritos' : 'Agregar a favoritos'}
                    >
                      {#if version.is_favorite}
                        <svg class="w-5 h-5 text-red-500 fill-current" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                        </svg>
                      {:else}
                        <svg class="w-5 h-5 text-gray-400 hover:text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                        </svg>
                      {/if}
                    </button>
                    <!-- Botón para editar nombre -->
                    <button
                      onclick={() => startEditingName(version)}
                      class="p-1.5 hover:bg-gray-100 rounded-full transition-colors"
                      title="Editar nombre"
                    >
                      <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                  {/if}
                </div>
                <p class="text-sm text-gray-500 mb-3">
                  {formatDate(version.created_at)}
                </p>
                <div class="flex gap-2">
                  <button
                    onclick={() => previewVersion(version.id)}
                    class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium"
                  >
                    Previsualizar
                  </button>
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
    </div>
  </div>

  <!-- Vista de Preview (se muestra cuando showPreview es true) -->
  <div 
    class="absolute inset-0 flex flex-col h-full bg-white transition-transform duration-300 ease-in-out {showPreview ? 'translate-x-0' : 'translate-x-full'}"
  >
    <!-- Header del Preview -->
    <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between flex-shrink-0 z-10">
      <button
        onclick={closePreview}
        class="p-2 hover:bg-gray-100 rounded-full transition-colors"
        aria-label="Volver"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <div class="flex-1 min-w-0 px-2">
        <h2 class="text-base md:text-lg font-semibold text-gray-900 truncate text-center">
          Previsualización - Versión {versions.find(v => v.id === previewVersionId)?.version_number || ''}
        </h2>
      </div>
      <div class="w-9"></div> <!-- Spacer para centrar -->
    </header>

    <!-- Contenido del Preview -->
    <div class="flex-1 overflow-hidden iframe-container relative" style="padding-bottom: {previewUrl ? '90px' : '0'};">
      {#if previewUrl}
        {#key iframeKey}
          <iframe
            src={previewUrl}
            class="w-full h-full border-0"
            title="Vista previa de la versión"
            loading="lazy"
            allow="camera; microphone; geolocation; autoplay; clipboard-write"
            sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation"
          ></iframe>
        {/key}
      {:else}
        <div class="flex items-center justify-center h-full">
          <div class="text-center">
            <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
            <p class="text-gray-600">Cargando vista previa...</p>
          </div>
        </div>
      {/if}
    </div>

    <!-- Botón flotante: "Usar esta versión" -->
    {#if previewUrl && previewVersionId}
      {@const previewVersion = versions.find(v => v.id === previewVersionId)}
      {#if previewVersion && previewVersion.id !== currentVersionId}
        <!-- Botón desktop -->
        <div class="hidden md:block bg-white border-t border-gray-200 px-4 py-3 flex-shrink-0 safe-area-inset-bottom shadow-lg">
          <button
            onclick={() => {
              if (previewVersionId) {
                activateVersion(previewVersionId)
              }
            }}
            disabled={activatingVersionId === previewVersionId || !previewVersionId}
            class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-3"
          >
            {#if activatingVersionId === previewVersionId}
              <div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
              <span>Activando...</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>Usar esta versión</span>
            {/if}
          </button>
        </div>

        <!-- Botón móvil -->
        <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-4 py-3 z-[100] safe-area-inset-bottom shadow-2xl md:hidden">
          <button
            onclick={() => {
              if (previewVersionId) {
                activateVersion(previewVersionId)
              }
            }}
            disabled={activatingVersionId === previewVersionId || !previewVersionId}
            class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-3"
          >
            {#if activatingVersionId === previewVersionId}
              <div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
              <span>Activando...</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>Usar esta versión</span>
            {/if}
          </button>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  /* Fix para scroll en iframes en mobile */
  .iframe-container {
    touch-action: pan-y pan-x;
    -webkit-overflow-scrolling: touch;
    position: relative;
  }
  
  .iframe-container iframe {
    touch-action: auto;
    pointer-events: auto;
  }

  /* Safe area para dispositivos con notch */
  .safe-area-inset-bottom {
    padding-bottom: env(safe-area-inset-bottom, 0.75rem);
  }
</style>
