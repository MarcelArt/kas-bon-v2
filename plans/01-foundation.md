# Phase 1 — Foundation

Authentication, app shell, dark mode, and the shared API layer.

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

## 3. Auth Server Functions

**File:** `src/lib/auth.server.ts`

Using `createServerFn` from TanStack Start:

- `loginFn(body: LoginInput)` — calls `POST /v1/users/login`, stores tokens in cookies, returns user
- `registerFn(body: UserInput)` — calls `POST /v1/users`, returns new user ID
- `refreshFn()` — calls `POST /v1/users/refresh` with refresh token from cookie
- `logoutFn()` — clears cookies
- `getCurrentUserFn()` — reads user from cookie/session (or re-validates token)
- `getUserPermissionsFn(userId, appId, domainId)` — calls `GET /v1/users/{id}/permissions`, returns permission tuples

## 4. Auth Context & State

**File:** `src/lib/auth.tsx`

React context/provider providing:
- `user: User | null`
- `permissions: Set<string>` — parsed from permission tuples
- `isSuperUser: boolean` — true if permissions has `"all#fullAccess"`
- `hasPermission(resource, action): boolean` — checks permission set
- `login(username, password, isRemember)`
- `register(username, email, password)`
- `logout()`
- `isLoading: boolean`

Permission check logic:
```
hasPermission(resource, action):
  if isSuperUser → return true
  return permissions.has(`${resource}#${action}`)
```

## 5. Route Structure

```
src/routes/
  __root.tsx          # Shell with dark class, sidebar layout
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

## 6. App Shell Layout

**File:** `src/components/layout/app-shell.tsx`

Structure for authenticated pages:
```
<AppShell>
  <Sidebar>          # shadcn sidebar component
    <SidebarHeader>  # App name / logo
    <SidebarContent>
      <SidebarGroup> # Navigation items (gated by permissions)
        <SidebarItem to="/dashboard" icon={House} permission="dashboard#read" />
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

**shadcn components needed:** `sidebar`, `separator`, `tooltip`

## 7. Auth Pages

### Login Page — `src/routes/_auth/login.tsx`

- Form with username, password, "Remember me" checkbox
- Calls `login()` from auth context
- On success: redirect to `/dashboard`
- On error: show error message
- Link to register page

### Register Page — `src/routes/_auth/register.tsx`

- Form with username, email, password, confirm password
- Calls `register()` from auth context
- On success: redirect to `/login` with success message
- On error: show validation errors

**shadcn components needed:** `card`, `input`, `label`, `checkbox`

## 8. Route Guards

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

## 9. shadcn Components to Install

```bash
cd web
bunx shadcn add sidebar separator tooltip card input label checkbox sonner
```

`sonner` for toast notifications (success/error feedback on mutations).

## 10. Files to Create/Modify

| File | Action |
|---|---|
| `src/routes/__root.tsx` | Modify: add `className="dark"` to `<html>`, update title |
| `src/lib/api.ts` | Create: API client with auth headers |
| `src/lib/api.types.ts` | Create: TypeScript types for all backend models |
| `src/lib/auth.server.ts` | Create: server functions for auth API calls |
| `src/lib/auth.tsx` | Create: auth context, provider, permission helpers |
| `src/lib/auth-guard.ts` | Create: route guard utilities |
| `src/routes/_auth.tsx` | Create: unauthenticated layout |
| `src/routes/_auth/login.tsx` | Create: login page |
| `src/routes/_auth/register.tsx` | Create: register page |
| `src/routes/_authenticated.tsx` | Create: authenticated layout with guard |
| `src/routes/_authenticated/dashboard.tsx` | Create: dashboard/landing page |
| `src/components/layout/app-shell.tsx` | Create: sidebar + main layout |

## Completion Criteria

- [ ] `bun run typecheck && bun run lint` passes
- [ ] Login and register forms work end-to-end against backend
- [ ] Dark mode is active by default
- [ ] Sidebar renders with navigation items
- [ ] Unauthenticated users are redirected to /login
- [ ] Authenticated users see the app shell
- [ ] Token refresh works transparently
