import React from 'react'
import { type ElectricAccountData, type ElectricTenantData } from '../db/collections'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/Card'
import { Building2, MapPin, Users, Calendar } from 'lucide-react'

interface TenantsListProps {
  account: ElectricAccountData
  tenants: ElectricTenantData[]
  token: string
}

const TenantsList: React.FC<TenantsListProps> = ({ account, tenants, token }) => {
  const handleTenantSelect = (tenant: ElectricTenantData) => {
    // Aquí puedes implementar la lógica para seleccionar un tenant
    console.log('Tenant seleccionado:', tenant)
    // Por ejemplo, redirigir a la aplicación principal con el tenant seleccionado
    // window.location.href = `/app?tenant=${tenant.id}`
  }

  const handleCreateNewOrganization = () => {
    // Redirigir al formulario de creación de organización
    window.location.href = '/create-organization'
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-4">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">
            Bienvenido, {account.email}
          </h1>
          <p className="text-gray-600">
            Selecciona una organización para continuar
          </p>
        </div>

        {/* Botón para crear nueva organización */}
        <div className="mb-8 text-center">
          <button
            onClick={handleCreateNewOrganization}
            className="inline-flex items-center px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors shadow-lg"
          >
            <Building2 className="w-5 h-5 mr-2" />
            Crear Nueva Organización
          </button>
        </div>

        {/* Lista de tenants */}
        {tenants.length === 0 ? (
          <Card className="max-w-2xl mx-auto">
            <CardHeader>
              <CardTitle className="text-center">No hay organizaciones</CardTitle>
              <CardDescription className="text-center">
                Aún no tienes organizaciones asociadas a tu cuenta
              </CardDescription>
            </CardHeader>
            <CardContent className="text-center">
              <button
                onClick={handleCreateNewOrganization}
                className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                Crear Primera Organización
              </button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {tenants.map((tenant) => (
              <div 
                key={tenant.id}
                className="cursor-pointer hover:shadow-xl transition-all duration-300 hover:scale-105"
                onClick={() => handleTenantSelect(tenant)}
              >
                <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <Building2 className="w-8 h-8 text-blue-500" />
                    <span className="text-sm text-gray-500">ID: {tenant.id.slice(0, 8)}...</span>
                  </div>
                  <CardTitle className="text-xl">{tenant.name}</CardTitle>
                  <CardDescription>
                    Organización en {tenant.country}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center text-gray-600">
                      <MapPin className="w-4 h-4 mr-2" />
                      <span className="text-sm">{tenant.country}</span>
                    </div>
                    <div className="flex items-center text-gray-600">
                      <Users className="w-4 h-4 mr-2" />
                      <span className="text-sm">Miembro activo</span>
                    </div>
                    {tenant.created_at && (
                      <div className="flex items-center text-gray-600">
                        <Calendar className="w-4 h-4 mr-2" />
                        <span className="text-sm">
                          Creado: {new Date(tenant.created_at).toLocaleDateString()}
                        </span>
                      </div>
                    )}
                    <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-2 mt-3">
                      <p className="text-yellow-800 text-xs font-medium text-center">
                        📋 Waiting List
                      </p>
                    </div>
                  </div>
                </CardContent>
                </Card>
              </div>
            ))}
          </div>
        )}

        {/* Información de debug (solo en desarrollo) */}
        {import.meta.env.DEV && (
          <div className="mt-8 p-4 bg-gray-100 rounded-lg">
            <h3 className="font-bold mb-2">Debug Info:</h3>
            <p><strong>Account ID:</strong> {account.id}</p>
            <p><strong>Email:</strong> {account.email}</p>
            <p><strong>Tenants count:</strong> {tenants.length}</p>
            <p><strong>Token:</strong> {token.slice(0, 20)}...</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default TenantsList
