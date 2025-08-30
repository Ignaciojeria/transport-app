import { createCollection } from '@tanstack/react-db'
import { electricCollectionOptions } from '@tanstack/electric-db-collection'
import type { Route } from '../../domain/route'

// Tipo para la estructura que devuelve Electric
type ElectricRouteData = {
  id: string
  raw: Route & { 
    // Incluir el id de Electric en el objeto raw para trabajar solo con él
    electricId: string 
  }
  reference_id?: string
  created_at?: Date
  updated_at?: Date
}

export type { ElectricRouteData }

// Helper para extraer solo el objeto raw con el id incluido
export type RouteWithElectricId = Route & { electricId: string }

// Función helper para transformar ElectricRouteData a RouteWithElectricId
export const extractRouteFromElectric = (electricData: ElectricRouteData): RouteWithElectricId => {
  return {
    ...electricData.raw,
    electricId: electricData.id
  }
}

// Factory para crear la colección inyectando el token
export const createRoutesCollection = (token: string, referenceId?: string) => {
  const url = (() => {
    const base = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=routes&columns=id,raw`
    return referenceId ? `${base}&where=reference_id='${referenceId}'` : base
  })()
  
  return createCollection(
    electricCollectionOptions({
      id: 'routes',
      shapeOptions: {
        url,
        headers: {
          'X-Access-Token': `Bearer ${token}`,
        },
        parser: {
          timestamptz: (iso: string) => new Date(iso),
          timestamp: (iso: string) => new Date(iso),
        },
      },
      getKey: (r: ElectricRouteData) => r.id,
      // No necesitas schema si usas tipos TypeScript
      
      async onInsert() {
        return { txid: [Date.now()] }
      },
      
      async onUpdate() {
        return { txid: [Date.now()] }
      },
      
      async onDelete() {
        return { txid: [Date.now()] }
      },
    })
  )
}

