<script lang="ts">
  import { supabase } from '../supabase'
  import { authState } from '../auth.svelte'

  interface PhotoUploadProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: PhotoUploadProps = $props()

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')

  let photos = $state<Array<{ id: number; image_url: string }>>([])
  let loadingPhotos = $state(false)
  let copySuccess = $state(false)

  // Cargar fotos cuando el componente se monta o cuando cambia el userId
  $effect(() => {
    if (userId) {
      loadPhotos()
    }
  })

  async function loadPhotos() {
    if (!userId) return

    loadingPhotos = true
    try {
      // Consultar solo imágenes del usuario actual
      const { data, error } = await supabase
        .from('catalog_images')
        .select('id, image_url')
        .eq('user_id', userId)
        .order('id', { ascending: false })
        .limit(100)

      if (error) {
        console.error('Error cargando fotos:', error)
        return
      }

      photos = (data || []).map((row) => ({
        id: row.id,
        image_url: row.image_url
      }))
    } catch (error) {
      console.error('Error inesperado cargando fotos:', error)
    } finally {
      loadingPhotos = false
    }
  }

  async function deletePhoto(photoId: number) {
    if (!userId) return

    if (!confirm('¿Estás seguro de que quieres eliminar esta foto?')) {
      return
    }

    try {
      const { error } = await supabase
        .from('catalog_images')
        .delete()
        .eq('id', photoId)

      if (error) {
        console.error('Error eliminando foto:', error)
        alert('Error al eliminar la foto')
        return
      }

      // Recargar la lista de fotos
      await loadPhotos()
    } catch (error) {
      console.error('Error inesperado eliminando foto:', error)
      alert('Error al eliminar la foto')
    }
  }

  function copyPhotoUrl(url: string) {
    navigator.clipboard.writeText(url)
    copySuccess = true
    setTimeout(() => {
      copySuccess = false
    }, 2000)
  }
</script>

<div class="flex flex-col h-screen h-[100dvh] bg-gray-50 overflow-hidden">
  <!-- Header con botón hamburguesa -->
  <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between md:hidden sticky top-0 z-10 flex-shrink-0">
    <button 
      onclick={onMenuClick}
      class="p-2 hover:bg-gray-100 rounded-full transition-colors" 
      aria-label="Abrir menú"
    >
      <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
    </button>
    <h1 class="text-lg font-medium text-gray-900">Galería</h1>
    <div class="w-10"></div> <!-- Spacer para centrar -->
  </header>

  <div class="flex-1 overflow-y-auto p-6 max-w-6xl mx-auto min-h-0">
    <!-- Galería de fotos -->
    <div class="bg-white rounded-lg shadow-lg p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-bold text-gray-900">Tus Fotos</h2>
        <button
          onclick={loadPhotos}
          disabled={loadingPhotos}
          class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg font-semibold hover:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          {#if loadingPhotos}
            <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>Cargando...</span>
          {:else}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Actualizar</span>
          {/if}
        </button>
      </div>

      {#if loadingPhotos && photos.length === 0}
        <div class="text-center py-12">
          <svg class="animate-spin h-8 w-8 text-gray-400 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="text-gray-500">Cargando fotos...</p>
        </div>
      {:else if photos.length === 0}
        <div class="text-center py-12">
          <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p class="text-gray-500">No hay fotos todavía. Sube tu primera foto arriba.</p>
        </div>
      {:else}
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
          {#each photos as photo}
            <div class="relative group">
              <div class="aspect-square rounded-lg overflow-hidden bg-gray-100 border border-gray-200">
                <img 
                  src={photo.image_url} 
                  alt={`Imagen ${photo.id}`}
                  class="w-full h-full object-cover"
                  loading="lazy"
                />
              </div>
              <!-- Overlay con acciones -->
              <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 transition-all duration-200 rounded-lg flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100">
                <button
                  onclick={() => copyPhotoUrl(photo.image_url)}
                  class="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
                  title="Copiar URL"
                >
                  <svg class="w-5 h-5 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                </button>
                <button
                  onclick={() => deletePhoto(photo.id)}
                  class="p-2 bg-red-600 rounded-full hover:bg-red-700 transition-colors"
                  title="Eliminar"
                >
                  <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</div>
