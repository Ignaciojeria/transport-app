<script>
  import WhatsAppIcon from '../atoms/WhatsAppIcon.svelte';
  
  const { 
    whatsapp = '',
    className = ''
  } = $props();
</script>

<div
  class={`flex items-center gap-3 px-6 py-3 sm:px-8 sm:py-4 lg:px-10 lg:py-5 bg-green-500 text-white rounded-lg ${className}`}
>
  <WhatsAppIcon className="w-6 h-6 sm:w-7 sm:h-7 lg:w-8 lg:h-8" />
  <span class="font-medium text-base sm:text-lg lg:text-xl">{whatsapp}</span>
</div>

