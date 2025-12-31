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
      upgradeModalContribution: 'Tu pago contribuye al crecimiento y mejora continua de la plataforma.',
      upgradeToPro: 'Actualizar a Pro',
      continueWithoutPayment: 'Continuar sin pagar',
      close: 'Cerrar',
      updateAddressMessage: 'Actualiza la dirección de mi carta',
      addDishesMessage: 'Agrega uno o varios platos a mi carta',
      addDessertsMessage: 'Agrega uno o varios postres a mi carta',
      updatePriceMessage: 'Actualiza el precio de uno de mis platos',
      deleteItemMessage: 'Elimina un item de mi carta',
      updateWhatsAppMessage: 'Actualiza mi número de WhatsApp'
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
      upgradeModalContribution: 'Seu pagamento contribui para o crescimento e melhoria contínua da plataforma.',
      upgradeToPro: 'Atualizar para Pro',
      continueWithoutPayment: 'Continuar sem pagar',
      close: 'Fechar',
      updateAddressMessage: 'Atualize o endereço da minha carta',
      addDishesMessage: 'Adicione um ou vários pratos à minha carta',
      addDessertsMessage: 'Adicione uma ou várias sobremesas à minha carta',
      updatePriceMessage: 'Atualize o preço de um dos meus pratos',
      deleteItemMessage: 'Exclua um item da minha carta',
      updateWhatsAppMessage: 'Atualize meu número do WhatsApp',
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
      upgradeModalContribution: 'Your payment contributes to the growth and continuous improvement of the platform.',
      upgradeToPro: 'Upgrade to Pro',
      continueWithoutPayment: 'Continue without payment',
      close: 'Close',
      updateAddressMessage: 'Update my menu address',
      addDishesMessage: 'Add one or more dishes to my menu',
      addDessertsMessage: 'Add one or more desserts to my menu',
      updatePriceMessage: 'Update the price of one of my dishes',
      deleteItemMessage: 'Delete an item from my menu',
      updateWhatsAppMessage: 'Update my WhatsApp number'
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

