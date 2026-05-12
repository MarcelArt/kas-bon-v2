import { createFileRoute, redirect } from "@tanstack/react-router"
import { AppShell } from "@/components/layout/app-shell"

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async () => {
    if (typeof window === "undefined") return

    const { useAuthStore } = await import("@/lib/stores/auth-store")
    const { user, accessToken, domainId } = useAuthStore.getState()

    if (!user || !accessToken) {
      throw redirect({ to: "/login" })
    }
    if (!domainId) {
      throw redirect({ to: "/select-organization" })
    }
  },
  component: AuthenticatedLayout,
})

function AuthenticatedLayout() {
  return (
    <AppShell />
  )
}
