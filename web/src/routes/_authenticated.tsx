import { Outlet, createFileRoute } from "@tanstack/react-router"
import { AppShell } from "@/components/layout/app-shell"
import { AuthProvider } from "@/lib/auth"
import { checkAuth } from "@/lib/auth-guard"

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async () => {
    const { user, permissions } = await checkAuth()
    return { user, permissions }
  },
  component: AuthenticatedLayout,
})

function AuthenticatedLayout() {
  const { user, permissions } = Route.useRouteContext()

  return (
    <AuthProvider initialUser={user} initialPermissions={permissions}>
      <AppShell>
        <Outlet />
      </AppShell>
    </AuthProvider>
  )
}
