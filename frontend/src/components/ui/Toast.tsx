import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import * as ToastPrimitive from '@radix-ui/react-toast'
import { Cross2Icon, CheckCircledIcon, ExclamationTriangleIcon, InfoCircledIcon } from '@radix-ui/react-icons'
import { cn } from '@/lib/cn'

type ToastVariant = 'success' | 'error' | 'info'

interface ToastInput {
  title: string
  description?: string
  variant?: ToastVariant
  duration?: number
}

interface ToastContextValue {
  toast: (input: ToastInput) => void
}

const ToastContext = createContext<ToastContextValue | null>(null)

type ToastRecord = {
  id: string
  title: string
  description?: string
  variant: ToastVariant
  duration: number
}

const VARIANT_STYLES: Record<ToastVariant, { icon: ReactNode; ring: string }> = {
  success: {
    icon: <CheckCircledIcon className="h-5 w-5 text-sev-baixa" />,
    ring: 'border-l-2 border-l-sev-baixa',
  },
  error: {
    icon: <ExclamationTriangleIcon className="h-5 w-5 text-danger" />,
    ring: 'border-l-2 border-l-danger',
  },
  info: {
    icon: <InfoCircledIcon className="h-5 w-5 text-accent" />,
    ring: 'border-l-2 border-l-accent',
  },
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastRecord[]>([])

  const toast = useCallback((input: ToastInput) => {
    const id = Math.random().toString(36).slice(2)
    setToasts((prev) => [
      ...prev,
      {
        id,
        title: input.title,
        description: input.description,
        variant: input.variant ?? 'info',
        duration: input.duration ?? 5000,
      },
    ])
  }, [])

  const remove = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const value = useMemo<ToastContextValue>(() => ({ toast }), [toast])

  return (
    <ToastContext.Provider value={value}>
      <ToastPrimitive.Provider swipeDirection="right" duration={5000}>
        {children}
        {toasts.map((t) => {
          const v = VARIANT_STYLES[t.variant]
          return (
            <ToastPrimitive.Root
              key={t.id}
              duration={t.duration}
              onOpenChange={(open) => {
                if (!open) remove(t.id)
              }}
              className={cn(
                'pointer-events-auto grid grid-cols-[auto_1fr_auto] items-start gap-3',
                'rounded-md border border-border bg-bg p-4 shadow-overlay',
                'data-[state=open]:animate-toast-in',
                'data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)]',
                'data-[swipe=cancel]:translate-x-0',
                'data-[swipe=end]:translate-x-[var(--radix-toast-swipe-end-x)]',
                v.ring,
              )}
            >
              <ToastPrimitive.Title className="flex items-start gap-3 text-sm font-medium text-ink">
                {v.icon}
              </ToastPrimitive.Title>
              <div className="flex flex-col gap-0.5">
                <ToastPrimitive.Title className="text-sm font-medium text-ink">
                  {t.title}
                </ToastPrimitive.Title>
                {t.description ? (
                  <ToastPrimitive.Description className="text-xs text-ink-secondary">
                    {t.description}
                  </ToastPrimitive.Description>
                ) : null}
              </div>
              <ToastPrimitive.Close
                aria-label="Fechar"
                className="rounded-sm p-1 text-ink-muted hover:bg-bg-subtle hover:text-ink"
              >
                <Cross2Icon />
              </ToastPrimitive.Close>
            </ToastPrimitive.Root>
          )
        })}
        {/* Viewport: top-center on mobile, bottom-right on >=768px */}
        <ToastPrimitive.Viewport
          className={cn(
            'fixed z-[60] flex flex-col gap-2 w-[calc(100%-1.5rem)] max-w-sm',
            'top-3 left-1/2 -translate-x-1/2',
            'md:top-auto md:bottom-3 md:left-auto md:right-3 md:translate-x-0',
          )}
        />
      </ToastPrimitive.Provider>
    </ToastContext.Provider>
  )
}

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext)
  if (!ctx) throw new Error('useToast must be used within <ToastProvider>')
  return ctx
}
