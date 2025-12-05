'use client'

import { useSearchParams } from 'next/navigation'
import Image from 'next/image'
import { motion } from 'framer-motion'
import { AlertTriangle, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/Button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card'

const errorMessages = {
  Configuration: 'Error de configuración del servidor.',
  AccessDenied: 'Acceso denegado. No tienes permisos para acceder a este sistema.',
  Verification: 'Token de verificación inválido o expirado.',
  Default: 'Ha ocurrido un error durante el proceso de autenticación.',
  DomainNotAllowed: 'Tu dominio de email no está autorizado para acceder a este sistema.',
  CredentialsSignin: 'Credenciales incorrectas. Verifica tu email y contraseña.',
}

export default function AuthError() {
  const searchParams = useSearchParams()
  const error = searchParams.get('error') as keyof typeof errorMessages

  const errorMessage = errorMessages[error] || errorMessages.Default

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
          className="w-full max-w-md"
        >
          {/* Logo */}
          <div className="flex items-center justify-center mb-8">
            <Image 
              src="/logo.png" 
              alt="MiCartaPro Logo" 
              width={200} 
              height={60}
              className="h-12 w-auto"
            />
          </div>

          <Card className="border-2 border-red-100 shadow-2xl">
            <CardHeader className="text-center space-y-4">
              <div className="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
                <AlertTriangle className="h-8 w-8 text-red-600" />
              </div>
              <CardTitle className="text-2xl font-bold text-gray-900">
                Error de Acceso
              </CardTitle>
            </CardHeader>

            <CardContent className="space-y-6 text-center">
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <p className="text-red-800 font-medium">
                  {errorMessage}
                </p>
              </div>

              {error === 'AccessDenied' && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                  <p className="text-blue-800 text-sm">
                    Si crees que deberías tener acceso, contacta al administrador del sistema
                    con tu dirección de email corporativo.
                  </p>
                </div>
              )}

              <div className="space-y-4">
                <Button
                  onClick={() => window.location.href = '/'}
                  className="w-full bg-blue-600 hover:bg-blue-700 text-white h-12"
                >
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Volver al inicio de sesión
                </Button>

                <div className="text-sm text-gray-600 space-y-2">
                  <p>¿Necesitas ayuda?</p>
                  <div className="flex justify-center space-x-4">
                    <a 
                      href="mailto:soporte@micartapro.com" 
                      className="text-blue-600 hover:text-blue-700"
                    >
                      Contactar Soporte
                    </a>
                    <span className="text-gray-400">•</span>
                    <a 
                      href="#" 
                      className="text-blue-600 hover:text-blue-700"
                    >
                      Documentación
                    </a>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Footer */}
          <div className="mt-8 text-center text-sm text-gray-500">
            <p>© 2024 MiCartaPro. Todos los derechos reservados.</p>
          </div>
        </motion.div>
      </div>
    </div>
  )
}

