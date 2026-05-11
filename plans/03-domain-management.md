# Phase 3 — Domain Management

Nested organization tree, domain CRUD, and tree navigation. Uses TanStack Start server functions, TanStack Query hooks, TanStack Form + Zod.

Domains support nesting via `parentId` — a domain can have a parent domain, forming a tree. Domains also have an `isOrganization` flag.

## Step 1: Domain Server Functions

**File:** `src/lib/server/domains.ts`

Server functions for all domain endpoints (see `.opencode/API.md` — Domains section):

```typescript
import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { Domain, PaginatedResponse } from "@/lib/api.types"

const authContextSchema = z.object({
  accessToken: z.string(),
  domainId: z.number().optional(),
  appId: z.number().optional(),
})

export const getDomainsFn = createServerFn()
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
    return serverApi<PaginatedResponse<Domain>>("GET", "/domains", { accessToken, domainId, appId }, undefined, params)
  })

export const getDomainFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<Domain>("GET", `/domains/${id}`, { accessToken, domainId, appId })
  })

export const createDomainFn = createServerFn()
  .validator(
    authContextSchema.extend({
      name: z.string(),
      description: z.string().optional().default(""),
      isOrganization: z.boolean().optional().default(false),
      parentId: z.number().nullable().optional(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, ...body } = data
    return serverApi<number>("POST", "/domains", { accessToken, domainId, appId }, body)
  })

export const updateDomainFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
      name: z.string(),
      description: z.string().optional().default(""),
      isOrganization: z.boolean().optional().default(false),
      parentId: z.number().nullable().optional(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id, ...body } = data
    return serverApi<null>("PUT", `/domains/${id}`, { accessToken, domainId, appId }, body)
  })

export const deleteDomainFn = createServerFn()
  .validator(
    authContextSchema.extend({
      id: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    const { accessToken, domainId, appId, id } = data
    return serverApi<null>("DELETE", `/domains/${id}`, { accessToken, domainId, appId })
  })
```

## Step 2: TanStack Query Hooks & Query Keys

**File:** `src/lib/queries/domains.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { getDomainsFn, getDomainFn, createDomainFn, updateDomainFn, deleteDomainFn } from "@/lib/server/domains"
import { useAuthContext, withAuthRetry } from "./auth-context"

export const domainKeys = {
  all: ["domains"] as const,
  lists: () => [...domainKeys.all, "list"] as const,
  list: (filters: object) => [...domainKeys.lists(), filters] as const,
  details: () => [...domainKeys.all, "detail"] as const,
  detail: (id: number) => [...domainKeys.details(), id] as const,
}

export function useDomains(filters?: { page?: number; size?: number; sort?: string; filters?: string }) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: domainKeys.list(filters ?? {}),
    queryFn: withAuthRetry(() => getDomainsFn({ data: { ...auth, ...filters } })),
    enabled: !!auth.accessToken,
  })
}

export function useDomain(id: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: domainKeys.detail(id),
    queryFn: withAuthRetry(() => getDomainFn({ data: { ...auth, id } })),
    enabled: !!auth.accessToken && !!id,
  })
}

export function useCreateDomain() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (body: { name: string; description?: string; isOrganization?: boolean; parentId?: number | null }) =>
      createDomainFn({ data: { ...auth, ...body } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: domainKeys.all }),
  })
}

export function useUpdateDomain() {
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
      isOrganization?: boolean
      parentId?: number | null
    }) => updateDomainFn({ data: { ...auth, id, ...body } }),
    onSuccess: (_, { id }) => {
      qc.invalidateQueries({ queryKey: domainKeys.detail(id) })
      qc.invalidateQueries({ queryKey: domainKeys.lists() })
    },
  })
}

export function useDeleteDomain() {
  const qc = useQueryClient()
  const auth = useAuthContext()
  return useMutation({
    mutationFn: (id: number) => deleteDomainFn({ data: { ...auth, id } }),
    onSuccess: () => qc.invalidateQueries({ queryKey: domainKeys.all }),
  })
}
```

## Step 3: Zod Schemas

**File:** `src/lib/schemas/domain.ts`

```typescript
import { z } from "zod"

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

function buildTree<T extends { ID: number; parentId: number | null }>(items: T[]): TreeNode<T>[]
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
| `src/lib/server/domains.ts` | Create: domain server functions |
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
- [ ] Create domain dialog works with TanStack Form + Zod → calls `createDomainFn`
- [ ] Edit domain works via `updateDomainFn`
- [ ] Delete domain with confirmation works via `deleteDomainFn`
- [ ] Breadcrumb navigation works on detail page
- [ ] Tree refreshes after CRUD operations (TanStack Query cache invalidation)
- [ ] Permission gating on create/edit/delete buttons
- [ ] `bun run typecheck && bun run lint` passes
