import React, { useEffect } from 'react'
import { useGoogleAuthFlow, useAuthRedirect } from '../hooks/useGoogleAuthFlow'
import CreateOrganization from './CreateOrganization'
import TenantsList from './TenantsList'
import LoadingSpinner from './ui/LoadingSpinner'

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
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-gray-600">
            {authResult.state === 'checking-account' && 'Verificando cuenta...'}
            {authResult.state === 'loading-tenants' && 'Cargando organizaciones...'}
          </p>
        </div>
      </div>
    )
  }

  // Mostrar error si ocurrió alguno
  if (authResult.state === 'error') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-red-50 via-pink-50 to-orange-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">Error de Autenticación</h2>
          <p className="text-gray-600 mb-4">{authResult.error}</p>
          <button 
            onClick={() => window.location.reload()}
            className="px-6 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          >
            Reintentar
          </button>
        </div>
      </div>
    )
  }

  // Si no existe la cuenta, mostrar formulario de creación de organización
  if (authResult.state === 'account-not-found') {
    return (
      <CreateOrganization 
        token={token}
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
