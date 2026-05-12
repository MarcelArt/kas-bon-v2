import { useMutation, useQuery } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { useAuthContext, withAuthRetry } from "./auth-context"
import type { LoginResponse } from "@/lib/api.types"
import { useAuthStore } from "@/lib/stores/auth-store"
import { getAppsFn, getOrganizationsFn, getPermissionsFn, loginFn } from "@/lib/server/auth"

const BOOTSTRAP_APP_ID = 1

export const authKeys = {
  user: ["auth", "user"] as const,
  organizations: (userId: number) => ["auth", "organizations", userId] as const,
  permissions: (userId: number) => ["auth", "permissions", userId] as const,
}

export function useLogin() {
  const { setUser, setTokens, setOrganizations, setDomain, setApp, setPermissions } = useAuthStore()
  const navigate = useNavigate()

  return useMutation<LoginResponse, Error, { username: string; password: string; isRemember?: boolean }>({
    mutationFn: (data) =>
      loginFn({ data }),
    onSuccess: async (loginData) => {
      setUser(loginData.user)
      setTokens(loginData.accessToken, loginData.refreshToken)

      const orgs = await getOrganizationsFn({
        data: { accessToken: loginData.accessToken, userId: loginData.user.ID },
      })
      setOrganizations(orgs)

      if (orgs.length === 0) {
        navigate({ to: "/no-access" })
      } else if (orgs.length === 1) {
        const domainId = orgs[0].ID
        setDomain(domainId)

        const apps = await getAppsFn({
          data: { accessToken: loginData.accessToken, domainId, appId: BOOTSTRAP_APP_ID },
        })
        if (apps.length > 0) {
          setApp(apps[0].ID)
        }

        const tuples = await getPermissionsFn({
          data: { accessToken: loginData.accessToken, userId: loginData.user.ID, domainId, appId: apps[0]?.ID ?? BOOTSTRAP_APP_ID },
        })
        setPermissions(tuples)
        navigate({ to: "/dashboard" })
      } else {
        navigate({ to: "/select-organization" })
      }
    },
  })
}

export function useUserOrganizations(userId: number) {
  const auth = useAuthContext()
  return useQuery({
    queryKey: authKeys.organizations(userId),
    queryFn: withAuthRetry(() =>
      getOrganizationsFn({ data: { accessToken: auth.accessToken!, userId } }),
    ),
    enabled: !!auth.accessToken && !!userId,
  })
}

export function useUserPermissions(userId: number) {
  const auth = useAuthContext()
  const { setPermissions } = useAuthStore()

  return useQuery({
    queryKey: authKeys.permissions(userId),
    queryFn: withAuthRetry(async () => {
      const tuples = await getPermissionsFn({
        data: { accessToken: auth.accessToken!, userId, domainId: auth.domainId, appId: auth.appId },
      })
      setPermissions(tuples)
      return tuples
    }),
    enabled: !!auth.accessToken && !!userId,
  })
}
