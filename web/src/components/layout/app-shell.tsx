import { Link, useRouterState } from "@tanstack/react-router"
import {
  AppWindow,
  Folders,
  House,
  Key,
  Shield,
  Users,
} from "@phosphor-icons/react"
import type { ReactNode } from "react"
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
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { Separator } from "@/components/ui/separator"
import { TooltipProvider } from "@/components/ui/tooltip"
import { useAuth } from "@/lib/auth"

interface AppShellProps {
  children: ReactNode
}

const navItems = [
  { title: "Dashboard", href: "/dashboard", icon: House, permission: null },
  { title: "Users", href: "/users", icon: Users, permission: "users#read" },
  { title: "Domains", href: "/domains", icon: Folders, permission: "domains#read" },
  { title: "Apps", href: "/apps", icon: AppWindow, permission: "apps#read" },
  { title: "Roles", href: "/roles", icon: Shield, permission: "roles#read" },
  { title: "Permissions", href: "/permissions", icon: Key, permission: "permissions#read" },
]

function AppSidebar() {
  const { hasPermission, user, logout } = useAuth()
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to="/dashboard">
                <div className="flex aspect-square size-8 items-center justify-center rounded-none bg-primary text-primary-foreground">
                  <span className="text-xs font-bold">KB</span>
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-medium">Kas Bon</span>
                  <span className="text-[10px] text-muted-foreground">Management</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => {
                if (item.permission && !hasPermission(item.permission.split("#")[0], item.permission.split("#")[1])) {
                  return null
                }
                return (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      asChild
                      isActive={currentPath.startsWith(item.href)}
                      tooltip={item.title}
                    >
                      <Link to={item.href}>
                        <item.icon />
                        <span>{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                )
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg">
              <div className="flex aspect-square size-8 items-center justify-center rounded-none bg-muted text-muted-foreground">
                <span className="text-xs font-medium">
                  {user ? user.username.slice(0, 2).toUpperCase() : "U"}
                </span>
              </div>
              <div className="flex flex-col gap-0.5 leading-none">
                <span className="text-xs font-medium">{user?.username ?? "User"}</span>
                <span className="text-[10px] text-muted-foreground">{user?.email ?? ""}</span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton onClick={logout} tooltip="Logout">
              <span>Logout</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}

export function AppShell({ children }: AppShellProps) {
  return (
    <SidebarProvider>
      <TooltipProvider>
        <AppSidebar />
        <SidebarInset>
          <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
          </header>
          <main className="flex-1 p-4">{children}</main>
        </SidebarInset>
      </TooltipProvider>
    </SidebarProvider>
  )
}
