<script>
  /** @type {string} - Título del item */
  export let title = '';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<h3 class={`text-xl sm:text-2xl lg:text-3xl font-bold text-gray-800 ${className}`}>
  {title}
</h3>

