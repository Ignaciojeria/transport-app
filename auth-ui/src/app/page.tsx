'use client'

import { useState } from 'react'
import { signIn, useSession } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { motion } from 'framer-motion'
import { 
  Truck, 
  ArrowRight,
  Activity,
  Shield,
  Zap,
  MapPin,
  Navigation
} from 'lucide-react'
import { Button } from '@/components/ui/Button'

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  
  const { data: session, status } = useSession()
  const router = useRouter()

  // Redirigir si ya está autenticado
  if (status === 'loading') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Cargando...</p>
        </div>
      </div>
    )
  }

  if (session) {
    // Redirigir a donde quieras después del login
    window.location.href = '/success' // Cambiar por tu URL de destino
    return null
  }

  const handleGoogleLogin = async () => {
    setIsLoading(true)
    setError('')
    try {
      const result = await signIn('google', { 
        redirect: false,
        callbackUrl: '/success' // Cambiar por tu URL de destino
      })
      
      if (result?.error) {
        setError('Error al iniciar sesión con Google')
      } else if (result?.url) {
        window.location.href = result.url
      }
    } catch (err) {
      setError('Error de conexión')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 particles-container flex flex-col">
      {/* Partículas de fondo */}
      <div className="auth-particles">
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
      </div>

      <div className="flex flex-col lg:flex-row flex-1 relative z-10 lg:h-screen">
        {/* Panel izquierdo - Información de TransportApp */}
        <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-blue-600 via-blue-700 to-indigo-800 text-white p-6 xl:p-8 flex-col justify-center overflow-y-auto">
          <motion.div
            initial={{ opacity: 0, x: -50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            className="max-w-lg mx-auto"
          >
            {/* Logo y título */}
            <div className="flex items-center space-x-4 mb-8">
              <div className="bg-white/20 p-4 rounded-2xl backdrop-blur-sm logo-pulse">
                <Truck className="h-12 w-12 text-white" />
              </div>
              <div>
                <h1 className="text-4xl font-bold leading-tight">TransportApp</h1>
                <p className="text-blue-200 text-lg leading-relaxed">Sistema de gestión logística</p>
              </div>
            </div>

            {/* Mensaje principal */}
            <div className="space-y-6">
              <h2 className="text-3xl font-bold leading-snug">
                Todo lo necesario para administrar tus procesos logísticos
              </h2>
              
              <p className="text-blue-100 text-xl leading-loose">
                Optimiza rutas, gestiona tu flota y maximiza la eficiencia 
                de tus entregas desde una sola plataforma.
              </p>

              {/* Módulos del panel */}
              <div className="space-y-4 mt-8">
                <div className="flex items-center space-x-4">
                  <div className="bg-green-500/20 p-3 rounded-xl backdrop-blur-sm">
                    <Activity className="h-6 w-6 text-green-300" />
                  </div>
                  <span className="text-blue-100 text-lg leading-relaxed">Panel de control y planificación central</span>
                </div>
                <div className="flex items-center space-x-4">
                  <div className="bg-purple-500/20 p-3 rounded-xl backdrop-blur-sm">
                    <Truck className="h-6 w-6 text-purple-300" />
                  </div>
                  <span className="text-blue-100 text-lg leading-relaxed">Administración de flotas y conductores</span>
                </div>
                <div className="flex items-center space-x-4">
                  <div className="bg-orange-500/20 p-3 rounded-xl backdrop-blur-sm">
                    <MapPin className="h-6 w-6 text-orange-300" />
                  </div>
                  <span className="text-blue-100 text-lg leading-relaxed">Monitoreo en vivo</span>
                </div>
              </div>
            </div>

            {/* Visualización de métricas */}
            <div className="mt-8 relative">
              <div className="absolute inset-0 bg-gradient-to-r from-blue-400/20 to-purple-400/20 rounded-3xl blur-xl"></div>
              <div className="relative bg-white/10 backdrop-blur-sm rounded-2xl p-6 text-center">
                <div className="text-5xl font-bold text-white mb-2 leading-none">24/7</div>
                <div className="text-blue-200 text-base leading-tight">Monitoreo continuo</div>
                <div className="text-blue-300 text-sm mt-1 leading-relaxed">Siempre optimizando</div>
              </div>
            </div>
          </motion.div>
        </div>

        {/* Panel derecho - Login simplificado */}
        <div className="w-full lg:w-1/2 flex items-center justify-center p-4 md:p-6 min-h-screen lg:min-h-0 lg:overflow-y-auto">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="w-full max-w-md"
          >
            {/* Logo móvil */}
            <div className="flex items-center justify-center mb-8 lg:hidden">
              <div className="flex items-center space-x-4">
                <div className="bg-blue-600 p-4 rounded-2xl">
                  <Truck className="h-10 w-10 text-white" />
                </div>
                <div>
                  <h1 className="text-3xl font-bold text-gray-900 leading-tight">TransportApp</h1>
                  <p className="text-gray-600 leading-relaxed">Acceso al sistema</p>
                </div>
              </div>
            </div>

            {/* Título principal */}
            <div className="text-center mb-8">
              <h2 className="text-3xl font-bold text-gray-900 mb-4 leading-tight">
                Iniciar Sesión
              </h2>
              <p className="text-gray-600 text-lg leading-relaxed">
                ¿Nuevo en TransportApp? <span className="text-blue-600 font-medium cursor-pointer hover:text-blue-700">Regístrate gratis</span>
              </p>
            </div>

            {/* Mensaje de error */}
            {error && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                className="bg-red-50 border border-red-200 rounded-xl p-4 mb-6"
              >
                <p className="text-sm text-red-800 text-center">{error}</p>
              </motion.div>
            )}

            {/* Botón de Google - Estilo minimalista */}
            <motion.div
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Button
                type="button"
                onClick={handleGoogleLogin}
                disabled={isLoading}
                className="w-full h-14 bg-gray-900 hover:bg-gray-800 text-white border-0 rounded-xl font-medium text-lg transition-all duration-200 shadow-lg hover:shadow-xl leading-relaxed"
              >
                <div className="flex items-center justify-center space-x-4">
                  <svg className="w-6 h-6" viewBox="0 0 24 24">
                    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                  </svg>
                  <span className="text-lg leading-relaxed">
                    {isLoading ? 'Conectando...' : 'Iniciar sesión con Google'}
                  </span>
                  {!isLoading && <ArrowRight className="h-5 w-5" />}
                  {isLoading && (
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white" />
                  )}
                </div>
              </Button>
            </motion.div>

            {/* Información adicional */}
            <div className="mt-6 space-y-4">
              <div className="bg-blue-50 border border-blue-200 rounded-xl p-4">
                <div className="flex items-start space-x-3">
                  <Shield className="h-5 w-5 text-blue-600 mt-0.5" />
                  <div>
                    <p className="text-sm font-medium text-blue-900 leading-tight">Acceso Seguro</p>
                    <p className="text-xs text-blue-700 mt-1 leading-relaxed">
                      Utiliza tu cuenta de Google corporativa para acceder de forma segura
                    </p>
                  </div>
                </div>
              </div>

              <p className="text-center text-sm text-gray-500 leading-relaxed">
                ¿Tienes problemas con tu cuenta?{' '}
                <span className="text-blue-600 font-medium cursor-pointer hover:text-blue-700">
                  Contactar soporte
                </span>
              </p>
            </div>

            {/* Footer */}
            <div className="mt-8 text-center text-sm text-gray-400">
              <p className="leading-relaxed">© 2024 TransportApp. Todos los derechos reservados.</p>
              <div className="flex justify-center space-x-6 mt-3">
                <span className="hover:text-blue-600 cursor-pointer leading-relaxed">Términos de Uso</span>
                <span className="hover:text-blue-600 cursor-pointer leading-relaxed">Política de Privacidad</span>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  )
}