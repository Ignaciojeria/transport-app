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
  
  // Delivery Modal
  deliveryModal: {
    title: string
    recipientNameLabel: string
    recipientNamePlaceholder: string
    documentLabel: string
    documentPlaceholder: string
    photoLabel: string
    activateCamera: string
    cancel: string
    confirmDelivery: string
    uploadingImage: string
  }
  
  // Non-Delivery Modal
  nonDeliveryModal: {
    title: string
    reasonLabel: string
    reasonPlaceholder: string
    observationsLabel: string
    observationsPlaceholder: string
    photoLabel: string
    activateCamera: string
    cancel: string
    confirmNonDelivery: string
    uploadingImage: string
    reasons: {
      clientRejects: string
      noResidents: string
      damagedProduct: string
      otherReason: string
    }
  }
  
  // Grouped Delivery Modal
  groupedDeliveryModal: {
    title: string
    unitsFor: string
    deliveryAddress: string
    unitsPending: string
    unitsToDeliver: string
    order: string
    items: string
    recipientInfo: string
    fullNameRequired: string
    fullNamePlaceholder: string
    documentOptional: string
    documentPlaceholder: string
    photographicEvidence: string
    captureGroupEvidence: string
    takePhoto: string
    actionWarning: string
    cancel: string
    processing: string
    delivered: string
    notDelivered: string
  }
  
  // Grouped Non-Delivery Modal
  groupedNonDeliveryModal: {
    title: string
    unitsFor: string
    deliveryAddress: string
    unitsPending: string
    unitsNotToDeliver: string
    order: string
    items: string
    nonDeliveryReason: string
    specifyReason: string
    reasonPlaceholder: string
    photographicEvidence: string
    captureGroupEvidence: string
    takePhoto: string
    actionWarning: string
    cancel: string
    processing: string
    markAsNotDelivered: string
    delivered: string
    notDelivered: string
    reasons: {
      recipientNotAvailable: string
      addressNotFound: string
      recipientRefused: string
      damagedPackage: string
      incorrectAddress: string
      securityIssue: string
      other: string
    }
  }
  
  // Sidebar Menu
  sidebar: {
    title: string
    closeMenu: string
    reports: string
    downloadReport: string
    connectionStatus: string
    internet: string
    connected: string
    disconnected: string
    synchronization: string
    markerSync: string
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
    deliveryModal: {
      title: "Evidencia de entrega",
      recipientNameLabel: "Nombre de quien recibe",
      recipientNamePlaceholder: "Nombre completo",
      documentLabel: "RUT / Documento",
      documentPlaceholder: "12.345.678-9",
      photoLabel: "Foto de evidencia",
      activateCamera: "Activar cámara",
      cancel: "Cancelar",
      confirmDelivery: "Confirmar entrega",
      uploadingImage: "Subiendo imagen..."
    },
    nonDeliveryModal: {
      title: "No entregado",
      reasonLabel: "Motivo",
      reasonPlaceholder: "Buscar/ingresar motivo",
      observationsLabel: "Observaciones",
      observationsPlaceholder: "Detalles adicionales (opcional)",
      photoLabel: "Foto de evidencia",
      activateCamera: "Activar cámara",
      cancel: "Cancelar",
      confirmNonDelivery: "Confirmar no entrega",
      uploadingImage: "Subiendo imagen...",
      reasons: {
        clientRejects: "cliente rechaza entrega",
        noResidents: "sin moradores",
        damagedProduct: "producto dañado",
        otherReason: "otro motivo"
      }
    },
    groupedDeliveryModal: {
      title: "Entregar todo",
      unitsFor: "unidades para",
      deliveryAddress: "Dirección de entrega",
      unitsPending: "pendientes",
      unitsToDeliver: "Unidades a entregar",
      order: "Orden",
      items: "ítems",
      recipientInfo: "Información del receptor",
      fullNameRequired: "Nombre completo *",
      fullNamePlaceholder: "Nombre del receptor",
      documentOptional: "RUT (opcional)",
      documentPlaceholder: "12.345.678-9",
      photographicEvidence: "Evidencia fotográfica",
      captureGroupEvidence: "Capturar evidencia de entrega grupal",
      takePhoto: "Tomar foto",
      actionWarning: "Esta acción marcará todas las unidades como entregadas",
      cancel: "Cancelar",
      processing: "Procesando...",
      delivered: "Entregado",
      notDelivered: "No entregado"
    },
    groupedNonDeliveryModal: {
      title: "No entregar todo",
      unitsFor: "unidades para",
      deliveryAddress: "Dirección de entrega",
      unitsPending: "pendientes",
      unitsNotToDeliver: "Unidades que no se entregarán",
      order: "Orden",
      items: "ítems",
      nonDeliveryReason: "Motivo de no entrega *",
      specifyReason: "Especificar motivo",
      reasonPlaceholder: "Describe el motivo específico...",
      photographicEvidence: "Evidencia fotográfica *",
      captureGroupEvidence: "Capturar evidencia de no entrega grupal",
      takePhoto: "Tomar foto",
      actionWarning: "Esta acción marcará todas las unidades como no entregadas",
      cancel: "Cancelar",
      processing: "Procesando...",
      markAsNotDelivered: "Marcar como no entregado",
      delivered: "Entregado",
      notDelivered: "No entregado",
      reasons: {
        recipientNotAvailable: "Receptor no disponible",
        addressNotFound: "Dirección no encontrada",
        recipientRefused: "Receptor rechazó la entrega",
        damagedPackage: "Paquete dañado",
        incorrectAddress: "Dirección incorrecta",
        securityIssue: "Problema de seguridad",
        other: "Otro"
      }
    },
    sidebar: {
      title: "Menú",
      closeMenu: "Cerrar menú",
      reports: "Reportes",
      downloadReport: "Descargar Reporte",
      connectionStatus: "Estado de Conexión",
      internet: "Internet",
      connected: "Conectado",
      disconnected: "Desconectado",
      synchronization: "Sincronización",
      markerSync: "Marcador Sync"
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
    deliveryModal: {
      title: "Evidência de entrega",
      recipientNameLabel: "Nome de quem recebe",
      recipientNamePlaceholder: "Nome completo",
      documentLabel: "CPF / Documento",
      documentPlaceholder: "123.456.789-01",
      photoLabel: "Foto de evidência",
      activateCamera: "Ativar câmera",
      cancel: "Cancelar",
      confirmDelivery: "Confirmar entrega",
      uploadingImage: "Enviando imagem..."
    },
    nonDeliveryModal: {
      title: "Não entregue",
      reasonLabel: "Motivo",
      reasonPlaceholder: "Buscar/inserir motivo",
      observationsLabel: "Observações",
      observationsPlaceholder: "Detalhes adicionais (opcional)",
      photoLabel: "Foto de evidência",
      activateCamera: "Ativar câmera",
      cancel: "Cancelar",
      confirmNonDelivery: "Confirmar não entrega",
      uploadingImage: "Enviando imagem...",
      reasons: {
        clientRejects: "cliente recusa entrega",
        noResidents: "sem moradores",
        damagedProduct: "produto danificado",
        otherReason: "outro motivo"
      }
    },
    groupedDeliveryModal: {
      title: "Entregar tudo",
      unitsFor: "unidades para",
      deliveryAddress: "Endereço de entrega",
      unitsPending: "pendentes",
      unitsToDeliver: "Unidades a entregar",
      order: "Pedido",
      items: "itens",
      recipientInfo: "Informações do destinatário",
      fullNameRequired: "Nome completo *",
      fullNamePlaceholder: "Nome do destinatário",
      documentOptional: "CPF (opcional)",
      documentPlaceholder: "123.456.789-01",
      photographicEvidence: "Evidência fotográfica",
      captureGroupEvidence: "Capturar evidência de entrega em grupo",
      takePhoto: "Tirar foto",
      actionWarning: "Esta ação marcará todas as unidades como entregues",
      cancel: "Cancelar",
      processing: "Processando...",
      delivered: "Entregue",
      notDelivered: "Não entregue"
    },
    groupedNonDeliveryModal: {
      title: "Não entregar tudo",
      unitsFor: "unidades para",
      deliveryAddress: "Endereço de entrega",
      unitsPending: "pendentes",
      unitsNotToDeliver: "Unidades que não serão entregues",
      order: "Pedido",
      items: "itens",
      nonDeliveryReason: "Motivo da não entrega *",
      specifyReason: "Especificar motivo",
      reasonPlaceholder: "Descreva o motivo específico...",
      photographicEvidence: "Evidência fotográfica *",
      captureGroupEvidence: "Capturar evidência de não entrega em grupo",
      takePhoto: "Tirar foto",
      actionWarning: "Esta ação marcará todas as unidades como não entregues",
      cancel: "Cancelar",
      processing: "Processando...",
      markAsNotDelivered: "Marcar como não entregue",
      delivered: "Entregue",
      notDelivered: "Não entregue",
      reasons: {
        recipientNotAvailable: "Destinatário não disponível",
        addressNotFound: "Endereço não encontrado",
        recipientRefused: "Destinatário recusou a entrega",
        damagedPackage: "Pacote danificado",
        incorrectAddress: "Endereço incorreto",
        securityIssue: "Problema de segurança",
        other: "Outro"
      }
    },
    sidebar: {
      title: "Menu",
      closeMenu: "Fechar menu",
      reports: "Relatórios",
      downloadReport: "Baixar Relatório",
      connectionStatus: "Status da Conexão",
      internet: "Internet",
      connected: "Conectado",
      disconnected: "Desconectado",
      synchronization: "Sincronização",
      markerSync: "Sincronização de Marcador"
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
      vehiclePlate: "PLATE"
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
    deliveryModal: {
      title: "Delivery Evidence",
      recipientNameLabel: "Recipient name",
      recipientNamePlaceholder: "Full name",
      documentLabel: "ID / Document",
      documentPlaceholder: "123-45-6789",
      photoLabel: "Evidence photo",
      activateCamera: "Activate camera",
      cancel: "Cancel",
      confirmDelivery: "Confirm delivery",
      uploadingImage: "Uploading image..."
    },
    nonDeliveryModal: {
      title: "Not delivered",
      reasonLabel: "Reason",
      reasonPlaceholder: "Search/enter reason",
      observationsLabel: "Observations",
      observationsPlaceholder: "Additional details (optional)",
      photoLabel: "Evidence photo",
      activateCamera: "Activate camera",
      cancel: "Cancel",
      confirmNonDelivery: "Confirm non-delivery",
      uploadingImage: "Uploading image...",
      reasons: {
        clientRejects: "customer refuses delivery",
        noResidents: "no residents",
        damagedProduct: "damaged product",
        otherReason: "other reason"
      }
    },
    groupedDeliveryModal: {
      title: "Deliver all",
      unitsFor: "units for",
      deliveryAddress: "Delivery address",
      unitsPending: "pending",
      unitsToDeliver: "Units to deliver",
      order: "Order",
      items: "items",
      recipientInfo: "Recipient information",
      fullNameRequired: "Full name *",
      fullNamePlaceholder: "Recipient name",
      documentOptional: "ID (optional)",
      documentPlaceholder: "123-45-6789",
      photographicEvidence: "Photographic evidence",
      captureGroupEvidence: "Capture group delivery evidence",
      takePhoto: "Take photo",
      actionWarning: "This action will mark all units as delivered",
      cancel: "Cancel",
      processing: "Processing...",
      delivered: "Delivered",
      notDelivered: "Not delivered"
    },
    groupedNonDeliveryModal: {
      title: "Don't deliver all",
      unitsFor: "units for",
      deliveryAddress: "Delivery address",
      unitsPending: "pending",
      unitsNotToDeliver: "Units not to deliver",
      order: "Order",
      items: "items",
      nonDeliveryReason: "Non-delivery reason *",
      specifyReason: "Specify reason",
      reasonPlaceholder: "Describe the specific reason...",
      photographicEvidence: "Photographic evidence *",
      captureGroupEvidence: "Capture group non-delivery evidence",
      takePhoto: "Take photo",
      actionWarning: "This action will mark all units as not delivered",
      cancel: "Cancel",
      processing: "Processing...",
      markAsNotDelivered: "Mark as not delivered",
      delivered: "Delivered",
      notDelivered: "Not delivered",
      reasons: {
        recipientNotAvailable: "Recipient not available",
        addressNotFound: "Address not found",
        recipientRefused: "Recipient refused delivery",
        damagedPackage: "Damaged package",
        incorrectAddress: "Incorrect address",
        securityIssue: "Security issue",
        other: "Other"
      }
    },
    sidebar: {
      title: "Menu",
      closeMenu: "Close menu",
      reports: "Reports",
      downloadReport: "Download Report",
      connectionStatus: "Connection Status",
      internet: "Internet",
      connected: "Connected",
      disconnected: "Disconnected",
      synchronization: "Synchronization",
      markerSync: "Marker Sync"
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
