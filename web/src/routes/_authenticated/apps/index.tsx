import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_authenticated/apps/")({
  component: AppsPage,
})

function AppsPage() {
  return (
    <div>
      <h1 className="text-lg font-medium">Apps</h1>
    </div>
  )
}
