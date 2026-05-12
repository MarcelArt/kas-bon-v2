import { useAuthStore } from "@/lib/stores/auth-store"

export function usePermission(resource: string, action: string): boolean {
  const permissions = useAuthStore((s) => s.permissions)
  if (permissions.has("all#fullAccess")) return true
  return permissions.has(`${resource}#${action}`)
}

export function useCanCreate(resource: string): boolean {
  return usePermission(resource, "create")
}

export function useCanEdit(resource: string): boolean {
  return usePermission(resource, "update")
}

export function useCanDelete(resource: string): boolean {
  return usePermission(resource, "delete")
}

export function useIsSuperUser(): boolean {
  const permissions = useAuthStore((s) => s.permissions)
  return permissions.has("all#fullAccess")
}
