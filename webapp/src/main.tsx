import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

// 🔐 AUTH FRAGMENT DECODER - Recibir tokens del auth-ui
function processAuthFragment() {
  const fragment = window.location.hash
  
  if (fragment.startsWith('#auth=')) {
    try {
      console.log('🔍 Fragment de auth detectado:', fragment.substring(0, 20) + '...')
      
      // Decodificar el payload
      const encodedData = fragment.substring(6) // Remove '#auth='
      const authData = JSON.parse(atob(encodedData))
      
      console.log('✅ AUTH PAYLOAD DECODIFICADO:', {
        access_token: authData.access_token ? `${authData.access_token.substring(0, 30)}...` : 'NONE',
        refresh_token: authData.refresh_token ? `${authData.refresh_token.substring(0, 30)}...` : 'NONE',
        user: authData.user,
        expires_in: authData.expires_in,
        timestamp: new Date(authData.timestamp).toISOString()
      })
      
      // Guardar tokens en localStorage (temporal - debería ir a un store apropiado)
      localStorage.setItem('transport_auth', JSON.stringify({
        access_token: authData.access_token,
        refresh_token: authData.refresh_token,
        user: authData.user,
        expires_at: Date.now() + (authData.expires_in * 1000),
        stored_at: Date.now()
      }))
      
      // Limpiar URL inmediatamente
      window.history.replaceState({}, '', window.location.pathname)
      
      console.log('🧹 Fragment limpiado de la URL')
      console.log('💾 Tokens guardados en localStorage como "transport_auth"')
      
      // Mostrar mensaje temporal en pantalla (FEO pero funcional)
      const authBanner = document.createElement('div')
      authBanner.id = 'auth-success-banner'
      authBanner.innerHTML = `
        <div style="position: fixed; top: 0; left: 0; right: 0; background: #10B981; color: white; padding: 12px; text-align: center; z-index: 9999; font-family: system-ui;">
          ✅ Autenticado como ${authData.user.name} (${authData.user.email}) - Token expires: ${new Date(Date.now() + authData.expires_in * 1000).toLocaleTimeString()}
          <button onclick="document.getElementById('auth-success-banner').remove()" style="margin-left: 10px; background: rgba(255,255,255,0.2); border: none; color: white; padding: 4px 8px; border-radius: 4px; cursor: pointer;">×</button>
        </div>
      `
      document.body.appendChild(authBanner)
      
      // Auto-remove después de 10 segundos
      setTimeout(() => {
        const banner = document.getElementById('auth-success-banner')
        if (banner) banner.remove()
      }, 10000)
      
    } catch (error) {
      console.error('❌ Error decodificando fragment de auth:', error)
    }
  }
}

// Procesar fragment al cargar
processAuthFragment()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
