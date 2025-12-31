<script lang="ts">
  import { onMount } from 'svelte'
  import confetti from 'canvas-confetti'
  import { t as tStore } from '../useLanguage'

  let countdown = 5
  let redirecting = false
  let t = $derived($tStore)

  onMount(() => {
    // Confeti moderado 🎉
    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8', '#F7DC6F', '#BB8FCE']
    
    // Explosión inicial desde arriba
    setTimeout(() => {
      confetti({
        particleCount: 80,
        angle: 90,
        spread: 60,
        origin: { x: 0.5, y: 0 },
        colors: colors
      })
    }, 100)

    // Una explosión más pequeña desde el centro después de medio segundo
    setTimeout(() => {
      confetti({
        particleCount: 50,
        spread: 70,
        origin: { x: 0.5, y: 0.5 },
        colors: colors
      })
    }, 600)

    // Countdown para redirección
    const countdownInterval = setInterval(() => {
      countdown--
      if (countdown <= 0) {
        clearInterval(countdownInterval)
        redirecting = true
        // Limpiar parámetros de la URL y redirigir
        const url = new URL(window.location.href)
        url.searchParams.delete('payment')
        url.searchParams.delete('success')
        url.hash = ''
        window.location.href = url.pathname
      }
    }, 1000)

    return () => {
      clearInterval(countdownInterval)
    }
  })

  function handleGoToConsole() {
    const url = new URL(window.location.href)
    url.searchParams.delete('payment')
    url.searchParams.delete('success')
    url.hash = ''
    window.location.href = url.pathname
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-green-50 via-blue-50 to-purple-50 flex items-center justify-center p-6">
  <div class="bg-white/90 backdrop-blur-sm rounded-3xl shadow-2xl border border-white/20 p-12 max-w-2xl w-full text-center relative overflow-hidden">
    <!-- Efecto de brillo animado -->
    <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent animate-shimmer"></div>
    
    <!-- Ícono de éxito grande -->
    <div class="relative z-10 mb-8">
      <div class="w-32 h-32 bg-gradient-to-br from-green-400 to-emerald-600 rounded-full flex items-center justify-center mx-auto shadow-lg animate-bounce-slow">
        <svg class="w-20 h-20 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
        </svg>
      </div>
    </div>

    <!-- Título -->
    <h1 class="text-5xl font-bold text-gray-900 mb-4 relative z-10 animate-fade-in">
      {t.paymentSuccess.title}
    </h1>

    <!-- Subtítulo -->
    <h2 class="text-2xl font-semibold text-gray-700 mb-6 relative z-10">
      {t.paymentSuccess.subtitle}
    </h2>

    <!-- Mensaje -->
    <p class="text-lg text-gray-600 mb-8 relative z-10 leading-relaxed">
      {t.paymentSuccess.message}
    </p>

    <!-- Botón -->
    <button
      on:click={handleGoToConsole}
      class="relative z-10 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold py-4 px-8 rounded-xl shadow-lg transform transition-all duration-200 hover:scale-105 hover:shadow-xl mb-6"
    >
      {t.paymentSuccess.goToConsole}
    </button>

    <!-- Countdown -->
    <p class="text-sm text-gray-500 relative z-10">
      {t.paymentSuccess.redirecting} {countdown > 0 ? `(${countdown}s)` : ''}
    </p>
  </div>
</div>

<style>
  @keyframes bounce-slow {
    0%, 100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  @keyframes shimmer {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  @keyframes fade-in {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-bounce-slow {
    animation: bounce-slow 2s ease-in-out infinite;
  }

  .animate-shimmer {
    animation: shimmer 3s ease-in-out infinite;
  }

  .animate-fade-in {
    animation: fade-in 0.6s ease-out;
  }
</style>

