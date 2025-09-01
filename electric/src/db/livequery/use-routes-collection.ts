import { useLiveQuery } from '@tanstack/react-db'
import { useMemo } from 'react'
import { createRoutesCollection, extractRouteFromElectric, type RouteWithElectricId } from '../collections/create-routes-collection'

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
  const { routes, error } = useRoutesCollection(token, referenceId)
  
  if (error) {
    console.error('Error cargando rutas:', error)
    return []
  }
  
  return routes
}
