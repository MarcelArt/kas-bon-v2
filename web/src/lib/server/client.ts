import type { JSONResponse } from "@/lib/api.types"
import { getApiBaseUrl } from "@/lib/env"

interface RequestOptions {
  accessToken?: string | null
  domainId?: number | null
  appId?: number | null
  customHeaders?: Record<string, string>
}

export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(`API_ERROR:${status}:${message}`)
    this.status = status
    this.name = "ApiError"
  }
}

export async function serverApi<T>(
  method: string,
  path: string,
  options?: RequestOptions,
  body?: unknown,
  params?: Record<string, string | number | undefined>,
): Promise<T> {
  const url = new URL(`/api/v1${path}`, getApiBaseUrl())
  if (params) {
    Object.entries(params).forEach(([k, v]) => {
      if (v != null) url.searchParams.set(k, String(v))
    })
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.customHeaders ?? {}),
  }
  if (options?.accessToken) headers["Authorization"] = `Bearer ${options.accessToken}`
  if (options?.domainId != null) headers["X-Domain-Id"] = String(options.domainId)
  if (options?.appId != null) headers["X-App-Id"] = String(options.appId)

  const res = await fetch(url.toString(), {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  })

  const json: JSONResponse<T> = await res.json()

  if (res.status === 401) {
    throw new ApiError("AUTH_EXPIRED:" + (json.message || "Token expired"), 401)
  }

  if (!res.ok) {
    throw new ApiError(json.message || res.statusText, res.status)
  }

  if (json.isSuccess === false) {
    throw new ApiError(json.message || "Request failed", res.status)
  }

  if ("isSuccess" in json) {
    return json.items as T
  }

  return json
}
