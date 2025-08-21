import { useLiveQuery } from '@tanstack/react-db'
import { useParams } from '@tanstack/react-router'
import { createRoutesCollection } from './db/create-routes-collection'
import { useMemo } from 'react'

// Componente para la ruta raíz
export function HomeComponent() {
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || ''
  const routes = useMemo(() => createRoutesCollection(token), [token])
  const { data } = useLiveQuery((query) => query.from({ route: routes }))

  return (
    <div>
      <h1>Rutas Disponibles</h1>
      <p>Accede a una ruta específica usando: /driver/routes/&#123;routeId&#125;#access_token=&#123;token&#125;</p>
      <pre>{JSON.stringify(data, (_key, value) => 
        typeof value === 'bigint' ? value.toString() : value, 2)}</pre>
    </div>
  )
}

// Componente para rutas específicas del driver
export function RouteComponent() {
  // Obtener el routeId de los parámetros de la ruta usando TanStack Router
  const { routeId } = useParams({ from: '/driver/routes/$routeId' })
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || ''
  const routes = useMemo(() => createRoutesCollection(token, routeId), [token, routeId])
  const { data } = useLiveQuery((query) => query.from({ route: routes }))

  return (
    <div>
      <h2>Driver - Ruta ID: {routeId}</h2>
      <pre>{JSON.stringify(data, (_key, value) => 
        typeof value === 'bigint' ? value.toString() : value, 2)}</pre>
    </div>
  )
}