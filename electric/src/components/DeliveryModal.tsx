import { useRef, useState, useEffect } from 'react'
import { CameraCapture } from './CameraCapture'
import type { DeliveryEvent } from '../domain/deliveries'

interface DeliveryModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (deliveryEvent: DeliveryEvent) => void
  initialDeliveryEvent?: DeliveryEvent // Para edición
  submitting?: boolean
}

export function DeliveryModal({
  isOpen,
  onClose,
  onSubmit,
  initialDeliveryEvent,
  submitting = false
}: DeliveryModalProps) {
  const [recipientName, setRecipientName] = useState('')
  const [recipientRut, setRecipientRut] = useState('')
  const [photoDataUrl, setPhotoDataUrl] = useState<string | null>(null)
  const nameInputRef = useRef<HTMLInputElement | null>(null)
  const rutInputRef = useRef<HTMLInputElement | null>(null)

  // Inicializar con datos existentes si los hay
  useEffect(() => {
    if (isOpen) {
      if (initialDeliveryEvent) {
        setRecipientName(initialDeliveryEvent.deliveryUnits[0]?.recipient?.fullName || '')
        setRecipientRut(initialDeliveryEvent.deliveryUnits[0]?.recipient?.nationalID || '')
        setPhotoDataUrl(initialDeliveryEvent.deliveryUnits[0]?.evidencePhotos[0]?.url || null)
      } else {
        // Limpiar para nueva entrega
        setRecipientName('')
        setRecipientRut('')
        setPhotoDataUrl(null)
      }
    }
  }, [isOpen, initialDeliveryEvent])

  const handleSubmit = () => {
    const trimmedName = recipientName.trim()
    const trimmedRut = recipientRut.trim()
    if (!trimmedName || !trimmedRut || !photoDataUrl) return
    
    // ✅ Crear y retornar un DeliveryEvent hidratado
    const hydratedDeliveryEvent: DeliveryEvent = {
      ...initialDeliveryEvent!,
      deliveryUnits: initialDeliveryEvent?.deliveryUnits.map(unit => ({
        ...unit,
        recipient: {
          fullName: trimmedName,
          nationalID: trimmedRut
        },
        evidencePhotos: [{
          takenAt: new Date().toISOString(),
          type: 'delivery_evidence',
          url: photoDataUrl
        }]
      })) || []
    }
    
    onSubmit(hydratedDeliveryEvent)
  }

  const handleClose = () => {
    setRecipientName('')
    setRecipientRut('')
    setPhotoDataUrl(null)
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={handleClose}></div>
      <div className="relative bg-white w-full max-w-md mx-auto rounded-xl shadow-xl border border-gray-200 p-4 max-h-[85vh] overflow-y-auto">
        <h3 className="text-base font-semibold text-gray-800 mb-3">Evidencia de entrega</h3>
        <div className="space-y-3">
          <div>
            <label className="block text-xs text-gray-600 mb-1">Nombre de quien recibe</label>
            <input
              type="text"
              value={recipientName}
              onChange={(e) => setRecipientName(e.target.value)}
              ref={nameInputRef}
              className="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Nombre completo"
            />
          </div>
          <div>
            <label className="block text-xs text-gray-600 mb-1">RUT / Documento</label>
            <input
              type="text"
              value={recipientRut}
              onChange={(e) => setRecipientRut(e.target.value)}
              ref={rutInputRef}
              className="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="12.345.678-9"
            />
          </div>
          <CameraCapture
            onCapture={(photoDataUrl) => setPhotoDataUrl(photoDataUrl)}
            title="Foto de evidencia"
            buttonText="Activar cámara"
          />
        </div>
        <div className="mt-4 flex items-center justify-end gap-2">
          <button
            onClick={handleClose}
            className="px-3 py-2 text-sm rounded-md border bg-white hover:bg-gray-50"
            disabled={submitting}
          >
            Cancelar
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting || !recipientName.trim() || !recipientRut.trim() || !photoDataUrl}
            className={`px-3 py-2 text-sm rounded-md text-white ${submitting || !recipientName.trim() || !recipientRut.trim() || !photoDataUrl ? 'bg-green-300' : 'bg-green-600 hover:bg-green-700'}`}
          >
            Confirmar entrega
          </button>
        </div>
      </div>
    </div>
  )
}
