// Sistema de logging configurable para controlar logs en desarrollo
class Logger {
  private isDebugMode: boolean
  private isProduction: boolean

  constructor() {
    // En desarrollo, solo mostrar logs si está habilitado explícitamente
    this.isDebugMode = import.meta.env.DEV && localStorage.getItem('debug-logs') === 'true'
    this.isProduction = import.meta.env.PROD
  }

  // Método para habilitar/deshabilitar logs de debug
  setDebugMode(enabled: boolean) {
    this.isDebugMode = enabled
    localStorage.setItem('debug-logs', enabled.toString())
  }

  // Logs de debug (solo en modo debug)
  debug(...args: any[]) {
    if (this.isDebugMode) {
      console.log('🐛 [DEBUG]', ...args)
    }
  }

  // Logs de información (solo en desarrollo)
  info(...args: any[]) {
    if (!this.isProduction) {
      console.log('ℹ️ [INFO]', ...args)
    }
  }

  // Logs de warning (siempre visibles)
  warn(...args: any[]) {
    console.warn('⚠️ [WARN]', ...args)
  }

  // Logs de error (siempre visibles)
  error(...args: any[]) {
    console.error('❌ [ERROR]', ...args)
  }

  // Logs de éxito (solo en desarrollo)
  success(...args: any[]) {
    if (!this.isProduction) {
      console.log('✅ [SUCCESS]', ...args)
    }
  }

  // Logs de operaciones de base de datos (configurable)
  db(...args: any[]) {
    if (this.isDebugMode || localStorage.getItem('db-logs') === 'true') {
      console.log('🗄️ [DB]', ...args)
    }
  }

  // Logs de estado (configurable)
  state(...args: any[]) {
    if (this.isDebugMode || localStorage.getItem('state-logs') === 'true') {
      console.log('🔄 [STATE]', ...args)
    }
  }
}

// Instancia global del logger
export const logger = new Logger()

// Función helper para habilitar logs específicos
export const enableLogs = {
  debug: () => logger.setDebugMode(true),
  db: () => localStorage.setItem('db-logs', 'true'),
  state: () => localStorage.setItem('state-logs', 'true'),
  all: () => {
    logger.setDebugMode(true)
    localStorage.setItem('db-logs', 'true')
    localStorage.setItem('state-logs', 'true')
  }
}

// Función helper para deshabilitar todos los logs
export const disableLogs = () => {
  logger.setDebugMode(false)
  localStorage.removeItem('db-logs')
  localStorage.removeItem('state-logs')
}

// Función para mostrar estado actual de logs
export const getLogStatus = () => ({
  debug: logger.isDebugMode,
  db: localStorage.getItem('db-logs') === 'true',
  state: localStorage.getItem('state-logs') === 'true'
})

export default logger
