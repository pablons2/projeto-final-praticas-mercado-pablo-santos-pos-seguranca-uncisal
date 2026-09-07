import { type ReactNode } from 'react'
import { BrandMark } from '@/components/BrandMark'

/**
 * Centered single-column card layout for auth screens.
 * Mobile: full-bleed with safe-area padding.
 * >=768px: vertically centered card, max-w-form.
 */
export function AuthLayout({
  title,
  subtitle,
  children,
  footer,
}: {
  title: string
  subtitle: string
  children: ReactNode
  footer: ReactNode
}) {
  return (
    <div className="flex min-h-dvh flex-col items-center justify-center bg-bg-surface px-4 py-8 pt-safe">
      <div className="w-full max-w-form">
        <div className="mb-6 flex justify-center">
          <BrandMark size="md" />
        </div>
        <div className="rounded-lg border border-border bg-bg p-6 shadow-sm sm:p-8">
          <h1 className="text-2xl font-semibold text-ink">{title}</h1>
          <p className="mt-1 text-sm text-ink-secondary">{subtitle}</p>
          <div className="mt-6">{children}</div>
        </div>
        <div className="mt-4 text-center text-sm text-ink-secondary">{footer}</div>
      </div>
    </div>
  )
}
