import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { EyeOpenIcon, EyeClosedIcon } from '@radix-ui/react-icons'
import { loginSchema, type LoginValues } from '../schemas/auth.schema'
import { useLogin } from '../hooks/useAuthMutations'
import { Button, Field, Input } from '@/components/ui'

export function LoginForm() {
  const [showPassword, setShowPassword] = useState(false)
  const login = useLogin()

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  })

  const onSubmit = handleSubmit(async (values) => {
    try {
      await login.mutateAsync(values)
    } catch (err) {
      // Field-level errors from a 422 (e.g. malformed payload caught server-side).
      if (err && typeof err === 'object' && 'fields' in err) {
        const fields = (err as { fields?: Record<string, string[]> }).fields
        if (fields) {
          for (const [key, msgs] of Object.entries(fields)) {
            if (key === 'email' || key === 'password') {
              setError(key, { message: msgs[0] })
            }
          }
        }
      }
    }
  })

  const pending = isSubmitting || login.isPending

  return (
    <form onSubmit={onSubmit} noValidate className="flex flex-col gap-4">
      <Field label="E-mail" required error={errors.email?.message}>
        {(props) => (
          <Input
            {...props}
            {...register('email')}
            type="email"
            inputMode="email"
            autoComplete="email"
            placeholder="voce@exemplo.com"
            disabled={pending}
          />
        )}
      </Field>

      <Field label="Senha" required error={errors.password?.message}>
        {(props) => (
          <div className="relative">
            <Input
              {...props}
              {...register('password')}
              type={showPassword ? 'text' : 'password'}
              autoComplete="current-password"
              placeholder="••••••••"
              disabled={pending}
              className="pr-11"
            />
            <button
              type="button"
              aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
              aria-pressed={showPassword}
              onClick={() => setShowPassword((v) => !v)}
              disabled={pending}
              className="absolute right-1 top-1 inline-flex h-9 w-9 items-center justify-center rounded-sm text-ink-muted hover:bg-bg-subtle hover:text-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
            >
              {showPassword ? <EyeClosedIcon /> : <EyeOpenIcon />}
            </button>
          </div>
        )}
      </Field>

      <Button type="submit" loading={pending} className="mt-2 w-full">
        Entrar
      </Button>
    </form>
  )
}
