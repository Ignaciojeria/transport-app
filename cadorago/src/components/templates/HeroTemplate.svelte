<script>
  import BrandHero from '../organisms/BrandHero.svelte';
  
  /** @type {string} - Background color class */
  export let bgColor = 'bg-[#E8E4D9]';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<div class={`min-h-screen ${bgColor} ${className}`}>
  <div class="container mx-auto px-4 sm:px-6 lg:px-12 py-4 sm:py-8 lg:py-16 max-w-[1600px]">
    <!-- Hero Section with BrandHero organism -->
    <section class="flex items-center justify-center min-h-[60vh] sm:min-h-[70vh] lg:min-h-[85vh]">
      <BrandHero />
    </section>
    
    <!-- Slot for additional content -->
    <slot />
  </div>
</div>
