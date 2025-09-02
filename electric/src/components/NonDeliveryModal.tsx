import { useRef, useState, useEffect } from 'react'
import { CameraCapture } from './CameraCapture'
import { processAndUploadImage, getUploadUrlFromRoute } from '../utils/imageUpload'
import type { DeliveryEvent } from '../domain/deliveries'
import type { Route as RouteType } from '../domain/route'

interface NonDeliveryModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (deliveryEvent: DeliveryEvent) => void
  initialDeliveryEvent?: DeliveryEvent // Para edición
  submitting?: boolean
  routeData?: RouteType // Para obtener URLs de evidencia
  visitIndex?: number
  orderIndex?: number
  unitIndex?: number
  isDemo?: boolean // Para modo demo
}

export function NonDeliveryModal({
  isOpen,
  onClose,
  onSubmit,
  initialDeliveryEvent,
  submitting = false,
  routeData,
  visitIndex,
  orderIndex,
  unitIndex,
  isDemo = false
}: NonDeliveryModalProps) {
  const [ndReasonQuery, setNdReasonQuery] = useState('')
  const [ndSelectedReason, setNdSelectedReason] = useState<string>('')
  const [ndObservations, setNdObservations] = useState<string>('')
  const [ndPhotoDataUrl, setNdPhotoDataUrl] = useState<string | null>(null)
  const [uploadingImage, setUploadingImage] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)
  const ndReasonInputRef = useRef<HTMLInputElement | null>(null)

  // Inicializar con datos existentes si los hay
  useEffect(() => {
    if (isOpen) {
      if (initialDeliveryEvent) {
        // Cargar datos existentes si estamos editando
        const failure = initialDeliveryEvent.deliveryUnits[0]?.delivery?.failure
        setNdSelectedReason(failure?.reason || '')
        setNdReasonQuery(failure?.reason || '')
        setNdObservations(failure?.detail || '')
        setNdPhotoDataUrl(initialDeliveryEvent.deliveryUnits[0]?.evidencePhotos[0]?.url || null)
      } else {
        // Limpiar para nueva no entrega
        setNdReasonQuery('')
        setNdSelectedReason('')
        setNdObservations('')
        setNdPhotoDataUrl(null)
      }
    }
  }, [isOpen, initialDeliveryEvent])

  const handleSubmit = async () => {
    const reason = (ndSelectedReason || ndReasonQuery || '').trim()
    if (!reason || !ndPhotoDataUrl) return
    
    try {
      setUploadingImage(true)
      setUploadError(null)
      
      let finalImageUrl = ndPhotoDataUrl
      
      if (!isDemo) {
        // Solo hacer upload si no es modo demo
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
        const { downloadUrl: uploadedDownloadUrl } = await processAndUploadImage(ndPhotoDataUrl, uploadUrl, downloadUrl || undefined)
        finalImageUrl = uploadedDownloadUrl
        console.log('✅ Imagen subida exitosamente:', finalImageUrl)
      } else {
        console.log('🎯 Modo demo: usando imagen local sin upload')
      }
      
      // ✅ Crear y retornar un DeliveryEvent hidratado
      const hydratedDeliveryEvent: DeliveryEvent = {
        ...initialDeliveryEvent!,
        deliveryUnits: initialDeliveryEvent?.deliveryUnits.map(unit => ({
          ...unit,
          delivery: {
            ...unit.delivery,
            status: 'not-delivered',
            failure: {
              reason,
              detail: ndObservations || '',
              referenceID: unit.orderReferenceID
            }
          },
          evidencePhotos: [{
            takenAt: new Date().toISOString(),
            type: 'non-delivery',
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
    setNdPhotoDataUrl(null)
    setNdReasonQuery('')
    setNdSelectedReason('')
    setNdObservations('')
    setUploadingImage(false)
    setUploadError(null)
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={handleClose}></div>
      <div className="relative bg-white w-full max-w-md mx-auto rounded-xl shadow-xl border border-gray-200 p-4 max-h-[85vh] overflow-y-auto">
        <h3 className="text-base font-semibold text-gray-800 mb-3">No entregado</h3>
        <div className="space-y-3">
          <div>
            <label className="block text-xs text-gray-600 mb-1">Motivo</label>
            <input
              type="text"
              value={ndReasonQuery}
              onChange={(e) => { setNdReasonQuery(e.target.value); setNdSelectedReason('') }}
              ref={ndReasonInputRef}
              className="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Buscar/ingresar motivo"
            />
            {/* Lista filtrada de motivos sugeridos */}
            {(() => {
              const base = ['cliente rechaza entrega', 'sin moradores', 'producto dañado', 'otro motivo']
              const q = ndReasonQuery.trim().toLowerCase()
              const items = base.filter((m) => m.includes(q))
              return (
                <div className="mt-2 max-h-40 overflow-auto border rounded-md">
                  {items.map((m) => (
                    <button
                      key={m}
                      type="button"
                      onClick={() => { setNdSelectedReason(m); setNdReasonQuery(m) }}
                      className={`w-full text-left px-3 py-2 text-sm hover:bg-gray-50 ${ndSelectedReason === m ? 'bg-indigo-50 text-indigo-700' : ''}`}
                    >
                      {m}
                    </button>
                  ))}
                </div>
              )
            })()}
          </div>
          <div>
            <label className="block text-xs text-gray-600 mb-1">Observaciones</label>
            <textarea
              value={ndObservations}
              onChange={(e) => setNdObservations(e.target.value)}
              rows={3}
              className="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Detalles adicionales (opcional)"
            />
          </div>
          <CameraCapture
            onCapture={(photoDataUrl) => setNdPhotoDataUrl(photoDataUrl)}
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
            disabled={submitting || uploadingImage || !(ndSelectedReason || ndReasonQuery).trim() || !ndPhotoDataUrl} 
            className={`px-3 py-2 text-sm rounded-md text-white ${submitting || uploadingImage || !(ndSelectedReason || ndReasonQuery).trim() || !ndPhotoDataUrl ? 'bg-red-300' : 'bg-red-600 hover:bg-red-700'}`}
          >
            {uploadingImage ? 'Subiendo imagen...' : 'Confirmar no entrega'}
          </button>
        </div>
      </div>
    </div>
  )
}
