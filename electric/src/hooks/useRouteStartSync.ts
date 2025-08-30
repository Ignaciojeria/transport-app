import { useRouteStartDomain } from '../db/gun/state/route-start-gun-state'
import type { RouteStart } from '../domain/route-start'

export const useRouteStartSync = (routeId: string) => {
  const { routeStart, loading } = useRouteStartDomain(routeId)
  
  return { routeStart, loading }
}
