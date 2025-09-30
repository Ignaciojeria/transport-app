import React, { useEffect, useState } from 'react'
import { useGoogleAuthFlow, useAuthRedirect } from '../hooks/useGoogleAuthFlow'
import CreateOrganization from './CreateOrganization'
import TenantsList from './TenantsList'
import LoadingSpinner from './ui/LoadingSpinner'
import SuccessNotification from './ui/SuccessNotification'
import { clearElectricCache, getElectricCacheInfo, forceAppReload } from '../utils/electricCacheUtils'
import { syncWithElectric } from '../utils/retryUtils'
import { isElectricSynced } from '../utils/electricSyncUtils'

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
  const [showSuccessNotification, setShowSuccessNotification] = useState(false)
  const [isRetrying, setIsRetrying] = useState(false)

  useEffect(() => {
    if (authResult.state === 'error' && onError) {
      onError(authResult.error || 'Error desconocido')
    }
  }, [authResult.state, authResult.error, onError])

  // Mostrar loading mientras se verifica la cuenta
  if (isLoading || isRetrying) {
    const cacheInfo = getElectricCacheInfo()
    
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-gray-600">
            {isRetrying && 'Sincronizando con Electric SQL...'}
            {!isRetrying && authResult.state === 'checking-account' && 'Verificando cuenta...'}
            {!isRetrying && authResult.state === 'loading-tenants' && 'Cargando organizaciones...'}
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
              <div className="space-y-2">
                <p className="text-xs text-gray-500 mb-2">
                  Con LiveQuery, la sincronización es automática. Usa estos botones solo si hay problemas:
                </p>
                <button 
                  onClick={async () => {
                    console.log('🔍 Verificando sincronización...')
                    const result = await isElectricSynced(email)
                    console.log('🔍 Estado de sincronización:', result)
                    alert(`Sincronización: ${result.synced ? 'SÍ' : 'NO'}\nMensaje: ${result.message}`)
                  }}
                  className="w-full px-3 py-1 bg-blue-500 text-white text-xs rounded hover:bg-blue-600 transition-colors"
                >
                  Verificar Estado
                </button>
                <button 
                  onClick={() => {
                    console.log('🧹 Limpiando caché local...')
                    clearElectricCache()
                    console.log('✅ Caché limpiado - LiveQuery se encargará de la sincronización')
                  }}
                  className="w-full px-3 py-1 bg-yellow-500 text-white text-xs rounded hover:bg-yellow-600 transition-colors"
                >
                  Limpiar Caché Local
                </button>
                <button 
                  onClick={() => {
                    console.log('🔄 Recarga completa como último recurso...')
                    forceAppReload()
                  }}
                  className="w-full px-3 py-1 bg-red-500 text-white text-xs rounded hover:bg-red-600 transition-colors"
                >
                  Recarga Completa
                </button>
              </div>
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
          console.log('✅ Organización creada exitosamente:', response)
          console.log('🔄 Recargando para verificar cuenta y cargar organizaciones...')
          
          // Mostrar notificación de éxito
          setShowSuccessNotification(true)
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

  // Mostrar notificación de éxito si está activa
  if (showSuccessNotification) {
    return (
      <SuccessNotification
        message="¡Organización creada exitosamente!"
        onComplete={async () => {
          console.log('🔄 Sincronizando con Electric SQL...')
          
          // Limpiar caché y reintentar verificación
          clearElectricCache()
          setIsRetrying(true)
          
          try {
            await syncWithElectric(async () => {
              console.log('🔄 Ejecutando sincronización...')
              await authResult.retry()
            })
            
            console.log('✅ Sincronización exitosa')
            setShowSuccessNotification(false)
            setIsRetrying(false)
          } catch (error) {
            console.error('❌ Error al sincronizar después de múltiples intentos:', error)
            setIsRetrying(false)
            
            // Mostrar error al usuario y ofrecer recargar
            if (onError) {
              onError('No se pudo sincronizar con la base de datos. Intenta recargar la página.')
            } else {
              // Fallback: recargar la página
              window.location.reload()
            }
          }
        }}
        duration={2000}
      />
    )
  }

  // Estado por defecto (no debería llegar aquí)
  return null
}

export default GoogleAuthHandler
