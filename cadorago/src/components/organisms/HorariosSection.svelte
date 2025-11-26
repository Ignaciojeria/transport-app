<script>
  /** @type {string[]} - Array de horarios de atención */
  export let horarios = [];
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<section class={`${className}`}>
  <p class="text-lg sm:text-xl lg:text-2xl text-gray-700 leading-relaxed">
    {#each horarios as horario, index}
      {horario}{#if index < horarios.length - 1}. {/if}
    {/each}
  </p>
</section>

