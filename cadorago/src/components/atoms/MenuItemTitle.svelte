<script>
  import { getMultilingualText } from '../../lib/multilingual';
  
  const { 
    title = '', 
    className = '' 
  } = $props();
  
  const displayTitle = $derived(getMultilingualText(title));
</script>

<h3 class={`text-xl sm:text-2xl lg:text-3xl font-bold text-gray-800 ${className}`}>
  {displayTitle}
</h3>

