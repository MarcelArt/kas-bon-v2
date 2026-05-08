import { z } from "zod"

export const loginInputSchema = z.object({
  username: z.string({ error: "Username is required" }),
  password: z.string({ error: "Password is required" }),
  isRemember: z.boolean(),
})

export const registerInputSchema = z
  .object({
    username: z.string().min(3, "Username must be at least 3 characters"),
    email: z.email("Invalid email address"),
    password: z.string().min(6, "Password must be at least 6 characters"),
    confirmPassword: z.string({ error: "Please confirm your password" }),
  })
  .refine((data) => data.password === data.confirmPassword, {
    error: "Passwords do not match",
    path: ["confirmPassword"],
  })

export const userInputSchema = z.object({
  username: z.string({ error: "Username is required" }),
  email: z.email("Invalid email address"),
  password: z.string({ error: "Password is required" }),
})

export const roleInputSchema = z.object({
  name: z.string({ error: "Name is required" }),
  description: z.string({ error: "Description is required" }),
  domainId: z.number({ error: "Domain ID is required" }),
})

export const permissionInputSchema = z.object({
  name: z.string({ error: "Name is required" }),
  description: z.string({ error: "Description is required" }),
  appId: z.number({ error: "App ID is required" }),
})

export const domainInputSchema = z.object({
  name: z.string({ error: "Name is required" }),
  description: z.string({ error: "Description is required" }),
  isOrganization: z.boolean(),
  parentId: z.number().nullable(),
})

export const appInputSchema = z.object({
  name: z.string({ error: "Name is required" }),
  description: z.string({ error: "Description is required" }),
})

export const loginResponseSchema = z.object({
  accessToken: z.string(),
  refreshToken: z.string(),
  user: z.object({
    id: z.number(),
    username: z.string(),
    email: z.string(),
    createdAt: z.string(),
    updatedAt: z.string(),
    deletedAt: z.string().nullable(),
  }),
})

export type LoginInput = z.infer<typeof loginInputSchema>
export type RegisterInput = z.infer<typeof registerInputSchema>
