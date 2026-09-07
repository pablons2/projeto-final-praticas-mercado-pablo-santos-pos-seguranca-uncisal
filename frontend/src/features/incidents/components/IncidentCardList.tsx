import { Pencil1Icon, TrashIcon } from '@radix-ui/react-icons'
import { SeverityBadge } from './SeverityBadge'
import { StatusBadge } from './StatusBadge'
import { CategoriaChip } from './CategoriaChip'
import { relativeTime } from '@/lib/format'
import type { Incident } from '@/lib/api-types'
import type { Categoria, Severidade, Status } from '@/lib/constants'

interface IncidentCardListProps {
  incidents: Incident[]
  onEdit: (incident: Incident) => void
  onDelete: (incident: Incident) => void
}

/** Mobile list — stacked cards with expandable details. */
export function IncidentCardList({ incidents, onEdit, onDelete }: IncidentCardListProps) {
  return (
    <ul className="flex flex-col gap-3" role="list">
      {incidents.map((inc) => (
        <li
          key={inc.id}
          className="rounded-md border border-border bg-bg p-4 shadow-sm"
        >
          <div className="flex items-start justify-between gap-3">
            <h3 className="text-sm font-semibold text-ink leading-snug">{inc.titulo}</h3>
            <div className="flex shrink-0 items-center gap-1">
              <button
                type="button"
                aria-label={`Editar ${inc.titulo}`}
                onClick={() => onEdit(inc)}
                className="inline-flex h-9 w-9 items-center justify-center rounded-sm text-ink-secondary hover:bg-bg-subtle hover:text-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
              >
                <Pencil1Icon className="h-4 w-4" />
              </button>
              <button
                type="button"
                aria-label={`Excluir ${inc.titulo}`}
                onClick={() => onDelete(inc)}
                className="inline-flex h-9 w-9 items-center justify-center rounded-sm text-ink-secondary hover:bg-danger-subtle hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
              >
                <TrashIcon className="h-4 w-4" />
              </button>
            </div>
          </div>

          <div className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-2">
            <CategoriaChip value={inc.categoria as Categoria} />
            <SeverityBadge value={inc.severidade as Severidade} />
            <StatusBadge value={inc.status as Status} />
          </div>

          <p className="mt-3 line-clamp-2 text-sm text-ink-secondary">{inc.descricao}</p>
          <p className="mt-2 text-xs text-ink-muted">{relativeTime(inc.updated_at)}</p>
        </li>
      ))}
    </ul>
  )
}
