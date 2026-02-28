# Guía de Uso de IOC

Esta documentación describe cómo usar el sistema de inyección de dependencias **github.com/Ignaciojeria/ioc** en el proyecto micartapro-backend, basado en los patrones establecidos en el código.

## Tabla de Contenidos

1. [Conceptos Básicos](#conceptos-básicos)
2. [Registro de Dependencias](#registro-de-dependencias)
3. [Patrones de Constructor](#patrones-de-constructor)
4. [Tipos Funcionales](#tipos-funcionales)
5. [Registro al Final del Ciclo](#registro-al-final-del-ciclo)
6. [Inicialización en main.go](#inicialización-en-maingo)
7. [Ejemplos Completos](#ejemplos-completos)
8. [Mejores Prácticas](#mejores-prácticas)

---

## Conceptos Básicos

### ¿Qué es ioc?

IOC es un contenedor de inyección de dependencias minimalista para Go que:
- **Infiere dependencias automáticamente** por tipos (parámetros → retornos)
- Registra constructores sin declarar dependencias manualmente
- Construye el grafo de dependencias al llamar `ioc.LoadDependencies()`

### Flujo de Trabajo

1. **Registro**: En la función `init()` de cada paquete, se registran los constructores
2. **Inicialización**: En `main.go`, se llama `ioc.LoadDependencies()` para construir todo el grafo
3. **Uso**: Las dependencias se inyectan automáticamente en los constructores

---

## Registro de Dependencias

### Sintaxis Básica

```go
func init() {
    ioc.Register(ConstructorFunction)  // Las dependencias se infieren por tipo
}
```

### Ejemplo 1: Registro Simple (Sin Dependencias)

```go
package configuration

import (
    ioc "github.com/Ignaciojeria/ioc"
)

func init() {
    ioc.Register(NewConf)
}

func NewConf() (Conf, error) {
    return Parse[Conf]()
}
```

### Ejemplo 2: Registro con Una Dependencia

```go
package supabasecli

import (
    "micartapro/app/shared/configuration"
    ioc "github.com/Ignaciojeria/ioc"
    supabase "github.com/supabase-community/supabase-go"
)

func init() {
    ioc.Register(NewSupabaseClient)  // configuration.Conf se infiere del parámetro
}

func NewSupabaseClient(conf configuration.Conf) (*supabase.Client, error) {
    return supabase.NewClient(
        conf.SUPABASE_PROJECT_URL,
        conf.SUPABASE_BACKEND_API_KEY,
        &supabase.ClientOptions{})
}
```

### Ejemplo 3: Registro con Múltiples Dependencias

```go
package storage

import (
    "context"
    "micartapro/app/events"
    "micartapro/app/shared/infrastructure/gcs"
    "micartapro/app/shared/infrastructure/observability"
    ioc "github.com/Ignaciojeria/ioc"
    "cloud.google.com/go/storage"
)

type GetLatestMenuById func(ctx context.Context, menuID string) (events.MenuCreateRequest, error)

func init() {
    ioc.Register(NewGetLatestMenuById)
}

func NewGetLatestMenuById(
    obs observability.Observability,
    gcs *storage.Client,
) GetLatestMenuById {
    return func(ctx context.Context, menuID string) (events.MenuCreateRequest, error) {
        // Implementación...
    }
}
```

### Ejemplo 4: Registro de Endpoints HTTP (Fuego API)

```go
package fuegoapi

import (
    "micartapro/app/adapter/out/storage"
    "micartapro/app/adapter/out/supabaserepo"
    "micartapro/app/shared/infrastructure/httpserver"
    "micartapro/app/shared/infrastructure/observability"
    ioc "github.com/Ignaciojeria/ioc"
    "github.com/go-fuego/fuego"
)

func init() {
    ioc.Register(searchMenuById)  // Dependencias inferidas por tipo
}

func searchMenuById(
    s httpserver.Server,
    obs observability.Observability,
    getMenuSlugBySlug supabaserepo.GetMenuSlugBySlug,
    getLatestMenuById storage.GetLatestMenuById,
) {
    fuego.Get(s.Manager, "/menu/slug/{slug}", /* ... */)
}
```

---

## Patrones de Constructor

### Patrón 1: Constructor que Retorna un Tipo

```go
type Server struct {
    Manager *fuego.Server
    conf    configuration.Conf
}

func New(conf configuration.Conf, requestLoggerMiddleware RequestLoggerMiddleware) Server {
    // Construcción del objeto
    return Server{
        Manager: fuego.NewServer(/* ... */),
        conf:    conf,
    }
}
```

### Patrón 2: Constructor que Retorna un Tipo Funcional

```go
type GetLatestMenuById func(ctx context.Context, menuID string) (events.MenuCreateRequest, error)

func NewGetLatestMenuById(
    obs observability.Observability,
    gcs *storage.Client,
) GetLatestMenuById {
    return func(ctx context.Context, menuID string) (events.MenuCreateRequest, error) {
        // Implementación de la función
    }
}
```

### Patrón 3: Constructor que Retorna un Error

```go
func NewSupabaseClient(conf configuration.Conf) (*supabase.Client, error) {
    return supabase.NewClient(
        conf.SUPABASE_PROJECT_URL,
        conf.SUPABASE_BACKEND_API_KEY,
        &supabase.ClientOptions{})
}
```

**Nota**: Si un constructor retorna un error, el IOC lo manejará durante `LoadDependencies()`.

---

## Tipos Funcionales

Un patrón común en el proyecto es usar **tipos funcionales** para definir interfaces de manera más flexible:

```go
// Definir el tipo funcional
type OnMenuInteractionRequest func(ctx context.Context, input events.MenuInteractionRequest) (string, error)

// Constructor que retorna el tipo funcional
func NewOnMenuInteractionRequest(
    obs observability.Observability,
    menuInteractionAgent agents.MenuInteractionAgent,
    publisherManager eventprocessing.PublisherManager,
    getLatestMenuById storage.GetLatestMenuById,
) OnMenuInteractionRequest {
    return func(ctx context.Context, input events.MenuInteractionRequest) (string, error) {
        // Lógica del caso de uso
        menu, err := getLatestMenuById(ctx, input.MenuID)
        // ...
    }
}
```

**Ventajas**:
- No requiere interfaces explícitas
- Más flexible y fácil de testear
- Permite inyección de dependencias sin acoplamiento

---

## Registro al Final del Ciclo

### `ioc.RegistryAtEnd`

Usa `ioc.RegistryAtEnd` cuando necesitas que una función se ejecute **después** de que todas las dependencias estén construidas. Útil para:
- Inicializar servidores HTTP
- Conectar a bases de datos
- Iniciar workers o subscribers

### Ejemplo: Iniciar Servidor HTTP

```go
package httpserver

func init() {
    ioc.Register(New)
    ioc.RegisterAtEnd(startAtEnd)
}

func New(conf configuration.Conf, requestLoggerMiddleware RequestLoggerMiddleware) Server {
    // Construir el servidor
    return Server{/* ... */}
}

// Esta función se ejecuta DESPUÉS de que todas las dependencias estén listas
func startAtEnd(e Server, obs observability.Observability) error {
    obs.Logger.Info(
        "http server started",
        "port", e.conf.PORT,
        "service", e.conf.PROJECT_NAME,
    )
    return e.Manager.Run() // Inicia el servidor
}
```

**Sintaxis**:
```go
ioc.RegisterAtEnd(
    FunctionToExecute,      // Función que se ejecutará al final
    Dependency1,           // Dependencia 1
    Dependency2,           // Dependencia 2
)
```

---

## Inicialización en main.go

El punto de entrada debe importar todos los paquetes que contienen registros de dependencias y luego llamar a `ioc.LoadDependencies()`:

```go
package main

import (
    _ "micartapro/app/shared/configuration"        // Importar para ejecutar init()
    _ "micartapro/app/adapter/in/fuegoapi"        // Importar para ejecutar init()
    _ "micartapro/app/shared/infrastructure/httpserver"
    _ "micartapro/app/adapter/out/storage"
    _ "micartapro/app/usecase/menu"
    // ... más imports

    ioc "github.com/Ignaciojeria/ioc"
    "log"
)

func main() {
    // Cargar todas las dependencias registradas
    if err := ioc.LoadDependencies(); err != nil {
        log.Fatal(err)
    }
    // El servidor se iniciará automáticamente gracias a RegistryAtEnd
}
```

**Importante**: 
- Usa `_` para importar paquetes solo por sus efectos secundarios (ejecutar `init()`)
- `ioc.LoadDependencies()` construye todo el grafo de dependencias
- Si hay errores en la construcción, se retornan aquí

---

## Ejemplos Completos

### Ejemplo 1: Caso de Uso Completo

```go
package menu

import (
    "context"
    "micartapro/app/adapter/out/storage"
    "micartapro/app/events"
    "micartapro/app/shared/infrastructure/observability"
    ioc "github.com/Ignaciojeria/ioc"
)

// 1. Definir el tipo funcional
type OnMenuCreateRequest func(ctx context.Context, input events.MenuCreateRequest) error

// 2. Registrar en init()
func init() {
    ioc.Register(NewOnMenuCreateRequest)
}

// 3. Implementar el constructor
func NewOnMenuCreateRequest(
    observability observability.Observability,
    saveMenu storage.SaveMenu,
) OnMenuCreateRequest {
    return func(ctx context.Context, input events.MenuCreateRequest) error {
        observability.Logger.Info("on_menu_create_request", "input", input)
        spanCtx, span := observability.Tracer.Start(ctx, "on_menu_create_request")
        defer span.End()
        
        return saveMenu(spanCtx, input)
    }
}
```

### Ejemplo 2: Repositorio con Cliente Externo

```go
package supabaserepo

import (
    "context"
    "micartapro/app/shared/infrastructure/supabasecli"
    ioc "github.com/Ignaciojeria/ioc"
    "github.com/supabase-community/supabase-go"
)

type GetMenuSlugBySlug func(ctx context.Context, slug string) (MenuSlugInfo, error)

func init() {
    ioc.Register(NewGetMenuSlugBySlug)
}

func NewGetMenuSlugBySlug(supabase *supabase.Client) GetMenuSlugBySlug {
    return func(ctx context.Context, slug string) (MenuSlugInfo, error) {
        // Usar el cliente de Supabase para hacer queries
        // ...
    }
}
```

### Ejemplo 3: Endpoint HTTP Completo

```go
package fuegoapi

import (
    "micartapro/app/adapter/out/storage"
    "micartapro/app/adapter/out/supabaserepo"
    "micartapro/app/events"
    "micartapro/app/shared/infrastructure/httpserver"
    "micartapro/app/shared/infrastructure/observability"
    ioc "github.com/Ignaciojeria/ioc"
    "github.com/go-fuego/fuego"
)

func init() {
    ioc.Register(searchMenuById)
    )
}

func searchMenuById(
    s httpserver.Server,
    obs observability.Observability,
    getMenuSlugBySlug supabaserepo.GetMenuSlugBySlug,
    getLatestMenuById storage.GetLatestMenuById,
) {
    fuego.Get(s.Manager, "/menu/slug/{slug}",
        func(c fuego.ContextNoBody) (events.MenuCreateRequest, error) {
            // Usar las dependencias inyectadas
            slug := c.PathParam("slug")
            slugInfo, err := getMenuSlugBySlug(c.Context(), slug)
            // ...
        },
    )
}
```

---

## Mejores Prácticas

### 1. Nomenclatura de Constructores

- **Siempre** usa el prefijo `New` para constructores: `NewService`, `NewRepository`, etc.
- El nombre debe reflejar qué construye: `NewSupabaseClient`, `NewGetLatestMenuById`

### 2. Inferencia de Dependencias

- Solo registra el constructor: `ioc.Register(NewService)`
- Las dependencias se infieren automáticamente por los tipos de los parámetros
- El orden de los parámetros del constructor define qué se inyecta

```go
// ✅ Correcto - las dependencias se infieren por tipo
ioc.Register(NewService)

func NewService(dep1 Dependency1, dep2 Dependency2) Service {
    // ...
}
```

### 3. Tipos Funcionales vs Interfaces

- **Usa tipos funcionales** para casos de uso, repositorios y servicios
- **Usa interfaces** solo cuando necesites múltiples implementaciones del mismo contrato

### 4. Manejo de Errores

- Si un constructor puede fallar, retorna `(T, error)`
- El IOC propagará el error durante `LoadDependencies()`
- No uses `panic` en constructores

### 5. Registro en init()

- **Siempre** registra dependencias en `func init()`
- No registres en funciones normales
- Importa los paquetes en `main.go` con `_` para ejecutar los `init()`

### 6. Dependencias Circulares

- **Evita** dependencias circulares
- Si es necesario, usa `ioc.RegistryAtEnd` para romper el ciclo
- Considera refactorizar si encuentras dependencias circulares

### 7. Testing

Para testear componentes que usan IOC:

```go
func TestMyService(t *testing.T) {
    // Crear mocks manualmente
    mockDep := &MockDependency{}
    
    // Llamar al constructor directamente
    service := NewMyService(mockDep)
    
    // Testear...
}
```

### 8. Estructura de Paquetes

```
app/
├── adapter/
│   ├── in/          # Entrada (HTTP, gRPC, etc.)
│   └── out/         # Salida (Repositorios, Clientes externos)
├── domain/          # Entidades de dominio
├── usecase/         # Casos de uso
└── shared/          # Infraestructura compartida
    ├── configuration/
    ├── infrastructure/
    └── ...
```

---

## Resumen de Comandos

| Comando | Descripción | Ejemplo |
|---------|-------------|---------|
| `ioc.Register(ctor)` | Registra un constructor (dependencias inferidas por tipo) | `ioc.Register(NewService)` |
| `ioc.RegisterAtEnd(ctor)` | Registra una función que se ejecuta al final | `ioc.RegisterAtEnd(startServer)` |
| `ioc.LoadDependencies()` | Construye todo el grafo de dependencias | Llamar en `main()` |

---

## Referencias

- Repositorio: [github.com/Ignaciojeria/einar-ioc](https://github.com/Ignaciojeria/einar-ioc)
- Versión usada: `v2`
- Ejemplos en el proyecto: Ver `micartapro-backend/app/`

---

## Preguntas Frecuentes

### ¿Cómo resuelve el IOC las dependencias?

El IOC analiza los tipos de los parámetros de los constructores y busca constructores registrados que retornen esos tipos. Si encuentra una coincidencia, inyecta esa dependencia.

### ¿Qué pasa si falta una dependencia?

`ioc.LoadDependencies()` retornará un error indicando qué dependencia no se pudo resolver.

### ¿Puedo tener múltiples constructores para el mismo tipo?

No directamente. El IOC usa el tipo de retorno para resolver dependencias, así que solo puede haber un constructor por tipo. Si necesitas múltiples implementaciones, usa interfaces.

### ¿Cómo testeo componentes con IOC?

Llama directamente al constructor pasando mocks manualmente. No necesitas usar el IOC en los tests.

---

**Última actualización**: Basado en los patrones de `micartapro-backend` - 2024
