import { useState, useEffect } from 'react'
import { getRouteStart, routeStartKey } from '../db/route-start-gun-state'
import { deliveriesData } from '../db/gun'
import type { RouteStart } from '../domain/route-start'

export const useRouteStartSync = (routeId: string) => {
  const [routeStart, setRouteStart] = useState<RouteStart | null>(null)
  const [loading, setLoading] = useState(true)
  
  useEffect(() => {
    if (!routeId) return
    
    setLoading(true)
    
    // Obtener estado inicial
    getRouteStart(routeId).then(setRouteStart).finally(() => setLoading(false))
    
    // Escuchar cambios
    const key = routeStartKey(routeId)
    const unsubscribe = deliveriesData.get(key).on((data) => {
      if (data && typeof data === 'string') {
        try {
          const parsed = JSON.parse(data)
          const { timestamp, deviceId, ...routeStartData } = parsed
          setRouteStart(routeStartData)
        } catch (error) {
          console.error('Error parseando datos de inicio de ruta:', error)
        }
      } else {
        setRouteStart(null)
      }
    })
    
    return () => {
      if (unsubscribe && typeof unsubscribe.off === 'function') {
        unsubscribe.off()
      }
    }
  }, [routeId])
  
  return { routeStart, loading }
}
