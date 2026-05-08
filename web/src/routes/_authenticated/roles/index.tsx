import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_authenticated/roles/")({
  component: RolesPage,
})

function RolesPage() {
  return (
    <div>
      <h1 className="text-lg font-medium">Roles</h1>
    </div>
  )
}
