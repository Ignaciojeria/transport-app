import { useLanguage } from '../hooks/useLanguage'

interface SidebarProps {
  isOpen: boolean
  onClose: () => void
  routeStarted: boolean
  onDownloadReport: () => void
  syncInfo?: {
    deviceId: string
  } | null
  markerPosition?: {
    visitIndex: number
    coordinates: [number, number]
    timestamp: number
    deviceId: string
  } | null
}

export function Sidebar({ 
  isOpen, 
  onClose, 
  routeStarted, 
  onDownloadReport,
  syncInfo,
  markerPosition
}: SidebarProps) {
  const { t } = useLanguage()
  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50">
      {/* Overlay */}
      <div 
        className="absolute inset-0 bg-black/40 backdrop-blur-sm transition-opacity duration-300"
        onClick={onClose}
      ></div>
      
      {/* Sidebar Panel */}
      <div className="absolute top-0 left-0 h-full w-80 bg-white shadow-2xl transform transition-transform duration-300 ease-out">
        {/* Header */}
        <div className="bg-gradient-to-r from-indigo-600 to-purple-600 text-white p-6">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold">{t.sidebar.title}</h2>
            <button 
              onClick={onClose}
              className="bg-white/20 hover:bg-white/30 rounded-lg p-2 transition-colors duration-200"
              aria-label={t.sidebar.closeMenu}
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        
        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Botón CSV */}
          {routeStarted && (
            <div className="space-y-2">
              <h3 className="text-sm font-semibold text-gray-600 uppercase tracking-wide">{t.sidebar.reports}</h3>
              <button
                onClick={() => {
                  onDownloadReport()
                  onClose()
                }}
                className="w-full bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white p-4 rounded-xl font-medium transition-all duration-200 shadow-lg hover:shadow-xl active:scale-95 flex items-center justify-center space-x-3"
                aria-label={t.sidebar.downloadReport}
              >
                <span className="text-2xl">📊</span>
                <span>{t.sidebar.downloadReport}</span>
              </button>
            </div>
          )}
          
          {/* Indicadores de conexión */}
          <div className="space-y-2">
            <h3 className="text-sm font-semibold text-gray-600 uppercase tracking-wide">{t.sidebar.connectionStatus}</h3>
            <div className="bg-gray-50 rounded-xl p-4 space-y-3">
              {/* Estado de conexión a internet */}
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-700">{t.sidebar.internet}</span>
                <div className="flex items-center space-x-2">
                  {navigator.onLine ? (
                    <>
                      <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                      <span className="text-sm text-green-600 font-medium">{t.sidebar.connected}</span>
                    </>
                  ) : (
                    <>
                      <div className="w-2 h-2 bg-red-500 rounded-full"></div>
                      <span className="text-sm text-red-600 font-medium">{t.sidebar.disconnected}</span>
                    </>
                  )}
                </div>
              </div>
              
              {/* Estado de sincronización GunJS */}
              {syncInfo && (
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-gray-700">{t.sidebar.synchronization}</span>
                  <div className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
                    <span className="text-sm text-blue-600 font-medium">
                      {syncInfo.deviceId.slice(-6)}
                    </span>
                  </div>
                </div>
              )}
              
              {/* Indicador de posición sincronizada */}
              {markerPosition && (Date.now() - markerPosition.timestamp) < 30000 && (
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-gray-700">{t.sidebar.markerSync}</span>
                  <div className="flex items-center space-x-2">
                    <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse"></div>
                    <span className="text-sm text-purple-600 font-medium">
                      📍 {markerPosition.deviceId.slice(-6)}
                    </span>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
