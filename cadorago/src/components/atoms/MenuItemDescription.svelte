<script>
  import { getMultilingualText } from '../../lib/multilingual';
  
  const { 
    description = [], 
    className = '',
    /** Si true, muestra cada elemento como ítem de lista; si false, como párrafo único */
    asList = true
  } = $props();
  
  const items = $derived(
    Array.isArray(description) 
      ? description.map(d => getMultilingualText(d)).filter(Boolean)
      : []
  );
</script>

{#if items.length > 0}
  {#if asList}
    <ul class={`list-disc list-inside space-y-2 text-sm sm:text-base text-gray-600 italic mt-2 sm:mt-3 ${className}`}>
      {#each items as item}
        <li>{item}</li>
      {/each}
    </ul>
  {:else}
    <p class={`text-sm sm:text-base text-gray-600 mt-2 sm:mt-3 ${className}`}>
      {items.join(' ')}
    </p>
  {/if}
{/if}

