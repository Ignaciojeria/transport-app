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

  async function handleMyPlan() {
    if (!session?.access_token) {
      alert(t.sidebar.errorNoSession)
      return
    }

    try {
      // Llamar al endpoint de customer portal
      const response = await fetch(`${API_BASE_URL}/customer-portal`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        if (response.status === 404) {
          alert(t.sidebar.errorNoActiveSubscription)
          return
        }
        const errorText = await response.text()
        console.error('Error obteniendo portal del consumidor:', errorText)
        alert(t.sidebar.errorGettingPortal)
        return
      }

      const data = await response.json()
      const portalUrl = data.customer_portal_link

      if (!portalUrl) {
        alert(t.sidebar.errorNoPortalUrl)
        return
      }

      // Redirigir al portal del consumidor
      window.open(portalUrl, '_blank')
    } catch (error) {
      console.error('Error accediendo al portal:', error)
      alert(t.sidebar.errorAccessingPortal)
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
      onclick={handleMyPlan}
      class="w-full flex items-center p-3 rounded-lg transition-all duration-200 mb-2 text-gray-300 hover:bg-gray-800 hover:text-white"
    >
      <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
      </svg>
      <span class="text-sm font-medium">{t.sidebar.myPlan}</span>
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
