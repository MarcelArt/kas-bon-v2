import { createServerFn } from "@tanstack/react-start"
import { deleteCookie, getCookie, setCookie } from "@tanstack/react-start/server"
import type { JSONResponse, LoginInput, LoginResponse, User, UserInput } from "@/lib/api.types"

const API_BASE_URL = process.env.API_BASE_URL || "http://localhost:8080/api/v1"

const COOKIE_OPTIONS = {
  path: "/",
  sameSite: "lax" as const,
  httpOnly: false,
}

const REFRESH_COOKIE_OPTIONS = {
  path: "/",
  sameSite: "lax" as const,
  httpOnly: true,
}

export const loginFn = createServerFn({ method: "POST" })
  .inputValidator((data: LoginInput) => data)
  .handler(async ({ data }) => {
    const response = await fetch(`${API_BASE_URL}/users/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })

    const json: JSONResponse<LoginResponse> = await response.json()

    if (!json.isSuccess || !json.items) {
      throw new Error(json.message || "Login failed")
    }

    const { accessToken, refreshToken, user } = json.items

    const maxAge = data.isRemember ? 30 * 24 * 60 * 60 : 24 * 60 * 60

    setCookie("kb_access_token", accessToken, {
      ...COOKIE_OPTIONS,
      maxAge,
    })
    setCookie("kb_refresh_token", refreshToken, {
      ...REFRESH_COOKIE_OPTIONS,
      maxAge: data.isRemember ? 30 * 24 * 60 * 60 : 7 * 24 * 60 * 60,
    })
    setCookie("kb_user", JSON.stringify(user), {
      ...COOKIE_OPTIONS,
      maxAge,
    })

    return user
  })

export const registerFn = createServerFn({ method: "POST" })
  .inputValidator((data: UserInput) => data)
  .handler(async ({ data }) => {
    const response = await fetch(`${API_BASE_URL}/users`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })

    const json: JSONResponse<number> = await response.json()

    if (!json.isSuccess) {
      throw new Error(json.message || "Registration failed")
    }

    return json.items
  })

export const refreshFn = createServerFn({ method: "POST" }).handler(
  async () => {
    const refreshToken = getCookie("kb_refresh_token")
    if (!refreshToken) {
      throw new Error("No refresh token")
    }

    const response = await fetch(`${API_BASE_URL}/users/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Refresh-Token": refreshToken,
      },
    })

    const json: JSONResponse<LoginResponse> = await response.json()

    if (!json.isSuccess || !json.items) {
      deleteCookie("kb_access_token", { path: "/" })
      deleteCookie("kb_refresh_token", { path: "/" })
      deleteCookie("kb_user", { path: "/" })
      throw new Error(json.message || "Token refresh failed")
    }

    const { accessToken, refreshToken: newRefreshToken, user } = json.items

    setCookie("kb_access_token", accessToken, { ...COOKIE_OPTIONS, maxAge: 24 * 60 * 60 })
    setCookie("kb_refresh_token", newRefreshToken, {
      ...REFRESH_COOKIE_OPTIONS,
      maxAge: 7 * 24 * 60 * 60,
    })
    setCookie("kb_user", JSON.stringify(user), {
      ...COOKIE_OPTIONS,
      maxAge: 24 * 60 * 60,
    })

    return user
  },
)

export const logoutFn = createServerFn({ method: "POST" }).handler(() => {
  deleteCookie("kb_access_token", { path: "/" })
  deleteCookie("kb_refresh_token", { path: "/" })
  deleteCookie("kb_user", { path: "/" })
  return true
})

export const getCurrentUserFn = createServerFn({ method: "GET" }).handler(
  () => {
    const userCookie = getCookie("kb_user")
    const accessToken = getCookie("kb_access_token")

    if (!userCookie || !accessToken) {
      return null
    }

    try {
      return JSON.parse(userCookie) as User
    } catch {
      return null
    }
  },
)

export const getUserPermissionsFn = createServerFn({ method: "POST" })
  .inputValidator(
    (data: { userId: number; appId: number; domainId: number }) => data,
  )
  .handler(async ({ data }) => {
    const accessToken = getCookie("kb_access_token")
    if (!accessToken) {
      return []
    }

    const response = await fetch(
      `${API_BASE_URL}/users/${data.userId}/permissions`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
          "X-App-Id": String(data.appId),
          "X-Domain-Id": String(data.domainId),
        },
      },
    )

    const json: JSONResponse<Array<Array<string>>> = await response.json()

    if (!json.isSuccess || !json.items) {
      return []
    }

    return json.items
  })
