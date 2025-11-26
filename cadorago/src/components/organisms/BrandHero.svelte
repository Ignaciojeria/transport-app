<script>
  import BrandLogo from '../molecules/BrandLogo.svelte';
  import Slogan from '../atoms/Slogan.svelte';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<div class={`grid grid-cols-1 min-[1180px]:grid-cols-2 gap-4 sm:gap-6 md:gap-8 min-[1180px]:gap-24 items-center ${className}`}>
  <!-- Primera columna: BrandLogo con tamaños responsivos -->
  <div class="flex justify-center min-[1180px]:justify-end w-full">
    <!-- Móvil muy pequeño: navbar (hasta 375px) -->
    <div class="block min-[376px]:hidden">
      <BrandLogo size="navbar" />
    </div>
    <!-- Móvil pequeño: footer (376px - 534px) -->
    <div class="hidden min-[376px]:block min-[535px]:hidden">
      <BrandLogo size="footer" />
    </div>
    <!-- Tablet pequeño: hero (535px - 1023px) -->
    <div class="hidden min-[535px]:block min-[1024px]:hidden">
      <BrandLogo size="hero" />
    </div>
    <!-- Desktop pequeño: xl (1024px - 1560px) -->
    <div class="hidden min-[1024px]:block min-[1561px]:hidden">
      <BrandLogo size="xl" />
    </div>
    <!-- Desktop grande: xxl (desde 1561px) -->
    <div class="hidden min-[1561px]:block">
      <BrandLogo size="xxl" />
    </div>
  </div>
  
  <!-- Segunda columna: Slogan con altura responsiva -->
  <div class="flex justify-center min-[1180px]:justify-start items-center w-full">
    <!-- Slogan con tamaños que coinciden con los breakpoints del BrandLogo -->
    <!-- Margen negativo moderado en mobile puro para evitar solapamiento, más agresivo en breakpoints intermedios -->
    <!-- Margen izquierdo negativo solo en layout de 2 columnas para moverlo a la izquierda -->
    <Slogan className="h-56 min-[376px]:h-72 min-[535px]:h-96 min-[1024px]:h-[28rem] min-[1561px]:h-[48rem] w-auto max-w-full -mt-28 min-[376px]:-mt-40 min-[535px]:-mt-60 min-[1024px]:-mt-80 min-[1180px]:-mt-[9rem] min-[1180px]:-ml-16" />
  </div>
</div>
