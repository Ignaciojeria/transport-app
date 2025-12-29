<script lang="ts">
  import { onMount } from 'svelte'
  import { v7 as uuidv7 } from 'uuid'
  import Message from './Message.svelte'
  import ChatInput from './ChatInput.svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, generateMenuUrl, pollUntilMenuUpdated, pollUntilMenuExists } from '../menuUtils'
  import { API_BASE_URL } from '../config'
  import { t as tStore, language } from '../useLanguage'

  interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: Date
  }

  let messages: ChatMessage[] = $state([])
  let isLoading = $state(false)
  let logoError = $state(false)
  let showPreview = $state(false)
  let menuUrl = $state<string | null>(null)
  let menuId = $state<string | null>(null)
  let copySuccess = $state(false)
  let menuReady = $state(false)
  let checkingMenu = $state(false)

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const session = $derived(authState.session)
  const currentLanguage = $derived($language)

  // Función para obtener mensajes de bienvenida
  function getWelcomeMessages() {
    return [
      {
        id: 'welcome-1',
        role: 'assistant' as const,
        content: $tStore.chat.welcomeMessage1,
        timestamp: new Date()
      },
      {
        id: 'welcome-2',
        role: 'assistant' as const,
        content: $tStore.chat.welcomeMessage2,
        timestamp: new Date()
      }
    ]
  }

  onMount(async () => {
    messages = getWelcomeMessages()
    
    // Actualizar menuUrl cuando cambie el idioma
    $effect(() => {
      if (menuUrl && userId && menuId) {
        // Reconstruir la URL con el idioma actual
        menuUrl = generateMenuUrl(userId, menuId, currentLanguage)
      }
    })
    
    // Cargar menuID al montar el componente
    if (userId) {
      try {
        const id = await getLatestMenuId(userId)
        if (id) {
          menuId = id
          
          // Hacer polling para validar que el menú exista en GCS
          checkingMenu = true
          const exists = await pollUntilMenuExists(userId, id)
          
          if (exists) {
            menuReady = true
          } else {
            // Si no se encontró después de todos los intentos, mostrar mensaje
            const errorMessage: ChatMessage = {
              id: `menu-not-found-${Date.now()}`,
              role: 'assistant',
              content: $tStore.chat.errorNoMenu,
              timestamp: new Date()
            }
            messages = [...messages, errorMessage]
          }
          
          checkingMenu = false
        } else {
          // No hay menuId en la base de datos
          menuReady = false
        }
      } catch (err) {
        console.error('Error cargando menuID:', err)
        checkingMenu = false
        menuReady = false
      }
    }
  })

  async function togglePreview() {
    if (showPreview) {
      showPreview = false
      return
    }

    if (!menuUrl && userId) {
      try {
        const menuId = await getLatestMenuId(userId)
        if (menuId) {
          menuUrl = generateMenuUrl(userId, menuId, currentLanguage)
        }
      } catch (err) {
        console.error('Error cargando menú:', err)
      }
    }
    
    showPreview = true
  }

  function copyToClipboard() {
    if (!menuUrl) return

    navigator.clipboard.writeText(menuUrl).then(() => {
      copySuccess = true
      setTimeout(() => {
        copySuccess = false
      }, 2000)
    }).catch((err) => {
      console.error('Error copiando al portapapeles:', err)
    })
  }

  async function handleSendMessage(content: string) {
    if (!content.trim() || isLoading) return

    // Validar que tenemos los datos necesarios
      if (!session?.access_token) {
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: $tStore.chat.errorNoSession,
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
      return
    }

    if (!menuId) {
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: $tStore.chat.errorNoMenu,
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
      return
    }

    // Agregar mensaje del usuario
    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: content.trim(),
      timestamp: new Date()
    }
    messages = [...messages, userMessage]

    // Mostrar indicador de carga
    isLoading = true

    try {
      // Generar idempotency key (UUID v7)
      const idempotencyKey = uuidv7()

      // Enviar POST al backend
      const response = await fetch(`${API_BASE_URL}/menu/interaction`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Idempotency-Key': idempotencyKey,
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
          menuId: menuId,
          message: content.trim()
        })
      })

      if (!response.ok) {
        throw new Error(`Error del servidor: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()

      // Agregar respuesta inicial del asistente
      const initialMessage: ChatMessage = {
        id: `assistant-${Date.now()}`,
        role: 'assistant',
        content: data.message || data.response || 'Procesando tu mensaje...',
        timestamp: new Date()
      }
      messages = [...messages, initialMessage]

      // Iniciar polling para esperar la actualización del menú
      if (userId && menuId) {
        try {
          const updatedMenu = await pollUntilMenuUpdated(userId, menuId, idempotencyKey)
          
          // Agregar mensaje de confirmación con el menú actualizado
          const successMessage: ChatMessage = {
            id: `success-${Date.now()}`,
            role: 'assistant',
            content: $tStore.chat.successUpdated,
            timestamp: new Date()
          }
          messages = [...messages, successMessage]
          
          // Abrir automáticamente la vista previa para mostrar los cambios
          if (!menuUrl) {
            menuUrl = generateMenuUrl(userId, menuId, currentLanguage)
          }
          showPreview = true
        } catch (pollError: any) {
          console.error('Error en polling:', pollError)
          
          // Agregar mensaje de error del polling
          const errorMessage: ChatMessage = {
            id: `poll-error-${Date.now()}`,
            role: 'assistant',
            content: $tStore.chat.errorPolling.replace('{message}', pollError.message || 'Error desconocido'),
            timestamp: new Date()
          }
          messages = [...messages, errorMessage]
        }
      }

    } catch (error: any) {
      console.error('Error enviando mensaje:', error)
      
      // Mostrar error al usuario
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: $tStore.chat.errorProcessing.replace('{message}', error.message || 'Error desconocido'),
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
    } finally {
      isLoading = false
    }
  }

  function scrollToBottom() {
    // Scroll automático al final
    setTimeout(() => {
      const container = document.getElementById('messages-container')
      if (container) {
        // Scroll suave al final, con un pequeño offset para asegurar visibilidad
        container.scrollTo({
          top: container.scrollHeight,
          behavior: 'smooth'
        })
      }
    }, 100)
  }
  
  // Función para hacer scroll cuando el teclado aparece en móviles
  function handleInputFocus() {
    // En móviles, hacer scroll para asegurar que el input sea visible
    if (window.innerWidth <= 768) {
      setTimeout(() => {
        scrollToBottom()
      }, 300) // Esperar a que el teclado aparezca
    }
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
      <button 
        onclick={togglePreview}
        class="p-2 hover:bg-gray-100 rounded-lg transition-colors flex items-center gap-2 {showPreview ? 'bg-blue-50 text-blue-600' : 'text-gray-600'}"
        aria-label="Ver vista previa"
        title="Ver vista previa de tu carta"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        <span class="text-xs font-medium hidden sm:inline">{$tStore.chat.previewButton}</span>
      </button>
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
    class="flex-1 overflow-y-auto px-4 py-6 pb-24 md:pb-6 space-y-6"
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
            {$tStore.chat.welcomeTitle}
          </h2>
          <p class="text-gray-600 mb-8">
            {$tStore.chat.welcomeSubtitle}
          </p>
        </div>

        <!-- Quick action buttons -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 w-full max-w-2xl">
          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage($tStore.chat.createMenuMessage)}
          >
            <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-green-200 transition-colors">
              <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">{$tStore.chat.createMenu}</p>
          </button>

          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage($tStore.chat.organizeDishesMessage)}
          >
            <div class="w-10 h-10 bg-orange-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-orange-200 transition-colors">
              <svg class="w-5 h-5 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">{$tStore.chat.organizeDishes}</p>
          </button>

          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage($tStore.chat.viewPricesMessage)}
          >
            <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-blue-200 transition-colors">
              <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">{$tStore.chat.viewPrices}</p>
          </button>

          <button 
            class="p-4 bg-gray-50 hover:bg-gray-100 rounded-xl transition-colors text-left group"
            onclick={() => handleSendMessage($tStore.chat.moreOptionsMessage)}
          >
            <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center mb-2 group-hover:bg-purple-200 transition-colors">
              <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-gray-900">{$tStore.chat.moreOptions}</p>
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
  <div class="border-t border-gray-200 bg-white px-4 py-3 sticky bottom-0 z-10 safe-area-inset-bottom">
    <div class="max-w-3xl mx-auto">
      {#if checkingMenu}
        <div class="flex items-center justify-center py-8">
          <div class="text-center">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 border-t-transparent mx-auto mb-4"></div>
            <p class="text-gray-600 text-sm">{$tStore.chat.checkingMenu}</p>
          </div>
        </div>
      {/if}
      
      <ChatInput onSend={handleSendMessage} disabled={isLoading || !menuReady || checkingMenu} onFocus={handleInputFocus} />
    </div>
  </div>
</div>

<!-- Modal de Vista Previa -->
{#if showPreview}
  <div 
    class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
    onclick={() => showPreview = false}
    role="dialog"
    aria-modal="true"
  >
    <div 
      class="bg-white rounded-lg shadow-xl w-full max-w-6xl h-[90vh] flex flex-col"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header del Modal -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200 gap-4">
        <h2 class="text-lg font-semibold text-gray-900 flex-1">{$tStore.chat.previewTitle}</h2>
        {#if menuUrl}
          <button
            onclick={copyToClipboard}
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center gap-2 whitespace-nowrap"
            title={$tStore.chat.copyLink}
          >
            {#if copySuccess}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-sm">{$tStore.chat.linkCopied}</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <span class="text-sm">{$tStore.chat.copyLink}</span>
            {/if}
          </button>
        {/if}
        <button
          onclick={() => showPreview = false}
          class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          aria-label="Cerrar"
        >
          <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Contenido del Modal -->
      <div class="flex-1 overflow-hidden iframe-container">
        {#if menuUrl}
          <iframe
            src={menuUrl}
            class="w-full h-full border-0"
            title="Vista previa de la carta"
            loading="lazy"
            allow="camera; microphone; geolocation; autoplay; clipboard-write"
            sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation"
          />
        {:else}
          <div class="flex items-center justify-center h-full">
            <div class="text-center">
              <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
              <p class="text-gray-600">{$tStore.chat.loadingPreview}</p>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  #messages-container {
    scroll-behavior: smooth;
  }
  
  /* Safe area para dispositivos con notch */
  .safe-area-inset-bottom {
    padding-bottom: env(safe-area-inset-bottom, 0.75rem);
  }
  
  /* Asegurar que el input sea visible en móviles */
  @media (max-width: 768px) {
    /* Aumentar padding inferior en móviles para que el input no oculte contenido */
    #messages-container {
      padding-bottom: 6rem;
    }
  }
  
  /* Fix para scroll en iframes en mobile */
  .iframe-container {
    touch-action: pan-y pan-x;
    -webkit-overflow-scrolling: touch;
    position: relative;
  }
  
  .iframe-container iframe {
    touch-action: auto;
    pointer-events: auto;
  }
</style>

