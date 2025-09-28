import { useState } from 'react'
import { Truck } from 'lucide-react'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from './ui/Card'

interface CreateOrganizationProps {
  onSubmit?: (data: { name: string; country: string }) => void
}

const countries = [
  'Argentina',
  'Bolivia',
  'Brasil',
  'Chile',
  'Colombia',
  'Ecuador',
  'Paraguay',
  'Perú',
  'Uruguay',
  'Venezuela'
]

export default function CreateOrganization({ onSubmit }: CreateOrganizationProps) {
  const [organizationName, setOrganizationName] = useState('')
  const [selectedCountry, setSelectedCountry] = useState('Chile')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!organizationName.trim()) return

    setIsSubmitting(true)
    
    try {
      await onSubmit?.({
        name: organizationName.trim(),
        country: selectedCountry
      })
    } catch (error) {
      console.error('Error creating organization:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50 flex items-center justify-center p-4 md:p-6 relative overflow-hidden">
      {/* Partículas de fondo más visibles y claras */}
      <div className="absolute inset-0 pointer-events-none">
        {/* Partículas grandes principales */}
        <div className="absolute top-20 left-10 w-6 h-6 bg-blue-400 rounded-full animate-float opacity-90 shadow-lg"></div>
        <div className="absolute top-40 right-20 w-7 h-7 bg-indigo-400 rounded-full animate-float opacity-95 shadow-lg" style={{animationDelay: '1s'}}></div>
        <div className="absolute bottom-40 left-20 w-5 h-5 bg-blue-500 rounded-full animate-float opacity-85 shadow-lg" style={{animationDelay: '2s'}}></div>
        <div className="absolute bottom-20 right-10 w-6 h-6 bg-indigo-500 rounded-full animate-float opacity-90 shadow-lg" style={{animationDelay: '3s'}}></div>
        
        {/* Partículas medianas */}
        <div className="absolute top-60 left-1/3 w-5 h-5 bg-blue-400 rounded-full animate-float opacity-80 shadow-md" style={{animationDelay: '0.5s'}}></div>
        <div className="absolute bottom-60 right-1/3 w-6 h-6 bg-indigo-400 rounded-full animate-float opacity-85 shadow-md" style={{animationDelay: '1.5s'}}></div>
        <div className="absolute top-1/2 left-1/4 w-4 h-4 bg-blue-300 rounded-full animate-float opacity-75 shadow-md" style={{animationDelay: '2.5s'}}></div>
        <div className="absolute bottom-1/3 right-1/4 w-5 h-5 bg-indigo-300 rounded-full animate-float opacity-80 shadow-md" style={{animationDelay: '0.8s'}}></div>
        <div className="absolute top-1/3 right-1/3 w-4 h-4 bg-blue-400 rounded-full animate-float opacity-70 shadow-md" style={{animationDelay: '1.2s'}}></div>
        <div className="absolute bottom-1/2 left-1/2 w-5 h-5 bg-indigo-400 rounded-full animate-float opacity-85 shadow-md" style={{animationDelay: '2.8s'}}></div>
        
        {/* Partículas pequeñas para detalle */}
        <div className="absolute top-1/4 left-1/5 w-3 h-3 bg-blue-500 rounded-full animate-float-small opacity-70 shadow-sm" style={{animationDelay: '1.8s'}}></div>
        <div className="absolute bottom-1/4 right-1/5 w-4 h-4 bg-indigo-500 rounded-full animate-float-small opacity-75 shadow-sm" style={{animationDelay: '3.5s'}}></div>
        <div className="absolute top-3/4 left-1/6 w-3 h-3 bg-blue-300 rounded-full animate-float-small opacity-65 shadow-sm" style={{animationDelay: '0.3s'}}></div>
        <div className="absolute bottom-3/4 right-1/6 w-4 h-4 bg-indigo-300 rounded-full animate-float-small opacity-70 shadow-sm" style={{animationDelay: '2.2s'}}></div>
        
        {/* Partículas adicionales para más densidad */}
        <div className="absolute top-1/6 left-2/3 w-3 h-3 bg-blue-400 rounded-full animate-float-small opacity-60 shadow-sm" style={{animationDelay: '1.7s'}}></div>
        <div className="absolute bottom-1/6 right-2/3 w-4 h-4 bg-indigo-400 rounded-full animate-float-small opacity-65 shadow-sm" style={{animationDelay: '3.1s'}}></div>
        <div className="absolute top-2/3 left-1/6 w-3 h-3 bg-blue-500 rounded-full animate-float-small opacity-55 shadow-sm" style={{animationDelay: '0.9s'}}></div>
        <div className="absolute bottom-2/3 right-1/6 w-4 h-4 bg-indigo-500 rounded-full animate-float-small opacity-60 shadow-sm" style={{animationDelay: '2.4s'}}></div>
      </div>
      <div className="w-full max-w-2xl relative z-10">
        {/* Logo y título */}
        <div className="text-center mb-10">
          <div className="flex items-center justify-center space-x-6 mb-8">
            <div className="bg-blue-600 p-5 rounded-2xl">
              <Truck className="h-12 w-12 text-white" />
            </div>
            <div>
              <h1 className="text-5xl font-bold text-gray-900">TransportApp</h1>
              <p className="text-gray-600 text-xl">Configuración inicial</p>
            </div>
          </div>
        </div>

        {/* Mensaje de bienvenida */}
        <div className="bg-green-50 border border-green-200 rounded-xl p-6 mb-10">
          <div className="flex items-center space-x-4">
            <div className="w-10 h-10 bg-green-500 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"/>
              </svg>
            </div>
            <div>
              <p className="text-green-800 font-semibold text-xl">¡Bienvenido, Usuario de Prueba!</p>
              <p className="text-green-700 text-lg">Autenticación exitosa</p>
            </div>
          </div>
        </div>

        {/* Card principal */}
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-5xl mb-6">
              Crear Organización
            </CardTitle>
            <CardDescription className="text-2xl">
              Configura tu organización para comenzar a gestionar tus operaciones logísticas
            </CardDescription>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-8">
              {/* Nombre de la organización */}
              <div>
                <label htmlFor="organizationName" className="block text-lg font-medium text-gray-700 mb-4">
                  Nombre de la Organización
                </label>
                <input
                  type="text"
                  id="organizationName"
                  value={organizationName}
                  onChange={(e) => setOrganizationName(e.target.value)}
                  placeholder="hass me :D"
                  className="w-full px-5 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors text-xl"
                  required
                />
              </div>

              {/* País */}
              <div>
                <label htmlFor="country" className="block text-lg font-medium text-gray-700 mb-4">
                  País de Operación Logística
                </label>
                <select
                  id="country"
                  value={selectedCountry}
                  onChange={(e) => setSelectedCountry(e.target.value)}
                  className="w-full px-5 py-4 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors appearance-none bg-white text-xl"
                >
                  {countries.map((country) => (
                    <option key={country} value={country}>
                      {country}
                    </option>
                  ))}
                </select>
              </div>

              {/* Botón de envío */}
              <button
                type="submit"
                disabled={isSubmitting || !organizationName.trim()}
                className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold py-5 px-8 rounded-xl transition-colors duration-200 flex items-center justify-center gap-3 text-xl"
              >
                {isSubmitting ? (
                  <>
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-white" />
                    <span>Creando...</span>
                  </>
                ) : (
                  <>
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M13 10V3L4 14h7v7l9-11h-7z"/>
                    </svg>
                    <span>Crear Organización</span>
                  </>
                )}
              </button>
            </form>
          </CardContent>

          <CardFooter>
            {/* Información adicional */}
            <div className="w-full bg-blue-50 border border-blue-200 rounded-xl p-6">
              <div className="flex items-start space-x-4">
                <svg className="w-6 h-6 text-blue-600 mt-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd"/>
                </svg>
                <div>
                  <p className="text-blue-800 font-medium text-lg">Configuración Inicial</p>
                  <p className="text-blue-700 mt-2 text-lg">
                    Esta información se utilizará para configurar tu espacio de trabajo y optimizar las rutas según tu región
                  </p>
                </div>
              </div>
            </div>
          </CardFooter>
        </Card>

        {/* Footer */}
        <div className="mt-8 text-center text-sm text-gray-400">
          <p>© 2024 TransportApp. Todos los derechos reservados.</p>
        </div>
      </div>
    </div>
  )
}
