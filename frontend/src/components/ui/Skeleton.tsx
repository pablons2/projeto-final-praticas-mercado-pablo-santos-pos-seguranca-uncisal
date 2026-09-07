import { cn } from '@/lib/cn'

interface SkeletonProps {
  className?: string
  /** When true, renders a rounded-full pill instead of a rounded-sm block. */
  circle?: boolean
}

/** Shimmer-free skeleton — flat bg-subtle with a subtle pulse. */
export function Skeleton({ className, circle }: SkeletonProps) {
  return (
    <div
      aria-hidden="true"
      className={cn(
        'bg-bg-subtle animate-pulse',
        circle ? 'rounded-full' : 'rounded-sm',
        className,
      )}
    />
  )
}
