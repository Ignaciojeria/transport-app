import { useState } from 'react'

interface CreateOrganizationProps {
  userEmail: string
}

export const CreateOrganizationSimple = ({ userEmail }: CreateOrganizationProps) => {
  const [organizationName, setOrganizationName] = useState('')
  const [country, setCountry] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validateForm = () => {
    const newErrors: Record<string, string> = {}
    
    if (!organizationName.trim()) {
      newErrors.organizationName = 'El nombre de la organización es requerido'
    } else if (organizationName.length > 64) {
      newErrors.organizationName = 'El nombre no puede exceder 64 caracteres'
    }
    
    if (!country) {
      newErrors.country = 'Debes seleccionar un país'
    }
    
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }
    
    setIsLoading(true)
    
    try {
      // Aquí iría la lógica para crear la organización
      console.log('Creando organización:', { organizationName, country, userEmail })
      
      // Simular API call
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      alert('Organización creada exitosamente!')
    } catch (error) {
      console.error('Error creando organización:', error)
      alert('Error al crear la organización')
    } finally {
      setIsLoading(false)
    }
  }

  const handleBack = () => {
    // Si estamos en modo demo, limpiar y volver
    if (localStorage.getItem('demo_mode')) {
      localStorage.removeItem('demo_mode');
    }
    window.location.href = '/'
  }

  const countries = [
    { value: 'AR', label: '🇦🇷 Argentina' },
    { value: 'CL', label: '🇨🇱 Chile' },
    { value: 'CO', label: '🇨🇴 Colombia' },
    { value: 'MX', label: '🇲🇽 México' },
    { value: 'PE', label: '🇵🇪 Perú' },
    { value: 'UY', label: '🇺🇾 Uruguay' }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-900 flex items-center justify-center p-4">
      {/* Header compacto */}
      <div className="absolute top-6 left-6 flex items-center gap-3 text-white">
        <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center">
          <span className="text-xl">🚛</span>
        </div>
        <div>
          <h1 className="font-bold text-lg">Transport APP</h1>
          <p className="text-white/70 text-sm">Gestión logística</p>
        </div>
      </div>

      {/* Usuario info */}
      <div className="absolute top-6 right-6 text-right text-white/80 text-sm">
        <p>{userEmail}</p>
        <p className="text-white/60">Configuración inicial</p>
      </div>

      {/* Formulario principal */}
      <div className="w-full max-w-md">
        <div className="bg-white/10 backdrop-blur-xl rounded-2xl border border-white/20 shadow-2xl p-8">
          {/* Header del formulario */}
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-gradient-to-br from-emerald-500 to-blue-600 rounded-xl flex items-center justify-center mx-auto mb-4 shadow-lg">
              <span className="text-white text-2xl">🏢</span>
            </div>
            <h2 className="text-2xl font-bold text-white mb-2">
              Crea tu organización
            </h2>
            <p className="text-white/70">
              Configura tu espacio de trabajo para gestionar tu logística
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Campo nombre organización */}
            <div>
              <label htmlFor="organizationName" className="block text-white font-medium mb-2">
                Nombre de tu organización
              </label>
              <input
                id="organizationName"
                type="text"
                value={organizationName}
                onChange={(e) => setOrganizationName(e.target.value)}
                className="w-full bg-white/5 border border-white/30 rounded-lg px-4 py-3 text-white placeholder-white/40 focus:outline-none focus:ring-2 focus:ring-emerald-400 focus:border-transparent transition-all"
                maxLength={64}
                placeholder="Ej: Transportes ABC"
                required
              />
              <div className="flex justify-between text-xs mt-1">
                <span className="text-white/50">
                  Visible para todos los miembros
                </span>
                <span className="text-white/50">
                  {organizationName.length}/64
                </span>
              </div>
              {errors.organizationName && (
                <p className="text-red-400 text-sm mt-1">
                  {errors.organizationName}
                </p>
              )}
            </div>

            {/* Campo país */}
            <div>
              <label htmlFor="country" className="block text-white font-medium mb-2">
                País de operación logística
              </label>
              <select
                id="country"
                value={country}
                onChange={(e) => setCountry(e.target.value)}
                className="w-full bg-white/5 border border-white/30 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-emerald-400 focus:border-transparent transition-all"
                required
              >
                <option value="" disabled className="bg-gray-800 text-gray-300">
                  Selecciona un país
                </option>
                {countries.map((countryOption) => (
                  <option 
                    key={countryOption.value} 
                    value={countryOption.value}
                    className="bg-gray-800 text-white"
                  >
                    {countryOption.label}
                  </option>
                ))}
              </select>
              <p className="text-white/50 text-xs mt-1">
                País principal de operaciones
              </p>
              {errors.country && (
                <p className="text-red-400 text-sm mt-1">
                  {errors.country}
                </p>
              )}
            </div>

            {/* Botones */}
            <div className="space-y-3 pt-4">
              <button
                type="submit"
                disabled={isLoading || !organizationName.trim() || !country}
                className="w-full bg-gradient-to-r from-emerald-600 to-blue-600 hover:from-emerald-700 hover:to-blue-700 disabled:from-gray-600 disabled:to-gray-700 text-white font-semibold py-3 px-6 rounded-lg transition-all duration-200 transform hover:scale-[1.02] disabled:hover:scale-100 shadow-lg disabled:cursor-not-allowed"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center gap-2">
                    <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Creando organización...
                  </div>
                ) : (
                  <div className="flex items-center justify-center gap-2">
                    <span>✨</span>
                    Crear Organización
                  </div>
                )}
              </button>

              <button
                type="button"
                onClick={handleBack}
                className="w-full text-white/70 hover:text-white hover:bg-white/5 py-2 px-4 rounded-lg transition-all text-sm"
              >
                ← Volver
              </button>
            </div>
          </form>
        </div>

        {/* Elementos decorativos */}
        <div className="absolute top-20 left-10 text-2xl animate-bounce opacity-60">🌟</div>
        <div className="absolute top-40 right-20 text-xl animate-pulse opacity-40">✨</div>
        <div className="absolute bottom-20 left-20 text-lg animate-bounce opacity-50">🚀</div>
        <div className="absolute bottom-40 right-10 text-xl animate-pulse opacity-30">💫</div>
      </div>
    </div>
  )
}
