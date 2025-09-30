import React, { useEffect, useState } from 'react'
import { CheckCircle, Loader2 } from 'lucide-react'

interface SuccessNotificationProps {
  message: string
  onComplete?: () => void
  duration?: number
}

const SuccessNotification: React.FC<SuccessNotificationProps> = ({ 
  message, 
  onComplete, 
  duration = 2000 
}) => {
  const [show, setShow] = useState(true)
  const [isProcessing, setIsProcessing] = useState(false)

  useEffect(() => {
    // Mostrar el mensaje de éxito
    const successTimer = setTimeout(() => {
      setIsProcessing(true)
    }, 1000)

    // Completar después del tiempo especificado
    const completeTimer = setTimeout(() => {
      setShow(false)
      onComplete?.()
    }, duration)

    return () => {
      clearTimeout(successTimer)
      clearTimeout(completeTimer)
    }
  }, [duration, onComplete])

  if (!show) return null

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-white rounded-xl shadow-2xl p-8 max-w-md mx-auto text-center">
        {!isProcessing ? (
          <>
            <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle className="w-8 h-8 text-white" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">¡Éxito!</h2>
            <p className="text-gray-600">{message}</p>
          </>
        ) : (
          <>
            <div className="w-16 h-16 bg-blue-500 rounded-full flex items-center justify-center mx-auto mb-4">
              <Loader2 className="w-8 h-8 text-white animate-spin" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Procesando...</h2>
            <p className="text-gray-600">Verificando cuenta y cargando organizaciones...</p>
          </>
        )}
      </div>
    </div>
  )
}

export default SuccessNotification
