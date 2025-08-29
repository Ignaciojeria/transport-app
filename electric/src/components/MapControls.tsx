import { Crosshair } from 'lucide-react'

interface MapControlsProps {
  gpsActive: boolean
  onGPSToggle: () => void
  onZoomToSelected: () => void
  onNavigate: (provider: 'google' | 'waze' | 'geo') => void
}

export function MapControls({ 
  gpsActive, 
  onGPSToggle, 
  onZoomToSelected, 
  onNavigate 
}: MapControlsProps) {
  return (
    <div className="absolute top-3 right-3 space-y-2" style={{ zIndex: 1000 }}>
      {/* Botón de GPS del conductor */}
      <button
        onClick={onGPSToggle}
        className={`w-10 h-10 rounded-lg shadow-lg flex items-center justify-center transition-all ${
          gpsActive 
            ? 'bg-green-500 text-white hover:bg-green-600' 
            : 'bg-white text-gray-700 hover:bg-gray-50'
        } hover:shadow-xl`}
        aria-label={gpsActive ? 'Desactivar GPS' : 'Activar GPS'}
        title={gpsActive ? 'Desactivar GPS' : 'Activar GPS del conductor'}
      >
        <div className={`w-5 h-5 ${gpsActive ? 'animate-pulse' : ''}`}>
          <svg 
            viewBox="0 0 24 24" 
            fill="none" 
            stroke="currentColor" 
            strokeWidth="2" 
            strokeLinecap="round" 
            strokeLinejoin="round"
            className="w-full h-full"
          >
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
            <circle cx="12" cy="9" r="2.5"/>
          </svg>
        </div>
      </button>
      
      {/* Botón de zoom al punto seleccionado */}
      <button
        onClick={onZoomToSelected}
        className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-gray-700 hover:bg-gray-50 hover:shadow-xl transition-all"
        aria-label="Zoom al punto seleccionado"
        title="Hacer zoom al punto seleccionado"
      >
        <Crosshair className="w-5 h-5" />
      </button>
      
      {/* Botones de navegación */}
      <div className="flex flex-col gap-2">
        <button
          onClick={() => onNavigate('google')}
          className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-blue-600 hover:bg-gray-50 hover:shadow-xl transition-all"
          aria-label="Navegar con Google Maps"
          title="Google Maps"
        >
          G
        </button>
        <button
          onClick={() => onNavigate('waze')}
          className="w-10 h-10 bg-white rounded-lg shadow-lg flex items-center justify-center text-indigo-600 hover:bg-gray-50 hover:shadow-xl transition-all"
          aria-label="Navegar con Waze"
          title="Waze"
        >
          W
        </button>
      </div>
    </div>
  )
}
