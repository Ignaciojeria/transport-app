import React, { useState } from 'react'
import GoogleAuthHandler from './GoogleAuthHandler'

const AuthExample: React.FC = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [token, setToken] = useState('')
  const [email, setEmail] = useState('')

  const handleGoogleLogin = async () => {
    try {
      // Aquí implementarías la lógica real de Google OAuth
      // Por ahora simulamos un login exitoso
      const mockToken = 'mock-jwt-token-12345'
      const mockEmail = 'usuario@ejemplo.com'
      
      setToken(mockToken)
      setEmail(mockEmail)
      setIsAuthenticated(true)
    } catch (error) {
      console.error('Error en login de Google:', error)
    }
  }

  // const handleLogout = () => {
  //   setToken('')
  //   setEmail('')
  //   setIsAuthenticated(false)
  // }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-800 mb-4">
            Transport App
          </h1>
          <p className="text-gray-600 mb-8">
            Inicia sesión con Google para continuar
          </p>
          <button
            onClick={handleGoogleLogin}
            className="px-8 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors shadow-lg"
          >
            Iniciar Sesión con Google
          </button>
        </div>
      </div>
    )
  }

  return (
    <GoogleAuthHandler 
      token={token}
      email={email}
      onError={(error) => {
        console.error('Error en el flujo de autenticación:', error)
        // Aquí puedes mostrar una notificación de error al usuario
      }}
    />
  )
}

export default AuthExample
