import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { EyeOpenIcon, EyeClosedIcon } from '@radix-ui/react-icons'
import { registerSchema, type RegisterValues } from '../schemas/auth.schema'
import { useRegister } from '../hooks/useAuthMutations'
import { PasswordStrengthHint } from './PasswordStrengthHint'
import { Button, Field, Input } from '@/components/ui'

export function RegisterForm() {
  const [showPassword, setShowPassword] = useState(false)
  const registerMutation = useRegister()

  const {
    register,
    handleSubmit,
    setError,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: '', email: '', password: '', confirmPassword: '' },
  })

  const passwordValue = watch('password')

  const onSubmit = handleSubmit(async (values) => {
    try {
      await registerMutation.mutateAsync(values)
    } catch (err) {
      if (err && typeof err === 'object' && 'fields' in err) {
        const fields = (err as { fields?: Record<string, string[]> }).fields
        if (fields) {
          for (const [key, msgs] of Object.entries(fields)) {
            if (
              key === 'name' ||
              key === 'email' ||
              key === 'password' ||
              key === 'confirmPassword'
            ) {
              setError(key, { message: msgs[0] })
            }
          }
        }
      }
    }
  })

  const pending = isSubmitting || registerMutation.isPending

  return (
    <form onSubmit={onSubmit} noValidate className="flex flex-col gap-4">
      <Field label="Nome" required error={errors.name?.message}>
        {(props) => (
          <Input
            {...props}
            {...register('name')}
            type="text"
            autoComplete="name"
            placeholder="Seu nome"
            disabled={pending}
          />
        )}
      </Field>

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

      <Field
        label="Senha"
        required
        error={errors.password?.message}
        hint={<PasswordStrengthHint value={passwordValue ?? ''} />}
      >
        {(props) => (
          <div className="relative">
            <Input
              {...props}
              {...register('password')}
              type={showPassword ? 'text' : 'password'}
              autoComplete="new-password"
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

      <Field label="Confirmar senha" required error={errors.confirmPassword?.message}>
        {(props) => (
          <Input
            {...props}
            {...register('confirmPassword')}
            type={showPassword ? 'text' : 'password'}
            autoComplete="new-password"
            placeholder="••••••••"
            disabled={pending}
          />
        )}
      </Field>

      <Button type="submit" loading={pending} className="mt-2 w-full">
        Criar conta
      </Button>
    </form>
  )
}
