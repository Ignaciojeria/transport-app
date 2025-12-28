<script>
  import BrandHero from '../organisms/BrandHero.svelte';
  
  /** @type {string} - Background color class */
  export let bgColor = 'bg-[#E8E4D9]';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<div class={`min-h-screen ${bgColor} ${className}`}>
  <!-- Hero Section with BrandHero organism - sin padding para que ocupe todo el ancho -->
  <section class="w-full">
    <BrandHero />
  </section>
  
  <!-- Contenedor para el contenido adicional con padding -->
  <div class="container mx-auto px-2 sm:px-4 lg:px-12 py-4 sm:py-8 lg:py-16 max-w-[1600px]">
    <!-- Slot for additional content -->
    <slot />
  </div>
</div>
