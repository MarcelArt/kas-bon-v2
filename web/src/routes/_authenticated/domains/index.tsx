import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_authenticated/domains/")({
  component: DomainsPage,
})

function DomainsPage() {
  return (
    <div>
      <h1 className="text-lg font-medium">Domains</h1>
    </div>
  )
}
