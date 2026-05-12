import { create } from "zustand"
import { persist } from "zustand/middleware"
import type { Domain, User } from "@/lib/api.types"
import { parsePermissionTuples } from "@/lib/permissions"

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  domainId: number | null
  appId: number | null
  organizations: Array<Domain>
  permissions: Set<string>

  setUser: (user: User | null) => void
  setTokens: (access: string, refresh: string) => void
  setDomain: (domainId: number) => void
  setApp: (appId: number) => void
  setOrganizations: (orgs: Array<Domain>) => void
  setPermissions: (tuples: Array<Array<string>>) => void
  hasPermission: (resource: string, action: string) => boolean
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      domainId: null,
      appId: null,
      organizations: [],
      permissions: new Set(),

      setUser: (user) => set({ user }),
      setTokens: (accessToken, refreshToken) => set({ accessToken, refreshToken }),
      setDomain: (domainId) => set({ domainId }),
      setApp: (appId) => set({ appId }),
      setOrganizations: (organizations) => set({ organizations }),
      setPermissions: (tuples) => {
        const permissions = parsePermissionTuples(tuples)
        set({ permissions })
      },
      hasPermission: (resource, action) => {
        const { permissions } = get()
        if (permissions.has("all#fullAccess")) return true
        return permissions.has(`${resource}#${action}`)
      },
      logout: () =>
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          domainId: null,
          appId: null,
          organizations: [],
          permissions: new Set(),
        }),
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        domainId: state.domainId,
        appId: state.appId,
      }),
    },
  ),
)
