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
    viewPricing: string
  }
  
  // Pricing Section
  pricing: {
    title: string
    subtitle: string
    startingFrom: string
    oneTimePayment: string
    firstYear: string
    firstYearFree: string
    renewal: string
    monthly: string
    fromSecondYear: string
    features: {
      customDesign: string
      customLogo: string
      exclusiveQR: string
      responsiveDesign: string
      shoppingCart: string
      whatsappIntegration: string
      costCalculation: string
      deliveryOptions: string
      support: string
    }
    cta: {
      title: string
      subtitle: string
      quoteButton: string
    }
    faq: {
      title: string
      questions: Array<{
        question: string
        answer: string
      }>
    }
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
      quoteButton: 'Quote Now',
      viewPricing: 'View Full Pricing'
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
    pricing: {
      title: 'Simple, Transparent Pricing',
      subtitle: 'Choose the perfect plan for your restaurant',
      startingFrom: 'Starting from',
      oneTimePayment: 'One-time payment',
      firstYear: 'First Year',
      firstYearFree: 'FREE',
      renewal: 'Renewal',
      monthly: 'per month',
      fromSecondYear: 'From the second year',
      features: {
        customDesign: '100% Custom Design',
        customLogo: 'Custom Logo',
        exclusiveQR: 'Exclusive QR Code',
        responsiveDesign: 'Responsive Design for All Devices',
        shoppingCart: 'Integrated Shopping Cart',
        whatsappIntegration: 'WhatsApp Integration',
        costCalculation: 'Automatic Cost Calculation',
        deliveryOptions: 'Delivery or In-Store Pickup',
        support: 'Customer Support'
      },
      cta: {
        title: 'Ready to Get Started?',
        subtitle: 'Contact us now and get your personalized digital menu',
        quoteButton: 'Get a Quote'
      },
      faq: {
        title: 'Frequently Asked Questions',
        questions: [
          {
            question: 'What is included in the initial payment?',
            answer: 'The initial payment of $150 USD includes the complete setup of your digital menu, custom design, logo integration, QR code generation, and the first year of service completely free.'
          },
          {
            question: 'When do I start paying the monthly fee?',
            answer: 'The monthly fee of $10 USD starts from the second year. The first year is completely free as part of our promotional offer.'
          },
          {
            question: 'Can I cancel my subscription?',
            answer: 'Yes, you can cancel your subscription at any time. The cancellation will take effect at the end of your current billing period. Please see our Refund Policy for more details.'
          },
          {
            question: 'What payment methods do you accept?',
            answer: 'We accept various payment methods. Contact us via WhatsApp to discuss payment options that work best for you.'
          },
          {
            question: 'Is there a setup fee?',
            answer: 'The initial payment of $150 USD covers the complete setup and first year of service. There are no additional setup fees.'
          }
        ]
      }
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
      quoteButton: 'Cotizar Ahora',
      viewPricing: 'Ver Precios Completos'
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
    pricing: {
      title: 'Precios Simples y Transparentes',
      subtitle: 'Elige el plan perfecto para tu restaurante',
      startingFrom: 'Desde',
      oneTimePayment: 'Pago único',
      firstYear: 'Primer Año',
      firstYearFree: 'GRATIS',
      renewal: 'Renovación',
      monthly: 'por mes',
      fromSecondYear: 'Desde el segundo año',
      features: {
        customDesign: 'Diseño 100% Personalizado',
        customLogo: 'Logo a Medida',
        exclusiveQR: 'Código QR Exclusivo',
        responsiveDesign: 'Diseño Responsivo para Todos los Dispositivos',
        shoppingCart: 'Carrito de Compras Integrado',
        whatsappIntegration: 'Integración con WhatsApp',
        costCalculation: 'Cálculo Automático de Costos',
        deliveryOptions: 'Envío o Retiro en Tienda',
        support: 'Soporte al Cliente'
      },
      cta: {
        title: '¿Listo para Empezar?',
        subtitle: 'Contáctanos ahora y obtén tu menú digital personalizado',
        quoteButton: 'Obtener Cotización'
      },
      faq: {
        title: 'Preguntas Frecuentes',
        questions: [
          {
            question: '¿Qué está incluido en el pago inicial?',
            answer: 'El pago inicial de $150 USD incluye la configuración completa de tu menú digital, diseño personalizado, integración de logo, generación de código QR y el primer año de servicio completamente gratis.'
          },
          {
            question: '¿Cuándo empiezo a pagar la tarifa mensual?',
            answer: 'La tarifa mensual de $10 USD comienza desde el segundo año. El primer año es completamente gratis como parte de nuestra oferta promocional.'
          },
          {
            question: '¿Puedo cancelar mi suscripción?',
            answer: 'Sí, puedes cancelar tu suscripción en cualquier momento. La cancelación tendrá efecto al final de tu período de facturación actual. Por favor, consulta nuestra Política de Reembolso para más detalles.'
          },
          {
            question: '¿Qué métodos de pago aceptan?',
            answer: 'Aceptamos varios métodos de pago. Contáctanos por WhatsApp para discutir las opciones de pago que mejor funcionen para ti.'
          },
          {
            question: '¿Hay una tarifa de configuración?',
            answer: 'El pago inicial de $150 USD cubre la configuración completa y el primer año de servicio. No hay tarifas de configuración adicionales.'
          }
        ]
      }
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
      quoteButton: 'Cotizar Agora',
      viewPricing: 'Ver Preços Completos'
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
    pricing: {
      title: 'Preços Simples e Transparentes',
      subtitle: 'Escolha o plano perfeito para seu restaurante',
      startingFrom: 'A partir de',
      oneTimePayment: 'Pagamento único',
      firstYear: 'Primeiro Ano',
      firstYearFree: 'GRÁTIS',
      renewal: 'Renovação',
      monthly: 'por mês',
      fromSecondYear: 'A partir do segundo ano',
      features: {
        customDesign: 'Design 100% Personalizado',
        customLogo: 'Logo Sob Medida',
        exclusiveQR: 'Código QR Exclusivo',
        responsiveDesign: 'Design Responsivo para Todos os Dispositivos',
        shoppingCart: 'Carrinho de Compras Integrado',
        whatsappIntegration: 'Integração com WhatsApp',
        costCalculation: 'Cálculo Automático de Custos',
        deliveryOptions: 'Entrega ou Retirada na Loja',
        support: 'Suporte ao Cliente'
      },
      cta: {
        title: 'Pronto para Começar?',
        subtitle: 'Entre em contato conosco agora e obtenha seu cardápio digital personalizado',
        quoteButton: 'Obter Cotização'
      },
      faq: {
        title: 'Perguntas Frequentes',
        questions: [
          {
            question: 'O que está incluído no pagamento inicial?',
            answer: 'O pagamento inicial de $150 USD inclui a configuração completa do seu cardápio digital, design personalizado, integração de logo, geração de código QR e o primeiro ano de serviço completamente grátis.'
          },
          {
            question: 'Quando começo a pagar a taxa mensal?',
            answer: 'A taxa mensal de $10 USD começa a partir do segundo ano. O primeiro ano é completamente grátis como parte da nossa oferta promocional.'
          },
          {
            question: 'Posso cancelar minha assinatura?',
            answer: 'Sim, você pode cancelar sua assinatura a qualquer momento. O cancelamento terá efeito no final do seu período de faturamento atual. Por favor, consulte nossa Política de Reembolso para mais detalhes.'
          },
          {
            question: 'Quais métodos de pagamento vocês aceitam?',
            answer: 'Aceitamos vários métodos de pagamento. Entre em contato conosco via WhatsApp para discutir as opções de pagamento que melhor funcionem para você.'
          },
          {
            question: 'Há uma taxa de configuração?',
            answer: 'O pagamento inicial de $150 USD cobre a configuração completa e o primeiro ano de serviço. Não há taxas de configuração adicionais.'
          }
        ]
      }
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

