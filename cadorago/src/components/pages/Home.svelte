<script>
  import HeroTemplate from '../templates/HeroTemplate.svelte';
  import HorariosSection from '../organisms/HorariosSection.svelte';
  import ContactSection from '../organisms/ContactSection.svelte';
  import CartaSection from '../organisms/CartaSection.svelte';
  import FloatingCart from '../organisms/FloatingCart.svelte';
  import Footer from '../organisms/Footer.svelte';
  import { getRestaurantData } from '../../services/restaurantData.js';
  
  const restaurantData = getRestaurantData();
</script>

<HeroTemplate>
  <!-- Sección de información del restaurante -->
  <section class="-mt-16 sm:-mt-20 lg:-mt-24 px-2 sm:px-4 lg:px-12">
    <div class="max-w-[1600px] mx-auto space-y-12 sm:space-y-16 lg:space-y-20">
      <!-- Horarios y Contacto en grid para breakpoints grandes -->
      <div class="grid grid-cols-1 min-[1180px]:grid-cols-2 gap-8 sm:gap-12 lg:gap-16">
        <!-- Horarios de atención -->
        <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-8 sm:p-10 lg:p-12">
          <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-6 sm:mb-8">
            Horarios de Atención
          </h2>
          <HorariosSection horarios={restaurantData.horariosDeAtencion} />
        </div>
        
        <!-- Contacto -->
        <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-8 sm:p-10 lg:p-12">
          <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-6 sm:mb-8">
            Contacto
          </h2>
          <ContactSection whatsapp={restaurantData.contacto.whatssap} />
        </div>
      </div>
      
      <!-- Carta -->
      <div class="bg-white/80 backdrop-blur-sm rounded-lg shadow-lg p-8 sm:p-10 lg:p-12 mb-24 sm:mb-28 lg:mb-32">
        <h2 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-8 sm:mb-10 lg:mb-12">
          Nuestra Carta
        </h2>
        <CartaSection carta={restaurantData.carta} />
      </div>
    </div>
  </section>
  
  <!-- Carrito flotante -->
  <FloatingCart />
</HeroTemplate>

<!-- Footer -->
<Footer />
