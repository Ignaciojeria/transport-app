

interface DownloadReportModalProps {
  isOpen: boolean
  onClose: () => void
  onDownloadReport: (format: 'csv' | 'excel') => void
}

export function DownloadReportModal({ isOpen, onClose, onDownloadReport }: DownloadReportModalProps) {
  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onClose}></div>
      <div className="relative bg-white w-full max-w-md mx-auto rounded-xl shadow-xl border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">Descargar Reporte de Ruta</h3>
        
        <p className="text-sm text-gray-600 mb-6">
          Selecciona el formato en el que deseas descargar el reporte con toda la información de la ruta:
        </p>

        <div className="space-y-3 mb-6">
          {/* Opción CSV */}
          <button
            onClick={() => onDownloadReport('csv')}
            className="w-full p-4 border-2 border-gray-200 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-all duration-200 text-left group"
          >
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center group-hover:bg-green-200 transition-colors">
                <span className="text-green-600 font-bold text-lg">CSV</span>
              </div>
              <div className="flex-1">
                <h4 className="font-medium text-gray-800">Archivo CSV</h4>
                <p className="text-sm text-gray-500">Compatible con Excel, Google Sheets, etc.</p>
                <p className="text-xs text-gray-400 mt-1">Formato estándar para datos tabulares</p>
              </div>
            </div>
          </button>

          {/* Opción Excel */}
          <button
            onClick={() => onDownloadReport('excel')}
            className="w-full p-4 border-2 border-gray-200 rounded-lg hover:border-blue-300 hover:bg-blue-50 transition-all duration-200 text-left group"
          >
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center group-hover:bg-blue-200 transition-colors">
                <span className="text-blue-600 font-bold text-lg">XLS</span>
              </div>
              <div className="flex-1">
                <h4 className="font-medium text-gray-800">Archivo Excel</h4>
                <p className="text-sm text-gray-500">Abre directamente en Microsoft Excel</p>
                <p className="text-xs text-gray-400 mt-1">Con formato y estilos incluidos</p>
              </div>
            </div>
          </button>
        </div>

        <div className="flex items-center justify-end gap-3">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm rounded-lg border border-gray-300 bg-white hover:bg-gray-50 text-gray-700 transition-colors"
          >
            Cancelar
          </button>
        </div>
      </div>
    </div>
  )
}
