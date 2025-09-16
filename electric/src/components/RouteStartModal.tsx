import { useRef, useEffect, useState } from 'react'
import { useLanguage } from '../hooks/useLanguage'

interface RouteStartModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: (license: string) => void
  defaultLicense?: string
}

export function RouteStartModal({ isOpen, onClose, onConfirm, defaultLicense }: RouteStartModalProps) {
  const { t } = useLanguage()
  const licenseInputRef = useRef<HTMLInputElement | null>(null)
  const [inputValue, setInputValue] = useState('')

  // Focus en el input cuando se abre el modal
  useEffect(() => {
    if (isOpen && licenseInputRef.current) {
      licenseInputRef.current.focus()
    }
    // Limpiar el input cuando se abre el modal
    if (isOpen) {
      setInputValue('')
    }
  }, [isOpen])

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onClose}></div>
      <div className="relative bg-white w-full max-w-md mx-auto rounded-xl shadow-xl border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">{t.routeStartModal.title}</h3>
        
        {defaultLicense && (
          <div className="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
            <p className="text-sm text-blue-800 mb-2">
              {t.routeStartModal.useAssignedPlate}
            </p>
            <div className="flex items-center gap-2">
              <span className="font-mono text-base font-semibold text-blue-900">
                {defaultLicense}
              </span>
              <button
                onClick={() => onConfirm(defaultLicense)}
                className="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
              >
                {t.routeStartModal.useThisPlate}
              </button>
                             <button
                 onClick={() => {
                   // Limpiar el input para que ingrese otra patente
                   setInputValue('')
                   if (licenseInputRef.current) {
                     licenseInputRef.current.focus()
                   }
                 }}
                 className="px-2 py-1 text-xs bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
               >
                 {t.routeStartModal.useOtherPlate}
               </button>
            </div>
          </div>
        )}
        
        <div className="mb-4">
          <label htmlFor="licenseInput" className="block text-sm font-medium text-gray-700 mb-2">
            {defaultLicense ? t.routeStartModal.differentPlateLabel : t.routeStartModal.licensePlateLabel}
          </label>
                     <input
             id="licenseInput"
             ref={licenseInputRef}
             type="text"
             value={inputValue}
             onChange={(e) => {
               const value = e.target.value.toUpperCase()
               setInputValue(value)
             }}
             onKeyPress={(e) => {
               if (e.key === 'Enter') {
                 const value = e.currentTarget.value.trim()
                 if (value) {
                   onConfirm(value)
                 }
               }
             }}
             placeholder={t.routeStartModal.plateExample}
             className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent uppercase"
             maxLength={8}
           />
          <p className="text-xs text-gray-500 mt-1">
            {defaultLicense 
              ? t.routeStartModal.plateHelpText
              : t.routeStartModal.plateHelpTextDefault
            }
          </p>
        </div>
        


        <div className="flex items-center justify-end gap-3">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm rounded-lg border border-gray-300 bg-white hover:bg-gray-50 text-gray-700 transition-colors"
          >
            {t.routeStartModal.cancelButton}
          </button>
                     <button
             onClick={() => {
               const value = inputValue.trim()
               if (value) {
                 onConfirm(value)
               }
             }}
             disabled={!inputValue.trim()}
             className={`px-4 py-2 text-sm rounded-lg text-white font-medium transition-colors ${
               !inputValue.trim()
                 ? 'bg-green-300 cursor-not-allowed'
                 : 'bg-green-600 hover:bg-green-700'
             }`}
           >
             {t.routeStartModal.startButton}
           </button>
        </div>
      </div>
    </div>
  )
}
