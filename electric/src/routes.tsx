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
      hash="access_token=test"
    />
  ),
})

// Ruta para rutas específicas del driver
const routeByIdRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/driver/routes/$routeId',
  component: RouteComponent,
})

// Crear el árbol de rutas
const routeTree = rootRoute.addChildren([indexRoute, routeByIdRoute])

// Crear el router
export const router = createRouter({ routeTree })

// Declarar tipos para TypeScript
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
