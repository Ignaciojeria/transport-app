# HTTP Server Infrastructure

This document defines the **HTTP Server** infrastructure dependency required by all Einar HTTP controllers.

## 1. Component Overview

The `httpserver.Server` component provides the core HTTP server functionality for the application, including:
*   **HTTP Server:** Wraps the `fuego` server.
*   **Health Check:** Provides a `/health` endpoint using `hellofresh/health-go`.
*   **Graceful Shutdown:** Handles system signals (SIGINT, SIGTERM) to shut down the server gracefully.

**Why it is required:** All HTTP controllers depend on this component to register their routes.

**IoC Requirements:**
*   `ioc.Registry(New, configuration.NewConf)`: Registers the server constructor.
*   `ioc.RegistryAtEnd(startAtEnd, New)`: Starts the server after all dependencies are loaded.

## 2. File Structure

*   **Directory:** `/app/shared/infrastructure/httpserver/`
*   **File:** `server.go`

## 3. Critical Rule: Automatic Generation

> [!IMPORTANT]
> **If `/app/shared/infrastructure/httpserver/server.go` does not exist when the LLM is asked to generate a controller, the LLM MUST create this file automatically using the reference implementation below.**

This is a core infrastructure dependency.

## 4. Reference Implementation

The following code is the **official specification** for `server.go`.

```go
package httpserver

import (
    "<module-name>/app/shared/configuration"
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    ioc "github.com/Ignaciojeria/einar-ioc/v2"
    "github.com/go-fuego/fuego"
    "github.com/go-fuego/fuego/option"
    "github.com/hellofresh/health-go/v5"
)

func init() {
    ioc.Registry(New, configuration.NewConf)
    ioc.RegistryAtEnd(startAtEnd, New)
}

type Server struct {
    Manager *fuego.Server
    conf    configuration.Conf
}

func New(conf configuration.Conf) Server {
    s := fuego.NewServer(fuego.WithAddr(":" + conf.PORT))
    server := Server{
        Manager: s,
        conf:    conf,
    }
    ctx, cancel := context.WithCancel(context.Background())
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second*2)
        defer shutdownCancel()
        if err := s.Shutdown(shutdownCtx); err != nil {
            fmt.Println("Failed to shutdown:", err)
        }
        cancel()
    }()
    server.healthCheck()
    return server
}

func startAtEnd(e Server) error {
    return e.Manager.Run()
}

func (s Server) healthCheck() error {
    h, err := health.New(
        health.WithComponent(health.Component{
            Name:    s.conf.PROJECT_NAME,
            Version: s.conf.VERSION,
        }), health.WithSystemInfo())
    if err != nil {
        return err
    }
    fuego.GetStd(s.Manager,
        "/health",
        h.Handler().ServeHTTP,
        option.Summary("healthCheck"))
    return nil
}

func WrapPostStd(s Server, path string, f func(w http.ResponseWriter, r *http.Request)) {
    fuego.PostStd(s.Manager, path, f)
}
```

## 5. Mandatory Registration in main.go

> [!IMPORTANT]
> **When creating the HTTP Server, the LLM MUST ensure that `main.go` contains a blank import for the server package.**

**Required Import:**

```go
_ "<module-name>/app/shared/infrastructure/httpserver"
```

**Why?**
Without this import, the server will not be registered in the IoC system (einar-ioc), and the application will fail to start or register controllers.
This requirement is **mandatory** and **non-optional**.