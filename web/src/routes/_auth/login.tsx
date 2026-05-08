import { useForm } from "@tanstack/react-form"
import { Link, createFileRoute } from "@tanstack/react-router"
import type { AnyFieldApi } from "@tanstack/react-form"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { loginInputSchema } from "@/lib/api.schemas"
import { useLoginMutation } from "@/lib/auth.query"

export const Route = createFileRoute("/_auth/login")({
  component: LoginPage,
})

function FieldInfo({ field }: { field: AnyFieldApi }) {
  return (
    <>
      {field.state.meta.isTouched && field.state.meta.errors.length > 0 ? (
        <em className="text-xs text-destructive">
          {field.state.meta.errors.map((err) => err.message).join(", ")}
        </em>
      ) : null}
    </>
  )
}

function LoginPage() {
  const loginMutation = useLoginMutation()

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
      isRemember: false,
    },
    validators: {
      onChange: loginInputSchema,
    },
    onSubmit: ({ value }) => {
      loginMutation.mutate(value)
    },
  })

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Login</CardTitle>
        <CardDescription>Enter your credentials to access the app</CardDescription>
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
          <form.Field
            name="username"
            children={(field) => (
              <div className="flex flex-col gap-2">
                <Label htmlFor={field.name}>Username</Label>
                <Input
                  id={field.name}
                  type="text"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  autoComplete="username"
                />
                <FieldInfo field={field} />
              </div>
            )}
          />
          <form.Field
            name="password"
            children={(field) => (
              <div className="flex flex-col gap-2">
                <Label htmlFor={field.name}>Password</Label>
                <Input
                  id={field.name}
                  type="password"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  autoComplete="current-password"
                />
                <FieldInfo field={field} />
              </div>
            )}
          />
          <form.Field
            name="isRemember"
            children={(field) => (
              <div className="flex items-center gap-2">
                <Checkbox
                  id="remember"
                  checked={field.state.value}
                  onCheckedChange={(checked) => field.handleChange(checked === true)}
                />
                <Label htmlFor="remember" className="cursor-pointer">
                  Remember me
                </Label>
              </div>
            )}
          />
          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <Button
                type="submit"
                className="w-full"
                disabled={!canSubmit || isSubmitting || loginMutation.isPending}
              >
                {isSubmitting || loginMutation.isPending
                  ? "Signing in..."
                  : "Sign in"}
              </Button>
            )}
          />
          <p className="text-center text-xs text-muted-foreground">
            Don&apos;t have an account?{" "}
            <Link
              to="/register"
              className="text-primary underline-offset-4 hover:underline"
            >
              Register
            </Link>
          </p>
        </form>
      </CardContent>
    </Card>
  )
}
