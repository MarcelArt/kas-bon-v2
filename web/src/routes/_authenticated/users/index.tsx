import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_authenticated/users/")({
  component: UsersPage,
})

function UsersPage() {
  return (
    <div>
      <h1 className="text-lg font-medium">Users</h1>
    </div>
  )
}
