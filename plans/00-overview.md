# Frontend Build Plan — Overview

Authorization management app for the kas-bon-v2 RBAC backend. Dark mode is the default theme.

## Phases

| Phase | Focus | Key Deliverables |
|---|---|---|
| [Phase 1](./01-foundation.md) | Foundation | API client, auth context/state, auth pages (login/register), app shell layout, dark mode default, route guards |
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

- **HTTP Client**: Axios — shared instance in `src/lib/api.ts` that attaches `Authorization`, `X-App-Id`, `X-Domain-Id` headers via interceptors. Intercepts 401 for token refresh. Do NOT use raw `fetch`.
- **Global State**: Zustand — auth store (`src/lib/stores/auth-store.ts`) holding user, tokens, selected org/app context, and permissions. Do NOT use React Context for global state.
- **Server State**: TanStack Query — all data fetching via `useQuery`, all mutations via `useMutation` with cache invalidation. Query key factories per resource in `src/lib/queries/`.
- **Forms**: TanStack Form + Zod — all forms use `useForm` from `@tanstack/react-form` with `@tanstack/zod-form-adapter`. Zod schemas for validation. Do NOT use manual `useState` for form fields.
- **Route loaders**: Use TanStack Router `loader` with `createServerFn` only for SSR-critical data (e.g., initial page load). For client-side data, use TanStack Query hooks.
- **Auth tokens**: Access/refresh tokens stored in Zustand (persisted to localStorage/sessionStorage). Axios interceptors read from store.
- **Path alias**: All imports via `@/`
- **Component library**: shadcn/ui (radix-lyra) — add components via `bunx shadcn add <name>`
- **Icons**: `@phosphor-icons/react` only
- **Font**: JetBrains Mono Variable (already configured)

## Post-Login Organization Selection Flow

After login, the app MUST call `GET /v1/users/{userId}/organizations` before showing any authenticated content:

1. **0 organizations** → show error page: "You don't have access to any organization"
2. **1 organization** → auto-select it, set `domainId` and first app's `appId` in Zustand store, redirect to dashboard
3. **2+ organizations** → redirect to `/select-organization` page, user picks an org

The selected `domainId` → `X-Domain-Id` header, and `appId` → `X-App-Id` header for ALL subsequent API calls. These are set in the Axios request interceptor by reading from the Zustand auth store.
