import { PlusIcon } from '@radix-ui/react-icons'

interface EmptyStateProps {
  title: string
  description: string
  actionLabel?: string
  onAction?: () => void
}

/** Friendly, illustration-free empty state. */
export function EmptyState({ title, description, actionLabel, onAction }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
      <span
        className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-bg-subtle text-ink-muted"
        aria-hidden="true"
      >
        <PlusIcon className="h-6 w-6" />
      </span>
      <div>
        <h3 className="text-base font-medium text-ink">{title}</h3>
        <p className="mt-1 text-sm text-ink-secondary">{description}</p>
      </div>
      {actionLabel && onAction ? (
        <button
          type="button"
          onClick={onAction}
          className="mt-2 inline-flex items-center gap-1.5 rounded-sm bg-accent px-4 h-11 text-sm font-medium text-white hover:bg-accent-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring focus-visible:ring-offset-2"
        >
          <PlusIcon className="h-4 w-4" />
          {actionLabel}
        </button>
      ) : null}
    </div>
  )
}
