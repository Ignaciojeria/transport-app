import { createRoute, createRootRoute, createRouter, Outlet, Navigate } from '@tanstack/react-router'
import { RouteComponent } from './App'

// Ruta raíz
const rootRoute = createRootRoute({
  component: () => (
    <>

      <Outlet />
    </>
  ),
})

// Ruta para la página principal
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: () => (
    <Navigate
      to="/driver/routes/$routeId"
      params={{ routeId: '123' }}
      hash="access_token=test&demo=true"
    />
  ),
})

// Ruta para rutas específicas del driver
const routeByIdRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/driver/routes/$routeId',
  component: RouteComponent,
})

// Ruta para la demo - redirige a la ruta del driver con demo=true
const demoRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/demo',
  component: () => {
    // Leer routeId de los query parameters
    const urlParams = new URLSearchParams(window.location.search)
    const routeId = urlParams.get('routeId') || '123' // Default a '123' si no se proporciona
    
    return (
      <Navigate
        to="/driver/routes/$routeId"
        params={{ routeId }}
        hash="access_token=test&demo=true"
      />
    )
  },
})

// Crear el árbol de rutas
const routeTree = rootRoute.addChildren([indexRoute, routeByIdRoute, demoRoute])

// Crear el router
export const router = createRouter({ routeTree })

// Declarar tipos para TypeScript
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
