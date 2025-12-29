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
    shoppingCart: string
    emptyCart: string
    emptyCartMessage: string
    clear: string
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
      loading: 'Cargando datos del restaurante...',
      errorLoading: 'Error al cargar los datos',
      ourMenu: 'Nuestra Carta',
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
      shoppingCart: 'Carrito de Compras',
      emptyCart: 'Tu carrito está vacío',
      emptyCartMessage: 'Agrega items del menú para comenzar',
      clear: 'Vaciar'
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
      shoppingCart: 'Carrinho de Compras',
      emptyCart: 'Seu carrinho está vazio',
      emptyCartMessage: 'Adicione itens do cardápio para começar',
      clear: 'Limpar'
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
      shoppingCart: 'Shopping Cart',
      emptyCart: 'Your cart is empty',
      emptyCartMessage: 'Add items from the menu to get started',
      clear: 'Clear'
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

