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
    clearOrder: string
    pickupInfo: string
    pickupName: string
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
    deliveryType: string
    delivery: string
    pickup: string
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
      redirectingWhatsApp: string
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
      noItems: 'No hay elementos'
    },
    cart: {
      yourOrder: 'Tu pedido',
      total: 'Total',
      quantity: 'Cantidad:',
      item: 'item',
      items: 'items',
      sendOrderWhatsApp: 'Enviar Pedido por WhatsApp',
      orderWhatsApp: 'Pedir por WhatsApp',
      clearOrder: 'Limpiar pedido',
      pickupInfo: 'Información de Retiro',
      pickupName: 'Nombre de quien va a retirar *',
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
      shoppingCart: 'Carrito de Compras',
      emptyCart: 'Tu carrito está vacío',
      emptyCartMessage: 'Agrega items del menú para comenzar',
      clear: 'Vaciar',
      preparingOrder: 'Preparando tu pedido...',
      redirectingWhatsApp: 'Redirigiendo a WhatsApp...',
      deliveryType: 'Tipo de entrega',
      delivery: 'Envío a domicilio',
      pickup: 'Retiro en tienda',
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
      noItems: 'Não há elementos'
    },
    cart: {
      yourOrder: 'Seu pedido',
      total: 'Total',
      quantity: 'Quantidade:',
      item: 'item',
      items: 'itens',
      sendOrderWhatsApp: 'Enviar Pedido por WhatsApp',
      orderWhatsApp: 'Pedir por WhatsApp',
      clearOrder: 'Limpar pedido',
      pickupInfo: 'Informações de Retirada',
      pickupName: 'Nome de quem vai retirar *',
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
      shoppingCart: 'Carrinho de Compras',
      emptyCart: 'Seu carrinho está vazio',
      emptyCartMessage: 'Adicione itens do cardápio para começar',
      clear: 'Limpar',
      preparingOrder: 'Preparando seu pedido...',
      redirectingWhatsApp: 'Redirecionando para o WhatsApp...',
      deliveryType: 'Tipo de entrega',
      delivery: 'Entrega em domicílio',
      pickup: 'Retirada na loja',
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
      noItems: 'No items'
    },
    cart: {
      yourOrder: 'Your order',
      total: 'Total',
      quantity: 'Quantity:',
      item: 'item',
      items: 'items',
      sendOrderWhatsApp: 'Send Order via WhatsApp',
      orderWhatsApp: 'Order via WhatsApp',
      clearOrder: 'Clear order',
      pickupInfo: 'Pickup Information',
      pickupName: 'Name of person picking up *',
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
      shoppingCart: 'Shopping Cart',
      emptyCart: 'Your cart is empty',
      emptyCartMessage: 'Add items from the menu to get started',
      clear: 'Clear',
      preparingOrder: 'Preparing your order...',
      redirectingWhatsApp: 'Redirecting to WhatsApp...',
      deliveryType: 'Delivery type',
      delivery: 'Home delivery',
      pickup: 'Store pickup',
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

