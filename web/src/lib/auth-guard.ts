import { redirect } from "@tanstack/react-router"
import type { User } from "@/lib/api.types"
import { getCurrentUserFn, getUserPermissionsFn } from "@/lib/auth.fns"

interface AuthGuardResult {
  user: User
  permissions: Array<Array<string>>
}

export async function checkAuth(): Promise<AuthGuardResult> {
  const user = await getCurrentUserFn()

  if (!user) {
    throw redirect({ to: "/login" })
  }

  const permissions = await getUserPermissionsFn({
    data: { userId: user.id, appId: 1, domainId: 1 },
  })

  return { user, permissions }
}
