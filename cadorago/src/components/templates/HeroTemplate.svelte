<script>
  import BrandHero from '../organisms/BrandHero.svelte';
  import './hero-template.css';

  /** @type {string} - Background color class */
  export let bgColor = 'bg-white';

  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<div class="hero-template min-h-screen {bgColor} {className}">
  <!-- Hero Section with BrandHero - ámbito solo de esta template -->
  <section class="hero-template__hero w-full">
    <BrandHero />
  </section>

  <!-- Contenedor del contenido adicional - ámbito solo de esta template -->
  <div class="hero-template__content container mx-auto px-2 sm:px-4 lg:px-12 pt-2 sm:pt-4 pb-4 sm:pb-8 lg:pb-16 max-w-[1600px]">
    <slot />
  </div>
</div>

