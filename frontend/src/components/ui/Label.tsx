import { forwardRef, type LabelHTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

export const Label = forwardRef<HTMLLabelElement, LabelHTMLAttributes<HTMLLabelElement>>(
  function Label({ className, ...rest }, ref) {
    return (
      <label
        ref={ref}
        className={cn(
          'block text-sm font-medium text-ink-secondary mb-1.5',
          className,
        )}
        {...rest}
      />
    )
  },
)
