import { useState, useEffect } from 'react'
import { useAccountByEmail } from './hooks/useAccountByEmail'
import { CreateOrganizationSimple } from './components/CreateOrganizationSimple'
import './App.css'

interface AuthData {
  access_token: string
  refresh_token: string
  user: {
    name: string
    email: string
  }
  expires_at: number
  stored_at: number
}

function App() {
  const [authData, setAuthData] = useState<AuthData | null>(null)
  const [isAuthLoading, setIsAuthLoading] = useState(true)

  // Cargar datos de autenticación del localStorage
  useEffect(() => {
    try {
      const storedAuth = localStorage.getItem('transport_auth')
      if (storedAuth) {
        const parsedAuth = JSON.parse(storedAuth) as AuthData
        
        // Verificar si el token no ha expirado
        if (parsedAuth.expires_at > Date.now()) {
          setAuthData(parsedAuth)
        } else {
          console.log('Token expirado, limpiando localStorage')
          localStorage.removeItem('transport_auth')
        }
      }
    } catch (error) {
      console.error('Error cargando auth data:', error)
      localStorage.removeItem('transport_auth')
    } finally {
      setIsAuthLoading(false)
    }
  }, [])

  // Verificar si estamos en modo demo
  const isDemoMode = localStorage.getItem('demo_mode') === 'true'

  // Hook para verificar si la cuenta existe en Electric (solo si no es demo)
  const { accountExists, isLoading: isAccountLoading, error } = useAccountByEmail(
    authData?.access_token || null,
    authData?.user?.email || null
  )

  // Estados de carga
  if (isAuthLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900 flex items-center justify-center text-white">
        <div className="text-center">
          <div className="text-5xl mb-4">🚛</div>
          <div className="text-lg">Cargando Transport APP...</div>
        </div>
      </div>
    )
  }

  // Si no hay datos de auth, mostrar mensaje para autenticarse
  if (!authData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900 flex items-center justify-center text-white">
        <div className="text-center max-w-md p-8">
          <div className="text-6xl mb-6">🚛</div>
          <h1 className="text-3xl font-bold mb-4">Transport APP</h1>
          <p className="text-lg opacity-90 mb-8">
            Para acceder a la aplicación, necesitas autenticarte primero.
          </p>
          <div className="space-y-4">
            <button 
              onClick={() => window.location.href = '/auth-ui'}
              className="w-full bg-white/20 hover:bg-white/30 border-2 border-white/30 text-white py-3 px-6 rounded-lg font-semibold backdrop-blur-sm transition-all duration-200 transform hover:scale-105"
            >
              Ir a Autenticación
            </button>
            
            <button 
              onClick={() => {
                // Simular que no hay cuenta para mostrar el formulario
                localStorage.removeItem('transport_auth');
                localStorage.setItem('demo_mode', 'true');
                window.location.reload();
              }}
              className="w-full bg-emerald-500/30 hover:bg-emerald-500/40 border-2 border-emerald-500/50 text-white py-2 px-5 rounded-lg font-medium backdrop-blur-sm transition-all text-sm"
            >
              🏢 Ver Demo: Crear Organización
            </button>
          </div>
        </div>
      </div>
    )
  }

  // Si está cargando la verificación de cuenta
  if (isAccountLoading) {
  return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        color: 'white',
        fontFamily: 'system-ui'
      }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>⚡</div>
          <div>Verificando cuenta de {authData.user.email}...</div>
        </div>
      </div>
    )
  }

  // Si hay error en la consulta
  if (error) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        color: 'white',
        fontFamily: 'system-ui'
      }}>
        <div style={{ textAlign: 'center', maxWidth: '400px', padding: '32px' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>❌</div>
          <h2>Error de conexión</h2>
          <p style={{ marginBottom: '24px', opacity: 0.9 }}>
            No se pudo verificar tu cuenta. Por favor, intenta de nuevo.
          </p>
          <button 
            onClick={() => window.location.reload()}
            style={{
              background: 'rgba(255, 255, 255, 0.2)',
              border: '2px solid rgba(255, 255, 255, 0.3)',
              color: 'white',
              padding: '12px 24px',
              borderRadius: '8px',
              cursor: 'pointer',
              fontSize: '16px',
              fontWeight: '600',
              backdropFilter: 'blur(10px)'
            }}
          >
            Reintentar
        </button>
        </div>
      </div>
    )
  }

  // Si estamos en modo demo O la cuenta NO existe, mostrar vista de crear organización
  if (isDemoMode || (!accountExists && authData)) {
    return <CreateOrganizationSimple userEmail={isDemoMode ? 'demo@example.com' : authData.user.email} />
  }

  // Si la cuenta SÍ existe, mostrar la aplicación principal
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900 flex items-center justify-center text-white">
      <div className="text-center max-w-md p-8">
        <div className="text-6xl mb-6">✅</div>
        <h1 className="text-3xl font-bold mb-4">¡Bienvenido a Transport APP!</h1>
        <p className="text-lg opacity-90 mb-6">
          Cuenta encontrada para {authData.user.email}
        </p>
        <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20">
          <p className="text-sm opacity-70 mb-4">
            🚛 Sistema de gestión logística
          </p>
          <p className="text-xs opacity-60">
            Aquí iría la aplicación principal de Transport APP
          </p>
        </div>
      </div>
    </div>
  )
}

export default App
