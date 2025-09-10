import { useState, useEffect } from 'react'
import { X, Package, User, MapPin, CheckCircle, AlertCircle } from 'lucide-react'
import { CameraCapture } from './CameraCapture'
import { IdentifierBadge } from './IdentifierBadge'
import type { DeliveryGroup } from './GroupedDeliveryUtils'
import type { DeliveryEvent } from '../domain/deliveries'

interface GroupedDeliveryModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (deliveryEvent: DeliveryEvent) => void
  group: DeliveryGroup | null
  visitIndex: number
  routeData: any
  submitting?: boolean
  isDemo?: boolean
}

export function GroupedDeliveryModal({
  isOpen,
  onClose,
  onSubmit,
  group,
  visitIndex,
  routeData,
  submitting = false,
  isDemo: _isDemo = false
}: GroupedDeliveryModalProps) {
  const [recipientName, setRecipientName] = useState('')
  const [recipientRut, setRecipientRut] = useState('')
  const [photoDataUrl, setPhotoDataUrl] = useState<string | null>(null)

  // Inicializar con datos del grupo
  useEffect(() => {
    if (group && isOpen) {
      setRecipientName(group.addressInfo.contact?.fullName || '')
      setRecipientRut('')
      setPhotoDataUrl(null)
    }
  }, [group, isOpen])

  const handlePhotoCapture = (capturedPhotoDataUrl: string) => {
    setPhotoDataUrl(capturedPhotoDataUrl)
  }

  const handleSubmit = () => {
    if (!group || !recipientName.trim() || !photoDataUrl) {
      return
    }

    // Crear DeliveryEvent para todas las unidades del grupo
    const deliveryEvent: DeliveryEvent = {
      carrier: {
        name: '',
        nationalID: ''
      },
      deliveryUnits: group.units.map(unit => ({
        businessIdentifiers: {
          commerce: '',
          consumer: ''
        },
        delivery: {
          status: 'pending',
          handledAt: new Date().toISOString(),
          location: { latitude: 0, longitude: 0 }
        },
        evidencePhotos: [{
          takenAt: new Date().toISOString(),
          type: 'delivery',
          url: photoDataUrl,
        }],
        items: unit.unit.items || [],
        lpn: unit.unit.lpn || '',
        orderReferenceID: `${routeData?.documentID || 'route'}-${visitIndex}-${unit.orderIndex}-${unit.uIdx}`,
        recipient: {
          fullName: recipientName.trim(),
          nationalID: recipientRut.trim()
        }
      })),
      driver: {
        email: '',
        nationalID: ''
      },
      route: {
        id: routeData?.id || 0,
        documentID: routeData?.documentID || '',
        referenceID: routeData?.referenceID || '',
        sequenceNumber: 0,
        startedAt: new Date().toISOString()
      },
      vehicle: {
        plate: routeData?.vehicle?.plate || ''
      }
    }

    onSubmit(deliveryEvent)
  }

  if (!isOpen || !group) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <div className="flex items-center space-x-3">
            <div className="bg-blue-100 p-2 rounded-lg">
              <Package className="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900">Entregar todo junto</h2>
              <p className="text-sm text-gray-600">{group.totalUnits} unidades para {group.addressInfo.contact?.fullName}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Contenido */}
        <div className="p-6 space-y-6">
          {/* Información del grupo */}
          <div className="bg-blue-50 rounded-lg p-4">
            <div className="flex items-start space-x-3">
              <MapPin className="w-5 h-5 text-blue-600 mt-0.5" />
              <div>
                <h3 className="font-semibold text-blue-900">Dirección de entrega</h3>
                <p className="text-blue-700">{group.addressInfo.addressLine1}</p>
                <p className="text-sm text-blue-600 mt-1">
                  {group.totalUnits} unidades • {group.pendingUnits} pendientes
                </p>
              </div>
            </div>
          </div>

          {/* Lista de unidades */}
          <div>
            <h3 className="font-semibold text-gray-900 mb-3">Unidades a entregar</h3>
            <div className="space-y-2">
              {group.units.map((unit, index) => (
                <div key={index} className="flex items-center justify-between bg-gray-50 rounded-lg p-3">
                  <div className="flex items-center space-x-3">
                    <div className="bg-orange-100 p-2 rounded-lg">
                      <Package className="w-4 h-4 text-orange-600" />
                    </div>
                    <div>
                      <div className="mb-2">
                        <IdentifierBadge 
                          lpn={unit.unit.lpn} 
                          code={unit.unit.code} 
                          size="sm"
                        />
                      </div>
                      <div className="flex items-center space-x-2 mb-1">
                        <span className="text-xs font-medium text-gray-600">
                          Orden {unit.orderIndex + 1}
                        </span>
                      </div>
                      <p className="text-sm text-gray-600">
                        {unit.unit.items?.map((item: any) => item.description).join(', ') || 'Sin descripción'}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-600">
                      {unit.unit.items?.length || 0} items
                    </p>
                    {unit.status && (
                      <span className={`text-xs px-2 py-1 rounded-full ${
                        unit.status === 'delivered' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {unit.status === 'delivered' ? 'Entregado' : 'No entregado'}
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Información del receptor */}
          <div className="space-y-4">
            <h3 className="font-semibold text-gray-900">Información del receptor</h3>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <User className="w-4 h-4 inline mr-2" />
                Nombre completo *
              </label>
              <input
                type="text"
                value={recipientName}
                onChange={(e) => setRecipientName(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Nombre del receptor"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                RUT (opcional)
              </label>
              <input
                type="text"
                value={recipientRut}
                onChange={(e) => setRecipientRut(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="12.345.678-9"
              />
            </div>
          </div>

          {/* Captura de foto */}
          <div className="space-y-4">
            <h3 className="font-semibold text-gray-900">Evidencia fotográfica</h3>
            
            <CameraCapture
              onCapture={handlePhotoCapture}
              title="Capturar evidencia de entrega grupal"
              buttonText="Tomar foto"
              className="w-full"
            />
            
            {photoDataUrl && (
              <div className="space-y-3">
                <div className="relative">
                  <img
                    src={photoDataUrl}
                    alt="Evidencia de entrega grupal"
                    className="w-full h-48 object-cover rounded-lg border"
                  />
                  <button
                    onClick={() => setPhotoDataUrl(null)}
                    className="absolute top-2 right-2 bg-red-500 hover:bg-red-600 text-white p-1 rounded-full"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between p-6 border-t bg-gray-50 rounded-b-xl">
            <div className="flex items-center space-x-2 text-sm text-gray-600">
              <AlertCircle className="w-4 h-4" />
              <span>Esta acción marcará todas las unidades como entregadas</span>
            </div>
          
          <div className="flex space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Cancelar
            </button>
            <button
              onClick={handleSubmit}
              disabled={!recipientName.trim() || !photoDataUrl || submitting}
              className="px-6 py-2 bg-green-500 hover:bg-green-600 disabled:bg-gray-300 text-white rounded-lg font-medium transition-colors flex items-center space-x-2"
            >
              {submitting ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  <span>Procesando...</span>
                </>
              ) : (
                <>
                  <CheckCircle className="w-4 h-4" />
                  <span>Entregar todo</span>
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
