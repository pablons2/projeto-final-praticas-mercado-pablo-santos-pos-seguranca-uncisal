import { forwardRef } from 'react'
import * as SelectPrimitive from '@radix-ui/react-select'
import { CheckIcon, ChevronDownIcon, ChevronUpIcon } from '@radix-ui/react-icons'
import { cn } from '@/lib/cn'

export interface SelectOption {
  value: string
  label: string
}

interface SelectProps {
  options: SelectOption[]
  value?: string
  onChange: (value: string) => void
  placeholder?: string
  invalid?: boolean
  disabled?: boolean
  ariaLabel?: string
  id?: string
  className?: string
}

/**
 * Minimal Radix Select — styled to match Input height/border.
 * Uses native-feeling chevron + check, no custom scroll on the viewport.
 */
export const Select = forwardRef<HTMLButtonElement, SelectProps>(function Select(
  {
    options,
    value,
    onChange,
    placeholder = 'Selecione…',
    invalid,
    disabled,
    ariaLabel,
    id,
    className,
  },
  ref,
) {
  return (
    <SelectPrimitive.Root
      value={value}
      onValueChange={onChange}
      disabled={disabled}
    >
      <SelectPrimitive.Trigger
        ref={ref}
        id={id}
        aria-label={ariaLabel}
        aria-invalid={invalid || undefined}
        className={cn(
          'inline-flex w-full items-center justify-between rounded-sm border bg-bg px-3',
          'h-11 text-sm text-ink min-h-touch',
          'transition-colors duration-150',
          'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring focus-visible:ring-offset-1',
          'disabled:cursor-not-allowed disabled:opacity-60',
          'data-[placeholder]:text-ink-muted',
          invalid ? 'border-danger' : 'border-border-strong',
          className,
        )}
      >
        <SelectPrimitive.Value placeholder={placeholder} />
        <SelectPrimitive.Icon className="ml-2 text-ink-muted">
          <ChevronDownIcon />
        </SelectPrimitive.Icon>
      </SelectPrimitive.Trigger>
      <SelectPrimitive.Portal>
        <SelectPrimitive.Content
          position="popper"
          sideOffset={4}
          className={cn(
            'overflow-hidden rounded-md border border-border bg-bg shadow-overlay z-50',
            'max-h-[var(--radix-select-content-available-height)] min-w-[var(--radix-select-trigger-width)]',
            'animate-fade-in',
          )}
        >
          <SelectPrimitive.ScrollUpButton className="flex h-6 items-center justify-center text-ink-muted">
            <ChevronUpIcon />
          </SelectPrimitive.ScrollUpButton>
          <SelectPrimitive.Viewport className="p-1">
            {options.map((opt) => (
              <SelectPrimitive.Item
                key={opt.value}
                value={opt.value}
                className={cn(
                  'relative flex h-9 cursor-pointer select-none items-center rounded-sm px-2 pl-8',
                  'text-sm text-ink outline-none',
                  'data-[highlighted]:bg-accent-subtle data-[highlighted]:text-ink',
                  'data-[state=checked]:font-medium',
                )}
              >
                <span className="absolute left-2 flex h-4 w-4 items-center justify-center">
                  <SelectPrimitive.ItemIndicator>
                    <CheckIcon />
                  </SelectPrimitive.ItemIndicator>
                </span>
                <SelectPrimitive.ItemText>{opt.label}</SelectPrimitive.ItemText>
              </SelectPrimitive.Item>
            ))}
          </SelectPrimitive.Viewport>
          <SelectPrimitive.ScrollDownButton className="flex h-6 items-center justify-center text-ink-muted">
            <ChevronDownIcon />
          </SelectPrimitive.ScrollDownButton>
        </SelectPrimitive.Content>
      </SelectPrimitive.Portal>
    </SelectPrimitive.Root>
  )
})
