import { Pencil1Icon, TrashIcon } from '@radix-ui/react-icons'
import { Tooltip } from '@/components/ui'
import { SeverityBadge } from './SeverityBadge'
import { StatusBadge } from './StatusBadge'
import { CategoriaChip } from './CategoriaChip'
import { relativeTime } from '@/lib/format'
import type { Incident } from '@/lib/api-types'
import type { Categoria, Severidade, Status } from '@/lib/constants'

interface IncidentTableProps {
  incidents: Incident[]
  onEdit: (incident: Incident) => void
  onDelete: (incident: Incident) => void
}

export function IncidentTable({ incidents, onEdit, onDelete }: IncidentTableProps) {
  return (
    <div className="overflow-x-auto rounded-md border border-border">
      <table className="w-full border-collapse text-sm">
        <thead>
          <tr className="border-b border-border bg-bg-subtle text-left">
            <th scope="col" className="px-4 py-2.5 font-medium text-ink-secondary">
              Título
            </th>
            <th scope="col" className="px-4 py-2.5 font-medium text-ink-secondary">
              Categoria
            </th>
            <th scope="col" className="px-4 py-2.5 font-medium text-ink-secondary">
              Severidade
            </th>
            <th scope="col" className="px-4 py-2.5 font-medium text-ink-secondary">
              Status
            </th>
            <th scope="col" className="px-4 py-2.5 font-medium text-ink-secondary">
              Atualizado
            </th>
            <th scope="col" className="px-4 py-2.5 text-right font-medium text-ink-secondary">
              Ações
            </th>
          </tr>
        </thead>
        <tbody>
          {incidents.map((inc) => (
            <tr
              key={inc.id}
              className="border-b border-border-hairline last:border-0 hover:bg-bg-subtle/60"
            >
              <td className="max-w-[280px] truncate px-4 py-3 font-medium text-ink">
                {inc.titulo}
              </td>
              <td className="px-4 py-3">
                <CategoriaChip value={inc.categoria as Categoria} />
              </td>
              <td className="px-4 py-3">
                <SeverityBadge value={inc.severidade as Severidade} />
              </td>
              <td className="px-4 py-3">
                <StatusBadge value={inc.status as Status} />
              </td>
              <td className="whitespace-nowrap px-4 py-3 text-ink-secondary">
                {relativeTime(inc.updated_at)}
              </td>
              <td className="px-4 py-3">
                <div className="flex items-center justify-end gap-1">
                  <Tooltip content="Editar">
                    <button
                      type="button"
                      aria-label={`Editar ${inc.titulo}`}
                      onClick={() => onEdit(inc)}
                      className="inline-flex h-9 w-9 items-center justify-center rounded-sm text-ink-secondary hover:bg-bg-subtle hover:text-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
                    >
                      <Pencil1Icon className="h-4 w-4" />
                    </button>
                  </Tooltip>
                  <Tooltip content="Excluir">
                    <button
                      type="button"
                      aria-label={`Excluir ${inc.titulo}`}
                      onClick={() => onDelete(inc)}
                      className="inline-flex h-9 w-9 items-center justify-center rounded-sm text-ink-secondary hover:bg-danger-subtle hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
                    >
                      <TrashIcon className="h-4 w-4" />
                    </button>
                  </Tooltip>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
