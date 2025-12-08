# Controller Creation Rules

This document outlines the strict rules for creating controllers in the Einar framework. These rules apply to both CLI-generated code and AI-generated code.

## 1. File Structure & Location

**Directory:** `/app/adapter/in/fuegoapi/`

**File Naming:** `snake_case.go` (e.g., `create_employees.go`, `search_employees.go`)

## 2. Code Structure

Every controller file must follow this exact pattern:

```go
package fuegoapi

import (
    "<module-name>/app/shared/infrastructure/httpserver"

    ioc "github.com/Ignaciojeria/einar-ioc/v2"
    "github.com/go-fuego/fuego"
    "github.com/go-fuego/fuego/option"
)

func init() {
    // Registration with IoC container
    ioc.Registry(functionName, httpserver.New)
}

func functionName(s httpserver.Server) {
    // Route definition using s.Manager
    fuego.Post(s.Manager, "/route-path",
        func(c fuego.ContextNoBody) (any, error) {
            return "unimplemented", nil
        },
        option.Summary("functionName"),
    )
}
```

### Key Elements:
*   **Package:** Must be `package fuegoapi`.
*   **Imports:** Must include `<module-name>/app/shared/infrastructure/httpserver`, `einar-ioc`, and `fuego` packages.
*   **Init Function:** Must register the controller using `ioc.Registry(functionName, httpserver.New)`.
*   **Controller Function:**
    *   Name must be `camelCase` (e.g., `createEmployees`).
    *   Must accept `s httpserver.Server`.
    *   Must define the route using `fuego.Get`, `fuego.Post`, etc., on `s.Manager`.

## 3. Configuration Update (.einar.cli.json)

Every new controller **MUST** be registered in the `.einar.cli.json` file in the project root.

**Location:** `.einar.cli.json`

**Update Rule:** Add a new object to the `components` array.

```json
{
    "kind": "post-controller", // or "get-controller", "put-controller", etc.
    "name": "component-name"   // kebab-case (e.g., "create-employees")
}
```

## 4. Naming Conventions Summary

| Entity | Convention | Example |
| :--- | :--- | :--- |
| **File Name** | snake_case | `create_employees.go` |
| **Component Name** | kebab-case | `create-employees` |
| **Go Function** | camelCase | `createEmployees` |
| **Route Path** | kebab-case | `/employees` |

## 6. Infrastructure Dependency: HTTP Server

All controllers depend on the `httpserver.Server` infrastructure component.
This component is responsible for providing the underlying `fuego.Server` instance used to register routes.
Without this server, controllers cannot be initialized or registered in the IoC container.

Controllers must import:

```go
import "<module-name>/app/shared/infrastructure/httpserver"
```

and `main.go` must include the blank import:

```go
_ "<module-name>/app/shared/infrastructure/httpserver"
```

Without this blank import, the server will not be registered in the IoC system (einar-ioc), and all controllers will fail.

For detailed server documentation, see [/llm/http_server.md](/llm/http_server.md).


## 7. Mandatory Registration in main.go

> [!IMPORTANT]
> **Every time a new controller is generated, the LLM MUST ensure that `main.go` contains a blank import for the controller package.**

**Required Import:**

```go
_ "<module-name>/app/adapter/in/fuegoapi"
```

**Why?**
Without this import, the controller will not be registered in the IoC container (einar-ioc), and the application will not load the controller.
This requirement is **mandatory** and **non-optional**.

## 5. Checklist

- [ ] File created in `app/adapter/in/fuegoapi/`.
- [ ] Filename is `snake_case.go`.
- [ ] Package is `fuegoapi`.
- [ ] `init()` registers with `ioc.Registry`.
- [ ] Controller function accepts `httpserver.Server`.
- [ ] Route defined on `s.Manager`.
- [ ] `.einar.cli.json` updated with correct `kind` and `name`.
- [ ] `main.go` includes blank import `_ "<module-name>/app/adapter/in/fuegoapi"`.
