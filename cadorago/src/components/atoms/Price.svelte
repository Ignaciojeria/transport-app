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

<span class={`font-semibold text-gray-800 text-lg sm:text-xl lg:text-2xl ${className}`}>
  {formattedPrice}
</span>

