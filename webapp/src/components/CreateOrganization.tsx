import { useState } from 'react'

interface CreateOrganizationProps {
  userEmail: string
}

export const CreateOrganization = ({ userEmail }: CreateOrganizationProps) => {
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
    // Lógica para volver (podría ser logout o redirect)
    console.log('Volver')
  }

  const countries = [
    { value: 'AR', label: 'Argentina' },
    { value: 'CL', label: 'Chile' },
    { value: 'CO', label: 'Colombia' },
    { value: 'MX', label: 'México' },
    { value: 'PE', label: 'Perú' },
    { value: 'UY', label: 'Uruguay' }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500">
      {/* Header con DaisyUI */}
      <div className="navbar bg-base-100/10 backdrop-blur-lg border-b border-white/10">
        <div className="navbar-start">
          <div className="flex items-center gap-2 text-white font-semibold text-lg">
            <span className="text-2xl">🚛</span>
            Transport APP
          </div>
        </div>
        <div className="navbar-end">
          <div className="flex items-center gap-4 text-white">
            <span className="text-sm">{userEmail}</span>
            <button className="btn btn-ghost btn-sm text-white">Ayuda</button>
            <button className="btn btn-ghost btn-sm text-white">Docs</button>
            <div className="avatar placeholder">
              <div className="bg-white/20 text-white rounded-full w-8">
                <span className="text-sm font-semibold">
                  {userEmail.charAt(0).toUpperCase()}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs con DaisyUI */}
      <div className="bg-white/5 px-6">
        <div className="tabs">
          <a className="tab tab-bordered tab-active text-white border-white">
            Organizaciones
          </a>
          <a className="tab tab-bordered text-white/70 hover:text-white">
            Uso
          </a>
        </div>
      </div>

      {/* Contenido principal */}
      <div className="flex justify-center items-center min-h-[calc(100vh-120px)] p-6">
        <div className="card w-full max-w-md bg-base-100 shadow-2xl relative">
          {/* Botón cerrar */}
          <button 
            className="btn btn-ghost btn-sm absolute top-4 right-4 text-base-content/60 hover:text-base-content"
            onClick={handleBack}
          >
            ✕
          </button>

          <div className="card-body">
            {/* Logo hexagonal */}
            <div className="flex justify-center mb-6">
              <div className="relative">
                <div className="w-20 h-20 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg">
                  <span className="text-white text-3xl font-bold">E</span>
                </div>
                <div className="absolute -top-1 -left-1 w-22 h-22 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl opacity-30 -z-10"></div>
              </div>
            </div>

            <h1 className="text-2xl font-bold text-center mb-8 text-base-content">
              Crea tu organización
            </h1>

            <form onSubmit={handleSubmit} className="space-y-6" noValidate>
              {/* Campo nombre organización */}
              <div className="form-control">
                <label htmlFor="organizationName" className="label">
                  <span className="label-text font-medium">Nombre de tu organización</span>
                  <span className="label-text-alt text-base-content/60">
                    {organizationName.length}/64
                  </span>
                </label>
                <input
                  id="organizationName"
                  name="organizationName"
                  type="text"
                  value={organizationName}
                  onChange={(e) => setOrganizationName(e.target.value)}
                  className={`input input-bordered w-full ${errors.organizationName ? 'input-error' : ''}`}
                  maxLength={64}
                  required
                  placeholder="Ingresa el nombre de tu organización"
                  aria-describedby="organizationName-help organizationName-error"
                  aria-invalid={!!errors.organizationName}
                />
                <label className="label">
                  <span id="organizationName-help" className="label-text-alt text-base-content/70">
                    Este nombre será visible para todos los miembros de tu organización
                  </span>
                </label>
                {errors.organizationName && (
                  <label className="label">
                    <span id="organizationName-error" className="label-text-alt text-error" role="alert">
                      {errors.organizationName}
                    </span>
                  </label>
                )}
              </div>

              {/* Campo país */}
              <div className="form-control">
                <label htmlFor="country" className="label">
                  <span className="label-text font-medium">País de operación logística</span>
                </label>
                <select
                  id="country"
                  name="country"
                  value={country}
                  onChange={(e) => setCountry(e.target.value)}
                  className={`select select-bordered w-full ${errors.country ? 'select-error' : ''}`}
                  required
                  aria-describedby="country-help country-error"
                  aria-invalid={!!errors.country}
                >
                  <option value="" disabled>
                    Selecciona un país
                  </option>
                  {countries.map((countryOption) => (
                    <option key={countryOption.value} value={countryOption.value}>
                      {countryOption.label}
                    </option>
                  ))}
                </select>
                <label className="label">
                  <span id="country-help" className="label-text-alt text-base-content/70">
                    Este será el país principal de operaciones de tu organización
                  </span>
                </label>
                {errors.country && (
                  <label className="label">
                    <span id="country-error" className="label-text-alt text-error" role="alert">
                      {errors.country}
                    </span>
                  </label>
                )}
              </div>

              {/* Botón submit */}
              <button 
                type="submit" 
                className={`btn btn-primary w-full ${isLoading ? 'loading' : ''}`}
                disabled={isLoading || !organizationName.trim() || !country}
                aria-describedby="submit-status"
              >
                {isLoading ? 'Creando...' : 'Crear Organización'}
              </button>

              {isLoading && (
                <div id="submit-status" className="sr-only" aria-live="polite">
                  Creando organización, por favor espera...
                </div>
              )}
            </form>

            {/* Botón volver */}
            <div className="text-center mt-6">
              <button 
                className="btn btn-ghost btn-sm text-base-content/70 hover:text-base-content"
                onClick={handleBack}
              >
                Volver
              </button>
            </div>
          </div>

          {/* Elementos decorativos */}
          <div className="absolute -top-5 -right-5 text-2xl animate-bounce">🌟</div>
          <div className="absolute -bottom-3 -left-3 text-xl animate-pulse">✨</div>
        </div>
      </div>
    </div>
  )
}
