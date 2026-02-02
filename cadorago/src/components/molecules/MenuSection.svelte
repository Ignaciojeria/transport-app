<script>
  import MenuItem from './MenuItem.svelte';
  import ModernMenuItem from '../templates/ModernMenuItem.svelte';
  import { getMultilingualText, getBaseText } from '../../lib/multilingual';
  
  const { 
    section = {
      title: '',
      items: []
    },
    className = '',
    useModernLayout = false
  } = $props();
  
  // Template modern: viene del contrato (presentationStyle) vía prop useModernLayout.
  // El query param ?template= está deprecado.
  const isModernTemplate = $derived(() => {
    if (typeof window === 'undefined') return useModernLayout;
    return useModernLayout || (typeof document !== 'undefined' && document.body.classList.contains('modern-template'));
  });
  
  const sectionTitle = $derived(getMultilingualText(section.title));
  // Slug para enlace y scroll desde la nav de categorías (mismo criterio que ModernTemplate)
  const sectionId = $derived(getBaseText(section.title).toLowerCase().replace(/\s+/g, '-'));
</script>

<div id={sectionId || undefined} class={`mb-8 sm:mb-10 section-anchor ${className}`}>
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

<style>
  /* Al hacer scroll desde la nav de categorías, dejar espacio para header + nav sticky */
  :global(.modern-template) .section-anchor {
    scroll-margin-top: 10rem;
  }
</style>
