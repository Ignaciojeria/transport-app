import { createRoute, createRootRoute, createRouter, Outlet } from '@tanstack/react-router'
import { HomeComponent, RouteComponent } from './App'

// Ruta raíz
const rootRoute = createRootRoute({
  component: () => (
    <>
      <div className="p-2 flex gap-2">
        <a href="/" className="font-bold">Home</a>
        <a href="/driver/routes/123#access_token=test" className="font-bold">Ejemplo Ruta</a>
      </div>
      <hr />
      <Outlet />
    </>
  ),
})

// Ruta para la página principal
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomeComponent,
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
