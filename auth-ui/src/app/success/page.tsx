'use client'

import { useSession, signOut } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { motion } from 'framer-motion'
import { 
  Truck, 
  CheckCircle, 
  LogOut, 
  User,
  Clock
} from 'lucide-react'
import { Button } from '@/components/ui/Button'

export default function SuccessPage() {
  const { data: session, status } = useSession()
  const router = useRouter()

  useEffect(() => {
    if (status === 'loading') return
    if (!session) {
      router.push('/')
    }
  }, [session, status, router])

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

  if (!session) {
    return null
  }

  const handleSignOut = async () => {
    await signOut({ callbackUrl: '/' })
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 particles-container">
      {/* Partículas de fondo */}
      <div className="auth-particles">
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
        <div className="particle"></div>
      </div>

      <div className="flex min-h-screen items-center justify-center p-8 relative z-10">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="w-full max-w-2xl text-center"
        >
          {/* Logo */}
          <div className="flex items-center justify-center mb-8">
            <div className="flex items-center space-x-4">
              <div className="bg-green-600 p-4 rounded-2xl">
                <Truck className="h-12 w-12 text-white" />
              </div>
              <div>
                <h1 className="text-3xl font-bold text-gray-900">TransportApp</h1>
                <p className="text-gray-600">Autenticación exitosa</p>
              </div>
            </div>
          </div>

          {/* Ícono de éxito */}
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="mx-auto w-24 h-24 bg-green-100 rounded-full flex items-center justify-center mb-8"
          >
            <CheckCircle className="h-12 w-12 text-green-600" />
          </motion.div>

          {/* Mensaje de bienvenida */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="space-y-6"
          >
            <h2 className="text-4xl font-bold text-gray-900">
              ¡Bienvenido, {session.user?.name?.split(' ')[0]}!
            </h2>
            
            <p className="text-xl text-gray-600 max-w-lg mx-auto">
              Has iniciado sesión exitosamente en TransportApp. 
              Tu cuenta está lista para usar.
            </p>
          </motion.div>

          {/* Información de la sesión */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.4 }}
            className="mt-12 bg-white rounded-2xl shadow-lg p-8 max-w-md mx-auto"
          >
            <h3 className="text-lg font-semibold text-gray-900 mb-6">
              Detalles de la sesión
            </h3>
            
            <div className="space-y-4 text-left">
              <div className="flex items-center space-x-3">
                <User className="h-5 w-5 text-blue-600" />
                <div>
                  <p className="text-sm font-medium text-gray-900">Usuario</p>
                  <p className="text-sm text-gray-600">{session.user?.email}</p>
                </div>
              </div>
              
              <div className="flex items-center space-x-3">
                <Clock className="h-5 w-5 text-blue-600" />
                <div>
                  <p className="text-sm font-medium text-gray-900">Inicio de sesión</p>
                  <p className="text-sm text-gray-600">
                    {new Date().toLocaleString('es-ES')}
                  </p>
                </div>
              </div>

              {session.user?.image && (
                <div className="flex items-center space-x-3">
                  <img
                    src={session.user.image}
                    alt="Profile"
                    className="h-10 w-10 rounded-full"
                  />
                  <div>
                    <p className="text-sm font-medium text-gray-900">Foto de perfil</p>
                    <p className="text-sm text-gray-600">Sincronizada con Google</p>
                  </div>
                </div>
              )}
            </div>
          </motion.div>

          {/* Mensaje informativo */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.5 }}
            className="mt-8 bg-blue-50 border border-blue-200 rounded-xl p-6"
          >
            <div className="flex items-start space-x-3">
              <CheckCircle className="h-6 w-6 text-blue-600 mt-0.5" />
              <div className="text-left">
                <h4 className="font-medium text-blue-900 mb-2">
                  Autenticación completada
                </h4>
                <p className="text-sm text-blue-700">
                  Tu sesión está activa y segura. Desde aquí puedes integrar con el 
                  sistema principal de TransportApp o redirigir a tu dashboard personalizado.
                </p>
              </div>
            </div>
          </motion.div>

          {/* Acciones */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.6 }}
            className="mt-8 space-y-4"
          >
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button
                onClick={() => alert('Aquí integrarías con tu dashboard principal')}
                className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3"
              >
                Ir al Dashboard
              </Button>
              
              <Button
                variant="outline"
                onClick={handleSignOut}
                className="border-2 border-gray-300 hover:border-red-300 hover:text-red-600 px-8 py-3"
              >
                <LogOut className="h-4 w-4 mr-2" />
                Cerrar sesión
              </Button>
            </div>

            <p className="text-sm text-gray-500">
              Esta es una página de ejemplo. Personaliza la redirección según tus necesidades.
            </p>
          </motion.div>

          {/* Footer */}
          <div className="mt-16 text-center text-sm text-gray-400">
            <p>© 2024 TransportApp. Sistema de autenticación.</p>
          </div>
        </motion.div>
      </div>
    </div>
  )
}
