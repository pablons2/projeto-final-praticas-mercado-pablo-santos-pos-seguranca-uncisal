import { forwardRef, type ReactNode } from 'react'
import * as DialogPrimitive from '@radix-ui/react-dialog'
import { Cross2Icon } from '@radix-ui/react-icons'
import { cn } from '@/lib/cn'

/**
 * Dialog — two variants:
 *  - "modal" (default): centered card, max-w-dialog, used on desktop.
 *  - "sheet": full-width, slides from bottom, used on mobile for forms.
 * The variant is chosen by the caller based on a media query.
 */

export const Dialog = DialogPrimitive.Root
export const DialogTrigger = DialogPrimitive.Trigger
export const DialogClose = DialogPrimitive.Close

interface DialogContentProps {
  title: string
  description?: string
  children: ReactNode
  /** "modal" | "sheet" — controls layout. */
  variant?: 'modal' | 'sheet'
  /** Footer slot, typically action buttons. */
  footer?: ReactNode
  /** Disable the X close button (e.g. while submitting). */
  hideClose?: boolean
  className?: string
}

export const DialogContent = forwardRef<HTMLDivElement, DialogContentProps>(
  function DialogContent(
    { title, description, children, variant = 'modal', footer, hideClose, className },
    ref,
  ) {
    return (
      <DialogPrimitive.Portal>
        <DialogPrimitive.Overlay
          className={cn(
            'fixed inset-0 z-40 bg-ink/40 backdrop-blur-[1px]',
            'animate-fade-in',
          )}
        />
        <DialogPrimitive.Content
          ref={ref}
          className={cn(
            'fixed z-50 bg-bg border border-border shadow-overlay',
            'flex flex-col outline-none',
            'data-[state=open]:animate-slide-up',
            variant === 'modal' &&
              'left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[calc(100%-2rem)] max-w-dialog max-h-[85vh] rounded-lg',
            variant === 'sheet' &&
              'inset-x-0 bottom-0 w-full max-h-[92vh] rounded-t-lg pb-safe data-[state=open]:animate-sheet-up',
            className,
          )}
        >
          <div
            className={cn(
              'flex items-start justify-between gap-4 border-b border-border px-5 py-4',
              variant === 'sheet' && 'pt-safe',
            )}
          >
            <div className="min-w-0">
              <DialogPrimitive.Title className="text-lg font-semibold text-ink leading-tight">
                {title}
              </DialogPrimitive.Title>
              {description ? (
                <DialogPrimitive.Description className="mt-1 text-sm text-ink-secondary">
                  {description}
                </DialogPrimitive.Description>
              ) : null}
            </div>
            {hideClose ? null : (
              <DialogPrimitive.Close
                aria-label="Fechar"
                className={cn(
                  'shrink-0 rounded-sm p-1 text-ink-muted',
                  'hover:bg-bg-subtle hover:text-ink transition-colors',
                  'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring',
                )}
              >
                <Cross2Icon />
              </DialogPrimitive.Close>
            )}
          </div>
          <div className="flex-1 overflow-y-auto px-5 py-4">{children}</div>
          {footer ? (
            <div className="flex items-center justify-end gap-2 border-t border-border px-5 py-3">
              {footer}
            </div>
          ) : null}
        </DialogPrimitive.Content>
      </DialogPrimitive.Portal>
    )
  },
)
