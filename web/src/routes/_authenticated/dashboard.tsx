import { createFileRoute } from "@tanstack/react-router"
import { useAuth } from "@/lib/auth"

export const Route = createFileRoute("/_authenticated/dashboard")({
  component: DashboardPage,
})

function DashboardPage() {
  const { user } = useAuth()

  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-lg font-medium">Dashboard</h1>
      <p className="text-sm text-muted-foreground">
        Welcome back, {user?.username ?? "User"}.
      </p>
    </div>
  )
}
