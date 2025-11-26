<script>
  import ContactInfo from '../molecules/ContactInfo.svelte';
  
  /** @type {string} - Número de WhatsApp */
  export let whatsapp = '';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<section class={`flex justify-center sm:justify-start ${className}`}>
  <ContactInfo {whatsapp} />
</section>

