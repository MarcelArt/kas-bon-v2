import { Outlet, createFileRoute, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/_auth")({
  beforeLoad: async () => {
    if (typeof window === "undefined") return

    const { user, accessToken } = await import("@/lib/stores/auth-store").then(
      (m) => {
        const s = m.useAuthStore.getState()
        return { user: s.user, accessToken: s.accessToken }
      },
    )
    if (user && accessToken) {
      throw redirect({ to: "/dashboard" })
    }
  },
  component: AuthLayout,
})

function AuthLayout() {
  return (
    <div className="flex min-h-svh items-center justify-center">
      <Outlet />
    </div>
  )
}
