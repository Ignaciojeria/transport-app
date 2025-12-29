<script>
  import { onMount } from 'svelte'
  import { initAuth, authState } from './lib/auth.svelte'
  import MenuChat from './lib/components/MenuChat.svelte'
  import MenuPreview from './lib/components/MenuPreview.svelte'

  // Usar valores derivados reactivos en el componente
  let user = $derived(authState.user)
  let session = $derived(authState.session)
  let loading = $derived(authState.loading)
  let currentView = $state('preview')

  onMount(() => {
    initAuth()
    
    // Mostrar fragment en consola para debugging
    const fragment = window.location.hash
    if (fragment.startsWith('#auth=')) {
      try {
        const encodedData = fragment.substring(6)
        const authData = JSON.parse(atob(encodedData))
        console.log('📦 Contenido completo del fragment decodificado:')
        console.log(JSON.stringify(authData, null, 2))
      } catch (e) {
        console.error('Error decodificando fragment:', e)
      }
    }
  })
</script>

<div class="min-h-screen bg-gray-50">
  {#if loading}
    <div class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
        <p class="text-gray-600">Cargando...</p>
      </div>
    </div>
  {:else if user}
    <!-- Navigation Tabs -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4">
        <nav class="flex space-x-8">
          <button
            onclick={() => currentView = 'preview'}
            class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentView === 'preview' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
          >
            Vista Previa de Carta
          </button>
          <button
            onclick={() => currentView = 'chat'}
            class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {currentView === 'chat' ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
          >
            Crear/Editar Menú
          </button>
        </nav>
      </div>
    </div>

    <!-- Content -->
    <div class="h-[calc(100vh-57px)]">
      {#if currentView === 'preview'}
        <MenuPreview />
      {:else}
        <MenuChat />
      {/if}
    </div>
  {:else}
    <div class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">
          No autenticado
        </h1>
        <p class="text-gray-600 mb-4">
          Por favor, inicia sesión en{' '}
          <a 
            href="http://localhost:3003" 
            class="text-blue-600 hover:text-blue-700 underline"
          >
            micartapro-auth-ui
          </a>
        </p>
      </div>
    </div>
  {/if}
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: system-ui, -apple-system, sans-serif;
  }
</style>
