<script>
  import HeroTemplate from '../templates/HeroTemplate.svelte';
</script>

<HeroTemplate>
  <!-- Contenido adicional opcional debajo del hero -->
  <section class="mt-8 sm:mt-12 lg:mt-16 text-center px-4">
    <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-4 sm:p-6 lg:p-8 max-w-2xl mx-auto">
      <h2 class="text-xl sm:text-2xl font-bold text-gray-800 mb-3 sm:mb-4">
        Bienvenido a Cadorago
      </h2>
      <p class="text-sm sm:text-base text-gray-600 leading-relaxed">
        Un viaje por la cocina chilena que combina tradición, sabor y pasión en cada plato.
      </p>
    </div>
  </section>
</HeroTemplate>
