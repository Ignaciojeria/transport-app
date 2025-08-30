import { useLiveQuery } from '@tanstack/react-db'
import { useMemo } from 'react'
import { createRoutesCollection, extractRouteFromElectric, type RouteWithElectricId } from '../collections/create-routes-collection'

// Hook personalizado que combina la collection con useLiveQuery
export const useRoutesCollection = (token: string, referenceId?: string) => {
  console.log('🚀🚀🚀 useRoutesCollection EJECUTÁNDOSE 🚀🚀🚀')
  console.log('🚀🚀🚀 Token:', token ? '✅' : '❌')
  console.log('🚀🚀🚀 ReferenceId:', referenceId)
  
  const collection = useMemo(() => createRoutesCollection(token, referenceId), [token, referenceId])
  
  console.log('🔍 useRoutesCollection - Token:', token ? '✅' : '❌')
  console.log('🔍 useRoutesCollection - ReferenceId:', referenceId)
  console.log('🔍 useRoutesCollection - Collection:', collection)
  
  const query = useLiveQuery((queryBuilder: any) => 
    queryBuilder.from({ route: collection })
  )
  
  console.log('🔍 useRoutesCollection - Query data:', query.data)
  console.log('🔍 useRoutesCollection - Query data type:', typeof query.data)
  console.log('🔍 useRoutesCollection - Query data isArray:', Array.isArray(query.data))
  if (query.data && query.data.length > 0) {
    console.log('🔍 useRoutesCollection - First item keys:', Object.keys(query.data[0]))
    console.log('🔍 useRoutesCollection - First item:', query.data[0])
  }
  console.log('🔍 useRoutesCollection - Query isLoading:', query.isLoading)
  console.log('🔍 useRoutesCollection - Query isError:', query.isError)
  
  // Transformar los datos para trabajar solo con el raw
  const routes = useMemo(() => {
    console.log('🔍 useRoutesCollection - Transformando datos...')
    console.log('🔍 useRoutesCollection - Query data completo:', query.data)
    
    // query.data es un array de items con {id, raw}
    if (Array.isArray(query.data) && query.data.length > 0) {
      console.log('🔍 useRoutesCollection - Procesando items...')
      
      const result = query.data.map((item: any) => {
        console.log('🔍 useRoutesCollection - Procesando item:', item)
        
        // Usar extractRouteFromElectric para transformar cada item
        if (item && item.id && item.raw) {
          return extractRouteFromElectric(item)
        }
        
        // Fallback: si no tiene la estructura esperada, devolver el item completo
        console.log('🔍 useRoutesCollection - Item no tiene estructura esperada, usando fallback')
        return item
      })
      
      console.log('🔍 useRoutesCollection - Rutas transformadas:', result)
      return result
    }
    
    console.log('🔍 useRoutesCollection - No hay datos para transformar')
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
  console.log('🎯🎯🎯 useRoutes EJECUTÁNDOSE 🎯🎯🎯')
  console.log('🎯🎯🎯 Token:', token ? '✅' : '❌')
  console.log('🎯🎯🎯 ReferenceId:', referenceId)
  
  const { routes, isLoading, error } = useRoutesCollection(token, referenceId)
  
  console.log('🎯🎯🎯 Routes result:', routes)
  console.log('🎯🎯🎯 IsLoading:', isLoading)
  console.log('🎯🎯🎯 Error:', error)
  
  if (error) {
    console.error('Error cargando rutas:', error)
    return []
  }
  
  return routes
}
