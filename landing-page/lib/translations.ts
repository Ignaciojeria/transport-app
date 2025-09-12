export type Language = 'CL' | 'BR' | 'EU'

export interface Translations {
  // Navigation
  nav: {
    howItWorks: string
    benefits: string
    freeEvaluation: string
  }
  
  // Hero Section
  hero: {
    title: string
    subtitle: string
    freeEvaluation: string
    personalizedConsultation: string
    vehicles: string
    deliveries: string
    routes: string
    vehicleDetails: {
      license: string
      weight: string
      volume: string
      insurance: string
      additionalVars: string
      moreVehicles: string
    }
    deliveryDetails: {
      client: string
      address: string
      weight: string
      volume: string
      price: string
      additionalVars: string
      moreDeliveries: string
    }
    routeDetails: {
      route: string
      assignedLicense: string
      deliveries: string
      additionalVars: string
      moreRoutes: string
    }
  }
  
  // How it Works Section
  howItWorks: {
    title: string
    subtitle: string
    step1: {
      title: string
      description: string
    }
    step2: {
      title: string
      description: string
    }
    step3: {
      title: string
      description: string
    }
  }
  
  // Benefits Section
  benefits: {
    title: string
    subtitle: string
    reduceCosts: {
      title: string
      description: string
    }
    balanceLoad: {
      title: string
      description: string
    }
    saveTime: {
      title: string
      description: string
    }
    realTimeTracking: {
      title: string
      description: string
    }
  }
  
  // CTA Section
  cta: {
    title: string
    subtitle: string
    freeEvaluation: string
    personalizedConsultation: string
  }
  
  // Demo Section
  demo: {
    title: string
    subtitle: string
    routeId: string
    interactiveDemo: string
    features: {
      startRoute: {
        title: string
        description: string
      }
      markDeliveries: {
        title: string
        description: string
      }
      navigateMaps: {
        title: string
        description: string
      }
      generateReports: {
        title: string
        description: string
      }
    }
    buttons: {
      viewDemo: string
      hideDemo: string
      openNewTab: string
    }
    preview: {
      title: string
      description: string
      features: {
        simulatedDeliveries: string
        fullFunctionality: string
        realisticData: string
      }
    }
  }
  
  // Footer
  footer: {
    description: string
    product: string
    support: string
    company: string
    howItWorks: string
    prices: string
    demo: string
    help: string
    contact: string
    documentation: string
    about: string
    blog: string
    careers: string
    copyright: string
  }
}

export const translations: Record<Language, Translations> = {
  CL: {
    nav: {
      howItWorks: "Cómo funciona",
      benefits: "Beneficios",
      freeEvaluation: "Evaluación Gratuita"
    },
    hero: {
      title: "Optimiza tu flota en",
      subtitle: "Optimiza rutas, genera enlaces para conductores y monitorea entregas en tiempo real.",
      freeEvaluation: "Evaluación Gratuita",
      personalizedConsultation: "Consulta Personalizada",
      vehicles: "Vehículos",
      deliveries: "Entregas",
      routes: "Rutas",
      vehicleDetails: {
        license: "Patente",
        weight: "Peso",
        volume: "Vol",
        insurance: "Seguro",
        additionalVars: "Variables adicionales...",
        moreVehicles: "+ 3 vehículos más..."
      },
      deliveryDetails: {
        client: "Cliente",
        address: "Dirección",
        weight: "Peso",
        volume: "Vol",
        price: "Precio",
        additionalVars: "Variables adicionales...",
        moreDeliveries: "+ 8 entregas más..."
      },
      routeDetails: {
        route: "Ruta",
        assignedLicense: "Patente asignada",
        deliveries: "entregas",
        additionalVars: "Variables adicionales...",
        moreRoutes: "+ 2 rutas más..."
      }
    },
    howItWorks: {
      title: "¿Cómo funciona?",
      subtitle: "En solo 3 pasos simples, optimiza tu logística y reduce costos",
      step1: {
        title: "Crea tus vehículos",
        description: "Configura y registra tus vehículos con sus capacidades, dimensiones y características. Define la flota que utilizarás para tus entregas."
      },
      step2: {
        title: "Carga las entregas",
        description: "Crea y configura tus entregas con destinos, productos y restricciones. Define todos los detalles necesarios para la planificación de rutas."
      },
      step3: {
        title: "Optimiza y ejecuta",
        description: "Llama a nuestro agente de optimización y recibe los enlaces para el conductor. Ejecuta las rutas optimizadas con trazabilidad en tiempo real."
      }
    },
    benefits: {
      title: "Beneficios para empresas",
      subtitle: "Optimiza tu operación logística y reduce costos significativamente",
      reduceCosts: {
        title: "Reduce costos logísticos",
        description: "Optimiza rutas y reduce combustible, tiempo y recursos."
      },
      balanceLoad: {
        title: "Balancea carga automáticamente",
        description: "Distribuye la carga de manera óptima entre vehículos, maximizando eficiencia y minimizando viajes vacíos"
      },
      saveTime: {
        title: "Ahorra tiempo planificando",
        description: "De 2 horas a 5 minutos. Planifica rutas automáticamente"
      },
      realTimeTracking: {
        title: "Trazabilidad en tiempo real",
        description: "Monitorea el progreso de entregas en tiempo real desde la web mobile, sin necesidad de software adicional"
      }
    },
    cta: {
      title: "¿Listo para optimizar tu flota?",
      subtitle: "Agenda una evaluación gratuita de tus procesos logísticos",
      freeEvaluation: "Evaluación Gratuita",
      personalizedConsultation: "Consulta Personalizada"
    },
    demo: {
      title: "Prueba la App en Vivo",
      subtitle: "Experimenta cómo funciona la app desde la perspectiva del conductor con datos de prueba",
      routeId: "ID de Ruta",
      interactiveDemo: "Demo Interactiva",
      features: {
        startRoute: {
          title: "Inicia una ruta",
          description: "Simula el inicio de una ruta con patente"
        },
        markDeliveries: {
          title: "Marca entregas",
          description: "Simula el proceso de entrega con evidencia"
        },
        navigateMaps: {
          title: "Navega con mapas",
          description: "Integración con Google Maps y Waze"
        },
        generateReports: {
          title: "Genera reportes",
          description: "Descarga reportes en CSV y Excel"
        }
      },
      buttons: {
        viewDemo: "Ver Demo",
        hideDemo: "Ocultar Demo",
        openNewTab: "Abrir en Nueva Pestaña"
      },
      preview: {
        title: "Demo Interactiva",
        description: "Haz clic en \"Ver Demo\" para experimentar la aplicación con datos simulados",
        features: {
          simulatedDeliveries: "9 entregas simuladas",
          fullFunctionality: "Funcionalidad completa",
          realisticData: "Datos realistas de Santiago"
        }
      }
    },
    footer: {
      description: "Optimiza rutas, genera enlaces para conductores y monitorea entregas en tiempo real. Todo desde una sola plataforma.",
      product: "Producto",
      support: "Soporte",
      company: "Empresa",
      howItWorks: "Cómo funciona",
      prices: "Precios",
      demo: "Demo",
      help: "Ayuda",
      contact: "Contacto",
      documentation: "Documentación",
      about: "Acerca de",
      blog: "Blog",
      careers: "Carreras",
      copyright: "© 2025 TransportApp. Todos los derechos reservados."
    }
  },
  
  BR: {
    nav: {
      howItWorks: "Como funciona",
      benefits: "Benefícios",
      freeEvaluation: "Avaliação Gratuita"
    },
    hero: {
      title: "Otimize sua frota em",
      subtitle: "Otimize rotas, gere links para motoristas e monitore entregas em tempo real.",
      freeEvaluation: "Avaliação Gratuita",
      personalizedConsultation: "Consulta Personalizada",
      vehicles: "Veículos",
      deliveries: "Entregas",
      routes: "Rotas",
      vehicleDetails: {
        license: "Placa",
        weight: "Peso",
        volume: "Vol",
        insurance: "Seguro",
        additionalVars: "Variáveis adicionais...",
        moreVehicles: "+ 3 veículos mais..."
      },
      deliveryDetails: {
        client: "Cliente",
        address: "Endereço",
        weight: "Peso",
        volume: "Vol",
        price: "Preço",
        additionalVars: "Variáveis adicionais...",
        moreDeliveries: "+ 8 entregas mais..."
      },
      routeDetails: {
        route: "Rota",
        assignedLicense: "Placa atribuída",
        deliveries: "entregas",
        additionalVars: "Variáveis adicionais...",
        moreRoutes: "+ 2 rotas mais..."
      }
    },
    howItWorks: {
      title: "Como funciona?",
      subtitle: "Em apenas 3 passos simples, otimize sua logística e reduza custos",
      step1: {
        title: "Crie seus veículos",
        description: "Configure e registre seus veículos com suas capacidades, dimensões e características. Defina a frota que utilizará para suas entregas."
      },
      step2: {
        title: "Carregue as entregas",
        description: "Crie e configure suas entregas com destinos, produtos e restrições. Defina todos os detalhes necessários para o planejamento de rotas."
      },
      step3: {
        title: "Otimize e execute",
        description: "Chame nosso agente de otimização e receba os links para o motorista. Execute as rotas otimizadas com rastreabilidade em tempo real."
      }
    },
    benefits: {
      title: "Benefícios para empresas",
      subtitle: "Otimize sua operação logística e reduza custos significativamente",
      reduceCosts: {
        title: "Reduz custos logísticos",
        description: "Otimize rotas e reduza combustível, tempo e recursos."
      },
      balanceLoad: {
        title: "Balanceia carga automaticamente",
        description: "Distribui a carga de maneira ótima entre veículos, maximizando eficiência e minimizando viagens vazias"
      },
      saveTime: {
        title: "Economiza tempo planejando",
        description: "De 2 horas a 5 minutos. Planeje rotas automaticamente"
      },
      realTimeTracking: {
        title: "Rastreabilidade em tempo real",
        description: "Monitore o progresso das entregas em tempo real a partir da web mobile, sem necessidade de software adicional"
      }
    },
    cta: {
      title: "Pronto para otimizar sua frota?",
      subtitle: "Agende uma avaliação gratuita de seus processos logísticos",
      freeEvaluation: "Avaliação Gratuita",
      personalizedConsultation: "Consulta Personalizada"
    },
    demo: {
      title: "Teste o App ao Vivo",
      subtitle: "Experimente como funciona o app na perspectiva do motorista com dados de teste",
      routeId: "ID da Rota",
      interactiveDemo: "Demo Interativa",
      features: {
        startRoute: {
          title: "Inicia uma rota",
          description: "Simula o início de uma rota com placa"
        },
        markDeliveries: {
          title: "Marca entregas",
          description: "Simula o processo de entrega com evidência"
        },
        navigateMaps: {
          title: "Navega com mapas",
          description: "Integração com Google Maps e Waze"
        },
        generateReports: {
          title: "Gera relatórios",
          description: "Baixa relatórios em CSV e Excel"
        }
      },
      buttons: {
        viewDemo: "Ver Demo",
        hideDemo: "Ocultar Demo",
        openNewTab: "Abrir em Nova Aba"
      },
      preview: {
        title: "Demo Interativa",
        description: "Clique em \"Ver Demo\" para experimentar a aplicação com dados simulados",
        features: {
          simulatedDeliveries: "9 entregas simuladas",
          fullFunctionality: "Funcionalidade completa",
          realisticData: "Dados realistas de Santiago"
        }
      }
    },
    footer: {
      description: "Otimize rotas, gere links para motoristas e monitore entregas em tempo real. Tudo em uma única plataforma.",
      product: "Produto",
      support: "Suporte",
      company: "Empresa",
      howItWorks: "Como funciona",
      prices: "Preços",
      demo: "Demo",
      help: "Ajuda",
      contact: "Contato",
      documentation: "Documentação",
      about: "Sobre",
      blog: "Blog",
      careers: "Carreiras",
      copyright: "© 2025 TransportApp. Todos os direitos reservados."
    }
  },
  
  EU: {
    nav: {
      howItWorks: "How it works",
      benefits: "Benefits",
      freeEvaluation: "Free Evaluation"
    },
    hero: {
      title: "Optimize your fleet in",
      subtitle: "Optimize routes, generate driver links and monitor deliveries in real time.",
      freeEvaluation: "Free Evaluation",
      personalizedConsultation: "Personalized Consultation",
      vehicles: "Vehicles",
      deliveries: "Deliveries",
      routes: "Routes",
      vehicleDetails: {
        license: "License",
        weight: "Weight",
        volume: "Vol",
        insurance: "Insurance",
        additionalVars: "Additional variables...",
        moreVehicles: "+ 3 more vehicles..."
      },
      deliveryDetails: {
        client: "Client",
        address: "Address",
        weight: "Weight",
        volume: "Vol",
        price: "Price",
        additionalVars: "Additional variables...",
        moreDeliveries: "+ 8 more deliveries..."
      },
      routeDetails: {
        route: "Route",
        assignedLicense: "Assigned license",
        deliveries: "deliveries",
        additionalVars: "Additional variables...",
        moreRoutes: "+ 2 more routes..."
      }
    },
    howItWorks: {
      title: "How it works?",
      subtitle: "In just 3 simple steps, optimize your logistics and reduce costs",
      step1: {
        title: "Create your vehicles",
        description: "Configure and register your vehicles with their capabilities, dimensions and characteristics. Define the fleet you will use for your deliveries."
      },
      step2: {
        title: "Load the deliveries",
        description: "Create and configure your deliveries with destinations, products and restrictions. Define all the necessary details for route planning."
      },
      step3: {
        title: "Optimize and execute",
        description: "Call our optimization agent and receive the links for the driver. Execute optimized routes with real-time traceability."
      }
    },
    benefits: {
      title: "Benefits for companies",
      subtitle: "Optimize your logistics operation and reduce costs significantly",
      reduceCosts: {
        title: "Reduce logistics costs",
        description: "Optimize routes and reduce fuel, time and resources."
      },
      balanceLoad: {
        title: "Automatically balance load",
        description: "Distributes the load optimally between vehicles, maximizing efficiency and minimizing empty trips"
      },
      saveTime: {
        title: "Save time planning",
        description: "From 2 hours to 5 minutes. Plan routes automatically"
      },
      realTimeTracking: {
        title: "Real-time traceability",
        description: "Monitor delivery progress in real time from the mobile web, without the need for additional software"
      }
    },
    cta: {
      title: "Ready to optimize your fleet?",
      subtitle: "Schedule a free evaluation of your logistics processes",
      freeEvaluation: "Free Evaluation",
      personalizedConsultation: "Personalized Consultation"
    },
    demo: {
      title: "Try the App Live",
      subtitle: "Experience how the app works from the driver's perspective with test data",
      routeId: "Route ID",
      interactiveDemo: "Interactive Demo",
      features: {
        startRoute: {
          title: "Start a route",
          description: "Simulate starting a route with license plate"
        },
        markDeliveries: {
          title: "Mark deliveries",
          description: "Simulate the delivery process with evidence"
        },
        navigateMaps: {
          title: "Navigate with maps",
          description: "Integration with Google Maps and Waze"
        },
        generateReports: {
          title: "Generate reports",
          description: "Download reports in CSV and Excel"
        }
      },
      buttons: {
        viewDemo: "View Demo",
        hideDemo: "Hide Demo",
        openNewTab: "Open in New Tab"
      },
      preview: {
        title: "Interactive Demo",
        description: "Click \"View Demo\" to experience the application with simulated data",
        features: {
          simulatedDeliveries: "9 simulated deliveries",
          fullFunctionality: "Full functionality",
          realisticData: "Realistic Santiago data"
        }
      }
    },
    footer: {
      description: "Optimize routes, generate driver links and monitor deliveries in real time. Everything from a single platform.",
      product: "Product",
      support: "Support",
      company: "Company",
      howItWorks: "How it works",
      prices: "Prices",
      demo: "Demo",
      help: "Help",
      contact: "Contact",
      documentation: "Documentation",
      about: "About",
      blog: "Blog",
      careers: "Careers",
      copyright: "© 2025 TransportApp. All rights reserved."
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
