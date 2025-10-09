import React, { useState } from 'react'
import { type ElectricAccountData, type ElectricTenantData } from '../db/collections'
import TenantsList from './TenantsList'
import SideNav from './SideNav'
import Catalog from './Catalog'

interface OrganizationDashboardProps {
  account: ElectricAccountData
  tenants: ElectricTenantData[]
  token: string
}

const OrganizationDashboard: React.FC<OrganizationDashboardProps> = ({ 
  account, 
  tenants, 
  token 
}) => {
  const [selectedTenant, setSelectedTenant] = useState<ElectricTenantData | null>(null)
  const [activeMenuItem, setActiveMenuItem] = useState<string>('catalog')
  const [isSideNavOpen, setIsSideNavOpen] = useState(false)

  const handleTenantSelect = (tenant: ElectricTenantData) => {
    setSelectedTenant(tenant)
    setIsSideNavOpen(true)
    setActiveMenuItem('catalog') // Por defecto mostrar catalog
  }

  const handleBackToOrganizations = () => {
    setSelectedTenant(null)
    setIsSideNavOpen(false)
    setActiveMenuItem('catalog')
  }

  const handleMenuSelect = (menuItem: string) => {
    setActiveMenuItem(menuItem)
  }

  const handleCloseSideNav = () => {
    setIsSideNavOpen(false)
  }

  // Si hay una organización seleccionada, mostrar la vista con sidenav
  if (selectedTenant) {
    return (
      <div className="min-h-screen bg-gray-50">
        <SideNav
          tenant={selectedTenant}
          isOpen={isSideNavOpen}
          onClose={handleCloseSideNav}
          onMenuSelect={handleMenuSelect}
          activeMenuItem={activeMenuItem}
        />
        
        {/* Contenido principal */}
        <div className={`
          transition-all duration-300 ease-in-out
          ${isSideNavOpen ? 'lg:ml-64 xl:ml-72' : 'ml-0'}
        `}>
          {/* Header */}
          <div className="bg-white shadow-sm border-b border-gray-200 p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <button
                  onClick={() => setIsSideNavOpen(!isSideNavOpen)}
                  className="p-2 hover:bg-gray-100 rounded-lg transition-colors lg:hidden mr-3"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                  </svg>
                </button>
                <div>
                  <h1 className="text-xl font-semibold text-gray-800">
                    {selectedTenant.name}
                  </h1>
                  <p className="text-sm text-gray-500">{selectedTenant.country}</p>
                </div>
              </div>
              
              <button
                onClick={handleBackToOrganizations}
                className="px-4 py-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg transition-colors"
              >
                ← Volver a Organizaciones
              </button>
            </div>
          </div>

          {/* Contenido según el menú seleccionado */}
          <div className="p-0">
            {activeMenuItem === 'catalog' && (
              <Catalog tenant={selectedTenant} />
            )}
            {activeMenuItem === 'users' && (
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Usuarios</h2>
                <p className="text-gray-600">Gestión de usuarios - Próximamente</p>
              </div>
            )}
            {activeMenuItem === 'analytics' && (
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Analytics</h2>
                <p className="text-gray-600">Análisis y estadísticas - Próximamente</p>
              </div>
            )}
            {activeMenuItem === 'reports' && (
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Reportes</h2>
                <p className="text-gray-600">Generación de reportes - Próximamente</p>
              </div>
            )}
            {activeMenuItem === 'settings' && (
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Configuración</h2>
                <p className="text-gray-600">Configuración de la organización - Próximamente</p>
              </div>
            )}
          </div>
        </div>
      </div>
    )
  }

  // Si no hay organización seleccionada, mostrar la lista de organizaciones
  return (
    <TenantsList 
      account={account}
      tenants={tenants}
      token={token}
      onTenantSelect={handleTenantSelect}
    />
  )
}

export default OrganizationDashboard
