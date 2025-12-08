# Einar Framework Coding Rules

## Overview
This documentation outlines the coding standards and architectural patterns for the **Einar** framework used in this project. Strict adherence to these rules is required to ensure compatibility with the Einar CLI and the project's dependency injection system.

## Documentation Index

*   **[Controller Creation](controllers_creation.md)**: Detailed rules for creating HTTP controllers, including file placement, naming conventions, code structure, and configuration updates.
*   **[HTTP Server Infrastructure](http_server.md)**: Documentation for the core `httpserver.Server` dependency and its automatic generation rule.
*   **[PostgreSQL Repository Creation](postgresql.md)**: Rules for creating PostgreSQL repositories and their infrastructure dependencies.

## General Principles

1.  **Virtual CLI Behavior:** When creating components manually or via AI, you must follow the exact same conventions as the Einar CLI.
2.  **IoC Container:** All components must be registered with the `einar-ioc` container in their `init()` functions.
3.  **Configuration:** The `.einar.cli.json` file serves as the registry for project components and must be kept in sync with the codebase.

## 4. Main Imports Strategy

> [!IMPORTANT]
> **Imports in `main.go` must be added ONLY when the corresponding component is present in the project.**

The LLM **MUST** inspect `main.go` and ensure the following imports exist **only if** the related infrastructure or adapter is being used:

### Core (Always Required)
```go
_ "<module-name>/app/shared/configuration"
```

### HTTP Server & Controllers (Required if HTTP is used)
```go
_ "<module-name>/app/shared/infrastructure/httpserver"
_ "<module-name>/app/adapter/in/fuegoapi"
```

### PostgreSQL (Required ONLY if PostgreSQL is used)
```go
_ "<module-name>/app/shared/infrastructure/postgresql"
```

**Note:** Replace `<module-name>` with the actual module name defined in `go.mod`.
