# Phase 5 — Permission Gating

Frontend-wide permission system that controls menu visibility, page access, and action availability based on user permissions from the backend. Integrates with Zustand auth store and TanStack Query.

## Step 1: Permission Data Flow

```
Login → Zustand store: user + tokens set
  → GET /v1/users/{userId}/organizations (after org selection)
  → Axios automatically attaches X-App-Id, X-Domain-Id
  → fetch GET /v1/users/{userId}/permissions
  → parse permission tuples: [[sub, app, dom, res, act], ...]
  → extract "resource#action" strings → Set<string>
  → store in Zustand auth store (setPermissions)

Every render:
  usePermission("users", "read") → reads from Zustand store → check Set or isSuperUser
```

## Step 2: Permission Parsing Utility

**File:** `src/lib/permissions.ts`

Already created in Phase 1. Verify it contains:

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

## Step 3: Permission Hooks

**File:** `src/hooks/use-permission.ts`

Already created in Phase 1. Uses Zustand store selectors:

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

## Step 4: Sidebar Menu Gating

**File:** Update `src/components/layout/app-shell.tsx`

Each sidebar navigation item declares its required permission:

```typescript
const navItems = [
  { to: "/dashboard", label: "Dashboard", icon: House, resource: null },      // always visible
  { to: "/users", label: "Users", icon: Users, resource: "users", action: "read" },
  { to: "/domains", label: "Domains", icon: Folders, resource: "domains", action: "read" },
  { to: "/apps", label: "Apps", icon: AppWindow, resource: "apps", action: "read" },
  { to: "/roles", label: "Roles", icon: Shield, resource: "roles", action: "read" },
  { to: "/permissions", label: "Permissions", icon: Key, resource: "permissions", action: "read" },
]
```

Rendering logic — use `usePermission` hook:
- If `resource` is null → always show (Dashboard)
- Otherwise → show only if `usePermission(resource, action)` is true (superuser check is inside)

## Step 5: Route-Level Permission Guard

**File:** Update `src/lib/auth-guard.ts`

Create a helper that generates `beforeLoad` for permission checking:

```typescript
import { redirect } from "@tanstack/react-router"
import { useAuthStore } from "@/lib/stores/auth-store"

export function requirePermission(resource: string, action: string) {
  return () => {
    const { hasPermission } = useAuthStore.getState()
    if (!hasPermission(resource, action)) {
      throw redirect({ to: "/dashboard" })
    }
  }
}
```

Usage in route definitions:
```typescript
// e.g., users/index.tsx
export const Route = createFileRoute("/_authenticated/users/")({
  beforeLoad: requirePermission("users", "read"),
  // loader can use TanStack Query or server fn for SSR data
})
```

## Step 6: Action-Level Permission Gating

All CRUD action buttons already wrapped with permission hooks from Phases 2-4. Verify consistency:

| Action | Permission | Hook |
|---|---|---|
| Create user | `users#create` | `useCanCreate("users")` |
| Edit user | `users#update` | `useCanEdit("users")` |
| Delete user | `users#delete` | `useCanDelete("users")` |
| Assign roles | `users#update` | `useCanEdit("users")` |
| Create domain | `domains#create` | `useCanCreate("domains")` |
| Edit domain | `domains#update` | `useCanEdit("domains")` |
| Delete domain | `domains#delete` | `useCanDelete("domains")` |
| Create app | `apps#create` | `useCanCreate("apps")` |
| Edit app | `apps#update` | `useCanEdit("apps")` |
| Delete app | `apps#delete` | `useCanDelete("apps")` |
| Create role | `roles#create` | `useCanCreate("roles")` |
| Edit role | `roles#update` | `useCanEdit("roles")` |
| Delete role | `roles#delete` | `useCanDelete("roles")` |
| Assign permissions | `roles#update` | `useCanEdit("roles")` |
| Create permission | `permissions#create` | `useCanCreate("permissions")` |
| Edit permission | `permissions#update` | `useCanEdit("permissions")` |
| Delete permission | `permissions#delete` | `useCanDelete("permissions")` |

## Step 7: Unauthorized Page

**File:** `src/routes/_authenticated/unauthorized.tsx`

Displayed when a user navigates to a page they don't have access to (fallback if redirect to dashboard isn't desired). Shows: "You don't have permission to view this page."

## Step 8: Permission Context Refresh

When permissions change (e.g., after role/permission mutation), refresh the permission set in Zustand:

- After `useAssignUserRoles()` mutation success → re-fetch permissions via TanStack Query and update Zustand store
- After `useAssignRolePermissions()` mutation success → re-fetch permissions (if current user is affected)
- Add `refreshPermissions()` function that calls the API and updates Zustand:

```typescript
// In a utility or hook
async function refreshPermissions(userId: number) {
  const tuples = await authApi.getPermissions(userId)
  useAuthStore.getState().setPermissions(tuples)
}
```

Call this in mutation `onSuccess` callbacks where role/permission assignments change.

## Step 9: Files to Create/Modify

| File | Action |
|---|---|
| `src/lib/permissions.ts` | Verify: permission parsing + constants (from Phase 1) |
| `src/hooks/use-permission.ts` | Verify: permission hooks (from Phase 1) |
| `src/components/layout/app-shell.tsx` | Modify: gate sidebar items by permission hooks |
| `src/lib/auth-guard.ts` | Modify: add `requirePermission` helper |
| `src/routes/_authenticated/unauthorized.tsx` | Create: unauthorized page |
| All route files in `_authenticated/` | Modify: add `beforeLoad: requirePermission(...)` |
| All mutation hooks with role/permission side effects | Modify: add permission refresh in `onSuccess` |

## Completion Criteria

- [ ] Sidebar menu items hidden when user lacks `resource#read` permission
- [ ] Direct URL navigation to unauthorized pages redirects or shows unauthorized
- [ ] All action buttons hidden when user lacks specific permission
- [ ] SuperUser (`all#fullAccess`) sees everything
- [ ] Permission hooks (`usePermission`, `useCanCreate`, etc.) work in all components
- [ ] Permissions refresh after role/permission mutations (Zustand store updated)
- [ ] `bun run typecheck && bun run lint` passes
