import { useState, useEffect } from 'react'
import CreateOrganization from './components/CreateOrganization'
import { type CreateOrganizationResponse } from './services/organizationService'
import { extractTokenFromFragment } from './utils/urlUtils'
import './App.css'

function App() {
  const [organizationCreated, setOrganizationCreated] = useState(false)
  const [organizationData, setOrganizationData] = useState<{name: string; country: string} | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Extraer token del fragment de la URL
  useEffect(() => {
    console.log('🚀 Iniciando extracción de token...')
    console.log('🚀 URL actual:', window.location.href)
    
    const extractedToken = extractTokenFromFragment()
    console.log('🚀 Token extraído:', extractedToken ? 'SÍ' : 'NO')
    
    if (extractedToken) {
      console.log('✅ Token encontrado, estableciendo...')
      setToken(extractedToken)
      // Limpiar la URL después de extraer el token
      const cleanUrl = window.location.origin + window.location.pathname
      console.log('🧹 Limpiando URL a:', cleanUrl)
      window.history.replaceState({}, document.title, cleanUrl)
    } else {
      console.warn('❌ No se encontró token en el fragment de la URL')
      // No establecer token - esto hará que se muestre el error de autenticación
      setToken(null)
    }
    
    console.log('🏁 Finalizando extracción de token')
    setIsLoading(false)
  }, [])

  const handleCreateOrganizationSuccess = (response: CreateOrganizationResponse) => {
    console.log('Organización creada exitosamente:', response)
    setOrganizationData({
      name: 'Organización Creada', // En una implementación real, obtendrías esto de la respuesta
      country: 'CL'
    })
    setOrganizationCreated(true)
    // setError(null)
  }

  const handleCreateOrganizationError = (error: string) => {
    console.error('Error al crear organización:', error)
    // setError(error)
  }

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

  return (
    <CreateOrganization 
      token={token}
      onSuccess={handleCreateOrganizationSuccess}
      onError={handleCreateOrganizationError}
    />
  )
}

export default App
