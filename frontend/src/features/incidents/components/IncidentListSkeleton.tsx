import { Skeleton } from '@/components/ui'

/** Desktop table skeleton — 5 rows. */
export function IncidentTableSkeleton() {
  return (
    <div className="overflow-hidden rounded-md border border-border">
      <div className="border-b border-border bg-bg-subtle px-4 py-2.5">
        <Skeleton className="h-4 w-32" />
      </div>
      <div className="divide-y divide-border-hairline">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="flex items-center gap-4 px-4 py-3">
            <Skeleton className="h-4 flex-1" />
            <Skeleton className="h-5 w-24" />
            <Skeleton className="h-5 w-16" />
            <Skeleton className="h-5 w-20" />
            <Skeleton className="h-4 w-16" />
            <Skeleton className="h-9 w-20" />
          </div>
        ))}
      </div>
    </div>
  )
}

/** Mobile card skeleton — 4 cards. */
export function IncidentCardListSkeleton() {
  return (
    <ul className="flex flex-col gap-3" role="list" aria-busy="true">
      {Array.from({ length: 4 }).map((_, i) => (
        <li key={i} className="rounded-md border border-border bg-bg p-4 shadow-sm">
          <div className="flex items-start justify-between gap-3">
            <Skeleton className="h-4 flex-1" />
            <Skeleton className="h-9 w-20" />
          </div>
          <div className="mt-3 flex gap-2">
            <Skeleton className="h-5 w-24" />
            <Skeleton className="h-5 w-16" />
          </div>
          <Skeleton className="mt-3 h-4 w-full" />
          <Skeleton className="mt-1 h-4 w-2/3" />
        </li>
      ))}
    </ul>
  )
}
