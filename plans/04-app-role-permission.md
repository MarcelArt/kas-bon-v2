# Phase 4 — Apps, Roles & Permissions

CRUD for apps, roles, and permissions, plus permission assignment to roles.

## 1. App Management

### Server Functions

**File:** `src/lib/apps.server.ts`

- `listAppsFn(opts)` — `GET /v1/apps`
- `getAppFn(id)` — `GET /v1/apps/{id}`
- `createAppFn(body: AppInput)` — `POST /v1/apps`
- `updateAppFn(id, body)` — `PUT /v1/apps/{id}`
- `deleteAppFn(id)` — `DELETE /v1/apps/{id}`

### App List Page

**File:** `src/routes/_authenticated/apps/index.tsx`

Simple table: Name, Description, Created At, Actions.
- Create: `apps#create`
- Edit: `apps#update`
- Delete: `apps#delete`

### App Detail Page

**File:** `src/routes/_authenticated/apps/$appId.tsx`

Display name, description, timestamps. Edit form and delete button.

### App CRUD Dialogs

- `src/components/apps/create-app-dialog.tsx`
- `src/components/apps/edit-app-dialog.tsx`
- `src/components/apps/delete-app-dialog.tsx`

---

## 2. Role Management

### Server Functions

**File:** `src/lib/roles.server.ts`

- `listRolesFn(opts)` — `GET /v1/roles`
- `getRoleFn(id)` — `GET /v1/roles/{id}`
- `createRoleFn(body: RoleInput)` — `POST /v1/roles`
- `updateRoleFn(id, body)` — `PUT /v1/roles/{id}`
- `deleteRoleFn(id)` — `DELETE /v1/roles/{id}`
- `getRolePermissionsFn(id)` — `GET /v1/roles/{id}/permissions`
- `assignRolePermissionsFn(id, permissionIds: number[])` — `PATCH /v1/roles/{id}/permissions`

### Role List Page

**File:** `src/routes/_authenticated/roles/index.tsx`

Table: Name, Description, Domain, Actions.
- Create: `roles#create`
- Edit: `roles#update`
- Delete: `roles#delete`

### Role Detail Page

**File:** `src/routes/_authenticated/roles/$roleId.tsx`

Three sections:

**Info Section**
- Name, description, domain (linked to domain detail)
- Edit form — `roles#update`
- Delete button — `roles#delete`

**Permissions Section**
- List of assigned permissions (from `getRolePermissionsFn`)
- "Assign Permissions" button — `roles#update`
- Assignment dialog shows all available permissions as checkboxes
- On submit, calls `assignRolePermissionsFn`

### Role CRUD Dialogs

- `src/components/roles/create-role-dialog.tsx`
- `src/components/roles/edit-role-dialog.tsx`
- `src/components/roles/delete-role-dialog.tsx`
- `src/components/roles/assign-permissions-dialog.tsx`

---

## 3. Permission Management

### Server Functions

**File:** `src/lib/permissions.server.ts`

- `listPermissionsFn(opts)` — `GET /v1/permissions`
- `getPermissionFn(id)` — `GET /v1/permissions/{id}`
- `createPermissionFn(body: PermissionInput)` — `POST /v1/permissions`
- `updatePermissionFn(id, body)` — `PUT /v1/permissions/{id}`
- `deletePermissionFn(id)` — `DELETE /v1/permissions/{id}`

### Permission List Page

**File:** `src/routes/_authenticated/permissions/index.tsx`

Table: Name (format: `resource#action`), Description, App, Actions.
- Create: `permissions#create`
- Edit: `permissions#update`
- Delete: `permissions#delete`

### Permission Detail Page

**File:** `src/routes/_authenticated/permissions/$permissionId.tsx`

Display name, description, app (linked). Edit form and delete button.

### Permission CRUD Dialogs

- `src/components/permissions/create-permission-dialog.tsx`
- `src/components/permissions/edit-permission-dialog.tsx`
- `src/components/permissions/delete-permission-dialog.tsx`

### Permission Naming Convention

Permission names follow `resource#action` format. When creating/editing:
- Resource field: dropdown or text input (e.g., `users`, `roles`, `domains`, `apps`, `permissions`, `all`)
- Action field: dropdown or text input (e.g., `read`, `create`, `update`, `delete`, `fullAccess`)
- Auto-combined into `resource#action` before sending to API

---

## 4. Assign Permissions to Role Dialog

**File:** `src/components/roles/assign-permissions-dialog.tsx`

- Fetches all permissions via `listPermissionsFn`
- Fetches currently assigned permissions via `getRolePermissionsFn`
- Displays permissions grouped by resource (parse `resource#action` name)
- Checkboxes for each permission, pre-checked if assigned
- On submit, sends array of selected permission IDs via `assignRolePermissionsFn`

## 5. Reusable CRUD Pattern

These three resources (Apps, Roles, Permissions) share the same CRUD pattern. Extract reusable components:

**File:** `src/components/shared/data-table.tsx`

Generic table component accepting:
- `columns: ColumnDef[]`
- `data: T[]`
- `pagination: PaginationState`
- `onPaginationChange: (state) => void`
- Actions column with configurable buttons

**File:** `src/components/shared/confirm-dialog.tsx`

Reusable confirmation dialog for deletes. Used by all resources.

## 6. shadcn Components to Install

```bash
cd web
bunx shadcn add select textarea form popover command
```

## 7. Files to Create

| File | Action |
|---|---|
| `src/lib/apps.server.ts` | Create: app server functions |
| `src/lib/roles.server.ts` | Create: role server functions |
| `src/lib/permissions.server.ts` | Create: permission server functions |
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
- [ ] All create/edit/delete buttons gated by permissions
- [ ] Reusable data-table component used across all list pages
- [ ] `bun run typecheck && bun run lint` passes
