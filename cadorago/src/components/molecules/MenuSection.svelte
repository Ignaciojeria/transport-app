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

<div class={`${className}`}>
  <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-6 sm:mb-8 capitalize">
    {section.title}
  </h2>
  <div class="space-y-4 sm:space-y-5 lg:space-y-6">
    {#each section.items as item}
      <MenuItem {item} />
    {/each}
  </div>
</div>

