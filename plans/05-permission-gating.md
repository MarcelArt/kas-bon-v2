# Phase 5 — Permission Gating

Frontend-wide permission system that controls menu visibility, page access, and action availability based on user permissions from the backend.

## 1. Permission Data Flow

```
Login → auth context stores user
  → fetch GET /v1/users/{userId}/permissions (with X-App-Id, X-Domain-Id)
  → parse permission tuples: [[sub, app, dom, res, act], ...]
  → extract "resource#action" strings → Set<string>
  → store in auth context

Every render:
  hasPermission("users", "read") → check Set or isSuperUser
```

## 2. Permission Parsing Utility

**File:** `src/lib/permissions.ts`

```typescript
interface PermissionTuple {
  sub: string
  app: string
  dom: string
  resource: string
  action: string
}

function parsePermissionTuples(tuples: string[][]): PermissionTuple[]
function buildPermissionSet(tuples: string[][]): Set<string>
function isSuperUser(permissions: Set<string>): boolean
function checkPermission(permissions: Set<string>, resource: string, action: string): boolean

// Resource-action constants
const RESOURCES = {
  USERS: "users",
  DOMAINS: "domains",
  APPS: "apps",
  ROLES: "roles",
  PERMISSIONS: "permissions",
  ALL: "all",
} as const

const ACTIONS = {
  READ: "read",
  CREATE: "create",
  UPDATE: "update",
  DELETE: "delete",
  FULL_ACCESS: "fullAccess",
} as const
```

## 3. Sidebar Menu Gating

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

Rendering logic:
- If `resource` is null → always show (Dashboard)
- If user is superUser (`all#fullAccess`) → always show
- Otherwise → show only if `hasPermission(resource, action)`

## 4. Route-Level Permission Guard

**File:** Update `src/lib/auth-guard.ts`

Add per-route permission checking. Each route declares its required permission via route context:

```typescript
// In route definition (e.g., users/index.tsx):
export const Route = createFileRoute("/_authenticated/users/")({
  beforeLoad: ({ context }) => {
    if (!context.auth.hasPermission("users", "read")) {
      throw redirect({ to: "/dashboard" })
    }
  },
  loader: async () => { ... },
})
```

Alternative: create a helper that generates the `beforeLoad`:

```typescript
function requirePermission(resource: string, action: string) {
  return ({ context }: { context: RouteContext }) => {
    if (!context.auth.hasPermission(resource, action)) {
      throw redirect({ to: "/dashboard" })
    }
  }
}

// Usage:
export const Route = createFileRoute("/_authenticated/users/")({
  beforeLoad: requirePermission("users", "read"),
})
```

## 5. Action-Level Permission Gating

All CRUD action buttons already wrapped with `hasPermission` checks from Phases 2-4. Verify consistency:

| Action | Permission | Button |
|---|---|---|
| Create user | `users#create` | "Create User" button on user list |
| Edit user | `users#update` | "Edit" button on user row + detail |
| Delete user | `users#delete` | "Delete" button on user row + detail |
| Assign roles | `users#update` | "Assign Roles" button on user detail |
| Create domain | `domains#create` | "Create Domain" button |
| Edit domain | `domains#update` | "Edit" button |
| Delete domain | `domains#delete` | "Delete" button |
| Create app | `apps#create` | "Create App" button |
| Edit app | `apps#update` | "Edit" button |
| Delete app | `apps#delete` | "Delete" button |
| Create role | `roles#create` | "Create Role" button |
| Edit role | `roles#update` | "Edit" button |
| Delete role | `roles#delete` | "Delete" button |
| Assign permissions | `roles#update` | "Assign Permissions" button on role detail |
| Create permission | `permissions#create` | "Create Permission" button |
| Edit permission | `permissions#update` | "Edit" button |
| Delete permission | `permissions#delete` | "Delete" button |

## 6. Permission Hook

**File:** `src/hooks/use-permission.ts`

Convenience hook for components:

```typescript
function usePermission(resource: string, action: string): boolean
function useCanCreate(resource: string): boolean  // action = "create"
function useCanEdit(resource: string): boolean    // action = "update"
function useCanDelete(resource: string): boolean  // action = "delete"
function useIsSuperUser(): boolean
```

These read from auth context and can be used in any component:

```tsx
const canCreate = useCanCreate("users")
return canCreate ? <Button>Create User</Button> : null
```

## 7. Unauthorized Page

**File:** `src/routes/_authenticated/unauthorized.tsx`

Displayed when a user navigates to a page they don't have access to (fallback if redirect to dashboard isn't desired). Shows a message: "You don't have permission to view this page."

## 8. Permission Context Refresh

When permissions change (e.g., after role/permission mutation), refresh the permission set:

- After `assignUserRolesFn` → re-fetch user permissions
- After `assignRolePermissionsFn` → re-fetch user permissions (if current user is affected)
- Provide a `refreshPermissions()` function in auth context

## 9. Files to Create/Modify

| File | Action |
|---|---|
| `src/lib/permissions.ts` | Create: permission parsing + constants |
| `src/hooks/use-permission.ts` | Create: permission hooks |
| `src/components/layout/app-shell.tsx` | Modify: gate sidebar items by permission |
| `src/lib/auth.tsx` | Modify: integrate permission parsing, add refreshPermissions |
| `src/lib/auth-guard.ts` | Modify: add requirePermission helper |
| `src/routes/_authenticated/unauthorized.tsx` | Create: unauthorized page |
| All route files in `_authenticated/` | Modify: add beforeLoad permission guards |

## Completion Criteria

- [ ] Sidebar menu items hidden when user lacks `resource#read` permission
- [ ] Direct URL navigation to unauthorized pages redirects or shows unauthorized
- [ ] All action buttons hidden when user lacks specific permission
- [ ] SuperUser (`all#fullAccess`) sees everything
- [ ] Permission hooks work in all components
- [ ] Permissions refresh after role/permission mutations
- [ ] `bun run typecheck && bun run lint` passes
