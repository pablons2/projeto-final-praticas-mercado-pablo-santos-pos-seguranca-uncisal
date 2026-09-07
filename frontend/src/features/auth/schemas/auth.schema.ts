import { z } from 'zod'

/**
 * Auth validation schemas — mirror the backend policy (docs/plano.md §5):
 * password >= 8 chars with at least one letter and one number.
 * Client-side validation is for fast feedback only; the server is
 * authoritative.
 */

const passwordPolicy = z
  .string()
  .min(8, 'Mínimo de 8 caracteres')
  .regex(/[a-zA-Z]/, 'Inclua ao menos uma letra')
  .regex(/[0-9]/, 'Inclua ao menos um número')

export const loginSchema = z.object({
  email: z.string().min(1, 'Informe o e-mail').email('E-mail inválido'),
  password: z.string().min(1, 'Informe a senha'),
})

export const registerSchema = z
  .object({
    name: z
      .string()
      .min(2, 'Nome muito curto')
      .max(80, 'Nome muito longo')
      .transform((v) => v.trim()),
    email: z.string().min(1, 'Informe o e-mail').email('E-mail inválido'),
    password: passwordPolicy,
    confirmPassword: z.string().min(1, 'Confirme a senha'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'As senhas não coincidem',
    path: ['confirmPassword'],
  })

export type LoginValues = z.infer<typeof loginSchema>
export type RegisterValues = z.infer<typeof registerSchema>

/** Deterministic rule list for the live password-strength hint. */
export const PASSWORD_RULES: { test: (v: string) => boolean; label: string }[] = [
  { test: (v) => v.length >= 8, label: 'Ao menos 8 caracteres' },
  { test: (v) => /[a-zA-Z]/.test(v), label: 'Ao menos uma letra' },
  { test: (v) => /[0-9]/.test(v), label: 'Ao menos um número' },
]
