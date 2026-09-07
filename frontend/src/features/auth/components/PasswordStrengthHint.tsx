import { PASSWORD_RULES } from '../schemas/auth.schema'
import { cn } from '@/lib/cn'

/** Live, deterministic password-strength hint — not a score meter. */
export function PasswordStrengthHint({ value }: { value: string }) {
  return (
    <ul className="mt-2 space-y-1" aria-live="polite">
      {PASSWORD_RULES.map((rule) => {
        const ok = rule.test(value)
        return (
          <li
            key={rule.label}
            className={cn(
              'flex items-center gap-1.5 text-xs',
              ok ? 'text-sev-baixa' : 'text-ink-muted',
            )}
          >
            <span
              className={cn(
                'inline-block h-1.5 w-1.5 rounded-full',
                ok ? 'bg-sev-baixa' : 'bg-border-strong',
              )}
              aria-hidden="true"
            />
            {rule.label}
          </li>
        )
      })}
    </ul>
  )
}
