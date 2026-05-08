# Phase 2 — User Management

User list, detail, CRUD, and role assignment.

## 1. Server Functions

**File:** `src/lib/users.server.ts`

Using `createServerFn`:

- `listUsersFn(opts: { page, size, sort, filters })` — `GET /v1/users`
- `getUserFn(id)` — `GET /v1/users/{id}`
- `createUserFn(body: UserInput)` — `POST /v1/users`
- `updateUserFn(id, body: Partial<User>)` — `PUT /v1/users/{id}`
- `deleteUserFn(id)` — `DELETE /v1/users/{id}`
- `getUserRolesFn(id)` — `GET /v1/users/{id}/roles`
- `assignUserRolesFn(id, roleIds: number[])` — `PATCH /v1/users/{id}/roles`
- `getUserPermissionsFn(id)` — `GET /v1/users/{id}/permissions`

All functions attach `X-App-Id` and `X-Domain-Id` headers from context.

## 2. User List Page

**File:** `src/routes/_authenticated/users/index.tsx`

- Route loader fetches paginated user list via `listUsersFn`
- Table with columns: Username, Email, Created At, Actions
- Actions column: View (link to detail), Edit, Delete — gated by permissions:
  - Edit button: `users#update`
  - Delete button: `users#delete`
- "Create User" button at top — gated by `users#create`
- Pagination controls at bottom
- Search/filter input (maps to `filters` query param)

**shadcn components needed:** `table`, `dialog`, `dropdown-menu`, `pagination`

## 3. User Detail Page

**File:** `src/routes/_authenticated/users/$userId.tsx`

Route loader fetches user via `getUserFn`.

Two sections:

### Profile Section
- Display: username, email, created/updated timestamps
- Edit form (inline or modal) — gated by `users#update`
- Delete button with confirmation — gated by `users#delete`

### Roles Section
- List of assigned roles (from `getUserRolesFn`)
- "Assign Roles" button opens a dialog — gated by `users#update`
- Dialog shows all available roles (from `listRolesFn`) as checkboxes
- On submit, calls `assignUserRolesFn`

### Permissions Section (read-only)
- List of effective permissions (from `getUserPermissionsFn`)
- Displayed as tags/badges: `resource#action`

**shadcn components needed:** `badge`, `tabs`

## 4. Create User Dialog

**File:** `src/components/users/create-user-dialog.tsx`

- Modal form with fields: username, email, password
- Calls `createUserFn`
- On success: toast notification, refresh user list
- On error: show field-level errors

## 5. Edit User Dialog

**File:** `src/components/users/edit-user-dialog.tsx`

- Modal form pre-filled with current user data
- Fields: username, email (password change is separate or omitted)
- Calls `updateUserFn`
- On success: toast, refresh data

## 6. Delete Confirmation Dialog

**File:** `src/components/users/delete-user-dialog.tsx`

- Reusable confirmation dialog
- Shows "Are you sure you want to delete user X?"
- Calls `deleteUserFn`
- On success: toast, navigate to user list

## 7. Role Assignment Dialog

**File:** `src/components/users/assign-roles-dialog.tsx`

- Fetches all roles (filtered by current domain) via `listRolesFn`
- Displays roles as a checkbox list, pre-checking currently assigned roles
- On submit, calls `assignUserRolesFn` with selected role IDs
- On success: toast, refresh roles section

## 8. Permission Gating Pattern

Each page/action checks permissions from auth context:

```typescript
const { hasPermission } = useAuth()

// In JSX:
{hasPermission("users", "create") && <Button>Create User</Button>}
{hasPermission("users", "update") && <Button>Edit</Button>}
{hasPermission("users", "delete") && <Button>Delete</Button>}
```

## 9. shadcn Components to Install

```bash
cd web
bunx shadcn add table dialog dropdown-menu badge tabs alert-dialog
```

## 10. Files to Create

| File | Action |
|---|---|
| `src/lib/users.server.ts` | Create: user server functions |
| `src/routes/_authenticated/users/index.tsx` | Create: user list page |
| `src/routes/_authenticated/users/$userId.tsx` | Create: user detail page |
| `src/components/users/create-user-dialog.tsx` | Create: create user dialog |
| `src/components/users/edit-user-dialog.tsx` | Create: edit user dialog |
| `src/components/users/delete-user-dialog.tsx` | Create: delete confirmation |
| `src/components/users/assign-roles-dialog.tsx` | Create: role assignment dialog |

## Completion Criteria

- [ ] User list loads with pagination
- [ ] Create user dialog works
- [ ] Edit user dialog works
- [ ] Delete user with confirmation works
- [ ] Role assignment dialog shows roles and saves correctly
- [ ] User permissions displayed as badges
- [ ] Create/Edit/Delete buttons hidden when user lacks permission
- [ ] `bun run typecheck && bun run lint` passes
