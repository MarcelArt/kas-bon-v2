# AGENTS.md

## Project Overview

Go backend API + React frontend for an RBAC-based authentication and authorization system.
Backend: Fiber v3, GORM, PostgreSQL, Casbin v3, JWT, Cobra CLI.
Frontend (`web/`): React 19, TanStack Start/Router/Query/Form, Zod, shadcn/ui, Tailwind CSS v4, Vitest.

## Build & Run Commands

### Backend (Go)

```bash
make dev             # Run with air hot-reload + swagger generation
make go              # Run without hot-reload (go run main.go serve)
make swag            # Generate swagger docs only
make migrate         # Run database migrations
make migrate-force   # Run migrations with table drop
go build -o bin/app  # Build binary
```

### Go Testing (no tests exist yet)

```bash
go test ./...                              # All tests
go test ./internal/v1/handlers/...         # Single package
go test -run TestFunctionName ./path/...   # Single test
go test -v ./...                           # Verbose output
```

### Frontend (web/)

```bash
cd web
bun run dev                                      # Dev server on port 3000
bun run build                                    # Production build
bun run test                                     # Run all tests (vitest run)
bun run test -- path/to/test.test.ts             # Run a single test file
bun run test -- -t "test name"                   # Run tests matching a name
bun run lint                                     # ESLint
bun run format                                   # Prettier
bun run typecheck                                # TypeScript check (tsc --noEmit)
```

## Architecture & Directory Structure

```
cmd/                        # Cobra CLI commands (root, serve, migrate, policy)
internal/
  configs/                  # DB, env, Casbin setup (global singletons: DB, Env)
  common/                   # Shared utilities (jwt, json_response, authz, permission)
  enums/                    # Constants (app, role, permission, time)
  v1/
    handlers/               # HTTP handlers (one file per resource)
    repositories/           # Data access layer (interface + implementation)
    routes/                 # Route registration (one file per resource)
    middlewares/             # Auth and Casbin middleware
    models/                 # GORM models and input DTOs
    usecases/               # Business logic (builder pattern)
pkg/arrays/                 # Shared array utilities
web/                        # React frontend (TanStack Start + shadcn/ui)
```

## Backend Code Style

### Layered Architecture

Each resource follows: **Route -> Handler -> Usecase/Repository -> Model**

- **Routes** (`internal/v1/routes/`): Wire handlers to endpoints with middleware.
- **Handlers** (`internal/v1/handlers/`): Parse request, call repo/usecase, return JSON.
- **Repositories** (`internal/v1/repositories/`): Interface + struct. Raw SQL for pagination, GORM for CRUD.
- **Usecases** (`internal/v1/usecases/`): Complex business logic using builder pattern.
- **Models** (`internal/v1/models/`): GORM models + input DTOs.

### Handlers

Struct with dependency injection, constructor returns pointer. Methods take `c fiber.Ctx`, return `error`. Add Swagger annotations above each method.

```go
type UserHandler struct {
    repo repositories.IUserRepo
    e    *casbin.Enforcer
}
func NewUserHandler(repo repositories.IUserRepo, e *casbin.Enforcer) *UserHandler {
    return &UserHandler{repo: repo, e: e}
}
```

### Repositories

Define interface first (`IXxxRepo`), then struct implementation. Constructor: `NewXxxRepo(db *gorm.DB)`. Pagination uses `morkid/paginate` with raw SQL stored in struct field.

### Models

DB models embed `gorm.Model`. Input DTOs embed `Input` (has ID, CreatedAt, UpdatedAt with `json:"-"`). Define a `const xxxTableName` and `TableName()` method on input types.

### Usecases (Builder Pattern)

Initialize with `InitXxxUsecase()`, set fields via chained setters, call `Execute()`:

```go
res, err := usecases.InitGenerateTokenPairUsecase().
    SetCtx(c).SetUser(user).SetIsRemember(isRemember).Execute()
```

### Error Handling

Always use `common.NewJSONResponse(dataOrErr, "message")` — it detects errors automatically. Use `common.StatusCodeFromError(err)` for dynamic status codes (maps `gorm.ErrRecordNotFound` to 404).

### Transactions

```go
tx := configs.DB.Begin()
defer tx.Rollback()
// ... operations using tx ...
tx.Commit()
```

### Imports

Group in order (separated by blank lines): standard library, project-internal packages, third-party packages.

### Naming Conventions

- Files: `snake_case` with dot separator (e.g., `user.handler.go`, `user.route.go`, `user.model.go`)
- Structs/Types: `PascalCase` (e.g., `UserHandler`)
- Interfaces: `IXxxRepo` with `I` prefix (e.g., `IUserRepo`)
- Variables: `camelCase` (e.g., `dRepo`, `pageQuery`)
- Constants: `PascalCase` (e.g., `AppName`, `ResourceAll`)
- Functions: `SetupXxxRoutes()`, `NewXxxHandler()`, `NewXxxRepo()`

### Authorization

Permissions use `"resource#action"` format (e.g., `"users#read"`). Middleware chain: `Authn()` for JWT, then `authz.HasPermission("resource#action")` for Casbin RBAC.

### Global Singletons

- `configs.DB` — GORM database instance
- `configs.Env` — Environment configuration

## Frontend (web/)

A `@frontend` subagent is configured in `.opencode/agents/frontend.md`. Use `@frontend` for any React/TypeScript work in `web/`. It enforces find-docs/ctx7 verification, shadcn/ui conventions, and TanStack Start patterns.

Run before finishing frontend work: `cd web && bun run typecheck && bun run lint`

## API Reference

A full API reference document is at `.opencode/API.md` — generated from `docs/swagger.json`. It covers all endpoints, request/response schemas, auth requirements, pagination, and error codes. Consult it when building API clients or modifying endpoints.

## General Guidelines

- Run `make swag` after modifying Swagger annotations in handler files.
- Environment variables loaded from `.env` via godotenv (never commit `.env`).
- Do not add comments unless asked.
