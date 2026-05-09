# Phase 3 — Domain Management

Nested organization tree, domain CRUD, and tree navigation. Uses Axios API functions, TanStack Query hooks, TanStack Form + Zod.

Domains support nesting via `parentId` — a domain can have a parent domain, forming a tree. Domains also have an `isOrganization` flag.

## Step 1: API Functions

**File:** `src/lib/api/domains.ts`

```typescript
import { api, unwrap } from "@/lib/api"
import type { Domain, JSONResponse, PaginatedResponse } from "@/lib/api.types"

export const domainApi = {
  list: (params?: { page?: number; size?: number; sort?: string; filters?: string }) =>
    api.get<JSONResponse<PaginatedResponse<Domain>>>("/v1/domains", { params }).then(unwrap),

  get: (id: number) =>
    api.get<JSONResponse<Domain>>(`/v1/domains/${id}`).then(unwrap),

  create: (body: { name: string; description: string; isOrganization: boolean; parentId?: number | null }) =>
    api.post<JSONResponse<number>>("/v1/domains", body).then(unwrap),

  update: (id: number, body: { name: string; description: string; isOrganization: boolean; parentId?: number | null }) =>
    api.put<JSONResponse<null>>(`/v1/domains/${id}`, body).then(unwrap),

  delete: (id: number) =>
    api.delete<JSONResponse<null>>(`/v1/domains/${id}`).then(unwrap),
}
```

## Step 2: TanStack Query Hooks & Query Keys

**File:** `src/lib/queries/domains.ts`

```typescript
export const domainKeys = {
  all: ["domains"] as const,
  lists: () => [...domainKeys.all, "list"] as const,
  list: (filters: object) => [...domainKeys.lists(), filters] as const,
  details: () => [...domainKeys.all, "detail"] as const,
  detail: (id: number) => [...domainKeys.details(), id] as const,
}

export function useDomains(filters?: object) {
  return useQuery({
    queryKey: domainKeys.list(filters ?? {}),
    queryFn: () => domainApi.list(filters),
  })
}

export function useDomain(id: number) {
  return useQuery({
    queryKey: domainKeys.detail(id),
    queryFn: () => domainApi.get(id),
    enabled: !!id,
  })
}

export function useCreateDomain() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: domainApi.create,
    onSuccess: () => qc.invalidateQueries({ queryKey: domainKeys.all }),
  })
}

export function useUpdateDomain() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, ...body }: { id: number } & Record<string, unknown>) =>
      domainApi.update(id, body),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: domainKeys.detail(id) })
      qc.invalidateQueries({ queryKey: domainKeys.lists() })
    },
  })
}

export function useDeleteDomain() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: domainApi.delete,
    onSuccess: () => qc.invalidateQueries({ queryKey: domainKeys.all }),
  })
}
```

## Step 3: Zod Schemas

**File:** `src/lib/schemas/domain.ts`

```typescript
export const domainSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string().optional().default(""),
  isOrganization: z.boolean().optional().default(false),
  parentId: z.number().nullable().optional(),
})
```

## Step 4: Domain List Page — Tree View

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
1. Fetch all domains via `useDomains()` TanStack Query hook (no pagination, or fetch all pages)
2. Build a tree structure client-side using `parentId` references
3. Render with a recursive tree component

Each tree node shows:
- Domain name
- Badge: "Organization" if `isOrganization === true`
- Actions: View, Edit, Delete — gated by `useCanEdit("domains")`, `useCanDelete("domains")`

"Create Domain" button at top — gated by `useCanCreate("domains")`

**shadcn components needed:** `collapsible` (or custom recursive component)

## Step 5: Domain Tree Component

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

Helper function in `src/lib/build-tree.ts`:
```typescript
interface TreeNode<T> {
  data: T
  children: TreeNode<T>[]
}

function buildTree<T extends { id: number; parentId: number | null }>(items: T[]): TreeNode<T>[]
```

## Step 6: Domain Detail Page

**File:** `src/routes/_authenticated/domains/$domainId.tsx`

Use `useDomain(id)` TanStack Query hook.

Display sections:
- **Info**: name, description, isOrganization flag, parent domain link
- **Children**: list of child domains (links to their detail pages)
- **Breadcrumb**: path from root to current domain for navigation

Edit form using TanStack Form + `domainSchema` — gated by `useCanEdit("domains")`

Delete button with confirmation — gated by `useCanDelete("domains")`

## Step 7: Create Domain Dialog

**File:** `src/components/domains/create-domain-dialog.tsx`

Form using TanStack Form with `zodValidator` and `domainSchema`:
- Name (required)
- Description
- Is Organization (checkbox/toggle)
- Parent Domain (select dropdown — fetch domains via `useDomains()` for `parentId`)

Calls `useCreateDomain()` mutation on submit.

## Step 8: Edit Domain Dialog

**File:** `src/components/domains/edit-domain-dialog.tsx`

Same fields as create, pre-filled with current domain data via `defaultValues`.

## Step 9: Delete Domain Dialog

**File:** `src/components/domains/delete-domain-dialog.tsx`

Confirmation dialog. Warn if domain has children.

## Step 10: Breadcrumb Navigation

**File:** `src/components/domains/domain-breadcrumb.tsx`

Given a domain, walk up the parent chain to build breadcrumb:
```
Root Organization > Engineering > Frontend Team
```

Each segment is a link to the domain detail page.

## Step 11: shadcn Components to Install

```bash
cd web
bunx shadcn add collapsible breadcrumb select switch
```

## Step 12: Files to Create

| File | Action |
|---|---|
| `src/lib/api/domains.ts` | Create: domain Axios API functions |
| `src/lib/queries/domains.ts` | Create: domain TanStack Query hooks + query keys |
| `src/lib/schemas/domain.ts` | Create: domain Zod schema |
| `src/lib/build-tree.ts` | Create: generic tree builder from flat list with parentId |
| `src/routes/_authenticated/domains/index.tsx` | Create: domain tree page |
| `src/routes/_authenticated/domains/$domainId.tsx` | Create: domain detail page |
| `src/components/domains/domain-tree.tsx` | Create: recursive tree component |
| `src/components/domains/domain-tree-node.tsx` | Create: tree node component |
| `src/components/domains/create-domain-dialog.tsx` | Create: create domain dialog (TanStack Form) |
| `src/components/domains/edit-domain-dialog.tsx` | Create: edit domain dialog (TanStack Form) |
| `src/components/domains/delete-domain-dialog.tsx` | Create: delete confirmation |
| `src/components/domains/domain-breadcrumb.tsx` | Create: breadcrumb navigation |

## Completion Criteria

- [ ] Domain tree renders with expand/collapse
- [ ] Nested domains display correctly
- [ ] Create domain dialog works with TanStack Form + Zod validation
- [ ] Edit domain works
- [ ] Delete domain with confirmation works
- [ ] Breadcrumb navigation works on detail page
- [ ] Tree refreshes after CRUD operations (TanStack Query cache invalidation)
- [ ] Permission gating on create/edit/delete buttons
- [ ] `bun run typecheck && bun run lint` passes
