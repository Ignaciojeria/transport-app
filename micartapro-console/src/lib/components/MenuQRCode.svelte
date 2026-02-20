<script lang="ts">
  import { onMount } from 'svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, getMenuSlugFromApi, generateMenuUrlFromSlug, generateMenuUrlFromMenuId } from '../menuUtils'
  import { language } from '../useLanguage'

  interface MenuQRCodeProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: MenuQRCodeProps = $props()

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const session = $derived(authState.session)
  const currentLanguage = $derived($language)

  let loading = $state(true)
  let error = $state<string | null>(null)
  let menuUrl = $state<string | null>(null)
  let qrCodeUrl = $state<string | null>(null)
  let slug = $state<string | null>(null)
  let menuWindowRef = $state<Window | null>(null)

  async function loadQRCode() {
    if (!userId || !session?.access_token) {
      error = 'No hay sesión activa'
      loading = false
      return
    }

    try {
      loading = true
      error = null

      // Obtener menuId
      const menuId = await getLatestMenuId(userId, session.access_token)
      if (!menuId) {
        error = 'No se encontró un menú'
        loading = false
      }

      // Obtener slug desde el backend (misma fuente que Compartir)
      if (menuId) {
        const menuSlug = await getMenuSlugFromApi(menuId, session.access_token)
        if (menuSlug && menuSlug.trim() !== '') {
          slug = menuSlug.trim()
          menuUrl = generateMenuUrlFromSlug(slug, currentLanguage)
        } else {
          // Sin slug: usar URL con menu_id para que igual pueda tener QR
          menuUrl = generateMenuUrlFromMenuId(menuId, currentLanguage)
        }
        const qrSize = 300
        qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=${qrSize}x${qrSize}&data=${encodeURIComponent(menuUrl!)}`
      }
    } catch (err: any) {
      console.error('Error cargando código QR:', err)
      error = err.message || 'Error al cargar el código QR'
    } finally {
      loading = false
    }
  }

  function openMenuInNewTab() {
    if (!menuUrl) return
    let url = menuUrl
    if (typeof window !== 'undefined' && (window !== window.top || new URLSearchParams(window.location.search).get('demo') === '1')) {
      url += (url.includes('?') ? '&' : '?') + 'demo=1'
    }
    openUrlInNewTab(url)
  }

  function openUrlInNewTab(url: string) {
    const w = window.open(url, '_blank')
    if (w) {
      menuWindowRef = w
      w.focus()
      setTimeout(() => { try { w.focus() } catch (_) {} }, 200)
      setTimeout(() => { try { w.focus() } catch (_) {} }, 500)
    }
  }

  function copyUrl() {
    if (menuUrl) {
      navigator.clipboard.writeText(menuUrl).then(() => {
        alert('URL copiada al portapapeles')
      }).catch(err => {
        console.error('Error copiando URL:', err)
        alert('Error al copiar la URL')
      })
    }
  }

  async function downloadQR() {
    if (!qrCodeUrl) return

    try {
      // Obtener la imagen como blob
      const response = await fetch(qrCodeUrl)
      if (!response.ok) {
        throw new Error('Error al obtener la imagen del QR')
      }
      
      const blob = await response.blob()
      
      // Crear una URL de objeto para el blob
      const blobUrl = URL.createObjectURL(blob)
      
      // Crear un enlace temporal para descargar
      const link = document.createElement('a')
      link.href = blobUrl
      link.download = `qr-menu-${slug || 'menu'}.png`
      document.body.appendChild(link)
      link.click()
      
      // Limpiar
      document.body.removeChild(link)
      URL.revokeObjectURL(blobUrl)
    } catch (err: any) {
      console.error('Error descargando QR:', err)
      alert('Error al descargar el código QR. Por favor, intenta de nuevo.')
    }
  }

  onMount(() => {
    loadQRCode()
  })
</script>

<div class="flex flex-col h-screen h-[100dvh] bg-white relative overflow-hidden">
  <div class="flex flex-col h-full">
  <!-- Header -->
  <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between flex-shrink-0 sticky top-0 z-20 bg-white">
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
    <h1 class="text-lg font-medium text-gray-900">Código QR</h1>
    <div class="w-9"></div>
  </header>

  <!-- Contenido -->
  <div class="flex-1 overflow-y-auto px-4 py-6">
    {#if loading}
      <div class="flex items-center justify-center h-full">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
          <p class="text-gray-600">Cargando código QR...</p>
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
            onclick={loadQRCode}
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            Reintentar
          </button>
        </div>
      </div>
    {:else if qrCodeUrl && menuUrl}
      <div class="max-w-2xl mx-auto">
        <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4 text-center">Código QR de tu Menú</h2>
          
          <!-- QR Code -->
          <div class="flex justify-center mb-6">
            <div class="bg-white p-4 rounded-xl border-2 border-gray-200">
              <img 
                src={qrCodeUrl} 
                alt="Código QR del menú" 
                class="w-64 h-64 md:w-80 md:h-80"
              />
            </div>
          </div>

          <!-- URL -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">URL del Menú</label>
            <div class="flex gap-2">
              <input
                type="text"
                value={menuUrl}
                readonly
                class="flex-1 px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 text-sm"
              />
              <button
                onclick={copyUrl}
                class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium"
              >
                Copiar
              </button>
            </div>
          </div>

          <!-- Botones de acción -->
          <div class="flex flex-col sm:flex-row gap-3">
            <button
              onclick={downloadQR}
              class="flex-1 px-6 py-3 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors font-medium flex items-center justify-center gap-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              Descargar QR
            </button>
            <button
              onclick={openMenuInNewTab}
              class="flex-1 px-6 py-3 bg-gray-600 hover:bg-gray-700 text-white rounded-lg transition-colors font-medium flex items-center justify-center gap-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              Abrir Menú
            </button>
          </div>

          <!-- Información -->
          <div class="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <p class="text-sm text-blue-800">
              <strong>💡 Consejo:</strong> Comparte este código QR con tus clientes para que puedan acceder fácilmente a tu menú digital desde sus teléfonos.
            </p>
          </div>
        </div>
      </div>
    {/if}
  </div>
  </div>
</div>
