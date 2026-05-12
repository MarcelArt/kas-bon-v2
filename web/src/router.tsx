import { createRouter as createTanStackRouter } from "@tanstack/react-router"
import { QueryClient, queryOptions } from "@tanstack/react-query"
import { routeTree } from "./routeTree.gen"

export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 1000 * 60 * 5,
      },
    },
  })
}

let queryClient: QueryClient | undefined

export function getQueryClient() {
  if (!queryClient) {
    queryClient = createQueryClient()
  }
  return queryClient
}

export function getRouter() {
  const qc = getQueryClient()

  const router = createTanStackRouter({
    routeTree,
    context: { queryClient: qc },
    scrollRestoration: true,
    defaultPreload: "intent",
    defaultPreloadStaleTime: 0,
  })

  return router
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}

export { queryOptions }
