<script>
  import MenuItem from './MenuItem.svelte';
  
  /** @type {Object} - Sección del menú */
  export let section = {
    title: '',
    items: []
  };
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<div class={`mb-8 sm:mb-10 ${className}`}>
  <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-4 sm:mb-5">
    {section.title}
  </h2>
  <div class="space-y-4 sm:space-y-5">
    {#each section.items as item}
      <MenuItem {item} />
    {/each}
  </div>
</div>

