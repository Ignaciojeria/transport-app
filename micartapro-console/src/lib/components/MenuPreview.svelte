<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, generateMenuUrl } from '../menuUtils'

  let menuUrl = $state<string | null>(null)
  let loading = $state(true)
  let error = $state<string | null>(null)
  let copySuccess = $state(false)

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')

  onMount(async () => {
    if (!userId) {
      error = 'Usuario no autenticado'
      loading = false
      return
    }

    try {
      const menuId = await getLatestMenuId(userId)
      
      if (!menuId) {
        error = 'No se encontró un menú para este usuario. Crea uno primero.'
        loading = false
        return
      }

      menuUrl = generateMenuUrl(userId, menuId)
      loading = false
    } catch (err: any) {
      console.error('Error cargando menú:', err)
      error = err.message || 'Error al cargar el menú'
      loading = false
    }
  })

  function copyToClipboard() {
    if (!menuUrl) return

    navigator.clipboard.writeText(menuUrl).then(() => {
      copySuccess = true
      setTimeout(() => {
        copySuccess = false
      }, 2000)
    }).catch((err) => {
      console.error('Error copiando al portapapeles:', err)
    })
  }

  function openInNewTab() {
    if (menuUrl) {
      window.open(menuUrl, '_blank')
    }
  }
</script>

<div class="flex flex-col h-full bg-gray-50">
  {#if loading}
    <div class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
        <p class="text-gray-600">Cargando tu carta...</p>
      </div>
    </div>
  {:else if error}
    <div class="flex-1 flex items-center justify-center p-6">
      <div class="text-center max-w-md">
        <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">Error</h3>
        <p class="text-gray-600 mb-4">{error}</p>
      </div>
    </div>
  {:else if menuUrl}
    <!-- Header con enlace -->
    <div class="bg-white border-b border-gray-200 px-4 py-3">
      <div class="max-w-4xl mx-auto">
        <h2 class="text-lg font-semibold text-gray-900 mb-3">Tu Carta Digital</h2>
        
        <!-- Enlace compartible -->
        <div class="flex items-center gap-2">
          <div class="flex-1 bg-gray-50 rounded-lg px-3 py-2 border border-gray-200">
            <p class="text-xs text-gray-500 mb-1">Enlace de tu carta:</p>
            <p class="text-sm text-gray-900 font-mono break-all">{menuUrl}</p>
          </div>
          
          <button
            onclick={copyToClipboard}
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center gap-2 whitespace-nowrap"
            title="Copiar enlace"
          >
            {#if copySuccess}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-sm">¡Copiado!</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <span class="text-sm">Copiar</span>
            {/if}
          </button>
          
          <button
            onclick={openInNewTab}
            class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors flex items-center gap-2 whitespace-nowrap"
            title="Abrir en nueva pestaña"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
            <span class="text-sm">Abrir</span>
          </button>
        </div>
      </div>
    </div>

    <!-- WebView con iframe -->
    <div class="flex-1 overflow-hidden bg-white">
      <iframe
        src={menuUrl}
        class="w-full h-full border-0"
        title="Vista previa de la carta"
        loading="lazy"
        allow="camera; microphone; geolocation; autoplay; clipboard-write"
        sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation"
      />
    </div>
  {/if}
</div>

<style>
  iframe {
    display: block;
  }
</style>

