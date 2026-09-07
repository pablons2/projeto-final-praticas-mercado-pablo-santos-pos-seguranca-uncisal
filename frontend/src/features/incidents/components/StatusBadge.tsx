import { STATUS_LABELS, type Status } from '@/lib/constants'
import { cn } from '@/lib/cn'

const DOT: Record<Status, string> = {
  aberto: 'bg-stat-aberto',
  em_analise: 'bg-stat-em_analise',
  resolvido: 'bg-stat-resolvido',
}

/** Status badge — colored dot + text label. */
export function StatusBadge({ value }: { value: Status }) {
  return (
    <span className="inline-flex items-center gap-1.5 text-sm text-ink-secondary">
      <span className={cn('inline-block h-2 w-2 rounded-full', DOT[value])} aria-hidden="true" />
      {STATUS_LABELS[value]}
    </span>
  )
}
