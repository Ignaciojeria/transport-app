<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  
  const { className = '' } = $props();
  
  const restaurantData = $derived(restaurantDataStore.value);
  const coverImage = $derived(restaurantData?.coverImage || '');
</script>

<div class={`w-full overflow-hidden ${className}`}>
  <!-- Cover Image como foto de portada ocupando todo el ancho sin recortarse -->
  {#if coverImage}
    <img 
      src={coverImage} 
      alt="Portada de Cadorago" 
      class="w-full h-auto"
      style="width: 100%; height: auto; display: block; object-fit: contain; object-position: center top;"
    />
  {/if}
</div>
