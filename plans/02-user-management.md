# Phase 2 — User Management

User list, detail, CRUD, and role assignment. Uses TanStack Start server functions, TanStack Query hooks, TanStack Form + Zod for forms.

## Step 1: User Server Functions

**File:** `src/lib/server/users.ts`

Server functions for all user endpoints (see `.opencode/API.md` — Users section):

```typescript
import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { User, PaginatedResponse } from "@/lib/api.types"

const authContextSchema = z.object({
  accessToken: z.string(),
  domainId: z.number().optional(),
  appId: z.number().optional(),
})

export const getUsersFn = createServerFn()
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
    return serverApi<PaginatedResponse<User>>("GET", "/users", { accessToken, domainId, appId }, undefined, params)
  })

export const getUserFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<User>("GET", `/users/${id}`, { accessToken, domainId, appId })
  })

export const createUserFn = createServerFn()
  .validator(
    authContextSchema.extend({
      username: z.string(),
      email: z.string().email(),
      password: z.string().min(6),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...body } = data
    return serverApi<number>("POST", "/users", { accessToken, domainId, appId }, body)
  })

export const updateUserFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      username: z.string(),
      email: z.string().email(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id, ...body } = data
    return serverApi<null>("PUT", `/users/${id}`, { accessToken, domainId, appId }, body)
  })

export const deleteUserFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<null>("DELETE", `/users/${id}`, { accessToken, domainId, appId })
  })

export const getUserRolesFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<string[]>("GET", `/users/${id}/roles`, { accessToken, domainId, appId })
  })

export const assignUserRolesFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      roleIds: z.array(z.number()),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id, roleIds } = data
    return serverApi<string[]>("PATCH", `/users/${id}/roles`, { accessToken, domainId, appId }, roleIds)
  })

export const getUserPermissionsFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<string[][]>("GET", `/users/${id}/permissions`, { accessToken, domainId, appId })
  })
```

## Step 2: TanStack Query Hooks & Query Keys

**File:** `src/lib/queries/users.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  getUsersFn,
  getUserFn,
  createUserFn,
  updateUserFn,
  deleteUserFn,
  getUserRolesFn,
  assignUserRolesFn,
  getUserPermissionsFn,
} from "@/lib/server/users"
import { useAuthContext, withAuthRetry } from "./auth-context"

export const userKeys = {
  all: ["users"] as const,
  lists: () => [...userKeys.all, "list"] as const,
  list: (filters: object) => [...userKeys.lists(), filters] as const,
  details: () => [...userKeys.all, "detail"] as const,
  detail: (id: number) => [...userKeys.details(), id] as const,
  roles: (id: number) => [...userKeys.detail(id), "roles"] as const,
  permissions: (id: number) => [...userKeys.detail(id), "permissions"] as const,
}

export function useUsers(filters: { page?: number; size?: number; sort?: string; filters?: string }) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: userKeys.list(filters),
    queryFn: withAuthRetry(() => getUsersFn({ data: { ...auth, ...filters } })),
    enabled: !!auth.accessToken,
  })
}

export function useUser(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: withAuthRetry(() => getUserFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useUserRoles(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: userKeys.roles(id),
    queryFn: withAuthRetry(() => getUserRolesFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useUserPermissions(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: userKeys.permissions(id),
    queryFn: withAuthRetry(() => getUserPermissionsFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useCreateUser() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (body: { username: string; email: string; password: string }) =>
      createUserFn({ data: { ...auth, ...body } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: userKeys.all }),
  })
}

export function useUpdateUser() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: ({ id, ...body }: { id: number; username: string; email: string }) =>
      updateUserFn({ data: { ...auth, id, ...body } }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: userKeys.detail(id) })
      qc.invalidateQueries({ queryKey: userKeys.lists() })
    },
  })
}

export function useDeleteUser() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (id: number) => deleteUserFn({ data: { ...auth, id } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: userKeys.all }),
  })
}

export function useAssignUserRoles() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: ({ userId, roleIds }: { userId: number; roleIds: number[] }) =>
      assignUserRolesFn({ data: { ...auth, id: userId, roleIds } }),
    onSuccess: (_, { userId }) => {
      qc.invalidateQueries({ queryKey: userKeys.roles(userId) })
      qc.invalidateQueries({ queryKey: userKeys.permissions(userId) })
    },
  })
}
```

## Step 3: Zod Schemas

**File:** `src/lib/schemas/user.ts`

```typescript
import { z } from "zod"

export const createUserSchema = z.object({
  username: z.string().min(1, "Username is required"),
  email: z.string().email("Invalid email"),
  password: z.string().min(6, "Password must be at least 6 characters"),
})

export const editUserSchema = z.object({
  username: z.string().min(1, "Username is required"),
  email: z.string().email("Invalid email"),
})

export type CreateUserFormData = z.infer<typeof createUserSchema>
export type EditUserFormData = z.infer<typeof editUserSchema>
```

## Step 4: User List Page

**File:** `src/routes/_authenticated/users/index.tsx`

- Use `useUsers(filters)` TanStack Query hook for data
- Table with columns: Username, Email, Created At, Actions
- Actions column: View (link to detail), Edit, Delete — gated by permissions:
  - Edit button: `useCanEdit("users")`
  - Delete button: `useCanDelete("users")`
- "Create User" button at top — gated by `useCanCreate("users")`
- Pagination controls at bottom
- Search/filter input (maps to `filters` query param)

**shadcn components needed:** `table`, `dialog`, `dropdown-menu`, `pagination`

## Step 5: User Detail Page

**File:** `src/routes/_authenticated/users/$userId.tsx`

Use `useUser(id)`, `useUserRoles(id)`, `useUserPermissions(id)` TanStack Query hooks.

### Profile Section
- Display: username, email, created/updated timestamps
- Edit form (inline or modal using TanStack Form + `editUserSchema`) — gated by `useCanEdit("users")`
- Delete button with confirmation — gated by `useCanDelete("users")`

### Roles Section
- List of assigned roles (from `useUserRoles`)
- "Assign Roles" button opens a dialog — gated by `useCanEdit("users")`
- Dialog shows all available roles as checkboxes
- On submit, calls `useAssignUserRoles()` mutation

### Permissions Section (read-only)
- List of effective permissions (from `useUserPermissions`)
- Displayed as tags/badges: `resource#action`

**shadcn components needed:** `badge`, `tabs`

## Step 6: Create User Dialog

**File:** `src/components/users/create-user-dialog.tsx`

- Modal form using TanStack Form with `zodValidator` and `createUserSchema`
- Fields: username, email, password
- Calls `useCreateUser()` mutation
- On success: `sonner` toast, query cache auto-invalidates
- On error: display field-level errors from Zod

## Step 7: Edit User Dialog

**File:** `src/components/users/edit-user-dialog.tsx`

- Modal form using TanStack Form with `zodValidator` and `editUserSchema`
- Pre-filled with current user data via `defaultValues`
- Calls `useUpdateUser()` mutation
- On success: toast, query cache auto-invalidates

## Step 8: Delete Confirmation Dialog

**File:** `src/components/users/delete-user-dialog.tsx`

- Reusable confirmation dialog
- Shows "Are you sure you want to delete user X?"
- Calls `useDeleteUser()` mutation
- On success: toast, navigate to user list

## Step 9: Role Assignment Dialog

**File:** `src/components/users/assign-roles-dialog.tsx`

- Fetches all roles (filtered by current domain) via `useRoles()` hook from Phase 4
- Displays roles as a checkbox list, pre-checking currently assigned roles
- On submit, calls `useAssignUserRoles()` mutation
- On success: toast, query cache auto-invalidates

## Step 10: Permission Gating Pattern

Each page/action uses permission hooks from `src/hooks/use-permission.ts`:

```typescript
const canCreate = useCanCreate("users")
const canEdit = useCanEdit("users")
const canDelete = useCanDelete("users")
```

## Step 11: shadcn Components to Install

```bash
cd web
bunx shadcn add table dialog dropdown-menu badge tabs alert-dialog
```

## Step 12: Files to Create

| File | Action |
|---|---|
| `src/lib/server/users.ts` | Create: user server functions |
| `src/lib/queries/users.ts` | Create: user TanStack Query hooks + query keys |
| `src/lib/schemas/user.ts` | Create: user Zod schemas |
| `src/routes/_authenticated/users/index.tsx` | Create: user list page |
| `src/routes/_authenticated/users/$userId.tsx` | Create: user detail page |
| `src/components/users/create-user-dialog.tsx` | Create: create user dialog (TanStack Form) |
| `src/components/users/edit-user-dialog.tsx` | Create: edit user dialog (TanStack Form) |
| `src/components/users/delete-user-dialog.tsx` | Create: delete confirmation |
| `src/components/users/assign-roles-dialog.tsx` | Create: role assignment dialog |

## Completion Criteria

- [ ] User list loads with pagination via TanStack Query → server functions
- [ ] Create user dialog works with TanStack Form + Zod → calls `createUserFn`
- [ ] Edit user dialog works with TanStack Form + Zod → calls `updateUserFn`
- [ ] Delete with confirmation works (mutation + cache invalidation → `deleteUserFn`)
- [ ] Role assignment dialog shows roles and saves correctly via `assignUserRolesFn`
- [ ] User permissions displayed as badges via `getUserPermissionsFn`
- [ ] Create/Edit/Delete buttons hidden when user lacks permission
- [ ] `bun run typecheck && bun run lint` passes
