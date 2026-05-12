import {
  AppWindow,
  Folders,
  House,
  Key,
  Shield,
  SignOut,
  Users,
} from "@phosphor-icons/react"
import { Outlet, useNavigate } from "@tanstack/react-router"
import { useEffect } from "react"
import { Separator } from "@/components/ui/separator"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
} from "@/components/ui/sidebar"
import { usePermission } from "@/hooks/use-permission"
import { useAuthStore } from "@/lib/stores/auth-store"
import { getPermissionsFn } from "@/lib/server/auth"

const BOOTSTRAP_APP_ID = 1

const NAV_ITEMS = [
  { to: "/dashboard", icon: House, label: "Dashboard", resource: null, action: null },
  { to: "/users", icon: Users, label: "Users", resource: "users", action: "read" },
  { to: "/domains", icon: Folders, label: "Domains", resource: "domains", action: "read" },
  { to: "/apps", icon: AppWindow, label: "Apps", resource: "apps", action: "read" },
  { to: "/roles", icon: Shield, label: "Roles", resource: "roles", action: "read" },
  { to: "/permissions", icon: Key, label: "Permissions", resource: "permissions", action: "read" },
] as const

function NavItem({
  to,
  icon: Icon,
  label,
  resource,
  action,
}: {
  to: string
  icon: React.ComponentType<{ className?: string }>
  label: string
  resource: string | null
  action: string | null
}) {
  const hasAccess = resource && action ? usePermission(resource, action) : true
  if (!hasAccess) return null

  return (
    <SidebarMenuItem>
      <SidebarMenuButton asChild>
        <a href={to}>
          <Icon className="size-4" />
          <span>{label}</span>
        </a>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )
}

function useEnsurePermissions() {
  const user = useAuthStore((s) => s.user)
  const accessToken = useAuthStore((s) => s.accessToken)
  const domainId = useAuthStore((s) => s.domainId)
  const appId = useAuthStore((s) => s.appId)
  const permissions = useAuthStore((s) => s.permissions)
  const setPermissions = useAuthStore((s) => s.setPermissions)
  const setApp = useAuthStore((s) => s.setApp)

  useEffect(() => {
    if (!user || !accessToken || !domainId || permissions.size > 0) return

    const effectiveAppId = appId ?? BOOTSTRAP_APP_ID

    getPermissionsFn({
      data: { accessToken, userId: user.ID, domainId, appId: effectiveAppId },
    })
      .then((tuples) => {
        setPermissions(tuples)
        if (!appId) setApp(effectiveAppId)
      })
      .catch((err) => {
        console.error("Failed to fetch permissions:", err)
      })
  }, [user, accessToken, domainId, appId, permissions.size, setPermissions, setApp])
}

export function AppShell() {
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const navigate = useNavigate()

  useEnsurePermissions()

  function handleLogout() {
    logout()
    navigate({ to: "/login" })
  }

  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarHeader className="p-4">
          <h2 className="text-lg font-semibold">KAS Bon</h2>
        </SidebarHeader>
        <Separator />
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {NAV_ITEMS.map((item) => (
                  <NavItem key={item.to} {...item} />
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <Separator />
          <div className="flex items-center justify-between p-2">
            <div className="flex flex-col text-sm">
              <span className="font-medium">{user?.username}</span>
              <span className="text-muted-foreground text-xs">{user?.email}</span>
            </div>
            <button
              onClick={handleLogout}
              className="text-muted-foreground hover:text-foreground"
            >
              <SignOut className="size-4" />
            </button>
          </div>
        </SidebarFooter>
        <SidebarRail />
      </Sidebar>
      <SidebarInset>
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  )
}
