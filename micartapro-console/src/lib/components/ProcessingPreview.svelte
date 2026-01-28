<script lang="ts">
  import { onMount } from 'svelte'
  import { t as tStore } from '../useLanguage'

  interface ProcessingPreviewProps {
    isVisible?: boolean
    blocking?: boolean // Si es true, muestra como modal bloqueante (fullscreen overlay)
    userInstruction?: string // Instrucción que envió el usuario (se muestra encima de la animación)
  }

  let { isVisible = true, blocking = false, userInstruction = '' }: ProcessingPreviewProps = $props()

  // Estados de progreso - se actualizan reactivamente con las traducciones
  const getProgressStates = () => {
    // Si está en modo bloqueante, mostrar solo el mensaje de preparación inicial
    if (blocking) {
      return [
        { 
          message: $tStore.processing.preparingInitialSetup,
          icon: '⚙️',
          duration: 2000
        }
      ]
    }
    
    // Estados normales para procesamiento de instrucciones
    return [
      { 
        message: $tStore.processing.understandingInstructions,
        icon: '💭',
        duration: 2000
      },
      { 
        message: $tStore.processing.creatingCatalog,
        icon: '📱',
        duration: 2500
      },
      { 
        message: $tStore.processing.preparingSite,
        icon: '⚙️',
        duration: 2000
      },
      { 
        message: $tStore.processing.validatingImages,
        icon: '🖼️',
        duration: 2000
      },
      { 
        message: $tStore.processing.improvingImages,
        icon: '✨',
        duration: 2000
      },
      { 
        message: $tStore.processing.finalizing,
        icon: '🎨',
        duration: 2000
      }
    ]
  }

  let currentStateIndex = $state(0)
  let currentMessage = $state(getProgressStates()[0].message)
  let currentIcon = $state(getProgressStates()[0].icon)
  let cycleInterval: ReturnType<typeof setTimeout> | null = null

  onMount(() => {
    if (!isVisible) return

    let stateIndex = 0
    const cycleStates = () => {
      if (!isVisible) {
        if (cycleInterval) {
          clearTimeout(cycleInterval)
          cycleInterval = null
        }
        return
      }
      
      const states = getProgressStates()
      stateIndex = (stateIndex + 1) % states.length
      currentStateIndex = stateIndex
      currentMessage = states[stateIndex].message
      currentIcon = states[stateIndex].icon
      
      cycleInterval = setTimeout(cycleStates, states[stateIndex].duration)
    }

    // Iniciar el ciclo después del primer estado
    const states = getProgressStates()
    cycleInterval = setTimeout(cycleStates, states[0].duration)

    return () => {
      if (cycleInterval) {
        clearTimeout(cycleInterval)
        cycleInterval = null
      }
    }
  })

  $effect(() => {
    if (isVisible) {
      // Resetear al primer estado cuando se vuelve a mostrar
      const states = getProgressStates()
      currentStateIndex = 0
      currentMessage = states[0].message
      currentIcon = states[0].icon
    } else {
      // Limpiar el intervalo cuando se oculta
      if (cycleInterval) {
        clearTimeout(cycleInterval)
        cycleInterval = null
      }
    }
  })
</script>

{#if isVisible}
  {#if blocking}
    <!-- Modal bloqueante: animación limpia sin celular -->
    <div class="fixed inset-0 bg-white z-50 flex items-center justify-center animate-fade-in">
      <div class="text-center px-6 max-w-md">
        {#if userInstruction}
          <p class="text-gray-600 text-sm mb-6 px-4 py-3 bg-gray-50 rounded-xl border border-gray-100 text-left max-h-20 overflow-y-auto">
            "{userInstruction}"
          </p>
        {/if}
        <div class="text-5xl mb-6 animate-pulse">{currentIcon}</div>
        <div class="processing-dots mb-5">
          <span></span><span></span><span></span>
        </div>
        <p class="text-lg font-medium text-gray-900 mb-2">{currentMessage}</p>
        <p class="text-sm text-gray-500">{$tStore.processing.pleaseWait}</p>
      </div>
    </div>
  {:else}
    <!-- Embebido en el chat: animación limpia sin celular -->
    <div class="flex flex-col items-center justify-center h-full w-full animate-fade-in py-8">
      <div class="text-center px-4 max-w-sm w-full">
        {#if userInstruction}
          <p class="text-gray-600 text-sm mb-4 px-3 py-2.5 bg-gray-50 rounded-lg border border-gray-100 text-left max-h-16 overflow-y-auto">
            "{userInstruction}"
          </p>
        {/if}
        <div class="processing-dots mb-4">
          <span></span><span></span><span></span>
        </div>
        <div class="flex items-center justify-center gap-2 mb-1">
          <div class="animate-spin rounded-full h-5 w-5 border-2 border-blue-500 border-t-transparent"></div>
          <p class="text-base font-medium text-gray-900">{currentMessage}</p>
        </div>
        <p class="text-xs text-gray-500">{$tStore.processing.pleaseWait}</p>
      </div>
    </div>
  {/if}
{/if}

<style>
  @keyframes fade-in {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes dot-pulse {
    0%, 80%, 100% {
      opacity: 0.4;
      transform: scale(0.9);
    }
    40% {
      opacity: 1;
      transform: scale(1.1);
    }
  }

  .animate-fade-in {
    animation: fade-in 0.3s ease-out;
  }

  .processing-dots {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .processing-dots span {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: linear-gradient(135deg, rgb(59 130 246), rgb(147 51 234));
    animation: dot-pulse 1.4s ease-in-out infinite both;
  }

  .processing-dots span:nth-child(1) { animation-delay: 0s; }
  .processing-dots span:nth-child(2) { animation-delay: 0.2s; }
  .processing-dots span:nth-child(3) { animation-delay: 0.4s; }
</style>
