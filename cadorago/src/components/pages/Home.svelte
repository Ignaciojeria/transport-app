<script>
  import { onMount } from 'svelte';
  import HeroTemplate from '../templates/HeroTemplate.svelte';
  import ModernTemplate from '../templates/ModernTemplate.svelte';
  import HorariosSection from '../organisms/HorariosSection.svelte';
  import ContactSection from '../organisms/ContactSection.svelte';
  import CartaSection from '../organisms/CartaSection.svelte';
  import FloatingCart from '../organisms/FloatingCart.svelte';
  import Footer from '../organisms/Footer.svelte';
  import MetaTags from '../organisms/MetaTags.svelte';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { initLanguage, t } from '../../lib/useLanguage';
  
  const restaurantData = $derived(restaurantDataStore.value);
  const loading = $derived(restaurantDataStore.loading);
  const error = $derived(restaurantDataStore.error);
  
  // Leer query param ?template= para seleccionar el template
  let templateName = $state('hero'); // Por defecto: hero (original)
  
  onMount(() => {
    initLanguage();
    
    // Leer query param de la URL
    if (typeof window !== 'undefined') {
      const params = new URLSearchParams(window.location.search);
      const templateParam = params.get('template');
      if (templateParam) {
        templateName = templateParam.toLowerCase();
      }
    }
  });
  
  // Seleccionar el componente template según el query param
  const TemplateComponent = $derived(() => {
    switch (templateName) {
      case 'modern':
        return ModernTemplate;
      case 'hero':
      default:
        return HeroTemplate;
    }
  });
  
  const CurrentTemplate = $derived(TemplateComponent());
</script>

<MetaTags />

{#if CurrentTemplate}
  {@const Template = CurrentTemplate}
  <Template>
  {#if loading}
    <!-- Estado de carga -->
    <section class="px-2 sm:px-4 lg:px-12 pt-8 sm:pt-12 lg:pt-16">
      <div class="max-w-[1600px] mx-auto text-center py-20">
        <p class="text-xl sm:text-2xl text-gray-600">{$t.home.loading}</p>
      </div>
    </section>
  {:else if error}
    <!-- Estado de error -->
    <section class="px-2 sm:px-4 lg:px-12 pt-8 sm:pt-12 lg:pt-16">
      <div class="max-w-[1600px] mx-auto text-center py-20">
        <p class="text-xl sm:text-2xl text-red-600 mb-4">{$t.home.errorLoading}</p>
        <p class="text-base sm:text-lg text-gray-600">{error}</p>
      </div>
    </section>
  {:else if restaurantData}
    <!-- Sección de información del restaurante -->
    <section class="px-4 sm:px-6 lg:px-8 pt-2 sm:pt-4">
      <div class="max-w-4xl mx-auto space-y-8 sm:space-y-10">
        <!-- Carta (Menú) - Al inicio -->
        <div>
          <CartaSection carta={restaurantData.menu || []} />
        </div>
        
        <!-- Horarios de atención y Contacto - Solo para template hero (no modern) -->
        {#if templateName !== 'modern'}
          <!-- Horarios de atención - En el medio -->
          <div>
            <h2 class="text-lg sm:text-xl font-semibold text-gray-900 mb-4 sm:mb-6 uppercase tracking-wide">
              {$t.home.businessHours}
            </h2>
            <HorariosSection horarios={restaurantData.businessInfo?.businessHours || []} />
          </div>
          
          <!-- Contacto - Al final -->
          <div class="mb-24 sm:mb-28 lg:mb-32">
            <h2 class="text-lg sm:text-xl font-semibold text-gray-900 mb-4 sm:mb-6 uppercase tracking-wide">
              {$t.home.contact}
            </h2>
            <ContactSection whatsapp={restaurantData.businessInfo?.whatsapp || ''} />
          </div>
        {/if}
      </div>
    </section>
  {/if}
  
  <!-- Carrito flotante -->
  <FloatingCart />
  </Template>
{/if}

<!-- Footer solo para template hero (ModernTemplate ya lo incluye) -->
{#if templateName === 'hero'}
  <Footer />
{/if}
