export interface User {
  id: number
  username: string
  email: string
  createdAt: string
  updatedAt: string
  deletedAt: string | null
}

export interface UserInput {
  username: string
  email: string
  password: string
}

export interface LoginInput {
  username: string
  password: string
  isRemember: boolean
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
  user: User
}

export interface Role {
  id: number
  name: string
  description: string
  domainId: number
  domain: Domain | null
  createdAt: string
  updatedAt: string
}

export interface RoleInput {
  name: string
  description: string
  domainId: number
}

export interface Permission {
  id: number
  name: string
  description: string
  appId: number
  app: App | null
  createdAt: string
  updatedAt: string
}

export interface PermissionInput {
  name: string
  description: string
  appId: number
}

export interface Domain {
  id: number
  name: string
  description: string
  isOrganization: boolean
  parentId: number | null
  parent: Domain | null
  createdAt: string
  updatedAt: string
}

export interface DomainInput {
  name: string
  description: string
  isOrganization: boolean
  parentId: number | null
}

export interface App {
  id: number
  name: string
  description: string
  createdAt: string
  updatedAt: string
}

export interface AppInput {
  name: string
  description: string
}

export interface JSONResponse<T> {
  items: T | null
  isSuccess: boolean
  message: string
}

export interface PaginatedResponse<T> {
  items: Array<T>
  page: number
  size: number
  total: number
  total_pages: number
}

export interface AccessControlEval {
  sub: string
  app: string
  dom: string
  obj: string
  act: string
}

export interface TokenEndpointRequest {
  appId: number
  domainId: number
  permission: string
}
