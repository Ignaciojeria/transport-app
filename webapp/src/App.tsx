import { useState, useEffect } from 'react'
import { useAccountByEmail } from './hooks/useAccountByEmail'
import { CreateOrganization } from './components/CreateOrganization'
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

  // Hook para verificar si la cuenta existe en Electric
  const { accountExists, isLoading: isAccountLoading, error } = useAccountByEmail(
    authData?.access_token || null,
    authData?.user?.email || null
  )

  // Estados de carga
  if (isAuthLoading) {
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
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>🚛</div>
          <div>Cargando Transport APP...</div>
        </div>
      </div>
    )
  }

  // Si no hay datos de auth, mostrar mensaje para autenticarse
  if (!authData) {
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
          <div style={{ fontSize: '64px', marginBottom: '24px' }}>🚛</div>
          <h1 style={{ marginBottom: '16px' }}>Transport APP</h1>
          <p style={{ marginBottom: '24px', opacity: 0.9 }}>
            Para acceder a la aplicación, necesitas autenticarte primero.
          </p>
          <button 
            onClick={() => window.location.href = '/auth-ui'}
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
            Ir a Autenticación
          </button>
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

  // Si la cuenta NO existe, mostrar vista de crear organización
  if (!accountExists) {
    return <CreateOrganization userEmail={authData.user.email} />
  }

  // Si la cuenta SÍ existe, mostrar la aplicación principal
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
        <div style={{ fontSize: '64px', marginBottom: '24px' }}>✅</div>
        <h1 style={{ marginBottom: '16px' }}>¡Bienvenido!</h1>
        <p style={{ marginBottom: '24px', opacity: 0.9 }}>
          Cuenta encontrada para {authData.user.email}
        </p>
        <p style={{ fontSize: '14px', opacity: 0.7 }}>
          Aquí iría la aplicación principal de Transport APP
        </p>
      </div>
    </div>
  )
}

export default App
