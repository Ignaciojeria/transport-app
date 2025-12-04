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
    saasDescription: string
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
      additionalServices: {
        title: string
        subtitle: string
        description: string
        startingFrom: string
        requiresQuote: string
        note: string
        button: string
        servicesTitle: string
        services: {
          customDesign: string
          menuSetup: string
          consulting: string
          migration: string
        }
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
      subtitle: 'Self-service SaaS platform for managing your digital menu',
      customDesign: '🎨 Customizable design options',
      customLogo: '✨ Upload your own logo',
      exclusiveQR: '🔗 Generate your QR code',
      responsiveDesign: '📱 Responsive design for all devices',
      pricingTitle: 'Simple monthly subscription — $15 USD per month',
      firstYearFree: '',
      renewalPrice: 'Cancel anytime, no long-term commitment',
      quoteButton: 'Get Started',
      viewPricing: 'View Full Pricing',
      saasDescription: 'MiCartaPro is a self-service SaaS platform that lets restaurants manage their digital menu. All plans include continuous improvements, new features, performance updates, and ongoing enhancements to the platform.'
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
      startingFrom: '',
      oneTimePayment: '',
      firstYear: '',
      firstYearFree: '',
      renewal: '',
      monthly: 'per month',
      fromSecondYear: '',
      features: {
        customDesign: 'Customizable Design Options',
        customLogo: 'Upload Your Own Logo',
        exclusiveQR: 'Generate Your QR Code',
        responsiveDesign: 'Responsive Design for All Devices',
        shoppingCart: 'Integrated Shopping Cart',
        whatsappIntegration: 'WhatsApp Integration',
        costCalculation: 'Automatic Cost Calculation',
        deliveryOptions: 'Delivery or In-Store Pickup',
        support: 'Platform Support & Updates'
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
            question: 'What is included in the subscription?',
            answer: 'Your $15 USD monthly subscription includes full access to the MiCartaPro self-service SaaS platform, where you can customize your digital menu design, upload your logo, generate your QR code, and manage your menu. All plans include continuous improvements, new features, and platform updates.'
          },
          {
            question: 'When do I start paying?',
            answer: 'You start paying $15 USD per month from the moment you subscribe. There are no setup fees or long-term commitments. You can cancel anytime.'
          },
          {
            question: 'Can I cancel my subscription?',
            answer: 'Yes, you can cancel your subscription at any time. The cancellation will take effect at the end of your current billing period. We offer a 14-day money-back guarantee on all purchases made via Paddle. Please see our Refund Policy for more details.'
          },
          {
            question: 'What is your refund policy?',
            answer: 'We offer a 14-day money-back guarantee on all purchases made via Paddle. If you\'re not satisfied with your purchase, you can request a full refund within 14 days from the date of purchase by contacting Paddle or our support team. After 14 days, refunds are handled at our and Paddle\'s discretion.'
          },
          {
            question: 'What payment methods do you accept?',
            answer: 'All payments are processed by Paddle, who acts as the Merchant of Record. Paddle accepts various payment methods including credit cards and other standard payment options.'
          },
          {
            question: 'Is there a setup fee?',
            answer: 'No, there are no setup fees. Your $15 USD monthly subscription includes everything you need to get started immediately.'
          }
        ]
      },
      additionalServices: {
        title: 'Add-On Services (Optional)',
        subtitle: 'Professional services billed separately',
        description: 'If you need professional custom design work, manual menu creation, or personalized consulting, we offer these services as optional add-ons. These services are billed separately and are NOT processed through Paddle.',
        startingFrom: 'Starting from $150 USD',
        requiresQuote: 'Requires a custom quote',
        note: 'These add-on services are optional, billed separately via bank transfer, PayPal, or invoice, and are NOT included in your SaaS subscription. They are NOT processed through Paddle. Contact us for a personalized quote based on your specific needs.',
        button: 'Request Custom Quote',
        servicesTitle: 'Available Add-On Services:',
        services: {
          customDesign: 'Custom Menu Design',
          menuSetup: 'Manual Menu Setup',
          consulting: 'Personalized Consulting',
          migration: 'Migration Support'
        }
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
      subtitle: 'Plataforma SaaS de autoservicio para gestionar tu menú digital',
      customDesign: '🎨 Opciones de diseño personalizables',
      customLogo: '✨ Sube tu propio logo',
      exclusiveQR: '🔗 Genera tu código QR',
      responsiveDesign: '📱 Diseño responsivo para todos los dispositivos',
      pricingTitle: 'Suscripción mensual simple — $15 USD por mes',
      firstYearFree: '',
      renewalPrice: 'Cancela en cualquier momento, sin compromiso a largo plazo',
      quoteButton: 'Comenzar',
      viewPricing: 'Ver Precios Completos',
      saasDescription: 'MiCartaPro es una plataforma SaaS de autoservicio que permite a los restaurantes gestionar su menú digital. Todos los planes incluyen mejoras continuas, nuevas funcionalidades, actualizaciones de rendimiento y mejoras constantes de la plataforma.'
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
      startingFrom: '',
      oneTimePayment: '',
      firstYear: '',
      firstYearFree: '',
      renewal: '',
      monthly: 'por mes',
      fromSecondYear: '',
      features: {
        customDesign: 'Opciones de Diseño Personalizables',
        customLogo: 'Sube Tu Propio Logo',
        exclusiveQR: 'Genera Tu Código QR',
        responsiveDesign: 'Diseño Responsivo para Todos los Dispositivos',
        shoppingCart: 'Carrito de Compras Integrado',
        whatsappIntegration: 'Integración con WhatsApp',
        costCalculation: 'Cálculo Automático de Costos',
        deliveryOptions: 'Envío o Retiro en Tienda',
        support: 'Soporte de Plataforma y Actualizaciones'
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
            question: '¿Qué está incluido en la suscripción?',
            answer: 'Tu suscripción mensual de $15 USD incluye acceso completo a la plataforma SaaS de autoservicio MiCartaPro, donde puedes personalizar el diseño de tu menú digital, subir tu logo, generar tu código QR y gestionar tu menú. Todos los planes incluyen mejoras continuas, nuevas funcionalidades y actualizaciones de la plataforma.'
          },
          {
            question: '¿Cuándo empiezo a pagar?',
            answer: 'Empiezas a pagar $15 USD por mes desde el momento en que te suscribes. No hay tarifas de configuración ni compromisos a largo plazo. Puedes cancelar en cualquier momento.'
          },
          {
            question: '¿Puedo cancelar mi suscripción?',
            answer: 'Sí, puedes cancelar tu suscripción en cualquier momento. La cancelación tendrá efecto al final de tu período de facturación actual. Ofrecemos una garantía de reembolso de 14 días en todas las compras realizadas a través de Paddle. Por favor, consulta nuestra Política de Reembolso para más detalles.'
          },
          {
            question: '¿Cuál es su política de reembolso?',
            answer: 'Ofrecemos una garantía de reembolso de 14 días en todas las compras realizadas a través de Paddle. Si no estás satisfecho con tu compra, puedes solicitar un reembolso completo dentro de los 14 días posteriores a la fecha de compra contactando a Paddle o a nuestro equipo de soporte. Después de 14 días, los reembolsos se manejan a nuestra discreción y la de Paddle.'
          },
          {
            question: '¿Qué métodos de pago aceptan?',
            answer: 'Todos los pagos son procesados por Paddle, quien actúa como Merchant of Record. Paddle acepta varios métodos de pago incluyendo tarjetas de crédito y otras opciones de pago estándar.'
          },
          {
            question: '¿Hay una tarifa de configuración?',
            answer: 'No, no hay tarifas de configuración. Tu suscripción mensual de $15 USD incluye todo lo que necesitas para comenzar inmediatamente.'
          }
        ]
      },
      additionalServices: {
        title: 'Servicios Adicionales (Opcionales)',
        subtitle: 'Servicios profesionales facturados por separado',
        description: 'Si necesitas trabajo de diseño personalizado profesional, creación manual de menús o consultoría personalizada, ofrecemos estos servicios como complementos opcionales. Estos servicios se facturan por separado y NO se procesan a través de Paddle.',
        startingFrom: 'Desde $150 USD',
        requiresQuote: 'Requiere cotización personalizada',
        note: 'Estos servicios adicionales son opcionales, se facturan por separado mediante transferencia bancaria, PayPal o factura, y NO están incluidos en tu suscripción SaaS. NO se procesan a través de Paddle. Contáctanos para una cotización personalizada según tus necesidades específicas.',
        button: 'Solicitar Cotización Personalizada',
        servicesTitle: 'Servicios Adicionales Disponibles:',
        services: {
          customDesign: 'Diseño Personalizado de Menú',
          menuSetup: 'Configuración Manual de Menú',
          consulting: 'Consultoría Personalizada',
          migration: 'Soporte de Migración'
        }
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
      subtitle: 'Plataforma SaaS de autoatendimento para gerenciar seu cardápio digital',
      customDesign: '🎨 Opções de design personalizáveis',
      customLogo: '✨ Envie seu próprio logo',
      exclusiveQR: '🔗 Gere seu código QR',
      responsiveDesign: '📱 Design responsivo para todos os dispositivos',
      pricingTitle: 'Assinatura mensal simples — $15 USD por mês',
      firstYearFree: '',
      renewalPrice: 'Cancele a qualquer momento, sem compromisso de longo prazo',
      quoteButton: 'Começar',
      viewPricing: 'Ver Preços Completos',
      saasDescription: 'MiCartaPro é uma plataforma SaaS de autoatendimento que permite aos restaurantes gerenciar seu cardápio digital. Todos os planos incluem melhorias contínuas, novos recursos, atualizações de desempenho e aprimoramentos constantes da plataforma.'
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
      startingFrom: '',
      oneTimePayment: '',
      firstYear: '',
      firstYearFree: '',
      renewal: '',
      monthly: 'por mês',
      fromSecondYear: '',
      features: {
        customDesign: 'Opções de Design Personalizáveis',
        customLogo: 'Envie Seu Próprio Logo',
        exclusiveQR: 'Gere Seu Código QR',
        responsiveDesign: 'Design Responsivo para Todos os Dispositivos',
        shoppingCart: 'Carrinho de Compras Integrado',
        whatsappIntegration: 'Integração com WhatsApp',
        costCalculation: 'Cálculo Automático de Custos',
        deliveryOptions: 'Entrega ou Retirada na Loja',
        support: 'Suporte da Plataforma e Atualizações'
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
            question: 'O que está incluído na assinatura?',
            answer: 'Sua assinatura mensal de $15 USD inclui acesso completo à plataforma SaaS de autoatendimento MiCartaPro, onde você pode personalizar o design do seu cardápio digital, enviar seu logo, gerar seu código QR e gerenciar seu cardápio. Todos os planos incluem melhorias contínuas, novos recursos e atualizações da plataforma.'
          },
          {
            question: 'Quando começo a pagar?',
            answer: 'Você começa a pagar $15 USD por mês a partir do momento em que se inscreve. Não há taxas de configuração nem compromissos de longo prazo. Você pode cancelar a qualquer momento.'
          },
          {
            question: 'Posso cancelar minha assinatura?',
            answer: 'Sim, você pode cancelar sua assinatura a qualquer momento. O cancelamento terá efeito no final do seu período de faturamento atual. Oferecemos uma garantia de reembolso de 14 dias em todas as compras realizadas através da Paddle. Por favor, consulte nossa Política de Reembolso para mais detalhes.'
          },
          {
            question: 'Qual é sua política de reembolso?',
            answer: 'Oferecemos uma garantia de reembolso de 14 dias em todas as compras realizadas através da Paddle. Se você não estiver satisfeito com sua compra, pode solicitar um reembolso completo dentro de 14 dias a partir da data de compra entrando em contato com a Paddle ou nossa equipe de suporte. Após 14 dias, os reembolsos são tratados a nosso critério e da Paddle.'
          },
          {
            question: 'Quais métodos de pagamento vocês aceitam?',
            answer: 'Todos os pagamentos são processados pela Paddle, que atua como Merchant of Record. A Paddle aceita vários métodos de pagamento, incluindo cartões de crédito e outras opções de pagamento padrão.'
          },
          {
            question: 'Há uma taxa de configuração?',
            answer: 'Não, não há taxas de configuração. Sua assinatura mensal de $15 USD inclui tudo que você precisa para começar imediatamente.'
          }
        ]
      },
      additionalServices: {
        title: 'Serviços Adicionais (Opcionais)',
        subtitle: 'Serviços profissionais faturados separadamente',
        description: 'Se você precisar de trabalho de design personalizado profissional, criação manual de cardápios ou consultoria personalizada, oferecemos esses serviços como complementos opcionais. Esses serviços são faturados separadamente e NÃO são processados através da Paddle.',
        startingFrom: 'A partir de $150 USD',
        requiresQuote: 'Requer cotização personalizada',
        note: 'Esses serviços adicionais são opcionais, são faturados separadamente via transferência bancária, PayPal ou fatura, e NÃO estão incluídos na sua assinatura SaaS. NÃO são processados através da Paddle. Entre em contato conosco para uma cotização personalizada com base em suas necessidades específicas.',
        button: 'Solicitar Cotização Personalizada',
        servicesTitle: 'Serviços Adicionais Disponíveis:',
        services: {
          customDesign: 'Design Personalizado de Cardápio',
          menuSetup: 'Configuração Manual de Cardápio',
          consulting: 'Consultoria Personalizada',
          migration: 'Suporte de Migração'
        }
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

