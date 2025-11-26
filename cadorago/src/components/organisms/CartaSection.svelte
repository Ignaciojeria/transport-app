<script>
  import MenuSection from '../molecules/MenuSection.svelte';
  
  /** @type {Array} - Array de secciones de la carta */
  export let carta = [];
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<section class={`space-y-8 ${className}`}>
  {#each carta as section}
    <MenuSection {section} />
  {/each}
</section>

