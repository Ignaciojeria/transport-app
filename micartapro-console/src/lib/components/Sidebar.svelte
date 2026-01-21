<script lang="ts">
  interface SidebarProps {
    activeSection: string
    onSectionChange: (section: string) => void
    isOpen?: boolean
    onClose?: () => void
  }

  let { activeSection, onSectionChange, isOpen = true, onClose }: SidebarProps = $props()
</script>

<div 
  class="w-64 bg-gray-900 text-white h-screen fixed left-0 top-0 z-40 shadow-xl transform transition-transform duration-300 ease-in-out md:translate-x-0 {isOpen ? 'translate-x-0' : '-translate-x-full'}"
>
  <div class="p-6 border-b border-gray-700 flex items-center justify-between">
    <h2 class="text-xl font-bold text-white">MiCartaPro</h2>
    <!-- Botón cerrar para móvil -->
    <button
      onclick={onClose}
      class="md:hidden p-2 hover:bg-gray-800 rounded-lg transition-colors"
      aria-label="Cerrar menú"
    >
      <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
  
  <nav class="px-4 py-4">
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
      <span class="text-sm font-medium">Menú</span>
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
      <span class="text-sm font-medium">Galería</span>
    </button>
  </nav>
</div>
