export type Language = 'EN' | 'ES' | 'PT'

export interface Translations {
  // Left Panel
  leftPanel: {
    title: string
    subtitle: string
    description: string
    menuControl: string
    templatesQR: string
    available24_7: string
    forYourClients: string
  }
  
  // Right Panel - Login
  login: {
    title: string
    signInWithGoogle: string
    connecting: string
    secureAccess: string
    secureAccessDescription: string
    contactSupport: string
    havingIssues: string
    copyright: string
    termsOfUse: string
    privacyPolicy: string
  }
  
  // Callback Page
  callback: {
    processing: string
    signingIn: string
    success: string
    redirecting: string
    welcome: string
    noSession: string
    error: string
    secureConnection: string
  }
  
  // Error Messages
  errors: {
    popupBlocked: string
    networkError: string
    authError: string
    processingError: string
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
    leftPanel: {
      title: 'MiCartaPro',
      subtitle: 'Your Digital Menu, Without Complications',
      description: 'Manage your digital menu and let sales flow.',
      menuControl: 'Menu control panel',
      templatesQR: 'Templates and QR codes',
      available24_7: '24/7',
      forYourClients: 'Available always for your clients'
    },
    login: {
      title: 'Sign In',
      signInWithGoogle: 'Sign in with Google',
      connecting: 'Connecting...',
      secureAccess: 'Secure Access',
      secureAccessDescription: 'Use your Google account to securely access your control panel',
      contactSupport: 'Contact support',
      havingIssues: 'Having issues with your account?',
      copyright: '© 2024 MiCartaPro. All rights reserved.',
      termsOfUse: 'Terms of Use',
      privacyPolicy: 'Privacy Policy'
    },
    callback: {
      processing: 'Processing authentication...',
      signingIn: 'Signing In',
      success: 'Authentication Successful!',
      redirecting: 'You will be redirected automatically...',
      welcome: 'Welcome {email}! Redirecting...',
      noSession: 'No session found. Redirecting...',
      error: 'Error: {message}',
      secureConnection: 'Connecting securely • SSL/TLS Encrypted'
    },
    errors: {
      popupBlocked: 'Popup window blocked. Please allow popups.',
      networkError: 'Connection error. Check your internet connection.',
      authError: 'Authentication error',
      processingError: 'Error processing authentication'
    }
  },
  ES: {
    leftPanel: {
      title: 'MiCartaPro',
      subtitle: 'Tu Menú Digital, Sin Complicaciones',
      description: 'Gestiona tu menú digital y deja que las ventas fluyan.',
      menuControl: 'Panel de control de tu menú',
      templatesQR: 'Plantillas y códigos QR',
      available24_7: '24/7',
      forYourClients: 'Disponible siempre para tus clientes'
    },
    login: {
      title: 'Iniciar Sesión',
      signInWithGoogle: 'Iniciar sesión con Google',
      connecting: 'Conectando...',
      secureAccess: 'Acceso Seguro',
      secureAccessDescription: 'Utiliza tu cuenta de Google para acceder de forma segura a tu panel de control',
      contactSupport: 'Contactar soporte',
      havingIssues: '¿Tienes problemas con tu cuenta?',
      copyright: '© 2024 MiCartaPro. Todos los derechos reservados.',
      termsOfUse: 'Términos de Uso',
      privacyPolicy: 'Política de Privacidad'
    },
    callback: {
      processing: 'Procesando autenticación...',
      signingIn: 'Iniciando Sesión',
      success: '¡Autenticación Exitosa!',
      redirecting: 'Serás redirigido automáticamente...',
      welcome: '¡Bienvenido {email}! Redirigiendo...',
      noSession: 'No se encontró sesión. Redirigiendo...',
      error: 'Error: {message}',
      secureConnection: 'Conectando de forma segura • SSL/TLS Encriptado'
    },
    errors: {
      popupBlocked: 'Ventana emergente bloqueada. Por favor, permite ventanas emergentes.',
      networkError: 'Error de conexión. Verifica tu conexión a internet.',
      authError: 'Error de autenticación',
      processingError: 'Error procesando autenticación'
    }
  },
  PT: {
    leftPanel: {
      title: 'MiCartaPro',
      subtitle: 'Seu Cardápio Digital, Sem Complicações',
      description: 'Gerencie seu cardápio digital e deixe as vendas fluírem.',
      menuControl: 'Painel de controle do seu cardápio',
      templatesQR: 'Modelos e códigos QR',
      available24_7: '24/7',
      forYourClients: 'Disponível sempre para seus clientes'
    },
    login: {
      title: 'Entrar',
      signInWithGoogle: 'Entrar com Google',
      connecting: 'Conectando...',
      secureAccess: 'Acesso Seguro',
      secureAccessDescription: 'Use sua conta do Google para acessar com segurança seu painel de controle',
      contactSupport: 'Contatar suporte',
      havingIssues: 'Está tendo problemas com sua conta?',
      copyright: '© 2024 MiCartaPro. Todos os direitos reservados.',
      termsOfUse: 'Termos de Uso',
      privacyPolicy: 'Política de Privacidade'
    },
    callback: {
      processing: 'Processando autenticação...',
      signingIn: 'Entrando',
      success: 'Autenticação Bem-Sucedida!',
      redirecting: 'Você será redirecionado automaticamente...',
      welcome: 'Bem-vindo {email}! Redirecionando...',
      noSession: 'Nenhuma sessão encontrada. Redirecionando...',
      error: 'Erro: {message}',
      secureConnection: 'Conectando com segurança • SSL/TLS Criptografado'
    },
    errors: {
      popupBlocked: 'Janela popup bloqueada. Por favor, permita popups.',
      networkError: 'Erro de conexão. Verifique sua conexão com a internet.',
      authError: 'Erro de autenticação',
      processingError: 'Erro ao processar autenticação'
    }
  }
}

