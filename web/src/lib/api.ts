import type { JSONResponse } from "@/lib/api.types"

const API_BASE_URL = "/api/v1"

function getAccessToken(): string | null {
  const match = document.cookie.match(/(?:^|;\s*)kb_access_token=([^;]*)/)
  return match ? decodeURIComponent(match[1]) : null
}

class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.name = "ApiError"
    this.status = status
  }
}

async function request<T>(
  method: string,
  path: string,
  options?: {
    body?: unknown
    headers?: Record<string, string>
  },
): Promise<T> {
  const token = getAccessToken()
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...options?.headers,
  }

  if (token) {
    headers["Authorization"] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers,
    body: options?.body ? JSON.stringify(options.body) : undefined,
  })

  if (response.status === 401) {
    window.location.href = "/login"
    throw new ApiError("Unauthorized", 401)
  }

  const data: JSONResponse<T> = await response.json()

  if (!data.isSuccess) {
    throw new ApiError(data.message || "Request failed", response.status)
  }

  return data.items as T
}

export const api = {
  get: <T>(path: string, headers?: Record<string, string>) =>
    request<T>("GET", path, { headers }),

  post: <T>(path: string, body?: unknown, headers?: Record<string, string>) =>
    request<T>("POST", path, { body, headers }),

  put: <T>(path: string, body?: unknown, headers?: Record<string, string>) =>
    request<T>("PUT", path, { body, headers }),

  patch: <T>(
    path: string,
    body?: unknown,
    headers?: Record<string, string>,
  ) => request<T>("PATCH", path, { body, headers }),

  delete: <T>(path: string, headers?: Record<string, string>) =>
    request<T>("DELETE", path, { headers }),
}

export { ApiError }
