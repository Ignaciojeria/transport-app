<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  
  const { className = '' } = $props();
  
  const restaurantData = $derived(restaurantDataStore.value);
  const footerImage = $derived(restaurantData?.footerImage || '');
</script>

<footer class={`bg-white pt-0 pb-8 sm:pb-12 lg:pb-16 ${className}`}>
  <div class="container mx-auto px-4 sm:px-6 lg:px-12 max-w-[1600px]">
    <div class="flex flex-col items-center -mt-8 sm:-mt-10 lg:-mt-12">
      {#if footerImage}
        <img 
          src={footerImage} 
          alt="Cadorago Logo" 
          class="h-64 w-auto"
        />
      {/if}
    </div>
  </div>
</footer>

