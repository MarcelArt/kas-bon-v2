# Frontend Build Plan — Overview

Authorization management app for the kas-bon-v2 RBAC backend. Dark mode is the default theme.

## Phases

| Phase | Focus | Key Deliverables |
|---|---|---|
| [Phase 1](./01-foundation.md) | Foundation | Type-safe env, server function API layer, auth context/state, auth pages (login/register), app shell layout, dark mode default, route guards |
| [Phase 2](./02-user-management.md) | User Management | User list, user detail, user CRUD, role assignment to users |
| [Phase 3](./03-domain-management.md) | Domain Management | Domain tree (nested organizations), domain CRUD, tree navigation |
| [Phase 4](./04-app-role-permission.md) | Apps, Roles & Permissions | App CRUD, role CRUD, permission CRUD, permission assignment to roles |
| [Phase 5](./05-permission-gating.md) | Permission Gating | Frontend permission system, menu/action visibility based on `users/{id}/permissions`, `all#fullAccess` superuser logic |

## Permission Model (Frontend)

The backend returns permission tuples as `string[][]` from `GET /v1/users/{id}/permissions`:
```
[[sub, app, dom, resource, action], ...]
```

Frontend permission checks:
- Parse tuples into a `Set<string>` of `"resource#action"` strings
- If the set contains `"all#fullAccess"`, the user has access to everything
- Otherwise, check specific permissions like `"users#read"`, `"roles#create"`, etc.

This gates:
- Sidebar menu item visibility
- Page access (route guards)
- Button/action visibility within pages (create, edit, delete buttons)

## Global Architecture Decisions

- **API Layer**: TanStack Start server functions (`createServerFn`) — all backend calls go through server functions in `src/lib/server/`. Server functions use `fetch` to call the backend API. Do NOT use Axios or raw `fetch` on the client.
- **Environment Variables**: Type-safe via `src/env.d.ts` TypeScript declarations for `NodeJS.ProcessEnv`. `API_URL` env var with fallback to `http://localhost:8080`. All env reads happen server-side only (inside server functions).
- **Auth Context**: Zustand auth store (`src/lib/stores/auth-store.ts`) holds user, tokens, selected org/app context, and permissions. Server functions receive auth context (`accessToken`, `domainId`, `appId`) as validated input from the client.
- **Token Refresh**: Handled client-side via `withAuthRetry` wrapper in TanStack Query hooks. Catches `AUTH_EXPIRED` errors, calls refresh server function, retries the original call.
- **Global State**: Zustand — auth store. Do NOT use React Context for global state.
- **Server State**: TanStack Query — all data fetching via `useQuery`, all mutations via `useMutation` with cache invalidation. Query key factories per resource in `src/lib/queries/`.
- **Forms**: TanStack Form + Zod — all forms use `useForm` from `@tanstack/react-form` with `@tanstack/zod-form-adapter`. Zod schemas for validation. Do NOT use manual `useState` for form fields.
- **Route loaders**: Use TanStack Router `loader` with `createServerFn` only for SSR-critical data (e.g., initial page load). For client-side data, use TanStack Query hooks.
- **Auth tokens**: Access/refresh tokens stored in Zustand (persisted to localStorage/sessionStorage). Client passes `accessToken` to server functions as input.
- **Path alias**: All imports via `@/`
- **Component library**: shadcn/ui (radix-lyra) — add components via `bunx shadcn add <name>`
- **Icons**: `@phosphor-icons/react` only
- **Font**: JetBrains Mono Variable (already configured)

## Server Function Architecture

All backend API calls use this pattern:

```
Component → TanStack Query hook → createServerFn (RPC) → Server handler → fetch to backend API
```

### Directory Structure

```
src/lib/
  env.ts                # Type-safe env vars (Zod + createServerOnlyFn)
  server/
    client.ts           # Base server API client (fetch wrapper, error handling)
    auth.ts             # Auth server functions (login, register, refresh, orgs, permissions)
    users.ts            # User server functions
    domains.ts          # Domain server functions
    apps.ts             # App server functions
    roles.ts            # Role server functions
    permissions.ts      # Permission server functions
  stores/
    auth-store.ts       # Zustand auth store
  queries/
    auth.ts             # Auth TanStack Query hooks
    users.ts            # User hooks
    ...
  schemas/
    auth.ts             # Zod schemas (shared between forms and server function validators)
    ...
```

### Server Function Pattern

```typescript
// src/lib/server/users.ts
import { createServerFn } from '@tanstack/react-start'
import { serverApi } from './client'

export const getUsersFn = createServerFn()
  .validator(z.object({
    accessToken: z.string(),
    domainId: z.number().optional(),
    appId: z.number().optional(),
    page: z.number().optional(),
    size: z.number().optional(),
  }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...params } = data
    return serverApi<PaginatedResponse<User>>('GET', '/users', { accessToken, domainId, appId }, undefined, params)
  })
```

### Query Hook Pattern

```typescript
// src/lib/queries/users.ts
export function useUsers(filters: object) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: userKeys.list(filters),
    queryFn: withAuthRetry(() => getUsersFn({ data: { ...auth, ...filters } })),
    enabled: !!auth.accessToken,
  })
}
```

### API Contracts

All API contracts are documented in `.opencode/API.md`. Consult it when building server functions or modifying endpoints. Key conventions:
- Base path: `/api/v1`
- Responses wrapped in `JSONResponse<T>` with `items`, `isSuccess`, `message` fields
- Auth: `Authorization: Bearer <token>` header
- Multi-tenancy: `X-App-Id` and `X-Domain-Id` headers
- Field naming: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` (PascalCase from gorm.Model)

## Post-Login Organization Selection Flow

After login, the app MUST call the organizations server function (which calls `GET /v1/users/{userId}/organizations`) before showing any authenticated content:

1. **0 organizations** → show error page: "You don't have access to any organization"
2. **1 organization** → auto-select it, set `domainId` and first app's `appId` in Zustand store, redirect to dashboard
3. **2+ organizations** → redirect to `/select-organization` page, user picks an org

The selected `domainId` and `appId` are stored in Zustand and passed to all subsequent server function calls as part of the auth context.
