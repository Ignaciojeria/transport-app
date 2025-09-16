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
  component: () => {
    // Preservar query parameters existentes como objeto
    const urlParams = new URLSearchParams(window.location.search)
    const searchParams: Record<string, string> = {}
    
    urlParams.forEach((value, key) => {
      searchParams[key] = value
    })
    
    console.log('🔄 IndexRoute: Redirecting with search params:', {
      originalSearch: window.location.search,
      searchParams,
      href: window.location.href
    })
    
    return (
      <Navigate
        to="/driver/routes/$routeId"
        params={{ routeId: '123' }}
        search={searchParams}
        hash="access_token=test&demo=true"
      />
    )
  },
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
    
    // Preservar todos los query parameters existentes como objeto
    const searchParams: Record<string, string> = {}
    urlParams.forEach((value, key) => {
      if (key !== 'routeId') { // No incluir routeId en search ya que va en params
        searchParams[key] = value
      }
    })
    
    console.log('🔄 DemoRoute: Redirecting with search params:', {
      originalSearch: window.location.search,
      searchParams,
      routeId,
      href: window.location.href
    })
    
    return (
      <Navigate
        to="/driver/routes/$routeId"
        params={{ routeId }}
        search={searchParams}
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
