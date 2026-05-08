# Phase 3 — Domain Management

Nested organization tree, domain CRUD, and tree navigation. Uses TanStack Query for data fetching and TanStack Form + Zod for forms.

Domains support nesting via `parentId` — a domain can have a parent domain, forming a tree. Domains also have an `isOrganization` flag.

## 1. Server Functions

**File:** `src/lib/domains.server.ts`

- `listDomainsFn(opts: { page, size, sort, filters })` — `GET /v1/domains`
- `getDomainFn(id)` — `GET /v1/domains/{id}`
- `createDomainFn(body: DomainInput)` — `POST /v1/domains`
- `updateDomainFn(id, body: Partial<Domain>)` — `PUT /v1/domains/{id}`
- `deleteDomainFn(id)` — `DELETE /v1/domains/{id}`

## 2. Domain Query Hooks (TanStack Query)

**File:** `src/lib/domains.query.ts`

- `useDomains(opts)` — `useQuery`, key: `["domains", opts]`
- `useDomain(id)` — `useQuery`, key: `["domains", id]`
- `useCreateDomainMutation()` — `useMutation`, invalidates `["domains"]`
- `useUpdateDomainMutation()` — `useMutation`, invalidates `["domains", id]`
- `useDeleteDomainMutation()` — `useMutation`, invalidates `["domains"]`

## 3. Domain List Page — Tree View

**File:** `src/routes/_authenticated/domains/index.tsx`

Instead of a flat table, render domains as a collapsible tree:

```
├── John's organization (org)
│   ├── Engineering (domain)
│   │   ├── Frontend Team (domain)
│   │   └── Backend Team (domain)
│   └── Marketing (domain)
├── Acme Corp (org)
│   └── ...
```

Implementation approach:
1. Fetch all domains via `listDomainsFn` (no pagination, or fetch all pages)
2. Build a tree structure client-side using `parentId` references
3. Render with a recursive tree component

Each tree node shows:
- Domain name
- Badge: "Organization" if `isOrganization === true`
- Actions: View, Edit, Delete — gated by permissions
  - Edit: `domains#update`
  - Delete: `domains#delete`

"Create Domain" button at top — gated by `domains#create`

**shadcn components needed:** `collapsible`, `tree` (or custom recursive component)

## 3. Domain Tree Component

**File:** `src/components/domains/domain-tree.tsx`

Recursive component:
```
DomainTree
  ├── DomainTreeNode
  │   ├── Collapsible trigger: domain name + badges
  │   ├── Collapsible content:
  │   │   ├── DomainTreeNode (for each child)
  │   │   └── ...
  │   └── Action buttons (view, edit, delete)
```

Helper function to build tree:
```typescript
function buildDomainTree(domains: Domain[]): DomainTreeNode[] {
  // Group children by parentId
  // Return root nodes (parentId === null)
}
```

## 4. Domain Detail Page

**File:** `src/routes/_authenticated/domains/$domainId.tsx`

Route loader fetches domain via `getDomainFn`.

Display sections:
- **Info**: name, description, isOrganization flag, parent domain link
- **Children**: list of child domains (links to their detail pages)
- **Breadcrumb**: path from root to current domain for navigation

Edit form — gated by `domains#update`

Delete button with confirmation — gated by `domains#delete`

## 5. Create Domain Dialog

**File:** `src/components/domains/create-domain-dialog.tsx`

Form fields:
- Name (required)
- Description
- Is Organization (checkbox/toggle)
- Parent Domain (select dropdown — list all domains for `parentId`)

Fetch domain list for parent dropdown via `listDomainsFn`.

Calls `createDomainFn` on submit.

## 6. Edit Domain Dialog

**File:** `src/components/domains/edit-domain-dialog.tsx`

Same fields as create, pre-filled with current domain data.

## 7. Delete Domain Dialog

**File:** `src/components/domains/delete-domain-dialog.tsx`

Confirmation dialog. Warn if domain has children.

## 8. Breadcrumb Navigation

**File:** `src/components/domains/domain-breadcrumb.tsx`

Given a domain, walk up the parent chain to build breadcrumb:
```
Root Organization > Engineering > Frontend Team
```

Each segment is a link to the domain detail page.

## 9. shadcn Components to Install

```bash
cd web
bunx shadcn add collapsible breadcrumb select switch
```

## 10. Files to Create

| File | Action |
|---|---|
| `src/lib/domains.server.ts` | Create: domain server functions |
| `src/lib/domains.query.ts` | Create: TanStack Query hooks for domains |
| `src/routes/_authenticated/domains/index.tsx` | Create: domain tree page |
| `src/routes/_authenticated/domains/$domainId.tsx` | Create: domain detail page |
| `src/components/domains/domain-tree.tsx` | Create: recursive tree component |
| `src/components/domains/domain-treeNode.tsx` | Create: tree node component |
| `src/components/domains/create-domain-dialog.tsx` | Create: create domain dialog |
| `src/components/domains/edit-domain-dialog.tsx` | Create: edit domain dialog |
| `src/components/domains/delete-domain-dialog.tsx` | Create: delete confirmation |
| `src/components/domains/domain-breadcrumb.tsx` | Create: breadcrumb navigation |
| `src/lib/build-tree.ts` | Create: generic tree builder from flat list with parentId |

## Completion Criteria

- [ ] Domain tree renders with expand/collapse
- [ ] Nested domains display correctly
- [ ] Create domain dialog with parent selector works
- [ ] Edit domain works
- [ ] Delete domain with confirmation works
- [ ] Breadcrumb navigation works on detail page
- [ ] Tree refreshes after CRUD operations
- [ ] Permission gating on create/edit/delete buttons
- [ ] `bun run typecheck && bun run lint` passes
