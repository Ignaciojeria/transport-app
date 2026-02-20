<script>
  const {
    message = 'Preparando tu pedido...',
    redirectingMessage = 'Tu pedido se ha enviado correctamente',
    trackingId = null
  } = $props();

  const trackingUrl = $derived(
    trackingId ? `${typeof window !== 'undefined' ? window.location.origin : ''}/track/${encodeURIComponent(trackingId)}` : null
  );
</script>

<div class="fixed inset-0 bg-black/60 z-[70] flex items-center justify-center p-4">
  <div class="bg-white rounded-lg shadow-2xl p-8 sm:p-10 max-w-sm w-full text-center">
    <!-- Spinner animado -->
    <div class="flex justify-center mb-6">
      <div class="relative w-16 h-16 sm:w-20 sm:h-20">
        <!-- Círculo exterior giratorio -->
        <div class="absolute inset-0 border-4 border-green-200 rounded-full"></div>
        <div class="absolute inset-0 border-4 border-green-500 border-t-transparent rounded-full animate-spin"></div>
        <!-- Icono de confirmación (checkmark) -->
        <div class="absolute inset-0 flex items-center justify-center">
          <svg class="w-8 h-8 sm:w-10 sm:h-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
      </div>
    </div>
    
    <!-- Mensaje -->
    <p class="text-lg sm:text-xl font-semibold text-gray-800 mb-2">
      {message}
    </p>
    <p class="text-sm text-gray-600 mb-3">
      {redirectingMessage}
    </p>
    {#if trackingId && trackingUrl}
      <p class="text-xs text-gray-500 mb-2">Código: <span class="font-mono font-semibold">{trackingId}</span></p>
      <a
        href={trackingUrl}
        target="_blank"
        rel="noopener noreferrer"
        class="inline-block text-sm text-emerald-600 hover:text-emerald-700 font-medium underline"
      >
        Ver estado del pedido →
      </a>
    {/if}
  </div>
</div>

<style>
  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
  
  .animate-spin {
    animation: spin 1s linear infinite;
  }
</style>
