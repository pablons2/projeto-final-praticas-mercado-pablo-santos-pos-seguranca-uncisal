import { MagnifyingGlassIcon, Cross2Icon } from '@radix-ui/react-icons'
import { Select, type SelectOption } from '@/components/ui'
import {
  CATEGORIAS,
  CATEGORIA_LABELS,
  SEVERIDADES,
  SEVERIDADE_LABELS,
  STATUS,
  STATUS_LABELS,
  type Categoria,
  type Severidade,
  type Status,
} from '@/lib/constants'
import type { IncidentFilters } from '../schemas/incident.schema'

const CATEGORIA_OPTIONS: SelectOption[] = [
  { value: '', label: 'Todas' },
  ...CATEGORIAS.map((c) => ({ value: c, label: CATEGORIA_LABELS[c] })),
]
const SEVERIDADE_OPTIONS: SelectOption[] = [
  { value: '', label: 'Todas' },
  ...SEVERIDADES.map((s) => ({ value: s, label: SEVERIDADE_LABELS[s] })),
]
const STATUS_OPTIONS: SelectOption[] = [
  { value: '', label: 'Todos' },
  ...STATUS.map((s) => ({ value: s, label: STATUS_LABELS[s] })),
]

interface FilterBarProps {
  filters: IncidentFilters
  onChange: (next: IncidentFilters) => void
  total: number
  shown: number
}

export function FilterBar({ filters, onChange, total, shown }: FilterBarProps) {
  const hasActiveFilters =
    filters.search !== '' ||
    filters.categoria !== '' ||
    filters.severidade !== '' ||
    filters.status !== ''

  return (
    <div className="flex flex-col gap-3">
      {/* Search */}
      <div className="relative">
        <MagnifyingGlassIcon
          className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-muted"
          aria-hidden="true"
        />
        <label htmlFor="incident-search" className="sr-only">
          Buscar incidentes
        </label>
        <input
          id="incident-search"
          type="search"
          inputMode="search"
          value={filters.search}
          onChange={(e) => onChange({ ...filters, search: e.target.value })}
          placeholder="Buscar por título ou descrição…"
          className="block h-11 w-full rounded-sm border border-border-strong bg-bg pl-10 pr-3 text-sm text-ink placeholder:text-ink-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring focus-visible:ring-offset-1"
        />
      </div>

      {/* Filters — wrap on mobile, row on >=640px */}
      <div className="flex flex-wrap items-center gap-2">
        <div className="min-w-[140px] flex-1">
          <Select
            ariaLabel="Filtrar por categoria"
            value={filters.categoria ?? ''}
            onChange={(v) => onChange({ ...filters, categoria: v as Categoria | '' })}
            options={CATEGORIA_OPTIONS}
          />
        </div>
        <div className="min-w-[140px] flex-1">
          <Select
            ariaLabel="Filtrar por severidade"
            value={filters.severidade ?? ''}
            onChange={(v) => onChange({ ...filters, severidade: v as Severidade | '' })}
            options={SEVERIDADE_OPTIONS}
          />
        </div>
        <div className="min-w-[140px] flex-1">
          <Select
            ariaLabel="Filtrar por status"
            value={filters.status ?? ''}
            onChange={(v) => onChange({ ...filters, status: v as Status | '' })}
            options={STATUS_OPTIONS}
          />
        </div>
        {hasActiveFilters ? (
          <button
            type="button"
            onClick={() =>
              onChange({ search: '', categoria: '', severidade: '', status: '' })
            }
            className="inline-flex h-11 items-center gap-1.5 rounded-sm px-3 text-sm text-ink-secondary hover:bg-bg-subtle hover:text-ink focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
          >
            <Cross2Icon className="h-4 w-4" />
            Limpar
          </button>
        ) : null}
      </div>

      <p className="text-xs text-ink-muted" aria-live="polite">
        {shown === total
          ? `${total} incidente${total === 1 ? '' : 's'}`
          : `${shown} de ${total}`}
      </p>
    </div>
  )
}
