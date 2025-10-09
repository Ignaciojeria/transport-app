import { useState, useEffect } from 'react'
import GoogleAuthHandler from './components/GoogleAuthHandler'
import { extractTokenEarly } from './utils/earlyTokenExtraction'
import './App.css'

// Extraer token INMEDIATAMENTE cuando se carga el módulo
const { token: earlyToken, email: earlyEmail } = extractTokenEarly()

function App() {
  const [organizationCreated, setOrganizationCreated] = useState(false)
  const [organizationData, setOrganizationData] = useState<{name: string; country: string} | null>(null)
  const [token, setToken] = useState<string | null>(earlyToken)
  const [email, setEmail] = useState<string | null>(earlyEmail)
  const [isLoading, setIsLoading] = useState(!earlyToken)

  // Si ya tenemos el token de la extracción temprana, no necesitamos hacer nada más
  useEffect(() => {
    if (earlyToken) {
      console.log('✅ Token ya extraído tempranamente:', earlyToken.substring(0, 20) + '...')
      console.log('✅ Email extraído tempranamente:', earlyEmail)
      setEmail(earlyEmail)
      setIsLoading(false)
      return
    }

    // Si no hay token temprano, intentar extraer del fragment o localStorage
    console.log('🚀 Iniciando extracción de token...')
    console.log('🚀 URL actual:', window.location.href)
    
    // Verificar localStorage como fallback
    console.log('🔍 Verificando localStorage para tokens guardados...')
    const storedAuth = localStorage.getItem('transport_auth')
    if (storedAuth) {
      try {
        const authData = JSON.parse(storedAuth)
        console.log('🔍 Auth data encontrada en localStorage:', authData)
        
        if (authData.access_token) {
          console.log('✅ Access token encontrado en localStorage')
          setToken(authData.access_token)
          // Extraer email del localStorage también
          if (authData.user?.email) {
            setEmail(authData.user.email)
          }
          setIsLoading(false)
          return
        }
      } catch (error) {
        console.error('❌ Error al parsear auth data del localStorage:', error)
      }
    }
    
    console.warn('❌ No se encontró token en la extracción temprana ni en localStorage')
    setToken(null)
    setEmail(null)
    setIsLoading(false)
  }, [])


  if (organizationCreated && organizationData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex flex-col items-center justify-center p-8">
        <div className="bg-white rounded-xl shadow-xl p-8 max-w-md w-full text-center">
          <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"/>
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">¡Organización Creada!</h2>
          <p className="text-gray-600 mb-4">
            La organización <strong>"{organizationData.name}"</strong> ha sido creada exitosamente en <strong>{organizationData.country}</strong>.
          </p>
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mb-4">
            <p className="text-yellow-800 text-sm font-medium">
              📋 Waiting List
            </p>
          </div>
          <button 
            onClick={() => {
              setOrganizationCreated(false)
              setOrganizationData(null)
            }}
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition-colors duration-200"
          >
            Crear Otra Organización
          </button>
        </div>
      </div>
    )
  }

  // Mostrar loading mientras se extrae el token
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Procesando autenticación...</p>
        </div>
      </div>
    )
  }

  // Si no hay token, mostrar error
  if (!token) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-red-50 via-pink-50 to-orange-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">🔐</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">Autenticación Requerida</h2>
          <p className="text-gray-600 mb-4">
            No se encontró token de autenticación en la URL.<br/>
            Por favor, inicia sesión con Google primero.
          </p>
          <div className="space-y-2">
            <button 
              onClick={() => window.location.reload()}
              className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors mr-2"
            >
              Reintentar
            </button>
            <button 
              onClick={() => {
                // Aquí podrías redirigir a tu sistema de autenticación
                console.log('Redirigir a autenticación...')
              }}
              className="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
            >
              Iniciar Sesión
            </button>
          </div>
        </div>
      </div>
    )
  }

  // Si no hay token o email, mostrar error de autenticación
  if (!token || !email) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-red-50 via-pink-50 to-orange-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">🔐</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">Autenticación Requerida</h2>
          <p className="text-gray-600 mb-4">
            No se encontró token de autenticación en la URL.<br/>
            Por favor, inicia sesión con Google primero.
          </p>
          <div className="space-y-2">
            <button 
              onClick={() => window.location.reload()}
              className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors mr-2"
            >
              Reintentar
            </button>
            <button 
              onClick={() => {
                // Aquí podrías redirigir a tu sistema de autenticación
                console.log('Redirigir a autenticación...')
              }}
              className="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
            >
              Iniciar Sesión
            </button>
          </div>
        </div>
      </div>
    )
  }

  // Usar GoogleAuthHandler para manejar el flujo completo
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

export default App
