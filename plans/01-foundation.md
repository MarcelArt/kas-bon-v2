# Phase 1 — Foundation

Authentication, app shell, dark mode, and the shared API layer.

## Step 0: Install Dependencies

```bash
cd web
bun add axios zustand @tanstack/react-query @tanstack/react-form @tanstack/zod-form-adapter zod
bunx shadcn add sidebar separator tooltip card input label checkbox sonner button
```

## Step 1: Dark Mode Default

**File:** `src/routes/__root.tsx`

- Add `className="dark"` to the `<html>` element in `RootDocument`
- This activates the `.dark` CSS custom properties already defined in `styles.css`
- Update page title from "TanStack Start Starter" to app name
- Wrap app in `QueryClientProvider` (TanStack Query)

## Step 2: API Types

**File:** `src/lib/api.types.ts`

TypeScript types mirroring all backend models:

```typescript
interface JSONResponse<T> {
  items: T | null
  isSuccess: boolean
  message: string
}

interface PaginatedResponse<T> {
  items: T[]
  page: number
  size: number
  total: number
}

interface User {
  ID: number
  username: string
  email: string
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
}

interface LoginInput {
  username: string
  password: string
  isRemember: boolean
}

interface LoginResponse {
  accessToken: string
  refreshToken: string
  user: User
}

interface Domain {
  ID: number
  name: string
  description: string
  isOrganization: boolean
  parentId: number | null
  parent: Domain | null
  CreatedAt: string
  UpdatedAt: string
}

interface App {
  ID: number
  name: string
  description: string
  CreatedAt: string
  UpdatedAt: string
}

interface Role {
  ID: number
  name: string
  description: string
  domainId: number
  domain: Domain | null
  CreatedAt: string
  UpdatedAt: string
}

interface Permission {
  ID: number
  name: string
  description: string
  appId: number
  app: App | null
  CreatedAt: string
  UpdatedAt: string
}
```

## Step 3: Axios API Client

**File:** `src/lib/api.ts`

Shared Axios instance for all backend calls:

```typescript
import axios from "axios"

export const api = axios.create({
  baseURL: "/api/v1",
})

// Request interceptor: attach auth headers from Zustand store
api.interceptors.request.use((config) => {
  const { accessToken, refreshToken } = useAuthStore.getState()
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`
  }
  const { domainId, appId } = useAuthStore.getState()
  if (domainId) config.headers["X-Domain-Id"] = domainId
  if (appId) config.headers["X-App-Id"] = appId
  return config
})

// Response interceptor: handle 401 → attempt refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const { refreshToken, setTokens, logout } = useAuthStore.getState()
      if (!refreshToken) {
        logout()
        window.location.href = "/login"
        return
      }
      try {
        const res = await axios.post("/api/v1/users/refresh", null, {
          headers: { "X-Refresh-Token": refreshToken },
        })
        const { accessToken: newAccess, refreshToken: newRefresh } = res.data.items
        setTokens(newAccess, newRefresh)
        error.config.headers.Authorization = `Bearer ${newAccess}`
        return api.request(error.config)
      } catch {
        logout()
        window.location.href = "/login"
      }
    }
    return Promise.reject(error)
  }
)
```

Helper to unwrap JSONResponse:
```typescript
export function unwrap<T>(res: { data: JSONResponse<T> }): T {
  if (!res.data.isSuccess) throw new Error(res.data.message)
  return res.data.items!
}
```

## Step 4: Zustand Auth Store

**File:** `src/lib/stores/auth-store.ts`

```typescript
import { create } from "zustand"
import { persist } from "zustand/middleware"

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  domainId: number | null
  appId: number | null
  organizations: Domain[]
  permissions: Set<string>

  setUser: (user: User | null) => void
  setTokens: (access: string, refresh: string) => void
  setDomain: (domainId: number) => void
  setApp: (appId: number) => void
  setOrganizations: (orgs: Domain[]) => void
  setPermissions: (tuples: string[][]) => void
  hasPermission: (resource: string, action: string) => boolean
  logout: () => void
}

function parsePermissions(tuples: string[][]): Set<string> {
  return new Set(tuples.map((t) => `${t[3]}#${t[4]}`))
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      domainId: null,
      appId: null,
      organizations: [],
      permissions: new Set(),

      setUser: (user) => set({ user }),
      setTokens: (accessToken, refreshToken) => set({ accessToken, refreshToken }),
      setDomain: (domainId) => set({ domainId }),
      setApp: (appId) => set({ appId }),
      setOrganizations: (organizations) => set({ organizations }),
      setPermissions: (tuples) => {
        const permissions = parsePermissions(tuples)
        set({ permissions })
      },
      hasPermission: (resource, action) => {
        const { permissions } = get()
        if (permissions.has("all#fullAccess")) return true
        return permissions.has(`${resource}#${action}`)
      },
      logout: () =>
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          domainId: null,
          appId: null,
          organizations: [],
          permissions: new Set(),
        }),
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        domainId: state.domainId,
        appId: state.appId,
      }),
    }
  )
)
```

**IMPORTANT:** The store is imported by `api.ts` for interceptors. Use `useAuthStore.getState()` in non-React contexts (interceptors) and `useAuthStore()` hook in components.

## Step 5: API Functions (Axios Wrappers)

**File:** `src/lib/api/auth.ts`

```typescript
import { api, unwrap } from "@/lib/api"
import type { LoginInput, LoginResponse, User, Domain } from "@/lib/api.types"

export const authApi = {
  login: (body: LoginInput) =>
    api.post<JSONResponse<LoginResponse>>("/v1/users/login", body).then(unwrap),

  register: (body: { username: string; email: string; password: string }) =>
    api.post<JSONResponse<number>>("/v1/users", body).then(unwrap),

  refresh: (refreshToken: string) =>
    api
      .post<JSONResponse<LoginResponse>>("/v1/users/refresh", null, {
        headers: { "X-Refresh-Token": refreshToken },
      })
      .then(unwrap),

  getOrganizations: (userId: number) =>
    api.get<JSONResponse<Domain[]>>(`/v1/users/${userId}/organizations`).then(unwrap),

  getPermissions: (userId: number) =>
    api.get<JSONResponse<string[][]>>(`/v1/users/${userId}/permissions`).then(unwrap),
}
```

## Step 6: TanStack Query Setup & Auth Hooks

**File:** `src/lib/queries/auth.ts`

```typescript
import { useMutation, useQuery } from "@tanstack/react-query"
import { authApi } from "@/lib/api/auth"
import { useAuthStore } from "@/lib/stores/auth-store"
import { useNavigate } from "@tanstack/react-router"

export const authKeys = {
  user: ["auth", "user"] as const,
  organizations: (userId: number) => ["auth", "organizations", userId] as const,
  permissions: (userId: number) => ["auth", "permissions", userId] as const,
}

export function useLogin() {
  const { setUser, setTokens, setOrganizations, setDomain, setApp } = useAuthStore()
  const navigate = useNavigate()

  return useMutation({
    mutationFn: authApi.login,
    onSuccess: async (data: LoginResponse) => {
      setUser(data.user)
      setTokens(data.accessToken, data.refreshToken)

      const orgs = await authApi.getOrganizations(data.user.ID)
      setOrganizations(orgs)

      if (orgs.length === 0) {
        navigate({ to: "/no-access" })
      } else if (orgs.length === 1) {
        setDomain(orgs[0].ID)
        // TODO: fetch default app for this org, then setApp(appId)
        navigate({ to: "/dashboard" })
      } else {
        navigate({ to: "/select-organization" })
      }
    },
  })
}

export function useUserOrganizations(userId: number) {
  return useQuery({
    queryKey: authKeys.organizations(userId),
    queryFn: () => authApi.getOrganizations(userId),
    enabled: !!userId,
  })
}

export function useUserPermissions(userId: number) {
  const { setPermissions } = useAuthStore()

  return useQuery({
    queryKey: authKeys.permissions(userId),
    queryFn: async () => {
      const tuples = await authApi.getPermissions(userId)
      setPermissions(tuples)
      return tuples
    },
    enabled: !!userId,
  })
}
```

## Step 7: Zod Validation Schemas

**File:** `src/lib/schemas/auth.ts`

```typescript
import { z } from "zod"

export const loginSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(6, "Password must be at least 6 characters"),
  isRemember: z.boolean().optional().default(false),
})

export const registerSchema = z
  .object({
    username: z.string().min(1, "Username is required"),
    email: z.string().email("Invalid email address"),
    password: z.string().min(6, "Password must be at least 6 characters"),
    confirmPassword: z.string().min(6, "Confirm your password"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
  })

export type LoginFormData = z.infer<typeof loginSchema>
export type RegisterFormData = z.infer<typeof registerSchema>
```

## Step 8: Route Structure

```
src/routes/
  __root.tsx              # Shell with dark class, QueryClientProvider, sidebar layout
  _auth.tsx               # Layout for unauthenticated pages (no sidebar)
  _auth/
    login.tsx             # /login
    register.tsx          # /register
  _authenticated.tsx      # Layout with sidebar, requires auth + org selected
  _authenticated/
    dashboard.tsx         # /dashboard (redirect here after login + org selection)
    select-organization.tsx  # /select-organization (pick org if multiple)
    no-access.tsx         # /no-access (user has 0 orgs)
    users/
      index.tsx           # /users
      $userId.tsx         # /users/:id
    domains/
      index.tsx           # /domains
      $domainId.tsx       # /domains/:id
    apps/
      index.tsx           # /apps
      $appId.tsx          # /apps/:id
    roles/
      index.tsx           # /roles
      $roleId.tsx         # /roles/:id
    permissions/
      index.tsx           # /permissions
      $permissionId.tsx   # /permissions/:id
```

Use TanStack Router layout routes (`_auth` and `_authenticated`) for shared layout wrappers.

## Step 9: App Shell Layout

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
    <SidebarFooter>  # User info + current org + logout
  </Sidebar>
  <main>             # Page content via <Outlet />
</AppShell>
```

**shadcn components needed:** `sidebar`, `separator`, `tooltip`

## Step 10: Auth Pages (TanStack Form + Zod)

### Login Page — `src/routes/_auth/login.tsx`

- Use `useForm` from `@tanstack/react-form` with `zodValidator` and `loginSchema`
- Form fields: username, password, "Remember me" checkbox
- On submit: call `useLogin()` mutation
- On success: mutation handles org selection flow + redirect automatically
- On error: display `mutation.error.message` via `sonner` toast
- Link to register page

### Register Page — `src/routes/_auth/register.tsx`

- Use `useForm` with `zodValidator` and `registerSchema`
- Form fields: username, email, password, confirm password
- On submit: call register API, on success redirect to `/login` with toast

## Step 11: Organization Selection Page

### Select Organization — `src/routes/_authenticated/select-organization.tsx`

This page is shown when a user has multiple organizations after login.

- Read `organizations` from Zustand auth store (already populated by login flow)
- Display as a card/list for user to pick one
- On select:
  1. Call `setDomain(selectedOrg.id)` in auth store
  2. Optionally fetch apps for this org to set a default `appId`
  3. Fetch permissions: `GET /v1/users/{userId}/permissions` (with new X-Domain-Id, X-App-Id)
  4. Store permissions in auth store via `setPermissions(tuples)`
  5. Navigate to `/dashboard`

### No Access — `src/routes/_authenticated/no-access.tsx`

Simple page: "You don't have access to any organization. Contact your administrator."

## Step 12: Route Guards

**File:** `src/lib/auth-guard.ts`

In `_authenticated.tsx` layout route, use `beforeLoad` to:
1. Check if user is authenticated (tokens exist in Zustand store)
2. If not → redirect to `/login`
3. Check if organization is selected (`domainId` exists in store)
4. If not → redirect to `/select-organization`
5. Load user permissions if not already loaded (or if domain changed)
6. Inject `permissions` and `hasPermission` into route context

```typescript
// _authenticated.tsx
export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async ({ context }) => {
    const { user, domainId, accessToken } = useAuthStore.getState()
    if (!user || !accessToken) {
      throw redirect({ to: "/login" })
    }
    if (!domainId) {
      throw redirect({ to: "/select-organization" })
    }
  },
})
```

## Step 13: Permission Utility

**File:** `src/lib/permissions.ts`

```typescript
export function parsePermissionTuples(tuples: string[][]): Set<string> {
  return new Set(tuples.map((t) => `${t[3]}#${t[4]}`))
}

export function isSuperUser(permissions: Set<string>): boolean {
  return permissions.has("all#fullAccess")
}

export function checkPermission(
  permissions: Set<string>,
  resource: string,
  action: string
): boolean {
  if (isSuperUser(permissions)) return true
  return permissions.has(`${resource}#${action}`)
}

export const RESOURCES = {
  USERS: "users",
  DOMAINS: "domains",
  APPS: "apps",
  ROLES: "roles",
  PERMISSIONS: "permissions",
  ALL: "all",
} as const

export const ACTIONS = {
  READ: "read",
  CREATE: "create",
  UPDATE: "update",
  DELETE: "delete",
  FULL_ACCESS: "fullAccess",
} as const
```

## Step 14: Permission Hooks

**File:** `src/hooks/use-permission.ts`

```typescript
import { useAuthStore } from "@/lib/stores/auth-store"

export function usePermission(resource: string, action: string): boolean {
  const hasPermission = useAuthStore((s) => s.hasPermission)
  return hasPermission(resource, action)
}

export function useCanCreate(resource: string): boolean {
  return usePermission(resource, "create")
}

export function useCanEdit(resource: string): boolean {
  return usePermission(resource, "update")
}

export function useCanDelete(resource: string): boolean {
  return usePermission(resource, "delete")
}

export function useIsSuperUser(): boolean {
  const permissions = useAuthStore((s) => s.permissions)
  return permissions.has("all#fullAccess")
}
```

## Step 15: Files to Create/Modify

| File | Action |
|---|---|
| `src/routes/__root.tsx` | Modify: add `className="dark"` to `<html>`, add `QueryClientProvider`, update title |
| `src/lib/api.types.ts` | Create: TypeScript types for all backend models |
| `src/lib/api.ts` | Create: Axios instance with auth interceptors |
| `src/lib/stores/auth-store.ts` | Create: Zustand auth store |
| `src/lib/api/auth.ts` | Create: auth API functions (login, register, refresh, orgs, permissions) |
| `src/lib/queries/auth.ts` | Create: TanStack Query hooks for auth (useLogin, useUserPermissions, etc.) |
| `src/lib/schemas/auth.ts` | Create: Zod validation schemas for login/register forms |
| `src/lib/permissions.ts` | Create: permission parsing utility |
| `src/hooks/use-permission.ts` | Create: permission hooks |
| `src/routes/_auth.tsx` | Create: unauthenticated layout |
| `src/routes/_auth/login.tsx` | Create: login page (TanStack Form + Zod) |
| `src/routes/_auth/register.tsx` | Create: register page (TanStack Form + Zod) |
| `src/routes/_authenticated.tsx` | Create: authenticated layout with guard (checks org selection) |
| `src/routes/_authenticated/dashboard.tsx` | Create: dashboard/landing page |
| `src/routes/_authenticated/select-organization.tsx` | Create: organization picker page |
| `src/routes/_authenticated/no-access.tsx` | Create: no-access page |
| `src/components/layout/app-shell.tsx` | Create: sidebar + main layout |

## Step 16: Install shadcn Components

```bash
cd web
bunx shadcn add sidebar separator tooltip card input label checkbox sonner button
```

## Completion Criteria

- [ ] `bun run typecheck && bun run lint` passes
- [ ] Login form works with TanStack Form + Zod validation against backend
- [ ] Register form works with TanStack Form + Zod validation against backend
- [ ] After login, organizations are fetched from `/v1/users/{id}/organizations`
- [ ] Single org → auto-selected, redirected to dashboard
- [ ] Multiple orgs → redirected to `/select-organization` picker
- [ ] Zero orgs → redirected to `/no-access` page
- [ ] Organization selection sets `X-Domain-Id` in Axios interceptor
- [ ] Permissions fetched and stored in Zustand after org selection
- [ ] Dark mode is active by default
- [ ] Sidebar renders with navigation items (gated by permissions)
- [ ] Unauthenticated users are redirected to `/login`
- [ ] Authenticated users without org selection are redirected to `/select-organization`
- [ ] Token refresh works transparently via Axios interceptor
