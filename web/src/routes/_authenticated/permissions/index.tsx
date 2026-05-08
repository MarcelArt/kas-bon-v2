import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_authenticated/permissions/")({
  component: PermissionsPage,
})

function PermissionsPage() {
  return (
    <div>
      <h1 className="text-lg font-medium">Permissions</h1>
    </div>
  )
}
