/**
 * Utilidades para manejo de menús en auth-ui
 */

// Detectar si estamos en desarrollo local
const isLocalDev = typeof window !== 'undefined' && 
  (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')

// URL del backend - se detecta automáticamente según el entorno
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 
  (isLocalDev ? 'http://localhost:8082' : 'https://micartapro-backend-27303662337.us-central1.run.app')

const SUPABASE_URL = 'https://rbpdhapfcljecofrscnj.supabase.co'
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJicGRoYXBmY2xqZWNvZnJzY25qIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ5NjY3NDMsImV4cCI6MjA4MDU0Mjc0M30.Ba-W2KHJS8U6OYVAjU98Y7JDn87gYPuhFvg_0vhcFfI'

/**
 * Verifica si el usuario tiene un menú activo
 * @param userId - ID del usuario
 * @param accessToken - Token de autenticación
 * @returns ID del menú si existe, null si no tiene menú
 */
export async function getUserActiveMenu(
  userId: string,
  accessToken: string
): Promise<string | null> {
  try {
    const { createClient } = await import('@supabase/supabase-js')
    const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
      global: {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      },
      auth: {
        autoRefreshToken: false,
        persistSession: false,
        detectSessionInUrl: false
      }
    })
    
    const { data, error } = await supabase
      .from('menus')
      .select('id')
      .eq('user_id', userId)
      .limit(1)
      .maybeSingle()
    
    if (error) {
      // Si es PGRST116, significa que no hay registros (no es un error real)
      if (error.code === 'PGRST116') {
        return null
      }
      console.error('Error verificando menú del usuario:', error)
      return null
    }
    
    return data?.id || null
  } catch (error) {
    console.error('Error en getUserActiveMenu:', error)
    return null
  }
}

/**
 * Crea un menú inicial para el usuario
 * @param menuId - ID del menú (UUID)
 * @param accessToken - Token de autenticación
 * @returns true si se creó correctamente, false si hubo error
 */
export async function createInitialMenu(
  menuId: string,
  accessToken: string
): Promise<boolean> {
  try {
    const payload = {
      id: menuId,
      coverImage: "https://storage.googleapis.com/micartapro-menus/core/micartaprov3.webp",
      businessInfo: {
        businessName: "cadorago",
        whatsapp: "+56957857558",
        businessHours: [] as string[]
      },
      menu: [] as any[],
      deliveryOptions: [
        {
          type: "DELIVERY" as const,
          requireTime: false
        },
        {
          type: "PICKUP" as const,
          requireTime: true,
          timeRequestType: "WINDOW" as const,
          timeWindows: [
            {
              start: "09:00",
              end: "23:59"
            }
          ]
        }
      ]
    }
    
    console.log('🔄 Creando menú inicial:', { menuId, payload })
    
    const response = await fetch(`${API_BASE_URL}/menu`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(payload)
    })
    
    if (!response.ok) {
      const errorText = await response.text()
      console.error('❌ Error creando menú inicial:', response.status, errorText)
      return false
    }
    
    const result = await response.json()
    console.log('✅ Menú inicial creado exitosamente:', result)
    return true
  } catch (error) {
    console.error('❌ Error en createInitialMenu:', error)
    return false
  }
}
