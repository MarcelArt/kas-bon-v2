import { useAuthStore } from "@/lib/stores/auth-store"
import { refreshTokenFn } from "@/lib/server/auth"

export function getAuthContext() {
  const { accessToken, domainId, appId } = useAuthStore.getState()
  return {
    accessToken: accessToken ?? undefined,
    domainId: domainId ?? undefined,
    appId: appId ?? undefined,
  }
}

export function useAuthContext() {
  const accessToken = useAuthStore((s) => s.accessToken)
  const domainId = useAuthStore((s) => s.domainId)
  const appId = useAuthStore((s) => s.appId)
  return {
    accessToken: accessToken ?? undefined,
    domainId: domainId ?? undefined,
    appId: appId ?? undefined,
  }
}

export function withAuthRetry<T>(queryFn: () => Promise<T>): () => Promise<T> {
  return async () => {
    try {
      return await queryFn()
    } catch (error) {
      if (error instanceof Error && error.message.includes("AUTH_EXPIRED:")) {
        const { refreshToken, setTokens, logout } = useAuthStore.getState()
        if (!refreshToken) {
          logout()
          window.location.href = "/login"
          throw error
        }
        try {
          const result = await refreshTokenFn({ data: { refreshToken } })
          setTokens(result.accessToken, result.refreshToken)
          return await queryFn()
        } catch {
          logout()
          window.location.href = "/login"
          throw error
        }
      }
      throw error
    }
  }
}
