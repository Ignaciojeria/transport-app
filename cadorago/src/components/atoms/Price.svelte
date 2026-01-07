<script>
  /** @type {number} - Precio a mostrar */
  export let price = 0;
  
  /** @type {string} - Additional custom classes */
  export let className = '';
  
  $: formattedPrice = new Intl.NumberFormat('es-CL', {
    style: 'currency',
    currency: 'CLP',
    minimumFractionDigits: 0
  }).format(price);
</script>

<span class={`font-bold text-gray-900 text-base sm:text-lg ${className}`}>
  {formattedPrice}
</span>

