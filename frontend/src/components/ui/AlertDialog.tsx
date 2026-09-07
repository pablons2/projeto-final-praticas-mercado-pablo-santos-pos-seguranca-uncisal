import { forwardRef, type ReactNode } from 'react'
import * as AlertDialogPrimitive from '@radix-ui/react-alert-dialog'
import { cn } from '@/lib/cn'
import { Button } from './Button'

export const AlertDialog = AlertDialogPrimitive.Root
export const AlertDialogTrigger = AlertDialogPrimitive.Trigger

interface AlertDialogContentProps {
  title: string
  description: ReactNode
  /** Label for the cancel button (safe default — receives initial focus). */
  cancelLabel?: string
  /** Label for the confirm button. */
  confirmLabel: string
  onConfirm: () => void
  /** Disable confirm + show loading (e.g. during mutation). */
  loading?: boolean
  /** Destructive styling for the confirm button. */
  destructive?: boolean
}

export const AlertDialogContent = forwardRef<
  HTMLDivElement,
  AlertDialogContentProps
>(function AlertDialogContent(
  {
    title,
    description,
    cancelLabel = 'Cancelar',
    confirmLabel,
    onConfirm,
    loading,
    destructive,
  },
  ref,
) {
  return (
    <AlertDialogPrimitive.Portal>
      <AlertDialogPrimitive.Overlay
        className="fixed inset-0 z-40 bg-ink/40 backdrop-blur-[1px] animate-fade-in"
      />
      <AlertDialogPrimitive.Content
        ref={ref}
        className={cn(
          'fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50',
          'w-[calc(100%-2rem)] max-w-[420px] rounded-lg bg-bg border border-border shadow-overlay',
          'p-5 outline-none data-[state=open]:animate-slide-up',
        )}
      >
        <AlertDialogPrimitive.Title className="text-lg font-semibold text-ink">
          {title}
        </AlertDialogPrimitive.Title>
        <AlertDialogPrimitive.Description className="mt-2 text-sm text-ink-secondary">
          {description}
        </AlertDialogPrimitive.Description>
        <div className="mt-5 flex items-center justify-end gap-2">
          <AlertDialogPrimitive.Cancel asChild>
            {/* Cancel receives initial focus — safe default per a11y best practice. */}
            <Button variant="secondary" autoFocus>
              {cancelLabel}
            </Button>
          </AlertDialogPrimitive.Cancel>
          <AlertDialogPrimitive.Action asChild>
            <Button
              variant={destructive ? 'danger' : 'primary'}
              loading={loading}
              onClick={(e) => {
                e.preventDefault()
                onConfirm()
              }}
            >
              {confirmLabel}
            </Button>
          </AlertDialogPrimitive.Action>
        </div>
      </AlertDialogPrimitive.Content>
    </AlertDialogPrimitive.Portal>
  )
})
