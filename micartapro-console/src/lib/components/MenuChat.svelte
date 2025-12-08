<script lang="ts">
  import { onMount } from 'svelte'
  import Message from './Message.svelte'
  import ChatInput from './ChatInput.svelte'

  interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: Date
  }

  let messages: ChatMessage[] = $state([])
  let isLoading = $state(false)
  let logoError = $state(false)

  const welcomeMessages = [
    {
      id: 'welcome-1',
      role: 'assistant' as const,
      content: '¡Hola! 👋 Soy tu asistente para crear menús digitales. Puedo ayudarte a armar tu carta de manera profesional.',
      timestamp: new Date()
    },
    {
      id: 'welcome-2',
      role: 'assistant' as const,
      content: 'Simplemente escribe tu menú y precios, y yo me encargaré de organizarlo y formatearlo para crear una carta atractiva.',
      timestamp: new Date()
    }
  ]

  onMount(() => {
    messages = [...welcomeMessages]
  })

  function handleSendMessage(content: string) {
    if (!content.trim() || isLoading) return

    // Agregar mensaje del usuario
    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: content.trim(),
      timestamp: new Date()
    }
    messages = [...messages, userMessage]

    // Simular procesamiento
    isLoading = true
    
    // Simular respuesta del asistente (aquí iría la lógica real de procesamiento)
    setTimeout(() => {
      const assistantMessage: ChatMessage = {
        id: `assistant-${Date.now()}`,
        role: 'assistant',
        content: generateResponse(content),
        timestamp: new Date()
      }
      messages = [...messages, assistantMessage]
      isLoading = false
    }, 1000)
  }

  function generateResponse(userInput: string): string {
    // Lógica temporal - aquí iría la integración con el backend
    return `He recibido tu menú. Estoy procesando la información para crear tu carta digital. Pronto tendrás una vista previa de cómo se verá.`
  }

  function scrollToBottom() {
    // Scroll automático al final
    setTimeout(() => {
      const container = document.getElementById('messages-container')
      if (container) {
        container.scrollTop = container.scrollHeight
      }
    }, 100)
  }

  $effect(() => {
    if (messages.length > 0) {
      scrollToBottom()
    }
  })
</script>

<div class="flex flex-col h-screen bg-white">
  <!-- Header -->
  <header class="border-b border-gray-200 bg-white px-2 flex items-center justify-between" style="padding-top: 0px; padding-bottom: 0px; min-height: auto;">
    <div class="flex items-center gap-1">
      <button class="p-1 hover:bg-gray-100 rounded-full transition-colors" aria-label="Menú">
        <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      {#if !logoError}
        <img 
          src="/logo.png" 
          alt="MiCartaPro Logo" 
          class="h-16 md:h-20 w-auto"
          style="margin: 0;"
          onerror={() => logoError = true}
        />
      {:else}
        <h1 class="text-2xl font-semibold text-gray-900" style="margin: 0; line-height: 1;">MiCartaPro</h1>
      {/if}
    </div>
    <div class="flex items-center gap-1">
      <button class="p-1 hover:bg-gray-100 rounded-full transition-colors" aria-label="Perfil">
        <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
      </button>
      <button class="p-1 hover:bg-gray-100 rounded-full transition-colors" aria-label="Nueva conversación">
        <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>
  </header>

  <!-- Messages Container -->
  <div 
    id="messages-container"
    class="flex-1 overflow-y-auto px-4 py-6 space-y-6"
  >
    {#if messages.length === 0}
      <div class="flex flex-col items-center justify-center h-full text-center px-4">
        <div class="mb-6">
          <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h2 class="text-2xl font-bold text-gray-900 mb-2">
            ¿En qué puedo ayudarte?
          </h2>
          <p class="text-gray-600 mb-8">
            Escribe tu menú y precios, y yo armaré tu carta digital
          </p>
        </div>

        <!-- Quick action buttons -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 w-full max-w-2xl">
          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage('Quiero crear un menú para un restaurante')}
          >
            <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-green-200 transition-colors">
              <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">Crear menú</p>
          </button>

          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage('Necesito ayuda para organizar mis platos')}
          >
            <div class="w-10 h-10 bg-orange-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-orange-200 transition-colors">
              <svg class="w-5 h-5 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">Organizar platos</p>
          </button>

          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage('¿Cómo funciona el sistema de precios?')}
          >
            <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-blue-200 transition-colors">
              <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">Ver precios</p>
          </button>

          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage('Más opciones')}
          >
            <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-purple-200 transition-colors">
              <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">Más</p>
          </button>
        </div>
      </div>
    {:else}
      <div class="max-w-3xl mx-auto space-y-6">
        {#each messages as message (message.id)}
          <Message message={message} />
        {/each}

        {#if isLoading}
          <div class="flex items-start gap-3">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div class="flex-1 bg-gray-50 rounded-2xl p-4">
              <div class="flex gap-1">
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0ms"></div>
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 150ms"></div>
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 300ms"></div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Chat Input -->
  <div class="border-t border-gray-200 bg-white px-4 py-3">
    <div class="max-w-3xl mx-auto">
      <ChatInput onSend={handleSendMessage} disabled={isLoading} />
    </div>
  </div>
</div>

<style>
  #messages-container {
    scroll-behavior: smooth;
  }
</style>

