<script lang="ts">
  import { onMount } from 'svelte'
  import { v7 as uuidv7 } from 'uuid'
  import Message from './Message.svelte'
  import ChatInput from './ChatInput.svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, generateMenuUrl, pollUntilMenuUpdated, pollUntilMenuExists, fetchEntitlement, calculateTrialDaysRemaining, getMenuSlug, createMenuSlug, generateSlugUrl, type Entitlement } from '../menuUtils'
  import { API_BASE_URL } from '../config'
  import { t as tStore, language } from '../useLanguage'

  interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: Date
    showExploreButton?: boolean
  }

  let messages: ChatMessage[] = $state([])
  let isLoading = $state(false)
  let logoError = $state(false)
  let showPreview = $state(false)
  let menuUrl = $state<string | null>(null)
  let menuId = $state<string | null>(null)
  let copySuccess = $state(false)
  let linkWasCopied = $state(false) // Estado para saber si el enlace fue copiado (persistente)
  let menuReady = $state(false)
  let checkingMenu = $state(false)
  let chatInputRef: any = $state(null)
  let showExamples = $state(false)
  let currentExampleType: 'address' | 'dishes' | 'desserts' | 'price' | 'delete' | 'whatsapp' | null = $state(null)
  let messageSent = $state(false) // Flag para indicar que se envió un mensaje
  let showUpgradeModal = $state(false) // Modal de upgrade al copiar link
  let entitlement = $state<Entitlement | null>(null) // Entitlement del usuario
  let iframeKey = $state(0) // Key para forzar recarga del iframe cuando el menú se actualiza
  let showSlugModal = $state(false) // Modal para crear slug
  let businessName = $state('') // Nombre del negocio
  let slugPreview = $state('') // Preview del slug
  let isCreatingSlug = $state(false) // Estado de carga al crear slug
  let showConfetti = $state(false) // Mostrar confeti
  let confettiParticles = $state<Array<{left: number, delay: number, color: string}>>([]) // Partículas de confeti

  const user = $derived(authState.user)
  const userId = $derived(user?.id || '')
  const userName = $derived(
    user?.user_metadata?.name || 
    user?.user_metadata?.full_name || 
    user?.email?.split('@')[0] || 
    'Usuario'
  )
  const userPicture = $derived(user?.user_metadata?.picture || user?.user_metadata?.avatar_url || null)
  const session = $derived(authState.session)
  const currentLanguage = $derived($language)
  
  // Días restantes del trial calculados desde el entitlement
  const trialDaysRemaining = $derived(
    entitlement && entitlement.status === 'trialing' 
      ? calculateTrialDaysRemaining(entitlement.ends_at)
      : null
  )

  async function handleUpgradeToPro() {
    // Cerrar el modal antes de redirigir
    showUpgradeModal = false
    
    // Validar que el token esté disponible
    if (!session?.access_token) {
      console.error('Error: No hay token de autenticación disponible')
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: $tStore.chat.errorNoSession,
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
      return
    }
    
    try {
      const checkoutUrl = `${API_BASE_URL}/checkout`
      
      // Hacer fetch con el token de Supabase para obtener la URL de checkout
      const response = await fetch(checkoutUrl, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        },
        credentials: 'include',
      })

      if (!response.ok) {
        const errorText = await response.text()
        console.error('Error del servidor:', response.status, errorText)
        const errorMessage: ChatMessage = {
          id: `error-${Date.now()}`,
          role: 'assistant',
          content: `Error al obtener la URL de checkout: ${response.status} ${response.statusText}`,
          timestamp: new Date()
        }
        messages = [...messages, errorMessage]
        return
      }

      // El backend devuelve un JSON con checkout_url
      const data = await response.json()
      if (data.checkout_url) {
        window.open(data.checkout_url, '_blank')
      } else {
        console.error('No se recibió checkout_url en la respuesta:', data)
        const errorMessage: ChatMessage = {
          id: `error-${Date.now()}`,
          role: 'assistant',
          content: 'Error: No se recibió la URL de checkout del servidor',
          timestamp: new Date()
        }
        messages = [...messages, errorMessage]
      }
    } catch (error) {
      console.error('Error al abrir checkout:', error)
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: $tStore.chat.errorProcessing.replace('{message}', error instanceof Error ? error.message : 'Error desconocido'),
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
    }
  }

  onMount(async () => {
    // No mostrar mensajes inicialmente, solo los botones
    messages = []
    
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

    if (!menuUrl && userId && menuId) {
      try {
        menuUrl = generateMenuUrl(userId, menuId, currentLanguage)
      } catch (err) {
        console.error('Error cargando menú:', err)
      }
    }
    
    // Cargar entitlement cuando se abre el modal
    if (userId && !entitlement) {
      try {
        entitlement = await fetchEntitlement(userId)
      } catch (err) {
        console.error('Error cargando entitlement:', err)
      }
    }
    
    showPreview = true
  }

  async function copyToClipboard() {
    if (!menuUrl) return

    // Siempre recargar entitlement para obtener la versión más reciente
    if (userId) {
      try {
        entitlement = await fetchEntitlement(userId)
      } catch (err) {
        console.error('Error cargando entitlement:', err)
      }
    }

    try {
      await navigator.clipboard.writeText(menuUrl)
      copySuccess = true
      linkWasCopied = true // Marcar que el enlace fue copiado
      setTimeout(() => {
        copySuccess = false
      }, 2000)
      
      // Mostrar modal de upgrade SOLO si el usuario está en trial (status === 'trialing')
      // NO mostrar si el status es 'active' o si tiene access=true (premium activo)
      console.log('Entitlement al copiar:', entitlement)
      console.log('Status:', entitlement?.status, 'Access:', entitlement?.access)
      
      // Solo mostrar si está explícitamente en trial
      const isInTrial = entitlement && entitlement.status === 'trialing'
      const hasPremium = entitlement && (entitlement.status === 'active' || entitlement.access === true)
      
      if (isInTrial && !hasPremium) {
        console.log('Mostrando modal de upgrade - usuario en trial')
        showUpgradeModal = true
      } else {
        // Usuario tiene premium (active) o no hay entitlement, no mostrar modal
        console.log('NO mostrando modal - isInTrial:', isInTrial, 'hasPremium:', hasPremium)
        showUpgradeModal = false
      }
    } catch (err) {
      console.error('Error copiando al portapapeles:', err)
    }
  }

  // Generar slug desde el nombre del negocio
  function generateSlugFromName(name: string): string {
    return name
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '') // Eliminar acentos
      .replace(/[^a-z0-9\s-]/g, '') // Eliminar caracteres especiales
      .trim()
      .replace(/\s+/g, '-') // Reemplazar espacios con guiones
      .replace(/-+/g, '-') // Reemplazar múltiples guiones con uno solo
      .replace(/^-|-$/g, '') // Eliminar guiones al inicio y final
  }

  // Actualizar preview del slug cuando cambia el nombre
  $effect(() => {
    if (businessName.trim()) {
      slugPreview = generateSlugFromName(businessName)
    } else {
      slugPreview = ''
    }
  })

  async function shareOnWhatsApp() {
    if (!menuId || !session?.access_token) {
      alert('No se puede compartir: falta información del menú')
      return
    }

    try {
      // Verificar si existe un slug para este menú
      const existingSlug = await getMenuSlug(menuId, session.access_token)
      
      if (!existingSlug) {
        // No existe slug, mostrar modal para crear uno
        businessName = ''
        slugPreview = ''
        showSlugModal = true
        return
      }

      // Si ya existe slug, proceder directamente a compartir
      await proceedWithShare(existingSlug)
    } catch (err) {
      console.error('Error compartiendo:', err)
      alert(
        currentLanguage === 'ES'
          ? 'Error al compartir el menú. Intenta nuevamente.'
          : currentLanguage === 'PT'
          ? 'Erro ao compartilhar o cardápio. Tente novamente.'
          : 'Error sharing menu. Please try again.'
      )
    }
  }

  async function createSlugAndShare() {
    if (!menuId || !session?.access_token) {
      return
    }

    if (!businessName.trim()) {
      const lang = currentLanguage
      alert(
        lang === 'ES'
          ? 'Por favor ingresa el nombre de tu negocio'
          : lang === 'PT'
          ? 'Por favor, insira o nome do seu negócio'
          : 'Please enter your business name'
      )
      return
    }

    if (!slugPreview) {
      const lang = currentLanguage
      alert(
        lang === 'ES'
          ? 'El nombre del negocio no genera un slug válido'
          : lang === 'PT'
          ? 'O nome do negócio não gera um slug válido'
          : 'Business name does not generate a valid slug'
      )
      return
    }

    isCreatingSlug = true

    try {
      const slug = await createMenuSlug(menuId, slugPreview, session.access_token)
      
      if (!slug) {
        const lang = currentLanguage
        alert(
          lang === 'ES'
            ? 'Error al crear el slug. Intenta nuevamente.'
            : lang === 'PT'
            ? 'Erro ao criar o slug. Tente novamente.'
            : 'Error creating slug. Please try again.'
        )
        isCreatingSlug = false
        return
      }

      // Generar partículas de confeti
      confettiParticles = Array.from({ length: 50 }, () => ({
        left: Math.random() * 100,
        delay: Math.random() * 0.5,
        color: ['#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8', '#F7DC6F', '#BB8FCE'][Math.floor(Math.random() * 7)]
      }))
      
      // Mostrar confeti
      showConfetti = true
      setTimeout(() => {
        showConfetti = false
        confettiParticles = []
      }, 3000)

      // Cerrar modal después de un breve delay
      setTimeout(() => {
        showSlugModal = false
        businessName = ''
        slugPreview = ''
        isCreatingSlug = false
      }, 1500)

      // Proceder a compartir
      await proceedWithShare(slug)
    } catch (error: any) {
      isCreatingSlug = false
      const lang = currentLanguage
      if (error.message === 'SLUG_EXISTS') {
        alert(
          lang === 'ES'
            ? 'Este slug ya está en uso por otro menú. Por favor, elige otro nombre.'
            : lang === 'PT'
            ? 'Este slug já está em uso por outro cardápio. Por favor, escolha outro nome.'
            : 'This slug is already in use by another menu. Please choose another name.'
        )
      } else {
        alert(
          lang === 'ES'
            ? 'Error al crear el slug. Intenta nuevamente.'
            : lang === 'PT'
            ? 'Erro ao criar o slug. Tente novamente.'
            : 'Error creating slug. Please try again.'
        )
      }
    }
  }

  async function proceedWithShare(slug: string) {
    const lang = currentLanguage

    // Generar URL con el slug
    const shareUrl = generateSlugUrl(slug, currentLanguage)

    // Copiar al portapapeles
    await navigator.clipboard.writeText(shareUrl)
    copySuccess = true
    linkWasCopied = true
    setTimeout(() => {
      copySuccess = false
    }, 2000)
    
    // Cargar entitlement para verificar si mostrar modal de upgrade
    if (userId && !entitlement) {
      try {
        entitlement = await fetchEntitlement(userId)
      } catch (err) {
        console.error('Error cargando entitlement:', err)
      }
    }
    
    // Mostrar modal de upgrade SOLO si el usuario está en trial
    const isInTrial = entitlement && entitlement.status === 'trialing'
    const hasPremium = entitlement && (entitlement.status === 'active' || entitlement.access === true)
    
    if (isInTrial && !hasPremium) {
      showUpgradeModal = true
    }

    // Mensaje por defecto según el idioma
    let message = ''
    if (lang === 'ES') {
      message = `¡Mira mi carta digital! ${shareUrl}`
    } else if (lang === 'PT') {
      message = `Confira meu cardápio digital! ${shareUrl}`
    } else {
      message = `Check out my digital menu! ${shareUrl}`
    }

    const encodedMessage = encodeURIComponent(message)
    const whatsappUrl = `https://wa.me/?text=${encodedMessage}`
    window.open(whatsappUrl, '_blank')
  }

  async function handleSendMessage(content: string) {
    if (!content.trim() || isLoading) return

    // Ocultar sugerencias cuando se envía un mensaje
    showExamples = false
    currentExampleType = null
    messageSent = true // Marcar que se envió un mensaje

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
            timestamp: new Date(),
            showExploreButton: true // Marcar que este mensaje debe mostrar el botón
          }
          messages = [...messages, successMessage]
          
          // Actualizar la URL del menú y forzar recarga del iframe
          menuUrl = generateMenuUrl(userId, menuId, currentLanguage)
          // Incrementar iframeKey para forzar recarga del iframe
          iframeKey++
          
          // Abrir automáticamente la vista previa para mostrar los cambios
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
    // Si hay ejemplos mostrándose, mantenerlos visibles cuando se hace focus
    if (showExamples && currentExampleType) {
      // Asegurar que los ejemplos sigan visibles
      setTimeout(() => {
        scrollToBottom()
      }, 100)
    }
    
    // En móviles, hacer scroll para asegurar que el input sea visible
    if (window.innerWidth <= 768) {
      setTimeout(() => {
        scrollToBottom()
      }, 300) // Esperar a que el teclado aparezca
    }
  }

  function handleButtonClick(type: 'address' | 'dishes' | 'desserts' | 'price' | 'delete' | 'whatsapp') {
    showExamples = true
    currentExampleType = type
    messages = []
    messageSent = false // Asegurar que el flag esté en false
    
    // En móvil, esperar un poco más para que el teclado se abra correctamente
    setTimeout(() => {
      chatInputRef?.focus()
      // Forzar focus nuevamente después de un pequeño delay para móviles
      if (window.innerWidth <= 768) {
        setTimeout(() => {
          chatInputRef?.focus()
        }, 300)
      }
    }, 100)
  }

  function getExampleMessages(): string[] {
    switch (currentExampleType) {
      case 'address':
        return [
          'Cambia la dirección a Avenida Siempre Viva 3151',
          'Actualiza la dirección del restaurante a Calle Principal 123'
        ]
      case 'dishes':
        return [
          'En la categoría "Platos Principales" añade empanadas de pollo a 3000 pesos',
          'Agrega pasta carbonara a 8500 pesos en la sección de platos principales'
        ]
      case 'desserts':
        return [
          'En la categoría "Postres" añade torta de chocolate a 4500 pesos',
          'Agrega helado de vainilla a 2500 pesos en postres'
        ]
      case 'price':
        return [
          'Cambia el precio de las empanadas a 3500 pesos',
          'Actualiza el precio del plato principal "Pasta Carbonara" a 9000 pesos'
        ]
      case 'delete':
        return [
          'Elimina las empanadas de pollo del menú',
          'Quita el plato "Pasta Carbonara" de la carta'
        ]
      case 'whatsapp':
        const lang = currentLanguage
        if (lang === 'ES') {
          return [
            'Cambia el número de WhatsApp a +56912345678'
          ]
        } else if (lang === 'PT') {
          return [
            'Mude o número do WhatsApp para +5511987654321'
          ]
        } else {
          return [
            'Change the WhatsApp number to +15551234567'
          ]
        }
      default:
        return []
    }
  }

  function handleInputBlur() {
    // Si se acaba de enviar un mensaje, no resetear
    if (messageSent) {
      messageSent = false // Resetear el flag
      return
    }
    
    // Si hay sugerencias mostrándose, mantenerlas visibles (no desaparecer al perder focus)
    // Solo desaparecerán cuando se envíe el prompt
    if (showExamples && currentExampleType) {
      return
    }
    
    // NO resetear si hay mensajes en el chat (el usuario ya está en una conversación)
    if (messages.length > 0) {
      return
    }
    
    // En móvil, esperar más tiempo para detectar si el teclado se cerró realmente
    const delay = window.innerWidth <= 768 ? 500 : 200
    
    setTimeout(() => {
      const textarea = chatInputRef?.textareaRef
      // Verificar si el textarea todavía tiene focus (en algunos móviles el blur puede ser temporal)
      const isStillFocused = document.activeElement === textarea
      
      if (!isStillFocused && textarea && !textarea.value?.trim()) {
        // Solo resetear si realmente perdió el focus y no hay texto
        // Y solo si no hay sugerencias mostrándose
        if (!showExamples) {
          showExamples = false
          currentExampleType = null
          messages = []
        }
      }
    }, delay)
  }

  function resetToOptions() {
    // Resetear para mostrar los botones de opciones nuevamente
    showExamples = false
    currentExampleType = null
    messages = []
  }

  $effect(() => {
    if (messages.length > 0) {
      scrollToBottom()
    }
  })
</script>

<div class="flex flex-col h-screen bg-white relative overflow-hidden">
  <!-- Vista de Chat (oculta cuando showPreview es true) -->
  <div 
    class="flex flex-col h-full transition-transform duration-300 ease-in-out {showPreview ? '-translate-x-full' : 'translate-x-0'}"
  >
    <!-- Header estilo Gemini -->
    <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between">
      <button class="p-2 hover:bg-gray-100 rounded-full transition-colors" aria-label="Menú">
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <h1 class="text-lg font-medium text-gray-900">MiCartaPro</h1>
      <div class="flex items-center gap-2">
        <button 
          onclick={togglePreview}
          class="p-2 hover:bg-gray-100 rounded-full transition-colors text-gray-600"
          aria-label="Ver vista previa"
          title="Ver vista previa de tu carta"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        </button>
        {#if userPicture}
          <img 
            src={userPicture} 
            alt={userName}
            class="w-8 h-8 rounded-full"
          />
        {:else}
          <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-medium text-sm">
            {userName.charAt(0).toUpperCase()}
          </div>
        {/if}
      </div>
    </header>

    <!-- Messages Container -->
    <div 
      id="messages-container"
      class="flex-1 overflow-y-auto px-4 py-6 pb-24 md:pb-6"
    >
    {#if messages.length === 0}
      <div class="flex flex-col h-full px-4 max-w-2xl mx-auto">
        <!-- Saludo personalizado estilo Gemini -->
        <div class="mt-8 mb-6 text-center">
          <h2 class="text-2xl font-normal text-gray-900 mb-2">
            {$tStore.chat.greeting}, {userName}
          </h2>
          <p class="text-lg text-gray-600">
            {$tStore.chat.greetingQuestion}
          </p>
        </div>

        <!-- Botones seleccionables - siempre visibles cuando no hay mensajes -->
        <div class="flex flex-col gap-3 w-full">
          <button 
            class="p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-2xl transition-all text-left group flex items-center gap-3"
            onclick={() => handleButtonClick('whatsapp')}
          >
            <span class="text-2xl">📱</span>
            <span class="text-base font-normal text-gray-900">{$tStore.chat.updateWhatsApp}</span>
          </button>

          <button 
            class="p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-2xl transition-all text-left group flex items-center gap-3"
            onclick={() => handleButtonClick('address')}
          >
            <span class="text-2xl">📍</span>
            <span class="text-base font-normal text-gray-900">{$tStore.chat.updateAddress}</span>
          </button>

          <button 
            class="p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-2xl transition-all text-left group flex items-center gap-3"
            onclick={() => handleButtonClick('dishes')}
          >
            <span class="text-2xl">🍽️</span>
            <span class="text-base font-normal text-gray-900">{$tStore.chat.addDishes}</span>
          </button>

          <button 
            class="p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-2xl transition-all text-left group flex items-center gap-3"
            onclick={() => handleButtonClick('desserts')}
          >
            <span class="text-2xl">🍰</span>
            <span class="text-base font-normal text-gray-900">{$tStore.chat.addDesserts}</span>
          </button>

          <button 
            class="p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-2xl transition-all text-left group flex items-center gap-3"
            onclick={() => handleButtonClick('price')}
          >
            <span class="text-2xl">💰</span>
            <span class="text-base font-normal text-gray-900">{$tStore.chat.updatePrice}</span>
          </button>

          <button 
            class="p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-2xl transition-all text-left group flex items-center gap-3"
            onclick={() => handleButtonClick('delete')}
          >
            <span class="text-2xl">🗑️</span>
            <span class="text-base font-normal text-gray-900">{$tStore.chat.deleteItem}</span>
          </button>
        </div>
      </div>
    {:else}
      <div class="max-w-3xl mx-auto space-y-6 pt-4">
        {#each messages as message (message.id)}
          <Message message={message} onExploreOptions={resetToOptions} />
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
    <div class="border-t border-gray-200 bg-white sticky bottom-0 z-10 safe-area-inset-bottom">
    <div class="max-w-3xl mx-auto px-4 py-3">
      {#if checkingMenu}
        <div class="flex items-center justify-center py-8">
          <div class="text-center">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 border-t-transparent mx-auto mb-4"></div>
            <p class="text-gray-600 text-sm">{$tStore.chat.checkingMenu}</p>
          </div>
        </div>
      {/if}
      
      <ChatInput bind:this={chatInputRef} onSend={handleSendMessage} disabled={isLoading || checkingMenu} onFocus={handleInputFocus} onBlur={handleInputBlur} />
      
      <!-- Sugerencias justo debajo del input -->
      {#if showExamples && currentExampleType}
        <div class="mt-3 px-2">
          <p class="text-xs text-gray-500 mb-2">{$tStore.chat.examplesLabel}</p>
          <div class="flex flex-col gap-2">
            {#each getExampleMessages() as example}
              <div class="text-sm font-normal text-gray-600 py-1.5">
                {example}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
    </div>
  </div>

  <!-- Vista de Preview (se muestra cuando showPreview es true) -->
  <div 
    class="absolute inset-0 flex flex-col h-full bg-white transition-transform duration-300 ease-in-out {showPreview ? 'translate-x-0' : 'translate-x-full'}"
  >
    <!-- Header del Preview -->
    <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between flex-shrink-0 z-10">
      <button
        onclick={() => {
          showPreview = false
          linkWasCopied = false
        }}
        class="p-2 hover:bg-gray-100 rounded-full transition-colors"
        aria-label="Volver"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <div class="flex-1 min-w-0 px-2">
        <h2 class="text-base md:text-lg font-semibold text-gray-900 truncate text-center">{$tStore.chat.previewTitle}</h2>
        {#if trialDaysRemaining !== null && entitlement?.status === 'trialing'}
          <p class="text-xs md:text-sm text-gray-600 mt-1 text-center">
            {$tStore.chat.trialDaysRemaining.replace('{days}', trialDaysRemaining.toString())}
          </p>
        {/if}
      </div>
      <div class="w-9"></div> <!-- Spacer para centrar -->
    </header>

    <!-- Contenido del Preview -->
    <div class="flex-1 overflow-hidden iframe-container relative" style="padding-bottom: {menuUrl ? '90px' : '0'};">
      {#if menuUrl}
        {#key iframeKey}
          <iframe
            src={menuUrl}
            class="w-full h-full border-0"
            title="Vista previa de la carta"
            loading="lazy"
            allow="camera; microphone; geolocation; autoplay; clipboard-write"
            sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-popups-to-escape-sandbox allow-presentation"
          ></iframe>
        {/key}
      {:else}
        <div class="flex items-center justify-center h-full">
          <div class="text-center">
            <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600 border-t-transparent mx-auto mb-4"></div>
            <p class="text-gray-600">{$tStore.chat.loadingPreview}</p>
          </div>
        </div>
      {/if}
    </div>

    <!-- Botón flotante de compartir (barra fija dentro del contenedor para desktop) -->
    {#if menuUrl}
      <div class="hidden md:block bg-white border-t border-gray-200 px-4 py-3 flex-shrink-0 safe-area-inset-bottom shadow-lg">
        <button
          onclick={shareOnWhatsApp}
          class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-3"
          title={$tStore.chat.shareLink}
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
          <span>{currentLanguage === 'ES' ? 'Compartir Carta' : currentLanguage === 'PT' ? 'Compartilhar Cardápio' : 'Share Menu'}</span>
        </button>
      </div>
    {/if}
  </div>

<!-- Botón flotante de compartir (fixed fuera del contenedor para móviles) -->
{#if menuUrl && showPreview && !showSlugModal}
  <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-4 py-3 z-[100] safe-area-inset-bottom shadow-2xl md:hidden">
    <button
      onclick={shareOnWhatsApp}
      class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-3"
      title={$tStore.chat.shareLink}
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
      </svg>
      <span>{currentLanguage === 'ES' ? 'Compartir Carta' : currentLanguage === 'PT' ? 'Compartilhar Cardápio' : 'Share Menu'}</span>
    </button>
  </div>
{/if}

  <!-- Vista de Crear Slug (se muestra cuando showSlugModal es true) -->
  <div 
    class="absolute inset-0 flex flex-col h-full bg-white transition-transform duration-300 ease-in-out {showSlugModal ? 'translate-x-0' : 'translate-x-full'}"
  >
    <!-- Confeti -->
    {#if showConfetti}
      <div class="absolute inset-0 pointer-events-none overflow-hidden z-50">
        {#each confettiParticles as particle, i}
          <div 
            class="absolute confetti"
            style="left: {particle.left}%; animation-delay: {particle.delay}s; background: {particle.color};"
          ></div>
        {/each}
      </div>
    {/if}

    <!-- Header -->
    <header class="border-b border-gray-200 bg-white px-4 py-3 flex items-center justify-between">
      <button
        onclick={() => {
          if (!isCreatingSlug) {
            showSlugModal = false
            businessName = ''
            slugPreview = ''
          }
        }}
        disabled={isCreatingSlug}
        class="p-2 hover:bg-gray-100 rounded-full transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        aria-label="Volver"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h2 class="text-lg font-semibold text-gray-900 flex-1 text-center">
        {currentLanguage === 'ES' 
          ? 'Crea tu enlace personalizado'
          : currentLanguage === 'PT'
          ? 'Crie seu link personalizado'
          : 'Create your custom link'}
      </h2>
      <div class="w-9"></div> <!-- Spacer para centrar -->
    </header>

    <!-- Contenido -->
    <div class="flex-1 overflow-y-auto px-4 py-6">
      <div class="max-w-md mx-auto">
        <!-- Icono y descripción -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-green-500 to-emerald-600 rounded-full mb-4">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            </svg>
          </div>
          <p class="text-gray-600 text-base">
            {currentLanguage === 'ES'
              ? 'Indica el nombre de tu negocio para generar tu enlace único'
              : currentLanguage === 'PT'
              ? 'Indique o nome do seu negócio para gerar seu link único'
              : 'Enter your business name to generate your unique link'}
          </p>
        </div>

        <!-- Input del nombre -->
        <div class="mb-6">
          <label for="business-name-input" class="block text-sm font-medium text-gray-700 mb-2">
            {currentLanguage === 'ES'
              ? 'Nombre de tu negocio'
              : currentLanguage === 'PT'
              ? 'Nome do seu negócio'
              : 'Business name'}
          </label>
          <input
            id="business-name-input"
            type="text"
            bind:value={businessName}
            placeholder={currentLanguage === 'ES' ? 'Ej: Mi Restaurante' : currentLanguage === 'PT' ? 'Ex: Meu Restaurante' : 'Ex: My Restaurant'}
            class="w-full px-4 py-3 border-2 border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent text-base transition-all"
            disabled={isCreatingSlug}
            onkeydown={(e) => {
              if (e.key === 'Enter' && !isCreatingSlug && businessName.trim()) {
                createSlugAndShare()
              }
            }}
          />
        </div>

        <!-- Preview de la URL -->
        <div class="mb-8">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            {currentLanguage === 'ES'
              ? 'Tu enlace será:'
              : currentLanguage === 'PT'
              ? 'Seu link será:'
              : 'Your link will be:'}
          </label>
          <div class="bg-gradient-to-r from-green-50 to-emerald-50 border-2 border-green-200 rounded-lg p-4">
            <p class="text-sm text-gray-800 font-mono break-all font-semibold">
              {slugPreview 
                ? `catalogo.micartapro.com/m/${slugPreview}`
                : currentLanguage === 'ES'
                  ? 'Ingresa el nombre de tu negocio...'
                  : currentLanguage === 'PT'
                  ? 'Digite o nome do seu negócio...'
                  : 'Enter your business name...'}
            </p>
          </div>
        </div>

        <!-- Botón de acción -->
        <button
          onclick={createSlugAndShare}
          disabled={isCreatingSlug || !businessName.trim() || !slugPreview}
          class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white rounded-lg transition-all font-semibold disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 text-base shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
        >
          {#if isCreatingSlug}
            <div class="animate-spin rounded-full h-6 w-6 border-2 border-white border-t-transparent"></div>
            <span>{currentLanguage === 'ES' ? 'Creando...' : currentLanguage === 'PT' ? 'Criando...' : 'Creating...'}</span>
          {:else}
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
            </svg>
            <span>{currentLanguage === 'ES' ? 'Crear y Compartir' : currentLanguage === 'PT' ? 'Criar e Compartilhar' : 'Create & Share'}</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
</div>

<!-- Modal de Upgrade al copiar link -->
{#if showUpgradeModal}
  <div 
    class="fixed inset-0 bg-black/50 backdrop-blur-sm z-[60] flex items-center justify-center p-4"
    onclick={() => showUpgradeModal = false}
    role="dialog"
    aria-modal="true"
  >
    <div 
      class="relative bg-gradient-to-br from-blue-50 via-white to-indigo-50 rounded-xl md:rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Decorative gradient background -->
      <div class="absolute top-0 left-0 right-0 h-32 bg-gradient-to-r from-blue-500/10 via-purple-500/10 to-indigo-500/10"></div>
      
      <!-- Close button -->
      <button
        onclick={() => showUpgradeModal = false}
        class="absolute top-4 right-4 p-2 hover:bg-white/50 rounded-full transition-colors z-10"
        aria-label={$tStore.chat.close}
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
      
      <div class="relative p-5 md:p-8">
        <!-- Header -->
        <div class="text-center mb-4 md:mb-6">
          <div class="inline-flex items-center justify-center w-12 h-12 md:w-16 md:h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full mb-3 md:mb-4">
            <svg class="w-6 h-6 md:w-8 md:h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h3 class="text-xl md:text-2xl font-bold text-gray-900 mb-2">{$tStore.chat.upgradeModalTitle}</h3>
        </div>
        
        <!-- Benefits -->
        <div class="bg-white/60 backdrop-blur-sm rounded-xl p-5 mb-6 border border-gray-200/50">
          <h4 class="text-sm font-semibold text-gray-900 mb-3">{$tStore.chat.upgradeModalBenefits}</h4>
          <div class="space-y-2">
            <div class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-sm text-gray-700">{$tStore.chat.upgradeModalBenefit1}</span>
            </div>
            <div class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-sm text-gray-700">{$tStore.chat.upgradeModalBenefit2}</span>
            </div>
            <div class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-sm text-gray-700">{$tStore.chat.upgradeModalBenefit3}</span>
            </div>
          </div>
        </div>
        
        <!-- Contribution message -->
        <div class="bg-blue-50/80 rounded-xl p-4 mb-6 border border-blue-200/50">
          <p class="text-sm text-blue-900 text-center">
            💙 {$tStore.chat.upgradeModalContribution}
          </p>
        </div>
        
        <!-- Action button -->
        <div>
          <button
            onclick={handleUpgradeToPro}
            class="w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white rounded-xl transition-all font-semibold shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
          >
            {$tStore.chat.upgradeToPro}
          </button>
        </div>
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
  
  /* Transición de vista para el preview */
  .transition-transform {
    will-change: transform;
  }

  /* Animación de confeti */
  .confetti {
    width: 10px;
    height: 10px;
    position: absolute;
    top: -10px;
    animation: confetti-fall 3s linear forwards;
    border-radius: 2px;
  }

  @keyframes confetti-fall {
    0% {
      transform: translateY(0) rotate(0deg);
      opacity: 1;
    }
    100% {
      transform: translateY(100vh) rotate(720deg);
      opacity: 0;
    }
  }
</style>

