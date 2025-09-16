import { Package, Tag } from 'lucide-react'

interface IdentifierBadgeProps {
  lpn?: string
  referenceID?: string
  className?: string
  size?: 'sm' | 'md' | 'lg'
  showIcons?: boolean
}

export function IdentifierBadge({ 
  lpn, 
  referenceID,
  className = '', 
  size = 'md',
  showIcons = true 
}: IdentifierBadgeProps) {
  const sizeClasses = {
    sm: 'text-xs px-2 py-1',
    md: 'text-sm px-3 py-1.5',
    lg: 'text-base px-4 py-2'
  }

  const iconSizes = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5'
  }

  const hasLPN = lpn && lpn.trim() !== ''
  const hasReferenceID = referenceID && referenceID.trim() !== ''

  // Lógica del contrato:
  // - Si hay LPN: mostrar REFERENCE y LPN
  // - Si no hay LPN pero hay referenceID: mostrar solo REFERENCE
  // - Si no hay ninguno: no mostrar nada

  if (!hasLPN && !hasReferenceID) {
    return null
  }

  return (
    <div className={`flex items-center space-x-2 ${className}`}>
      {/* REFERENCE Badge - siempre se muestra si hay referenceID */}
      {hasReferenceID && (
        <div className="flex items-center space-x-1 bg-gradient-to-r from-orange-400 to-red-500 text-white rounded-lg font-medium shadow-sm">
          {showIcons && <Tag className={iconSizes[size]} />}
          <span className={sizeClasses[size]}>
            REFERENCE: {referenceID}
          </span>
        </div>
      )}

      {/* LPN Badge - solo se muestra si hay LPN */}
      {hasLPN && (
        <div className="flex items-center space-x-1 bg-gradient-to-r from-orange-400 to-red-500 text-white rounded-lg font-medium shadow-sm">
          {showIcons && <Package className={iconSizes[size]} />}
          <span className={sizeClasses[size]}>
            LPN: {lpn}
          </span>
        </div>
      )}
    </div>
  )
}
