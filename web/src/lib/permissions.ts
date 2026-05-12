export function parsePermissionTuples(tuples: Array<Array<string>>): Set<string> {
  return new Set(tuples.map((t) => `${t[3]}#${t[4]}`))
}

export function isSuperUser(permissions: Set<string>): boolean {
  return permissions.has("all#fullAccess")
}

export function checkPermission(
  permissions: Set<string>,
  resource: string,
  action: string,
): boolean {
  if (isSuperUser(permissions)) return true
  return permissions.has(`${resource}#${action}`)
}

export const RESOURCES = {
  USERS: "users",
  DOMAINS: "domains",
  APPS: "apps",
  ROLES: "roles",
  PERMISSIONS: "permissions",
  ALL: "all",
} as const

export const ACTIONS = {
  READ: "read",
  CREATE: "create",
  UPDATE: "update",
  DELETE: "delete",
  FULL_ACCESS: "fullAccess",
} as const
