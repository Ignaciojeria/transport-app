import { useRef, useEffect } from 'react'

interface RouteStartModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: (license: string) => void
  defaultLicense?: string
}

export function RouteStartModal({ isOpen, onClose, onConfirm, defaultLicense }: RouteStartModalProps) {
  const licenseInputRef = useRef<HTMLInputElement | null>(null)

  // Focus en el input cuando se abre el modal
  useEffect(() => {
    if (isOpen && licenseInputRef.current) {
      licenseInputRef.current.focus()
    }
  }, [isOpen])

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onClose}></div>
      <div className="relative bg-white w-full max-w-md mx-auto rounded-xl shadow-xl border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">Ingresar patente del vehículo</h3>
        
        {defaultLicense && (
          <div className="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
            <p className="text-sm text-blue-800 mb-2">
              ¿Usar la patente asignada a esta ruta?
            </p>
            <div className="flex items-center gap-2">
              <span className="font-mono text-base font-semibold text-blue-900">
                {defaultLicense}
              </span>
              <button
                onClick={() => onConfirm(defaultLicense)}
                className="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
              >
                Usar esta
              </button>
              <button
                onClick={() => {
                  // Limpiar el input para que ingrese otra patente
                  if (licenseInputRef.current) {
                    licenseInputRef.current.value = ''
                    licenseInputRef.current.focus()
                  }
                }}
                className="px-2 py-1 text-xs bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
              >
                Usar otra
              </button>
            </div>
          </div>
        )}
        
        <div className="mb-4">
          <label htmlFor="licenseInput" className="block text-sm font-medium text-gray-700 mb-2">
            Patente del vehículo:
          </label>
          <input
            id="licenseInput"
            ref={licenseInputRef}
            type="text"
            onChange={(e) => {
              e.target.value = e.target.value.toUpperCase()
            }}
            onKeyPress={(e) => {
              if (e.key === 'Enter') {
                const value = e.currentTarget.value.trim()
                if (value) {
                  onConfirm(value)
                }
              }
            }}
            placeholder="Ej: ABC123"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent uppercase"
            maxLength={8}
          />
          <p className="text-xs text-gray-500 mt-1">
            Puedes ingresar cualquier patente para esta ruta
          </p>
        </div>

        <div className="flex items-center justify-end gap-3">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm rounded-lg border border-gray-300 bg-white hover:bg-gray-50 text-gray-700 transition-colors"
          >
            Cancelar
          </button>
          <button
            onClick={() => {
              const value = licenseInputRef.current?.value.trim()
              if (value) {
                onConfirm(value)
              }
            }}
            disabled={!licenseInputRef.current?.value.trim()}
            className={`px-4 py-2 text-sm rounded-lg text-white font-medium transition-colors ${
              !licenseInputRef.current?.value.trim()
                ? 'bg-green-300 cursor-not-allowed'
                : 'bg-green-600 hover:bg-green-700'
            }`}
          >
            Iniciar ruta
          </button>
        </div>
      </div>
    </div>
  )
}
