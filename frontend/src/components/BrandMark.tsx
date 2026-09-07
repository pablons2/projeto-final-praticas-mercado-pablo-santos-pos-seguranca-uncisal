import { cn } from '@/lib/cn'

/** Wordmark — no image asset, keeps the bundle tiny. */
export function BrandMark({
  className,
  size = 'md',
}: {
  className?: string
  size?: 'sm' | 'md'
}) {
  return (
    <div className={cn('flex items-center gap-2 select-none', className)}>
      <span
        className={cn(
          'inline-flex items-center justify-center rounded-md bg-accent text-white',
          size === 'sm' ? 'h-7 w-7' : 'h-8 w-8',
        )}
        aria-hidden="true"
      >
        <svg viewBox="0 0 24 24" className={size === 'sm' ? 'h-4 w-4' : 'h-5 w-5'} fill="none">
          <path
            d="M12 4a7 7 0 0 0-7 7v3a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-3a7 7 0 0 0-7-7Z"
            stroke="currentColor"
            strokeWidth="1.6"
          />
          <circle cx="12" cy="10.5" r="1.2" fill="currentColor" />
          <path d="M12 11.8v2.2" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
        </svg>
      </span>
      <span
        className={cn(
          'font-semibold tracking-tight text-ink',
          size === 'sm' ? 'text-sm' : 'text-base',
        )}
      >
        IncidentTrack
      </span>
    </div>
  )
}
