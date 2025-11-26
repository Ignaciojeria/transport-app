<script>
  /** @type {string} - Descripción del item */
  export let description = '';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

{#if description}
  <p class={`text-base sm:text-lg lg:text-xl text-gray-600 mt-2 sm:mt-3 ${className}`}>
    {description}
  </p>
{/if}

