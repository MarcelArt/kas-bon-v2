import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { useRouter } from "@tanstack/react-router"
import { toast } from "sonner"
import type { LoginInput } from "@/lib/api.schemas"
import type { User, UserInput } from "@/lib/api.types"
import {
  getCurrentUserFn,
  getUserPermissionsFn,
  loginFn,
  logoutFn,
  registerFn,
} from "@/lib/auth.fns"

export function useCurrentUser() {
  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: () => getCurrentUserFn(),
  })
}

export function useUserPermissions(userId: number | undefined) {
  return useQuery({
    queryKey: ["auth", "permissions", userId],
    queryFn: () =>
      getUserPermissionsFn({
        data: { userId: userId!, appId: 1, domainId: 1 },
      }),
    enabled: userId !== undefined,
  })
}

export function useLoginMutation() {
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation({
    mutationFn: (data: LoginInput) => loginFn({ data }),
    onSuccess: async (user: User) => {
      queryClient.setQueryData(["auth", "me"], user)
      await queryClient.invalidateQueries({ queryKey: ["auth"] })
      await router.navigate({ to: "/dashboard" })
    },
    onError: (error: Error) => {
      toast.error(error.message || "Login failed")
    },
  })
}

export function useRegisterMutation() {
  const router = useRouter()

  return useMutation({
    mutationFn: (data: UserInput) => registerFn({ data }),
    onSuccess: () => {
      toast.success("Registration successful! Please sign in.")
      router.navigate({ to: "/login" })
    },
    onError: (error: Error) => {
      toast.error(error.message || "Registration failed")
    },
  })
}

export function useLogoutMutation() {
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation({
    mutationFn: () => logoutFn(),
    onSuccess: async () => {
      await queryClient.clear()
      await router.navigate({ to: "/login" })
    },
    onError: (error: Error) => {
      toast.error(error.message || "Logout failed")
    },
  })
}
