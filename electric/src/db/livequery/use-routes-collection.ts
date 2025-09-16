import { useLiveQuery } from '@tanstack/react-db'
import { useMemo } from 'react'
import { createRoutesCollection, extractRouteFromElectric, type RouteWithElectricId } from '../collections/create-routes-collection'
import { getMockDataByLanguage } from '../../demo/mockDataLoader'

// Función utilitaria para detectar modo demo
export const isDemoMode = (): boolean => {
  const urlParams = new URLSearchParams(window.location.search)
  const hashParams = new URLSearchParams(window.location.hash.slice(1))
  return urlParams.get('demoId') !== null || 
         window.location.pathname.includes('/demo') ||
         hashParams.get('demo') === 'true' ||
         hashParams.get('demoId') !== null ||
         urlParams.get('demold') !== null // Para el typo en la URL
}



// Hook personalizado que combina la collection con useLiveQuery
export const useRoutesCollection = (token: string, referenceId?: string) => {
  const collection = useMemo(() => createRoutesCollection(token, referenceId), [token, referenceId])
  
  const query = useLiveQuery((queryBuilder: any) => 
    queryBuilder.from({ route: collection })
  )
  
  // Transformar los datos para trabajar solo con el raw
  const routes = useMemo(() => {
    // query.data es un array de items con {id, raw}
    if (Array.isArray(query.data) && query.data.length > 0) {
      return query.data.map((item: any) => {
        // Usar extractRouteFromElectric para transformar cada item
        if (item && item.id && item.raw) {
          return extractRouteFromElectric(item)
        }
        
        // Fallback: si no tiene la estructura esperada, devolver el item completo
        return item
      })
    }
    
    return []
  }, [query.data])
  
  return {
    collection,
    query,
    routes,
    isLoading: query.isLoading,
    error: query.isError,
    // Métodos de la collection para mutaciones
    insert: collection.insert,
    update: collection.update,
    delete: collection.delete,
  }
}

// Hook simplificado que solo devuelve las rutas transformadas
export const useRoutes = (token: string, referenceId?: string): RouteWithElectricId[] => {
  // Detectar si estamos en modo demo
  const isDemo = useMemo(() => {
    const demo = isDemoMode()
    
    console.log('🔍 Demo detection:', {
      pathname: window.location.pathname,
      search: window.location.search,
      hash: window.location.hash,
      isDemo: demo
    })
    
    return demo
  }, [])
  
  // Si es demo, retornar datos mock sin consultar Electric SQL
  if (isDemo) {
    // Detectar idioma desde URL
    const urlParams = new URLSearchParams(window.location.search)
    const language = (urlParams.get('lang') || 'CL') as 'CL' | 'BR' | 'EU'
    
    // Obtener mock data basado en el idioma
    const mockRouteData = getMockDataByLanguage(language)
    
    // Usar el referenceId (que es el routeId de la ruta) en lugar del hardcodeado
    const actualRouteId = referenceId || '123' // Fallback si no hay referenceId
    const mockRoute = {
      ...mockRouteData,
      electricId: actualRouteId
    } as RouteWithElectricId
    
    console.log('🎭 Returning mock data with routeId:', actualRouteId, 'language:', language, mockRoute)
    return [mockRoute]
  }
  
  // Si no es demo, usar la query normal de Electric SQL
  const { routes, error } = useRoutesCollection(token, referenceId)
  
  if (error) {
    console.error('Error cargando rutas:', error)
    return []
  }
  
  return routes
}
