export type Language = 'EN' | 'ES' | 'PT'

export interface Translations {
  // Home.svelte
  home: {
    loading: string
    errorLoading: string
    ourMenu: string
    businessHours: string
    contact: string
  }
  
  // HorariosSection.svelte
  hours: {
    notSpecified: string
  }
  
  // CartaSection.svelte
  menu: {
    noItems: string
    addToCart: string
    addedToCart: string
    add1ToCart: string
    viewOptions: string
    fromPrice: string
    chooseOption: string
    optional: string
    allergens: {
      gluten: string
      egg: string
      lactose: string
      seafood: string
      nuts: string
      dairy: string
      soy: string
      vegan: string
      vegetarian: string
      spicy: string
      alcohol: string
    }
  }
  
  // FloatingCart.svelte & ShoppingCart.svelte
  cart: {
    yourOrder: string
    total: string
    quantity: string
    item: string
    items: string
    sendOrderWhatsApp: string
    orderWhatsApp: string
    placeOrder: string
    clearOrder: string
    pickupInfo: string
    pickupName: string
    nameLabel: string
    pickupTime: string
    timeFormat: string
    cancel: string
    sendOrder: string
    confirmClear: string
    confirmClearMessage: string
    yes: string
    no: string
    completeFields: string
    invalidTime: string
    timeFormatExample: string
    nameFormatExample: string
    phone: string
    phonePlaceholder: string
    email: string
    emailPlaceholder: string
    deliveryType: string
    delivery: string
    pickup: string
    digital: string
    digitalInfo: string
      deliveryAddress: string
      deliveryAddressPlaceholder: string
      searchingAddress: string
      addressNotFound: string
      selectDeliveryType: string
      addressNumber: string
      addressNumberPlaceholder: string
      addressNotes: string
      addressNotesPlaceholder: string
      confirmAddress: string
      next: string
      back: string
      step1Title: string
      step2Title: string
      step1Label: string
      step2Label: string
    shoppingCart: string
      emptyCart: string
      emptyCartMessage: string
      clear: string
      preparingOrder: string
      orderSent: string
      viewYourOrder: string
      viewYourOrders: string
    }
  
  // TrackingView.svelte
  tracking: {
    shareTracking: string
    shareCopied: string
    contactStore: string
    productCancelled: string
    orderCancelled: string
    backToMenu: string
    activeOrders: string
    activeOrdersTitle: string
    recentOrders: string
    orderTracking: string
    yourOrdersInProgress: string
    loading: string
    consult: string
    consultByCode: string
    codePlaceholder: string
    searching: string
    loadingOrder: string
    delivery: string
    pickup: string
    digital: string
    statusConfirmed: string
    statusPreparing: string
    statusOnTheWay: string
    statusReadyForPickup: string
    statusDelivered: string
    statusCancelled: string
    statusConsult: string
    progress: string
    received: string
    preparing: string
    detail: string
    total: string
    deliveredCheck: string
    whatsappGreeting: string
    orderLabel: string
    orderReceived: string
    dateLabel: string
    codeLabel: string
    statusLabel: string
    errorFetching: string
    statusPendingNoJourney: string
    orderInfo: string
    customerLabel: string
    deliveryInfo: string
    unit: string
    notes: string
  }

  // TemplateSelector (preview mode)
  preview: {
    templateLabel: string
    templateHero: string
    templateModern: string
    useThisDesign: string
  }

  // WhatsApp messages
  whatsapp: {
    greeting: string
    orderItem: string
    each: string
    itemTotal: string
    orderTotal: string
    pickupInfoLabel: string
    pickupNameLabel: string
    pickupTimeLabel: string
  }
}

export const translations: Record<Language, Translations> = {
  ES: {
    home: {
      loading: 'Cargando catálogo...',
      errorLoading: 'Error al cargar los datos',
      ourMenu: 'Catálogo',
      businessHours: 'Horarios de Atención',
      contact: 'Contacto'
    },
    hours: {
      notSpecified: 'No especificado'
    },
    menu: {
      noItems: 'No hay elementos',
      addToCart: 'Agregar',
      addedToCart: '✓ Agregado',
      add1ToCart: 'Agregar 1',
      viewOptions: 'Ver opciones',
      fromPrice: 'Desde',
      chooseOption: 'Elige tu opción',
      optional: 'Opcional',
      allergens: {
        gluten: 'Contiene Gluten',
        egg: 'Contiene Huevo',
        lactose: 'Contiene Lactosa',
        seafood: 'Contiene Mariscos',
        nuts: 'Contiene Frutos Secos',
        dairy: 'Contiene Lácteos',
        soy: 'Contiene Soja',
        vegan: 'Vegano',
        vegetarian: 'Vegetariano',
        spicy: 'Picante',
        alcohol: 'Contiene Alcohol'
      }
    },
    cart: {
      yourOrder: 'Tu pedido',
      total: 'Total',
      quantity: 'Cantidad:',
      item: 'item',
      items: 'items',
      sendOrderWhatsApp: 'Enviar Pedido por WhatsApp',
      orderWhatsApp: 'Pedir por WhatsApp',
      placeOrder: 'Realizar pedido',
      clearOrder: 'Limpiar pedido',
      pickupInfo: 'Información de Retiro',
      pickupName: 'Nombre de quien va a retirar *',
      nameLabel: 'Nombre *',
      pickupTime: 'Hora de retiro *',
      timeFormat: 'Formato: HH:MM (24 horas, ejemplo: 14:30 para 2:30 PM)',
      cancel: 'Cancelar',
      sendOrder: 'Enviar Pedido',
      confirmClear: '¿Limpiar pedido?',
      confirmClearMessage: '¿Estás seguro de que deseas limpiar todo el pedido?',
      yes: 'Sí',
      no: 'No',
      completeFields: 'Por favor completa todos los campos',
      invalidTime: 'Por favor ingresa una hora válida (formato: HH:MM, ejemplo: 14:30)',
      timeFormatExample: 'Ej: 14:30',
      nameFormatExample: 'Ej: Juan Perez',
      phone: 'Teléfono',
      phonePlaceholder: 'Ej: +56912345678',
      email: 'Email',
      emailPlaceholder: 'Ej: correo@ejemplo.com',
      shoppingCart: 'Carrito de Compras',
      emptyCart: 'Tu carrito está vacío',
      emptyCartMessage: 'Agrega items del menú para comenzar',
      clear: 'Vaciar',
      preparingOrder: 'Preparando tu pedido...',
      orderSent: 'Tu pedido se ha enviado correctamente',
      viewYourOrder: 'Ver tu pedido',
      viewYourOrders: 'Ver tus pedidos',
      deliveryType: 'Tipo de entrega',
      delivery: 'Envío a domicilio',
      pickup: 'Retiro en tienda',
      digital: 'Producto digital',
      digitalInfo: 'Información para entrega digital',
      deliveryAddress: 'Dirección de envío *',
      deliveryAddressPlaceholder: 'Ej: Av. Providencia 123, Santiago',
      searchingAddress: 'Buscando dirección...',
      addressNotFound: 'No se encontró la dirección. Intenta con otra.',
      selectDeliveryType: 'Selecciona el tipo de entrega',
      addressNumber: 'Número de casa/departamento',
      addressNumberPlaceholder: 'Ej: 123, Depto 4B',
      addressNotes: 'Indicaciones adicionales (opcional)',
      addressNotesPlaceholder: 'Ej: Portón azul, timbre 2',
      confirmAddress: 'Confirmar dirección',
      next: 'Siguiente',
      back: 'Atrás',
      step1Title: 'Paso 1: Dirección de entrega',
      step2Title: 'Paso 2: Información de contacto',
      step1Label: 'Dirección',
      step2Label: 'Contacto'
    },
    tracking: {
      shareTracking: 'Compartir seguimiento',
      shareCopied: 'Enlace copiado',
      contactStore: 'Contactar con la tienda',
      productCancelled: 'Producto cancelado',
      orderCancelled: 'Pedido cancelado',
      backToMenu: '← Volver al menú',
      activeOrders: '← Pedidos activos',
      activeOrdersTitle: 'Pedidos activos',
      recentOrders: 'Pedidos recientes',
      orderTracking: 'Seguimiento de pedido',
      yourOrdersInProgress: 'Tus pedidos en curso y recientes',
      loading: '→ Cargando...',
      consult: 'Consultar',
      consultByCode: 'Consultar pedido por código',
      codePlaceholder: 'Ej: ABC12345',
      searching: 'Buscando...',
      loadingOrder: 'Cargando pedido...',
      delivery: 'Envío a domicilio',
      pickup: 'Retiro en local',
      digital: 'Producto digital',
      statusConfirmed: 'Pedido confirmado',
      statusPreparing: 'En preparación',
      statusOnTheWay: 'En camino',
      statusReadyForPickup: 'Listo para entregar',
      statusDelivered: 'Entregado',
      statusCancelled: 'Cancelado',
      statusConsult: 'Consultar',
      progress: 'Progreso',
      received: 'Pedido recibido',
      preparing: 'En preparación',
      detail: 'Detalle',
      total: 'Total',
      deliveredCheck: 'Entregado ✓',
      whatsappGreeting: 'Hola, tengo una consulta sobre mi pedido:',
      orderLabel: 'Pedido',
      orderReceived: 'Pedido recibido',
      dateLabel: 'Fecha',
      codeLabel: 'Código',
      statusLabel: 'Estado',
      errorFetching: 'Error al consultar el pedido',
      statusPendingNoJourney: 'Pedido recibido. Se preparará cuando el negocio abra',
      orderInfo: 'Datos del pedido',
      deliveryInfo: 'Dirección de entrega',
      unit: 'Depto/Unidad',
      notes: 'Notas'
    },
    preview: {
      templateLabel: 'Estilo',
      templateHero: 'Hero',
      templateModern: 'Modern',
      useThisDesign: 'Usar este diseño'
    },
    whatsapp: {
      greeting: '¡Hola! Me gustaría hacer el siguiente pedido:\n\n',
      orderItem: 'Pedido',
      each: 'c/u',
      itemTotal: 'Total',
      orderTotal: 'Total',
      pickupInfoLabel: 'Información de retiro:',
      pickupNameLabel: '👤 Nombre:',
      pickupTimeLabel: '🕐 Hora de retiro:'
    }
  },
  PT: {
    home: {
      loading: 'Carregando dados do restaurante...',
      errorLoading: 'Erro ao carregar os dados',
      ourMenu: 'Nosso Cardápio',
      businessHours: 'Horários de Funcionamento',
      contact: 'Contato'
    },
    hours: {
      notSpecified: 'Não especificado'
    },
    menu: {
      noItems: 'Não há elementos',
      addToCart: 'Adicionar',
      addedToCart: '✓ Adicionado',
      viewOptions: 'Ver opções',
      fromPrice: 'A partir de',
      chooseOption: 'Escolha sua opção',
      allergens: {
        gluten: 'Contém Glúten',
        egg: 'Contém Ovo',
        lactose: 'Contém Lactose',
        seafood: 'Contém Frutos do Mar',
        nuts: 'Contém Frutos Secos',
        dairy: 'Contém Lácteos',
        soy: 'Contém Soja',
        vegan: 'Vegano',
        vegetarian: 'Vegetariano',
        spicy: 'Picante',
        alcohol: 'Contém Álcool'
      }
    },
    cart: {
      yourOrder: 'Seu pedido',
      total: 'Total',
      quantity: 'Quantidade:',
      item: 'item',
      items: 'itens',
      sendOrderWhatsApp: 'Enviar Pedido por WhatsApp',
      orderWhatsApp: 'Pedir por WhatsApp',
      placeOrder: 'Realizar pedido',
      clearOrder: 'Limpar pedido',
      pickupInfo: 'Informações de Retirada',
      pickupName: 'Nome de quem vai retirar *',
      nameLabel: 'Nome *',
      pickupTime: 'Hora de retirada *',
      timeFormat: 'Formato: HH:MM (24 horas, exemplo: 14:30 para 2:30 PM)',
      cancel: 'Cancelar',
      sendOrder: 'Enviar Pedido',
      confirmClear: 'Limpar pedido?',
      confirmClearMessage: 'Tem certeza de que deseja limpar todo o pedido?',
      yes: 'Sim',
      no: 'Não',
      completeFields: 'Por favor, preencha todos os campos',
      invalidTime: 'Por favor, insira uma hora válida (formato: HH:MM, exemplo: 14:30)',
      timeFormatExample: 'Ex: 14:30',
      nameFormatExample: 'Ex: João Silva',
      phone: 'Telefone',
      phonePlaceholder: 'Ex: +5511999999999',
      email: 'Email',
      emailPlaceholder: 'Ex: email@exemplo.com',
      shoppingCart: 'Carrinho de Compras',
      emptyCart: 'Seu carrinho está vazio',
      emptyCartMessage: 'Adicione itens do cardápio para começar',
      clear: 'Limpar',
      preparingOrder: 'Preparando seu pedido...',
      orderSent: 'Seu pedido foi enviado com sucesso',
      viewYourOrder: 'Ver seu pedido',
      viewYourOrders: 'Ver seus pedidos',
      deliveryType: 'Tipo de entrega',
      delivery: 'Entrega em domicílio',
      pickup: 'Retirada na loja',
      digital: 'Produto digital',
      digitalInfo: 'Informações para entrega digital',
      deliveryAddress: 'Endereço de entrega *',
      deliveryAddressPlaceholder: 'Ex: Av. Paulista 123, São Paulo',
      searchingAddress: 'Buscando endereço...',
      addressNotFound: 'Endereço não encontrado. Tente outro.',
      selectDeliveryType: 'Selecione o tipo de entrega',
      addressNumber: 'Número da casa/apartamento',
      addressNumberPlaceholder: 'Ex: 123, Apt 4B',
      addressNotes: 'Instruções adicionais (opcional)',
      addressNotesPlaceholder: 'Ex: Portão azul, interfone 2',
      confirmAddress: 'Confirmar endereço',
      next: 'Próximo',
      back: 'Voltar',
      step1Title: 'Passo 1: Endereço de entrega',
      step2Title: 'Passo 2: Informações de contato',
      step1Label: 'Endereço',
      step2Label: 'Contato'
    },
    tracking: {
      shareTracking: 'Compartir rastreamento',
      shareCopied: 'Link copiado',
      contactStore: 'Contatar a loja',
      productCancelled: 'Produto cancelado',
      orderCancelled: 'Pedido cancelado',
      backToMenu: '← Voltar ao cardápio',
      activeOrders: '← Pedidos ativos',
      activeOrdersTitle: 'Pedidos ativos',
      recentOrders: 'Pedidos recentes',
      orderTracking: 'Rastreamento de pedido',
      yourOrdersInProgress: 'Seus pedidos em andamento e recentes',
      loading: '→ Carregando...',
      consult: 'Consultar',
      consultByCode: 'Consultar pedido por código',
      codePlaceholder: 'Ex: ABC12345',
      searching: 'Buscando...',
      loadingOrder: 'Carregando pedido...',
      delivery: 'Entrega em domicílio',
      pickup: 'Retirada na loja',
      digital: 'Produto digital',
      statusConfirmed: 'Pedido confirmado',
      statusPreparing: 'Em preparação',
      statusOnTheWay: 'A caminho',
      statusReadyForPickup: 'Pronto para retirar',
      statusDelivered: 'Entregue',
      statusCancelled: 'Cancelado',
      statusConsult: 'Consultar',
      progress: 'Progresso',
      received: 'Pedido recebido',
      preparing: 'Em preparação',
      detail: 'Detalhe',
      total: 'Total',
      deliveredCheck: 'Entregue ✓',
      whatsappGreeting: 'Olá, tenho uma dúvida sobre meu pedido:',
      orderLabel: 'Pedido',
      orderReceived: 'Pedido recebido',
      dateLabel: 'Data',
      codeLabel: 'Código',
      statusLabel: 'Status',
      errorFetching: 'Erro ao consultar o pedido',
      statusPendingNoJourney: 'Pedido recebido. Será preparado quando o negócio abrir',
      orderInfo: 'Dados do pedido',
      customerLabel: 'Cliente',
      deliveryInfo: 'Endereço de entrega',
      unit: 'Apto/Unidade',
      notes: 'Notas'
    },
    preview: {
      templateLabel: 'Estilo',
      templateHero: 'Hero',
      templateModern: 'Modern',
      useThisDesign: 'Usar este estilo'
    },
    whatsapp: {
      greeting: 'Olá! Gostaria de fazer o seguinte pedido:\n\n',
      orderItem: 'Pedido',
      each: 'c/u',
      itemTotal: 'Total',
      orderTotal: 'Total',
      pickupInfoLabel: 'Informações de retirada:',
      pickupNameLabel: '👤 Nome:',
      pickupTimeLabel: '🕐 Hora de retirada:'
    }
  },
  EN: {
    home: {
      loading: 'Loading restaurant data...',
      errorLoading: 'Error loading data',
      ourMenu: 'Our Menu',
      businessHours: 'Business Hours',
      contact: 'Contact'
    },
    hours: {
      notSpecified: 'Not specified'
    },
    menu: {
      noItems: 'No items',
      addToCart: 'Add',
      addedToCart: '✓ Added',
      add1ToCart: 'Add 1',
      viewOptions: 'View options',
      fromPrice: 'From',
      chooseOption: 'Choose your option',
      optional: 'Optional',
      allergens: {
        gluten: 'Contains Gluten',
        egg: 'Contains Egg',
        lactose: 'Contains Lactose',
        seafood: 'Contains Seafood',
        nuts: 'Contains Nuts',
        dairy: 'Contains Dairy',
        soy: 'Contains Soy',
        vegan: 'Vegan',
        vegetarian: 'Vegetarian',
        spicy: 'Spicy',
        alcohol: 'Contains Alcohol'
      }
    },
    cart: {
      yourOrder: 'Your order',
      total: 'Total',
      quantity: 'Quantity:',
      item: 'item',
      items: 'items',
      sendOrderWhatsApp: 'Send Order via WhatsApp',
      orderWhatsApp: 'Order via WhatsApp',
      placeOrder: 'Place order',
      clearOrder: 'Clear order',
      pickupInfo: 'Pickup Information',
      pickupName: 'Name of person picking up *',
      nameLabel: 'Name *',
      pickupTime: 'Pickup time *',
      timeFormat: 'Format: HH:MM (24 hours, example: 14:30 for 2:30 PM)',
      cancel: 'Cancel',
      sendOrder: 'Send Order',
      confirmClear: 'Clear order?',
      confirmClearMessage: 'Are you sure you want to clear the entire order?',
      yes: 'Yes',
      no: 'No',
      completeFields: 'Please complete all fields',
      invalidTime: 'Please enter a valid time (format: HH:MM, example: 14:30)',
      timeFormatExample: 'Ex: 14:30',
      nameFormatExample: 'Ex: John Doe',
      phone: 'Phone',
      phonePlaceholder: 'Ex: +1234567890',
      email: 'Email',
      emailPlaceholder: 'Ex: email@example.com',
      shoppingCart: 'Shopping Cart',
      emptyCart: 'Your cart is empty',
      emptyCartMessage: 'Add items from the menu to get started',
      clear: 'Clear',
      preparingOrder: 'Preparing your order...',
      orderSent: 'Your order has been sent successfully',
      viewYourOrder: 'View your order',
      viewYourOrders: 'View your orders',
      deliveryType: 'Delivery type',
      delivery: 'Home delivery',
      pickup: 'Store pickup',
      digital: 'Digital product',
      digitalInfo: 'Information for digital delivery',
      deliveryAddress: 'Delivery address *',
      deliveryAddressPlaceholder: 'Ex: 123 Main St, New York',
      searchingAddress: 'Searching address...',
      addressNotFound: 'Address not found. Try another one.',
      selectDeliveryType: 'Select delivery type',
      addressNumber: 'House/apartment number',
      addressNumberPlaceholder: 'Ex: 123, Apt 4B',
      addressNotes: 'Additional instructions (optional)',
      addressNotesPlaceholder: 'Ex: Blue gate, buzzer 2',
      confirmAddress: 'Confirm address',
      next: 'Next',
      back: 'Back',
      step1Title: 'Step 1: Delivery address',
      step2Title: 'Step 2: Contact information',
      step1Label: 'Address',
      step2Label: 'Contact'
    },
    tracking: {
      shareTracking: 'Share tracking',
      shareCopied: 'Link copied',
      contactStore: 'Contact store',
      productCancelled: 'Product cancelled',
      orderCancelled: 'Order cancelled',
      backToMenu: '← Back to menu',
      activeOrders: '← Active orders',
      activeOrdersTitle: 'Active orders',
      recentOrders: 'Recent orders',
      orderTracking: 'Order tracking',
      yourOrdersInProgress: 'Your orders in progress and recent',
      loading: '→ Loading...',
      consult: 'Search',
      consultByCode: 'Search order by code',
      codePlaceholder: 'E.g. ABC12345',
      searching: 'Searching...',
      loadingOrder: 'Loading order...',
      delivery: 'Home delivery',
      pickup: 'Store pickup',
      digital: 'Digital product',
      statusConfirmed: 'Order confirmed',
      statusPreparing: 'Preparing',
      statusOnTheWay: 'On the way',
      statusReadyForPickup: 'Ready for pickup',
      statusDelivered: 'Delivered',
      statusCancelled: 'Cancelled',
      statusConsult: 'Check status',
      progress: 'Progress',
      received: 'Order received',
      preparing: 'Preparing',
      detail: 'Detail',
      total: 'Total',
      deliveredCheck: 'Delivered ✓',
      whatsappGreeting: 'Hello, I have a question about my order:',
      orderLabel: 'Order',
      orderReceived: 'Order received',
      dateLabel: 'Date',
      codeLabel: 'Code',
      statusLabel: 'Status',
      errorFetching: 'Error fetching order',
      statusPendingNoJourney: 'Order received. It will be prepared when the business opens',
      orderInfo: 'Order details',
      customerLabel: 'Customer',
      deliveryInfo: 'Delivery address',
      unit: 'Unit',
      notes: 'Notes'
    },
    preview: {
      templateLabel: 'Style',
      templateHero: 'Hero',
      templateModern: 'Modern',
      useThisDesign: 'Use this design'
    },
    whatsapp: {
      greeting: 'Hello! I would like to place the following order:\n\n',
      orderItem: 'Order',
      each: 'each',
      itemTotal: 'Total',
      orderTotal: 'Total',
      pickupInfoLabel: 'Pickup information:',
      pickupNameLabel: '👤 Name:',
      pickupTimeLabel: '🕐 Pickup time:'
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

