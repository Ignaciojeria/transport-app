import { useState } from 'react'
import CreateOrganization from './components/CreateOrganization'
import { type CreateOrganizationResponse } from './services/organizationService'
import './App.css'

function App() {
  const [organizationCreated, setOrganizationCreated] = useState(false)
  const [organizationData, setOrganizationData] = useState<{name: string; country: string} | null>(null)
  // const [error, setError] = useState<string | null>(null)

  // Token simulado para pruebas
  const mockToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzdWFyaW9AZWplbXBsby5jb20iLCJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlVzdWFyaW8gZGUgUHJ1ZWJhIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'

  const handleCreateOrganizationSuccess = (response: CreateOrganizationResponse) => {
    console.log('Organización creada exitosamente:', response)
    setOrganizationData({
      name: 'Organización Creada', // En una implementación real, obtendrías esto de la respuesta
      country: 'CL'
    })
    setOrganizationCreated(true)
    // setError(null)
  }

  const handleCreateOrganizationError = (error: string) => {
    console.error('Error al crear organización:', error)
    // setError(error)
  }

  if (organizationCreated && organizationData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex flex-col items-center justify-center p-8">
        <div className="bg-white rounded-xl shadow-xl p-8 max-w-md w-full text-center">
          <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd"/>
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">¡Organización Creada!</h2>
          <p className="text-gray-600 mb-4">
            La organización <strong>"{organizationData.name}"</strong> ha sido creada exitosamente en <strong>{organizationData.country}</strong>.
          </p>
          <button 
            onClick={() => {
              setOrganizationCreated(false)
              setOrganizationData(null)
            }}
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition-colors duration-200"
          >
            Crear Otra Organización
          </button>
        </div>
      </div>
    )
  }

  return (
    <CreateOrganization 
      token={mockToken}
      onSuccess={handleCreateOrganizationSuccess}
      onError={handleCreateOrganizationError}
    />
  )
}

export default App
