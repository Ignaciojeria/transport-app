<script lang="ts">
  import { signOut, authState } from '../auth.svelte'
  import { API_BASE_URL } from '../config'
  import { t as tStore } from '../useLanguage'
  
  interface SidebarProps {
    activeSection: string
    onSectionChange: (section: string) => void
    isOpen?: boolean
    onClose?: () => void
  }

  let { activeSection, onSectionChange, isOpen = true, onClose }: SidebarProps = $props()
  
  const session = $derived(authState.session)
  const t = $derived($tStore)
  let showBuyCreditsModal = $state(false)
  
  async function handleSignOut() {
    if (confirm(t.sidebar.confirmSignOut)) {
      try {
        await signOut()
      } catch (error) {
        console.error('Error al cerrar sesión:', error)
        alert(t.sidebar.errorNoSession)
      }
    }
  }

  function handleBuyCreditsClick() {
    showBuyCreditsModal = true
  }

  function closeBuyCreditsModal() {
    showBuyCreditsModal = false
  }

  async function handleBuyCredits() {
    if (!session?.access_token) {
      alert(t.sidebar.errorNoSession)
      return
    }

    try {
      // Llamar al endpoint de checkout de MercadoPago
      const response = await fetch(`${API_BASE_URL}/checkout/mercadopago`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        const errorText = await response.text()
        console.error('Error obteniendo checkout URL:', errorText)
        alert('Error al obtener la URL de checkout. Por favor, intenta de nuevo.')
        return
      }

      const data = await response.json()
      const checkoutUrl = data.checkout_url

      if (!checkoutUrl) {
        alert('No se recibió la URL de checkout. Por favor, intenta de nuevo.')
        return
      }

      // Cerrar el modal
      closeBuyCreditsModal()

      // Redirigir a la URL de checkout
      window.open(checkoutUrl, '_blank')
    } catch (error) {
      console.error('Error comprando créditos:', error)
      alert('Error al comprar créditos. Por favor, intenta de nuevo.')
    }
  }
</script>

<div 
  class="w-64 bg-gray-900 text-white fixed left-0 top-0 z-40 shadow-xl transform transition-transform duration-300 ease-in-out md:translate-x-0 {isOpen ? 'translate-x-0' : '-translate-x-full'} flex flex-col"
  style="height: 100dvh; max-height: 100dvh;"
>
  <div class="p-6 border-b border-gray-700 flex items-center justify-between flex-shrink-0">
    <h2 class="text-xl font-bold text-white">MiCartaPro</h2>
    <!-- Botón cerrar para móvil -->
    <button
      onclick={onClose}
      class="md:hidden p-2 hover:bg-gray-800 rounded-lg transition-colors"
      aria-label={t.sidebar.closeMenu}
    >
      <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
  
  <nav class="px-4 py-4 flex-1 overflow-y-auto min-h-0">
    <button
      onclick={() => onSectionChange('menu')}
      class={`w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 ${
        activeSection === 'menu' 
          ? 'bg-blue-600 text-white shadow-md' 
          : 'text-gray-300 hover:bg-gray-800 hover:text-white'
      }`}
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.aiAssistant}</span>
    </button>

    <button
      onclick={() => onSectionChange('historial')}
      class={`w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 ${
        activeSection === 'historial' 
          ? 'bg-blue-600 text-white shadow-md' 
          : 'text-gray-300 hover:bg-gray-800 hover:text-white'
      }`}
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.history}</span>
    </button>

    <button
      onclick={() => onSectionChange('galeria')}
      class={`w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 ${
        activeSection === 'galeria' 
          ? 'bg-blue-600 text-white shadow-md' 
          : 'text-gray-300 hover:bg-gray-800 hover:text-white'
      }`}
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.gallery}</span>
    </button>

    <button
      onclick={() => onSectionChange('qr')}
      class={`w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 ${
        activeSection === 'qr' 
          ? 'bg-blue-600 text-white shadow-md' 
          : 'text-gray-300 hover:bg-gray-800 hover:text-white'
      }`}
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.qrCode}</span>
    </button>

    <button
      onclick={() => onSectionChange('ordenes')}
      class={`w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 ${
        activeSection === 'ordenes' 
          ? 'bg-blue-600 text-white shadow-md' 
          : 'text-gray-300 hover:bg-gray-800 hover:text-white'
      }`}
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.orders}</span>
    </button>

    <button
      onclick={handleBuyCreditsClick}
      class="w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 text-gray-300 hover:bg-gray-800 hover:text-white"
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span class="text-sm font-medium">Comprar créditos</span>
    </button>
  </nav>
  
  <!-- Botón de cerrar sesión al final - siempre visible -->
  <div class="p-4 border-t border-gray-700 flex-shrink-0 bg-gray-900 sticky bottom-0 md:relative" style="bottom: env(safe-area-inset-bottom, 0);">
    <button
      onclick={handleSignOut}
      class="w-full flex items-center p-3 rounded-lg transition-all duration-200 text-gray-300 hover:bg-gray-800 hover:text-white"
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.signOut}</span>
    </button>
  </div>
</div>

<!-- Modal de compra de créditos - Fuera del sidebar para que aparezca sobre toda la aplicación -->
{#if showBuyCreditsModal}
  <!-- Overlay con transición -->
  <div 
    class="fixed inset-0 bg-black bg-opacity-50 z-[100] transition-opacity duration-300"
    onclick={closeBuyCreditsModal}
    role="dialog"
    aria-modal="true"
    aria-labelledby="buy-credits-title"
  >
    <!-- Modal -->
    <div 
      class="fixed inset-0 flex items-center justify-center p-4 z-[100]"
      onclick={(e) => e.stopPropagation()}
    >
      <div 
        class="bg-white rounded-2xl shadow-2xl max-w-md w-full transform transition-all duration-300 scale-100"
        onclick={(e) => e.stopPropagation()}
      >
        <!-- Header -->
        <div class="px-6 pt-6 pb-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h2 id="buy-credits-title" class="text-2xl font-bold text-gray-900">
              💳 Comprar Créditos
            </h2>
            <button
              onclick={closeBuyCreditsModal}
              class="p-2 hover:bg-gray-100 rounded-full transition-colors"
              aria-label="Cerrar"
            >
              <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Contenido -->
        <div class="px-6 py-6">
          <div class="text-center mb-6">
            <div class="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full mb-4">
              <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">
              $3.500 CLP
            </h3>
            <p class="text-3xl font-bold text-blue-600 mb-2">
              25 Créditos
            </p>
            <p class="text-gray-600 text-sm">
              25 ediciones de menú desde nuestro agente
            </p>
          </div>

          <!-- Información adicional -->
          <div class="bg-blue-50 rounded-xl p-4 mb-6 border border-blue-100">
            <p class="text-sm text-gray-700 text-center">
              Cada crédito te permite realizar una interacción con el agente para generar o editar tu carta digital.
            </p>
          </div>

          <!-- Botones -->
          <div class="flex flex-col gap-3">
            <button
              onclick={handleBuyCredits}
              class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-2"
            >
              <span>💳</span>
              <span>Comprar ahora</span>
            </button>
            <button
              onclick={closeBuyCreditsModal}
              class="w-full px-6 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-xl transition-all font-medium text-sm"
            >
              Cancelar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
