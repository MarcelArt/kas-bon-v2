import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { getAppsFn, getPermissionsFn } from "@/lib/server/auth"
import { useAuthStore } from "@/lib/stores/auth-store"

const BOOTSTRAP_APP_ID = 1

export const Route = createFileRoute("/_authenticated/select-organization")({
  component: SelectOrganizationPage,
})

function SelectOrganizationPage() {
  const organizations = useAuthStore((s) => s.organizations)
  const accessToken = useAuthStore((s) => s.accessToken)
  const user = useAuthStore((s) => s.user)
  const setDomain = useAuthStore((s) => s.setDomain)
  const setApp = useAuthStore((s) => s.setApp)
  const setPermissions = useAuthStore((s) => s.setPermissions)
  const navigate = useNavigate()

  async function handleSelect(domainId: number) {
    setDomain(domainId)

    if (accessToken && user) {
      try {
        const apps = await getAppsFn({
          data: { accessToken, domainId, appId: BOOTSTRAP_APP_ID },
        })
        const appId = apps[0]?.ID ?? BOOTSTRAP_APP_ID
        setApp(appId)

        const tuples = await getPermissionsFn({
          data: { accessToken, userId: user.ID, domainId, appId },
        })
        setPermissions(tuples)
      } catch {
        toast.error("Failed to load permissions")
      }
    }

    navigate({ to: "/dashboard" })
  }

  return (
    <div className="flex min-h-svh items-center justify-center">
      <div className="w-full max-w-md space-y-4 p-6">
        <h1 className="text-2xl font-semibold text-center">Select Organization</h1>
        <p className="text-muted-foreground text-center text-sm">
          You belong to multiple organizations. Choose one to continue.
        </p>
        <div className="flex flex-col gap-3">
          {organizations.map((org) => (
            <Card key={org.ID} className="cursor-pointer transition-colors hover:bg-accent">
              <CardHeader className="pb-2">
                <CardTitle className="text-base">{org.name}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground text-sm">{org.description}</p>
                <Button
                  className="mt-3"
                  size="sm"
                  onClick={() => handleSelect(org.ID)}
                >
                  Select
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}
