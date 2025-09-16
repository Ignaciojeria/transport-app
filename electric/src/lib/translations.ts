export type Language = 'CL' | 'BR' | 'EU'

export interface MobileTranslations {
  // Header
  header: {
    routeId: string
    vehiclePlate: string
  }
  
  // Navigation
  navigation: {
    map: string
    list: string
    start: string
  }
  
  // Next Visit Card
  nextVisit: {
    title: string
    multipleClients: string
    client: string
    clients: string
    selectClient: string
  }
  
  // Map View
  mapView: {
    multipleClientsAtLocation: string
    selectClientToDeliver: string
  }
  
  // Visit Card
  visitCard: {
    sequence: string
    deliveryUnits: string
    deliveryUnit: string
    orders: string
    order: string
    quantity: string
    weight: string
    volume: string
    units: string
    unit: string
    selectedClient: string
    selectedVisit: string
    nextToDeliver: string
    instructions: string
  }
  
  // Delivery Actions
  delivery: {
    deliver: string
    deliverAll: string
    deliverRemaining: string
    notDeliver: string
    notDeliverAll: string
    notDeliverRemaining: string
    delivered: string
    notDelivered: string
    next: string
    groupActions: string
    actionsForRemaining: string
    changeToDelivered: string
    changeToNotDelivered: string
    singleUnitPending: string
  }
  
  // Status
  status: {
    pending: string
    completed: string
    delivered: string
    pendingUnits: string
  }
  
  // Tabs
  tabs: {
    enRoute: string
    delivered: string
    notDelivered: string
  }
  
  // Route Start Modal
  routeStartModal: {
    title: string
    licensePlateLabel: string
    licensePlatePlaceholder: string
    startButton: string
    cancelButton: string
    validationRequired: string
    useAssignedPlate: string
    useThisPlate: string
    useOtherPlate: string
    differentPlateLabel: string
    plateExample: string
    plateHelpText: string
    plateHelpTextDefault: string
  }
}

export const translations: Record<Language, MobileTranslations> = {
  CL: {
    header: {
      routeId: "ID RUTA",
      vehiclePlate: "PATENTE"
    },
    navigation: {
      map: "Mapa",
      list: "Lista", 
      start: "Iniciar"
    },
    nextVisit: {
      title: "Siguiente Visita",
      multipleClients: "Múltiples clientes",
      client: "Cliente",
      clients: "clientes",
      selectClient: "Selecciona un cliente en el selector de arriba"
    },
    mapView: {
      multipleClientsAtLocation: "Múltiples clientes en esta ubicación:",
      selectClientToDeliver: "Selecciona el cliente al que vas a entregar"
    },
    visitCard: {
      sequence: "Secuencia",
      deliveryUnits: "Unidades de Entrega",
      deliveryUnit: "Unidad de Entrega",
      orders: "órdenes",
      order: "Orden",
      quantity: "Cant.",
      weight: "Peso",
      volume: "Vol",
      units: "unidades",
      unit: "unidad",
      selectedClient: "Cliente seleccionado",
      selectedVisit: "Visita seleccionada",
      nextToDeliver: "Siguiente a entregar",
      instructions: "Instrucciones:"
    },
    delivery: {
      deliver: "Entregar",
      deliverAll: "Entregar todo",
      deliverRemaining: "Entregar restantes",
      notDeliver: "No entregar",
      notDeliverAll: "No entregar todo", 
      notDeliverRemaining: "No entregar restantes",
      delivered: "Entregado",
      notDelivered: "No entregado",
      next: "Siguiente a entregar",
      groupActions: "Acciones grupales:",
      actionsForRemaining: "Acciones para restantes:",
      changeToDelivered: "Cambiar a entregado",
      changeToNotDelivered: "Cambiar a no entregado",
      singleUnitPending: "Solo queda 1 unidad pendiente. Usa los botones individuales abajo."
    },
    status: {
      pending: "Pendiente",
      completed: "Completado",
      delivered: "entregadas",
      pendingUnits: "pendientes"
    },
    tabs: {
      enRoute: "En ruta",
      delivered: "Entregados",
      notDelivered: "No entregados"
    },
    routeStartModal: {
      title: "Iniciar Ruta",
      licensePlateLabel: "Patente del Vehículo",
      licensePlatePlaceholder: "Ingresa la patente...",
      startButton: "Iniciar ruta",
      cancelButton: "Cancelar",
      validationRequired: "La patente es requerida",
      useAssignedPlate: "¿Usar la patente asignada a esta ruta?",
      useThisPlate: "Usar esta",
      useOtherPlate: "Usar otra",
      differentPlateLabel: "O ingresa una patente diferente:",
      plateExample: "Ej: ABC123",
      plateHelpText: "Escribe aquí si quieres usar una patente diferente a la asignada",
      plateHelpTextDefault: "Puedes ingresar cualquier patente para esta ruta"
    }
  },
  
  BR: {
    header: {
      routeId: "ID DA ROTA",
      vehiclePlate: "PLACA"
    },
    navigation: {
      map: "Mapa",
      list: "Lista",
      start: "Iniciar"
    },
    nextVisit: {
      title: "Próxima Visita",
      multipleClients: "Múltiplos clientes",
      client: "Cliente",
      clients: "clientes",
      selectClient: "Selecione um cliente no seletor acima"
    },
    mapView: {
      multipleClientsAtLocation: "Múltiplos clientes nesta localização:",
      selectClientToDeliver: "Selecione o cliente ao qual vai entregar"
    },
    visitCard: {
      sequence: "Sequência",
      deliveryUnits: "Unidades de Entrega",
      deliveryUnit: "Unidade de Entrega",
      orders: "pedidos",
      order: "Pedido",
      quantity: "Qtd.",
      weight: "Peso",
      volume: "Vol",
      units: "unidades",
      unit: "unidade",
      selectedClient: "Cliente selecionado",
      selectedVisit: "Visita selecionada",
      nextToDeliver: "Próximo a entregar",
      instructions: "Instruções:"
    },
    delivery: {
      deliver: "Entregar",
      deliverAll: "Entregar tudo",
      deliverRemaining: "Entregar restantes",
      notDeliver: "Não entregar",
      notDeliverAll: "Não entregar tudo",
      notDeliverRemaining: "Não entregar restantes", 
      delivered: "Entregue",
      notDelivered: "Não entregue",
      next: "Próximo a entregar",
      groupActions: "Ações em grupo:",
      actionsForRemaining: "Ações para restantes:",
      changeToDelivered: "Mudar para entregue",
      changeToNotDelivered: "Mudar para não entregue",
      singleUnitPending: "Resta apenas 1 unidade pendente. Use os botões individuais abaixo."
    },
    status: {
      pending: "Pendente",
      completed: "Concluído",
      delivered: "entregues",
      pendingUnits: "pendentes"
    },
    tabs: {
      enRoute: "Em rota",
      delivered: "Entregues",
      notDelivered: "Não entregues"
    },
    routeStartModal: {
      title: "Iniciar Rota",
      licensePlateLabel: "Placa do Veículo",
      licensePlatePlaceholder: "Digite a placa...",
      startButton: "Iniciar rota",
      cancelButton: "Cancelar",
      validationRequired: "A placa é obrigatória",
      useAssignedPlate: "Usar a placa atribuída a esta rota?",
      useThisPlate: "Usar esta",
      useOtherPlate: "Usar outra",
      differentPlateLabel: "Ou digite uma placa diferente:",
      plateExample: "Ex: ABC123",
      plateHelpText: "Digite aqui se quiser usar uma placa diferente da atribuída",
      plateHelpTextDefault: "Você pode inserir qualquer placa para esta rota"
    }
  },
  
  EU: {
    header: {
      routeId: "ROUTE ID",
      vehiclePlate: "LICENSE PLATE"
    },
    navigation: {
      map: "Map",
      list: "List",
      start: "Start"
    },
    nextVisit: {
      title: "Next Visit",
      multipleClients: "Multiple clients",
      client: "Client",
      clients: "clients",
      selectClient: "Select a client in the selector above"
    },
    mapView: {
      multipleClientsAtLocation: "Multiple clients at this location:",
      selectClientToDeliver: "Select the client you want to deliver to"
    },
    visitCard: {
      sequence: "Sequence",
      deliveryUnits: "Delivery Units",
      deliveryUnit: "Delivery Unit",
      orders: "orders",
      order: "Order",
      quantity: "Qty.",
      weight: "Weight",
      volume: "Vol",
      units: "units",
      unit: "unit",
      selectedClient: "Selected client",
      selectedVisit: "Selected visit",
      nextToDeliver: "Next to deliver",
      instructions: "Instructions:"
    },
    delivery: {
      deliver: "Deliver",
      deliverAll: "Deliver all",
      deliverRemaining: "Deliver remaining",
      notDeliver: "Don't deliver",
      notDeliverAll: "Don't deliver all",
      notDeliverRemaining: "Don't deliver remaining",
      delivered: "Delivered",
      notDelivered: "Not delivered",
      next: "Next to deliver",
      groupActions: "Group actions:",
      actionsForRemaining: "Actions for remaining:",
      changeToDelivered: "Change to delivered",
      changeToNotDelivered: "Change to not delivered",
      singleUnitPending: "Only 1 unit pending. Use individual buttons below."
    },
    status: {
      pending: "Pending",
      completed: "Completed",
      delivered: "delivered",
      pendingUnits: "pending"
    },
    tabs: {
      enRoute: "In route",
      delivered: "Delivered",
      notDelivered: "Not delivered"
    },
    routeStartModal: {
      title: "Start Route",
      licensePlateLabel: "Vehicle License Plate",
      licensePlatePlaceholder: "Enter license plate...",
      startButton: "Start route",
      cancelButton: "Cancel",
      validationRequired: "License plate is required",
      useAssignedPlate: "Use the assigned license plate for this route?",
      useThisPlate: "Use this",
      useOtherPlate: "Use other",
      differentPlateLabel: "Or enter a different license plate:",
      plateExample: "Ex: ABC123",
      plateHelpText: "Type here if you want to use a different license plate from the assigned one",
      plateHelpTextDefault: "You can enter any license plate for this route"
    }
  }
}

export const languageNames: Record<Language, string> = {
  CL: 'Español',
  BR: 'Português', 
  EU: 'English'
}

export const languageFlags: Record<Language, string> = {
  CL: '🇨🇱',
  BR: '🇧🇷',
  EU: '🇪🇺'
}
