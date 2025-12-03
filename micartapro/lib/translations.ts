export type Language = 'EN' | 'ES' | 'PT'

export interface Translations {
  // Navigation
  nav: {
    service: string
    benefits: string
    demo: string
    quote: string
  }
  
  // Hero Section
  hero: {
    title: string
    subtitle: string
    description: string
    quoteNow: string
    viewDemo: string
  }
  
  // Service Section
  service: {
    title: string
    subtitle: string
    customDesign: string
    customLogo: string
    exclusiveQR: string
    responsiveDesign: string
    pricingTitle: string
    firstYearFree: string
    renewalPrice: string
    quoteButton: string
  }
  
  // Benefits Section
  benefits: {
    title: string
    autoCostCalculation: {
      title: string
      description: string
    }
    shoppingCart: {
      title: string
      description: string
    }
    deliveryOptions: {
      title: string
      description: string
    }
    whatsappOrders: {
      title: string
      description: string
    }
  }
  
  // Demo Section
  demo: {
    title: string
    subtitle: string
    interactiveDemo: string
    selectDish: {
      title: string
      description: string
    }
    chooseSide: {
      title: string
      description: string
    }
    enterName: {
      title: string
      description: string
    }
    selectTime: {
      title: string
      description: string
    }
    sendWhatsApp: {
      title: string
      description: string
    }
    viewDemo: string
    hideDemo: string
    openNewTab: string
    preview: {
      title: string
      description: string
      features: {
        digitalMenu: string
        shoppingCart: string
        whatsapp: string
      }
    }
  }
  
  // CTA Section
  cta: {
    title: string
    subtitle: string
    quoteButton: string
  }
  
  // Footer
  footer: {
    description: string
    service: string
    ourService: string
    benefits: string
    demo: string
    contact: string
    quoteWhatsApp: string
    privacy: string
    terms: string
    refund: string
    copyright: string
  }
}

export const languageNames: Record<Language, string> = {
  EN: 'English',
  ES: 'Español',
  PT: 'Português'
}

export const languageFlags: Record<Language, string> = {
  EN: '🇬🇧',
  ES: '🇪🇸',
  PT: '🇧🇷'
}

export const translations: Record<Language, Translations> = {
  EN: {
    nav: {
      service: 'Service',
      benefits: 'Benefits',
      demo: 'Demo',
      quote: 'Quote'
    },
    hero: {
      title: 'MiCartaPro',
      subtitle: 'Your Digital Menu, Without Complications',
      description: 'Manage your digital menu and let sales flow.',
      quoteNow: 'Quote Now',
      viewDemo: 'View Demo'
    },
    service: {
      title: '🚀 Our Service',
      subtitle: 'At MiCartaPro we do it differently',
      customDesign: '🎨 100% custom design',
      customLogo: '✨ Custom logo',
      exclusiveQR: '🔗 Exclusive QR code',
      responsiveDesign: '📱 Responsive design for all devices',
      pricingTitle: 'Limited-time offer — starting from $150 USD',
      firstYearFree: '✅ First year free',
      renewalPrice: 'Renewal from the second year: $10 USD monthly',
      quoteButton: 'Quote Now'
    },
    benefits: {
      title: '🎯 Included Benefits',
      autoCostCalculation: {
        title: '💰 Automatic cost calculation',
        description: 'Forget manual calculations. Your menu processes and displays the total cost of each dish.'
      },
      shoppingCart: {
        title: '🛒 Integrated shopping cart',
        description: 'Allow your customers to build their order simply, organized, and quickly.'
      },
      deliveryOptions: {
        title: '🚚 Delivery or in-store pickup',
        description: 'Your menu automatically asks for the necessary details to complete the order.'
      },
      whatsappOrders: {
        title: '📩 Order reception via WhatsApp',
        description: 'Receive orders in an organized, clear, and transparent way for both the kitchen and your customers.'
      }
    },
    demo: {
      title: '🧪 Try MiCartaPro live',
      subtitle: 'Experience the app as if you were a real customer',
      interactiveDemo: 'Interactive Demo',
      selectDish: {
        title: 'Select a dish',
        description: 'Explore the menu and choose your favorite dish'
      },
      chooseSide: {
        title: 'Choose the side',
        description: 'Customize your order with available sides'
      },
      enterName: {
        title: 'Enter the name of who will pick up',
        description: 'Complete the necessary details for pickup'
      },
      selectTime: {
        title: 'Select pickup time (24-hour format)',
        description: 'Indicate when you want to pick up your order'
      },
      sendWhatsApp: {
        title: 'Send your order directly via WhatsApp',
        description: 'The order is automatically sent to the restaurant'
      },
      viewDemo: 'View Demo',
      hideDemo: 'Hide Demo',
      openNewTab: 'Open in new tab',
      preview: {
        title: 'Preview',
        description: 'Click "View Demo" to experience the full application',
        features: {
          digitalMenu: 'Interactive digital menu',
          shoppingCart: 'Functional shopping cart',
          whatsapp: 'WhatsApp integration'
        }
      }
    },
    cta: {
      title: 'Ready to digitize your menu?',
      subtitle: 'Contact us now and get your personalized digital menu',
      quoteButton: 'Quote Now'
    },
    footer: {
      description: 'Your digital menu, without complications. Manage your restaurant and let sales flow.',
      service: 'Service',
      ourService: 'Our Service',
      benefits: 'Benefits',
      demo: 'Demo',
      contact: 'Contact',
      quoteWhatsApp: 'Quote via WhatsApp',
      privacy: 'Privacy Policy',
      terms: 'Terms and Conditions',
      refund: 'Refund Policy',
      copyright: 'All rights reserved.'
    }
  },
  ES: {
    nav: {
      service: 'Servicio',
      benefits: 'Beneficios',
      demo: 'Demo',
      quote: 'Cotizar'
    },
    hero: {
      title: 'MiCartaPro',
      subtitle: 'Tu Menú Digital, Sin Complicaciones',
      description: 'Gestiona tu menú digital y deja que las ventas fluyan.',
      quoteNow: 'Cotizar Ahora',
      viewDemo: 'Ver Demo'
    },
    service: {
      title: '🚀 Nuestro Servicio',
      subtitle: 'En MiCartaPro lo hacemos diferente',
      customDesign: '🎨 Diseño 100% personalizado',
      customLogo: '✨ Logo a medida',
      exclusiveQR: '🔗 Código QR exclusivo',
      responsiveDesign: '📱 Diseño responsivo para todos los dispositivos',
      pricingTitle: 'Oferta única con cupos limitados — desde $150 USD',
      firstYearFree: '✅ Primer año gratis',
      renewalPrice: 'Renovación desde el segundo año: $10 USD mensuales',
      quoteButton: 'Cotizar Ahora'
    },
    benefits: {
      title: '🎯 Beneficios incluidos',
      autoCostCalculation: {
        title: '💰 Cálculo de costos automático',
        description: 'Olvídate del cálculo manual. Tu carta procesa y muestra el costo total de cada plato.'
      },
      shoppingCart: {
        title: '🛒 Carrito de compras integrado',
        description: 'Permite que tus clientes armen su pedido de manera simple, ordenada y rápida.'
      },
      deliveryOptions: {
        title: '🚚 Envío o retiro en tienda',
        description: 'Tu carta pregunta automáticamente por los detalles necesarios para completar el pedido.'
      },
      whatsappOrders: {
        title: '📩 Recepción de pedidos por WhatsApp',
        description: 'Recibe los pedidos de forma ordenada, clara y transparente tanto para la cocina como para tus clientes.'
      }
    },
    demo: {
      title: '🧪 Prueba MiCartaPro en vivo',
      subtitle: 'Experimenta la app como si fueras un cliente real',
      interactiveDemo: 'Demo Interactiva',
      selectDish: {
        title: 'Selecciona un plato',
        description: 'Explora el menú y elige tu plato favorito'
      },
      chooseSide: {
        title: 'Elige el acompañamiento',
        description: 'Personaliza tu pedido con los acompañamientos disponibles'
      },
      enterName: {
        title: 'Ingresa el nombre de quien retira',
        description: 'Completa los datos necesarios para el retiro'
      },
      selectTime: {
        title: 'Selecciona la hora de retiro (formato 24 h)',
        description: 'Indica cuándo quieres retirar tu pedido'
      },
      sendWhatsApp: {
        title: 'Envía tu pedido directamente por WhatsApp',
        description: 'El pedido se envía automáticamente al restaurante'
      },
      viewDemo: 'Ver Demo',
      hideDemo: 'Ocultar Demo',
      openNewTab: 'Abrir en nueva pestaña',
      preview: {
        title: 'Vista Previa',
        description: 'Haz clic en "Ver Demo" para experimentar la aplicación completa',
        features: {
          digitalMenu: 'Menú digital interactivo',
          shoppingCart: 'Carrito de compras funcional',
          whatsapp: 'Integración con WhatsApp'
        }
      }
    },
    cta: {
      title: '¿Listo para digitalizar tu menú?',
      subtitle: 'Contáctanos ahora y obtén tu menú digital personalizado',
      quoteButton: 'Cotizar Ahora'
    },
    footer: {
      description: 'Tu menú digital, sin complicaciones. Gestiona tu restaurante y deja que las ventas fluyan.',
      service: 'Servicio',
      ourService: 'Nuestro Servicio',
      benefits: 'Beneficios',
      demo: 'Demo',
      contact: 'Contacto',
      quoteWhatsApp: 'Cotizar por WhatsApp',
      privacy: 'Política de Privacidad',
      terms: 'Términos y Condiciones',
      refund: 'Política de Reembolso',
      copyright: 'Todos los derechos reservados.'
    }
  },
  PT: {
    nav: {
      service: 'Serviço',
      benefits: 'Benefícios',
      demo: 'Demo',
      quote: 'Cotizar'
    },
    hero: {
      title: 'MiCartaPro',
      subtitle: 'Seu Cardápio Digital, Sem Complicações',
      description: 'Gerencie seu cardápio digital e deixe as vendas fluírem.',
      quoteNow: 'Cotizar Agora',
      viewDemo: 'Ver Demo'
    },
    service: {
      title: '🚀 Nosso Serviço',
      subtitle: 'Na MiCartaPro fazemos diferente',
      customDesign: '🎨 Design 100% personalizado',
      customLogo: '✨ Logo sob medida',
      exclusiveQR: '🔗 Código QR exclusivo',
      responsiveDesign: '📱 Design responsivo para todos os dispositivos',
      pricingTitle: 'Oferta única com vagas limitadas — a partir de $150 USD',
      firstYearFree: '✅ Primeiro ano grátis',
      renewalPrice: 'Renovação a partir do segundo ano: $10 USD mensais',
      quoteButton: 'Cotizar Agora'
    },
    benefits: {
      title: '🎯 Benefícios incluídos',
      autoCostCalculation: {
        title: '💰 Cálculo automático de custos',
        description: 'Esqueça o cálculo manual. Seu cardápio processa e exibe o custo total de cada prato.'
      },
      shoppingCart: {
        title: '🛒 Carrinho de compras integrado',
        description: 'Permita que seus clientes montem seu pedido de forma simples, organizada e rápida.'
      },
      deliveryOptions: {
        title: '🚚 Entrega ou retirada na loja',
        description: 'Seu cardápio pergunta automaticamente pelos detalhes necessários para completar o pedido.'
      },
      whatsappOrders: {
        title: '📩 Recepção de pedidos via WhatsApp',
        description: 'Receba os pedidos de forma organizada, clara e transparente tanto para a cozinha quanto para seus clientes.'
      }
    },
    demo: {
      title: '🧪 Experimente MiCartaPro ao vivo',
      subtitle: 'Experimente o aplicativo como se fosse um cliente real',
      interactiveDemo: 'Demo Interativa',
      selectDish: {
        title: 'Selecione um prato',
        description: 'Explore o cardápio e escolha seu prato favorito'
      },
      chooseSide: {
        title: 'Escolha o acompanhamento',
        description: 'Personalize seu pedido com os acompanhamentos disponíveis'
      },
      enterName: {
        title: 'Digite o nome de quem vai retirar',
        description: 'Complete os dados necessários para retirada'
      },
      selectTime: {
        title: 'Selecione o horário de retirada (formato 24h)',
        description: 'Indique quando deseja retirar seu pedido'
      },
      sendWhatsApp: {
        title: 'Envie seu pedido diretamente via WhatsApp',
        description: 'O pedido é enviado automaticamente para o restaurante'
      },
      viewDemo: 'Ver Demo',
      hideDemo: 'Ocultar Demo',
      openNewTab: 'Abrir em nova aba',
      preview: {
        title: 'Visualização',
        description: 'Clique em "Ver Demo" para experimentar o aplicativo completo',
        features: {
          digitalMenu: 'Cardápio digital interativo',
          shoppingCart: 'Carrinho de compras funcional',
          whatsapp: 'Integração com WhatsApp'
        }
      }
    },
    cta: {
      title: 'Pronto para digitalizar seu cardápio?',
      subtitle: 'Entre em contato conosco agora e obtenha seu cardápio digital personalizado',
      quoteButton: 'Cotizar Agora'
    },
    footer: {
      description: 'Seu cardápio digital, sem complicações. Gerencie seu restaurante e deixe as vendas fluírem.',
      service: 'Serviço',
      ourService: 'Nosso Serviço',
      benefits: 'Benefícios',
      demo: 'Demo',
      contact: 'Contato',
      quoteWhatsApp: 'Cotizar via WhatsApp',
      privacy: 'Política de Privacidade',
      terms: 'Termos e Condições',
      refund: 'Política de Reembolso',
      copyright: 'Todos os direitos reservados.'
    }
  }
}

