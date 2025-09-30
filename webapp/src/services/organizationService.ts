/**
 * Servicio para manejar operaciones relacionadas con organizaciones
 */

export interface CreateOrganizationRequest {
  email: string
  organizationName: string
  country: string
}

export interface CreateOrganizationResponse {
  success: boolean
  message?: string
  organizationId?: string
  error?: string
}

/**
 * Crea una nueva organización
 * @param data - Datos de la organización a crear
 * @returns Respuesta del servidor
 */
export const createOrganization = async (data: CreateOrganizationRequest): Promise<CreateOrganizationResponse> => {
  try {
    const response = await fetch('https://einar-main-f0820bc.d2.zuplo.dev/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: data.email,
        organizationName: data.organizationName,
        country: data.country
      })
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      return {
        success: false,
        error: errorData.message || `Error del servidor: ${response.status}`
      }
    }

    const responseData = await response.json()
    
    return {
      success: true,
      message: 'Organización creada exitosamente',
      organizationId: responseData.organizationId || responseData.id
    }
  } catch (error) {
    console.error('Error al crear organización:', error)
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Error desconocido al crear la organización'
    }
  }
}

/**
 * Valida los datos de entrada para crear una organización
 * @param data - Datos a validar
 * @returns Objeto con errores de validación
 */
export const validateOrganizationData = (data: Partial<CreateOrganizationRequest>) => {
  const errors: Record<string, string> = {}

  if (!data.email || !data.email.trim()) {
    errors.email = 'El email es requerido'
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
    errors.email = 'El email no tiene un formato válido'
  }

  if (!data.organizationName || !data.organizationName.trim()) {
    errors.organizationName = 'El nombre de la organización es requerido'
  } else if (data.organizationName.trim().length < 2) {
    errors.organizationName = 'El nombre debe tener al menos 2 caracteres'
  }

  if (!data.country || !data.country.trim()) {
    errors.country = 'El país es requerido'
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  }
}
