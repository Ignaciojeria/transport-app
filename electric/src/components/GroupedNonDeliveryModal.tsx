import { useState, useEffect } from 'react'
import { X, Package, AlertTriangle, MapPin, CheckCircle } from 'lucide-react'
import { CameraCapture } from './CameraCapture'
import { IdentifierBadge } from './IdentifierBadge'
import { useLanguage } from '../hooks/useLanguage'
import type { DeliveryGroup } from './GroupedDeliveryUtils'
import type { DeliveryEvent } from '../domain/deliveries'

interface GroupedNonDeliveryModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (deliveryEvent: DeliveryEvent) => void
  group: DeliveryGroup | null
  visitIndex: number
  routeData: any
  submitting?: boolean
  isDemo?: boolean
}


export function GroupedNonDeliveryModal({
  isOpen,
  onClose,
  onSubmit,
  group,
  visitIndex,
  routeData,
  submitting = false,
  isDemo: _isDemo = false
}: GroupedNonDeliveryModalProps) {
  const { t } = useLanguage()
  const [selectedReason, setSelectedReason] = useState('')
  const [customReason, setCustomReason] = useState('')
  const [photoDataUrl, setPhotoDataUrl] = useState<string | null>(null)

  // Motivos de no entrega traducidos
  const nonDeliveryReasons = [
    { value: 'recipient_not_available', label: t.groupedNonDeliveryModal.reasons.recipientNotAvailable },
    { value: 'address_not_found', label: t.groupedNonDeliveryModal.reasons.addressNotFound },
    { value: 'recipient_refused', label: t.groupedNonDeliveryModal.reasons.recipientRefused },
    { value: 'damaged_package', label: t.groupedNonDeliveryModal.reasons.damagedPackage },
    { value: 'incorrect_address', label: t.groupedNonDeliveryModal.reasons.incorrectAddress },
    { value: 'security_issue', label: t.groupedNonDeliveryModal.reasons.securityIssue },
    { value: 'other', label: t.groupedNonDeliveryModal.reasons.other }
  ]

  // Inicializar con datos del grupo
  useEffect(() => {
    if (group && isOpen) {
      setSelectedReason('')
      setCustomReason('')
      setPhotoDataUrl(null)
    }
  }, [group, isOpen])

  const handlePhotoCapture = (capturedPhotoDataUrl: string) => {
    setPhotoDataUrl(capturedPhotoDataUrl)
  }

  const handleSubmit = () => {
    if (!group || !selectedReason || !photoDataUrl) {
      return
    }

    const finalReason = selectedReason === 'other' ? customReason : selectedReason
    if (!finalReason.trim()) return

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
          location: { latitude: 0, longitude: 0 },
          failure: {
            reason: finalReason.trim(),
            detail: `No entrega grupal: ${finalReason.trim()}`,
            referenceID: `${routeData?.documentID || 'route'}-${visitIndex}-${unit.orderIndex}-${unit.uIdx}`
          }
        },
        evidencePhotos: [{
          takenAt: new Date().toISOString(),
          type: 'non-delivery',
          url: photoDataUrl,
        }],
        items: unit.unit.items || [],
        lpn: unit.unit.lpn || '',
        orderReferenceID: `${routeData?.documentID || 'route'}-${visitIndex}-${unit.orderIndex}-${unit.uIdx}`,
        recipient: {
          fullName: group.addressInfo.contact?.fullName || '',
          nationalID: ''
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
            <div className="bg-red-100 p-2 rounded-lg">
              <AlertTriangle className="w-6 h-6 text-red-600" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-gray-900">{t.delivery.notDeliverAll}</h2>
              <p className="text-sm text-gray-600">{group.totalUnits} {t.groupedNonDeliveryModal.unitsFor} {group.addressInfo.contact?.fullName}</p>
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
          <div className="bg-red-50 rounded-lg p-4">
            <div className="flex items-start space-x-3">
              <MapPin className="w-5 h-5 text-red-600 mt-0.5" />
              <div>
                <h3 className="font-semibold text-red-900">{t.groupedNonDeliveryModal.deliveryAddress}</h3>
                <p className="text-red-700">{group.addressInfo.addressLine1}</p>
                <p className="text-sm text-red-600 mt-1">
{group.totalUnits} {t.visitCard.units} • {group.pendingUnits} {t.groupedNonDeliveryModal.unitsPending}
                </p>
              </div>
            </div>
          </div>

          {/* Lista de unidades */}
          <div>
            <h3 className="font-semibold text-gray-900 mb-3">{t.groupedNonDeliveryModal.unitsNotToDeliver}</h3>
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
{t.groupedNonDeliveryModal.order} {unit.orderIndex + 1}
                        </span>
                      </div>
                      <p className="text-sm text-gray-600">
                        {unit.unit.items?.map((item: any) => item.description).join(', ') || 'Sin descripción'}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-600">
{unit.unit.items?.length || 0} {t.groupedNonDeliveryModal.items}
                    </p>
                    {unit.status && (
                      <span className={`text-xs px-2 py-1 rounded-full ${
                        unit.status === 'delivered' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
{unit.status === 'delivered' ? t.groupedNonDeliveryModal.delivered : t.groupedNonDeliveryModal.notDelivered}
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Motivo de no entrega */}
          <div className="space-y-4">
            <h3 className="font-semibold text-gray-900">{t.groupedNonDeliveryModal.nonDeliveryReason}</h3>
            
            <div className="grid grid-cols-1 gap-3">
              {nonDeliveryReasons.map((reason) => (
                <label key={reason.value} className="flex items-center space-x-3 p-3 border rounded-lg hover:bg-gray-50 cursor-pointer">
                  <input
                    type="radio"
                    name="reason"
                    value={reason.value}
                    checked={selectedReason === reason.value}
                    onChange={(e) => setSelectedReason(e.target.value)}
                    className="w-4 h-4 text-red-600 focus:ring-red-500"
                  />
                  <span className="text-sm font-medium text-gray-900">{reason.label}</span>
                </label>
              ))}
            </div>

            {selectedReason === 'other' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
{t.groupedNonDeliveryModal.specifyReason}
                </label>
                <textarea
                  value={customReason}
                  onChange={(e) => setCustomReason(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
                  placeholder={t.groupedNonDeliveryModal.reasonPlaceholder}
                  rows={3}
                />
              </div>
            )}
          </div>

          {/* Captura de foto */}
          <div className="space-y-4">
            <h3 className="font-semibold text-gray-900">{t.groupedNonDeliveryModal.photographicEvidence}</h3>
            
            <CameraCapture
              onCapture={handlePhotoCapture}
              title={t.groupedNonDeliveryModal.captureGroupEvidence}
              buttonText={t.groupedNonDeliveryModal.takePhoto}
              className="w-full"
            />
            
            {photoDataUrl && (
              <div className="space-y-3">
                <div className="relative">
                  <img
                    src={photoDataUrl}
                    alt="Evidencia de no entrega grupal"
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
            <AlertTriangle className="w-4 h-4" />
            <span>{t.groupedNonDeliveryModal.actionWarning}</span>
          </div>
          
          <div className="flex space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
{t.groupedNonDeliveryModal.cancel}
            </button>
            <button
              onClick={handleSubmit}
              disabled={!selectedReason || (selectedReason === 'other' && !customReason.trim()) || !photoDataUrl || submitting}
              className="px-6 py-2 bg-red-500 hover:bg-red-600 disabled:bg-gray-300 text-white rounded-lg font-medium transition-colors flex items-center space-x-2"
            >
              {submitting ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  <span>{t.groupedNonDeliveryModal.processing}</span>
                </>
              ) : (
                <>
                  <CheckCircle className="w-4 h-4" />
                  <span>{t.groupedNonDeliveryModal.markAsNotDelivered}</span>
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
