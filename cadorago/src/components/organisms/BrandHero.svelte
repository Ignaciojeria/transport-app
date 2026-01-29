<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import ImageWithLoader from '../atoms/ImageWithLoader.svelte';
  
  const { className = '' } = $props();
  
  const restaurantData = $derived(restaurantDataStore.value);
  const coverImage = $derived(restaurantData?.coverImage || '');
</script>

<div class={`w-full overflow-hidden ${className}`}>
  <!-- Cover Image como foto de portada ocupando todo el ancho sin recortarse -->
  {#if coverImage}
    <div class="w-full">
      <ImageWithLoader 
        src={coverImage} 
        alt="Portada de Cadorago" 
        className="w-full h-auto"
        fallbackIcon={false}
      />
    </div>
  {/if}
</div>
