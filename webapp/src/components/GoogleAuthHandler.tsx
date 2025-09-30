import React, { useEffect } from 'react'
import { useGoogleAuthFlow, useAuthRedirect } from '../hooks/useGoogleAuthFlow'
import CreateOrganization from './CreateOrganization'
import TenantsList from './TenantsList'
import LoadingSpinner from './ui/LoadingSpinner'
import { clearElectricCache, getElectricCacheInfo } from '../utils/electricCacheUtils'

interface GoogleAuthHandlerProps {
  token: string
  email: string
  onError?: (error: string) => void
}

const GoogleAuthHandler: React.FC<GoogleAuthHandlerProps> = ({ 
  token, 
  email, 
  onError 
}) => {
  const authResult = useGoogleAuthFlow(token, email)
  const { isLoading } = useAuthRedirect(authResult)

  useEffect(() => {
    if (authResult.state === 'error' && onError) {
      onError(authResult.error || 'Error desconocido')
    }
  }, [authResult.state, authResult.error, onError])

  // Mostrar loading mientras se verifica la cuenta
  if (isLoading) {
    const cacheInfo = getElectricCacheInfo()
    
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-gray-600">
            {authResult.state === 'checking-account' && 'Verificando cuenta...'}
            {authResult.state === 'loading-tenants' && 'Cargando organizaciones...'}
          </p>
          
          {/* Información de debug en modo desarrollo */}
          {import.meta.env.DEV && (
            <div className="mt-6 bg-white/80 backdrop-blur-sm rounded-lg p-4 max-w-md mx-auto">
              <h3 className="font-semibold text-gray-800 mb-2">🐛 Debug Info</h3>
              <p className="text-sm text-gray-600 mb-1">
                <strong>Estado:</strong> {authResult.state}
              </p>
              <p className="text-sm text-gray-600 mb-1">
                <strong>Email:</strong> {email}
              </p>
              <p className="text-sm text-gray-600 mb-1">
                <strong>Caché Electric:</strong> {cacheInfo.keys.length} claves
              </p>
              <button 
                onClick={() => {
                  clearElectricCache()
                  console.log('🧹 Caché limpiado manualmente')
                }}
                className="mt-2 px-3 py-1 bg-yellow-500 text-white text-xs rounded hover:bg-yellow-600 transition-colors"
              >
                Limpiar Caché
              </button>
            </div>
          )}
        </div>
      </div>
    )
  }

  // Mostrar error si ocurrió alguno
  if (authResult.state === 'error') {
    const cacheInfo = getElectricCacheInfo()
    
    return (
      <div className="min-h-screen bg-gradient-to-br from-red-50 via-pink-50 to-orange-50 flex items-center justify-center">
        <div className="text-center max-w-2xl">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">Error de Autenticación</h2>
          <p className="text-gray-600 mb-4">{authResult.error}</p>
          
          {/* Información de debug del caché */}
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4 text-left">
            <h3 className="font-semibold text-yellow-800 mb-2">🐛 Debug - Información del Caché</h3>
            <p className="text-sm text-yellow-700 mb-2">
              <strong>Claves de Electric SQL:</strong> {cacheInfo.keys.length}
            </p>
            <p className="text-sm text-yellow-700 mb-2">
              <strong>Tamaño del caché:</strong> {cacheInfo.size} bytes
            </p>
            {cacheInfo.keys.length > 0 && (
              <div className="text-xs text-yellow-600">
                <strong>Claves encontradas:</strong>
                <ul className="list-disc list-inside mt-1">
                  {cacheInfo.keys.map(key => (
                    <li key={key}>{key}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>
          
          <div className="space-y-2">
            <button 
              onClick={() => window.location.reload()}
              className="px-6 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors mr-2"
            >
              Reintentar
            </button>
            <button 
              onClick={() => {
                clearElectricCache()
                window.location.reload()
              }}
              className="px-6 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition-colors"
            >
              Limpiar Caché y Reintentar
            </button>
          </div>
        </div>
      </div>
    )
  }

  // Si no existe la cuenta, mostrar formulario de creación de organización
  if (authResult.state === 'account-not-found') {
    return (
      <CreateOrganization 
        email={email}
        onSuccess={(response) => {
          console.log('Organización creada exitosamente:', response)
          // Aquí puedes redirigir o actualizar el estado
        }}
        onError={(error) => {
          console.error('Error al crear organización:', error)
          onError?.(error)
        }}
      />
    )
  }

  // Si la cuenta existe y se cargaron los tenants, mostrar la lista
  if (authResult.state === 'tenants-loaded') {
    return (
      <TenantsList 
        account={authResult.account!}
        tenants={authResult.tenants}
        token={token}
      />
    )
  }

  // Estado por defecto (no debería llegar aquí)
  return null
}

export default GoogleAuthHandler
