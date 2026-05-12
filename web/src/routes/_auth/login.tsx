import { Link, createFileRoute } from "@tanstack/react-router"
import { useForm } from "@tanstack/react-form"
import { toast } from "sonner"
import type { AnyFieldApi } from "@tanstack/react-form"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useLogin } from "@/lib/queries/auth"
import { loginSchema } from "@/lib/schemas/auth"

export const Route = createFileRoute("/_auth/login")({
  component: LoginPage,
})

function FieldInfo({ field }: { field: AnyFieldApi }) {
  return (
    <>
      {field.state.meta.isTouched && !field.state.meta.isValid
        ? field.state.meta.errors.map((err) => (
            <em key={String(err)} className="text-destructive text-xs">
              {String(err.message ?? err)}
            </em>
          ))
        : null}
    </>
  )
}

function LoginPage() {
  const login = useLogin()

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
      isRemember: false,
    },
    validators: {
      onChange: loginSchema,
    },
    onSubmit: ({ value }) => {
      login.mutate(value, {
        onError: (error) => {
          toast.error(error.message)
        },
      })
    },
  })

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle className="text-2xl">Login</CardTitle>
        <CardDescription>Enter your credentials to access the system.</CardDescription>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault()
            e.stopPropagation()
            form.handleSubmit()
          }}
          className="flex flex-col gap-4"
        >
          <form.Field name="username">
            {(field) => (
              <div className="flex flex-col gap-2">
                <Label htmlFor={field.name}>Username</Label>
                <Input
                  id={field.name}
                  name={field.name}
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
                <FieldInfo field={field} />
              </div>
            )}
          </form.Field>

          <form.Field name="password">
            {(field) => (
              <div className="flex flex-col gap-2">
                <Label htmlFor={field.name}>Password</Label>
                <Input
                  id={field.name}
                  name={field.name}
                  type="password"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
                <FieldInfo field={field} />
              </div>
            )}
          </form.Field>

          <form.Field name="isRemember">
            {(field) => (
              <div className="flex items-center gap-2">
                <Checkbox
                  id={field.name}
                  checked={field.state.value}
                  onCheckedChange={(checked) => field.handleChange(checked === true)}
                  onBlur={field.handleBlur}
                />
                <Label htmlFor={field.name} className="text-sm">
                  Remember me
                </Label>
              </div>
            )}
          </form.Field>

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
          >
            {([canSubmit, isSubmitting]) => (
              <Button type="submit" className="w-full" disabled={!canSubmit || login.isPending}>
                {isSubmitting || login.isPending ? "Logging in..." : "Login"}
              </Button>
            )}
          </form.Subscribe>
        </form>

        <div className="mt-4 text-center text-sm">
          Don&apos;t have an account?{" "}
          <Link to="/register" className="text-primary underline">
            Register
          </Link>
        </div>
      </CardContent>
    </Card>
  )
}
