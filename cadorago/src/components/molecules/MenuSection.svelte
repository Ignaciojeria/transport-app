<script>
  import MenuItem from './MenuItem.svelte';
  import ModernMenuItem from '../templates/ModernMenuItem.svelte';
  import { getMultilingualText } from '../../lib/multilingual';
  
  const { 
    section = {
      title: '',
      items: []
    },
    className = '',
    useModernLayout = false
  } = $props();
  
  // Detectar si estamos en template modern por la clase del body o query param
  const isModernTemplate = $derived(() => {
    if (typeof window === 'undefined') return useModernLayout;
    const params = new URLSearchParams(window.location.search);
    return useModernLayout || params.get('template') === 'modern' || (typeof document !== 'undefined' && document.body.classList.contains('modern-template'));
  });
  
  const sectionTitle = $derived(getMultilingualText(section.title));
</script>

<div class={`mb-8 sm:mb-10 ${className}`}>
  <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-4 sm:mb-5 {isModernTemplate() ? 'menu-category-title' : ''}">
    {sectionTitle}
  </h2>
  {#if isModernTemplate()}
    <!-- Layout moderno: grid con cards estilo chef -->
    <div class="items-grid">
      {#each section.items as item}
        <ModernMenuItem {item} />
      {/each}
    </div>
  {:else}
    <!-- Layout original: lista vertical -->
    <div class="space-y-4 sm:space-y-5">
      {#each section.items as item}
        <MenuItem {item} />
      {/each}
    </div>
  {/if}
</div>

