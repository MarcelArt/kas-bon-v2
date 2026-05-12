import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_authenticated/no-access")({
  component: NoAccessPage,
})

function NoAccessPage() {
  return (
    <div className="flex min-h-svh items-center justify-center">
      <div className="text-center">
        <h1 className="text-2xl font-semibold">No Access</h1>
        <p className="text-muted-foreground mt-2">
          You don&apos;t have access to any organization. Contact your administrator.
        </p>
      </div>
    </div>
  )
}
