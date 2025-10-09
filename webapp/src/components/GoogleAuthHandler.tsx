import React from 'react'
import { useAccountData } from '../hooks/useAccountData'
import CreateOrganization from './CreateOrganization'
import OrganizationDashboard from './OrganizationDashboard'
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
  const { account, tenants, isLoading, error } = useAccountData(token, email)

  // Mostrar loading mientras se cargan los datos
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-gray-600">Cargando datos...</p>
        </div>
      </div>
    )
  }

  // Mostrar error si ocurrió alguno
  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-red-50 via-pink-50 to-orange-50 flex items-center justify-center">
        <div className="text-center max-w-2xl">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">Error de Carga</h2>
          <p className="text-gray-600 mb-4">{error || 'Error desconocido'}</p>
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
  if (!account) {
    return (
      <CreateOrganization 
        email={email}
        onSuccess={(response) => {
          console.log('✅ Organización creada exitosamente:', response)
          console.log('🔄 Recargando para verificar cuenta y cargar organizaciones...')
          window.location.reload()
        }}
        onError={(error) => {
          console.error('Error al crear organización:', error)
          onError?.(error)
        }}
      />
    )
  }

  // Si la cuenta existe, mostrar el dashboard de organizaciones
  return (
    <OrganizationDashboard 
      account={account}
      tenants={tenants}
      token={token}
    />
  )
}

export default GoogleAuthHandler