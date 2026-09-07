import { SEVERIDADE_LABELS, type Severidade } from '@/lib/constants'
import { cn } from '@/lib/cn'

const DOT: Record<Severidade, string> = {
  baixa: 'bg-sev-baixa',
  media: 'bg-sev-media',
  alta: 'bg-sev-alta',
  critica: 'bg-sev-critica',
}

/** Severity badge — colored dot + text label. Color is never the only signal. */
export function SeverityBadge({ value }: { value: Severidade }) {
  return (
    <span className="inline-flex items-center gap-1.5 text-sm text-ink-secondary">
      <span className={cn('inline-block h-2 w-2 rounded-full', DOT[value])} aria-hidden="true" />
      {SEVERIDADE_LABELS[value]}
    </span>
  )
}
