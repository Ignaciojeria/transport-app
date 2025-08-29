import React, { useState, useEffect } from 'react'
import { 
  setDeliveryStatus, 
  getDeliveryStatus, 
  useDriverState,
  syncLocalBackupToGun,
  clearLocalBackup
} from '../db/driver-gun-state'

export function PersistenceTest() {
  const [testRouteId] = useState('test-route-123')
  const [testVisitIndex] = useState(0)
  const [testOrderIndex] = useState(0)
  const [testUnitIndex] = useState(0)
  const [currentStatus, setCurrentStatus] = useState<string>('')
  const [message, setMessage] = useState('')
  
  const { data: localState } = useDriverState()

  // Cargar estado actual al montar el componente
  useEffect(() => {
    const status = getDeliveryStatus(testRouteId, testVisitIndex, testOrderIndex, testUnitIndex)
    setCurrentStatus(status || 'no-status')
  }, [localState])

  const handleSetDelivered = () => {
    try {
      setDeliveryStatus(testRouteId, testVisitIndex, testOrderIndex, testUnitIndex, 'delivered')
      setMessage('✅ Estado marcado como entregado y guardado localmente')
      setTimeout(() => setMessage(''), 3000)
    } catch (error) {
      setMessage(`❌ Error: ${error}`)
      setTimeout(() => setMessage(''), 3000)
    }
  }

  const handleSetNotDelivered = () => {
    try {
      setDeliveryStatus(testRouteId, testVisitIndex, testOrderIndex, testUnitIndex, 'not-delivered')
      setMessage('✅ Estado marcado como no entregado y guardado localmente')
      setTimeout(() => setMessage(''), 3000)
    } catch (error) {
      setMessage(`❌ Error: ${error}`)
      setTimeout(() => setMessage(''), 3000)
    }
  }

  const handleSyncLocalBackup = () => {
    try {
      const count = syncLocalBackupToGun()
      setMessage(`🔄 Sincronizados ${count} elementos del respaldo local`)
      setTimeout(() => setMessage(''), 3000)
    } catch (error) {
      setMessage(`❌ Error sincronizando: ${error}`)
      setTimeout(() => setMessage(''), 3000)
    }
  }

  const handleClearBackup = () => {
    try {
      clearLocalBackup()
      setMessage('🗑️ Respaldo local limpiado')
      setTimeout(() => setMessage(''), 3000)
    } catch (error) {
      setMessage(`❌ Error limpiando: ${error}`)
      setTimeout(() => setMessage(''), 3000)
    }
  }

  const handleRefresh = () => {
    window.location.reload()
  }

  return (
    <div className="p-6 bg-white rounded-lg shadow-md max-w-md mx-auto">
      <h2 className="text-xl font-bold mb-4 text-gray-800">🧪 Prueba de Persistencia Local</h2>
      
      <div className="space-y-4">
        <div className="bg-gray-50 p-3 rounded">
          <p className="text-sm text-gray-600">Ruta de prueba: <span className="font-mono">{testRouteId}</span></p>
          <p className="text-sm text-gray-600">Índices: {testVisitIndex}-{testOrderIndex}-{testUnitIndex}</p>
        </div>

        <div className="bg-blue-50 p-3 rounded">
          <p className="text-sm text-blue-800">
            Estado actual: <span className="font-semibold">{currentStatus}</span>
          </p>
        </div>

        <div className="space-y-2">
          <button
            onClick={handleSetDelivered}
            className="w-full bg-green-500 hover:bg-green-600 text-white font-medium py-2 px-4 rounded transition-colors"
          >
            📦 Marcar como Entregado
          </button>
          
          <button
            onClick={handleSetNotDelivered}
            className="w-full bg-red-500 hover:bg-red-600 text-white font-medium py-2 px-4 rounded transition-colors"
          >
            ❌ Marcar como No Entregado
          </button>
        </div>

        <div className="border-t pt-4 space-y-2">
          <button
            onClick={handleSyncLocalBackup}
            className="w-full bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded transition-colors"
          >
            🔄 Sincronizar Respaldo Local
          </button>
          
          <button
            onClick={handleClearBackup}
            className="w-full bg-yellow-500 hover:bg-yellow-600 text-white font-medium py-2 px-4 rounded transition-colors"
          >
            🗑️ Limpiar Respaldo
          </button>
          
          <button
            onClick={handleRefresh}
            className="w-full bg-gray-500 hover:bg-gray-600 text-white font-medium py-2 px-4 rounded transition-colors"
          >
            🔄 Refrescar Página
          </button>
        </div>

        {message && (
          <div className="mt-4 p-3 bg-blue-100 border border-blue-300 rounded text-blue-800">
            {message}
          </div>
        )}

        <div className="mt-4 p-3 bg-gray-50 rounded text-sm text-gray-600">
          <p className="font-medium mb-2">📋 Instrucciones de prueba:</p>
          <ol className="list-decimal list-inside space-y-1">
            <li>Marca un estado (entregado/no entregado)</li>
            <li>Refresca la página</li>
            <li>Verifica que el estado se mantenga</li>
            <li>Usa "Sincronizar Respaldo Local" si es necesario</li>
          </ol>
        </div>
      </div>
    </div>
  )
}
