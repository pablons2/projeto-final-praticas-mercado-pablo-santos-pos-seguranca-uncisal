import { useId, type ReactNode } from 'react'
import { cn } from '@/lib/cn'
import { Label } from './Label'

interface FieldProps {
  label: string
  /** Optional, e.g. "(obrigatório)". */
  hint?: ReactNode
  /** Error message — rendered with role="alert" when present. */
  error?: string
  /** Marks the field as required (adds the asterisk + aria). */
  required?: boolean
  /** Right-aligned slot under the input, e.g. a char counter. */
  footer?: ReactNode
  /** The control. Receives id, aria-describedby, aria-invalid. */
  children: (props: {
    id: string
    'aria-describedby'?: string
    'aria-invalid': boolean | undefined
    invalid: boolean
  }) => ReactNode
  className?: string
}

/**
 * Field — composes label + control + hint + error with correct a11y wiring.
 * The control is rendered via a render-prop so we can inject id/aria
 * without each input component needing to know about Field.
 */
export function Field({
  label,
  hint,
  error,
  required,
  footer,
  children,
  className,
}: FieldProps) {
  const id = useId()
  const hintId = `${id}-hint`
  const errorId = `${id}-error`
  const describedBy = [error ? errorId : null, hint ? hintId : null]
    .filter(Boolean)
    .join(' ') || undefined

  return (
    <div className={cn('flex flex-col', className)}>
      <Label htmlFor={id}>
        {label}
        {required ? <span className="text-danger ml-0.5" aria-hidden="true">*</span> : null}
      </Label>
      {children({
        id,
        'aria-describedby': describedBy,
        'aria-invalid': error ? true : undefined,
        invalid: Boolean(error),
      })}
      {hint ? (
        <p id={hintId} className="mt-1.5 text-xs text-ink-muted">
          {hint}
        </p>
      ) : null}
      {error ? (
        <p id={errorId} role="alert" className="mt-1.5 text-xs text-danger">
          {error}
        </p>
      ) : null}
      {footer ? <div className="mt-1.5 text-right text-xs text-ink-muted">{footer}</div> : null}
    </div>
  )
}
