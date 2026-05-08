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

- **Server functions**: Use `createServerFn` from TanStack Start for all API calls (keeps tokens server-side)
- **Data fetching**: TanStack Query (`@tanstack/react-query`) for all client-side data fetching, caching, and mutations. Use `@tanstack/react-router-ssr-query` (already installed) to integrate Query with Router. Define query keys and query/mutation functions in server files.
- **Forms**: TanStack Form (`@tanstack/react-form`) for all form state management. Use Zod (`zod`) for form validation schemas. Do NOT use react-hook-form or formik.
- **Validation**: Zod (`zod`) for all runtime validation — form inputs, API response parsing, and schema definitions.
- **Route loaders**: Use TanStack Query's router integration for SSR data loading (query hydration from server to client)
- **Auth state**: Access/refresh tokens stored in cookies (httpOnly preferred); user + permissions in context
- **API client**: Shared fetch wrapper in `src/lib/api.ts` that attaches `Authorization`, `X-App-Id`, `X-Domain-Id` headers
- **Path alias**: All imports via `@/`
- **Component library**: shadcn/ui (radix-lyra) — add components via `bunx shadcn add <name>`
- **Icons**: `@phosphor-icons/react` only
- **Font**: JetBrains Mono Variable (already configured)

### Libraries to Install (Phase 1)

```bash
cd web
bun add @tanstack/react-query @tanstack/react-form @tanstack/zod-form-adapter zod
```
