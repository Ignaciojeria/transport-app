<script lang="ts">
  import { onMount } from 'svelte'
  import { initAuth, authState } from './lib/auth.svelte'
  import { initLanguage, t as tStore, loading as langLoadingStore } from './lib/useLanguage'
  import MenuChat from './lib/components/MenuChat.svelte'
  import PaymentSuccess from './lib/components/PaymentSuccess.svelte'
  import Sidebar from './lib/components/Sidebar.svelte'
  import PhotoUpload from './lib/components/PhotoUpload.svelte'
  import MenuHistory from './lib/components/MenuHistory.svelte'
  import MenuQRCode from './lib/components/MenuQRCode.svelte'

  // Usar valores derivados reactivos en el componente
  let user = $derived(authState.user)
  let session = $derived(authState.session)
  let authLoading = $derived(authState.loading)

  // Estado de la sección activa
  let activeSection = $state('menu')
  
  // Estado del sidebar (abierto/cerrado en móvil, siempre abierto en desktop)
  let sidebarOpen = $state(false)
  
  // Función para verificar si estamos en móvil
  function isMobile(): boolean {
    if (typeof window === 'undefined') return false
    return window.innerWidth < 768
  }

  // Función para cambiar de sección
  function handleSectionChange(section: string) {
    activeSection = section
    // Cerrar sidebar en móvil después de seleccionar una sección
    if (isMobile()) {
      sidebarOpen = false
    }
  }
  
  // Función para toggle del sidebar
  function toggleSidebar() {
    sidebarOpen = !sidebarOpen
  }
  
  // Función para cerrar el sidebar
  function closeSidebar() {
    sidebarOpen = false
  }
  
  // Determinar si el sidebar debe mostrarse (siempre en desktop, condicional en móvil)
  const isSidebarVisible = $derived(!isMobile() || sidebarOpen)

  // Detectar si estamos en la página de éxito de pago
  let isPaymentSuccess = $derived(() => {
    if (typeof window === 'undefined') return false
    const urlParams = new URLSearchParams(window.location.search)
    return urlParams.get('payment') === 'success' || urlParams.get('success') === 'true'
  })

  onMount(() => {
    initLanguage()
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
    
    // Si no hay usuario autenticado, redirigir a auth-ui
    // Esto se ejecuta después de que initAuth termine de cargar
    setTimeout(() => {
      if (!user && !authLoading) {
        const isLocalDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
        const authUiUrl = isLocalDev ? 'http://localhost:3003' : 'https://auth.micartapro.com'
        window.location.replace(authUiUrl)
      }
    }, 100)
  })
</script>

<div class="min-h-screen bg-gray-50">
  {#if isPaymentSuccess()}
    <!-- Página de éxito de pago -->
    <PaymentSuccess />
  {:else if authLoading || $langLoadingStore}
    <div class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
        <p class="text-gray-600">{$tStore.app.loading}</p>
      </div>
    </div>
  {:else if user}
    <!-- Content con Sidebar -->
    <div class="flex h-screen h-[100dvh] relative">
      <!-- Overlay para móvil -->
      {#if sidebarOpen}
        <div 
          class="fixed inset-0 bg-black bg-opacity-50 z-30 md:hidden"
          onclick={closeSidebar}
          role="button"
          tabindex="0"
          onkeydown={(e) => e.key === 'Escape' && closeSidebar()}
        ></div>
      {/if}
      
      <!-- Sidebar -->
      <Sidebar 
        activeSection={activeSection} 
        onSectionChange={handleSectionChange}
        isOpen={isSidebarVisible}
        onClose={closeSidebar}
      />
      
      <!-- Contenido principal -->
      <div class="flex-1 md:ml-64 overflow-hidden bg-gray-50">
        {#if activeSection === 'menu'}
          <MenuChat onMenuClick={toggleSidebar} />
        {:else if activeSection === 'historial'}
          <MenuHistory onMenuClick={toggleSidebar} />
        {:else if activeSection === 'galeria'}
          <PhotoUpload onMenuClick={toggleSidebar} />
        {:else if activeSection === 'qr'}
          <MenuQRCode onMenuClick={toggleSidebar} />
        {/if}
      </div>
    </div>
  {:else}
    <!-- Vista de no autenticado - redirigir automáticamente a auth-ui -->
    <div class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
        <p class="text-gray-600">Redirigiendo al inicio de sesión...</p>
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
