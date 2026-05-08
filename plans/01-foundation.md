# Phase 1 — Foundation

Authentication, app shell, dark mode, shared API layer, data fetching, and form handling.

## 0. Install Dependencies

```bash
cd web
bun add @tanstack/react-query @tanstack/react-form @tanstack/zod-form-adapter zod
bunx shadcn add sidebar separator tooltip card input label checkbox sonner
```

## 1. Dark Mode Default

**File:** `src/routes/__root.tsx`

- Add `className="dark"` to the `<html>` element in `RootDocument`
- This activates the `.dark` CSS custom properties already defined in `styles.css`
- Update page title from "TanStack Start Starter" to app name

## 2. API Client

**File:** `src/lib/api.ts`

Shared fetch wrapper for all backend calls:

```
- Base URL: /api/v1 (from env or hardcoded)
- Auto-attaches headers:
  - Authorization: Bearer <accessToken> (from cookie or context)
  - X-App-Id (from app context)
  - X-Domain-Id (from domain context)
- Parses JSONResponse envelope:
  - If isSuccess === true, return items
  - If isSuccess === false, throw with message
- Handles 401 by triggering token refresh or redirect to /login
```

**File:** `src/lib/api.types.ts`

TypeScript types mirroring all backend models:
- `User`, `UserInput`, `LoginInput`, `LoginResponse`
- `Role`, `RoleInput`
- `Permission`, `PermissionInput`
- `Domain`, `DomainInput`
- `App`, `AppInput`
- `JSONResponse<T>`, `PaginatedResponse<T>`
- `AccessControlEval`, `TokenEndpointRequest`

**File:** `src/lib/api.schemas.ts`

Zod schemas for all API types (used by both TanStack Form validation and API response parsing):
- `userInputSchema`, `loginInputSchema`, `registerInputSchema`
- `roleInputSchema`, `permissionInputSchema`, `domainInputSchema`, `appInputSchema`
- `loginResponseSchema`

## 3. TanStack Query Setup

**File:** `src/lib/query-client.tsx`

- Create a `QueryClient` instance with sensible defaults (staleTime, retry)
- Wrap the app in `QueryClientProvider` (in `__root.tsx` or a dedicated provider component)
- The project already has `@tanstack/react-router-ssr-query` installed — use it to integrate Query with TanStack Router for SSR hydration

## 4. Auth Server Functions

**File:** `src/lib/auth.server.ts`

Using `createServerFn` from TanStack Start:

- `loginFn(body: LoginInput)` — calls `POST /v1/users/login`, stores tokens in cookies, returns user
- `registerFn(body: UserInput)` — calls `POST /v1/users`, returns new user ID
- `refreshFn()` — calls `POST /v1/users/refresh` with refresh token from cookie
- `logoutFn()` — clears cookies
- `getCurrentUserFn()` — reads user from cookie/session (or re-validates token)
- `getUserPermissionsFn(userId, appId, domainId)` — calls `GET /v1/users/{id}/permissions`, returns permission tuples

## 5. Auth Hooks (TanStack Query)

**File:** `src/lib/auth.query.ts`

TanStack Query hooks for auth state:

- `useCurrentUser()` — `useQuery` that calls `getCurrentUserFn`, key: `["auth", "me"]`
- `useUserPermissions(userId)` — `useQuery` that calls `getUserPermissionsFn`, key: `["auth", "permissions", userId]`
- `useLoginMutation()` — `useMutation` wrapping `loginFn`, on success: set cookies + invalidate auth queries + redirect
- `useRegisterMutation()` — `useMutation` wrapping `registerFn`, on success: redirect to login
- `useLogoutMutation()` — `useMutation` wrapping `logoutFn`, on success: clear queries + redirect
- `useRefreshMutation()` — `useMutation` wrapping `refreshFn` (used internally for 401 handling)

## 6. Auth Context & State

**File:** `src/lib/auth.tsx`

React context/provider providing:
- `user: User | null`
- `permissions: Set<string>` — parsed from permission tuples
- `isSuperUser: boolean` — true if permissions has `"all#fullAccess"`
- `hasPermission(resource, action): boolean` — checks permission set
- `isLoading: boolean`

This context reads from TanStack Query's cache (useCurrentUser, useUserPermissions) and derives the permission set.

Permission check logic:
```
hasPermission(resource, action):
  if isSuperUser → return true
  return permissions.has(`${resource}#${action}`)
```

## 7. Route Structure

```
src/routes/
  __root.tsx          # Shell with dark class, QueryClientProvider
  _auth.tsx           # Layout for unauthenticated pages (no sidebar)
  _auth/
    login.tsx         # /login
    register.tsx      # /register
  _authenticated.tsx  # Layout with sidebar, requires auth
  _authenticated/
    dashboard.tsx     # / (redirect here after login)
    users/
      index.tsx       # /users
      $userId.tsx     # /users/:id
    domains/
      index.tsx       # /domains
      $domainId.tsx   # /domains/:id
    apps/
      index.tsx       # /apps
      $appId.tsx      # /apps/:id
    roles/
      index.tsx       # /roles
      $roleId.tsx     # /roles/:id
    permissions/
      index.tsx       # /permissions
      $permissionId.tsx # /permissions/:id
```

Use TanStack Router layout routes (`_auth` and `_authenticated`) for shared layout wrappers.

## 8. App Shell Layout

**File:** `src/components/layout/app-shell.tsx`

Structure for authenticated pages:
```
<AppShell>
  <Sidebar>          # shadcn sidebar component
    <SidebarHeader>  # App name / logo
    <SidebarContent>
      <SidebarGroup> # Navigation items (gated by permissions)
        <SidebarItem to="/dashboard" icon={House} />
        <SidebarItem to="/users" icon={Users} permission="users#read" />
        <SidebarItem to="/domains" icon={Folders} permission="domains#read" />
        <SidebarItem to="/apps" icon={AppWindow} permission="apps#read" />
        <SidebarItem to="/roles" icon={Shield} permission="roles#read" />
        <SidebarItem to="/permissions" icon={Key} permission="permissions#read" />
      </SidebarGroup>
    </SidebarContent>
    <SidebarFooter>  # User info + logout
  </Sidebar>
  <main>             # Page content via <Outlet />
</AppShell>
```

## 9. Auth Pages (TanStack Form + Zod)

### Login Page — `src/routes/_auth/login.tsx`

- Use `useForm` from `@tanstack/react-form` with `@tanstack/zod-form-adapter`
- Zod schema `loginInputSchema` validates: username (required), password (required), isRemember (boolean)
- On valid submit, call `useLoginMutation().mutate(body)`
- On success: redirect to `/dashboard`
- On error: show server error via sonner toast
- Link to register page

### Register Page — `src/routes/_auth/register.tsx`

- Use `useForm` from `@tanstack/react-form` with `@tanstack/zod-form-adapter`
- Zod schema `registerInputSchema` validates: username (required, min 3), email (required, email format), password (required, min 6), confirmPassword (must match password)
- On valid submit, call `useRegisterMutation().mutate(body)` (omit confirmPassword before sending)
- On success: redirect to `/login` with success toast
- On error: show validation / server errors

## 10. Route Guards

**File:** `src/lib/auth-guard.ts`

In `_authenticated.tsx` layout route, use `beforeLoad` to:
1. Check if user is authenticated (token exists and valid)
2. If not, redirect to `/login`
3. Load user permissions if not already loaded
4. Inject `permissions` and `hasPermission` into route context

```typescript
// _authenticated.tsx
export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: "/login" })
    }
  },
})
```

## 11. Files to Create/Modify

| File | Action |
|---|---|
| `src/routes/__root.tsx` | Modify: add `className="dark"` to `<html>`, update title, wrap in QueryClientProvider |
| `src/lib/api.ts` | Create: API client with auth headers |
| `src/lib/api.types.ts` | Create: TypeScript types for all backend models |
| `src/lib/api.schemas.ts` | Create: Zod schemas for all API types |
| `src/lib/query-client.tsx` | Create: QueryClient instance + provider |
| `src/lib/auth.server.ts` | Create: server functions for auth API calls |
| `src/lib/auth.query.ts` | Create: TanStack Query hooks for auth |
| `src/lib/auth.tsx` | Create: auth context, provider, permission helpers |
| `src/lib/auth-guard.ts` | Create: route guard utilities |
| `src/routes/_auth.tsx` | Create: unauthenticated layout |
| `src/routes/_auth/login.tsx` | Create: login page (TanStack Form + Zod) |
| `src/routes/_auth/register.tsx` | Create: register page (TanStack Form + Zod) |
| `src/routes/_authenticated.tsx` | Create: authenticated layout with guard |
| `src/routes/_authenticated/dashboard.tsx` | Create: dashboard/landing page |
| `src/components/layout/app-shell.tsx` | Create: sidebar + main layout |

## Completion Criteria

- [ ] `bun run typecheck && bun run lint` passes
- [ ] Login and register forms validate with Zod via TanStack Form
- [ ] Login and register forms work end-to-end against backend
- [ ] TanStack Query manages auth state (current user, permissions)
- [ ] Dark mode is active by default
- [ ] Sidebar renders with navigation items
- [ ] Unauthenticated users are redirected to /login
- [ ] Authenticated users see the app shell
- [ ] Token refresh works transparently
