<script lang="ts">
  import { onMount } from 'svelte'
  import { v7 as uuidv7 } from 'uuid'
  import Message from './Message.svelte'
  import ChatInput from './ChatInput.svelte'
  import { authState } from '../auth.svelte'
  import { getLatestMenuId, generateMenuUrl, pollUntilMenuUpdated, pollUntilMenuExists, pollUntilVersionExists, getMenuSlug, createMenuSlug, generateSlugUrl, updateCurrentVersionId, getUserCredits, hasEnoughCredits } from '../menuUtils'
  import { API_BASE_URL } from '../config'
  import { t as tStore, language } from '../useLanguage'
  import { supabase } from '../supabase'
  import { resizeImage } from '../utils/imageResize'

  interface MenuChatProps {
    onMenuClick?: () => void
  }

  let { onMenuClick }: MenuChatProps = $props()

  interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: Date
    showExploreButton?: boolean
    imageUrl?: string
    isPreview?: boolean // Indica si es un preview pendiente de enviar
    pendingVersionId?: string // ID de la versión pendiente de activar
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
  let iframeKey = $state(0) // Key para forzar recarga del iframe cuando el menú se actualiza
  let showSlugModal = $state(false) // Modal para crear slug
  let businessName = $state('') // Nombre del negocio
  let slugPreview = $state('') // Preview del slug
  let isCreatingSlug = $state(false) // Estado de carga al crear slug
  let showConfetti = $state(false) // Mostrar confeti
  let confettiParticles = $state<Array<{left: number, delay: number, color: string}>>([]) // Partículas de confeti
  let showSubscriptionPromo = $state(false) // Modal promocional de suscripción
  let showCreditsPromo = $state(false) // Modal promocional de créditos
  let userCredits = $state<number | null>(null) // Créditos del usuario
  let showCamera = $state(false) // Mostrar cámara para tomar foto
  let stream: MediaStream | null = $state(null) // Stream de la cámara
  let videoRef: HTMLVideoElement | null = $state(null) // Referencia al video
  let uploadingPhoto = $state(false) // Estado de subida de foto
  let pendingPhotoUrl = $state<string | null>(null) // URL de la foto pendiente de enviar
  let isActivatingVersion = $state(false) // Estado de activación de versión
  let versionActivated = $state(false) // Flag para saber si la versión ya fue activada
  let showDiscardModal = $state(false) // Modal de confirmación para descartar menú
  let pendingNavigation: (() => void) | null = $state(null) // Callback para ejecutar después de confirmar descarte

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
  


  // Función para verificar si hay un menú pendiente
  function hasPendingMenu(): boolean {
    const hasPending = messages.some(msg => msg.pendingVersionId && !versionActivated)
    console.log('🔍 hasPendingMenu - resultado:', hasPending, {
      messages: messages.map(m => ({ id: m.id, pendingVersionId: m.pendingVersionId })),
      versionActivated
    })
    return hasPending
  }

  // Función para manejar la navegación hacia atrás
  function handleBackNavigation(event?: PopStateEvent) {
    if (hasPendingMenu()) {
      // Prevenir la navegación por defecto volviendo a la misma página
      if (event) {
        // Volver a la misma página para interceptar el retroceso
        window.history.pushState(null, '', window.location.href)
      }
      
      // Mostrar modal de confirmación
      showDiscardModal = true
      pendingNavigation = async () => {
        // Limpiar el pendingVersionId del mensaje
        messages = messages.map(msg => ({
          ...msg,
          pendingVersionId: undefined
        }))
        versionActivated = false
        showDiscardModal = false
        pendingNavigation = null
        
        // Actualizar menuUrl para que no incluya version_id (mostrar versión actual)
        if (menuId && session?.access_token) {
          try {
            const url = await generateMenuUrl(menuId, session.access_token, currentLanguage)
            if (url) {
              menuUrl = url
              // Forzar recarga del iframe
              iframeKey++
            }
          } catch (err) {
            console.error('Error actualizando URL del menú después de descartar:', err)
          }
        }
        
        // Ejecutar la navegación después de limpiar
        setTimeout(() => {
          window.history.back()
        }, 100)
      }
    }
  }

  // Función para confirmar descarte
  async function confirmDiscard() {
    if (pendingNavigation) {
      await pendingNavigation()
      pendingNavigation = null
    }
  }

  // Función para cancelar descarte
  function cancelDiscard() {
    showDiscardModal = false
    pendingNavigation = null
  }

  onMount(() => {
    // No mostrar mensajes inicialmente, solo los botones
    messages = []
    
    // Agregar un estado inicial al historial para poder interceptar el retroceso
    window.history.pushState(null, '', window.location.href)
    
    // Listener para detectar navegación hacia atrás
    const handlePopState = (event: PopStateEvent) => {
      handleBackNavigation(event)
    }
    
    // Listener para detectar cierre de página/pestaña
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      if (hasPendingMenu()) {
        event.preventDefault()
        event.returnValue = '' // Chrome requiere returnValue
        return '' // Algunos navegadores requieren return
      }
    }
    
    window.addEventListener('popstate', handlePopState)
    window.addEventListener('beforeunload', handleBeforeUnload)
    
    // Actualizar menuUrl cuando cambie el idioma
    $effect(() => {
      if (menuUrl && userId && menuId && session?.access_token) {
        // Reconstruir la URL con el idioma actual
        // Extraer version_id de la URL actual si existe
        const currentUrl = new URL(menuUrl)
        const currentVersionId = currentUrl.searchParams.get('version_id')
        generateMenuUrl(menuId, session.access_token, currentLanguage, currentVersionId || undefined)
          .then(newUrl => {
            if (newUrl) {
              menuUrl = newUrl
            }
          })
          .catch(err => {
            console.error('Error actualizando URL del menú:', err)
          })
      }
    })
    
    // Cargar menuID y créditos al montar el componente
    ;(async () => {
      if (userId && session?.access_token) {
        try {
          // Cargar créditos del usuario
          const credits = await getUserCredits(session.access_token)
          if (credits !== null) {
            userCredits = credits.balance
          }

          const id = await getLatestMenuId(userId, session.access_token)
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
    })()
    
    // Cleanup al desmontar
    return () => {
      window.removeEventListener('popstate', handlePopState)
      window.removeEventListener('beforeunload', handleBeforeUnload)
    }
  })

  async function togglePreview() {
    if (showPreview) {
      // Verificar si hay un menú pendiente
      const hasPending = hasPendingMenu()
      console.log('🔍 togglePreview - Cerrando vista previa. hasPendingMenu:', hasPending)
      console.log('🔍 togglePreview - messages:', messages.map(m => ({ id: m.id, pendingVersionId: m.pendingVersionId })))
      console.log('🔍 togglePreview - versionActivated:', versionActivated)
      
      // Si hay un menú pendiente, mostrar modal de confirmación
      if (hasPending) {
        console.log('⚠️ togglePreview - Mostrando modal de descarte')
        showDiscardModal = true
        pendingNavigation = async () => {
          console.log('✅ togglePreview - Confirmando descarte')
          // Limpiar el pendingVersionId del mensaje
          messages = messages.map(msg => ({
            ...msg,
            pendingVersionId: undefined
          }))
          versionActivated = false
          showDiscardModal = false
          pendingNavigation = null
          
          // Actualizar menuUrl para que no incluya version_id (mostrar versión actual)
          if (menuId && session?.access_token) {
            try {
              const url = await generateMenuUrl(menuId, session.access_token, currentLanguage)
              if (url) {
                menuUrl = url
                // Forzar recarga del iframe
                iframeKey++
              }
            } catch (err) {
              console.error('Error actualizando URL del menú después de descartar:', err)
            }
          }
          
          // Cerrar la vista previa
          showPreview = false
        }
        return
      }
      // Si no hay menú pendiente, cerrar directamente
      console.log('✅ togglePreview - No hay menú pendiente, cerrando directamente')
      showPreview = false
      return
    }

    if (!menuUrl && userId && menuId && session?.access_token) {
      try {
        const url = await generateMenuUrl(menuId, session.access_token, currentLanguage)
        if (url) {
          menuUrl = url
        }
      } catch (err) {
        console.error('Error cargando menú:', err)
      }
    }
    
    showPreview = true
  }

  async function copyToClipboard() {
    if (!menuUrl) return

    try {
      await navigator.clipboard.writeText(menuUrl)
      copySuccess = true
      linkWasCopied = true // Marcar que el enlace fue copiado
      setTimeout(() => {
        copySuccess = false
      }, 2000)
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

  async function handleUseThisMenu(versionId: string) {
    console.log('🚀 handleUseThisMenu llamado con:', { versionId, menuId, hasAccessToken: !!session?.access_token })
    
    if (!menuId || !session?.access_token || !versionId) {
      console.error('❌ Faltan parámetros:', { menuId, hasAccessToken: !!session?.access_token, versionId })
      const lang = currentLanguage
      alert(
        lang === 'ES'
          ? 'No se puede activar: falta información del menú'
          : lang === 'PT'
          ? 'Não é possível ativar: falta informação do cardápio'
          : 'Cannot activate: missing menu information'
      )
      return
    }

    isActivatingVersion = true

    try {
      console.log('🔄 Iniciando actualización de current_version_id:', { menuId, versionId, accessTokenLength: session.access_token.length })
      const success = await updateCurrentVersionId(menuId, versionId, session.access_token)
      console.log('📊 Resultado de updateCurrentVersionId:', success)
      
      if (success) {
        versionActivated = true
        
        // Actualizar el mensaje para quitar el botón "Usar este menú"
        messages = messages.map(msg => {
          if (msg.pendingVersionId === versionId) {
            return {
              ...msg,
              pendingVersionId: undefined,
              content: currentLanguage === 'ES' 
                ? '✅ Menú actualizado y activado. Ya puedes compartirlo.'
                : currentLanguage === 'PT'
                ? '✅ Cardápio atualizado e ativado. Agora você pode compartilhar.'
                : '✅ Menu updated and activated. You can now share it.'
            }
          }
          return msg
        })

        // Actualizar la URL del menú para mostrar la versión activada (sin version_id, usa current_version_id)
        if (userId && menuId && session?.access_token) {
          const url = await generateMenuUrl(menuId, session.access_token, currentLanguage)
          if (url) {
            menuUrl = url
          }
        }
        
        // Forzar recarga del iframe para mostrar la versión activada
        iframeKey++
        
        // Mostrar mensaje de éxito
        const lang = currentLanguage
        const successMsg = lang === 'ES'
          ? '¡Menú activado exitosamente!'
          : lang === 'PT'
          ? 'Cardápio ativado com sucesso!'
          : 'Menu activated successfully!'
        
        const activatedMessage: ChatMessage = {
          id: `activated-${Date.now()}`,
          role: 'assistant',
          content: successMsg,
          timestamp: new Date()
        }
        messages = [...messages, activatedMessage]
      } else {
        throw new Error('Failed to update version')
      }
    } catch (error: any) {
      console.error('Error activando versión:', error)
      const lang = currentLanguage
      alert(
        lang === 'ES'
          ? 'Error al activar el menú. Intenta nuevamente.'
          : lang === 'PT'
          ? 'Erro ao ativar o cardápio. Tente novamente.'
          : 'Error activating menu. Please try again.'
      )
    } finally {
      isActivatingVersion = false
    }
  }

  async function shareOnWhatsApp() {
    if (!menuId || !session?.access_token || !userId) {
      alert('No se puede compartir: falta información del menú')
      return
    }

    try {
      // Verificar si existe un slug para este menú (siempre permitir compartir)
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

  async function handleActivatePlan() {
    if (!session?.access_token) {
      alert('No hay sesión activa. Por favor, inicia sesión nuevamente.')
      return
    }

    try {
      showSubscriptionPromo = false
      showCreditsPromo = false
      
      // Llamar al endpoint de checkout de MercadoPago del backend
      const response = await fetch(`${API_BASE_URL}/checkout/mercadopago`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        const errorText = await response.text()
        console.error('Error obteniendo checkout URL:', errorText)
        alert('Error al obtener la URL de checkout. Por favor, intenta de nuevo.')
        return
      }

      const data = await response.json()
      const checkoutUrl = data.checkout_url

      if (!checkoutUrl) {
        alert('No se recibió la URL de checkout. Por favor, intenta de nuevo.')
        return
      }

      // Redirigir a la URL de checkout
      window.open(checkoutUrl, '_blank')
      
      // Refrescar créditos después de un momento (el webhook puede tardar)
      setTimeout(async () => {
        await refreshCredits()
      }, 3000)
    } catch (error) {
      console.error('Error activando plan:', error)
      alert('Error al activar el plan. Por favor, intenta de nuevo.')
    }
  }

  function handleBackToEdit() {
    showSubscriptionPromo = false
    showCreditsPromo = false
    // Cerrar preview si está abierto
    if (showPreview) {
      showPreview = false
    }
  }

  async function refreshCredits() {
    if (session?.access_token) {
      const credits = await getUserCredits(session.access_token)
      if (credits !== null) {
        userCredits = credits.balance
      }
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
    // Solo permitir enviar si hay texto escrito
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

    // Guardar la URL de la foto antes de limpiarla
    const photoUrlToSend = pendingPhotoUrl

    // Remover el preview de la foto del chat
    messages = messages.filter(msg => !msg.isPreview)

    // Agregar mensaje del usuario (con foto si existe)
    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: content.trim() || (photoUrlToSend ? '📸 Foto' : ''),
      timestamp: new Date(),
      imageUrl: photoUrlToSend || undefined
    }
    messages = [...messages, userMessage]

    // Limpiar la foto pendiente
    pendingPhotoUrl = null

    // Mostrar indicador de carga
    isLoading = true

    try {
      // Generar idempotency key (UUID v7)
      const idempotencyKey = uuidv7()

      // Preparar el body del request
      const requestBody: any = {
        menuId: menuId,
        message: content.trim() || ''
      }

      // Agregar photoUrl si existe
      if (photoUrlToSend) {
        requestBody.photoUrl = photoUrlToSend
      }

      // Enviar POST al backend
      const response = await fetch(`${API_BASE_URL}/menu/interaction`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
          'Idempotency-Key': idempotencyKey,
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify(requestBody)
      })

      if (!response.ok) {
        // Si el error es 402 (Payment Required), significa que no hay suficientes créditos
        if (response.status === 402) {
          isLoading = false
          // Obtener créditos actuales para mostrar en el modal
          const credits = await getUserCredits(session.access_token)
          userCredits = credits?.balance || 0
          showCreditsPromo = true
          return
        }
        throw new Error(`Error del servidor: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('📥 Respuesta completa del backend:', JSON.stringify(data, null, 2))
      console.log('📥 versionId recibido:', data.versionId)
      console.log('📥 Tipo de data.versionId:', typeof data.versionId)
      console.log('📥 data.versionId es truthy?', !!data.versionId)

      // Agregar respuesta inicial del asistente
      const initialMessage: ChatMessage = {
        id: `assistant-${Date.now()}`,
        role: 'assistant',
        content: data.message || data.response || 'Procesando tu mensaje...',
        timestamp: new Date()
      }
      messages = [...messages, initialMessage]

      // Iniciar polling para esperar la actualización del menú usando versionID
      if (data.versionId && session?.access_token) {
        console.log('✅ Entrando en bloque con versionId, iniciando polling...')
        try {
          // Usar el versionID retornado para consultar directamente Supabase
          // Aumentamos los intentos a 120 (2 minutos) para procesos que pueden tardar más
          const updatedMenu = await pollUntilVersionExists(data.versionId, session.access_token, 120, 1000)
          
          // Agregar mensaje de confirmación con el menú actualizado
          // Incluir el versionId pendiente para mostrar el botón "Usar este menú"
          console.log('✅ Creando mensaje de éxito con versionId:', data.versionId)
          console.log('✅ Tipo de versionId:', typeof data.versionId)
          console.log('✅ versionId es válido?', !!data.versionId && data.versionId.length > 0)
          
          const successMessage: ChatMessage = {
            id: `success-${Date.now()}`,
            role: 'assistant',
            content: $tStore.chat.successUpdated,
            timestamp: new Date(),
            pendingVersionId: data.versionId // Versión pendiente de activar
          }
          console.log('✅ Mensaje de éxito creado:', JSON.stringify(successMessage, null, 2))
          messages = [...messages, successMessage]
          console.log('✅ Total de mensajes después de agregar éxito:', messages.length)
          console.log('✅ Último mensaje tiene pendingVersionId:', messages[messages.length - 1].pendingVersionId)
          console.log('✅ Valor exacto de pendingVersionId:', messages[messages.length - 1].pendingVersionId)
          versionActivated = false // Resetear el flag cuando hay nueva versión
          
          // Refrescar créditos después de consumir uno
          await refreshCredits()
          
          // Actualizar la URL del menú y forzar recarga del iframe
          // Incluir version_id en la URL para mostrar la versión específica (interacción)
          if (userId && menuId && session?.access_token) {
            const url = await generateMenuUrl(menuId, session.access_token, currentLanguage, data.versionId)
            if (url) {
              menuUrl = url
            }
            // Incrementar iframeKey para forzar recarga del iframe
            iframeKey++
          }
          
          // Abrir automáticamente la vista previa para mostrar los cambios
          showPreview = true
        } catch (pollError: any) {
          console.error('❌ Error en polling:', pollError)
          
          // Agregar mensaje de error del polling
          const errorMessage: ChatMessage = {
            id: `poll-error-${Date.now()}`,
            role: 'assistant',
            content: $tStore.chat.errorPolling.replace('{message}', pollError.message || 'Error desconocido'),
            timestamp: new Date()
          }
          messages = [...messages, errorMessage]
        }
      } else {
        console.log('⚠️ No hay versionId o access_token. versionId:', data.versionId, 'access_token:', !!session?.access_token)
        // Fallback: si no viene versionId, usar el método anterior (GCS)
        if (userId && menuId) {
          try {
            const updatedMenu = await pollUntilMenuUpdated(userId, menuId, idempotencyKey)
            
            // Agregar mensaje de confirmación con el menú actualizado
            // Nota: En el fallback no tenemos versionId, así que no podemos mostrar el botón "Usar este menú"
            const successMessage: ChatMessage = {
              id: `success-${Date.now()}`,
              role: 'assistant',
              content: $tStore.chat.successUpdated,
              timestamp: new Date(),
              showExploreButton: true
            }
            messages = [...messages, successMessage]
            versionActivated = false // Resetear el flag
            
            // Actualizar la URL del menú y forzar recarga del iframe
            // Sin version_id para mostrar la versión actual (visualización simple)
            if (session?.access_token) {
              const url = await generateMenuUrl(menuId, session.access_token, currentLanguage)
              if (url) {
                menuUrl = url
              }
            }
            iframeKey++
            
            // Abrir automáticamente la vista previa para mostrar los cambios
            showPreview = true
          } catch (pollError: any) {
            console.error('Error en polling (fallback):', pollError)
            
            const errorMessage: ChatMessage = {
              id: `poll-error-${Date.now()}`,
              role: 'assistant',
              content: $tStore.chat.errorPolling.replace('{message}', pollError.message || 'Error desconocido'),
              timestamp: new Date()
            }
            messages = [...messages, errorMessage]
          }
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
    
    // Verificar si hay texto en el input antes de resetear
    const textarea = chatInputRef?.textareaRef
    if (textarea && textarea.value?.trim()) {
      // Hay texto, no resetear
      return
    }
    
    // En móvil, esperar más tiempo para detectar si el teclado se cerró realmente
    const delay = window.innerWidth <= 768 ? 500 : 200
    
    setTimeout(() => {
      // Verificar nuevamente si el textarea todavía tiene focus (en algunos móviles el blur puede ser temporal)
      const isStillFocused = document.activeElement === textarea
      
      // Verificar nuevamente si hay texto (puede haber cambiado mientras esperábamos)
      const hasText = textarea && textarea.value?.trim()
      
      if (!isStillFocused && !hasText) {
        // Solo resetear si realmente perdió el focus, no hay texto, y no hay sugerencias mostrándose
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

  // Bucket de Supabase Storage para las fotos
  const BUCKET_NAME = 'menu-photos'

  function handleTakePhoto() {
    // Abrir cámara directamente
    startCamera()
  }

  function handleUploadPhoto() {
    // Abrir selector de archivo
    const input = document.getElementById('photo-file-input') as HTMLInputElement
    if (input) {
      input.click()
    }
  }

  async function startCamera() {
    try {
      // Solicitar acceso a la cámara
      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: 'environment',
          width: { ideal: 1280 },
          height: { ideal: 720 }
        }
      })
      
      stream = mediaStream
      showCamera = true
      
      // Esperar un momento para que el DOM se actualice
      await new Promise(resolve => setTimeout(resolve, 100))
      
      // Esperar a que el video esté listo
      if (videoRef) {
        videoRef.srcObject = mediaStream
        
        // Esperar a que el video esté cargado y reproduciéndose
        await new Promise((resolve, reject) => {
          if (!videoRef) {
            reject(new Error('Video ref no disponible'))
            return
          }
          
          const onLoadedMetadata = () => {
            videoRef?.removeEventListener('loadedmetadata', onLoadedMetadata)
            videoRef?.play()
              .then(() => {
                setTimeout(resolve, 200)
              })
              .catch(reject)
          }
          
          videoRef.addEventListener('loadedmetadata', onLoadedMetadata)
          
          setTimeout(() => {
            videoRef?.removeEventListener('loadedmetadata', onLoadedMetadata)
            reject(new Error('Timeout esperando video'))
          }, 5000)
        })
      }
    } catch (error) {
      console.error('Error accediendo a la cámara:', error)
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: 'No se pudo acceder a la cámara. Asegúrate de dar permisos de cámara.',
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
      stopCamera()
    }
  }

  function stopCamera() {
    if (stream) {
      stream.getTracks().forEach(track => track.stop())
      stream = null
    }
    if (videoRef) {
      videoRef.srcObject = null
    }
    showCamera = false
  }

  function capturePhoto() {
    if (!videoRef) return

    if (videoRef.readyState !== videoRef.HAVE_ENOUGH_DATA) {
      return
    }

    if (videoRef.videoWidth === 0 || videoRef.videoHeight === 0) {
      return
    }

    try {
      const canvas = document.createElement('canvas')
      canvas.width = videoRef.videoWidth
      canvas.height = videoRef.videoHeight
      
      const ctx = canvas.getContext('2d')
      if (!ctx) return

      ctx.drawImage(videoRef, 0, 0, canvas.width, canvas.height)
      
      canvas.toBlob(async (blob) => {
        if (!blob) return

        const file = new File([blob], `photo-${Date.now()}.jpg`, {
          type: 'image/jpeg',
          lastModified: Date.now()
        })

        stopCamera()
        await uploadAndSendPhoto(file)
      }, 'image/jpeg', 0.9)
    } catch (error) {
      console.error('Error capturando foto:', error)
      stopCamera()
    }
  }

  async function handleFileSelect(event: Event) {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    
    if (!file || !file.type.startsWith('image/')) return
    
    await uploadAndSendPhoto(file)
    
    // Resetear el input
    if (target) target.value = ''
  }

  function handlePastePhoto(file: File) {
    // Procesar el archivo pegado de la misma manera que un archivo seleccionado
    uploadAndSendPhoto(file)
  }

  async function uploadAndSendPhoto(file: File) {
    if (!userId) return

    uploadingPhoto = true

    try {
      // Redimensionar la imagen
      const resizedBlob = await resizeImage(file, 400, 400, 0.8)

      // Generar nombre único para el archivo
      const fileExt = file.name.split('.').pop() || 'jpg'
      const fileName = `${userId}/${Date.now()}-${Math.random().toString(36).substring(7)}.${fileExt}`
      
      // Subir a Supabase Storage
      const { error } = await supabase.storage
        .from(BUCKET_NAME)
        .upload(fileName, resizedBlob, {
          contentType: resizedBlob.type,
          upsert: false
        })

      if (error) {
        throw error
      }

      // Obtener la URL pública de la imagen
      const { data: urlData } = supabase.storage
        .from(BUCKET_NAME)
        .getPublicUrl(fileName)

      const publicUrl = urlData.publicUrl
      
      // Imprimir la URL en la consola
      console.log('✅ Foto subida exitosamente!')
      console.log('📸 URL de la foto:', publicUrl)

      // Guardar la URL de la foto pendiente (no enviar mensaje todavía)
      pendingPhotoUrl = publicUrl

      // Mostrar preview de la foto en el chat
      const photoPreview: ChatMessage = {
        id: `photo-preview-${Date.now()}`,
        role: 'user',
        content: '📸 Foto lista para enviar',
        timestamp: new Date(),
        imageUrl: publicUrl,
        isPreview: true
      }
      messages = [...messages, photoPreview]

      // Agregar mensaje del asistente preguntando a qué parte del catálogo pertenece
      const questionMessage: ChatMessage = {
        id: `question-${Date.now()}`,
        role: 'assistant',
        content: '¿A qué parte de tu catálogo pertenece esta foto?',
        timestamp: new Date()
      }
      messages = [...messages, questionMessage]

    } catch (error) {
      console.error('Error subiendo foto:', error)
      const errorMessage: ChatMessage = {
        id: `error-${Date.now()}`,
        role: 'assistant',
        content: 'Error al subir la foto. Intenta de nuevo.',
        timestamp: new Date()
      }
      messages = [...messages, errorMessage]
    } finally {
      uploadingPhoto = false
    }
  }

  // Limpiar cámara al desmontar
  $effect(() => {
    return () => {
      stopCamera()
    }
  })

  $effect(() => {
    if (messages.length > 0) {
      scrollToBottom()
    }
  })
</script>

<div class="flex flex-col h-screen h-[100dvh] bg-white relative overflow-hidden">
  <!-- Vista de Chat (oculta cuando showPreview, showSlugModal o showSubscriptionPromo es true) -->
  <div
    class="flex flex-col h-full transition-transform duration-300 ease-in-out {(showPreview || showSlugModal || showSubscriptionPromo || showCreditsPromo) ? '-translate-x-full' : 'translate-x-0'} min-h-0"
    style="max-height: 100dvh;"
  >
    <!-- Header estilo Gemini -->
    <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between flex-shrink-0 sticky top-0 z-20 bg-white">
      <!-- Botón hamburguesa para móvil -->
      <button 
        onclick={onMenuClick}
        class="md:hidden p-2 hover:bg-gray-100 rounded-full transition-colors" 
        aria-label="Abrir menú"
      >
        <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      <div class="hidden md:block w-9"></div> <!-- Spacer para desktop -->
      <h1 class="text-lg font-medium text-gray-900">MiCartaPro</h1>
      <div class="flex items-center gap-3">
        <!-- Contador de créditos -->
        <div class="flex items-center gap-2 px-3 py-1.5 bg-blue-50 rounded-full border border-blue-200">
          <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {#if userCredits !== null}
            <span class="text-sm font-semibold text-blue-700">{userCredits}</span>
          {:else}
            <span class="text-sm text-blue-500">-</span>
          {/if}
        </div>
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
      class="flex-1 px-4 py-6 min-h-0 overflow-y-auto"
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
        <div class="grid grid-cols-2 gap-2 md:gap-3 w-full">
          <button 
            class="p-2.5 md:p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-xl md:rounded-2xl transition-all text-left group flex items-center gap-2 md:gap-3"
            onclick={() => handleButtonClick('whatsapp')}
          >
            <span class="text-xl md:text-2xl">📱</span>
            <span class="text-sm md:text-base font-normal text-gray-900">{$tStore.chat.updateWhatsApp}</span>
          </button>

          <button 
            class="p-2.5 md:p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-xl md:rounded-2xl transition-all text-left group flex items-center gap-2 md:gap-3"
            onclick={() => handleButtonClick('address')}
          >
            <span class="text-xl md:text-2xl">📍</span>
            <span class="text-sm md:text-base font-normal text-gray-900">{$tStore.chat.updateAddress}</span>
          </button>

          <button 
            class="p-2.5 md:p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-xl md:rounded-2xl transition-all text-left group flex items-center gap-2 md:gap-3"
            onclick={() => handleButtonClick('dishes')}
          >
            <span class="text-xl md:text-2xl">🍽️</span>
            <span class="text-sm md:text-base font-normal text-gray-900">{$tStore.chat.addDishes}</span>
          </button>

          <button 
            class="p-2.5 md:p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-xl md:rounded-2xl transition-all text-left group flex items-center gap-2 md:gap-3"
            onclick={() => handleButtonClick('desserts')}
          >
            <span class="text-xl md:text-2xl">🍰</span>
            <span class="text-sm md:text-base font-normal text-gray-900">{$tStore.chat.addDesserts}</span>
          </button>

          <button 
            class="p-2.5 md:p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-xl md:rounded-2xl transition-all text-left group flex items-center gap-2 md:gap-3"
            onclick={() => handleButtonClick('price')}
          >
            <span class="text-xl md:text-2xl">💰</span>
            <span class="text-sm md:text-base font-normal text-gray-900">{$tStore.chat.updatePrice}</span>
          </button>

          <button 
            class="p-2.5 md:p-4 bg-white hover:bg-gray-50 border border-gray-200 rounded-xl md:rounded-2xl transition-all text-left group flex items-center gap-2 md:gap-3"
            onclick={() => handleButtonClick('delete')}
          >
            <span class="text-xl md:text-2xl">🗑️</span>
            <span class="text-sm md:text-base font-normal text-gray-900">{$tStore.chat.deleteItem}</span>
          </button>
        </div>
      </div>
    {:else}
      <div class="max-w-3xl mx-auto space-y-6 pt-4">
        {#each messages as message (message.id)}
          <Message 
            message={message} 
            onExploreOptions={resetToOptions}
          />
        {/each}

        {#if isLoading || uploadingPhoto}
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
    <div class="border-t border-gray-200 bg-white flex-shrink-0 z-10 md:relative sticky bottom-0" style="bottom: env(safe-area-inset-bottom, 0px); background: white; position: -webkit-sticky; position: sticky;">
    <div class="max-w-3xl mx-auto px-4 py-3">
      {#if checkingMenu}
        <div class="flex items-center justify-center py-8">
          <div class="text-center">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 border-t-transparent mx-auto mb-4"></div>
            <p class="text-gray-600 text-sm">{$tStore.chat.checkingMenu}</p>
          </div>
        </div>
      {/if}
      
      <!-- Preview de foto pendiente -->
      {#if pendingPhotoUrl}
        <div class="mb-2 px-2">
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-3 flex items-start gap-3">
            <img 
              src={pendingPhotoUrl} 
              alt="Foto pendiente" 
              class="w-16 h-16 object-cover rounded-lg flex-shrink-0"
            />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-700 mb-1">Foto lista para enviar</p>
              <p class="text-xs text-gray-500 truncate">{pendingPhotoUrl}</p>
            </div>
            <button
              onclick={() => {
                pendingPhotoUrl = null
                messages = messages.filter(msg => !msg.isPreview)
              }}
              class="p-1 text-gray-400 hover:text-red-600 transition-colors flex-shrink-0"
              title="Eliminar foto"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      {/if}

      <ChatInput 
        bind:this={chatInputRef} 
        onSend={handleSendMessage} 
        disabled={isLoading || checkingMenu || uploadingPhoto} 
        onFocus={handleInputFocus} 
        onBlur={handleInputBlur}
        onTakePhoto={handleTakePhoto}
        onUploadPhoto={handleUploadPhoto}
        onPastePhoto={handlePastePhoto}
        hasPendingPhoto={!!pendingPhotoUrl}
      />
      
      <!-- Input oculto para subir archivo -->
      <input
        id="photo-file-input"
        type="file"
        accept="image/*"
        onchange={handleFileSelect}
        class="hidden"
        disabled={uploadingPhoto}
      />
      
      <!-- Sugerencias justo debajo del input -->
      {#if showExamples && currentExampleType}
        <div class="mt-3 px-2">
          <p class="text-xs text-gray-500 mb-2 md:text-xs">{$tStore.chat.examplesLabel}</p>
          <div class="flex flex-col gap-1.5 md:gap-2">
            {#each getExampleMessages() as example}
              <div class="text-xs md:text-sm font-normal text-gray-600 py-1 md:py-1.5">
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
    class="absolute inset-0 flex flex-col h-full bg-white transition-transform duration-300 ease-in-out {(showPreview && !showSlugModal && !showSubscriptionPromo && !showCreditsPromo) ? 'translate-x-0' : 'translate-x-full'}"
  >
    <!-- Header del Preview -->
    <header class="border-b border-gray-200 bg-white px-4 py-2 flex items-center justify-between flex-shrink-0 z-10">
      <button
        onclick={() => {
          linkWasCopied = false
          togglePreview()
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
      </div>
      <div class="w-9"></div> <!-- Spacer para centrar -->
    </header>

    <!-- Contenido del Preview -->
    <div class="flex-1 overflow-hidden iframe-container relative min-h-0">
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

    <!-- Botón flotante: "Usar este menú" o "Compartir" según el estado -->
    {#if menuUrl}
      {@const pendingVersion = messages.find(msg => msg.pendingVersionId && !versionActivated)}
      {#if pendingVersion}
        {console.log('🔍 Botón "Usar este menú" visible - pendingVersion:', { id: pendingVersion.id, pendingVersionId: pendingVersion.pendingVersionId, versionActivated })}
        <!-- Debug: {console.log('🔍 Botón "Usar este menú" - pendingVersion:', pendingVersion)} -->
        <!-- Botón "Usar este menú" cuando hay versión pendiente (desktop) -->
        <div class="hidden md:block bg-white border-t border-gray-200 px-4 py-3 flex-shrink-0 safe-area-inset-bottom shadow-lg">
          <button
            onclick={() => {
              console.log('🔘 Click en "Usar este menú" - pendingVersionId:', pendingVersion.pendingVersionId)
              handleUseThisMenu(pendingVersion.pendingVersionId!)
            }}
            disabled={isActivatingVersion}
            class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-3"
            title={$tStore.chat.useThisMenu || 'Usar este menú'}
          >
            {#if isActivatingVersion}
              <div class="animate-spin rounded-full h-6 w-6 border-2 border-white border-t-transparent"></div>
              <span>{$tStore.chat.activating || 'Activando...'}</span>
            {:else}
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>{$tStore.chat.useThisMenu || 'Usar este menú'}</span>
            {/if}
          </button>
        </div>
      {:else}
        <!-- Botón "Compartir" cuando no hay versión pendiente (desktop) -->
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
    {/if}
  </div>

<!-- Botón flotante: "Usar este menú" o "Compartir" según el estado (móviles) -->
{#if menuUrl && showPreview && !showSlugModal && !showSubscriptionPromo && !showCreditsPromo}
  {@const pendingVersion = messages.find(msg => msg.pendingVersionId && !versionActivated)}
      {#if pendingVersion}
        <!-- Botón "Usar este menú" cuando hay versión pendiente (móvil) -->
        <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-4 py-3 z-[100] safe-area-inset-bottom shadow-2xl md:hidden">
          <button
            onclick={() => {
              console.log('🔘 Click en "Usar este menú" (móvil) - pendingVersionId:', pendingVersion.pendingVersionId)
              handleUseThisMenu(pendingVersion.pendingVersionId!)
            }}
        disabled={isActivatingVersion}
        class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-3"
        title={$tStore.chat.useThisMenu || 'Usar este menú'}
      >
        {#if isActivatingVersion}
          <div class="animate-spin rounded-full h-6 w-6 border-2 border-white border-t-transparent"></div>
          <span>{$tStore.chat.activating || 'Activando...'}</span>
        {:else}
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <span>{$tStore.chat.useThisMenu || 'Usar este menú'}</span>
        {/if}
      </button>
    </div>
  {:else}
    <!-- Botón "Compartir" cuando no hay versión pendiente (móvil) -->
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
{/if}

  <!-- Vista de Crear Slug (se muestra cuando showSlugModal es true) -->
  <div 
    class="absolute inset-0 flex flex-col h-full bg-white transition-transform duration-300 ease-in-out {(showSlugModal && !showSubscriptionPromo) ? 'translate-x-0' : 'translate-x-full'}"
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

<!-- Modal de confirmación para descartar menú -->
{#if showDiscardModal}
  <div 
    class="fixed inset-0 bg-black/50 backdrop-blur-sm z-[70] flex items-center justify-center p-4"
    onclick={cancelDiscard}
    role="dialog"
    aria-modal="true"
  >
    <div 
      class="relative bg-white rounded-xl md:rounded-2xl shadow-2xl w-full max-w-md"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="p-6 border-b border-gray-200">
        <h3 class="text-xl font-bold text-gray-900">{$tStore.chat.discardMenuTitle}</h3>
      </div>
      
      <!-- Content -->
      <div class="p-6">
        <p class="text-gray-700 mb-6">{$tStore.chat.discardMenuMessage}</p>
        
        <!-- Action buttons -->
        <div class="flex gap-3">
          <button
            onclick={cancelDiscard}
            class="flex-1 px-4 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors font-medium"
          >
            {$tStore.chat.discardMenuCancel}
          </button>
          <button
            onclick={confirmDiscard}
            class="flex-1 px-4 py-3 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors font-medium"
          >
            {$tStore.chat.discardMenuConfirm}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  #messages-container {
    scroll-behavior: smooth;
    -webkit-overflow-scrolling: touch;
    overflow-y: auto;
  }
  
  /* Safe area para dispositivos con notch */
  .safe-area-inset-bottom {
    padding-bottom: env(safe-area-inset-bottom, 0.75rem);
  }
  
  /* Asegurar que el input sea visible en móviles */
  @media (max-width: 768px) {
    /* El padding se maneja con inline style para considerar safe areas */
    #messages-container {
      -webkit-overflow-scrolling: touch;
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

  /* Modal de cámara */
  .camera-modal {
    z-index: 9999;
  }
</style>

<!-- Modal de cámara (si se quiere usar cámara directa en desktop) -->
{#if showCamera}
  <div class="fixed inset-0 bg-black z-50 flex items-center justify-center camera-modal">
    <div class="relative w-full max-w-2xl aspect-video bg-black">
      <video
        bind:this={videoRef}
        autoplay
        playsinline
        muted
        class="w-full h-full object-cover"
        aria-label="Vista previa de la cámara"
      ></video>
      
      <!-- Controles de cámara -->
      <div class="absolute bottom-4 left-0 right-0 flex justify-center items-center gap-6 z-10">
        <button
          onclick={stopCamera}
          class="px-4 py-2 bg-red-600/90 backdrop-blur-sm text-white rounded-full font-semibold hover:bg-red-700 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          Cancelar
        </button>
        
        <button
          onclick={capturePhoto}
          class="w-20 h-20 bg-white rounded-full border-4 border-gray-200 hover:border-gray-300 transition-colors flex items-center justify-center shadow-2xl"
          aria-label="Tomar foto"
        >
          <div class="w-16 h-16 bg-white rounded-full border-2 border-gray-400"></div>
        </button>
        
        <div class="w-20"></div> <!-- Spacer para centrar el botón de captura -->
      </div>
    </div>
  </div>
{/if}

  <!-- Vista Promocional de Suscripción (se muestra cuando showSubscriptionPromo es true) -->
  {#if showSubscriptionPromo}
  <div 
    class="absolute inset-0 flex flex-col h-full bg-white transition-transform duration-300 ease-in-out z-30 translate-x-0"
  >
    <!-- Header -->
    <header class="border-b border-gray-200 bg-white px-4 py-3 flex items-center justify-between flex-shrink-0">
      <button
        onclick={handleBackToEdit}
        class="p-2 hover:bg-gray-100 rounded-full transition-colors"
        aria-label="Volver"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h2 class="text-lg font-semibold text-gray-900 flex-1 text-center">
        🎉 Tu carta está lista
      </h2>
      <div class="w-9"></div> <!-- Spacer para centrar -->
    </header>

    <!-- Contenido -->
    <div class="flex-1 overflow-y-auto px-4 py-6">
      <div class="max-w-lg mx-auto">
        <!-- Título -->
        <h2 class="text-2xl md:text-3xl font-bold text-gray-900 mb-4 text-center">
          🎉 Tu carta está lista para compartirse
        </h2>

        <!-- Texto principal -->
        <div class="space-y-4 mb-6">
          <p class="text-gray-700 text-base leading-relaxed">
            Estás a un paso de publicar tu carta digital y empezar a recibir pedidos.
          </p>
          
          <p class="text-gray-700 text-base leading-relaxed">
            Normalmente MiCartaPro cuesta <span class="font-semibold text-gray-900">USD 15</span>, pero por tiempo limitado puedes activar tu plan por solo <span class="font-bold text-green-600 text-lg">USD 3.5</span>.
          </p>
        </div>

        <!-- Lista de beneficios -->
        <div class="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl p-5 mb-6 border border-blue-100">
          <p class="font-semibold text-gray-900 mb-3 text-sm">Activa tu plan ahora y desbloquea:</p>
          <ul class="space-y-2">
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-gray-700 text-sm">Compartir tu enlace público</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-gray-700 text-sm">Código QR listo para imprimir</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-gray-700 text-sm">Recepción de pedidos por WhatsApp</span>
            </li>
            <li class="flex items-start gap-2">
              <svg class="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-gray-700 text-sm">Tu carta siempre disponible online</span>
            </li>
          </ul>
        </div>

        <!-- Botones -->
        <div class="flex flex-col gap-3">
          <!-- CTA Principal -->
          <button
            onclick={handleActivatePlan}
            class="w-full px-6 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white rounded-xl shadow-lg hover:shadow-xl transition-all font-semibold text-base flex items-center justify-center gap-2"
          >
            <span>👉</span>
            <span>Activar plan por USD 3.5</span>
          </button>

          <!-- CTA Secundario -->
          <button
            onclick={handleBackToEdit}
            class="w-full px-6 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-xl transition-all font-medium text-sm"
          >
            Volver a editar mi carta
          </button>
        </div>
      </div>
    </div>
  </div>
  {/if}

