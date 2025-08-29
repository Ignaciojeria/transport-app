import { createCollection } from '@tanstack/react-db'
import { electricCollectionOptions } from '@tanstack/electric-db-collection'
import { z } from 'zod'
import type { Route } from '../domain/route'

// Schema para las rutas
const RouteSchema = z.object({
  id: z.string(),
  raw: z.any(), // Mantener z.any() por ahora para evitar conflictos de tipos
  reference_id: z.string().optional(),
  created_at: z.date().optional(),
  updated_at: z.date().optional(),
})

// Tipos derivados de los schemas
type RouteType = z.infer<typeof RouteSchema>

export type { RouteType }

// Factory para crear la colección inyectando el token
export const createRoutesCollection = (token: string, referenceId?: string) =>
  createCollection(
    electricCollectionOptions({
      id: 'routes',
      shapeOptions: {
        url: (() => {
          const base = `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=routes&columns=id,raw`
          return referenceId ? `${base}&where=reference_id='${referenceId}'` : base
        })(),
        headers: {
          'X-Access-Token': `Bearer ${token}`,
        },
        parser: {
          timestamptz: (iso: string) => new Date(iso),
          timestamp: (iso: string) => new Date(iso),
        },
      },
      getKey: (r) => r.id,
      schema: RouteSchema,
      
      async onInsert({ transaction }) {
        const { modified: newRoute } = transaction.mutations[0]
        console.log('Inserting route:', newRoute)
        return { txid: [Date.now()] }
      },
      
      async onUpdate({ transaction }) {
        const { original, modified } = transaction.mutations[0]
        console.log('Updating route:', { original, modified })
        return { txid: [Date.now()] }
      },
      
      async onDelete({ transaction }) {
        const { original } = transaction.mutations[0]
        console.log('Deleting route:', original)
        return { txid: [Date.now()] }
      },
    })
  )

