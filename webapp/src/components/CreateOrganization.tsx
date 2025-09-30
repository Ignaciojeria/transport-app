import { useState } from 'react'
import { Truck, AlertCircle, CheckCircle } from 'lucide-react'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from './ui/Card'
import { createOrganization, validateOrganizationData, type CreateOrganizationResponse } from '../services/organizationService'

interface CreateOrganizationProps {
  email: string
  onSuccess?: (response: CreateOrganizationResponse) => void
  onError?: (error: string) => void
}

const countries = [
  { name: 'Argentina', code: 'AR' },
  { name: 'Brasil', code: 'BR' },
  { name: 'Chile', code: 'CL' },
  { name: 'Colombia', code: 'CO' },
  { name: 'Ecuador', code: 'EC' },
  { name: 'Paraguay', code: 'PY' },
  { name: 'Perú', code: 'PE' },
  { name: 'Uruguay', code: 'UY' },
  { name: 'Venezuela', code: 'VE' }
]

export default function CreateOrganization({ email, onSuccess, onError }: CreateOrganizationProps) {
  const [organizationName, setOrganizationName] = useState('')
  const [selectedCountry, setSelectedCountry] = useState('CL')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})
  const [submitResponse, setSubmitResponse] = useState<CreateOrganizationResponse | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!email) {
      onError?.('Email no disponible')
      return
    }

    const formData = {
      email,
      organizationName: organizationName.trim(),
      country: selectedCountry
    }

    // Validar datos
    const validation = validateOrganizationData(formData)
    if (!validation.isValid) {
      setValidationErrors(validation.errors)
      return
    }

    setValidationErrors({})
    setIsSubmitting(true)
    setSubmitResponse(null)
    
    try {
      const response = await createOrganization(formData)
      setSubmitResponse(response)
      
      if (response.success) {
        onSuccess?.(response)
      } else {
        onError?.(response.error || 'Error al crear la organización')
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Error desconocido'
      onError?.(errorMessage)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50 flex items-center justify-center p-4 md:p-6 relative overflow-hidden">
      {/* Partículas de fondo más pequeñas y numerosas */}
      <div className="absolute inset-0 pointer-events-none">
        {/* Partículas medianas principales */}
        <div className="absolute top-20 left-10 w-3 h-3 bg-blue-400 rounded-full animate-float opacity-80 shadow-md"></div>
        <div className="absolute top-40 right-20 w-4 h-4 bg-indigo-400 rounded-full animate-float opacity-85 shadow-md" style={{animationDelay: '1s'}}></div>
        <div className="absolute bottom-40 left-20 w-3 h-3 bg-blue-500 rounded-full animate-float opacity-75 shadow-md" style={{animationDelay: '2s'}}></div>
        <div className="absolute bottom-20 right-10 w-3 h-3 bg-indigo-500 rounded-full animate-float opacity-80 shadow-md" style={{animationDelay: '3s'}}></div>
        
        {/* Partículas pequeñas principales */}
        <div className="absolute top-60 left-1/3 w-2 h-2 bg-blue-400 rounded-full animate-float-small opacity-70 shadow-sm" style={{animationDelay: '0.5s'}}></div>
        <div className="absolute bottom-60 right-1/3 w-3 h-3 bg-indigo-400 rounded-full animate-float-small opacity-75 shadow-sm" style={{animationDelay: '1.5s'}}></div>
        <div className="absolute top-1/2 left-1/4 w-2 h-2 bg-blue-300 rounded-full animate-float-small opacity-65 shadow-sm" style={{animationDelay: '2.5s'}}></div>
        <div className="absolute bottom-1/3 right-1/4 w-3 h-3 bg-indigo-300 rounded-full animate-float-small opacity-70 shadow-sm" style={{animationDelay: '0.8s'}}></div>
        <div className="absolute top-1/3 right-1/3 w-2 h-2 bg-blue-400 rounded-full animate-float-small opacity-60 shadow-sm" style={{animationDelay: '1.2s'}}></div>
        <div className="absolute bottom-1/2 left-1/2 w-3 h-3 bg-indigo-400 rounded-full animate-float-small opacity-75 shadow-sm" style={{animationDelay: '2.8s'}}></div>
        
        {/* Partículas muy pequeñas para detalle */}
        <div className="absolute top-1/4 left-1/5 w-1.5 h-1.5 bg-blue-500 rounded-full animate-float-small opacity-60 shadow-sm" style={{animationDelay: '1.8s'}}></div>
        <div className="absolute bottom-1/4 right-1/5 w-2 h-2 bg-indigo-500 rounded-full animate-float-small opacity-65 shadow-sm" style={{animationDelay: '3.5s'}}></div>
        <div className="absolute top-3/4 left-1/6 w-1.5 h-1.5 bg-blue-300 rounded-full animate-float-small opacity-55 shadow-sm" style={{animationDelay: '0.3s'}}></div>
        <div className="absolute bottom-3/4 right-1/6 w-2 h-2 bg-indigo-300 rounded-full animate-float-small opacity-60 shadow-sm" style={{animationDelay: '2.2s'}}></div>
        
        {/* Partículas adicionales para más densidad */}
        <div className="absolute top-1/6 left-2/3 w-1.5 h-1.5 bg-blue-400 rounded-full animate-float-small opacity-50 shadow-sm" style={{animationDelay: '1.7s'}}></div>
        <div className="absolute bottom-1/6 right-2/3 w-2 h-2 bg-indigo-400 rounded-full animate-float-small opacity-55 shadow-sm" style={{animationDelay: '3.1s'}}></div>
        <div className="absolute top-2/3 left-1/6 w-1.5 h-1.5 bg-blue-500 rounded-full animate-float-small opacity-45 shadow-sm" style={{animationDelay: '0.9s'}}></div>
        <div className="absolute bottom-2/3 right-1/6 w-2 h-2 bg-indigo-500 rounded-full animate-float-small opacity-50 shadow-sm" style={{animationDelay: '2.4s'}}></div>
        
        {/* Más partículas pequeñas para mayor densidad */}
        <div className="absolute top-1/8 left-1/8 w-1 h-1 bg-blue-400 rounded-full animate-float-small opacity-40 shadow-sm" style={{animationDelay: '0.7s'}}></div>
        <div className="absolute bottom-1/8 right-1/8 w-1.5 h-1.5 bg-indigo-400 rounded-full animate-float-small opacity-45 shadow-sm" style={{animationDelay: '2.9s'}}></div>
        <div className="absolute top-5/8 left-3/8 w-1 h-1 bg-blue-300 rounded-full animate-float-small opacity-35 shadow-sm" style={{animationDelay: '1.1s'}}></div>
        <div className="absolute bottom-5/8 right-3/8 w-1.5 h-1.5 bg-indigo-300 rounded-full animate-float-small opacity-40 shadow-sm" style={{animationDelay: '3.3s'}}></div>
        <div className="absolute top-7/8 left-5/8 w-1 h-1 bg-blue-500 rounded-full animate-float-small opacity-30 shadow-sm" style={{animationDelay: '0.4s'}}></div>
        <div className="absolute bottom-7/8 right-5/8 w-1.5 h-1.5 bg-indigo-500 rounded-full animate-float-small opacity-35 shadow-sm" style={{animationDelay: '2.6s'}}></div>
        
        {/* Partículas extra para completar el efecto */}
        <div className="absolute top-3/8 left-7/8 w-1 h-1 bg-blue-400 rounded-full animate-float-small opacity-25 shadow-sm" style={{animationDelay: '1.4s'}}></div>
        <div className="absolute bottom-3/8 right-7/8 w-1.5 h-1.5 bg-indigo-400 rounded-full animate-float-small opacity-30 shadow-sm" style={{animationDelay: '3.7s'}}></div>
        <div className="absolute top-6/8 left-2/8 w-1 h-1 bg-blue-300 rounded-full animate-float-small opacity-20 shadow-sm" style={{animationDelay: '0.6s'}}></div>
        <div className="absolute bottom-6/8 right-2/8 w-1.5 h-1.5 bg-indigo-300 rounded-full animate-float-small opacity-25 shadow-sm" style={{animationDelay: '2.1s'}}></div>
      </div>
      <div className="w-full max-w-xl relative z-10 px-4">
        {/* Logo y título */}
        <div className="text-center mb-3 sm:mb-4">
          <div className="flex flex-col sm:flex-row items-center justify-center space-y-2 sm:space-y-0 sm:space-x-3 mb-3 sm:mb-4">
            <div className="bg-blue-600 p-2 rounded-lg">
              <Truck className="h-6 w-6 sm:h-7 sm:w-7 text-white" />
            </div>
            <div className="text-center sm:text-left">
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">TransportApp</h1>
              <p className="text-gray-600 text-sm sm:text-base">Configuración inicial</p>
            </div>
          </div>
        </div>

        {/* Mensaje de bienvenida */}
        <div className="bg-green-50 border border-green-200 rounded-lg p-2 sm:p-3 mb-3 sm:mb-4">
          <div className="flex flex-col sm:flex-row items-center sm:items-start space-y-1 sm:space-y-0 sm:space-x-2">
            <div className="w-6 h-6 sm:w-7 sm:h-7 bg-green-500 rounded-full flex items-center justify-center flex-shrink-0">
              <svg className="w-3 h-3 sm:w-4 sm:h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"/>
              </svg>
            </div>
            <div className="text-center sm:text-left">
              <p className="text-green-800 font-semibold text-sm sm:text-base">
                ¡Bienvenido, {email || 'Usuario'}!
              </p>
              <p className="text-green-700 text-sm">Autenticación exitosa</p>
            </div>
          </div>
        </div>

        {/* Card principal */}
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl sm:text-2xl mb-2 sm:mb-3">
              Crear Organización
            </CardTitle>
            <CardDescription className="text-sm sm:text-base">
              Configura tu organización para comenzar a gestionar tus operaciones logísticas
            </CardDescription>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-3 sm:space-y-4">
              {/* Nombre de la organización */}
              <div>
                <label htmlFor="organizationName" className="block text-sm sm:text-base font-medium text-gray-700 mb-1 sm:mb-2">
                  Nombre de la Organización
                </label>
                <input
                  type="text"
                  id="organizationName"
                  value={organizationName}
                  onChange={(e) => setOrganizationName(e.target.value)}
                  placeholder="Ingresa el nombre de tu organización"
                  className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors text-sm sm:text-base ${
                    validationErrors.organizationName ? 'border-red-500' : 'border-gray-300'
                  }`}
                  required
                />
                {validationErrors.organizationName && (
                  <p className="mt-1 text-sm text-red-600 flex items-center gap-1">
                    <AlertCircle className="w-4 h-4" />
                    {validationErrors.organizationName}
                  </p>
                )}
              </div>

              {/* País */}
              <div>
                <label htmlFor="country" className="block text-sm sm:text-base font-medium text-gray-700 mb-1 sm:mb-2">
                  País de Operación Logística
                </label>
                <select
                  id="country"
                  value={selectedCountry}
                  onChange={(e) => setSelectedCountry(e.target.value)}
                  className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors appearance-none bg-white text-sm sm:text-base ${
                    validationErrors.country ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  {countries.map((country) => (
                    <option key={country.code} value={country.code}>
                      {country.name}
                    </option>
                  ))}
                </select>
                {validationErrors.country && (
                  <p className="mt-1 text-sm text-red-600 flex items-center gap-1">
                    <AlertCircle className="w-4 h-4" />
                    {validationErrors.country}
                  </p>
                )}
              </div>

              {/* Botón de envío */}
              <button
                type="submit"
                disabled={isSubmitting || !organizationName.trim()}
                className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold py-2 sm:py-3 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center gap-2 text-sm sm:text-base"
              >
                {isSubmitting ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                    <span>Creando...</span>
                  </>
                ) : (
                  <>
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M13 10V3L4 14h7v7l9-11h-7z"/>
                    </svg>
                    <span>Crear Organización</span>
                  </>
                )}
              </button>
            </form>

            {/* Mensaje de respuesta */}
            {submitResponse && (
              <div className={`mt-4 p-3 rounded-lg border ${
                submitResponse.success 
                  ? 'bg-green-50 border-green-200 text-green-800' 
                  : 'bg-red-50 border-red-200 text-red-800'
              }`}>
                <div className="flex items-center gap-2">
                  {submitResponse.success ? (
                    <CheckCircle className="w-5 h-5 text-green-600" />
                  ) : (
                    <AlertCircle className="w-5 h-5 text-red-600" />
                  )}
                  <span className="font-medium">
                    {submitResponse.success ? '¡Éxito!' : 'Error'}
                  </span>
                </div>
                <p className="mt-1 text-sm">
                  {submitResponse.message || submitResponse.error}
                </p>
              </div>
            )}
          </CardContent>

          <CardFooter>
            {/* Información adicional */}
            <div className="w-full bg-blue-50 border border-blue-200 rounded-lg p-2 sm:p-3">
              <div className="flex flex-col sm:flex-row items-start space-y-1 sm:space-y-0 sm:space-x-2">
                <svg className="w-4 h-4 sm:w-5 sm:h-5 text-blue-600 mt-1 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd"/>
                </svg>
                <div className="text-center sm:text-left">
                  <p className="text-blue-800 font-medium text-sm sm:text-base">Configuración Inicial</p>
                  <p className="text-blue-700 mt-1 text-sm">
                    Esta información se utilizará para configurar tu espacio de trabajo y optimizar las rutas según tu región
                  </p>
                </div>
              </div>
            </div>
          </CardFooter>
        </Card>

        {/* Footer */}
        <div className="mt-2 sm:mt-3 text-center text-sm text-gray-400">
          <p>© 2024 TransportApp. Todos los derechos reservados.</p>
        </div>
      </div>
    </div>
  )
}
