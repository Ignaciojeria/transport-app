<script>
  import { t } from '../../lib/useLanguage';
  
  const { horarios = [], className = '' } = $props();
  
  // Filtrar horarios que tengan contenido (no vacíos)
  const horariosConContenido = $derived(
    horarios?.filter(h => h && typeof h === 'string' && h.trim().length > 0) || []
  );
  
  const tieneHorarios = $derived(horariosConContenido.length > 0);
</script>

<section class={`${className}`}>
  {#if tieneHorarios}
    <ul class="space-y-2 sm:space-y-3 list-none">
      {#each horariosConContenido as horario}
        <li class="text-lg sm:text-xl lg:text-2xl text-gray-700 leading-relaxed flex items-start">
          <span class="mr-2 sm:mr-3 flex-shrink-0">•</span>
          <span>{horario}</span>
        </li>
      {/each}
    </ul>
  {:else}
    <p class="text-lg sm:text-xl lg:text-2xl text-gray-500 italic">
      {$t.hours.notSpecified}
    </p>
  {/if}
</section>

