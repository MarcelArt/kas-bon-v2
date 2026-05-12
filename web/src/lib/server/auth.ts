import { createServerFn } from "@tanstack/react-start"
import { z } from "zod"
import { serverApi } from "./client"
import type { App, Domain, LoginResponse } from "@/lib/api.types"

export const loginFn = createServerFn({ method: "POST" })
  .inputValidator(
    z.object({
      username: z.string(),
      password: z.string(),
      isRemember: z.boolean().optional().default(false),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<LoginResponse>("POST", "/users/login", undefined, data)
  })

export const registerFn = createServerFn({ method: "POST" })
  .inputValidator(
    z.object({
      username: z.string(),
      email: z.string().email(),
      password: z.string().min(6),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<number>("POST", "/users", undefined, data)
  })

export const refreshTokenFn = createServerFn({ method: "POST" })
  .inputValidator(
    z.object({
      refreshToken: z.string(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<LoginResponse>("POST", "/users/refresh", {
      customHeaders: { "X-Refresh-Token": data.refreshToken },
    })
  })

export const getOrganizationsFn = createServerFn({ method: "GET" })
  .inputValidator(
    z.object({
      accessToken: z.string(),
      userId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<Array<Domain>>(
      "GET",
      `/users/${data.userId}/organizations`,
      { accessToken: data.accessToken },
    )
  })

export const getAppsFn = createServerFn({ method: "GET" })
  .inputValidator(
    z.object({
      accessToken: z.string(),
      domainId: z.number(),
      appId: z.number(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<Array<App>>(
      "GET",
      `/apps`,
      { accessToken: data.accessToken, domainId: data.domainId, appId: data.appId },
    )
  })

export const getPermissionsFn = createServerFn({ method: "GET" })
  .inputValidator(
    z.object({
      accessToken: z.string(),
      userId: z.number(),
      domainId: z.number().optional(),
      appId: z.number().optional(),
    }),
  )
  .handler(async ({ data }) => {
    return serverApi<Array<Array<string>>>(
      "GET",
      `/users/${data.userId}/permissions`,
      { accessToken: data.accessToken, domainId: data.domainId, appId: data.appId },
    )
  })
