# AGENTS.md

## Project Overview

Go backend API + React frontend for an RBAC-based authentication and authorization system.
Backend: Fiber v3, GORM, PostgreSQL, Casbin v3, JWT, Cobra CLI.
Frontend (`web/`): React 19, TanStack Start/Router, shadcn/ui, Tailwind CSS v4, Vitest.

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

### Frontend Libraries

These libraries MUST be used as described. Verify with find-docs/ctx7 before using any API:

- **Axios** — HTTP client. All API calls go through a shared Axios instance (`src/lib/api.ts`) that attaches `Authorization`, `X-App-Id`, `X-Domain-Id` headers and intercepts 401s for token refresh. Do NOT use raw `fetch`.
- **Zustand** — Global state management. Used for auth store (`useAuthStore`) holding user, tokens, selected org/app context, permissions. Do NOT use React Context for global state.
- **TanStack Query** — Server state management. All data fetching (lists, detail views) uses `useQuery`. All mutations (create, update, delete) use `useMutation` with cache invalidation. Define query keys factory per resource in `src/lib/queries/`.
- **TanStack Form** — Form state management. All forms (login, register, CRUD dialogs) use `useForm` from `@tanstack/react-form`. Do NOT use manual `useState` for form fields.
- **Zod** — Schema validation. All form validation schemas defined with Zod. Used with TanStack Form via `@tanstack/zod-form-adapter`. API response types validated with Zod schemas where needed.

### Post-Login Organization Selection Flow

After login, the app MUST call `GET /v1/users/{id}/organizations` to get the user's organizations:
1. If user has **0 organizations** → show error or redirect to a "no access" page
2. If user has **1 organization** → auto-select it, set `X-Domain-Id` in Axios defaults, proceed to dashboard
3. If user has **multiple organizations** → show organization picker page/route before dashboard

The selected organization's ID becomes `X-Domain-Id` and a default app's ID becomes `X-App-Id` for all subsequent API calls. This selection is stored in the Zustand auth store.

## API Reference

A full API reference document is at `.opencode/API.md` — generated from `docs/swagger.json`. It covers all endpoints, request/response schemas, auth requirements, pagination, and error codes. Consult it when building API clients or modifying endpoints.

**JSON field naming:** All DB models embed `gorm.Model` (no `json` tags), so Go serializes as PascalCase: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`. Other fields have explicit lowercase `json` tags (e.g., `username`, `domainId`). TypeScript types must match: `ID` not `id`, `CreatedAt` not `createdAt`.

## General Guidelines

- Run `make swag` after modifying Swagger annotations in handler files.
- Environment variables loaded from `.env` via godotenv (never commit `.env`).
- Do not add comments unless asked.
