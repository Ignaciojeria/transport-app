<script>
  import { getMultilingualText } from '../../lib/multilingual';
  
  const { 
    description = '', 
    className = '' 
  } = $props();
  
  const displayDescription = $derived(getMultilingualText(description));
</script>

{#if displayDescription}
  <p class={`text-base sm:text-lg lg:text-xl text-gray-600 mt-2 sm:mt-3 ${className}`}>
    {displayDescription}
  </p>
{/if}

