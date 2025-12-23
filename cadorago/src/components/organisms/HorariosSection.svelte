<script>
  /** @type {string[]} - Array de horarios de atención */
  export let horarios = [];
  
  /** @type {string} - Additional custom classes */
  export let className = '';
</script>

<section class={`${className}`}>
  <ul class="space-y-2 sm:space-y-3 list-none">
    {#each horarios as horario}
      <li class="text-lg sm:text-xl lg:text-2xl text-gray-700 leading-relaxed flex items-start">
        <span class="mr-2 sm:mr-3 flex-shrink-0">•</span>
        <span>{horario}</span>
      </li>
    {/each}
  </ul>
</section>

