export type Language = 'EN' | 'ES' | 'PT'

export interface Translations {
  // App.svelte
  app: {
    loading: string
    notAuthenticated: string
    pleaseSignIn: string
    signInLink: string
  }
  
  // MenuChat.svelte
  chat: {
    greeting: string
    greetingQuestion: string
    welcomeTitle: string
    welcomeSubtitle: string
    previewButton: string
    previewTitle: string
    copyLink: string
    linkCopied: string
    shareLink: string
    loadingPreview: string
    placeholder: string
    sendButton: string
    
    // Quick actions
    updateAddress: string
    addDishes: string
    addDesserts: string
    updatePrice: string
    deleteItem: string
    updateWhatsApp: string
    
    // Messages
    errorNoSession: string
    errorNoMenu: string
    checkingMenu: string
    errorProcessing: string
    errorPolling: string
    successUpdated: string
    exploreOptions: string
    examplesLabel: string
    trialDaysRemaining: string
    trialDaysRemainingLabel: string
    upgradeModalTitle: string
    upgradeModalMessage: string
    upgradeModalBenefits: string
    upgradeModalBenefit1: string
    upgradeModalBenefit2: string
    upgradeModalBenefit3: string
    upgradeModalContribution: string
    upgradeToPro: string
    continueWithoutPayment: string
    close: string
    
    // Quick action messages (sent to backend)
    updateAddressMessage: string
    addDishesMessage: string
    addDessertsMessage: string
    updatePriceMessage: string
    deleteItemMessage: string
    updateWhatsAppMessage: string
      useThisMenu: string
      activating: string
      discardMenuTitle: string
      discardMenuMessage: string
      discardMenuConfirm: string
      discardMenuCancel: string
    }
  
  // MenuPreview.svelte
  preview: {
    loading: string
    error: string
    errorNoMenu: string
    errorLoading: string
    title: string
    linkLabel: string
    copyButton: string
    copied: string
    openButton: string
  }
  
  // PaymentSuccess.svelte
  paymentSuccess: {
    title: string
    subtitle: string
    message: string
    redirecting: string
    goToConsole: string
  }
  
  // Sidebar.svelte
  sidebar: {
    aiAssistant: string
    history: string
    gallery: string
    qrCode: string
    orders: string
    kitchen: string
    myPlan: string
    signOut: string
    closeMenu: string
    confirmSignOut: string
    errorNoSession: string
    errorNoActiveSubscription: string
    errorGettingPortal: string
    errorNoPortalUrl: string
    errorAccessingPortal: string
  }
  
  // ProcessingPreview.svelte
  processing: {
    understandingInstructions: string
    creatingCatalog: string
    preparingSite: string
    validatingImages: string
    improvingImages: string
    finalizing: string
    pleaseWait: string
    preparingInitialSetup: string // Mensaje específico para la creación inicial del menú
  }
  // MenuOrders.svelte
  orders?: {
    noSession: string
    noMenu: string
    empty: string
    subtitle: string
    delivery: string
    pickup: string
    viewAsPaper?: string
    print?: string
    printThermal?: string
    eventPayload: string
    itemsToPrepare: string
    forTime: string
    remainingIn: string
    late: string
    itemsCount: string
    markAsReady: string
    kitchenMode: string
    exitKitchenMode: string
    statusPending: string
    statusPreparing: string
    statusDone: string
    tabPending: string
    tabPreparing: string
    tabDone: string
    statusGeneralLabel: string
    readyToDeliver: string
    startPreparing: string
    reload: string
    filterAll: string
    filterKitchen: string
    filterBar: string
    showQR: string
    emptyForStation: string
    startToReady?: string
    deliver?: string
  }
}

export const translations: Record<Language, Translations> = {
  ES: {
    app: {
      loading: 'Cargando...',
      notAuthenticated: 'No autenticado',
      pleaseSignIn: 'Por favor, inicia sesión en',
      signInLink: 'micartapro-auth-ui'
    },
    chat: {
      greeting: 'Hola',
      greetingQuestion: '¿Qué quieres en tu carta?',
      welcomeTitle: '¿En qué puedo ayudarte?',
      welcomeSubtitle: 'Escribe tu menú y precios, y yo armaré tu carta digital',
      previewButton: 'Vista Previa',
      previewTitle: 'Vista Previa de tu Carta',
      copyLink: 'Copiar enlace',
      linkCopied: '¡Copiado!',
      shareLink: 'Compartir enlace',
      loadingPreview: 'Cargando vista previa...',
      placeholder: 'Escribe tu solicitud aquí...',
      sendButton: 'Enviar',
      updateAddress: 'Actualiza la dirección de mi carta',
      addDishes: 'Agrega uno o varios platos a mi carta',
      addDesserts: 'Agrega uno o varios postres a mi carta',
      updatePrice: 'Actualiza el precio de uno de mis platos',
      deleteItem: 'Elimina un item de mi carta',
      updateWhatsApp: 'Actualiza mi número de WhatsApp',
      errorNoSession: 'Error: No hay sesión activa. Por favor, inicia sesión nuevamente.',
      errorNoMenu: 'Error: No se encontró un menú. Por favor, crea un menú primero.',
      checkingMenu: 'Verificando que tu menú esté listo...',
      errorProcessing: 'Error al procesar tu mensaje: {message}',
      errorPolling: 'El mensaje fue procesado, pero hubo un problema al verificar la actualización: {message}',
      successUpdated: '¡Tu carta ha sido actualizada exitosamente! El menú se ha guardado con los cambios solicitados.',
      exploreOptions: 'Explorar opciones',
      examplesLabel: 'Ejemplos:',
      trialDaysRemaining: '{days} días restantes',
      trialDaysRemainingLabel: 'Días de prueba restantes',
      upgradeModalTitle: '¡Actualiza a Pro!',
      upgradeModalMessage: 'Te quedan {days} días de prueba. Actualiza a Pro para seguir disfrutando de todos los beneficios.',
      upgradeModalBenefits: 'Con Pro obtienes: acceso ilimitado, sin límites de edición, soporte prioritario y más.',
      upgradeModalBenefit1: 'Acceso ilimitado a todas las funciones',
      upgradeModalBenefit2: 'Sin límites de edición',
      upgradeModalBenefit3: 'Soporte prioritario',
      upgradeModalContribution: 'Tu pago contribuye al crecimiento y mejora continua de la plataforma.',
      upgradeToPro: 'Actualizar a Pro - 14 días gratuitos',
      continueWithoutPayment: 'Continuar sin pagar',
      close: 'Cerrar',
      updateAddressMessage: 'Actualiza la dirección de mi carta',
      addDishesMessage: 'Agrega uno o varios platos a mi carta',
      addDessertsMessage: 'Agrega uno o varios postres a mi carta',
      updatePriceMessage: 'Actualiza el precio de uno de mis platos',
      deleteItemMessage: 'Elimina un item de mi carta',
      updateWhatsAppMessage: 'Actualiza mi número de WhatsApp',
      useThisMenu: 'Usar este menú',
      activating: 'Activando...',
      discardMenuTitle: '¿Descartar este menú?',
      discardMenuMessage: 'Tienes un menú pendiente de aceptar. Si retrocedes, este menú no se activará.',
      discardMenuConfirm: 'Sí, descartar',
      discardMenuCancel: 'Cancelar'
    },
    preview: {
      loading: 'Cargando tu carta...',
      error: 'Error',
      errorNoMenu: 'No se encontró un menú para este usuario. Crea uno primero.',
      errorLoading: 'Error al cargar el menú',
      title: 'Tu Carta Digital',
      linkLabel: 'Enlace de tu carta:',
      copyButton: 'Copiar',
      copied: '¡Copiado!',
      openButton: 'Abrir'
    },
    paymentSuccess: {
      title: '¡Pago Exitoso! 🎉',
      subtitle: '¡Bienvenido a Pro!',
      message: 'Tu suscripción Pro ha sido activada. Ahora puedes disfrutar de todas las funciones premium.',
      redirecting: 'Redirigiendo a la consola...',
      goToConsole: 'Ir a la consola'
    },
    sidebar: {
      aiAssistant: 'Asistente IA',
      history: 'Historial',
      gallery: 'Galería',
      qrCode: 'Código QR',
      orders: 'Órdenes',
      kitchen: 'Cocina',
      myPlan: 'Mi Plan',
      signOut: 'Cerrar sesión',
      closeMenu: 'Cerrar menú',
      confirmSignOut: '¿Estás seguro de que deseas cerrar sesión?',
      errorNoSession: 'No hay sesión activa. Por favor, inicia sesión nuevamente.',
      errorNoActiveSubscription: 'No se encontró una suscripción activa. Por favor, activa tu plan primero.',
      errorGettingPortal: 'Error al obtener el portal del consumidor. Por favor, intenta de nuevo.',
      errorNoPortalUrl: 'No se recibió la URL del portal. Por favor, intenta de nuevo.',
      errorAccessingPortal: 'Error al acceder al portal del consumidor. Por favor, intenta de nuevo.'
    },
    processing: {
      understandingInstructions: 'Entendiendo tus instrucciones...',
      creatingCatalog: 'Creando tu catálogo...',
      preparingSite: 'Preparando tu sitio...',
      validatingImages: 'Validando si necesitas imágenes...',
      improvingImages: 'Mejorando tus imágenes...',
      finalizing: 'Finalizando tu menú...',
      pleaseWait: 'Por favor espera, esto puede tomar unos momentos',
      preparingInitialSetup: 'Preparando todo para que puedas crear tu catálogo...'
    },
    orders: {
      noSession: 'No hay sesión activa',
      noMenu: 'No se encontró un menú',
      empty: 'No hay órdenes aún.',
      subtitle: 'Ordenado por hora comprometida. Vista orientada a cocina.',
      delivery: 'Envío',
      pickup: 'Retiro',
      viewAsPaper: 'Ver como hoja',
      print: 'Imprimir',
      printThermal: 'Imprimir en térmica',
      eventPayload: 'event_payload',
      itemsToPrepare: 'Qué preparar',
      forTime: 'Para',
      remainingIn: 'En {min} min',
      late: 'Atrasado {min} min',
      itemsCount: '{count} ítems',
      markAsReady: 'LISTO',
      kitchenMode: 'Modo full',
      exitKitchenMode: 'Salir modo full',
      statusPending: 'Pendiente',
      statusPreparing: 'En preparación',
      statusDone: 'LISTO',
      tabPending: 'Pendientes',
      tabPreparing: 'En preparación',
      tabDone: 'Listos',
      statusGeneralLabel: 'Estado general',
      readyToDeliver: 'Listo para entregar',
      startPreparing: 'Iniciar preparación',
      reload: 'Recargar',
      filterAll: 'Entrega',
      filterKitchen: 'Cocina',
      filterBar: 'Barra',
      showQR: 'QR',
      emptyForStation: 'No hay órdenes para esta estación.',
      startToReady: 'INICIAR → LISTO',
      deliver: 'ENTREGAR'
    }
  },
  PT: {
    app: {
      loading: 'Carregando...',
      notAuthenticated: 'Não autenticado',
      pleaseSignIn: 'Por favor, faça login em',
      signInLink: 'micartapro-auth-ui'
    },
    chat: {
      greeting: 'Olá',
      greetingQuestion: 'O que você quer no seu cardápio?',
      welcomeTitle: 'Como posso ajudá-lo?',
      welcomeSubtitle: 'Escreva seu cardápio e preços, e eu criarei sua carta digital',
      previewButton: 'Visualizar',
      previewTitle: 'Visualização da sua Carta',
      copyLink: 'Copiar link',
      linkCopied: 'Copiado!',
      shareLink: 'Compartilhar link',
      loadingPreview: 'Carregando visualização...',
      placeholder: 'Escreva sua solicitação aqui...',
      sendButton: 'Enviar',
      updateAddress: 'Atualize o endereço da minha carta',
      addDishes: 'Adicione um ou vários pratos à minha carta',
      addDesserts: 'Adicione uma ou várias sobremesas à minha carta',
      updatePrice: 'Atualize o preço de um dos meus pratos',
      deleteItem: 'Exclua um item da minha carta',
      updateWhatsApp: 'Atualize meu número do WhatsApp',
      createMenu: 'Criar cardápio',
      organizeDishes: 'Organizar pratos',
      viewPrices: 'Ver preços',
      moreOptions: 'Mais',
      welcomeMessage1: 'Olá! 👋 Sou seu assistente para criar cardápios digitais. Posso ajudá-lo a criar sua carta de forma profissional.',
      welcomeMessage2: 'Simplesmente escreva seu cardápio e preços, e eu me encarregarei de organizá-los e formatá-los para criar uma carta atraente.',
      errorNoSession: 'Erro: Não há sessão ativa. Por favor, faça login novamente.',
      errorNoMenu: 'Erro: Nenhum cardápio encontrado. Por favor, crie um primeiro.',
      checkingMenu: 'Verificando se seu cardápio está pronto...',
      errorProcessing: 'Erro ao processar sua mensagem: {message}',
      errorPolling: 'A mensagem foi processada, mas houve um problema ao verificar a atualização: {message}',
      successUpdated: 'Sua carta foi atualizada com sucesso! O cardápio foi salvo com as alterações solicitadas.',
      exploreOptions: 'Explorar opções',
      examplesLabel: 'Exemplos:',
      trialDaysRemaining: '{days} dias restantes',
      trialDaysRemainingLabel: 'Dias de teste restantes',
      upgradeModalTitle: 'Atualize para Pro!',
      upgradeModalMessage: 'Você tem {days} dias de teste restantes. Atualize para Pro para continuar aproveitando todos os benefícios.',
      upgradeModalBenefits: 'Com Pro você obtém: acesso ilimitado, sem limites de edição, suporte prioritário e muito mais.',
      upgradeModalBenefit1: 'Acesso ilimitado a todas as funções',
      upgradeModalBenefit2: 'Sem limites de edição',
      upgradeModalBenefit3: 'Suporte prioritário',
      upgradeModalContribution: 'Seu pagamento contribui para o crescimento e melhoria contínua da plataforma.',
      upgradeToPro: 'Atualizar para Pro - 14 dias grátis',
      continueWithoutPayment: 'Continuar sem pagar',
      close: 'Fechar',
      updateAddressMessage: 'Atualize o endereço da minha carta',
      addDishesMessage: 'Adicione um ou vários pratos à minha carta',
      addDessertsMessage: 'Adicione uma ou várias sobremesas à minha carta',
      updatePriceMessage: 'Atualize o preço de um dos meus pratos',
      deleteItemMessage: 'Exclua um item da minha carta',
      updateWhatsAppMessage: 'Atualize meu número do WhatsApp',
      useThisMenu: 'Usar este cardápio',
      activating: 'Ativando...',
      discardMenuTitle: 'Descartar este cardápio?',
      discardMenuMessage: 'Você tem um cardápio pendente de aceitar. Se voltar, este cardápio não será ativado.',
      discardMenuConfirm: 'Sim, descartar',
      discardMenuCancel: 'Cancelar',
      createMenuMessage: 'Quero criar um cardápio para um restaurante',
      organizeDishesMessage: 'Preciso de ajuda para organizar meus pratos',
      viewPricesMessage: 'Como funciona o sistema de preços?',
      moreOptionsMessage: 'Mais opções'
    },
    preview: {
      loading: 'Carregando sua carta...',
      error: 'Erro',
      errorNoMenu: 'Nenhum cardápio encontrado para este usuário. Crie um primeiro.',
      errorLoading: 'Erro ao carregar o cardápio',
      title: 'Sua Carta Digital',
      linkLabel: 'Link da sua carta:',
      copyButton: 'Copiar',
      copied: 'Copiado!',
      openButton: 'Abrir'
    },
    paymentSuccess: {
      title: 'Pagamento Bem-sucedido! 🎉',
      subtitle: 'Bem-vindo ao Pro!',
      message: 'Sua assinatura Pro foi ativada. Agora você pode desfrutar de todos os recursos premium.',
      redirecting: 'Redirecionando para o console...',
      goToConsole: 'Ir para o console'
    },
    sidebar: {
      aiAssistant: 'Assistente IA',
      history: 'Histórico',
      gallery: 'Galeria',
      qrCode: 'Código QR',
      orders: 'Pedidos',
      myPlan: 'Meu Plano',
      signOut: 'Sair',
      closeMenu: 'Fechar menu',
      confirmSignOut: 'Tem certeza de que deseja sair?',
      errorNoSession: 'Não há sessão ativa. Por favor, faça login novamente.',
      errorNoActiveSubscription: 'Nenhuma assinatura ativa encontrada. Por favor, ative seu plano primeiro.',
      errorGettingPortal: 'Erro ao obter o portal do consumidor. Por favor, tente novamente.',
      errorNoPortalUrl: 'A URL do portal não foi recebida. Por favor, tente novamente.',
      errorAccessingPortal: 'Erro ao acessar o portal do consumidor. Por favor, tente novamente.'
    },
    processing: {
      understandingInstructions: 'Entendendo suas instruções...',
      creatingCatalog: 'Criando seu catálogo...',
      preparingSite: 'Preparando seu site...',
      validatingImages: 'Validando se você precisa de imagens...',
      improvingImages: 'Melhorando suas imagens...',
      finalizing: 'Finalizando seu cardápio...',
      pleaseWait: 'Por favor aguarde, isso pode levar alguns momentos',
      preparingInitialSetup: 'Preparando tudo para que você possa criar seu catálogo...'
    },
    orders: {
      noSession: 'Não há sessão ativa',
      noMenu: 'Nenhum cardápio encontrado',
      empty: 'Ainda não há pedidos.',
      subtitle: 'Ordenado por hora comprometida. Vista orientada à cozinha.',
      delivery: 'Entrega',
      pickup: 'Retirada',
      viewAsPaper: 'Ver como folha',
      print: 'Imprimir',
      printThermal: 'Imprimir em térmica',
      eventPayload: 'event_payload',
      itemsToPrepare: 'O que preparar',
      forTime: 'Para',
      remainingIn: 'Em {min} min',
      late: 'Atrasado {min} min',
      itemsCount: '{count} itens',
      markAsReady: 'PRONTO',
      kitchenMode: 'Modo full',
      exitKitchenMode: 'Sair do modo full',
      statusPending: 'Pendente',
      statusPreparing: 'Em preparação',
      statusDone: 'PRONTO',
      tabPending: 'Pendentes',
      tabPreparing: 'Em preparação',
      tabDone: 'Prontos',
      statusGeneralLabel: 'Estado geral',
      readyToDeliver: 'Pronto para entregar',
      startPreparing: 'Iniciar preparação',
      reload: 'Atualizar',
      filterAll: 'Entrega',
      filterKitchen: 'Cozinha',
      filterBar: 'Bar',
      showQR: 'QR',
      emptyForStation: 'Nenhum pedido para esta estação.',
      startToReady: 'INICIAR → PRONTO',
      deliver: 'ENTREGAR'
    }
  },
  EN: {
    app: {
      loading: 'Loading...',
      notAuthenticated: 'Not authenticated',
      pleaseSignIn: 'Please sign in at',
      signInLink: 'micartapro-auth-ui'
    },
    chat: {
      greeting: 'Hello',
      greetingQuestion: 'What do you want in your menu?',
      welcomeTitle: 'How can I help you?',
      welcomeSubtitle: 'Write your menu and prices, and I will create your digital menu',
      previewButton: 'Preview',
      previewTitle: 'Preview of your Menu',
      copyLink: 'Copy link',
      linkCopied: 'Copied!',
      shareLink: 'Share link',
      loadingPreview: 'Loading preview...',
      placeholder: 'Write your request here...',
      sendButton: 'Send',
      updateAddress: 'Update my menu address',
      addDishes: 'Add one or more dishes to my menu',
      addDesserts: 'Add one or more desserts to my menu',
      updatePrice: 'Update the price of one of my dishes',
      deleteItem: 'Delete an item from my menu',
      updateWhatsApp: 'Update my WhatsApp number',
      errorNoSession: 'Error: No active session. Please sign in again.',
      errorNoMenu: 'Error: No menu found. Please create one first.',
      checkingMenu: 'Checking that your menu is ready...',
      errorProcessing: 'Error processing your message: {message}',
      errorPolling: 'The message was processed, but there was a problem verifying the update: {message}',
      successUpdated: 'Your menu has been successfully updated! The menu has been saved with the requested changes.',
      exploreOptions: 'Explore options',
      examplesLabel: 'Examples:',
      trialDaysRemaining: '{days} days remaining',
      trialDaysRemainingLabel: 'Trial days remaining',
      upgradeModalTitle: 'Upgrade to Pro!',
      upgradeModalMessage: 'You have {days} trial days remaining. Upgrade to Pro to continue enjoying all the benefits.',
      upgradeModalBenefits: 'With Pro you get: unlimited access, no editing limits, priority support, and more.',
      upgradeModalBenefit1: 'Unlimited access to all features',
      upgradeModalBenefit2: 'No editing limits',
      upgradeModalBenefit3: 'Priority support',
      upgradeModalContribution: 'Your payment contributes to the growth and continuous improvement of the platform.',
      upgradeToPro: 'Upgrade to Pro - 14 days free',
      continueWithoutPayment: 'Continue without payment',
      close: 'Close',
      updateAddressMessage: 'Update my menu address',
      addDishesMessage: 'Add one or more dishes to my menu',
      addDessertsMessage: 'Add one or more desserts to my menu',
      updatePriceMessage: 'Update the price of one of my dishes',
      deleteItemMessage: 'Delete an item from my menu',
      updateWhatsAppMessage: 'Update my WhatsApp number',
      useThisMenu: 'Use this menu',
      activating: 'Activating...',
      discardMenuTitle: 'Discard this menu?',
      discardMenuMessage: 'You have a menu pending acceptance. If you go back, this menu will not be activated.',
      discardMenuConfirm: 'Yes, discard',
      discardMenuCancel: 'Cancel'
    },
    preview: {
      loading: 'Loading your menu...',
      error: 'Error',
      errorNoMenu: 'No menu found for this user. Create one first.',
      errorLoading: 'Error loading menu',
      title: 'Your Digital Menu',
      linkLabel: 'Your menu link:',
      copyButton: 'Copy',
      copied: 'Copied!',
      openButton: 'Open'
    },
    paymentSuccess: {
      title: 'Payment Successful! 🎉',
      subtitle: 'Welcome to Pro!',
      message: 'Your Pro subscription has been activated. You can now enjoy all premium features.',
      redirecting: 'Redirecting to console...',
      goToConsole: 'Go to console'
    },
    sidebar: {
      aiAssistant: 'AI Assistant',
      history: 'History',
      gallery: 'Gallery',
      qrCode: 'QR Code',
      orders: 'Orders',
      kitchen: 'Kitchen',
      myPlan: 'My Plan',
      signOut: 'Sign Out',
      closeMenu: 'Close menu',
      confirmSignOut: 'Are you sure you want to sign out?',
      errorNoSession: 'No active session. Please sign in again.',
      errorNoActiveSubscription: 'No active subscription found. Please activate your plan first.',
      errorGettingPortal: 'Error getting customer portal. Please try again.',
      errorNoPortalUrl: 'Portal URL was not received. Please try again.',
      errorAccessingPortal: 'Error accessing customer portal. Please try again.'
    },
    processing: {
      understandingInstructions: 'Understanding your instructions...',
      creatingCatalog: 'Creating your catalog...',
      preparingSite: 'Preparing your site...',
      validatingImages: 'Validating if you need images...',
      improvingImages: 'Improving your images...',
      finalizing: 'Finalizing your menu...',
      pleaseWait: 'Please wait, this may take a few moments',
      preparingInitialSetup: 'Preparing everything so you can create your catalog...'
    },
    orders: {
      noSession: 'No active session',
      noMenu: 'No menu found',
      empty: 'No orders yet.',
      subtitle: 'Sorted by requested time. Kitchen-oriented view.',
      delivery: 'Delivery',
      pickup: 'Pickup',
      viewAsPaper: 'View as sheet',
      print: 'Print',
      printThermal: 'Print to thermal',
      eventPayload: 'event_payload',
      itemsToPrepare: 'What to prepare',
      forTime: 'For',
      remainingIn: 'In {min} min',
      late: '{min} min late',
      itemsCount: '{count} items',
      markAsReady: 'READY',
      kitchenMode: 'Full mode',
      exitKitchenMode: 'Exit full mode',
      statusPending: 'Pending',
      statusPreparing: 'Preparing',
      statusDone: 'READY',
      tabPending: 'Pending',
      tabPreparing: 'Preparing',
      tabDone: 'Ready',
      statusGeneralLabel: 'Overall status',
      readyToDeliver: 'Ready for delivery',
      startPreparing: 'Start preparation',
      reload: 'Refresh',
      filterAll: 'Delivery',
      filterKitchen: 'Kitchen',
      filterBar: 'Bar',
      showQR: 'QR',
      emptyForStation: 'No orders for this station.',
      startToReady: 'START → READY',
      deliver: 'DELIVER'
    }
  }
}

export const languageNames: Record<Language, string> = {
  ES: 'Español',
  PT: 'Português',
  EN: 'English'
}

export const languageFlags: Record<Language, string> = {
  ES: '🇪🇸',
  PT: '🇧🇷',
  EN: '🇺🇸'
}

