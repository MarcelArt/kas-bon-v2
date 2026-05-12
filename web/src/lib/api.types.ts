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
}

export interface User {
  ID: number
  username: string
  email: string
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
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

export interface Domain {
  ID: number
  name: string
  description: string
  isOrganization: boolean
  parentId: number | null
  parent: Domain | null
  CreatedAt: string
  UpdatedAt: string
}

export interface App {
  ID: number
  name: string
  description: string
  CreatedAt: string
  UpdatedAt: string
}

export interface Role {
  ID: number
  name: string
  description: string
  domainId: number
  domain: Domain | null
  CreatedAt: string
  UpdatedAt: string
}

export interface Permission {
  ID: number
  name: string
  description: string
  appId: number
  app: App | null
  CreatedAt: string
  UpdatedAt: string
}
