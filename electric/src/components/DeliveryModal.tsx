import { useRef, useState, useEffect } from 'react'
import { CameraCapture } from './CameraCapture'
import { processAndUploadImage, getUploadUrlFromRoute } from '../utils/imageUpload'
import type { DeliveryEvent } from '../domain/deliveries'
import type { Route as RouteType } from '../domain/route'

interface DeliveryModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (deliveryEvent: DeliveryEvent) => void
  initialDeliveryEvent?: DeliveryEvent // Para edición
  submitting?: boolean
  routeData?: RouteType // Para obtener URLs de evidencia
  visitIndex?: number
  orderIndex?: number
  unitIndex?: number
}

export function DeliveryModal({
  isOpen,
  onClose,
  onSubmit,
  initialDeliveryEvent,
  submitting = false,
  routeData,
  visitIndex,
  orderIndex,
  unitIndex
}: DeliveryModalProps) {
  const [recipientName, setRecipientName] = useState('')
  const [recipientRut, setRecipientRut] = useState('')
  const [photoDataUrl, setPhotoDataUrl] = useState<string | null>(null)
  const [uploadingImage, setUploadingImage] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)
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

  const handleSubmit = async () => {
    const trimmedName = recipientName.trim()
    const trimmedRut = recipientRut.trim()
    if (!trimmedName || !trimmedRut || !photoDataUrl) return
    
    try {
      setUploadingImage(true)
      setUploadError(null)
      
      // Obtener URLs del contrato de ruta
      const { uploadUrl, downloadUrl } = getUploadUrlFromRoute(
        routeData, 
        visitIndex || 0, 
        orderIndex || 0, 
        unitIndex || 0
      )
      
      if (!uploadUrl) {
        throw new Error('No se encontró uploadUrl en el contrato de ruta')
      }
      
      console.log('📤 Subiendo imagen usando URL firmada del contrato...')
      const { downloadUrl: uploadedDownloadUrl } = await processAndUploadImage(photoDataUrl, uploadUrl, downloadUrl || undefined)
      const finalImageUrl = uploadedDownloadUrl
      console.log('✅ Imagen subida exitosamente:', finalImageUrl)
      
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
            url: finalImageUrl
          }]
        })) || []
      }
      
      onSubmit(hydratedDeliveryEvent)
      
    } catch (error) {
      console.error('❌ Error subiendo imagen:', error)
      setUploadError(error instanceof Error ? error.message : 'Error subiendo imagen')
    } finally {
      setUploadingImage(false)
    }
  }

  const handleClose = () => {
    setRecipientName('')
    setRecipientRut('')
    setPhotoDataUrl(null)
    setUploadingImage(false)
    setUploadError(null)
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
        {uploadError && (
          <div className="mt-3 p-2 bg-red-50 border border-red-200 rounded-md">
            <p className="text-xs text-red-600">{uploadError}</p>
          </div>
        )}
        <div className="mt-4 flex items-center justify-end gap-2">
          <button
            onClick={handleClose}
            className="px-3 py-2 text-sm rounded-md border bg-white hover:bg-gray-50"
            disabled={submitting || uploadingImage}
          >
            Cancelar
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting || uploadingImage || !recipientName.trim() || !recipientRut.trim() || !photoDataUrl}
            className={`px-3 py-2 text-sm rounded-md text-white ${submitting || uploadingImage || !recipientName.trim() || !recipientRut.trim() || !photoDataUrl ? 'bg-green-300' : 'bg-green-600 hover:bg-green-700'}`}
          >
            {uploadingImage ? 'Subiendo imagen...' : 'Confirmar entrega'}
          </button>
        </div>
      </div>
    </div>
  )
}
