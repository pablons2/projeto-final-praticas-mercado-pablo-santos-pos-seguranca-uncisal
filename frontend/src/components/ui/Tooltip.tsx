import { forwardRef, type ReactNode } from 'react'
import * as TooltipPrimitive from '@radix-ui/react-tooltip'
import { cn } from '@/lib/cn'

export const TooltipProvider = TooltipPrimitive.Provider

interface TooltipProps {
  children: ReactNode
  content: ReactNode
  side?: 'top' | 'right' | 'bottom' | 'left'
}

/** Minimal tooltip — used for icon-button labels on the dashboard table. */
export const Tooltip = forwardRef<HTMLDivElement, TooltipProps>(function Tooltip(
  { children, content, side = 'top' },
  ref,
) {
  return (
    <TooltipPrimitive.Root delayDuration={200}>
      <TooltipPrimitive.Trigger asChild>
        <span ref={ref} className="inline-flex">
          {children}
        </span>
      </TooltipPrimitive.Trigger>
      <TooltipPrimitive.Portal>
        <TooltipPrimitive.Content
          side={side}
          sideOffset={6}
          className={cn(
            'z-50 rounded-sm bg-ink px-2 py-1 text-xs text-white shadow-overlay',
            'animate-fade-in pointer-events-none select-none',
            'max-w-[220px]',
          )}
        >
          {content}
          <TooltipPrimitive.Arrow className="fill-ink" />
        </TooltipPrimitive.Content>
      </TooltipPrimitive.Portal>
    </TooltipPrimitive.Root>
  )
})
