# Phase 4 — Apps, Roles & Permissions

CRUD for apps, roles, and permissions, plus permission assignment to roles. Uses Axios API functions, TanStack Query hooks, TanStack Form + Zod.

## Step 1: App Management

### API Functions

**File:** `src/lib/api/apps.ts`

```typescript
import { api, unwrap } from "@/lib/api"
import type { App, JSONResponse, PaginatedResponse } from "@/lib/api.types"

export const appApi = {
  list: (params?: { page?: number; size?: number; sort?: string; filters?: string }) =>
    api.get<JSONResponse<PaginatedResponse<App>>>("/v1/apps", { params }).then(unwrap),

  get: (id: number) =>
    api.get<JSONResponse<App>>(`/v1/apps/${id}`).then(unwrap),

  create: (body: { name: string; description: string }) =>
    api.post<JSONResponse<number>>("/v1/apps", body).then(unwrap),

  update: (id: number, body: { name: string; description: string }) =>
    api.put<JSONResponse<null>>(`/v1/apps/${id}`, body).then(unwrap),

  delete: (id: number) =>
    api.delete<JSONResponse<null>>(`/v1/apps/${id}`).then(unwrap),
}
```

### Query Keys & Hooks

**File:** `src/lib/queries/apps.ts`

```typescript
export const appKeys = {
  all: ["apps"] as const,
  lists: () => [...appKeys.all, "list"] as const,
  list: (filters: object) => [...appKeys.lists(), filters] as const,
  details: () => [...appKeys.all, "detail"] as const,
  detail: (id: number) => [...appKeys.details(), id] as const,
}

export function useApps(filters?: object) { ... }
export function useApp(id: number) { ... }
export function useCreateApp() { ... }
export function useUpdateApp() { ... }
export function useDeleteApp() { ... }
```

### Zod Schema

**File:** `src/lib/schemas/app.ts`

```typescript
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

### API Functions

**File:** `src/lib/api/roles.ts`

```typescript
export const roleApi = {
  list: (params?) => ...,
  get: (id) => ...,
  create: (body: { name: string; description: string; domainId: number }) => ...,
  update: (id, body) => ...,
  delete: (id) => ...,
  getPermissions: (id: number) =>
    api.get<JSONResponse<string[][]>>(`/v1/roles/${id}/permissions`).then(unwrap),
  assignPermissions: (id: number, permissionIds: number[]) =>
    api.patch<JSONResponse<string[]>>(`/v1/roles/${id}/permissions`, permissionIds).then(unwrap),
}
```

### Query Keys & Hooks

**File:** `src/lib/queries/roles.ts`

```typescript
export const roleKeys = {
  all: ["roles"] as const,
  lists: () => [...roleKeys.all, "list"] as const,
  list: (filters: object) => [...roleKeys.lists(), filters] as const,
  details: () => [...roleKeys.all, "detail"] as const,
  detail: (id: number) => [...roleKeys.details(), id] as const,
  permissions: (id: number) => [...roleKeys.detail(id), "permissions"] as const,
}

export function useRoles(filters?: object) { ... }
export function useRole(id: number) { ... }
export function useRolePermissions(id: number) { ... }
export function useCreateRole() { ... }
export function useUpdateRole() { ... }
export function useDeleteRole() { ... }
export function useAssignRolePermissions() { ... }  // invalidates role permissions cache
```

### Zod Schema

**File:** `src/lib/schemas/role.ts`

```typescript
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

### API Functions

**File:** `src/lib/api/permissions.ts`

```typescript
export const permissionApi = {
  list: (params?) => ...,
  get: (id) => ...,
  create: (body: { name: string; description: string; appId: number }) => ...,
  update: (id, body) => ...,
  delete: (id) => ...,
}
```

### Query Keys & Hooks

**File:** `src/lib/queries/permissions.ts`

```typescript
export const permissionKeys = {
  all: ["permissions"] as const,
  lists: () => [...permissionKeys.all, "list"] as const,
  list: (filters: object) => [...permissionKeys.lists(), filters] as const,
  details: () => [...permissionKeys.all, "detail"] as const,
  detail: (id: number) => [...permissionKeys.details(), id] as const,
}

export function usePermissions(filters?: object) { ... }
export function usePermission(id: number) { ... }
export function useCreatePermission() { ... }
export function useUpdatePermission() { ... }
export function useDeletePermission() { ... }
```

### Zod Schema

**File:** `src/lib/schemas/permission.ts`

```typescript
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
| `src/lib/api/apps.ts` | Create: app Axios API functions |
| `src/lib/api/roles.ts` | Create: role Axios API functions |
| `src/lib/api/permissions.ts` | Create: permission Axios API functions |
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

- [ ] App list/detail/CRUD works end-to-end
- [ ] Role list/detail/CRUD works end-to-end
- [ ] Permission list/detail/CRUD works end-to-end
- [ ] Permission assignment to roles works via dialog
- [ ] All create/edit/delete buttons gated by permission hooks
- [ ] Reusable data-table component used across all list pages
- [ ] TanStack Query cache invalidation works on all mutations
- [ ] All forms use TanStack Form + Zod validation
- [ ] `bun run typecheck && bun run lint` passes
