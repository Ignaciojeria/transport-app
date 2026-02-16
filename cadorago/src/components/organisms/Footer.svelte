<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  
  const { className = '' } = $props();
  
  const restaurantData = $derived(restaurantDataStore.value);
  const footerImage = $derived(restaurantData?.footerImage || '');
  const businessName = $derived(restaurantData?.businessInfo?.businessName || '');
  const whatsapp = $derived(restaurantData?.businessInfo?.whatsapp || '');
  const businessHours = $derived(restaurantData?.businessInfo?.businessHours || []);
  
  // Formatear horarios para mostrar en el footer
  const formattedHours = $derived(
    businessHours.length === 0 ? '' : businessHours.join(' | ')
  );
  
  // Obtener año actual
  const currentYear = $derived(new Date().getFullYear());
</script>

<footer class={`bg-white pt-0 pb-4 sm:pb-6 ${className}`}>
  <div class="container mx-auto px-4 sm:px-6 lg:px-12 max-w-[1600px]">
    {#if footerImage}
      <div class="flex flex-col items-center -mt-8 sm:-mt-10 lg:-mt-12 mb-4">
        <img 
          src={footerImage} 
          alt="Cadorago Logo" 
          class="h-64 w-auto"
        />
      </div>
    {/if}
    
    <!-- Copyright -->
    {#if businessName}
      <div class="text-center text-gray-400 text-xs sm:text-sm mb-3">
        © {currentYear} {businessName}. Todos os direitos reservados.
      </div>
    {/if}
    
    <!-- Horarios y Contacto -->
    <div class="flex flex-wrap items-center justify-center gap-3 sm:gap-4 text-gray-400 text-xs sm:text-sm">
      {#if formattedHours}
        <div class="flex items-center gap-1.5">
          <span>⏰</span>
          <span>{formattedHours}</span>
        </div>
      {/if}
      
      {#if whatsapp}
        <div class="flex items-center gap-1.5">
          <span>💬</span>
          <span>WhatsApp</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span>📞</span>
          <span>Telefone</span>
        </div>
      {/if}
      
      <!-- Ubicación - puedes agregar el campo de dirección cuando esté disponible -->
      <div class="flex items-center gap-1.5">
        <span>📍</span>
        <span>Santiago, Chile</span>
      </div>
      <a href={restaurantData?.id ? `/track?m=${encodeURIComponent(restaurantData.id)}` : '/track'} class="flex items-center gap-1.5 text-emerald-600 hover:text-emerald-700 transition-colors">
        <span>📦</span>
        <span>Seguimiento de pedido</span>
      </a>
    </div>
  </div>
</footer>

<style>
  footer {
    font-family: inherit;
  }
</style>

