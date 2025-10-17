import React from 'react'
import { type ElectricTenantData } from '../db/collections'
import { 
  Package, 
  Users, 
  Settings, 
  BarChart3, 
  FileText,
  X
} from 'lucide-react'

interface SideNavProps {
  tenant: ElectricTenantData
  isOpen: boolean
  onClose: () => void
  onMenuSelect: (menuItem: string) => void
  activeMenuItem: string
}

const SideNav: React.FC<SideNavProps> = ({ 
  tenant, 
  isOpen, 
  onClose, 
  onMenuSelect, 
  activeMenuItem 
}) => {
  const menuItems = [
    { id: 'catalog', label: 'Catalogo', icon: Package },
    { id: 'users', label: 'Usuarios', icon: Users },
    { id: 'analytics', label: 'Analytics', icon: BarChart3 },
    { id: 'reports', label: 'Reportes', icon: FileText },
    { id: 'settings', label: 'Configuración', icon: Settings },
  ]

  return (
    <>
      {/* Overlay para móviles */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* SideNav */}
      <div className={`
        fixed top-0 left-0 h-full bg-white shadow-xl z-50 transform transition-transform duration-300 ease-in-out
        ${isOpen ? 'translate-x-0' : '-translate-x-full'}
        w-64 lg:w-72
      `}>
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <div>
            <h2 className="text-lg font-semibold text-gray-800">{tenant.name}</h2>
            <p className="text-sm text-gray-500">{tenant.country}</p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors lg:hidden"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Menu Items */}
        <nav className="mt-4">
          {menuItems.map((item) => {
            const Icon = item.icon
            const isActive = activeMenuItem === item.id
            
            return (
              <button
                key={item.id}
                onClick={() => onMenuSelect(item.id)}
                className={`
                  w-full flex items-center px-4 py-3 text-left transition-colors
                  ${isActive 
                    ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-700' 
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-800'
                  }
                `}
              >
                <Icon className="w-5 h-5 mr-3" />
                <span className="font-medium">{item.label}</span>
              </button>
            )
          })}
        </nav>

        {/* Footer */}
        <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200">
          <div className="text-xs text-gray-500 text-center">
            <p>ID: {tenant.id.slice(0, 8)}...</p>
            <p className="mt-1">Waiting List</p>
          </div>
        </div>
      </div>
    </>
  )
}

export default SideNav
