<script>
  import HeroTemplate from '../templates/HeroTemplate.svelte';
  import HorariosSection from '../organisms/HorariosSection.svelte';
  import ContactSection from '../organisms/ContactSection.svelte';
  import CartaSection from '../organisms/CartaSection.svelte';
  import FloatingCart from '../organisms/FloatingCart.svelte';
  import Footer from '../organisms/Footer.svelte';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  
  const restaurantData = $derived(restaurantDataStore.value);
  const loading = $derived(restaurantDataStore.loading);
  const error = $derived(restaurantDataStore.error);
</script>

<HeroTemplate>
  {#if loading}
    <!-- Estado de carga -->
    <section class="px-2 sm:px-4 lg:px-12 pt-8 sm:pt-12 lg:pt-16">
      <div class="max-w-[1600px] mx-auto text-center py-20">
        <p class="text-xl sm:text-2xl text-gray-600">Cargando datos del restaurante...</p>
      </div>
    </section>
  {:else if error}
    <!-- Estado de error -->
    <section class="px-2 sm:px-4 lg:px-12 pt-8 sm:pt-12 lg:pt-16">
      <div class="max-w-[1600px] mx-auto text-center py-20">
        <p class="text-xl sm:text-2xl text-red-600 mb-4">Error al cargar los datos</p>
        <p class="text-base sm:text-lg text-gray-600">{error}</p>
      </div>
    </section>
  {:else if restaurantData}
    <!-- Sección de información del restaurante -->
    <section class="px-2 sm:px-4 lg:px-12 pt-8 sm:pt-12 lg:pt-16">
      <div class="max-w-[1600px] mx-auto space-y-12 sm:space-y-16 lg:space-y-20">
        <!-- Carta (Menú) - Al inicio -->
        <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-8 sm:p-10 lg:p-12">
          <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-8 sm:mb-10 lg:mb-12">
            Nuestra Carta
          </h2>
          <CartaSection carta={restaurantData.menu || []} />
        </div>
        
        <!-- Horarios de atención - En el medio -->
        <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-8 sm:p-10 lg:p-12">
          <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-6 sm:mb-8">
            Horarios de Atención
          </h2>
          <HorariosSection horarios={restaurantData.businessInfo?.businessHours || []} />
        </div>
        
        <!-- Contacto - Al final -->
        <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-8 sm:p-10 lg:p-12 mb-24 sm:mb-28 lg:mb-32">
          <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-6 sm:mb-8">
            Contacto
          </h2>
          <ContactSection whatsapp={restaurantData.businessInfo?.whatsapp || ''} />
        </div>
      </div>
    </section>
  {/if}
  
  <!-- Carrito flotante -->
  <FloatingCart />
</HeroTemplate>

<!-- Footer -->
<Footer />
