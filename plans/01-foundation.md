# Phase 1 — Foundation

Authentication, app shell, dark mode, type-safe environment variables, and server function API layer.

## Step 0: Install Dependencies

```bash
cd web
bun add zustand @tanstack/react-query @tanstack/react-form @tanstack/zod-form-adapter zod
bunx shadcn add sidebar separator tooltip card input label checkbox sonner button
```

No Axios needed — server functions use native `fetch`.

## Step 1: Dark Mode Default

**File:** `src/routes/__root.tsx`

- Add `className="dark"` to the `<html>` element in `RootDocument`
- This activates the `.dark` CSS custom properties already defined in `styles.css`
- Update page title from "TanStack Start Starter" to app name
- Wrap app in `QueryClientProvider` (TanStack Query)

## Step 2: API Types

**File:** `src/lib/api.types.ts`

TypeScript types mirroring all backend models (see `.opencode/API.md` for contracts):

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

## Step 3: Type-Safe Environment Variables

### TypeScript Declarations

**File:** `src/env.d.ts`

Provides type safety for server-side environment variables via `NodeJS.ProcessEnv` declarations:

```typescript
/// <reference types="vite/client" />

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      readonly API_URL?: string
      readonly NODE_ENV: "development" | "production" | "test"
    }
  }
}

export {}
```

This gives full autocomplete and type-checking for `process.env.API_URL` inside server functions.

### Server-side env helper

**File:** `src/lib/env.ts`

Reads `process.env.API_URL` with type safety from `env.d.ts`, falls back to `http://localhost:8080`. Only accessed inside server functions (which are server-only by default).

```typescript
export function getApiBaseUrl(): string {
  return process.env.API_URL || "http://localhost:8080"
}
```

Usage in server functions:
```typescript
import { getApiBaseUrl } from "@/lib/env"

const url = `${getApiBaseUrl()}/api/v1/users`
```

## Step 4: Server API Client

**File:** `src/lib/server/client.ts`

Base fetch wrapper used by all server function handlers. Reads `API_BASE_URL` from type-safe env, unwraps `JSONResponse`, handles errors with status codes.

```typescript
import { getApiBaseUrl } from "@/lib/env"
import type { JSONResponse } from "@/lib/api.types"

interface RequestOptions {
  accessToken?: string | null
  domainId?: number | null
  appId?: number | null
  customHeaders?: Record<string, string>
}

export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(`API_ERROR:${status}:${message}`)
    this.status = status
    this.name = "ApiError"
  }
}

export async function serverApi<T>(
  method: string,
  path: string,
  options?: RequestOptions,
  body?: unknown,
  params?: Record<string, string | number | undefined>,
): Promise<T> {
  const url = new URL(`/api/v1${path}`, getApiBaseUrl())
  if (params) {
    Object.entries(params).forEach(([k, v]) => {
      if (v != null) url.searchParams.set(k, String(v))
    })
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.customHeaders ?? {}),
  }
  if (options?.accessToken) headers["Authorization"] = `Bearer ${options.accessToken}`
  if (options?.domainId != null) headers["X-Domain-Id"] = String(options.domainId)
  if (options?.appId != null) headers["X-App-Id"] = String(options.appId)

  const res = await fetch(url.toString(), {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  })

  const json: JSONResponse<T> = await res.json()

  if (res.status === 401) {
    throw new ApiError("AUTH_EXPIRED:" + (json.message || "Token expired"), 401)
  }

  if (!res.ok) {
    throw new ApiError(json.message || res.statusText, res.status)
  }

  if (!json.isSuccess) {
    throw new ApiError(json.message, res.status)
  }

  return json.items as T
}
```

Error conventions:
- `AUTH_EXPIRED:...` → 401, client should attempt token refresh
- `API_ERROR:403:...` → forbidden
- `API_ERROR:404:...` → not found

## Step 5: Auth Server Functions

**File:** `src/lib/server/auth.ts`

All auth endpoints from `.opencode/API.md`:

```typescript
import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { LoginResponse, User, Domain } from "@/lib/api.types"

export const loginFn = createServerFn()
  .validator(
    z.object({
      username: z.string(),
      password: z.string(),
      isRemember: z.boolean().optional().default(false),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<LoginResponse>("POST", "/users/login", undefined, data)
  })

export const registerFn = createServerFn()
  .validator(
    z.object({
      username: z.string(),
      email: z.string().email(),
      password: z.string().min(6),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<number>("POST", "/users", undefined, data)
  })

export const refreshTokenFn = createServerFn()
  .validator(
    z.object({
      refreshToken: z.string(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<LoginResponse>("POST", "/users/refresh", undefined, undefined, undefined, {
      "X-Refresh-Token": data.refreshToken,
    })
  })

export const getOrganizationsFn = createServerFn()
  .validator(
    z.object({
      accessToken: z.string(),
      userId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<Domain[]>(
      "GET",
      `/users/${data.userId}/organizations`,
      { accessToken: data.accessToken },
    )
  })

export const getPermissionsFn = createServerFn()
  .validator(
    z.object({
      accessToken: z.string(),
      userId: z.number(),
      domainId: z.number().optional(),
      appId: z.number().optional(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<string[][]>(
      "GET",
      `/users/${data.userId}/permissions`,
      { accessToken: data.accessToken, domainId: data.domainId, appId: data.appId },
    )
  })
```

## Step 6: Zustand Auth Store

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
    },
  ),
)
```

## Step 7: Auth Context Helper & Retry Wrapper

**File:** `src/lib/queries/auth-context.ts`

Helper to extract auth context from Zustand store for server function calls, plus a retry wrapper for handling expired tokens:

```typescript
import { useAuthStore } from "@/lib/stores/auth-store"
import { refreshTokenFn } from "@/lib/server/auth"

export function getAuthContext() {
  const { accessToken, domainId, appId } = useAuthStore.getState()
  return {
    accessToken: accessToken ?? undefined,
    domainId: domainId ?? undefined,
    appId: appId ?? undefined,
  }
}

export function useAuthContext() {
  const accessToken = useAuthStore((s) => s.accessToken)
  const domainId = useAuthStore((s) => s.domainId)
  const appId = useAuthStore((s) => s.appId)
  return {
    accessToken: accessToken ?? undefined,
    domainId: domainId ?? undefined,
    appId: appId ?? undefined,
  }
}

export function withAuthRetry<T>(queryFn: () => Promise<T>): () => Promise<T> {
  return async () => {
    try {
      return await queryFn()
    } catch (error) {
      if (error instanceof Error && error.message.includes("AUTH_EXPIRED:")) {
        const { refreshToken, setTokens, logout } = useAuthStore.getState()
        if (!refreshToken) {
          logout()
          window.location.href = "/login"
          throw error
        }
        try {
          const result = await refreshTokenFn({ data: { refreshToken } })
          setTokens(result.accessToken, result.refreshToken)
          return await queryFn()
        } catch {
          logout()
          window.location.href = "/login"
          throw error
        }
      }
      throw error
    }
  }
}
```

**IMPORTANT:** `getAuthContext()` (uses `getState()`) is for non-React contexts. `useAuthContext()` (hook) is for React components and TanStack Query hooks.

## Step 8: TanStack Query Setup & Auth Hooks

**File:** `src/lib/queries/auth.ts`

```typescript
import { useMutation, useQuery } from "@tanstack/react-query"
import { loginFn, getOrganizationsFn, getPermissionsFn } from "@/lib/server/auth"
import { useAuthStore } from "@/lib/stores/auth-store"
import { useAuthContext, withAuthRetry } from "./auth-context"
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
    mutationFn: (data: { username: string; password: string; isRemember?: boolean }) =>
      loginFn({ data }),
    onSuccess: async (loginData) => {
      setUser(loginData.user)
      setTokens(loginData.accessToken, loginData.refreshToken)

      const orgs = await getOrganizationsFn({
        data: { accessToken: loginData.accessToken, userId: loginData.user.ID },
      })
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
  const auth = useAuthContext()
  return useQuery({
    queryKey: authKeys.organizations(userId),
    queryFn: withAuthRetry(() =>
      getOrganizationsFn({ data: { accessToken: auth.accessToken!, userId } }),
    ),
    enabled: !!auth.accessToken && !!userId,
  })
}

export function useUserPermissions(userId: number) {
  const auth = useAuthContext()
  const { setPermissions } = useAuthStore()

  return useQuery({
    queryKey: authKeys.permissions(userId),
    queryFn: withAuthRetry(async () => {
      const tuples = await getPermissionsFn({
        data: { accessToken: auth.accessToken!, userId, domainId: auth.domainId, appId: auth.appId },
      })
      setPermissions(tuples)
      return tuples
    }),
    enabled: !!auth.accessToken && !!userId,
  })
}
```

## Step 9: Zod Validation Schemas

**File:** `src/lib/schemas/auth.ts`

Shared between TanStack Form and server function validators:

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

## Step 10: Route Structure

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

## Step 11: App Shell Layout

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

## Step 12: Auth Pages (TanStack Form + Zod)

### Login Page — `src/routes/_auth/login.tsx`

- Use `useForm` from `@tanstack/react-form` with `zodValidator` and `loginSchema`
- Form fields: username, password, "Remember me" checkbox
- On submit: call `useLogin()` mutation (which calls `loginFn` server function)
- On success: mutation handles org selection flow + redirect automatically
- On error: display `mutation.error.message` via `sonner` toast
- Link to register page

### Register Page — `src/routes/_auth/register.tsx`

- Use `useForm` with `zodValidator` and `registerSchema`
- Form fields: username, email, password, confirm password
- On submit: call `registerFn` server function, on success redirect to `/login` with toast

## Step 13: Organization Selection Page

### Select Organization — `src/routes/_authenticated/select-organization.tsx`

This page is shown when a user has multiple organizations after login.

- Read `organizations` from Zustand auth store (already populated by login flow)
- Display as a card/list for user to pick one
- On select:
  1. Call `setDomain(selectedOrg.ID)` in auth store
  2. Optionally fetch apps for this org to set a default `appId`
  3. Fetch permissions via `getPermissionsFn` server function (with auth context)
  4. Store permissions in auth store via `setPermissions(tuples)`
  5. Navigate to `/dashboard`

### No Access — `src/routes/_authenticated/no-access.tsx`

Simple page: "You don't have access to any organization. Contact your administrator."

## Step 14: Route Guards

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

## Step 15: Permission Utility

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
  action: string,
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

## Step 16: Permission Hooks

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

## Step 17: Files to Create/Modify

| File | Action |
|---|---|
| `src/routes/__root.tsx` | Modify: add `className="dark"` to `<html>`, add `QueryClientProvider`, update title |
| `src/lib/api.types.ts` | Create: TypeScript types for all backend models |
| `src/env.d.ts` | Create: TypeScript declarations for `NodeJS.ProcessEnv` (type-safe env vars) |
| `src/lib/env.ts` | Create: `getApiBaseUrl()` helper, fallback `http://localhost:8080` |
| `src/lib/server/client.ts` | Create: base server API client (fetch wrapper, error handling) |
| `src/lib/server/auth.ts` | Create: auth server functions (login, register, refresh, orgs, permissions) |
| `src/lib/stores/auth-store.ts` | Create: Zustand auth store |
| `src/lib/queries/auth-context.ts` | Create: auth context helper + `withAuthRetry` wrapper |
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

## Step 18: Install shadcn Components

```bash
cd web
bunx shadcn add sidebar separator tooltip card input label checkbox sonner button
```

## Completion Criteria

- [ ] `bun run typecheck && bun run lint` passes
- [ ] `src/env.d.ts` provides TypeScript declarations for `process.env.API_URL`
- [ ] `getApiBaseUrl()` returns `process.env.API_URL` with fallback to `http://localhost:8080`
- [ ] Login form works with TanStack Form + Zod → calls `loginFn` server function
- [ ] Register form works with TanStack Form + Zod → calls `registerFn` server function
- [ ] After login, organizations are fetched via `getOrganizationsFn` server function
- [ ] Single org → auto-selected, redirected to dashboard
- [ ] Multiple orgs → redirected to `/select-organization` picker
- [ ] Zero orgs → redirected to `/no-access` page
- [ ] Organization selection updates auth context in Zustand
- [ ] Permissions fetched via `getPermissionsFn` server function and stored in Zustand
- [ ] Dark mode is active by default
- [ ] Sidebar renders with navigation items (gated by permissions)
- [ ] Unauthenticated users are redirected to `/login`
- [ ] Authenticated users without org selection are redirected to `/select-organization`
- [ ] Token refresh works via `withAuthRetry` wrapper (calls `refreshTokenFn` server function)
