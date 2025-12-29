<script>
  import MenuSection from '../molecules/MenuSection.svelte';
  import { t } from '../../lib/useLanguage';
  
  const { carta = [], className = '' } = $props();
  
  // Filtrar secciones que tengan contenido (título y al menos un item con contenido)
  const seccionesConContenido = $derived(
    carta?.filter(section => {
      if (!section || !section.title || typeof section.title !== 'string' || !section.title.trim()) {
        return false;
      }
      // Verificar que tenga items con contenido
      const itemsConContenido = section.items?.filter(item => {
        return item && item.title && typeof item.title === 'string' && item.title.trim().length > 0;
      }) || [];
      return itemsConContenido.length > 0;
    }) || []
  );
  
  const tieneElementos = $derived(seccionesConContenido.length > 0);
</script>

<section class={`space-y-8 ${className}`}>
  {#if tieneElementos}
    {#each seccionesConContenido as section}
      <MenuSection {section} />
    {/each}
  {:else}
    <p class="text-lg sm:text-xl lg:text-2xl text-gray-500 italic">
      {$t.menu.noItems}
    </p>
  {/if}
</section>

