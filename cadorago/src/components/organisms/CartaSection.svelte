<script>
  import MenuSection from '../molecules/MenuSection.svelte';
  import { t } from '../../lib/useLanguage';
  import { getBaseText } from '../../lib/multilingual';
  
  const { carta = [], templateName = 'hero', className = '' } = $props();
  
  // Filtrar secciones que tengan contenido (título y al menos un item con contenido)
  const seccionesConContenido = $derived.by(() => {
    if (!carta || !Array.isArray(carta) || carta.length === 0) {
      if (typeof window !== 'undefined') {
        console.log('[CartaSection] Carta vacía o no es array:', carta);
      }
      return [];
    }
    
    if (typeof window !== 'undefined') {
      console.log('[CartaSection] Procesando', carta.length, 'secciones');
    }
    
    const resultado = carta.filter(section => {
      if (!section) {
        return false;
      }
      
      // Verificar que tenga título (puede ser string o MultilingualText)
      if (!section.title) {
        if (typeof window !== 'undefined') {
          console.log('[CartaSection] Sección sin title:', section);
        }
        return false;
      }
      
      const titleBase = getBaseText(section.title);
      if (typeof window !== 'undefined') {
        console.log('[CartaSection] Title base:', titleBase, 'de:', section.title);
      }
      
      if (!titleBase || !titleBase.trim()) {
        return false;
      }
      
      // Verificar que tenga items con contenido
      if (!section.items || !Array.isArray(section.items) || section.items.length === 0) {
        if (typeof window !== 'undefined') {
          console.log('[CartaSection] Sección sin items válidos:', section.items);
        }
        return false;
      }
      
      if (typeof window !== 'undefined') {
        console.log('[CartaSection] Sección tiene', section.items.length, 'items');
      }
      
      const itemsConContenido = section.items.filter(item => {
        if (!item) {
          return false;
        }
        
        // Verificar que tenga título (puede ser string o MultilingualText)
        if (!item.title) {
          if (typeof window !== 'undefined') {
            console.log('[CartaSection] Item sin title:', item);
          }
          return false;
        }
        
        const itemTitleBase = getBaseText(item.title);
        if (typeof window !== 'undefined') {
          console.log('[CartaSection] Item title base:', itemTitleBase, 'de:', item.title);
        }
        const tieneTitulo = itemTitleBase && itemTitleBase.trim().length > 0;
        if (typeof window !== 'undefined' && !tieneTitulo) {
          console.log('[CartaSection] Item sin título válido');
        }
        return tieneTitulo;
      });
      
      if (typeof window !== 'undefined') {
        console.log('[CartaSection] Items con contenido:', itemsConContenido.length);
      }
      return itemsConContenido.length > 0;
    });
    
    if (typeof window !== 'undefined') {
      console.log('[CartaSection] Secciones filtradas:', resultado.length, 'de', carta.length);
    }
    return resultado;
  });
  
  const tieneElementos = $derived(seccionesConContenido.length > 0);
  
  // Debug: mostrar información en consola
  $effect(() => {
    if (typeof window !== 'undefined') {
      console.log('[CartaSection] Props carta recibidos:', carta);
      console.log('[CartaSection] Secciones con contenido:', seccionesConContenido);
      console.log('[CartaSection] Tiene elementos:', tieneElementos);
    }
  });
</script>

<section class={`space-y-10 sm:space-y-12 ${className}`}>
  {#if tieneElementos}
    {#each seccionesConContenido as section}
      <MenuSection {section} useModernLayout={templateName === 'modern'} />
    {/each}
  {:else}
    <p class="text-base sm:text-lg text-gray-500">
      {$t.menu.noItems}
    </p>
  {/if}
</section>

