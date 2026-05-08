import {
  createContext,
  useCallback,
  useContext,
  useMemo,
} from "react"
import type { User } from "@/lib/api.types"
import { useLogoutMutation, useUserPermissions } from "@/lib/auth.query"

interface AuthContextValue {
  user: User | null
  permissions: Set<string>
  isSuperUser: boolean
  hasPermission: (resource: string, action: string) => boolean
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextValue | null>(null)

interface AuthProviderProps {
  children: React.ReactNode
  initialUser: User | null
  initialPermissions: Array<Array<string>>
}

function parsePermissions(tuples: Array<Array<string>>): Set<string> {
  const permissions = new Set<string>()
  for (const tuple of tuples) {
    if (tuple.length >= 5) {
      const resource = tuple[3]
      const action = tuple[4]
      permissions.add(`${resource}#${action}`)
    }
  }
  return permissions
}

export function AuthProvider({
  children,
  initialUser,
  initialPermissions,
}: AuthProviderProps) {
  const logoutMutation = useLogoutMutation()

  const permissionsData = useUserPermissions(
    initialUser?.id,
  )

  const user = initialUser

  const permissions = useMemo(() => {
    const raw = permissionsData.data ?? initialPermissions
    return parsePermissions(raw)
  }, [permissionsData.data, initialPermissions])

  const isSuperUser = permissions.has("all#fullAccess")

  const hasPermission = useCallback(
    (resource: string, action: string) => {
      if (isSuperUser) return true
      return permissions.has(`${resource}#${action}`)
    },
    [isSuperUser, permissions],
  )

  const logout = useCallback(() => {
    logoutMutation.mutate()
  }, [logoutMutation])

  const isLoading = logoutMutation.isPending

  return (
    <AuthContext.Provider
      value={{ user, permissions, isSuperUser, hasPermission, logout, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
