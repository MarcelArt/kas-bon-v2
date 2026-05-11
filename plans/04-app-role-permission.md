# Phase 4 — Apps, Roles & Permissions

CRUD for apps, roles, and permissions, plus permission assignment to roles. Uses TanStack Start server functions, TanStack Query hooks, TanStack Form + Zod.

## Step 1: App Management

### App Server Functions

**File:** `src/lib/server/apps.ts`

Server functions for app endpoints (see `.opencode/API.md` — Apps section):

```typescript
import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { App, PaginatedResponse } from "@/lib/api.types"

const authContextSchema = z.object({
  accessToken: z.string(),
  domainId: z.number().optional(),
  appId: z.number().optional(),
})

export const getAppsFn = createServerFn()
  .validator(
    authContextSchema.extend({
      page: z.number().optional(),
      size: z.number().optional(),
      sort: z.string().optional(),
      filters: z.string().optional(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...params } = data
    return serverApi<PaginatedResponse<App>>("GET", "/apps", { accessToken, domainId, appId }, undefined, params)
  })

export const getAppFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<App>("GET", `/apps/${id}`, { accessToken, domainId, appId })
  })

export const createAppFn = createServerFn()
  .validator(
    authContextSchema.extend({
      name: z.string(),
      description: z.string().optional().default(""),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...body } = data
    return serverApi<number>("POST", "/apps", { accessToken, domainId, appId }, body)
  })

export const updateAppFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      name: z.string(),
      description: z.string().optional().default(""),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id, ...body } = data
    return serverApi<null>("PUT", `/apps/${id}`, { accessToken, domainId, appId }, body)
  })

export const deleteAppFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<null>("DELETE", `/apps/${id}`, { accessToken, domainId, appId })
  })
```

### Query Keys & Hooks

**File:** `src/lib/queries/apps.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { getAppsFn, getAppFn, createAppFn, updateAppFn, deleteAppFn } from "@/lib/server/apps"
import { useAuthContext, withAuthRetry } from "./auth-context"

export const appKeys = {
  all: ["apps"] as const,
  lists: () => [...appKeys.all, "list"] as const,
  list: (filters: object) => [...appKeys.lists(), filters] as const,
  details: () => [...appKeys.all, "detail"] as const,
  detail: (id: number) => [...appKeys.details(), id] as const,
}

export function useApps(filters?: { page?: number; size?: number; sort?: string; filters?: string }) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: appKeys.list(filters ?? {}),
    queryFn: withAuthRetry(() => getAppsFn({ data: { ...auth, ...filters } })),
    enabled: !!auth.accessToken,
  })
}

export function useApp(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: appKeys.detail(id),
    queryFn: withAuthRetry(() => getAppFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useCreateApp() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (body: { name: string; description?: string }) =>
      createAppFn({ data: { ...auth, ...body } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: appKeys.all }),
  })
}

export function useUpdateApp() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: ({ id, ...body }: { id: number; name: string; description?: string }) =>
      updateAppFn({ data: { ...auth, id, ...body } }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: appKeys.detail(id) })
      qc.invalidateQueries({ queryKey: appKeys.lists() })
    },
  })
}

export function useDeleteApp() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (id: number) => deleteAppFn({ data: { ...auth, id } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: appKeys.all }),
  })
}
```

### Zod Schema

**File:** `src/lib/schemas/app.ts`

```typescript
import { z } from "zod"

export const appSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string().optional().default(""),
})
```

### App List Page

**File:** `src/routes/_authenticated/apps/index.tsx`

Simple table: Name, Description, Created At, Actions. Permission gating with hooks.
- Create: `useCanCreate("apps")`
- Edit: `useCanEdit("apps")`
- Delete: `useCanDelete("apps")`

### App Detail Page

**File:** `src/routes/_authenticated/apps/$appId.tsx`

Display name, description, timestamps. Edit form (TanStack Form + `appSchema`) and delete button.

### App CRUD Dialogs (TanStack Form + Zod)

- `src/components/apps/create-app-dialog.tsx`
- `src/components/apps/edit-app-dialog.tsx`
- `src/components/apps/delete-app-dialog.tsx`

---

## Step 2: Role Management

### Role Server Functions

**File:** `src/lib/server/roles.ts`

Server functions for role endpoints (see `.opencode/API.md` — Roles section):

```typescript
import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { Role, PaginatedResponse } from "@/lib/api.types"

const authContextSchema = z.object({
  accessToken: z.string(),
  domainId: z.number().optional(),
  appId: z.number().optional(),
})

export const getRolesFn = createServerFn()
  .validator(
    authContextSchema.extend({
      page: z.number().optional(),
      size: z.number().optional(),
      sort: z.string().optional(),
      filters: z.string().optional(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...params } = data
    return serverApi<PaginatedResponse<Role>>("GET", "/roles", { accessToken, domainId, appId }, undefined, params)
  })

export const getRoleFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<Role>("GET", `/roles/${id}`, { accessToken, domainId, appId })
  })

export const createRoleFn = createServerFn()
  .validator(
    authContextSchema.extend({
      name: z.string(),
      description: z.string().optional().default(""),
      domainId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, appId, ...body } = data
    return serverApi<number>("POST", "/roles", { accessToken, domainId: body.domainId, appId }, body)
  })

export const updateRoleFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      name: z.string(),
      description: z.string().optional().default(""),
      domainId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, appId, id, ...body } = data
    return serverApi<null>("PUT", `/roles/${id}`, { accessToken, domainId: body.domainId, appId }, body)
  })

export const deleteRoleFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<null>("DELETE", `/roles/${id}`, { accessToken, domainId, appId })
  })

export const getRolePermissionsFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<string[][]>("GET", `/roles/${id}/permissions`, { accessToken, domainId, appId })
  })

export const assignRolePermissionsFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      permissionIds: z.array(z.number()),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id, permissionIds } = data
    return serverApi<string[]>("PATCH", `/roles/${id}/permissions`, { accessToken, domainId, appId }, permissionIds)
  })
```

### Query Keys & Hooks

**File:** `src/lib/queries/roles.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  getRolesFn,
  getRoleFn,
  createRoleFn,
  updateRoleFn,
  deleteRoleFn,
  getRolePermissionsFn,
  assignRolePermissionsFn,
} from "@/lib/server/roles"
import { useAuthContext, withAuthRetry } from "./auth-context"

export const roleKeys = {
  all: ["roles"] as const,
  lists: () => [...roleKeys.all, "list"] as const,
  list: (filters: object) => [...roleKeys.lists(), filters] as const,
  details: () => [...roleKeys.all, "detail"] as const,
  detail: (id: number) => [...roleKeys.details(), id] as const,
  permissions: (id: number) => [...roleKeys.detail(id), "permissions"] as const,
}

export function useRoles(filters?: { page?: number; size?: number; sort?: string; filters?: string }) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: roleKeys.list(filters ?? {}),
    queryFn: withAuthRetry(() => getRolesFn({ data: { ...auth, ...filters } })),
    enabled: !!auth.accessToken,
  })
}

export function useRole(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: roleKeys.detail(id),
    queryFn: withAuthRetry(() => getRoleFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useRolePermissions(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: roleKeys.permissions(id),
    queryFn: withAuthRetry(() => getRolePermissionsFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useCreateRole() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (body: { name: string; description?: string; domainId: number }) =>
      createRoleFn({ data: { ...auth, ...body } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: roleKeys.all }),
  })
}

export function useUpdateRole() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: ({
      id,
      ...body
    }: {
      id: number
      name: string
      description?: string
      domainId: number
    }) => updateRoleFn({ data: { ...auth, id, ...body } }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: roleKeys.detail(id) })
      qc.invalidateQueries({ queryKey: roleKeys.lists() })
    },
  })
}

export function useDeleteRole() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (id: number) => deleteRoleFn({ data: { ...auth, id } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: roleKeys.all }),
  })
}

export function useAssignRolePermissions() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: ({ roleId, permissionIds }: { roleId: number; permissionIds: number[] }) =>
      assignRolePermissionsFn({ data: { ...auth, id: roleId, permissionIds } }),
    onSuccess: (_, { roleId }) => {
      qc.invalidateQueries({ queryKey: roleKeys.permissions(roleId) })
    },
  })
}
```

### Zod Schema

**File:** `src/lib/schemas/role.ts`

```typescript
import { z } from "zod"

export const roleSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string().optional().default(""),
  domainId: z.number().min(1, "Domain is required"),
})
```

### Role List Page

**File:** `src/routes/_authenticated/roles/index.tsx`

Table: Name, Description, Domain, Actions. Permission gating with hooks.

### Role Detail Page

**File:** `src/routes/_authenticated/roles/$roleId.tsx`

Three sections:

**Info Section**
- Name, description, domain (linked to domain detail)
- Edit form (TanStack Form + `roleSchema`) — `useCanEdit("roles")`
- Delete button — `useCanDelete("roles")`

**Permissions Section**
- List of assigned permissions from `useRolePermissions(id)`
- "Assign Permissions" button — `useCanEdit("roles")`
- Assignment dialog shows all available permissions as checkboxes
- On submit, calls `useAssignRolePermissions()` mutation

### Role CRUD Dialogs (TanStack Form + Zod)

- `src/components/roles/create-role-dialog.tsx`
- `src/components/roles/edit-role-dialog.tsx`
- `src/components/roles/delete-role-dialog.tsx`
- `src/components/roles/assign-permissions-dialog.tsx`

---

## Step 3: Permission Management

### Permission Server Functions

**File:** `src/lib/server/permissions.ts`

Server functions for permission endpoints (see `.opencode/API.md` — Permissions section):

```typescript
import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { Permission, PaginatedResponse } from "@/lib/api.types"

const authContextSchema = z.object({
  accessToken: z.string(),
  domainId: z.number().optional(),
  appId: z.number().optional(),
})

export const getPermissionsFn = createServerFn()
  .validator(
    authContextSchema.extend({
      page: z.number().optional(),
      size: z.number().optional(),
      sort: z.string().optional(),
      filters: z.string().optional(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...params } = data
    return serverApi<PaginatedResponse<Permission>>("GET", "/permissions", { accessToken, domainId, appId }, undefined, params)
  })

export const getPermissionFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<Permission>("GET", `/permissions/${id}`, { accessToken, domainId, appId })
  })

export const createPermissionFn = createServerFn()
  .validator(
    authContextSchema.extend({
      name: z.string(),
      description: z.string().optional().default(""),
      appId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, ...body } = data
    return serverApi<number>("POST", "/permissions", { accessToken, domainId, appId: body.appId }, body)
  })

export const updatePermissionFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      name: z.string(),
      description: z.string().optional().default(""),
      appId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, id, ...body } = data
    return serverApi<null>("PUT", `/permissions/${id}`, { accessToken, domainId, appId: body.appId }, body)
  })

export const deletePermissionFn = createServerFn()
  .validator(authContextSchema.extend({ id: z.number() }))
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<null>("DELETE", `/permissions/${id}`, { accessToken, domainId, appId })
  })
```

### Query Keys & Hooks

**File:** `src/lib/queries/permissions.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  getPermissionsFn,
  getPermissionFn,
  createPermissionFn,
  updatePermissionFn,
  deletePermissionFn,
} from "@/lib/server/permissions"
import { useAuthContext, withAuthRetry } from "./auth-context"

export const permissionKeys = {
  all: ["permissions"] as const,
  lists: () => [...permissionKeys.all, "list"] as const,
  list: (filters: object) => [...permissionKeys.lists(), filters] as const,
  details: () => [...permissionKeys.all, "detail"] as const,
  detail: (id: number) => [...permissionKeys.details(), id] as const,
}

export function usePermissions(filters?: { page?: number; size?: number; sort?: string; filters?: string }) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: permissionKeys.list(filters ?? {}),
    queryFn: withAuthRetry(() => getPermissionsFn({ data: { ...auth, ...filters } })),
    enabled: !!auth.accessToken,
  })
}

export function usePermission(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: permissionKeys.detail(id),
    queryFn: withAuthRetry(() => getPermissionFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useCreatePermission() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (body: { name: string; description?: string; appId: number }) =>
      createPermissionFn({ data: { ...auth, ...body } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: permissionKeys.all }),
  })
}

export function useUpdatePermission() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: ({
      id,
      ...body
    }: {
      id: number
      name: string
      description?: string
      appId: number
    }) => updatePermissionFn({ data: { ...auth, id, ...body } }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: permissionKeys.detail(id) })
      qc.invalidateQueries({ queryKey: permissionKeys.lists() })
    },
  })
}

export function useDeletePermission() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (id: number) => deletePermissionFn({ data: { ...auth, id } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: permissionKeys.all }),
  })
}
```

### Zod Schema

**File:** `src/lib/schemas/permission.ts`

```typescript
import { z } from "zod"

export const permissionSchema = z.object({
  name: z.string().regex(/^[a-zA-Z]+#[a-zA-Z]+$/, "Must be resource#action format"),
  description: z.string().optional().default(""),
  appId: z.number().min(1, "App is required"),
})
```

### Permission List Page

**File:** `src/routes/_authenticated/permissions/index.tsx`

Table: Name (format: `resource#action`), Description, App, Actions.

### Permission Detail Page

**File:** `src/routes/_authenticated/permissions/$permissionId.tsx`

Display name, description, app (linked). Edit form and delete button.

### Permission CRUD Dialogs (TanStack Form + Zod)

- `src/components/permissions/create-permission-dialog.tsx`
- `src/components/permissions/edit-permission-dialog.tsx`
- `src/components/permissions/delete-permission-dialog.tsx`

### Permission Naming Convention

Permission names follow `resource#action` format. When creating/editing:
- Resource field: dropdown or text input (e.g., `users`, `roles`, `domains`, `apps`, `permissions`, `all`)
- Action field: dropdown or text input (e.g., `read`, `create`, `update`, `delete`, `fullAccess`)
- Auto-combined into `resource#action` before sending to API

---

## Step 4: Assign Permissions to Role Dialog

**File:** `src/components/roles/assign-permissions-dialog.tsx`

- Fetches all permissions via `usePermissions()`
- Fetches currently assigned permissions via `useRolePermissions(roleId)`
- Displays permissions grouped by resource (parse `resource#action` name)
- Checkboxes for each permission, pre-checked if assigned
- On submit, sends array of selected permission IDs via `useAssignRolePermissions()` mutation

## Step 5: Reusable Components

### Data Table

**File:** `src/components/shared/data-table.tsx`

Generic table component accepting TanStack Query paginated data:
- `columns: ColumnDef[]`
- `data: T[]`
- `pagination: PaginationState`
- `onPaginationChange: (state) => void`
- Actions column with configurable buttons

### Confirm Dialog

**File:** `src/components/shared/confirm-dialog.tsx`

Reusable confirmation dialog for deletes. Used by all resources.

## Step 6: shadcn Components to Install

```bash
cd web
bunx shadcn add select textarea form popover command
```

## Step 7: Files to Create

| File | Action |
|---|---|
| `src/lib/server/apps.ts` | Create: app server functions |
| `src/lib/server/roles.ts` | Create: role server functions |
| `src/lib/server/permissions.ts` | Create: permission server functions |
| `src/lib/queries/apps.ts` | Create: app TanStack Query hooks + query keys |
| `src/lib/queries/roles.ts` | Create: role TanStack Query hooks + query keys |
| `src/lib/queries/permissions.ts` | Create: permission TanStack Query hooks + query keys |
| `src/lib/schemas/app.ts` | Create: app Zod schema |
| `src/lib/schemas/role.ts` | Create: role Zod schema |
| `src/lib/schemas/permission.ts` | Create: permission Zod schema |
| `src/routes/_authenticated/apps/index.tsx` | Create: app list page |
| `src/routes/_authenticated/apps/$appId.tsx` | Create: app detail page |
| `src/routes/_authenticated/roles/index.tsx` | Create: role list page |
| `src/routes/_authenticated/roles/$roleId.tsx` | Create: role detail page |
| `src/routes/_authenticated/permissions/index.tsx` | Create: permission list page |
| `src/routes/_authenticated/permissions/$permissionId.tsx` | Create: permission detail page |
| `src/components/apps/create-app-dialog.tsx` | Create |
| `src/components/apps/edit-app-dialog.tsx` | Create |
| `src/components/apps/delete-app-dialog.tsx` | Create |
| `src/components/roles/create-role-dialog.tsx` | Create |
| `src/components/roles/edit-role-dialog.tsx` | Create |
| `src/components/roles/delete-role-dialog.tsx` | Create |
| `src/components/roles/assign-permissions-dialog.tsx` | Create |
| `src/components/permissions/create-permission-dialog.tsx` | Create |
| `src/components/permissions/edit-permission-dialog.tsx` | Create |
| `src/components/permissions/delete-permission-dialog.tsx` | Create |
| `src/components/shared/data-table.tsx` | Create: reusable table component |
| `src/components/shared/confirm-dialog.tsx` | Create: reusable confirmation dialog |

## Completion Criteria

- [ ] App list/detail/CRUD works end-to-end via server functions
- [ ] Role list/detail/CRUD works end-to-end via server functions
- [ ] Permission list/detail/CRUD works end-to-end via server functions
- [ ] Permission assignment to roles works via dialog → `assignRolePermissionsFn`
- [ ] All create/edit/delete buttons gated by permission hooks
- [ ] Reusable data-table component used across all list pages
- [ ] TanStack Query cache invalidation works on all mutations
- [ ] All forms use TanStack Form + Zod validation
- [ ] `bun run typecheck && bun run lint` passes
