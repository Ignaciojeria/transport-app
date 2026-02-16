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
      templateLabel: string
      templateHero: string
      templateModern: string
      useThisDesign: string
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
    negocios?: string
    aiAssistant: string
    history: string
    gallery: string
    qrCode: string
    orders: string
    jornada: string
    reportes: string
    cost: string
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
  
  // CostView.svelte
  cost?: {
    title: string
    subtitle: string
    unit: string
    baseUnit: string
    priceSale: string
    cost: string
    save: string
    saved: string
    reload: string
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
    emptyCancelled?: string
    startToReady?: string
    deliver?: string
    dispatch?: string
    delivered?: string
    dispatched?: string
    cancelled?: string
    cancelOrder?: string
    cancelModalTitle?: string
    cancelModalSubtitle?: string
    cancelModalReasonLabel?: string
    cancelReasons?: Record<string, string>
    cancelModalCommentLabel?: string
    cancelModalCommentPlaceholder?: string
    cancelModalBack?: string
    cancelModalConfirm?: string
    viewVertical?: string
    viewThreeColumns?: string
    // Vista Pendientes (órdenes sin jornada)
    pendingView?: string
    pendingFrom?: string
    pendingTo?: string
    pendingSelectAll?: string
    pendingSelected?: string
    pendingAssignToJourney?: string
    pendingCancelOrders?: string
    pendingClearSelection?: string
    pendingViewDetails?: string
    pendingAssign?: string
    pendingEmpty?: string
    pendingNoActiveJourney?: string
  }
  jornada?: {
    title: string
    active: string
    date: string
    openedSince: string
    summary: string
    totalOrders: string
    delivered: string
    cancelled: string
    pending: string
    closeWorkday: string
    closeModalTitle: string
    closeModalMessage: string
    closeModalCancel: string
    closeModalCancelPending?: string
    closeModalCancelPendingHint?: string
    closeModalKeepForNext?: string
    closeModalKeepForNextHint?: string
    closeModalConfirm: string
    closing: string
    success: string
    comingSoon: string
    noMenu: string
    noSession: string
    errorLoading: string
    noActiveJourney?: string
    openJourney?: string
    openingJourney?: string
    errorCreatingJourney?: string
    openJourneyReason?: string
    reports?: string
    downloadExcel?: string
    reportGenerating?: string
    stats?: string
    totalRevenue?: string
    revenueConcreted?: string
    revenueDelivered?: string
    revenueDispatched?: string
    revenuePending?: string
    revenueCancelled?: string
    statsOrders?: string
    averageTicket?: string
    topByRevenue?: string
    topByQuantity?: string
    chartByRevenue?: string
    chartByProfit?: string
    chartByQuantity?: string
    topProducts?: string
    noStats?: string
    errorLoadingStats?: string
    prev?: string
    next?: string
    page?: string
    back?: string
    workdayReport?: string
  }
  negocios?: {
    title: string
    subtitle: string
    createNew: string
    creating: string
    yourBusinesses: string
    noBusinesses: string
    business: string
    active: string
    select: string
    errorLoading: string
    errorCreating: string
    slugLabel?: string
    slugPlaceholder?: string
    slugHint?: string
    create?: string
    cancel?: string
    newMenuTitle?: string
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
      discardMenuCancel: 'Cancelar',
      templateLabel: 'Estilo',
      templateHero: 'Hero',
      templateModern: 'Modern',
      useThisDesign: 'Usar este diseño'
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
    cost: {
      title: 'Costos del menú',
      subtitle: 'Configura el costo de cada ítem y acompañamiento. Se usa para márgenes y reportes.',
      unit: 'Unidad de venta',
      baseUnit: 'Base',
      priceSale: 'Precio venta',
      cost: 'Costo',
      save: 'Guardar',
      saved: 'Precio y costo guardados correctamente',
      reload: 'Recargar'
    },
    sidebar: {
      negocios: 'Negocios',
      aiAssistant: 'Asistente IA',
      history: 'Historial',
      gallery: 'Galería',
      qrCode: 'Código QR',
      orders: 'Kanban',
      jornada: 'Jornada',
      reportes: 'Reportes',
      cost: 'Costos',
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
      noSession: 'No Hay Sesión Activa',
      noMenu: 'No Se Encontró Un Menú',
      empty: 'No Hay Órdenes Aún.',
      subtitle: 'Ordenado Por Hora Comprometida. Vista Orientada A Cocina.',
      delivery: 'Envío',
      pickup: 'Retiro',
      viewAsPaper: 'Ver Como Hoja',
      print: 'Imprimir',
      printThermal: 'Imprimir En Térmica',
      eventPayload: 'event_payload',
      itemsToPrepare: 'Qué Preparar',
      forTime: 'Para',
      remainingIn: 'En {min} Min',
      late: 'Atrasado {min} Min',
      itemsCount: '{count} Ítems',
      markAsReady: 'Listo',
      kitchenMode: 'Modo Full',
      exitKitchenMode: 'Salir Modo Full',
      statusPending: 'Pendiente',
      statusPreparing: 'En Preparación',
      statusDone: 'Listo',
      tabPending: 'Pendientes',
      tabPreparing: 'En Preparación',
      tabDone: 'Listos',
      statusGeneralLabel: 'Estado General',
      readyToDeliver: 'Listo Para Entregar',
      startPreparing: 'Iniciar Preparación',
      reload: 'Recargar',
      filterAll: 'Entrega',
      filterKitchen: 'Cocina',
      filterBar: 'Barra',
      showQR: 'QR',
      emptyForStation: 'No Hay Órdenes Para Esta Estación.',
      emptyCancelled: 'No Hay Órdenes Canceladas.',
      startToReady: 'Iniciar → Listo',
      deliver: 'Entregar',
      dispatch: 'Despachar',
      delivered: 'Entregado',
      dispatched: 'Despachado',
      cancelled: 'Cancelado',
      cancelOrder: 'Cancelar Pedido',
      cancelModalTitle: 'Cancelar Pedido',
      cancelModalSubtitle: 'Esta Acción No Se Puede Deshacer.',
      cancelModalReasonLabel: 'Motivo (Elige Uno):',
      cancelReasons: {
        outOfStock: 'Falta De Stock',
        orderError: 'Error De Pedido',
        customerLeft: 'Cliente Se Fue',
        paymentIssue: 'Problema De Pago',
        other: 'Otro'
      },
      cancelModalCommentLabel: 'Comentario (Opcional):',
      cancelModalCommentPlaceholder: 'Ej: Cliente No Contestó...',
      cancelModalBack: 'Volver',
      cancelModalConfirm: 'Confirmar Cancelación',
      viewVertical: 'Vertical',
      viewThreeColumns: '3 Columnas',
      pendingView: 'Pendientes',
      pendingFrom: 'Desde',
      pendingTo: 'Hasta',
      pendingSelectAll: 'Seleccionar todas',
      pendingSelected: 'órdenes seleccionadas',
      pendingAssignToJourney: 'Asignar a jornada activa',
      pendingCancelOrders: 'Cancelar órdenes',
      pendingClearSelection: 'Limpiar selección',
      pendingViewDetails: 'Ver',
      pendingAssign: 'Asignar',
      pendingEmpty: 'No hay órdenes pendientes.',
      pendingNoActiveJourney: 'Abre una jornada para asignar órdenes.'
    },
    jornada: {
      title: 'Jornada',
      active: 'Jornada Activa',
      date: 'Fecha',
      openedSince: 'Abierta Desde',
      summary: 'Resumen Rápido',
      totalOrders: 'Órdenes Totales',
      delivered: 'Entregadas',
      cancelled: 'Canceladas',
      pending: 'Pendientes',
      closeWorkday: 'Cerrar Jornada',
      closeModalTitle: 'Cerrar Jornada',
      closeModalMessage: 'Estás Por Cerrar La Jornada Del {date}. ¿Qué Deseas Hacer Con Las Órdenes Pendientes?',
      closeModalCancel: 'Cancelar',
      closeModalCancelPending: 'Cancelar órdenes pendientes',
      closeModalCancelPendingHint: 'estado cancelado',
      closeModalKeepForNext: 'Mantener para próxima jornada',
      closeModalKeepForNextHint: 'estado inicial',
      closeModalConfirm: 'Cerrar Jornada',
      closing: 'Cerrando...',
      success: 'Jornada Cerrada Correctamente.',
      comingSoon: 'El Cierre De Jornada Estará Disponible En Una Próxima Actualización.',
      noMenu: 'No Se Encontró Un Menú',
      noSession: 'No Hay Sesión Activa',
      errorLoading: 'Error Al Cargar Los Datos De La Jornada.',
      noActiveJourney: 'No tienes una jornada abierta. Abre una para comenzar a registrar órdenes del día.',
      openJourney: 'Abrir jornada',
      openingJourney: 'Abriendo jornada...',
      errorCreatingJourney: 'Error al abrir la jornada. Intenta de nuevo.',
      openJourneyReason: 'Apertura manual',
      reports: 'Reportes de Jornadas',
      downloadExcel: 'Descargar Excel',
      reportGenerating: 'Generando...',
      stats: 'Estadísticas',
      totalRevenue: 'Ventas totales',
      revenueConcreted: 'Ventas concretadas',
      revenueDelivered: 'Entregadas',
      revenueDispatched: 'Despachadas',
      revenuePending: 'Ventas pendientes',
      revenueCancelled: 'Canceladas',
      statsOrders: 'Órdenes',
      averageTicket: 'Ticket promedio',
      topByRevenue: 'Top ventas',
      topByQuantity: 'Top unidades',
      chartByRevenue: 'Por ventas',
      chartByProfit: 'Top productos por ganancias',
      chartByQuantity: 'Por unidades',
      topProducts: 'Productos más vendidos',
      noStats: 'No hay datos de ventas para esta jornada.',
      errorLoadingStats: 'Error al cargar estadísticas.',
      prev: 'Anterior',
      next: 'Siguiente',
      page: 'Página',
      workdayReport: 'Reporte de Jornada',
      duration: 'Duración',
      itemsSold: 'Ítems vendidos',
      itemsOrdered: 'Ítems ordenados',
      downloadCSV: 'Descargar CSV',
      productsTable: 'Productos',
      productName: 'Producto',
      quantity: 'Cant.',
      revenue: 'Ventas',
      cost: 'Costo',
      margin: 'Margen',
      profit: 'Ganancias',
      marginPercent: 'margen',
      totalCost: 'Costo total',
      totalProfit: 'Ganancias totales',
      noCostConfigured: 'Configura costos en Cost para ver el margen',
      totalMargin: 'Margen',
      topProduct: 'Producto más vendido'
    },
    negocios: {
      title: 'Negocios',
      subtitle: 'Selecciona tu negocio activo o crea uno nuevo para trabajar.',
      createNew: 'Crear nuevo menú',
      creating: 'Creando...',
      yourBusinesses: 'Tus negocios',
      noBusinesses: 'Aún no tienes negocios. Crea uno arriba.',
      business: 'Negocio',
      active: 'Activo',
      select: 'Seleccionar',
      errorLoading: 'Error al cargar los negocios.',
      errorCreating: 'Error al crear el nuevo menú.',
      slugLabel: 'Identificador de tu negocio (slug)',
      slugPlaceholder: 'ej: mi-restaurante',
      slugHint: 'Solo letras minúsculas, números y guiones. Será la URL de tu carta.',
      create: 'Crear',
      cancel: 'Cancelar',
      newMenuTitle: 'Nuevo negocio'
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
      templateLabel: 'Estilo',
      templateHero: 'Hero',
      templateModern: 'Modern',
      useThisDesign: 'Usar este design',
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
    cost: {
      title: 'Custos do cardápio',
      subtitle: 'Configure o custo de cada item e acompanhamento. Usado para margens e relatórios.',
      unit: 'Unidade de venda',
      baseUnit: 'Base',
      priceSale: 'Preço venda',
      cost: 'Custo',
      save: 'Salvar custos',
      reload: 'Recarregar'
    },
    sidebar: {
      negocios: 'Negócios',
      aiAssistant: 'Assistente IA',
      history: 'Histórico',
      gallery: 'Galeria',
      qrCode: 'Código QR',
      orders: 'Kanban',
      jornada: 'Jornada',
      reportes: 'Relatórios',
      cost: 'Custos',
      kitchen: 'Cozinha',
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
      noSession: 'Não Há Sessão Ativa',
      noMenu: 'Nenhum Cardápio Encontrado',
      empty: 'Ainda Não Há Pedidos.',
      subtitle: 'Ordenado Por Hora Comprometida. Vista Orientada À Cozinha.',
      delivery: 'Entrega',
      pickup: 'Retirada',
      viewAsPaper: 'Ver Como Folha',
      print: 'Imprimir',
      printThermal: 'Imprimir Em Térmica',
      eventPayload: 'event_payload',
      itemsToPrepare: 'O Que Preparar',
      forTime: 'Para',
      remainingIn: 'Em {min} Min',
      late: 'Atrasado {min} Min',
      itemsCount: '{count} Itens',
      markAsReady: 'Pronto',
      kitchenMode: 'Modo Full',
      exitKitchenMode: 'Sair Do Modo Full',
      statusPending: 'Pendente',
      statusPreparing: 'Em Preparação',
      statusDone: 'Pronto',
      tabPending: 'Pendentes',
      tabPreparing: 'Em Preparação',
      tabDone: 'Prontos',
      statusGeneralLabel: 'Estado Geral',
      readyToDeliver: 'Pronto Para Entregar',
      startPreparing: 'Iniciar Preparação',
      reload: 'Atualizar',
      filterAll: 'Entrega',
      filterKitchen: 'Cozinha',
      filterBar: 'Bar',
      showQR: 'QR',
      emptyForStation: 'Nenhum Pedido Para Esta Estação.',
      emptyCancelled: 'Nenhum Pedido Cancelado.',
      startToReady: 'Iniciar → Pronto',
      deliver: 'Entregar',
      dispatch: 'Despachar',
      delivered: 'Entregado',
      dispatched: 'Despachado',
      cancelled: 'Cancelado',
      cancelOrder: 'Cancelar Pedido',
      cancelModalTitle: 'Cancelar Pedido',
      cancelModalSubtitle: 'Esta Ação Não Pode Ser Desfeita.',
      cancelModalReasonLabel: 'Motivo (Escolha Um):',
      cancelReasons: {
        outOfStock: 'Falta De Estoque',
        orderError: 'Erro No Pedido',
        customerLeft: 'Cliente Foi Embora',
        paymentIssue: 'Problema De Pagamento',
        other: 'Outro'
      },
      cancelModalCommentLabel: 'Comentário (Opcional):',
      cancelModalCommentPlaceholder: 'Ex: Cliente Não Atendeu...',
      cancelModalBack: 'Voltar',
      cancelModalConfirm: 'Confirmar Cancelamento',
      viewVertical: 'Vertical',
      viewThreeColumns: '3 Colunas',
      pendingView: 'Pendentes',
      pendingFrom: 'De',
      pendingTo: 'Até',
      pendingSelectAll: 'Selecionar todas',
      pendingSelected: 'pedidos selecionados',
      pendingAssignToJourney: 'Atribuir ao turno ativo',
      pendingCancelOrders: 'Cancelar pedidos',
      pendingClearSelection: 'Limpar seleção',
      pendingViewDetails: 'Ver',
      pendingAssign: 'Atribuir',
      pendingEmpty: 'Não há pedidos pendentes.',
      pendingNoActiveJourney: 'Abra um turno para atribuir pedidos.'
    },
    jornada: {
      title: 'Jornada',
      active: 'Jornada Ativa',
      date: 'Data',
      openedSince: 'Aberta Desde',
      summary: 'Resumo Rápido',
      totalOrders: 'Pedidos Totais',
      delivered: 'Entregues',
      cancelled: 'Cancelados',
      pending: 'Pendentes',
      closeWorkday: 'Fechar Jornada',
      closeModalTitle: 'Fechar Jornada',
      closeModalMessage: 'Você Está Prestes A Fechar A Jornada De {date}. O Que Deseja Fazer Com Os Pedidos Pendentes?',
      closeModalCancel: 'Cancelar',
      closeModalCancelPending: 'Cancelar pedidos pendentes',
      closeModalCancelPendingHint: 'estado cancelado',
      closeModalKeepForNext: 'Manter para próxima jornada',
      closeModalKeepForNextHint: 'estado inicial',
      closeModalConfirm: 'Fechar Jornada',
      closing: 'Fechando...',
      success: 'Jornada Fechada Com Sucesso.',
      comingSoon: 'O Fechamento De Jornada Estará Disponível Em Uma Próxima Atualização.',
      noMenu: 'Nenhum Cardápio Encontrado',
      noSession: 'Não Há Sessão Ativa',
      errorLoading: 'Erro Ao Carregar Os Dados Da Jornada.',
      reports: 'Relatórios de Jornadas',
      downloadExcel: 'Baixar Excel',
      reportGenerating: 'Gerando...',
      stats: 'Estatísticas',
      totalRevenue: 'Vendas totais',
      revenueDelivered: 'Vendas entregues',
      revenuePending: 'Vendas pendentes',
      revenueCancelled: 'Vendas canceladas',
      statsOrders: 'Pedidos',
      averageTicket: 'Ticket médio',
      topByRevenue: 'Top vendas',
      topByQuantity: 'Top unidades',
      chartByRevenue: 'Por vendas',
      chartByProfit: 'Top produtos por lucro',
      chartByQuantity: 'Por unidades',
      topProducts: 'Produtos mais vendidos',
      noStats: 'Não há dados de vendas para esta jornada.',
      errorLoadingStats: 'Erro ao carregar estatísticas.',
      prev: 'Anterior',
      next: 'Próximo',
      page: 'Página',
      back: 'Voltar aos relatórios',
      workdayReport: 'Relatório de Jornada',
      duration: 'Duração',
      itemsSold: 'Itens vendidos',
      itemsOrdered: 'Itens pedidos',
      downloadCSV: 'Baixar CSV',
      productsTable: 'Produtos',
      productName: 'Produto',
      quantity: 'Qtd.',
      revenue: 'Vendas',
      cost: 'Custo',
      margin: 'Margem',
      totalCost: 'Custo total',
      totalProfit: 'Lucro total',
      noCostConfigured: 'Configure custos em Cost para ver a margem',
      totalMargin: 'Margem',
      topProduct: 'Produto mais vendido'
    },
    negocios: {
      title: 'Negócios',
      subtitle: 'Selecione seu negócio ativo ou crie um novo para trabalhar.',
      createNew: 'Criar novo cardápio',
      creating: 'Criando...',
      yourBusinesses: 'Seus negócios',
      noBusinesses: 'Você ainda não tem negócios. Crie um acima.',
      business: 'Negócio',
      active: 'Ativo',
      select: 'Selecionar',
      errorLoading: 'Erro ao carregar os negócios.',
      errorCreating: 'Erro ao criar o novo cardápio.',
      slugLabel: 'Identificador do seu negócio (slug)',
      slugPlaceholder: 'ex: meu-restaurante',
      slugHint: 'Apenas letras minúsculas, números e hífens. Será a URL da sua carta.',
      create: 'Criar',
      cancel: 'Cancelar',
      newMenuTitle: 'Novo negócio'
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
      discardMenuCancel: 'Cancel',
      templateLabel: 'Style',
      templateHero: 'Hero',
      templateModern: 'Modern',
      useThisDesign: 'Use this design'
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
    cost: {
      title: 'Menu costs',
      subtitle: 'Configure the cost of each item and side. Used for margins and reports.',
      unit: 'Unit of sale',
      baseUnit: 'Base',
      priceSale: 'Sale price',
      cost: 'Cost',
      save: 'Save',
      saved: 'Price and cost saved successfully',
      reload: 'Reload'
    },
    sidebar: {
      negocios: 'Businesses',
      aiAssistant: 'AI Assistant',
      history: 'History',
      gallery: 'Gallery',
      qrCode: 'QR Code',
      orders: 'Kanban',
      jornada: 'Workday',
      reportes: 'Reports',
      cost: 'Cost',
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
      noSession: 'No Active Session',
      noMenu: 'No Menu Found',
      empty: 'No Orders Yet.',
      subtitle: 'Sorted By Requested Time. Kitchen-Oriented View.',
      delivery: 'Delivery',
      pickup: 'Pickup',
      viewAsPaper: 'View As Sheet',
      print: 'Print',
      printThermal: 'Print To Thermal',
      eventPayload: 'event_payload',
      itemsToPrepare: 'What To Prepare',
      forTime: 'For',
      remainingIn: 'In {min} Min',
      late: '{min} Min Late',
      itemsCount: '{count} Items',
      markAsReady: 'Ready',
      kitchenMode: 'Full Mode',
      exitKitchenMode: 'Exit Full Mode',
      statusPending: 'Pending',
      statusPreparing: 'Preparing',
      statusDone: 'Ready',
      tabPending: 'Pending',
      tabPreparing: 'Preparing',
      tabDone: 'Ready',
      statusGeneralLabel: 'Overall Status',
      readyToDeliver: 'Ready For Delivery',
      startPreparing: 'Start Preparation',
      reload: 'Refresh',
      filterAll: 'Delivery',
      filterKitchen: 'Kitchen',
      filterBar: 'Bar',
      showQR: 'QR',
      emptyForStation: 'No Orders For This Station.',
      emptyCancelled: 'No Cancelled Orders.',
      startToReady: 'Start → Ready',
      deliver: 'Deliver',
      dispatch: 'Dispatch',
      delivered: 'Delivered',
      dispatched: 'Dispatched',
      cancelled: 'Cancelled',
      cancelOrder: 'Cancel Order',
      cancelModalTitle: 'Cancel Order',
      cancelModalSubtitle: 'This Action Cannot Be Undone.',
      cancelModalReasonLabel: 'Reason (Choose One):',
      cancelReasons: {
        outOfStock: 'Out Of Stock',
        orderError: 'Order Error',
        customerLeft: 'Customer Left',
        paymentIssue: 'Payment Issue',
        other: 'Other'
      },
      cancelModalCommentLabel: 'Comment (Optional):',
      cancelModalCommentPlaceholder: 'E.g. Customer Did Not Answer...',
      cancelModalBack: 'Back',
      cancelModalConfirm: 'Confirm Cancellation',
      viewVertical: 'Vertical',
      viewThreeColumns: '3 Columns',
      pendingView: 'Pending',
      pendingFrom: 'From',
      pendingTo: 'To',
      pendingSelectAll: 'Select all',
      pendingSelected: 'orders selected',
      pendingAssignToJourney: 'Assign to active workday',
      pendingCancelOrders: 'Cancel orders',
      pendingClearSelection: 'Clear selection',
      pendingViewDetails: 'View',
      pendingAssign: 'Assign',
      pendingEmpty: 'No pending orders.',
      pendingNoActiveJourney: 'Open a workday to assign orders.'
    },
    jornada: {
      title: 'Workday',
      active: 'Active Workday',
      date: 'Date',
      openedSince: 'Opened Since',
      summary: 'Quick Summary',
      totalOrders: 'Total Orders',
      delivered: 'Delivered',
      cancelled: 'Cancelled',
      pending: 'Pending',
      closeWorkday: 'Close Workday',
      closeModalTitle: 'Close Workday',
      closeModalMessage: 'You Are About To Close The Workday Of {date}. What Do You Want To Do With Pending Orders?',
      closeModalCancel: 'Cancel',
      closeModalCancelPending: 'Cancel pending orders',
      closeModalCancelPendingHint: 'cancelled status',
      closeModalKeepForNext: 'Keep for next journey',
      closeModalKeepForNextHint: 'initial status',
      closeModalConfirm: 'Close Workday',
      closing: 'Closing...',
      success: 'Workday Closed Successfully.',
      comingSoon: 'Workday Closure Will Be Available In An Upcoming Update.',
      noMenu: 'No Menu Found',
      noSession: 'No Active Session',
      errorLoading: 'Error Loading Workday Data.',
      noActiveJourney: "You don't have an open workday. Open one to start recording today's orders.",
      openJourney: 'Open workday',
      openingJourney: 'Opening workday...',
      errorCreatingJourney: 'Error opening workday. Please try again.',
      openJourneyReason: 'Manual opening',
      reports: 'Workday Reports',
      downloadExcel: 'Download Excel',
      reportGenerating: 'Generating...',
      stats: 'Statistics',
      totalRevenue: 'Total sales',
      revenueConcreted: 'Concreted sales',
      revenueDelivered: 'Delivered',
      revenueDispatched: 'Dispatched',
      revenuePending: 'Pending sales',
      revenueCancelled: 'Cancelled',
      statsOrders: 'Orders',
      averageTicket: 'Average ticket',
      topByRevenue: 'Top by revenue',
      topByQuantity: 'Top by quantity',
      chartByRevenue: 'By revenue',
      chartByProfit: 'Top products by profit',
      chartByQuantity: 'By quantity',
      topProducts: 'Top products',
      noStats: 'No sales data for this workday.',
      errorLoadingStats: 'Error loading statistics.',
      prev: 'Previous',
      next: 'Next',
      page: 'Page',
      back: 'Back to reports',
      workdayReport: 'Workday Report',
      duration: 'Duration',
      itemsSold: 'Items sold',
      itemsOrdered: 'Items ordered',
      downloadCSV: 'Download CSV',
      productsTable: 'Products',
      productName: 'Product',
      quantity: 'Qty',
      revenue: 'Revenue',
      cost: 'Cost',
      margin: 'Margin',
      profit: 'Profit',
      marginPercent: 'margin',
      totalCost: 'Total cost',
      totalProfit: 'Total profit',
      noCostConfigured: 'Configure costs in Cost to see margin',
      totalMargin: 'Margin',
      topProduct: 'Top product'
    },
    negocios: {
      title: 'Businesses',
      subtitle: 'Select your active business or create a new one to work on.',
      createNew: 'Create new menu',
      creating: 'Creating...',
      yourBusinesses: 'Your businesses',
      noBusinesses: "You don't have any businesses yet. Create one above.",
      business: 'Business',
      active: 'Active',
      select: 'Select',
      errorLoading: 'Error loading businesses.',
      errorCreating: 'Error creating new menu.',
      slugLabel: 'Business identifier (slug)',
      slugPlaceholder: 'e.g. my-restaurant',
      slugHint: 'Lowercase letters, numbers and hyphens only. Will be your menu URL.',
      create: 'Create',
      cancel: 'Cancel',
      newMenuTitle: 'New business'
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

