import { useMemo, useState } from 'react'
import { PlusIcon, ExclamationTriangleIcon } from '@radix-ui/react-icons'
import { TopBar } from '@/components/TopBar'
import { Button, TooltipProvider } from '@/components/ui'
import { useIncidents } from '@/features/incidents/hooks/useIncidents'
import { FilterBar } from '@/features/incidents/components/FilterBar'
import { IncidentTable } from '@/features/incidents/components/IncidentTable'
import { IncidentCardList } from '@/features/incidents/components/IncidentCardList'
import {
  IncidentTableSkeleton,
  IncidentCardListSkeleton,
} from '@/features/incidents/components/IncidentListSkeleton'
import { EmptyState } from '@/features/incidents/components/EmptyState'
import { IncidentFormDialog } from '@/features/incidents/components/IncidentFormDialog'
import { DeleteConfirmDialog } from '@/features/incidents/components/DeleteConfirmDialog'
import { useIsDesktop } from '@/lib/useMediaQuery'
import type { Incident } from '@/lib/api-types'
import type { IncidentFilters } from '@/features/incidents/schemas/incident.schema'

const EMPTY_FILTERS: IncidentFilters = {
  search: '',
  categoria: '',
  severidade: '',
  status: '',
}

export function DashboardPage() {
  const isDesktop = useIsDesktop()
  const { data: incidents, isLoading, isError, refetch, isFetching } = useIncidents()

  const [filters, setFilters] = useState<IncidentFilters>(EMPTY_FILTERS)
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Incident | undefined>(undefined)
  const [deleting, setDeleting] = useState<Incident | null>(null)

  const filtered = useMemo(() => {
    const list = incidents ?? []
    const q = filters.search.trim().toLowerCase()
    return list.filter((inc) => {
      if (q) {
        const hay = `${inc.titulo} ${inc.descricao}`.toLowerCase()
        if (!hay.includes(q)) return false
      }
      if (filters.categoria && inc.categoria !== filters.categoria) return false
      if (filters.severidade && inc.severidade !== filters.severidade) return false
      if (filters.status && inc.status !== filters.status) return false
      return true
    })
  }, [incidents, filters])

  const total = incidents?.length ?? 0

  const openCreate = () => {
    setEditing(undefined)
    setFormOpen(true)
  }
  const openEdit = (inc: Incident) => {
    setEditing(inc)
    setFormOpen(true)
  }

  return (
    <TooltipProvider delayDuration={200}>
      <div className="min-h-dvh bg-bg-surface">
        <TopBar />

        {/* Skip link for keyboard users */}
        <a
          href="#main"
          className="sr-only focus:not-sr-only focus:fixed focus:left-3 focus:top-3 focus:z-[70] focus:rounded-sm focus:bg-ink focus:px-3 focus:py-2 focus:text-sm focus:text-white"
        >
          Pular para o conteúdo
        </a>

        <main id="main" className="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 sm:py-8">
          {/* Page header */}
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h1 className="text-3xl font-semibold tracking-tight text-ink">Incidentes</h1>
              <p className="mt-1 text-sm text-ink-secondary">
                {total === 0
                  ? 'Nenhum incidente ainda.'
                  : `${total} registrado${total === 1 ? '' : 's'}`}
              </p>
            </div>
            <Button onClick={openCreate} icon={PlusIcon} className="w-full sm:w-auto">
              Novo incidente
            </Button>
          </div>

          {/* Filters */}
          <div className="mt-6">
            <FilterBar
              filters={filters}
              onChange={setFilters}
              total={total}
              shown={filtered.length}
            />
          </div>

          {/* Content */}
          <div className="mt-6">
            {isError ? (
              <div className="flex flex-col items-center gap-3 rounded-md border border-border bg-bg px-6 py-12 text-center">
                <ExclamationTriangleIcon className="h-8 w-8 text-danger" />
                <div>
                  <p className="text-sm font-medium text-ink">Falha ao carregar</p>
                  <p className="mt-1 text-sm text-ink-secondary">
                    Não foi possível buscar seus incidentes.
                  </p>
                </div>
                <Button variant="secondary" onClick={() => refetch()}>
                  Tentar novamente
                </Button>
              </div>
            ) : isLoading ? (
              isDesktop ? (
                <IncidentTableSkeleton />
              ) : (
                <IncidentCardListSkeleton />
              )
            ) : total === 0 ? (
              <div className="rounded-md border border-border bg-bg">
                <EmptyState
                  title="Nenhum incidente registrado"
                  description="Comece a registrar incidentes de segurança para acompanhá-los aqui."
                  actionLabel="Registrar primeiro incidente"
                  onAction={openCreate}
                />
              </div>
            ) : filtered.length === 0 ? (
              <div className="rounded-md border border-border bg-bg">
                <EmptyState
                  title="Nada encontrado"
                  description="Ajuste os filtros ou limpe a busca para ver seus incidentes."
                />
              </div>
            ) : isDesktop ? (
              <div className={isFetching ? 'opacity-60 transition-opacity' : undefined}>
                <IncidentTable incidents={filtered} onEdit={openEdit} onDelete={setDeleting} />
              </div>
            ) : (
              <div className={isFetching ? 'opacity-60 transition-opacity' : undefined}>
                <IncidentCardList
                  incidents={filtered}
                  onEdit={openEdit}
                  onDelete={setDeleting}
                />
              </div>
            )}
          </div>
        </main>

        {/* Dialogs */}
        <IncidentFormDialog
          open={formOpen}
          onOpenChange={setFormOpen}
          incident={editing}
        />
        <DeleteConfirmDialog
          open={deleting !== null}
          onOpenChange={(o) => !o && setDeleting(null)}
          incident={deleting}
        />
      </div>
    </TooltipProvider>
  )
}
